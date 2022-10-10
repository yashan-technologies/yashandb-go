#include "yapi_inc.h"
#include <stdlib.h>
#include <stdarg.h>
#include <stdio.h>

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

static YapiSymbols yapiSymbols = {NULL};
static void*       yapiLibHandle = NULL;

#ifdef _WIN32

static YapiResult yapiGetWindowsError(DWORD errNum, YapiErrorMsg* error)
{
    TCHAR* errBUf = NULL;
    DWORD  status =
        FormatMessage(FORMAT_MESSAGE_FROM_SYSTEM | FORMAT_MESSAGE_IGNORE_INSERTS | FORMAT_MESSAGE_ALLOCATE_BUFFER, NULL,
                      errNum, MAKELANGID(LANG_ENGLISH, SUBLANG_ENGLISH_US), (LPTSTR)&errBUf, 0, NULL);
    return YAPI_ERROR;
}

YapiResult yapiOpenDynamicLib(char* libName, YapiPointer* handler, YapiErrorMsg* error)
{
    *handler = LoadLibrary(libName);
    if (*handler != NULL) {
        yapiLibHandle = handler;
        return YAPI_SUCCESS;
    }

    DWORD errNum = GetLastError();
    // otherwise, attempt to get the error message
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

YapiResult yapiCliAllocHandle(YapiHandleType type, YacHandle input, YacHandle* output)
{
    YapiErrorMsg* error = NULL;
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacAllocHandle", yapiSymbols.fnAllocHandle)
    ret = (*yapiSymbols.fnAllocHandle)(type, input, output);
    return ret;
}

YapiResult yapiCliFreeHandle(YapiHandleType type, YacHandle handle)
{
    YapiErrorMsg* error = NULL;
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacFreeHandle", yapiSymbols.fnHandleFree)
    ret = (*yapiSymbols.fnHandleFree)(type, handle);
    return ret;
}

YapiResult yapiCliGetVersion(char** version)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacGetGetVersion", yapiSymbols.fnGetVersion)
    *version = (*yapiSymbols.fnGetVersion)();
    return YAPI_SUCCESS;
}

YapiResult yapiCliGetLastError(int32_t* errCode, char** message, char** sqlState, YapiTextPos* pos)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacGetLastError", yapiSymbols.fnGetLastError)
    (*yapiSymbols.fnGetLastError)(errCode, message, sqlState, pos);
    return YAPI_SUCCESS;
}

YapiResult yapiCliGetEnvAttr(YacHandle hEnv, YapiEnvAttr attr, void* value, int32_t bufLength, int32_t* stringLength)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacGetEnvAttr", yapiSymbols.fnGetEnvAttr)
    return (*yapiSymbols.fnGetEnvAttr)(hEnv, attr, value, bufLength, stringLength);
}

YapiResult yapiCliConnect(YacHandle hConn, const char* url, int16_t urlLength, const char* user, int16_t userLength,
                          const char* password, int16_t passwordLength)
{
    YapiErrorMsg* error = NULL;
    YapiResult    ret;

    YAPI_LOAD_SYMBOL("yacConnect", yapiSymbols.fnConnect)
    ret = (*yapiSymbols.fnConnect)(hConn, url, urlLength, user, userLength, password, passwordLength);
    return ret;
}

YapiResult yapiCliDisconnect(YacHandle hConn)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacDisconnect", yapiSymbols.fnDisconnect)
    (*yapiSymbols.fnDisconnect)(hConn);
    return YAPI_SUCCESS;
}

YapiResult yapiCliSetConnAttr(YacHandle hConn, YapiConnAttr attr, void* value, int32_t length)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacSetConnAttr", yapiSymbols.fnSetConnAttr)
    return (*yapiSymbols.fnSetConnAttr)(hConn, attr, value, length);
}

YapiResult yapiCliGetConnAttr(YacHandle hConn, YapiConnAttr attr, void* value, int32_t bufLength, int32_t* stringLength)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacGetConnAttr", yapiSymbols.fnGetConnAttr)
    return (*yapiSymbols.fnGetConnAttr)(hConn, attr, value, bufLength, stringLength);
}

YapiResult yapiCliCommit(YacHandle hConn)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacCommit", yapiSymbols.fnCommit)
    return (*yapiSymbols.fnCommit)(hConn);
}

YapiResult yapiCliRollback(YacHandle hConn)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacRollback", yapiSymbols.fnRollback)
    return (*yapiSymbols.fnRollback)(hConn);
}

YapiResult yapiCliCancel(YacHandle hConn)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacCancel", yapiSymbols.fnCancel)
    return (*yapiSymbols.fnCancel)(hConn);
}

YapiResult yapiCliDirectExecute(YacHandle hStmt, const char* sql, int32_t sqlLength)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacDirectExecute", yapiSymbols.fnDirectExecute)
    return (*yapiSymbols.fnDirectExecute)(hStmt, sql, sqlLength);
}

YapiResult yapiCliPrepare(YacHandle hStmt, const char* sql, int32_t sqlLength)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacPrepare", yapiSymbols.fnPrepare)
    return (*yapiSymbols.fnPrepare)(hStmt, sql, sqlLength);
}

YapiResult yapiCliExecute(YacHandle hStmt)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacExecute", yapiSymbols.fnExecute)
    return (*yapiSymbols.fnExecute)(hStmt);
}

YapiResult yapiCliSetStmtAttr(YacHandle hStmt, YapiStmtAttr attr, void* value, int32_t length)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacSetStmtAttr", yapiSymbols.fnSetStmtAttr)
    return (*yapiSymbols.fnSetStmtAttr)(hStmt, attr, value, length);
}

YapiResult yapiCliGetStmtAttr(YacHandle hStmt, YapiStmtAttr attr, void* value, int32_t bufLength, int32_t* stringLength)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacGetStmtAttr", yapiSymbols.fnGetStmtAttr)
    return (*yapiSymbols.fnGetStmtAttr)(hStmt, attr, value, bufLength, stringLength);
}

YapiResult yapiCliFetch(YacHandle hStmt, uint32_t* rows)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacFetch", yapiSymbols.fnFetch)
    return (*yapiSymbols.fnFetch)(hStmt, rows);
}

YapiResult yapiCliDescribeCol2(YacHandle hStmt, uint16_t id, YapiColumnDesc* desc)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacDescribeCol2", yapiSymbols.fnDescribeCol2)
    return (*yapiSymbols.fnDescribeCol2)(hStmt, id, desc);
}

YapiResult yapiCliBindColumn(YacHandle hStmt, uint16_t id, YapiType type, YapiPointer value, int32_t bufLen,
                             int32_t* indicator)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacBindColumn", yapiSymbols.fnBindColumn)
    return (*yapiSymbols.fnBindColumn)(hStmt, id, type, value, bufLen, indicator);
}

YapiResult yapiCliBindParameter(YacHandle hStmt, uint16_t id, YapiParamDirection direction, YapiType bindType,
                                YapiPointer value, uint32_t bindSize, int32_t bufLength, int32_t* indicator)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacBindParameter", yapiSymbols.fnBindParameter)
    return (*yapiSymbols.fnBindParameter)(hStmt, id, direction, bindType, value, bindSize, bufLength, indicator);
}

YapiResult yapiCliBindParameterByName(YacHandle hStmt, char* name, YapiParamDirection direction, YapiType bindType,
                                      YapiPointer value, uint32_t bindSize, int32_t bufLength, int32_t* indicator)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacBindParameterByName", yapiSymbols.fnBindParameterByName)
    return (*yapiSymbols.fnBindParameterByName)(hStmt, name, direction, bindType, value, bindSize, bufLength,
                                                indicator);
}

YapiResult yapiCliNumResultCols(YacHandle hStmt, int16_t* count)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacNumResultCols", yapiSymbols.fnNumResultCols)
    return (*yapiSymbols.fnNumResultCols)(hStmt, count);
}

YapiResult yapiCliGetDateStruct(YapiDate date, YapiDateStruct* ds)
{
    YapiErrorMsg* error = NULL;

    YAPI_LOAD_SYMBOL("yacGetDateStruct", yapiSymbols.fnGetDateStruct)
    return (*yapiSymbols.fnGetDateStruct)(date, ds);
}

YapiResult yapiCliLobDescAlloc(YapiConnect* hConn, YapiType type, void** desc)
{
    YapiErrorMsg* error = NULL;
    YAPI_LOAD_SYMBOL("yacLobDescAlloc", yapiSymbols.fnLobDescAlloc)
    return (*yapiSymbols.fnLobDescAlloc)(hConn, type, desc);
}
YapiResult yapiCliLobDescFree(void* desc, YapiType type)
{
    YapiErrorMsg* error = NULL;
    YAPI_LOAD_SYMBOL("yacLobDescFree", yapiSymbols.fnLobDescFree)
    return (*yapiSymbols.fnLobDescFree)(desc, type);
}
YapiResult yapiCliLobGetChunkSize(YapiConnect* hConn, YapiLobLocator* locator, uint16_t* chunkSize)
{
    YapiErrorMsg* error = NULL;
    YAPI_LOAD_SYMBOL("yacLobGetChunkSize", yapiSymbols.fnLobGetChunkSize)
    return (*yapiSymbols.fnLobGetChunkSize)(hConn, locator, chunkSize);
}
YapiResult yapiCliLobGetLength(YapiConnect* hConn, YapiLobLocator* locator, uint64_t* length)
{
    YapiErrorMsg* error = NULL;
    YAPI_LOAD_SYMBOL("yacLobGetLength", yapiSymbols.fnLobGetLength)
    return (*yapiSymbols.fnLobGetLength)(hConn, locator, length);
}
YapiResult yapiCliLobRead(YapiConnect* hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf, uint64_t bufLen)
{
    YapiErrorMsg* error = NULL;
    YAPI_LOAD_SYMBOL("yacLobRead", yapiSymbols.fnLobRead)
    return (*yapiSymbols.fnLobRead)(hConn, loc, bytes, buf, bufLen);
}
YapiResult yapiCliLobWrite(YapiConnect* hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf, uint64_t bufLen)
{
    YapiErrorMsg* error = NULL;
    YAPI_LOAD_SYMBOL("yacLobWrite", yapiSymbols.fnLobWrite)
    return (*yapiSymbols.fnLobWrite)(hConn, loc, bytes, buf, bufLen);
}
YapiResult yapiCliLobCreateTemporary(YapiConnect* hConn, YapiLobLocator* loc)
{
    YapiErrorMsg* error = NULL;
    YAPI_LOAD_SYMBOL("yacLobCreateTemporaryt", yapiSymbols.fnLobCreateTemporary)
    return (*yapiSymbols.fnLobCreateTemporary)(hConn, loc);
}
YapiResult yapiCliLobFreeTemporary(YapiConnect* hConn, YapiLobLocator* loc)
{
    YapiErrorMsg* error = NULL;
    YAPI_LOAD_SYMBOL("yacLobFreeTemporary", yapiSymbols.fnLobFreeTemporary)
    return (*yapiSymbols.fnLobFreeTemporary)(hConn, loc);
}