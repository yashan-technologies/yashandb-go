#ifndef ANC_H
#define ANC_H

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef void*               AncHandle;
typedef char                AncInt8;
typedef unsigned char       AncUint8;
typedef short               AncInt16;
typedef unsigned short      AncUint16;
typedef int                 AncInt32;
typedef unsigned int        AncUint32;
typedef int64_t             AncInt64;
typedef uint64_t            AncUint64;
typedef char                AncChar;
typedef void*               AncPointer;
typedef bool                AncBool;
typedef AncInt64            AncDate;
typedef AncInt64            AncShortTime;
typedef AncInt32            AncYMInterval;
typedef AncInt64            AncDSInterval;

typedef enum EnAncResult
{
    ANC_SUCCESS           = 0,
    ANC_SUCCESS_WITH_INFO = 1,
    ANC_ERROR             = -1
}AncResult;

typedef enum EnAncHandleType
{
    ANC_HANDLE_UNKNOWN = 0,
    ANC_HANDLE_ENV     = 1,
    ANC_HANDLE_DBC     = 2,
    ANC_HANDLE_STMT    = 3,
    ANC_HANDLE_DESC    = 4,
}AncHandleType;

typedef struct StAncTextPos
{
    AncInt32 line;
    AncInt32 column;
}AncTextPos;

typedef enum EnAncType
{
    ANC_TYPE_BOOL           = 1,
    ANC_TYPE_TINYINT        = 2,
    ANC_TYPE_SMALLINT       = 3,
    ANC_TYPE_INTEGER        = 4,
    ANC_TYPE_BIGINT         = 5,
    ANC_TYPE_UTINYINT       = 6,
    ANC_TYPE_USMALLINT      = 7,
    ANC_TYPE_UINTEGER       = 8,
    ANC_TYPE_UBIGINT        = 9,
    ANC_TYPE_FLOAT          = 10,
    ANC_TYPE_DOUBLE         = 11,
    ANC_TYPE_NUMBER         = 12,
    ANC_TYPE_DATE           = 13,
    ANC_TYPE_SHORTDATE      = 14,
    ANC_TYPE_SHORTTIME      = 15,
    ANC_TYPE_TIMESTAMP      = 16,
    ANC_TYPE_TIMESTAMP_TZ   = 17,
    ANC_TYPE_TIMESTAMP_LTZ  = 18,
    ANC_TYPE_YM_INTERVAL    = 19,
    ANC_TYPE_DS_INTERVAL    = 20,
    ANC_TYPE_CHAR           = 24,
    ANC_TYPE_NCHAR          = 25,
    ANC_TYPE_VARCHAR        = 26,
    ANC_TYPE_NVARCHAR       = 27,
    ANC_TYPE_BINARY         = 28,
    ANC_TYPE_CLOB           = 29,
    ANC_TYPE_BLOB           = 30,
    ANC_TYPE_BIT            = 31,
    ANC_TYPE_ROWID          = 32,
    ANC_TYPES_COUNT
}AncType;

typedef enum EnAncParamDirection
{
    ANC_PARAM_INPUT   = 0,
    ANC_PARAM_OUTPUT  = 1,
    ANC_PARAM_INOUT   = 2,
}AncParamDirection;

#define ANC_TYPE_SRTING              ANC_TYPE_VARCHAR
#define ANC_PARAM_NAME_BUFFER_SIZE   32
#define ANC_MAX_PARAM_NAME_LEN       30
#define ANC_MAX_ERROR_MSG_LEN        256

typedef enum EnAncConnAttr
{
    ANC_ATTR_LOGONINFO_PTR = 1,
    ANC_ATTR_ASYNC_ENABLE,
    ANC_ATTR_AUTOCOMMIT,
    ANC_ATTR_LOGIN_TIMEOUT,
    ANC_ATTR_STMTS,
    ANC_ATTR_PACKET_SIZE,
    ANC_ATTR_TXN_ISOLATION,
    ANC_ATTR_SERVEROUTPUT,
    ANC_ATTR_NUMWIDTH,
}AncConnAttr;

typedef enum EnAncEnvAttr
{
    ANC_ATTR_DATE_FORMAT = 1,
}AncEnvAttr;

typedef enum EnAncSQLType
{
    ANC_SQLTYPE_SELECT           = 1,
    ANC_SQLTYPE_INSERT           = 2,
    ANC_SQLTYPE_UPDATE           = 3,
    ANC_SQLTYPE_DELETE           = 4,
    ANC_SQLTYPE_MERGE            = 5,
    ANC_SQLTYPE_WITH             = 6,
    ANC_SQLTYPE_ANONYMOUS_BLOCK  = 7,
    ANC_SQLTYPE_DML_CEIL,

    ANC_SQLTYPE_CREATE_DATABASE  = 10,
    ANC_SQLTYPE_CREATE_TABLE     = 11,
    ANC_SQLTYPE_CREATE_INDEX     = 12,
    ANC_SQLTYPE_CREATE_VIEW      = 13,
    ANC_SQLTYPE_CREATE_PROCEDURE = 14,
    ANC_SQLTYPE_CREATE_FUNCTION  = 15,
    ANC_SQLTYPE_ALTER_TABLE      = 16,

    ANC_SQLTYPE_ALTER_SESSION    = 34,
}AncSQLType;

typedef enum EnAncStmtAttr
{
    ANC_ATTR_PARAMSET_SIZE = 100,
    ANC_ATTR_ROWSET_SIZE,
    ANC_ATTR_ROWS_FETECHED,
    ANC_ATTR_ROWS_AFFECTED,
    ANC_ATTR_CURSOR_EOF,
    ANC_ATTR_SQLTYPE,
    ANC_ATTR_IS_BATCHROWS,
    ANC_ATTR_IS_BATCH_ERRORS,
    ANC_ATTR_ACK_BATCHROWS_SIZE,
    ANC_ATTR_ACK_BATCH_ERRORS_SIZE,
    ANC_ATTR_ACK_BATCHROWS,
    ANC_ATTR_ACK_BATCH_ERRORS,
}AncStmtAttr;

typedef struct StAncColumnDesc
{
    AncChar*    name;
    AncUint32   size;
    AncUint8    type;
    AncUint8    precision;
    AncInt8     scale;
    AncUint8    nullable;
}AncColumnDesc;

typedef struct StAncBatchError {
    AncUint32  rowNum;
    AncUint32  errCode;
    AncUint32  msgLen;
    AncChar    msg[ANC_MAX_ERROR_MSG_LEN];
} AncBatchError;

#pragma pack(4)
typedef struct StAncNumber {
    AncUint64 item[2];
    AncInt8   sign;
    AncUint8  unused;
    AncInt16  exp;
} AncNumber;

typedef struct StAncTimestamp
{
    AncInt64 stamp;
    AncInt16 bias;   //minutes
    AncInt16 unused;
}AncTimestamp;
#pragma pack()



typedef AncResult (*AncServerOutput)(AncHandle hConn, const AncChar* message);

#define ANC_NULL_DATA    ((AncInt32)-1)

AncChar*  ancGetVersion();
void      ancGetLastError(AncInt32* errCode, AncChar** message, AncChar** sqlState, AncTextPos* pos);
AncResult ancAllocHandle(AncHandleType type, AncHandle input, AncHandle* output);
AncResult ancFreeHandle(AncHandleType type, AncHandle handle);
AncResult ancConnect(AncHandle hConn, const AncChar* url, const AncChar* user, const AncChar* password);
AncResult ancCancel(AncHandle hConn);
void      ancDisconnect(AncHandle hConn);
AncResult ancPrepare(AncHandle hStmt, const AncChar* sql);
AncResult ancExecute(AncHandle hStmt);
AncResult ancDirectExecute(AncHandle hStmt, const AncChar* sql);
AncResult ancFetch(AncHandle hStmt, AncUint32* rows);
AncResult ancCommit(AncHandle hConn);
AncResult ancRollback(AncHandle hConn);
AncResult ancSetEnvAttr(AncHandle hEnv, AncEnvAttr attr, void* value, AncInt32 len);
AncResult ancGetEnvAttr(AncHandle hEnv, AncEnvAttr attr, void* value, AncInt32 bufLen);
AncResult ancSetConnAttr(AncHandle hConn, AncConnAttr attr, void* value, AncInt32 len);
AncResult ancGetConnAttr(AncHandle hConn, AncConnAttr attr, void* value, AncInt32 bufLen);
AncResult ancSetStmtAttr(AncHandle hStmt, AncStmtAttr attr, void* value, AncInt32 len);
AncResult ancGetStmtAttr(AncHandle hStmt, AncStmtAttr attr, void* value, AncInt32 bufLen);
AncResult ancNumResultCols(AncHandle hStmt, AncInt16* count);
AncResult ancDescribeCol(AncHandle hStmt, AncUint16 id, AncChar* name, AncInt16 bufLen, AncInt16* nameLen,
    AncInt16* dataType, AncUint32* size, AncInt16* numDigits, AncInt16* nullable);
AncResult ancDescribeCol2(AncHandle hStmt, AncUint16 id, AncColumnDesc* desc);
AncResult ancBindColumn(AncHandle hStmt, AncUint16 id,
    AncType type, AncPointer value, AncInt32 bufLen, AncInt32* indicator);
AncResult ancBindParameter(AncHandle hStmt, AncUint16 id,
    AncParamDirection direction, AncType type, AncPointer value, AncInt32 size, AncInt32* indicator);
void ancPrintPadded(const AncChar* str, AncChar padChar, AncInt32 width);
void ancNumberFromInt32(AncNumber* n, AncInt32 v);
AncDate   ancNow();
AncResult ancText2Timestamp(AncChar* text, AncChar* format, AncDate* stamp, AncInt16* bias);
AncResult ancText2YMInterval(AncChar* str, AncYMInterval* interval);
AncResult ancText2DSInterval(AncChar* str, AncDSInterval* interval);
AncResult ancText2ShortTime(AncChar* str, AncChar* format, AncShortTime* shortTime);

#define ANC_CALL(proc)                           \
    do {                                         \
        if ((AncResult)(proc) != ANC_SUCCESS) { \
            return ANC_ERROR;                    \
        }                                        \
    } while (0)

#ifdef __cplusplus
}
#endif
#endif