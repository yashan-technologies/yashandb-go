#ifndef YACLIC_H
#define YACLIC_H

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

#define YAC_TRUE true
#define YAC_FALSE false
#define YAC_MAX_COL_LEN 1024
#define YAC_MAX_SQL_LEN 4096
#define YAC_PARAM_NAME_BUFFER_SIZE 32
#define YAC_MAX_ERROR_MSG_LEN 256
#define YAC_MIN_PACKET_SIZE KB(64)
#define YAC_MAX_PACKET_SIZE MB(32)
#define YAC_NULL_DATA 0xFFFFFFFF
#define YAC_NULL_TERM_STR 0xFFFFFFFE
#define YAC_CALL(proc)                          \
    do {                                        \
        if ((YacResult)(proc) != YAC_SUCCESS) { \
            return YAC_ERROR;                   \
        }                                       \
    } while (0)

typedef char           YacInt8;
typedef unsigned char  YacUint8;
typedef short          YacInt16;
typedef unsigned short YacUint16;
typedef int            YacInt32;
typedef unsigned int   YacUint32;
typedef int64_t        YacInt64;
typedef uint64_t       YacUint64;
typedef char           YacChar;
typedef bool           YacBool;
typedef double         YacDouble;
typedef float          YacFloat;
typedef YacInt64       YacDate;
typedef YacInt64       YacShortTime;
typedef YacInt32       YacYMInterval;
typedef YacInt64       YacDSInterval;
typedef void           YacVoid;
typedef YacVoid*       YacPointer;
typedef YacVoid*       YacHandle;

#pragma pack(4)
typedef struct StYacNumber {
    YacUint64 item[2];
    YacInt8   sign;
    YacUint8  unused;
    YacInt16  exp;
} YacNumber;

typedef struct StYacTimestamp {
    YacInt64 stamp;
    YacInt16 bias;  // minutes
    YacInt16 unused;
} YacTimestamp;

#pragma pack()

typedef struct StYacDateStruct {
    YacUint16 year;
    YacUint8  month;
    YacUint8  day;
    YacUint8  hour;
    YacUint8  minute;
    YacUint8  second;
    YacUint8  dayOfWeek;
    YacUint8  weekName;
    YacUint16 dayOfYear;
    YacUint8  unused[5];
    YacUint32 fraction;
    YacUint32 secondOfDay;
} YacDateStruct;

typedef enum EnYacType {
    YAC_TYPE_UNKNOWN = 0,
    YAC_TYPE_BOOL = 1,
    YAC_TYPE_TINYINT = 2,
    YAC_TYPE_SMALLINT = 3,
    YAC_TYPE_INTEGER = 4,
    YAC_TYPE_BIGINT = 5,
    YAC_TYPE_UTINYINT = 6,
    YAC_TYPE_USMALLINT = 7,
    YAC_TYPE_UINTEGER = 8,
    YAC_TYPE_UBIGINT = 9,
    YAC_TYPE_FLOAT = 10,
    YAC_TYPE_DOUBLE = 11,
    YAC_TYPE_NUMBER = 12,
    YAC_TYPE_DATE = 13,
    YAC_TYPE_SHORTDATE = 14,
    YAC_TYPE_SHORTTIME = 15,
    YAC_TYPE_TIMESTAMP = 16,
    YAC_TYPE_TIMESTAMP_TZ = 17,
    YAC_TYPE_TIMESTAMP_LTZ = 18,
    YAC_TYPE_YM_INTERVAL = 19,
    YAC_TYPE_DS_INTERVAL = 20,
    // 21-23 reversed
    YAC_TYPE_CHAR = 24,
    YAC_TYPE_NCHAR = 25,
    YAC_TYPE_VARCHAR = 26,
    YAC_TYPE_NVARCHAR = 27,
    YAC_TYPE_BINARY = 28,
    YAC_TYPE_CLOB = 29,
    YAC_TYPE_BLOB = 30,
    YAC_TYPE_BIT = 31,
    YAC_TYPE_ROWID = 32,
    YAC_TYPE_NCLOB = 33,
    YAC_TYPE_CURSOR = 34,
    __YAC_TYPES_COUNT__
} YacType;

typedef enum EnYacResult { YAC_SUCCESS = 0, YAC_SUCCESS_WITH_INFO = 1, YAC_ERROR = -1 } YacResult;

typedef struct StYacTextPos {
    YacInt32 line;
    YacInt32 column;
} YacTextPos;

typedef struct StYacColumnDesc {
    YacChar*  name;
    YacUint32 size;
    YacUint8  type;
    YacUint8  precision;
    YacInt8   scale;
    YacUint8  nullable;
} YacColumnDesc;

typedef enum EnYacBindType {
    YAC_BIND_COLUMN = 0,
    YAC_BIND_PARAM = 1,
    __YAC_BIND_TYPE_COUNT__,
} YacBindType;

typedef enum EnYacParamDirection {
    YAC_PARAM_INPUT = 1,
    YAC_PARAM_OUTPUT = 2,
    YAC_PARAM_INOUT = 3,
} YacParamDirection;

typedef enum EnYacHandleType {
    YAC_HANDLE_UNKNOWN = 0,
    YAC_HANDLE_ENV = 1,
    YAC_HANDLE_DBC = 2,
    YAC_HANDLE_STMT = 3,
    YAC_HANDLE_DESC = 4,
    YAC_HANDLE_PUMP = 5,
    __YAC_HANDLE_COUNT__
} YacHandleType;

typedef enum EnYacEnvAttr {
    __YAC_ENV_ATTR_BEGIN__ = 60,
    YAC_ATTR_DATE_FORMAT = 60,
    YAC_ATTR_CHARSET = 61,
    __YAC_ENV_ATTR_END__
} YacEnvAttr;

typedef enum EnYacConnAttr {
    __YAC_CONN_ATTR_BEGIN__ = 1,
    YAC_ATTR_LOGONINFO_PTR = 1,
    YAC_ATTR_ASYNC_ENABLE,
    YAC_ATTR_AUTOCOMMIT,
    YAC_ATTR_LOGIN_TIMEOUT,
    YAC_ATTR_STMTS,
    YAC_ATTR_PACKET_SIZE,
    YAC_ATTR_TXN_ISOLATION,
    YAC_ATTR_SERVEROUTPUT,
    YAC_ATTR_NUMWIDTH,
    YAC_ATTR_AUTOTRACE,
    __YAC_CONN_ATTR_END__
} YacConnAttr;

typedef enum EnYacStmtAttr {
    __YAC_STMT_ATTR_BEGIN__ = 100,
    YAC_ATTR_PARAMSET_SIZE = 100,
    YAC_ATTR_ROWSET_SIZE,
    YAC_ATTR_ROWS_FETECHED,
    YAC_ATTR_ROWS_AFFECTED,
    YAC_ATTR_CURSOR_EOF,
    YAC_ATTR_SQLTYPE,
    YAC_ATTR_IS_BATCHROWS,
    YAC_ATTR_IS_BATCH_ERRORS,
    YAC_ATTR_ACK_BATCHROWS_SIZE,
    YAC_ATTR_ACK_BATCH_ERRORS_SIZE,
    YAC_ATTR_ACK_BATCHROWS,
    YAC_ATTR_ACK_BATCH_ERRORS,
    __YAC_STMT_ATTR_END__
} YacStmtAttr;

typedef struct StYacBatchError {
    YacUint32 rowNum;
    YacUint32 errCode;
    YacUint32 msgLen;
    YacChar   msg[YAC_MAX_ERROR_MSG_LEN];
} YacBatchError;

typedef struct StYacLobLocator YacLobLocator;

typedef enum EnYacSQLType {
    YAC_SQLTYPE_QUERY = 1,
    YAC_SQLTYPE_INSERT,
    YAC_SQLTYPE_UPDATE,
    YAC_SQLTYPE_DELETE,
    YAC_SQLTYPE_MERGE,
    YAC_SQLTYPE_WITH,
    YAC_SQLTYPE_ANONYMOUS_BLOCK,
    YAC_SQLTYPE_DML_CEIL = 10,

    YAC_SQLTYPE_CREATE_DATABASE = 11,
    YAC_SQLTYPE_CREATE_DATASPACE,
    YAC_SQLTYPE_CREATE_TABLESPACE_SET,
    YAC_SQLTYPE_CREATE_TABLESPACE,
    YAC_SQLTYPE_CREATE_TABLE,
    YAC_SQLTYPE_CREATE_SHARDED_TABLE,
    YAC_SQLTYPE_CREATE_DUPLICATED_TABLE,
    YAC_SQLTYPE_CREATE_TEMP_TABLE,
    YAC_SQLTYPE_CREATE_INDEX,
    YAC_SQLTYPE_CREATE_AC,
    YAC_SQLTYPE_CREATE_VIEW,
    YAC_SQLTYPE_CREATE_SYNONYM,
    YAC_SQLTYPE_CREATE_PROCEDURE,
    YAC_SQLTYPE_CREATE_FUNCTION,
    YAC_SQLTYPE_CREATE_TRIGGER,
    YAC_SQLTYPE_CREATE_PACKAGE,
    YAC_SQLTYPE_CREATE_OR_REPLACE,
    YAC_SQLTYPE_CREATE_SEQUENCE,
    YAC_SQLTYPE_CREATE_USER,
    YAC_SQLTYPE_ALTER_DATABASE,
    YAC_SQLTYPE_ALTER_TABLE,
    YAC_SQLTYPE_ALTER_INDEX,
    YAC_SQLTYPE_ALTER_TABLESPACE_SET,
    YAC_SQLTYPE_ALTER_TABLESPACE,
    YAC_SQLTYPE_ALTER_SEQUENCE,
    YAC_SQLTYPE_ALTER_SYSTEM,
    YAC_SQLTYPE_ALTER_SESSION,
    YAC_SQLTYPE_DROP_DATASPACE,
    YAC_SQLTYPE_DROP_TABLESPACE_SET,
    YAC_SQLTYPE_DROP_TABLE,
    YAC_SQLTYPE_DROP_INDEX,
    YAC_SQLTYPE_DROP_AC,
    YAC_SQLTYPE_DROP_SEQUENCE,
    YAC_SQLTYPE_DROP_TABLESPACE,
    YAC_SQLTYPE_DROP_VIEW,
    YAC_SQLTYPE_DROP_SYNONYM,
    YAC_SQLTYPE_DROP_USER,
    YAC_SQLTYPE_DROP_PROCEDURE,
    YAC_SQLTYPE_DROP_FUNCTION,
    YAC_SQLTYPE_DROP_PACKAGE,
    YAC_SQLTYPE_TRUNCATE_TABLE,
    YAC_SQLTYPE_BACKUP_DATABASE,
    YAC_SQLTYPE_RESTORE_DATABASE,
    YAC_SQLTYPE_RECOVER_DATABASE,
    YAC_SQLTYPE_BUILD_DATABASE,
    YAC_SQLTYPE_SET_TRANSACTION,
    YAC_SQLTYPE_PREPARE_TRANSACTION,
    YAC_SQLTYPE_REPLACE_VIEW,
    YAC_SQLTYPE_REPLACE_SYNONYM,
    YAC_SQLTYPE_REPLACE_FUNCTION,
    YAC_SQLTYPE_REPLACE_PROCEDURE,
    YAC_SQLTYPE_REPLACE_PACKAGE,
    YAC_SQLTYPE_FLASHBACK_TABLE,
    YAC_SQLTYPE_COMMENT,
    YAC_SQLTYPE_PURGE_RECYCLEBIN,
    YAC_SQLTYPE_CREATE_ROLE,
    YAC_SQLTYPE_DROP_ROLE,
    YAC_SQLTYPE_GRANT,
    YAC_SQLTYPE_REVOKE,
    YAC_SQLTYPE_ALTER_USER,
    YAC_SQLTYPE_ALTER_TRIGGER,
    YAC_SQLTYPE_DROP_TRIGGER,
    YAC_SQLTYPE_REPLACE_TRIGGER,
    YAC_SQLTYPE_REPLACE_OUTLINE,
    YAC_SQLTYPE_CREATE_AUDIT_POLICY,
    YAC_SQLTYPE_ALTER_AUDIT_POLICY,
    YAC_SQLTYPE_DROP_AUDIT_POLICY,
    YAC_SQLTYPE_AUDIT_POLICY,
    YAC_SQLTYPE_NOAUDIT_POLICY,
    YAC_SQLTYPE_CREATE_OUTLINE,
    YAC_SQLTYPE_ALTER_OUTLINE,
    YAC_SQLTYPE_DROP_OUTLINE,
    YAC_SQLTYPE_ANAYLZE_TABLE,
    YAC_SQLTYPE_DDL_CEIL = 128,

    YAC_SQLTYPE_COMMIT,
    YAC_SQLTYPE_ROLLBACK,
    YAC_SQLTYPE_EXPLAIN,
    YAC_SQLTYPE_SAVEPOINT,
    YAC_SQLTYPE_SHUTDOWN,
    YAC_SQLTYPE_RELEASE_SAVEPOINT,
    YAC_SQLTYPE_SO_ENTITY,
    YAC_SQLTYPE_PACK_ENTITY,

    YAC_SQLTYPE_BATCH_INSERT,
    /* distributed inner type */
    YAC_SQLTYPE_DXG_QUERY,

    __YAC_SQLTYPE_COUNT__ = 255
} YacSQLType;

typedef enum EnYacPumpParam {
    YAC_EXP_FILE = 0,
    YAC_EXP_FULL,
    YAC_EXP_OWNER,
    YAC_EXP_TABLES,

    YAC_IMP_FILE,
    YAC_IMP_FROMUSER,
    YAC_IMP_FULL,
    YAC_IMP_TABLES,
    __YAC_PUMP_PARAM_COUNT__
} YacPumpParam;

typedef enum EnPumpFormat {
    PUMP_FORMAT_BINARY = 0,
} PumpFormat;

YacChar* yacGetVersion();
YacVoid  yacGetLastError(YacInt32* errCode, YacChar** message, YacChar** sqlState, YacTextPos* pos);
YacVoid  yacPrintPadded(const YacChar* str, YacChar padChar, YacInt32 width);

YacResult yacAllocHandle(YacHandleType type, YacHandle input, YacHandle* output);
YacResult yacFreeHandle(YacHandleType type, YacHandle handle);
YacResult yacConnect(YacHandle hConn, const YacChar* url, YacInt16 urlLength, const YacChar* user, YacInt16 userLength,
                     const YacChar* password, YacInt16 passwordLength);
YacVoid   yacDisconnect(YacHandle hConn);
YacResult yacCancel(YacHandle hConn);
YacResult yacDirectExecute(YacHandle hStmt, const YacChar* sql, YacInt32 sqlLength);
YacResult yacPrepare(YacHandle hStmt, const YacChar* sql, YacInt32 sqlLength);
YacResult yacExecute(YacHandle hStmt);
YacResult yacFetch(YacHandle hStmt, YacUint32* rows);
YacResult yacCommit(YacHandle hConn);
YacResult yacRollback(YacHandle hConn);

YacResult yacGetEnvAttr(YacHandle hEnv, YacEnvAttr attr, YacVoid* value, YacInt32 bufLength, YacInt32* stringLength);
YacResult yacSetConnAttr(YacHandle hConn, YacConnAttr attr, YacVoid* value, YacInt32 length);
YacResult yacGetConnAttr(YacHandle hConn, YacConnAttr attr, YacVoid* value, YacInt32 bufLength, YacInt32* stringLength);
YacResult yacSetStmtAttr(YacHandle hStmt, YacStmtAttr attr, YacVoid* value, YacInt32 length);
YacResult yacGetStmtAttr(YacHandle hStmt, YacStmtAttr attr, YacVoid* value, YacInt32 bufLength, YacInt32* stringLength);

YacResult yacDescribeCol2(YacHandle hStmt, YacUint16 id, YacColumnDesc* desc);
YacResult yacBindColumn(YacHandle hStmt, YacUint16 id, YacType type, YacPointer value, YacInt32 bufLen,
                        YacInt32* indicator);
YacResult yacBindParameter(YacHandle hStmt, YacUint16 id, YacParamDirection direction, YacType bindType,
                           YacPointer value, YacUint32 bindSize, YacInt32 bufLength, YacInt32* indicator);
YacResult yacBindParameterByName(YacHandle hStmt, YacChar* name, YacParamDirection direction, YacType bindType,
                                 YacPointer value, YacUint32 bindSize, YacInt32 bufLength, YacInt32* indicator);
YacResult yacNumResultCols(YacHandle hStmt, YacInt16* count);

YacDate   yacNow();
YacVoid   yacNumberFromInt32(YacNumber* n, YacInt32 v);
YacResult yacText2Timestamp(YacChar* text, YacChar* format, YacDate* stamp, YacInt16* bias);
YacResult yacText2YMInterval(YacChar* str, YacYMInterval* interval);
YacResult yacText2DSInterval(YacChar* str, YacDSInterval* interval);
YacResult yacText2ShortTime(YacChar* str, YacChar* format, YacShortTime* shortTime);
YacResult yacGetDateStruct(YacDate date, YacDateStruct* ds);

// multi insert API in "batch" mode
YacResult yacBatchInsertPrepare(YacHandle hStmt, YacChar* tableName);
YacResult yacBatchInsertExecute(YacHandle hStmt);

// lob API
YacResult yacLobDescAlloc(YacHandle hConn, YacType type, YacVoid** desc);
YacResult yacLobDescFree(YacVoid* desc, YacType type);
YacResult yacLobGetChunkSize(YacHandle hConn, YacLobLocator* locator, YacUint16* chunkSize);
YacResult yacLobGetLength(YacHandle hConn, YacLobLocator* locator, YacUint64* length);
YacResult yacLobRead(YacHandle hConn, YacLobLocator* loc, YacUint64* bytes, YacUint8* buf, YacUint64 bufLen);
YacResult yacLobWrite(YacHandle hConn, YacLobLocator* loc, YacUint64* bytes, YacUint8* buf, YacUint64 bufLen);
YacResult yacLobCreateTemporary(YacHandle hConn, YacLobLocator* loc);
YacResult yacLobFreeTemporary(YacHandle hConn, YacLobLocator* loc);

typedef YacResult (*YacServerOutput)(YacHandle hConn, const YacChar* message);

// exp/imp API
typedef YacVoid (*YacPumpMessageCB)(YacVoid* context, const YacChar* message);
YacResult yacExport(YacHandle hPump, const YacChar* url, YacChar* user, const YacChar* password);
YacResult yacImport(YacHandle hPump, const YacChar* url, YacChar* user, const YacChar* password);
YacResult yacSetPumpParam(YacHandle hPump, YacPumpParam paramId, YacVoid* value, YacInt32 len);
YacVoid   yacSetMessageCallback(YacHandle hPump, YacPumpMessageCB cb, YacVoid* context);
YacBool   yacIsPumpWithWarning(YacHandle hPump);

#ifdef __cplusplus
}
#endif
#endif
