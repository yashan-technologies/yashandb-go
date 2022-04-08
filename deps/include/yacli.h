#ifndef YACLIC_H
#define YACLIC_H

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef void*               YacHandle;
typedef char                YacInt8;
typedef unsigned char       YacUint8;
typedef short               YacInt16;
typedef unsigned short      YacUint16;
typedef int                 YacInt32;
typedef unsigned int        YacUint32;
typedef int64_t             YacInt64;
typedef uint64_t            YacUint64;
typedef char                YacChar;
typedef void*               YacPointer;
typedef bool                YacBool;
typedef YacInt64            YacDate;
typedef YacInt64            YacShortTime;
typedef YacInt32            YacYMInterval;
typedef YacInt64            YacDSInterval;


typedef enum EnYacResult
{
    YAC_SUCCESS           = 0,
    YAC_SUCCESS_WITH_INFO = 1,
    YAC_ERROR             = -1
} YacResult;

typedef enum EnYacHandleType
{
    YAC_HANDLE_UNKNOWN = 0,
    YAC_HANDLE_ENV     = 1,
    YAC_HANDLE_DBC     = 2,
    YAC_HANDLE_STMT    = 3,
    YAC_HANDLE_DESC    = 4,
    __YAC_HANDLE_COUNT__
} YacHandleType;

typedef struct StYacTextPos
{
    YacInt32 line;
    YacInt32 column;
}YacTextPos;

typedef enum EnYacType
{
    YAC_TYPE_BOOL           = 1,
    YAC_TYPE_TINYINT        = 2,
    YAC_TYPE_SMALLINT       = 3,
    YAC_TYPE_INTEGER        = 4,
    YAC_TYPE_BIGINT         = 5,
    YAC_TYPE_UTINYINT       = 6,
    YAC_TYPE_USMALLINT      = 7,
    YAC_TYPE_UINTEGER       = 8,
    YAC_TYPE_UBIGINT        = 9,
    YAC_TYPE_FLOAT          = 10,
    YAC_TYPE_DOUBLE         = 11,
    YAC_TYPE_NUMBER         = 12,
    YAC_TYPE_DATE           = 13,
    YAC_TYPE_SHORTDATE      = 14,
    YAC_TYPE_SHORTTIME      = 15,
    YAC_TYPE_TIMESTAMP      = 16,
    YAC_TYPE_TIMESTAMP_TZ   = 17,
    YAC_TYPE_TIMESTAMP_LTZ  = 18,
    YAC_TYPE_YM_INTERVAL    = 19,
    YAC_TYPE_DS_INTERVAL    = 20,
    YAC_TYPE_CHAR           = 24,
    YAC_TYPE_NCHAR          = 25,
    YAC_TYPE_VARCHAR        = 26,
    YAC_TYPE_NVARCHAR       = 27,
    YAC_TYPE_BINARY         = 28,
    YAC_TYPE_CLOB           = 29,
    YAC_TYPE_BLOB           = 30,
    YAC_TYPE_BIT            = 31,
    YAC_TYPE_ROWID          = 32,
    YAC_TYPE_NCLOB          = 33,
    YAC_TYPE_CURSOR         = 34,
    YAC_TYPES_COUNT
}YacType;

typedef enum EnYacParamDirection
{
    YAC_PARAM_INPUT   = 1,
    YAC_PARAM_OUTPUT  = 2,
    YAC_PARAM_INOUT   = 3,
}YacParamDirection;

#define YAC_TYPE_SRTING              YAC_TYPE_VARCHAR
#define YAC_PARAM_NAME_BUFFER_SIZE   32
#define YAC_MAX_PARAM_NAME_LEN       30
#define YAC_MAX_ERROR_MSG_LEN        256

//should be consistent with ANR_MIN_PACKET_SIZE
#define YAC_MIN_PACKET_SIZE          KB(64)
//should be consistent with ANR_MAX_PACKET_SIZE
#define YAC_MAX_PACKET_SIZE          MB(32)

typedef enum EnYacConnAttr
{
    YAC_ATTR_LOGONINFO_PTR = 1,
    YAC_ATTR_ASYNC_ENABLE,
    YAC_ATTR_AUTOCOMMIT,
    YAC_ATTR_LOGIN_TIMEOUT,
    YAC_ATTR_STMTS,
    YAC_ATTR_PACKET_SIZE,
    YAC_ATTR_TXN_ISOLATION,
    YAC_ATTR_SERVEROUTPUT,
    YAC_ATTR_NUMWIDTH,
}YacConnAttr;

typedef enum EnYacEnvAttr
{
    YAC_ATTR_DATE_FORMAT = 1,
}YacEnvAttr;

typedef enum EnYacSQLType
{
    YAC_SQLTYPE_QUERY = 1,
    YAC_SQLTYPE_INSERT,
    YAC_SQLTYPE_UPDATE,
    YAC_SQLTYPE_DELETE,
    YAC_SQLTYPE_MERGE,
    YAC_SQLTYPE_WITH,
    YAC_SQLTYPE_ANONYMOUS_BLOCK,
    YAC_SQLTYPE_DML_CEIL = 10,

    YAC_SQLTYPE_CREATE_DATABASE = 11,
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
    YAC_SQLTYPE_ALTER_TABLESPACE,
    YAC_SQLTYPE_ALTER_SEQUENCE,
    YAC_SQLTYPE_ALTER_SYSTEM,
    YAC_SQLTYPE_ALTER_SESSION,
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
    YAC_SQLTYPE_DDL_CEIL = 128,

    YAC_SQLTYPE_COMMIT,
    YAC_SQLTYPE_ROLLBACK,
    YAC_SQLTYPE_EXPLAIN,
    YAC_SQLTYPE_SAVEPOINT,
    YAC_SQLTYPE_SHUTDOWN,
    YAC_SQLTYPE_SO_ENTITY,

    /* distributed inner type */
    YAC_SQLTYPE_DXG_QUERY,

    __YAC_SQLTYPE_COUNT__ = 255
}YacSQLType;

typedef enum EnYacStmtAttr
{
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
}YacStmtAttr;

typedef struct StYacColumnDesc
{
    YacChar*    name;
    YacUint32   size;
    YacUint8    type;
    YacUint8    precision;
    YacInt8     scale;
    YacUint8    nullable;
}YacColumnDesc;

typedef struct StYacBatchError {
    YacUint32  rowNum;
    YacUint32  errCode;
    YacUint32  msgLen;
    YacChar    msg[YAC_MAX_ERROR_MSG_LEN];
} YacBatchError;

#pragma pack(4)
typedef struct StYacNumber {
    YacUint64 item[2];
    YacInt8   sign;
    YacUint8  unused;
    YacInt16  exp;
} YacNumber;

typedef struct StYacTimestamp
{
    YacInt64 stamp;
    YacInt16 bias;   //minutes
    YacInt16 unused;
}YacTimestamp;
#pragma pack()

typedef struct {
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

typedef struct StYacLobLocator YacLobLocator;


typedef YacResult (*YacServerOutput)(YacHandle hConn, const YacChar* message);

#define YAC_NULL_DATA    ((YacInt32)-1)

YacChar*  yacGetVersion();
void      yacGetLastError(YacInt32* errCode, YacChar** message, YacChar** sqlState, YacTextPos* pos);
YacResult yacAllocHandle(YacHandleType type, YacHandle input, YacHandle* output);
YacResult yacFreeHandle(YacHandleType type, YacHandle handle);
YacResult yacConnect(YacHandle hConn, const YacChar* url, const YacChar* user, const YacChar* password);
YacResult yacCyacel(YacHandle hConn);
void      yacDisconnect(YacHandle hConn);
YacResult yacPrepare(YacHandle hStmt, const YacChar* sql);
YacResult yacExecute(YacHandle hStmt);
YacResult yacDirectExecute(YacHandle hStmt, const YacChar* sql);
YacResult yacFetch(YacHandle hStmt, YacUint32* rows);
YacResult yacCommit(YacHandle hConn);
YacResult yacRollback(YacHandle hConn);
YacResult yacSetEnvAttr(YacHandle hEnv, YacEnvAttr attr, void* value, YacInt32 len);
YacResult yacGetEnvAttr(YacHandle hEnv, YacEnvAttr attr, void* value, YacInt32 bufLen);
YacResult yacSetConnAttr(YacHandle hConn, YacConnAttr attr, void* value, YacInt32 len);
YacResult yacGetConnAttr(YacHandle hConn, YacConnAttr attr, void* value, YacInt32 bufLen);
YacResult yacSetStmtAttr(YacHandle hStmt, YacStmtAttr attr, void* value, YacInt32 len);
YacResult yacGetStmtAttr(YacHandle hStmt, YacStmtAttr attr, void* value, YacInt32 bufLen);
YacResult yacNumResultCols(YacHandle hStmt, YacInt16* count);
YacResult yacDescribeCol(YacHandle hStmt, YacUint16 id, YacChar* name, YacInt16 bufLen, YacInt16* nameLen,
    YacInt16* dataType, YacUint32* size, YacInt16* numDigits, YacInt16* nullable);
YacResult yacDescribeCol2(YacHandle hStmt, YacUint16 id, YacColumnDesc* desc);
YacResult yacBindColumn(YacHandle hStmt, YacUint16 id,
    YacType type, YacPointer value, YacInt32 bufLen, YacInt32* indicator);
YacResult yacBindParameter(YacHandle hStmt, YacUint16 id,
    YacParamDirection direction, YacType type, YacPointer value, YacInt32 size, YacInt32* indicator);
void yacPrintPadded(const YacChar* str, YacChar padChar, YacInt32 width);
void yacNumberFromInt32(YacNumber* n, YacInt32 v);
YacDate   yacNow();
YacResult yacText2Timestamp(YacChar* text, YacChar* format, YacDate* stamp, YacInt16* bias);
YacResult yacText2YMInterval(YacChar* str, YacYMInterval* interval);
YacResult yacText2DSInterval(YacChar* str, YacDSInterval* interval);
YacResult yacText2ShortTime(YacChar* str, YacChar* format, YacShortTime* shortTime);
YacResult yacGetDateStruct(YacDate date, YacDateStruct* ds);

#define YAC_TRUE true
#define YAC_FALSE false

#define YAC_MAX_COL_LEN 1024
#define YAC_MAX_SQL_LEN 4096
YacResult yacBatchInsertPrepare(YacHandle hStmt, YacChar* tableName);
YacResult yacBatchInsertExecute(YacHandle hStmt);


/* lob */
YacResult yacDescriptAlloc(YacHandle hConn, YacType type, void** desc);
YacResult yacDescriptFree(void* desc, YacType type);
YacResult yacLobGetChunkSize(YacHandle hConn, YacLobLocator* locator, YacUint16* chunkSize);
YacResult yacLobGetLength(YacHandle hConn, YacLobLocator* locator, YacUint64* length);
YacResult yacLobRead(YacHandle hConn, YacLobLocator* locator, YacUint64* bytes, YacUint64 offset, YacUint8* buf,
                     YacUint64 bufLen);
YacResult yacLobWrite(YacHandle hConn, YacLobLocator* locator, YacUint64* bytes, YacUint64 offset, YacUint8* buf,
                      YacUint64 bufLen);

/*!
 *
 * @param hConn the handle of connection
 * @param locator locatot of lob
 * @param lobtype CLOB BLOB or NCLOB
 * @return
 */
YacResult yacLobCreateTemporary(YacHandle* hConn, YacLobLocator* locator, YacUint8 lobtype);
YacResult yacLobFreeTemporary(YacHandle hConn, YacLobLocator* locator);

#define YAC_CALL(proc)                           \
    do {                                         \
        if ((YacResult)(proc) != YAC_SUCCESS) { \
            return YAC_ERROR;                    \
        }                                        \
    } while (0)

#ifdef __cplusplus
}
#endif
#endif