#ifndef YAPI_INC_H
#define YAPI_INC_H

#include "yacapi.h"

#ifdef __cplusplus
extern "C" {
#endif

#define YAPI_CALL(yapiFunc)      \
    do {                         \
        YapiResult r = yapiFunc; \
        if (r == YAPI_ERROR) {   \
            return r;            \
        }                        \
    } while (0)

typedef enum EnYacResult { YAC_SUCCESS = 0, YAC_SUCCESS_WITH_INFO = 1, YAC_ERROR = -1 } YacResult;
typedef void* YacHandle;

typedef YacResult (*yapiFuncAllocHandle)(YapiHandleType type, YacHandle input, YacHandle* output);
typedef YacResult (*yapiFuncFreeHandler)(YapiHandleType type, YacHandle handle);
typedef char* (*yapiFuncGetVersion)();
typedef void (*yapiFuncGetLastError)(int32_t* errCode, char** message, char** sqlState, YapiTextPos* pos);
typedef YacResult (*yapiFuncSetEnvAttr)(YacHandle hEnv, YapiEnvAttr attr, void* value, int32_t length);
typedef YacResult (*yapiFuncGetEnvAttr)(YacHandle hEnv, YapiEnvAttr attr, void* value, int32_t bufLength,
                                        int32_t* stringLength);

typedef YacResult (*yapiFuncConnect)(YacHandle hConn, const char* url, int16_t urlLength, const char* user,
                                     int16_t userLength, const char* password, int16_t passwordLength);
typedef void (*yapiFuncDisconnect)(YacHandle hConn);
typedef YacResult (*yapiFuncCancel)(YacHandle hConn);
typedef YacResult (*yapiFuncDirectExecute)(YacHandle hStmt, const char* sql, int32_t sqlLength);
typedef YacResult (*yapiFuncPrepare)(YacHandle hStmt, const char* sql, int32_t sqlLength);
typedef YacResult (*yapiFuncExecute)(YacHandle hStmt);
typedef YacResult (*yapiFuncFetch)(YacHandle hStmt, uint32_t* rows);
typedef YacResult (*yapiFuncCommit)(YacHandle hConn);
typedef YacResult (*yapiFuncRollback)(YacHandle hConn);

typedef YacResult (*yapiFuncSetConnAttr)(YacHandle hConn, YapiConnAttr attr, void* value, int32_t length);
typedef YacResult (*yapiFuncGetConnAttr)(YacHandle hConn, YapiConnAttr attr, void* value, int32_t bufLength,
                                         int32_t* stringLength);
typedef YacResult (*yapiFuncSetStmtAttr)(YacHandle hStmt, YapiStmtAttr attr, void* value, int32_t length);
typedef YacResult (*yapiFuncGetStmtAttr)(YacHandle hStmt, YapiStmtAttr attr, void* value, int32_t bufLength,
                                         int32_t* stringLength);

typedef YacResult (*yapiFuncDescribeCol2)(YacHandle hStmt, uint16_t id, YapiColumnDesc* desc);
typedef YacResult (*yapiFuncBindColumn)(YacHandle hStmt, uint16_t id, YapiType type, void* value, int32_t bufLen,
                                        int32_t* indicator);
typedef YacResult (*yapiFuncBindParameter)(YacHandle hStmt, uint16_t id, YapiParamDirection direction,
                                           YapiType bindType, void* value, int32_t bindSize, int32_t bufLength,
                                           int32_t* indicator);
typedef YacResult (*yapiFuncBindParameterByName)(YacHandle hStmt, char* name, YapiParamDirection direction,
                                                 YapiType bindType, void* value, int32_t bindSize, int32_t bufLength,
                                                 int32_t* indicator);
typedef YacResult (*yapiFuncNumResultCols)(YacHandle hStmt, int16_t* count);
typedef YacResult (*yapiFuncColAttribute)(YacHandle hStmt, uint16_t id, YapiColAttr attr, void* value, int32_t bufLen,
                                          int32_t* stringLength);
typedef YacResult (*yapiFuncNumParams)(YacHandle hStmt, int16_t* count);

typedef YapiDate (*yapiFuncNow)();
typedef void (*yapiFuncNumberFromInt32)(YapiNumber* n, int32_t v);
typedef YacResult (*yapiFuncText2Timestamp)(char* text, char* format, YapiDate* stamp, int16_t* bias);
typedef YacResult (*yapiFuncText2YMInterval)(char* str, YapiYMInterval* interval);
typedef YacResult (*yapiFuncText2DSInterval)(char* str, YapiDSInterval* interval);
typedef YacResult (*yapiFuncText2ShortTime)(char* str, char* format, YapiShortTime* shortTime);
typedef YacResult (*yapiFuncGetDateStruct)(YapiDate date, YapiDateStruct* ds);

// multi insert API in "batch" mode
typedef YacResult (*yapiFuncBatchInsertPrepare)(YacHandle hStmt, char* tableName);
typedef YacResult (*yapiFuncBatchInsertExecute)(YacHandle hStmt);

// lob API
typedef YacResult (*yapiFuncLobDescAlloc)(YacHandle hConn, YapiType type, void** desc);
typedef YacResult (*yapiFuncLobDescFree)(void* desc, YapiType type);
typedef YacResult (*yapiFuncLobGetChunkSize)(YacHandle hConn, YapiLobLocator* locator, uint16_t* chunkSize);
typedef YacResult (*yapiFuncLobGetLength)(YacHandle hConn, YapiLobLocator* locator, uint64_t* length);
typedef YacResult (*yapiFuncLobRead)(YacHandle hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf,
                                     uint64_t bufLen);
typedef YacResult (*yapiFuncLobWrite)(YacHandle hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf,
                                      uint64_t bufLen);
typedef YacResult (*yapiFuncLobCreateTemporary)(YacHandle hConn, YapiLobLocator* loc);
typedef YacResult (*yapiFuncLobFreeTemporary)(YacHandle hConn, YapiLobLocator* loc);

// dataType API
typedef YacResult (*yapiFuncDateGetDate)(const YapiDate date, int16_t* year, uint8_t* month, uint8_t* day);
typedef YacResult (*yapiFuncShortTimeGetShortTime)(const YapiShortTime time, uint8_t* hour, uint8_t* minute,
                                                   uint8_t* second, uint32_t* fraction);
typedef YacResult (*yapiFuncTimestampGetTimestamp)(const YapiTimestamp timestamp, int16_t* year, uint8_t* month,
                                                   uint8_t* day, uint8_t* hour, uint8_t* minute, uint8_t* second,
                                                   uint32_t* fraction);
typedef YacResult (*yapiFuncYMIntervalGetYearMonth)(const YapiYMInterval ymInterval, int32_t* year, int32_t* month);
typedef YacResult (*yapiFuncDSIntervalGetDaySecond)(const YapiDSInterval dsInterval, int32_t* day, int32_t* hour,
                                                    int32_t* minute, int32_t* second, int32_t* fraction);

typedef YacResult (*yapiFuncDateSetDate)(YapiDate* date, int16_t year, uint8_t month, uint8_t day);
typedef YacResult (*yapiFuncShortTimeSetShortTime)(YapiShortTime* time, uint8_t hour, uint8_t minute, uint8_t second,
                                                   uint32_t fraction);
typedef YacResult (*yapiFuncTimestampSetTimestamp)(YapiTimestamp* timestamp, int16_t year, uint8_t month, uint8_t day,
                                                   uint8_t hour, uint8_t minute, uint8_t second, uint32_t fraction);
typedef YacResult (*yapiFuncYMIntervalSetYearMonth)(YapiYMInterval* ymInterval, int32_t year, int32_t month);
typedef YacResult (*yapiFuncDSIntervalSetDaySecond)(YapiDSInterval* dsInterval, int32_t day, int32_t hour,
                                                    int32_t minute, int32_t second, int32_t fraction);

typedef YacResult (*yapiFuncNumberRound)(YapiNumber* n, int32_t precision, int32_t scale);

typedef YacResult (*yapiFuncPdbgStart)(YacHandle stmt, char* procName, uint32_t procNameLen);

typedef YacResult (*yapiFuncPdbgAbort)(YacHandle stmt);

typedef YacResult (*yapiFuncPdbgContinue)(YacHandle stmt);

typedef YacResult (*yapiFuncPdbgStepInto)(YacHandle stmt);

typedef YacResult (*yapiFuncPdbgStepOut)(YacHandle stmt);

typedef YacResult (*yapiFuncPdbgStepNext)(YacHandle stmt);

typedef YacResult (*yapiFuncPdbgShowSource)(YacHandle stmt);

typedef YacResult (*yapiFuncPdbgDeleteAllBreakpoints)(YacHandle stmt);

typedef YacResult (*yapiFuncPdbgAddBreakpoint)(YacHandle stmt, int lineNum, uint32_t* bpID);

typedef YacResult (*yapiFuncPdbgDeleteBreakpoint)(YacHandle stmt, uint32_t bpID);

typedef YacResult (*yapiFuncPdbgShowBreakpoints)(YacHandle stmt);

typedef YacResult (*yapiFuncPdbgShowFrameVariables)(YacHandle stmt);

typedef YacResult (*yapiFuncPdbgShowFrames)(YacHandle stmt);

typedef struct StYapiSymbols {
    yapiFuncAllocHandle fnAllocHandle;
    yapiFuncFreeHandler fnHandleFree;

    yapiFuncConnect    fnConnect;
    yapiFuncDisconnect fnDisconnect;

    yapiFuncGetVersion   fnGetVersion;
    yapiFuncGetLastError fnGetLastError;

    yapiFuncCancel        fnCancel;
    yapiFuncDirectExecute fnDirectExecute;
    yapiFuncPrepare       fnPrepare;
    yapiFuncExecute       fnExecute;
    yapiFuncFetch         fnFetch;
    yapiFuncCommit        fnCommit;
    yapiFuncRollback      fnRollback;

    yapiFuncSetEnvAttr  fnSetEnvAttr;
    yapiFuncGetEnvAttr  fnGetEnvAttr;
    yapiFuncSetConnAttr fnSetConnAttr;
    yapiFuncGetConnAttr fnGetConnAttr;
    yapiFuncSetStmtAttr fnSetStmtAttr;
    yapiFuncGetStmtAttr fnGetStmtAttr;

    yapiFuncDescribeCol2        fnDescribeCol2;
    yapiFuncBindColumn          fnBindColumn;
    yapiFuncBindParameter       fnBindParameter;
    yapiFuncBindParameterByName fnBindParameterByName;
    yapiFuncNumResultCols       fnNumResultCols;
    yapiFuncColAttribute        fnColAttribute;
    yapiFuncNumParams           fnNumParams;

    yapiFuncNow             fnNow;
    yapiFuncNumberFromInt32 fnNumberFromInt32;
    yapiFuncText2Timestamp  fnText2Timestamp;
    yapiFuncText2YMInterval fnText2YMInterval;
    yapiFuncText2DSInterval fnText2DSInterval;
    yapiFuncText2ShortTime  fnText2ShortTime;
    yapiFuncGetDateStruct   fnGetDateStruct;

    yapiFuncBatchInsertPrepare fnBatchInsertPrepare;
    yapiFuncBatchInsertExecute fnBatchInsertExecute;

    yapiFuncLobDescAlloc       fnLobDescAlloc;
    yapiFuncLobDescFree        fnLobDescFree;
    yapiFuncLobGetChunkSize    fnLobGetChunkSize;
    yapiFuncLobGetLength       fnLobGetLength;
    yapiFuncLobRead            fnLobRead;
    yapiFuncLobWrite           fnLobWrite;
    yapiFuncLobCreateTemporary fnLobCreateTemporary;
    yapiFuncLobFreeTemporary   fnLobFreeTemporary;

    yapiFuncDateGetDate            fnDateGetDate;
    yapiFuncShortTimeGetShortTime  fnShortTimeGetShortTime;
    yapiFuncTimestampGetTimestamp  fnTimestampGetTimestamp;
    yapiFuncYMIntervalGetYearMonth fnYMIntervalGetYearMonth;
    yapiFuncDSIntervalGetDaySecond fnDSIntervalGetDaySecond;

    yapiFuncDateSetDate            fnDateSetDate;
    yapiFuncShortTimeSetShortTime  fnShortTimeSetShortTime;
    yapiFuncTimestampSetTimestamp  fnTimestampSetTimestamp;
    yapiFuncYMIntervalSetYearMonth fnYMIntervalSetYearMonth;
    yapiFuncDSIntervalSetDaySecond fnDSIntervalSetDaySecond;

    yapiFuncNumberRound fnNumberRound;

    yapiFuncPdbgStart    fnPdbgStart;
    yapiFuncPdbgAbort    fnPdbgAbort;
    yapiFuncPdbgContinue fnPdbgContinue;
    yapiFuncPdbgStepInto fnPdbgStepInto;
    yapiFuncPdbgStepOut  fnPdbgStepOut;
    yapiFuncPdbgStepNext fnPdbgStepNext;

    yapiFuncPdbgShowSource           fnPdbgShowSource;
    yapiFuncPdbgDeleteAllBreakpoints fnPdbgDeleteAllBreakpoints;
    yapiFuncPdbgAddBreakpoint        fnPdbgAddBreakpoint;
    yapiFuncPdbgDeleteBreakpoint     fnPdbgDeleteBreakpoin;
    yapiFuncPdbgShowBreakpoints      fnPdbgShowBreakpoints;
    yapiFuncPdbgShowFrameVariables   fnPdbgShowFrameVariables;
    yapiFuncPdbgShowFrames           fnPdbgShowFrames;

} YapiSymbols;

#define T2S_BUFFER_SIZE 4096

typedef struct StYapiErrorBuffer {
    int32_t     code;
    uint32_t    messageLen;
    char        message[T2S_BUFFER_SIZE];
    char        sqlState[YAPI_MAX_SQLSTAT_LEN];
    YapiTextPos pos;
} YapiErrorBuffer;

typedef struct StYapiErrorMsg {
    YapiErrorBuffer* buf;
} YapiErrorMsg;

typedef enum {
    YAPI_ERR_NO_ERR = 20000,
    YAPI_ERR_LOAD_SYMBOL,
    YAPI_ERR_ALLOC_MEM,
} yapiErrorNum;

typedef struct StYapiEnv {
    uint32_t  version;
    YacHandle envHandler;
} YapiEnv;

typedef struct StYapiConnect {
    YapiEnv*  env;
    YacHandle connHandler;
} YapiConnect;

typedef struct StYapiStmt {
    YapiConnect* conn;
    YacHandle    stmtHandler;
} YapiStmt;

YapiResult yapiOpenDynamicLib(char* libName, YapiPointer* handler, YapiErrorMsg* error);
void       yapiSetError(YapiErrorMsg* error, yapiErrorNum errorNum, const char* format, ...);

YapiResult yapiCliAllocHandle(YapiHandleType type, YacHandle input, YacHandle* output, YapiErrorMsg* error);
YapiResult yapiCliFreeHandle(YapiHandleType type, YacHandle handle, YapiErrorMsg* error);
YapiResult yapiCliGetVersion(char** version, YapiErrorMsg* error);
YapiResult yapiCliGetLastError(YapiErrorMsg* error);

YapiResult yapiCliSetEnvAttr(YapiEnv* hEnv, YapiEnvAttr attr, void* value, int32_t length, YapiErrorMsg* error);
YapiResult yapiCliGetEnvAttr(YacHandle hEnv, YapiEnvAttr attr, void* value, int32_t bufLength, int32_t* stringLength,
                             YapiErrorMsg* error);

YapiResult yapiCliConnect(YacHandle hConn, const char* url, int16_t urlLength, const char* user, int16_t userLength,
                          const char* password, int16_t passwordLength, YapiErrorMsg* error);
YapiResult yapiCliDisconnect(YacHandle hConn, YapiErrorMsg* error);
YapiResult yapiCliCommit(YacHandle hConn, YapiErrorMsg* error);
YapiResult yapiCliRollback(YacHandle hConn, YapiErrorMsg* error);
YapiResult yapiCliSetConnAttr(YacHandle hConn, YapiConnAttr attr, void* value, int32_t length, YapiErrorMsg* error);
YapiResult yapiCliGetConnAttr(YacHandle hConn, YapiConnAttr attr, void* value, int32_t bufLength, int32_t* stringLength,
                              YapiErrorMsg* error);
YapiResult yapiCliCancel(YacHandle hConn, YapiErrorMsg* error);

YapiResult yapiCliDirectExecute(YacHandle hStmt, const char* sql, int32_t sqlLength, YapiErrorMsg* error);
YapiResult yapiCliPrepare(YacHandle hStmt, const char* sql, int32_t sqlLength, YapiErrorMsg* error);
YapiResult yapiCliExecute(YacHandle hStmt, YapiErrorMsg* error);
YapiResult yapiCliSetStmtAttr(YacHandle hStmt, YapiStmtAttr attr, void* value, int32_t length, YapiErrorMsg* error);
YapiResult yapiCliGetStmtAttr(YacHandle hStmt, YapiStmtAttr attr, void* value, int32_t bufLength, int32_t* stringLength,
                              YapiErrorMsg* error);
YapiResult yapiCliFetch(YacHandle hStmt, uint32_t* rows, YapiErrorMsg* error);
YapiResult yapiCliDescribeCol2(YacHandle hStmt, uint16_t id, YapiColumnDesc* desc, YapiErrorMsg* error);
YapiResult yapiCliBindColumn(YacHandle hStmt, uint16_t id, YapiType type, YapiPointer value, int32_t bufLen,
                             int32_t* indicator, YapiErrorMsg* error);
YapiResult yapiCliBindParameter(YacHandle hStmt, uint16_t id, YapiParamDirection direction, YapiType bindType,
                                YapiPointer value, int32_t bindSize, int32_t bufLength, int32_t* indicator,
                                YapiErrorMsg* error);
YapiResult yapiCliBindParameterByName(YacHandle hStmt, char* name, YapiParamDirection direction, YapiType bindType,
                                      YapiPointer value, int32_t bindSize, int32_t bufLength, int32_t* indicator,
                                      YapiErrorMsg* error);
YapiResult yapiCliNumResultCols(YacHandle hStmt, int16_t* count, YapiErrorMsg* error);
YapiResult yapiCliColAttribute(YacHandle hStmt, uint16_t id, YapiColAttr attr, void* value, int32_t bufLen,
                               int32_t* stringLength, YapiErrorMsg* error);
YapiResult yapiCliNumParams(YacHandle hStmt, int16_t* count, YapiErrorMsg* error);

YapiResult yapiCliGetDateStruct(YapiDate date, YapiDateStruct* ds, YapiErrorMsg* error);

YapiResult yapiCliLobDescAlloc(YacHandle* hConn, YapiType type, void** desc, YapiErrorMsg* error);
YapiResult yapiCliLobDescFree(void* desc, YapiType type, YapiErrorMsg* error);
YapiResult yapiCliLobGetChunkSize(YacHandle* hConn, YapiLobLocator* locator, uint16_t* chunkSize, YapiErrorMsg* error);
YapiResult yapiCliLobGetLength(YacHandle* hConn, YapiLobLocator* locator, uint64_t* length, YapiErrorMsg* error);
YapiResult yapiCliLobRead(YacHandle* hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf, uint64_t bufLen,
                          YapiErrorMsg* error);
YapiResult yapiCliLobWrite(YacHandle* hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf, uint64_t bufLen,
                           YapiErrorMsg* error);
YapiResult yapiCliLobCreateTemporary(YacHandle* hConn, YapiLobLocator* loc, YapiErrorMsg* error);
YapiResult yapiCliLobFreeTemporary(YacHandle* hConn, YapiLobLocator* loc, YapiErrorMsg* error);

YapiResult yapiCliDateGetDate(const YapiDate date, int16_t* year, uint8_t* month, uint8_t* day, YapiErrorMsg* error);
YapiResult yapiCliShortTimeGetShortTime(const YapiShortTime time, uint8_t* hour, uint8_t* minute, uint8_t* second,
                                        uint32_t* fraction, YapiErrorMsg* error);
YapiResult yapiCliTimestampGetTimestamp(const YapiTimestamp timestamp, int16_t* year, uint8_t* month, uint8_t* day,
                                        uint8_t* hour, uint8_t* minute, uint8_t* second, uint32_t* fraction,
                                        YapiErrorMsg* error);
YapiResult yapiCliYMIntervalGetYearMonth(const YapiYMInterval ymInterval, int32_t* year, int32_t* month,
                                         YapiErrorMsg* error);
YapiResult yapiCliDSIntervalGetDaySecond(const YapiDSInterval dsInterval, int32_t* day, int32_t* hour, int32_t* minute,
                                         int32_t* second, int32_t* fraction, YapiErrorMsg* error);

YapiResult yapiCliDateSetDate(YapiDate* date, int16_t year, uint8_t month, uint8_t day, YapiErrorMsg* error);
YapiResult yapiCliShortTimeSetShortTime(YapiShortTime* time, uint8_t hour, uint8_t minute, uint8_t second,
                                        uint32_t fraction, YapiErrorMsg* error);
YapiResult yapiCliTimestampSetTimestamp(YapiTimestamp* timestamp, int16_t year, uint8_t month, uint8_t day,
                                        uint8_t hour, uint8_t minute, uint8_t second, uint32_t fraction,
                                        YapiErrorMsg* error);
YapiResult yapiCliYMIntervalSetYearMonth(YapiYMInterval* ymInterval, int32_t year, int32_t month, YapiErrorMsg* error);
YapiResult yapiCliDSIntervalSetDaySecond(YapiDSInterval* dsInterval, int32_t day, int32_t hour, int32_t minute,
                                         int32_t second, int32_t fraction, YapiErrorMsg* error);

YapiResult yapiCliNumberRound(YapiNumber* n, int32_t precision, int32_t scale, YapiErrorMsg* error);

void yapiInitError(YapiErrorMsg* error);
void yapiGetErrorInfo(YapiErrorMsg* error, YapiErrorInfo* info);
void yapiGetCliError(YapiErrorMsg* error);

YapiResult yapiAllocMem(const char* name, size_t numMembers, size_t memberSize, void** ptr, YapiErrorMsg* error);
void       yapiFreeMem(void* ptr);

YapiResult yapiCiPdbgStart(YacHandle hStmt, char* procName, uint32_t procNameLen, YapiErrorMsg* error);
YapiResult yapiCiPdbgAbort(YacHandle hStmt, YapiErrorMsg* error);
YapiResult yapiCiPdbgContinue(YacHandle hStmt, YapiErrorMsg* error);
YapiResult yapiCiPdbgStepInto(YacHandle hStmt, YapiErrorMsg* error);
YapiResult yapiCiPdbgStepOut(YacHandle hStmt, YapiErrorMsg* error);
YapiResult yapiCiPdbgStepNext(YacHandle hStmt, YapiErrorMsg* error);
YapiResult yapiCiPdbgShowSource(YacHandle hStmt, YapiErrorMsg* error);
YapiResult yapiCiPdbgDeleteAllBreakpoints(YacHandle hStmt, YapiErrorMsg* error);
YapiResult yapiCiPdbgAddBreakpoint(YacHandle hStmt, int lineNum, uint32_t* bpID, YapiErrorMsg* error);
YapiResult yapiCiPdbgDeleteBreakpoint(YacHandle hStmt, uint32_t bpID, YapiErrorMsg* error);
YapiResult yapiCiPdbgShowBreakpoints(YacHandle hStmt, YapiErrorMsg* error);
YapiResult yapiCiPdbgShowFrameVariables(YacHandle hStmt, YapiErrorMsg* error);
YapiResult yapiCiPdbgShowFrames(YacHandle hStmt, YapiErrorMsg* error);

#ifdef __cplusplus
}
#endif

#endif