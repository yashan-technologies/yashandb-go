#include "yapi_inc.h"
#include "stdlib.h"
#include "string.h"

YapiResult yapiConnect(YapiEnv* env, const char* url, int16_t urlLength, const char* user, int16_t userLength,
                       const char* password, int16_t passwordLength, YapiConnect** hConn)
{
    YapiErrorMsg error;
    yapiInitError(&error);

    YapiConnect* conn;
    if (yapiAllocMem("Connection", 1, sizeof(YapiConnect), (void**)&conn, &error) != YAPI_SUCCESS) {
        return YAPI_ERROR;
    }
    if (yapiCliAllocHandle(YAPI_HANDLE_DBC, env->envHandler, &conn->connHandler, &error) != YAPI_SUCCESS) {
        yapiFreeMem(conn);
        return YAPI_ERROR;
    }
    if (yapiCliConnect(conn->connHandler, url, urlLength, user, userLength, password, passwordLength, &error) !=
        YAPI_SUCCESS) {
        yapiCliFreeHandle(YAPI_HANDLE_DBC, conn->connHandler, &error);
        yapiFreeMem(conn);
        return YAPI_ERROR;
    }
    *hConn = conn;
    return YAPI_SUCCESS;
}

YapiResult yapiDisconnect(YapiConnect* hConn)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliDisconnect(hConn->connHandler, &error);
}

YapiResult yapiReleaseConn(YapiConnect* hConn)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliFreeHandle(YAPI_HANDLE_DBC, hConn->connHandler, &error);
}

YapiResult yapiCancel(YapiConnect* hConn)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return YAPI_ERROR;
}

YapiResult yapiCommit(YapiConnect* hConn)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliCommit(hConn->connHandler, &error);
}

YapiResult yapiRollback(YapiConnect* hConn)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliRollback(hConn->connHandler, &error);
}

YapiResult yapiSetConnAttr(YapiConnect* hConn, YapiConnAttr attr, void* value, int32_t length)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliSetConnAttr(hConn->connHandler, attr, value, length, &error);
}

YapiResult yapiGetConnAttr(YapiConnect* hConn, YapiConnAttr attr, void* value, int32_t bufLength, int32_t* stringLength)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliGetConnAttr(hConn->connHandler, attr, value, bufLength, stringLength, &error);
}

void yapiGetLastError(YapiErrorInfo* info)
{
    YapiErrorMsg error;

    yapiInitError(&error);
    yapiGetErrorInfo(&error, info);
}

YapiResult yapiLobDescAlloc(YapiConnect* hConn, YapiType type, void** desc)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliLobDescAlloc(hConn->connHandler, type, desc, &error);
}

YapiResult yapiLobDescFree(void* desc, YapiType type)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliLobDescFree(desc, type, &error);
}

YapiResult yapiLobGetChunkSize(YapiConnect* hConn, YapiLobLocator* locator, uint16_t* chunkSize)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliLobGetChunkSize(hConn->connHandler, locator, chunkSize, &error);
}

YapiResult yapiLobGetLength(YapiConnect* hConn, YapiLobLocator* locator, uint64_t* length)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliLobGetLength(hConn->connHandler, locator, length, &error);
}

YapiResult yapiLobRead(YapiConnect* hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf, uint64_t bufLen)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliLobRead(hConn->connHandler, loc, bytes, buf, bufLen, &error);
}

YapiResult yapiLobWrite(YapiConnect* hConn, YapiLobLocator* loc, uint64_t* bytes, uint8_t* buf, uint64_t bufLen)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliLobWrite(hConn->connHandler, loc, bytes, buf, bufLen, &error);
}

YapiResult yapiLobCreateTemporary(YapiConnect* hConn, YapiLobLocator* loc)
{
    YapiErrorMsg error;
    yapiInitError(&error);

    return yapiCliLobCreateTemporary(hConn->connHandler, loc, &error);
}

YapiResult yapiLobFreeTemporary(YapiConnect* hConn, YapiLobLocator* loc)
{
    YapiErrorMsg error;
    yapiInitError(&error);

    return yapiCliLobFreeTemporary(hConn->connHandler, loc, &error);
}
