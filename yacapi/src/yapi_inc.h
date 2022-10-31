#ifndef YAPI_INC_H
#define YAPI_INC_H

#include "yacapi.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef enum EnYacResult { YA_SUCCESS = 0, YAC_SUCCESS_WITH_INFO = 1, YAC_ERROR = -1 } YacResult;
typedef void* YacHandle;

typedef YacResult (*yapiFuncAllocHandle)(YapiHandleType type, YacHandle input, YacHandle* output);
typedef YacResult (*yapiFuncFreeHandler)(YapiHandleType type, YacHandle handle);
typedef char* (*yapiFuncGetVersion)();
typedef void (*yapiFuncGetLastError)(int32_t* errCode, char** message, char** sqlState, YapiTextPos* pos);
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
                                           YapiType bindType, void* value, uint32_t bindSize, int32_t bufLength,
                                           int32_t* indicator);
typedef YacResult (*yapiFuncBindParameterByName)(YacHandle hStmt, char* name, YapiParamDirection direction,
                                                 YapiType bindType, void* value, uint32_t bindSize, int32_t bufLength,
                                                 int32_t* indicator);
typedef YacResult (*yapiFuncNumResultCols)(YacHandle hStmt, int16_t* count);

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
} YapiSymbols;

#define T2S_BUFFER_SIZE 4096

typedef struct StYapiErrorMsg {
    uint32_t code;
    uint32_t messageLen;
    char     message[T2S_BUFFER_SIZE];
} YapiErrorMsg;

typedef enum {
    YAPI_ERR_NO_ERR = 20000,
    YAPI_ERR_LOAD_SYMBOL,
} yapiErrorNum;

typedef struct StYapiEnv {
    uint32_t version;
    void*    envHandler;
} YapiEnv;

typedef struct StYapiConnect {
    YapiEnv* env;
    void*    connHandler;
} YapiConnect;

typedef struct StYapiStmt {
    YapiConnect* conn;
    void*        stmtHandler;
} YapiStmt;

YapiResult yapiOpenDynamicLib(char* libName, YapiPointer* handler, YapiErrorMsg* error);
void       yapiSetError(YapiErrorMsg* error, yapiErrorNum errorNum, const char* format, ...);

YapiResult yapiCliAllocHandle(YapiHandleType type, YacHandle input, YacHandle* output);
YapiResult yapiCliFreeHandle(YapiHandleType type, YacHandle handle);
YapiResult yapiCliGetVersion(char** version);
YapiResult yapiCliGetLastError(int32_t* errCode, char** message, char** sqlState, YapiTextPos* pos);

YapiResult yapiCliGetEnvAttr(YacHandle hEnv, YapiEnvAttr attr, void* value, int32_t bufLength, int32_t* stringLength);

YapiResult yapiCliConnect(YacHandle hConn, const char* url, int16_t urlLength, const char* user, int16_t userLength,
                          const char* password, int16_t passwordLength);
YapiResult yapiCliDisconnect(YacHandle hConn);
YapiResult yapiCliCommit(YacHandle hConn);
YapiResult yapiCliRollback(YacHandle hConn);
YapiResult yapiCliSetConnAttr(YacHandle hConn, YapiConnAttr attr, void* value, int32_t length);
YapiResult yapiCliGetConnAttr(YacHandle hConn, YapiConnAttr attr, void* value, int32_t bufLength,
                              int32_t* stringLength);
YapiResult yapiCliCancel(YacHandle hConn);

YapiResult yapiCliDirectExecute(YacHandle hStmt, const char* sql, int32_t sqlLength);
YapiResult yapiCliPrepare(YacHandle hStmt, const char* sql, int32_t sqlLength);
YapiResult yapiCliExecute(YacHandle hStmt);
YapiResult yapiCliSetStmtAttr(YacHandle hStmt, YapiStmtAttr attr, void* value, int32_t length);
YapiResult yapiCliGetStmtAttr(YacHandle hStmt, YapiStmtAttr attr, void* value, int32_t bufLength,
                              int32_t* stringLength);
YapiResult yapiCliFetch(YacHandle hStmt, uint32_t* rows);
YapiResult yapiCliDescribeCol2(YacHandle hStmt, uint16_t id, YapiColumnDesc* desc);
YapiResult yapiCliBindColumn(YacHandle hStmt, uint16_t id, YapiType type, YapiPointer value, int32_t bufLen,
                             int32_t* indicator);
YapiResult yapiCliBindParameter(YacHandle hStmt, uint16_t id, YapiParamDirection direction, YapiType bindType,
                                YapiPointer value, uint32_t bindSize, int32_t bufLength, int32_t* indicator);
YapiResult yapiCliBindParameterByName(YacHandle hStmt, char* name, YapiParamDirection direction, YapiType bindType,
                                      YapiPointer value, uint32_t bindSize, int32_t bufLength, int32_t* indicator);
YapiResult yapiCliNumResultCols(YacHandle hStmt, int16_t* count);

YapiResult yapiCliGetDateStruct(YapiDate date, YapiDateStruct* ds);

YapiResult yapiCliLobDescAlloc(YapiConnect* hConn, YapiType type, void** desc);
YapiResult yapiCliLobDescFree(void* desc, YapiType type);
YapiResult yapiCliLobGetChunkSize(YapiConnect* hConn, YapiLobLocator* locator, uint16_t* chunkSize);
YapiResult yapiCliLobGetLength(YapiConnect* hConn, YapiLobLocator* locator, uint64_t* length);
YapiResult yapiCliLobRead(YapiConnect* hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf, uint64_t bufLen);
YapiResult yapiCliLobWrite(YapiConnect* hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf, uint64_t bufLen);
YapiResult yapiCliLobCreateTemporary(YapiConnect* hConn, YapiLobLocator* loc);
YapiResult yapiCliLobFreeTemporary(YapiConnect* hConn, YapiLobLocator* loc);

#ifdef __cplusplus
}
#endif

#endif