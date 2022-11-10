#include "yapi_inc.h"
#include <stdarg.h>
#include <stdio.h>
#include <string.h>

#ifdef _WIN32
#include <windows.h>
#else
#include <errno.h>
#include <pthread.h>
#include <sys/time.h>
#include <dlfcn.h>
#endif
#ifdef __linux
#include <unistd.h>
#include <sys/syscall.h>
#endif

// macro to simplify code for loading each symbol
#define YAPI_LOAD_SYMBOL(symbolName, symbol)                                       \
    if (!symbol && yapiLoadSymbol(symbolName, (YapiPointer*)&symbol, error) < 0) { \
        return YAPI_ERROR;                                                         \
    }

#define YAPI_CHECK_CLI_RETURN() \
    if (ret != YAPI_SUCCESS) {   \
        yapiGetCliError(error); \
    }                           \
    return ret;

static YapiSymbols yapiSymbols = {NULL};
static void*       yapiLibHandle = NULL;

#ifdef _WIN32

static YapiResult yapiGetWindowsError(DWORD errNum, YapiErrorMsg* error)
{
    char*    fallbackErrorFormat = "failed to get message for Windows Error %d";
    wchar_t* errBuf = NULL;
    DWORD    length = 0;

    DWORD status =
        FormatMessageW(FORMAT_MESSAGE_FROM_SYSTEM | FORMAT_MESSAGE_IGNORE_INSERTS | FORMAT_MESSAGE_ALLOCATE_BUFFER,
                       NULL, errNum, MAKELANGID(LANG_ENGLISH, SUBLANG_ENGLISH_US), (LPWSTR)&errBuf, 0, NULL);
    if (!status && GetLastError() == ERROR_MUI_FILE_NOT_FOUND)
        FormatMessageW(FORMAT_MESSAGE_FROM_SYSTEM | FORMAT_MESSAGE_IGNORE_INSERTS | FORMAT_MESSAGE_ALLOCATE_BUFFER,
                       NULL, errNum, MAKELANGID(LANG_NEUTRAL, SUBLANG_DEFAULT), (LPWSTR)&errBuf, 0, NULL);

    if (errBuf == NULL) {
        return YAPI_SUCCESS;
    }

    // strip trailing period and carriage return from message, if needed
    length = (DWORD)wcslen(errBuf);
    errBuf[length] = L'\0';

    // convert to UTF-8 encoding
    if (length > 0) {
        length = WideCharToMultiByte(CP_UTF8, 0, errBuf, -1, error->buf->message, T2S_BUFFER_SIZE, NULL, NULL);
    }
    LocalFree(errBuf);
    return YAPI_SUCCESS;
}

YapiResult yapiOpenDynamicLib(char* libName, YapiPointer* handler, YapiErrorMsg* error)
{
    *handler = LoadLibraryA(libName);
    if (*handler != NULL) {
        yapiLibHandle = *handler;
        return YAPI_SUCCESS;
    }

    DWORD errNum = GetLastError();
    return yapiGetWindowsError(errNum, error);
}

YapiResult yapiCloseDynamicLib(YapiPointer* handler, YapiErrorMsg* error)
{
    BOOL ret = FreeLibrary(*handler);
    if (ret) {
        *handler = NULL;
        return YAPI_SUCCESS;
    }
    DWORD errNum = GetLastError();
    yapiSetError(error, errNum, "");
    return YAPI_ERROR;
}

static int yapiLoadSymbol(const char* symbolName, void** symbol, YapiErrorMsg* error)
{
    *symbol = GetProcAddress(yapiLibHandle, symbolName);
    if (*symbol != NULL) {
        return YAPI_SUCCESS;
    }
    yapiSetError(error, YAPI_ERR_LOAD_SYMBOL, "symbol %s not found in yacli library", symbolName);
    return YAPI_ERROR;
}

#else

YapiResult yapiOpenDynamicLib(char* libName, YapiPointer* handler, YapiErrorMsg* error)
{
    *handler = dlopen(libName, RTLD_LAZY);
    if (!*handler) {
        char* errMsg = dlerror();
        yapiSetError(error, YAPI_ERR_LOAD_SYMBOL, "load yacli library error [%s]", errMsg);
        return YAPI_ERROR;
    }

    yapiLibHandle = *handler;
    return YAPI_SUCCESS;
}

YapiResult yapiCloseDynamicLib(YapiPointer* handler, YapiErrorMsg* error)
{
    int32_t ret = dlclose(*handler);
    if (ret != YAPI_SUCCESS) {
        char* errMsg = dlerror();
        yapiSetError(error, errno, errMsg);
        return YAPI_ERROR;
    }
    *handler = NULL;
    return YAPI_SUCCESS;
}

static int yapiLoadSymbol(const char* symbolName, void** symbol, YapiErrorMsg* error)
{
    *symbol = dlsym(yapiLibHandle, symbolName);
    if (!*symbol) {
        yapiSetError(error, YAPI_ERR_LOAD_SYMBOL, "symbol %s not found in yacli library", symbolName);
        return YAPI_ERROR;
    }
    return YAPI_SUCCESS;
}
#endif

YapiResult yapiCliAllocHandle(YapiHandleType type, YacHandle input, YacHandle* output, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacAllocHandle", yapiSymbols.fnAllocHandle)
    ret = (*yapiSymbols.fnAllocHandle)(type, input, output);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliFreeHandle(YapiHandleType type, YacHandle handle, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacFreeHandle", yapiSymbols.fnHandleFree)
    ret = (*yapiSymbols.fnHandleFree)(type, handle);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliGetVersion(char** version, YapiErrorMsg* error)
{
    YAPI_LOAD_SYMBOL("yacGetGetVersion", yapiSymbols.fnGetVersion)
    *version = (*yapiSymbols.fnGetVersion)();
    return YAPI_SUCCESS;
}

YapiResult yapiCliGetLastError(YapiErrorMsg* error)
{
    char *msg;
    char *stat;
    YAPI_LOAD_SYMBOL("yacGetLastError", yapiSymbols.fnGetLastError)
    (*yapiSymbols.fnGetLastError)(&error->buf->code, &msg, &stat, &error->buf->pos);
    strcpy(error->buf->message, msg);
    strcpy(error->buf->sqlState, stat);
    return YAPI_SUCCESS;
}

YapiResult yapiCliGetEnvAttr(YacHandle hEnv, YapiEnvAttr attr, void* value, int32_t bufLength, int32_t* stringLength,
                             YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacGetEnvAttr", yapiSymbols.fnGetEnvAttr)
    ret = (*yapiSymbols.fnGetEnvAttr)(hEnv, attr, value, bufLength, stringLength);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliConnect(YacHandle hConn, const char* url, int16_t urlLength, const char* user, int16_t userLength,
                          const char* password, int16_t passwordLength, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacConnect", yapiSymbols.fnConnect)
    ret = (*yapiSymbols.fnConnect)(hConn, url, urlLength, user, userLength, password, passwordLength);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliDisconnect(YacHandle hConn, YapiErrorMsg* error)
{
    YAPI_LOAD_SYMBOL("yacDisconnect", yapiSymbols.fnDisconnect)
    (*yapiSymbols.fnDisconnect)(hConn);
    return YAPI_SUCCESS;
}

YapiResult yapiCliSetConnAttr(YacHandle hConn, YapiConnAttr attr, void* value, int32_t length, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacSetConnAttr", yapiSymbols.fnSetConnAttr)
    ret = (*yapiSymbols.fnSetConnAttr)(hConn, attr, value, length);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliGetConnAttr(YacHandle hConn, YapiConnAttr attr, void* value, int32_t bufLength, int32_t* stringLength,
                              YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacGetConnAttr", yapiSymbols.fnGetConnAttr)
    ret = (*yapiSymbols.fnGetConnAttr)(hConn, attr, value, bufLength, stringLength);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliCommit(YacHandle hConn, YapiErrorMsg* error)
{
       YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacCommit", yapiSymbols.fnCommit)
    ret = (*yapiSymbols.fnCommit)(hConn);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliRollback(YacHandle hConn, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacRollback", yapiSymbols.fnRollback)
    ret = (*yapiSymbols.fnRollback)(hConn);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliCancel(YacHandle hConn, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacCancel", yapiSymbols.fnCancel)
    ret = (*yapiSymbols.fnCancel)(hConn);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliDirectExecute(YacHandle hStmt, const char* sql, int32_t sqlLength, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacDirectExecute", yapiSymbols.fnDirectExecute)
    ret = (*yapiSymbols.fnDirectExecute)(hStmt, sql, sqlLength);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliPrepare(YacHandle hStmt, const char* sql, int32_t sqlLength, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacPrepare", yapiSymbols.fnPrepare)
    ret = (*yapiSymbols.fnPrepare)(hStmt, sql, sqlLength);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliExecute(YacHandle hStmt, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacExecute", yapiSymbols.fnExecute)
    ret = (*yapiSymbols.fnExecute)(hStmt);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliSetStmtAttr(YacHandle hStmt, YapiStmtAttr attr, void* value, int32_t length, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacSetStmtAttr", yapiSymbols.fnSetStmtAttr)
    ret = (*yapiSymbols.fnSetStmtAttr)(hStmt, attr, value, length);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliGetStmtAttr(YacHandle hStmt, YapiStmtAttr attr, void* value, int32_t bufLength, int32_t* stringLength, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacGetStmtAttr", yapiSymbols.fnGetStmtAttr)
    ret = (*yapiSymbols.fnGetStmtAttr)(hStmt, attr, value, bufLength, stringLength);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliFetch(YacHandle hStmt, uint32_t* rows, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacFetch", yapiSymbols.fnFetch)
    ret = (*yapiSymbols.fnFetch)(hStmt, rows);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliDescribeCol2(YacHandle hStmt, uint16_t id, YapiColumnDesc* desc, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacDescribeCol2", yapiSymbols.fnDescribeCol2)
    ret = (*yapiSymbols.fnDescribeCol2)(hStmt, id, desc);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliBindColumn(YacHandle hStmt, uint16_t id, YapiType type, YapiPointer value, int32_t bufLen,
                             int32_t* indicator, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacBindColumn", yapiSymbols.fnBindColumn)
    ret = (*yapiSymbols.fnBindColumn)(hStmt, id, type, value, bufLen, indicator);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliBindParameter(YacHandle hStmt, uint16_t id, YapiParamDirection direction, YapiType bindType,
                                YapiPointer value, int32_t bindSize, int32_t bufLength, int32_t* indicator, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacBindParameter", yapiSymbols.fnBindParameter)
    ret = (*yapiSymbols.fnBindParameter)(hStmt, id, direction, bindType, value, bindSize, bufLength, indicator);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliBindParameterByName(YacHandle hStmt, char* name, YapiParamDirection direction, YapiType bindType,
                                      YapiPointer value, int32_t bindSize, int32_t bufLength, int32_t* indicator, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacBindParameterByName", yapiSymbols.fnBindParameterByName)
    ret = (*yapiSymbols.fnBindParameterByName)(hStmt, name, direction, bindType, value, bindSize, bufLength, indicator);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliNumResultCols(YacHandle hStmt, int16_t* count, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacNumResultCols", yapiSymbols.fnNumResultCols)
    ret = (*yapiSymbols.fnNumResultCols)(hStmt, count);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliColAttribute(YacHandle hStmt, uint16_t id, YapiColAttr attr, void* value, int32_t bufLen,
                              int32_t* stringLength, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacColAttribute", yapiSymbols.fnColAttribute)
    ret = (*yapiSymbols.fnColAttribute)(hStmt, id, attr, value, bufLen, stringLength);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliNumParams(YacHandle hStmt, int16_t* count, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacNumParams", yapiSymbols.fnNumParams)
    ret = (*yapiSymbols.fnNumParams)(hStmt, count);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliGetDateStruct(YapiDate date, YapiDateStruct* ds, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacGetDateStruct", yapiSymbols.fnGetDateStruct)
    ret = (*yapiSymbols.fnGetDateStruct)(date, ds);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliLobDescAlloc(YacHandle* hConn, YapiType type, void** desc, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacLobDescAlloc", yapiSymbols.fnLobDescAlloc)
    ret = (*yapiSymbols.fnLobDescAlloc)(hConn, type, desc);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliLobDescFree(void* desc, YapiType type, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacLobDescFree", yapiSymbols.fnLobDescFree)
    ret = (*yapiSymbols.fnLobDescFree)(desc, type);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliLobGetChunkSize(YacHandle* hConn, YapiLobLocator* locator, uint16_t* chunkSize, YapiErrorMsg* error)
{
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacLobGetChunkSize", yapiSymbols.fnLobGetChunkSize)
    ret = (*yapiSymbols.fnLobGetChunkSize)(hConn, locator, chunkSize);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliLobGetLength(YacHandle* hConn, YapiLobLocator* locator, uint64_t* length, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacLobGetLength", yapiSymbols.fnLobGetLength)
    ret = (*yapiSymbols.fnLobGetLength)(hConn, locator, length);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliLobRead(YacHandle* hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf, uint64_t bufLen,
                          YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacLobRead", yapiSymbols.fnLobRead)
    ret = (*yapiSymbols.fnLobRead)(hConn, loc, bytes, buf, bufLen);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliLobWrite(YacHandle* hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf, uint64_t bufLen,
                           YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacLobWrite", yapiSymbols.fnLobWrite)
    ret = (*yapiSymbols.fnLobWrite)(hConn, loc, bytes, buf, bufLen);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliLobCreateTemporary(YacHandle* hConn, YapiLobLocator* loc, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacLobCreateTemporary", yapiSymbols.fnLobCreateTemporary)
    ret = (*yapiSymbols.fnLobCreateTemporary)(hConn, loc);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliLobFreeTemporary(YacHandle* hConn, YapiLobLocator* loc, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacLobFreeTemporary", yapiSymbols.fnLobFreeTemporary)
    ret = (*yapiSymbols.fnLobFreeTemporary)(hConn, loc);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliDateGetDate(const YapiDate date, int16_t* year, uint8_t* month, uint8_t* day, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacDateGetDate", yapiSymbols.fnDateGetDate)
    ret = (*yapiSymbols.fnDateGetDate)(date, year, month, day);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliShortTimeGetShortTime(const YapiShortTime time, uint8_t* hour, uint8_t* minute, uint8_t* second,
                                       uint32_t* fraction, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacShortTimeGetShortTime", yapiSymbols.fnShortTimeGetShortTime)
    ret = (*yapiSymbols.fnShortTimeGetShortTime)(time, hour, minute, second, fraction);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliTimestampGetTimestamp(const YapiTimestamp timestamp, int16_t* year, uint8_t* month, uint8_t* day,
                                       uint8_t* hour, uint8_t* minute, uint8_t* second, uint32_t* fraction,
                                       YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacTimestampGetTimestamp", yapiSymbols.fnTimestampGetTimestamp)
    ret = (*yapiSymbols.fnTimestampGetTimestamp)(timestamp, year, month, day, hour, minute, second, fraction);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliYMIntervalGetYearMonth(const YapiYMInterval ymInterval, int32_t* year, int32_t* month,
                                        YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacYMIntervalGetYearMonth", yapiSymbols.fnYMIntervalGetYearMonth)
    ret = (*yapiSymbols.fnYMIntervalGetYearMonth)(ymInterval, year, month);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliDSIntervalGetDaySecond(const YapiDSInterval dsInterval, int32_t* day, int32_t* hour, int32_t* minute,
                                        int32_t* second, int32_t* fraction, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacDSIntervalGetDaySecond", yapiSymbols.fnDSIntervalGetDaySecond)
    ret = (*yapiSymbols.fnDSIntervalGetDaySecond)(dsInterval, day, hour, minute, second, fraction);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliDateSetDate(YapiDate* date, int16_t year, uint8_t month, uint8_t day, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacDateSetDate", yapiSymbols.fnDateSetDate)
    ret = (*yapiSymbols.fnDateSetDate)(date, year, month, day);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliShortTimeSetShortTime(YapiShortTime* time, uint8_t hour, uint8_t minute, uint8_t second,
                                       uint32_t fraction, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacShortTimeSetShortTime", yapiSymbols.fnShortTimeSetShortTime)
    ret = (*yapiSymbols.fnShortTimeSetShortTime)(time, hour, minute, second, fraction);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliTimestampSetTimestamp(YapiTimestamp* timestamp, int16_t year, uint8_t month, uint8_t day,
                                        uint8_t hour,
                                       uint8_t minute, uint8_t second, uint32_t fraction, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacTimestampSetTimestamp", yapiSymbols.fnTimestampSetTimestamp)
    ret = (*yapiSymbols.fnTimestampSetTimestamp)(timestamp, year, month, day, hour, minute, second, fraction);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliYMIntervalSetYearMonth(YapiYMInterval* ymInterval, int32_t year, int32_t month, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacYMIntervalSetYearMonth", yapiSymbols.fnYMIntervalSetYearMonth)
    ret = (*yapiSymbols.fnYMIntervalSetYearMonth)(ymInterval, year, month);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliDSIntervalSetDaySecond(YapiDSInterval* dsInterval, int32_t day, int32_t hour, int32_t minute,
                                        int32_t second, int32_t fraction, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacDSIntervalSetDaySecond", yapiSymbols.fnDSIntervalSetDaySecond)
    ret = (*yapiSymbols.fnDSIntervalSetDaySecond)(dsInterval, day, hour, minute, second, fraction);
    YAPI_CHECK_CLI_RETURN();
}

YapiResult yapiCliNumberRound(YapiNumber* n, int32_t precision, int32_t scale, YapiErrorMsg* error)
{
    YapiResult ret;

    YAPI_LOAD_SYMBOL("yacNumberRound", yapiSymbols.fnNumberRound)
    ret = (*yapiSymbols.fnNumberRound)(n, precision, scale);
    YAPI_CHECK_CLI_RETURN();
}