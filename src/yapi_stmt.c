#include "yapi_inc.h"
#include "stdlib.h"

YapiResult yapiPrepare(YapiConnect* hConn, const char* sql, int32_t sqlLength, YapiStmt** hStmt)
{
    YapiStmt* stmt = malloc(sizeof(YapiStmt));
    if (stmt == NULL) {
        return YAPI_ERROR;
    }
    if (yapiCliAllocHandle(YAPI_HANDLE_STMT, hConn->connHandler, &stmt->stmtHandler) != YAPI_SUCCESS) {
        return YAPI_ERROR;
    }
    if (yapiCliPrepare(stmt->stmtHandler, sql, sqlLength) != YAPI_SUCCESS) {
        return YAPI_ERROR;
    }
    *hStmt = stmt;
    return YAPI_SUCCESS;
}

YapiResult yapiReleaseStmt(YapiStmt* hStmt)
{
    yapiCliFreeHandle(YAPI_HANDLE_STMT, hStmt->stmtHandler);
    return YAPI_SUCCESS;
}

YapiResult yapiExecute(YapiStmt* hStmt)
{
    return yapiCliExecute(hStmt->stmtHandler);
}

YapiResult yapiFetch(YapiStmt* hStmt, uint32_t* rows)
{
    return yapiCliFetch(hStmt->stmtHandler, rows);
}

YapiResult yapiDirectExecute(YapiStmt* hStmt, const char* sql, int32_t sqlLength)
{
    return yapiCliDirectExecute(hStmt->stmtHandler, sql, sqlLength);
}

YapiResult yapiDescribeCol2(YapiStmt* hStmt, uint16_t id, YapiColumnDesc* desc)
{
    return yapiCliDescribeCol2(hStmt->stmtHandler, id, desc);
}

YapiResult yapiBindColumn(YapiStmt* hStmt, uint16_t id, YapiType type, YapiPointer value, int32_t bufLen,
                        int32_t* indicator)
{
    return yapiCliBindColumn(hStmt->stmtHandler, id, type, value, bufLen, indicator);
}

YapiResult yapiBindParameter(YapiStmt* hStmt, uint16_t id, YapiParamDirection direction, YapiType bindType,
                           YapiPointer value, uint32_t bindSize, int32_t bufLength, int32_t* indicator)
{
    return yapiCliBindParameter(hStmt->stmtHandler, id, direction, bindType, value, bindSize, bufLength, indicator);
}

YapiResult yapiBindParameterByName(YapiStmt* hStmt, char* name, YapiParamDirection direction, YapiType bindType,
                                 YapiPointer value, uint32_t bindSize, int32_t bufLength, int32_t* indicator)
{
    return yapiCliBindParameterByName(hStmt->stmtHandler, name, direction, bindType, value, bindSize, bufLength, indicator);
}

YapiResult yapiNumResultCols(YapiStmt* hStmt, int16_t* count)
{
    return yapiCliNumResultCols(hStmt->stmtHandler, count);
}

YapiResult yapiSetStmtAttr(YapiStmt* hStmt, YapiStmtAttr attr, void* value, int32_t length)
{
    return yapiCliSetStmtAttr(hStmt->stmtHandler, attr, value, length);
}

YapiResult yapiGetStmtAttr(YapiStmt* hStmt, YapiStmtAttr attr, void* value, int32_t bufLength, int32_t* stringLength)
{
    return yapiCliGetStmtAttr(hStmt->stmtHandler, attr, value, bufLength, stringLength);
}
