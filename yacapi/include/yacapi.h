#ifndef YAPI_API_H
#define YAPI_API_H

#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

#define YAPI_TRUE true
#define YAPI_FALSE false

#define YAPI_PARAM_NAME_BUFFER_SIZE 32
#define YAPI_MIN_PACKET_SIZE KB(64)
#define YAPI_MAX_PACKET_SIZE MB(32)
#define YAPI_NULL_DATA -1
#define YAPI_NULL_TERM_STR -2
#define YAPI_MAX_ERROR_MSG_LEN 4096
#define YAPI_MAX_SQLSTAT_LEN 16

typedef int64_t YapiDate;
typedef int64_t YapiShortTime;
typedef int32_t YapiYMInterval;
typedef int64_t YapiDSInterval;
typedef void*   YapiPointer;

typedef struct StYapiConnect YapiConnect;
typedef struct StYapiStmt    YapiStmt;
typedef struct StYapiEnv     YapiEnv;

#pragma pack(4)
typedef struct StYapiNumber {
    uint64_t item[2];
    int8_t   sign;
    uint8_t  unused;
    int16_t  exp;
} YapiNumber;

typedef struct StYapiTimestamp {
    int64_t stamp;
    int16_t bias;  // minutes
    int16_t unused;
} YapiTimestamp;
#pragma pack()

typedef struct StYapiDateStruct {
    uint16_t year;
    uint8_t  month;
    uint8_t  day;
    uint8_t  hour;
    uint8_t  minute;
    uint8_t  second;
    uint8_t  dayOfWeek;
    uint8_t  weekName;
    uint16_t dayOfYear;
    uint8_t  unused[5];
    uint32_t fraction;
    uint32_t secondOfDay;
} YapiDateStruct;

typedef enum EnYapiType {
    YAPI_TYPE_UNKNOWN = 0,
    YAPI_TYPE_BOOL = 1,
    YAPI_TYPE_TINYINT = 2,
    YAPI_TYPE_SMALLINT = 3,
    YAPI_TYPE_INTEGER = 4,
    YAPI_TYPE_BIGINT = 5,
    YAPI_TYPE_UTINYINT = 6,
    YAPI_TYPE_USMALLINT = 7,
    YAPI_TYPE_UINTEGER = 8,
    YAPI_TYPE_UBIGINT = 9,
    YAPI_TYPE_FLOAT = 10,
    YAPI_TYPE_DOUBLE = 11,
    YAPI_TYPE_NUMBER = 12,
    YAPI_TYPE_DATE = 13,
    YAPI_TYPE_SHORTDATE = 14,
    YAPI_TYPE_SHORTTIME = 15,
    YAPI_TYPE_TIMESTAMP = 16,
    YAPI_TYPE_TIMESTAMP_TZ = 17,
    YAPI_TYPE_TIMESTAMP_LTZ = 18,
    YAPI_TYPE_YM_INTERVAL = 19,
    YAPI_TYPE_DS_INTERVAL = 20,
    // 21-23 reversed
    YAPI_TYPE_CHAR = 24,
    YAPI_TYPE_NCHAR = 25,
    YAPI_TYPE_VARCHAR = 26,
    YAPI_TYPE_NVARCHAR = 27,
    YAPI_TYPE_BINARY = 28,
    YAPI_TYPE_CLOB = 29,
    YAPI_TYPE_BLOB = 30,
    YAPI_TYPE_BIT = 31,
    YAPI_TYPE_ROWID = 32,
    YAPI_TYPE_NCLOB = 33,
    YAPI_TYPE_CURSOR = 34,
    __YAPI_TYPES_COUNT__
} YapiType;

typedef enum EnYapiResult { YAPI_SUCCESS = 0, YAPI_SUCCESS_WITH_INFO = 1, YAPI_ERROR = -1 } YapiResult;

typedef struct StYapiTextPos {
    int32_t line;
    int32_t column;
} YapiTextPos;

typedef struct StYapiErrorInfo {
    int32_t     errCode;
    char        message[YAPI_MAX_ERROR_MSG_LEN];
    char        sqlState[YAPI_MAX_SQLSTAT_LEN];
    YapiTextPos pos;
} YapiErrorInfo;

typedef struct StYapiColumnDesc {
    char*    name;
    uint32_t size;
    uint8_t  type;
    uint8_t  precision;
    int8_t   scale;
    uint8_t  nullable;
} YapiColumnDesc;

typedef enum EnYapiBindType {
    YAPI_BIND_COLUMN = 0,
    YAPI_BIND_PARAM = 1,
    __YAPI_BIND_TYPE_COUNT__,
} YapiBindType;

typedef enum EnYapiParamDirection {
    YAPI_PARAM_INPUT = 1,
    YAPI_PARAM_OUTPUT = 2,
    YAPI_PARAM_INOUT = 3,
} YapiParamDirection;

typedef enum EnYapiHandleType {
    YAPI_HANDLE_UNKNOWN = 0,
    YAPI_HANDLE_ENV = 1,
    YAPI_HANDLE_DBC = 2,
    YAPI_HANDLE_STMT = 3,
    YAPI_HANDLE_DESC = 4,
    YAPI_HANDLE_PUMP = 5,
    __YAPI_HANDLE_COUNT__
} YapiHandleType;

typedef enum EnYapiEnvAttr {
    __YAPI_ENV_ATTR_BEGIN__ = 60,
    YAPI_ATTR_DATE_FORMAT = 60,
    YAPI_ATTR_CHARSET = 61,
    __YAPI_ENV_ATTR_END__
} YapiEnvAttr;

typedef enum EnYapiConnAttr {
    __YAPI_CONN_ATTR_BEGIN__ = 1,
    YAPI_ATTR_LOGONINFO_PTR = 1,
    YAPI_ATTR_ASYNC_ENABLE,
    YAPI_ATTR_AUTOCOMMIT,
    YAPI_ATTR_LOGIN_TIMEOUT,
    YAPI_ATTR_STMTS,
    YAPI_ATTR_PACKET_SIZE,
    YAPI_ATTR_TXN_ISOLATION,
    YAPI_ATTR_SERVEROUTPUT,
    YAPI_ATTR_NUMWIDTH,
    YAPI_ATTR_AUTOTRACE,
    __YAPI_CONN_ATTR_END__
} YapiConnAttr;

typedef enum EnYapiStmtAttr {
    __YAPI_STMT_ATTR_BEGIN__ = 100,
    YAPI_ATTR_PARAMSET_SIZE = 100,
    YAPI_ATTR_ROWSET_SIZE,
    YAPI_ATTR_ROWS_FETECHED,
    YAPI_ATTR_ROWS_AFFECTED,
    YAPI_ATTR_CURSOR_EOF,
    YAPI_ATTR_SQLTYPE,
    YAPI_ATTR_IS_BATCHROWS,
    YAPI_ATTR_IS_BATCH_ERRORS,
    YAPI_ATTR_ACK_BATCHROWS_SIZE,
    YAPI_ATTR_ACK_BATCH_ERRORS_SIZE,
    YAPI_ATTR_ACK_BATCHROWS,
    YAPI_ATTR_ACK_BATCH_ERRORS,
    __YAPI_STMT_ATTR_END__
} YapiStmtAttr;

typedef struct StYapiBatchError {
    uint32_t rowNum;
    uint32_t errCode;
    uint32_t msgLen;
    char     msg[YAPI_MAX_ERROR_MSG_LEN];
} YapiBatchError;

typedef struct StYapiLobLocator YapiLobLocator;

typedef enum EnYapiSQLType {
    YAPI_SQLTYPE_QUERY = 1,
    YAPI_SQLTYPE_INSERT,
    YAPI_SQLTYPE_UPDATE,
    YAPI_SQLTYPE_DELETE,
    YAPI_SQLTYPE_MERGE,
    YAPI_SQLTYPE_WITH,
    YAPI_SQLTYPE_ANONYMOUS_BLOCK,
    YAPI_SQLTYPE_DML_CEIL = 10,

    YAPI_SQLTYPE_CREATE_DATABASE = 11,
    YAPI_SQLTYPE_CREATE_DATASPACE,
    YAPI_SQLTYPE_CREATE_TABLESPACE_SET,
    YAPI_SQLTYPE_CREATE_TABLESPACE,
    YAPI_SQLTYPE_CREATE_TABLE,
    YAPI_SQLTYPE_CREATE_SHARDED_TABLE,
    YAPI_SQLTYPE_CREATE_DUPLICATED_TABLE,
    YAPI_SQLTYPE_CREATE_TEMP_TABLE,
    YAPI_SQLTYPE_CREATE_INDEX,
    YAPI_SQLTYPE_CREATE_AC,
    YAPI_SQLTYPE_CREATE_VIEW,
    YAPI_SQLTYPE_CREATE_SYNONYM,
    YAPI_SQLTYPE_CREATE_PROCEDURE,
    YAPI_SQLTYPE_CREATE_FUNCTION,
    YAPI_SQLTYPE_CREATE_TRIGGER,
    YAPI_SQLTYPE_CREATE_PACKAGE,
    YAPI_SQLTYPE_CREATE_OR_REPLACE,
    YAPI_SQLTYPE_CREATE_SEQUENCE,
    YAPI_SQLTYPE_CREATE_USER,
    YAPI_SQLTYPE_ALTER_DATABASE,
    YAPI_SQLTYPE_ALTER_TABLE,
    YAPI_SQLTYPE_ALTER_INDEX,
    YAPI_SQLTYPE_ALTER_TABLESPACE_SET,
    YAPI_SQLTYPE_ALTER_TABLESPACE,
    YAPI_SQLTYPE_ALTER_SEQUENCE,
    YAPI_SQLTYPE_ALTER_SYSTEM,
    YAPI_SQLTYPE_ALTER_SESSION,
    YAPI_SQLTYPE_DROP_DATASPACE,
    YAPI_SQLTYPE_DROP_TABLESPACE_SET,
    YAPI_SQLTYPE_DROP_TABLE,
    YAPI_SQLTYPE_DROP_INDEX,
    YAPI_SQLTYPE_DROP_AC,
    YAPI_SQLTYPE_DROP_SEQUENCE,
    YAPI_SQLTYPE_DROP_TABLESPACE,
    YAPI_SQLTYPE_DROP_VIEW,
    YAPI_SQLTYPE_DROP_SYNONYM,
    YAPI_SQLTYPE_DROP_USER,
    YAPI_SQLTYPE_DROP_PROCEDURE,
    YAPI_SQLTYPE_DROP_FUNCTION,
    YAPI_SQLTYPE_DROP_PACKAGE,
    YAPI_SQLTYPE_TRUNCATE_TABLE,
    YAPI_SQLTYPE_BACKUP_DATABASE,
    YAPI_SQLTYPE_RESTORE_DATABASE,
    YAPI_SQLTYPE_RECOVER_DATABASE,
    YAPI_SQLTYPE_BUILD_DATABASE,
    YAPI_SQLTYPE_SET_TRANSACTION,
    YAPI_SQLTYPE_PREPARE_TRANSACTION,
    YAPI_SQLTYPE_REPLACE_VIEW,
    YAPI_SQLTYPE_REPLACE_SYNONYM,
    YAPI_SQLTYPE_REPLACE_FUNCTION,
    YAPI_SQLTYPE_REPLACE_PROCEDURE,
    YAPI_SQLTYPE_REPLACE_PACKAGE,
    YAPI_SQLTYPE_FLASHBACK_TABLE,
    YAPI_SQLTYPE_COMMENT,
    YAPI_SQLTYPE_PURGE_RECYCLEBIN,
    YAPI_SQLTYPE_CREATE_ROLE,
    YAPI_SQLTYPE_DROP_ROLE,
    YAPI_SQLTYPE_GRANT,
    YAPI_SQLTYPE_REVOKE,
    YAPI_SQLTYPE_ALTER_USER,
    YAPI_SQLTYPE_ALTER_TRIGGER,
    YAPI_SQLTYPE_DROP_TRIGGER,
    YAPI_SQLTYPE_REPLACE_TRIGGER,
    YAPI_SQLTYPE_REPLACE_OUTLINE,
    YAPI_SQLTYPE_CREATE_AUDIT_POLICY,
    YAPI_SQLTYPE_ALTER_AUDIT_POLICY,
    YAPI_SQLTYPE_DROP_AUDIT_POLICY,
    YAPI_SQLTYPE_AUDIT_POLICY,
    YAPI_SQLTYPE_NOAUDIT_POLICY,
    YAPI_SQLTYPE_CREATE_OUTLINE,
    YAPI_SQLTYPE_ALTER_OUTLINE,
    YAPI_SQLTYPE_DROP_OUTLINE,
    YAPI_SQLTYPE_ANAYLZE_TABLE,
    YAPI_SQLTYPE_DDL_CEIL = 128,

    YAPI_SQLTYPE_COMMIT,
    YAPI_SQLTYPE_ROLLBACK,
    YAPI_SQLTYPE_EXPLAIN,
    YAPI_SQLTYPE_SAVEPOINT,
    YAPI_SQLTYPE_SHUTDOWN,
    YAPI_SQLTYPE_RELEASE_SAVEPOINT,
    YAPI_SQLTYPE_SO_ENTITY,
    YAPI_SQLTYPE_PACK_ENTITY,

    YAPI_SQLTYPE_BATCH_INSERT,
    /* distributed inner type */
    YAPI_SQLTYPE_DXG_QUERY,

    __YAPI_SQLTYPE_COUNT__ = 255
} YapiSQLType;

char* yapiGetVersion(YapiEnv* inst);
void  yapiGetLastError(YapiErrorInfo* info);

//-----------------------------------------------------------------------------
// Enviment Function
//-----------------------------------------------------------------------------
YapiResult yapiAllocEnv(YapiEnv** inst);
YapiResult yapiReleaseEnv(YapiEnv* inst);

//-----------------------------------------------------------------------------
// Session Function
//-----------------------------------------------------------------------------
YapiResult yapiConnect(YapiEnv* env, const char* url, int16_t urlLength, const char* user, int16_t userLength,
                       const char* password, int16_t passwordLength, YapiConnect** hConn);
YapiResult yapiDisconnect(YapiConnect* hConn);
YapiResult yapiReleaseConn(YapiConnect* hConn);
YapiResult yapiCancel(YapiConnect* hConn);
YapiResult yapiCommit(YapiConnect* hConn);
YapiResult yapiRollback(YapiConnect* hConn);
YapiResult yapiSetConnAttr(YapiConnect* hConn, YapiConnAttr attr, void* value, int32_t length);
YapiResult yapiGetConnAttr(YapiConnect* hConn, YapiConnAttr attr, void* value, int32_t bufLength,
                           int32_t* stringLength);

//-----------------------------------------------------------------------------
// Statment Function
//-----------------------------------------------------------------------------
YapiResult yapiPrepare(YapiConnect* hConn, const char* sql, int32_t sqlLength, YapiStmt** hStmt);
YapiResult yapiExecute(YapiStmt* hStmt);
YapiResult yapiFetch(YapiStmt* hStmt, uint32_t* rows);
YapiResult yapiDirectExecute(YapiStmt* hStmt, const char* sql, int32_t sqlLength);
YapiResult yapiDescribeCol2(YapiStmt* hStmt, uint16_t id, YapiColumnDesc* desc);
YapiResult yapiBindColumn(YapiStmt* hStmt, uint16_t id, YapiType type, YapiPointer value, int32_t bufLen,
                          int32_t* indicator);
YapiResult yapiBindParameter(YapiStmt* hStmt, uint16_t id, YapiParamDirection direction, YapiType bindType,
                             YapiPointer value, uint32_t bindSize, int32_t bufLength, int32_t* indicator);
YapiResult yapiBindParameterByName(YapiStmt* hStmt, char* name, YapiParamDirection direction, YapiType bindType,
                                   YapiPointer value, uint32_t bindSize, int32_t bufLength, int32_t* indicator);
YapiResult yapiNumResultCols(YapiStmt* hStmt, int16_t* count);
YapiResult yapiSetStmtAttr(YapiStmt* hStmt, YapiStmtAttr attr, void* value, int32_t length);
YapiResult yapiGetStmtAttr(YapiStmt* hStmt, YapiStmtAttr attr, void* value, int32_t bufLength, int32_t* stringLength);
YapiResult yapiReleaseStmt(YapiStmt* hStmt);

//-----------------------------------------------------------------------------
// Enviment Function
//-----------------------------------------------------------------------------
YapiResult yapiGetEnvAttr(YapiEnv* hEnv, YapiEnvAttr attr, void* value, int32_t bufLength, int32_t* stringLength);

void       yacNumberFromInt32(YapiNumber* n, int32_t v);
YapiResult yacText2Timestamp(char* text, char* format, YapiDate* stamp, int16_t* bias);
YapiResult yacText2YMInterval(char* str, YapiYMInterval* interval);
YapiResult yacText2DSInterval(char* str, YapiDSInterval* interval);
YapiResult yacText2ShortTime(char* str, char* format, YapiShortTime* shortTime);
YapiResult yapiGetDateStruct(YapiDate date, YapiDateStruct* ds);

//-----------------------------------------------------------------------------
// Enviment Function
//-----------------------------------------------------------------------------
YapiResult yapiLobDescAlloc(YapiConnect* hConn, YapiType type, void** desc);
YapiResult yapiLobDescFree(void* desc, YapiType type);
YapiResult yapiLobGetChunkSize(YapiConnect* hConn, YapiLobLocator* locator, uint16_t* chunkSize);
YapiResult yapiLobGetLength(YapiConnect* hConn, YapiLobLocator* locator, uint64_t* length);
YapiResult yapiLobRead(YapiConnect* hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf, uint64_t bufLen);
YapiResult yapiLobWrite(YapiConnect* hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf, uint64_t bufLen);
YapiResult yapiLobCreateTemporary(YapiConnect* hConn, YapiLobLocator* loc);
YapiResult yapiLobFreeTemporary(YapiConnect* hConn, YapiLobLocator* loc);

#ifdef __cplusplus
}
#endif

#endif
