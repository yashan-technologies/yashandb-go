#include "yapi_inc.h"
#include "stdlib.h"
#include "inttypes.h"

#ifdef _WIN32
#define YACLI_LIB_NAME "yascli.dll"
#else
#define YACLI_LIB_NAME "libyascli.so"
#endif

YapiResult yapiAllocEnv(YapiEnv** inst)
{
    YapiErrorMsg error;

    void* handle;
    yapiInitError(&error);
    if (yapiOpenDynamicLib(YACLI_LIB_NAME, &handle, &error) == YAPI_ERROR) {
        return YAPI_ERROR;
    }

    YapiEnv* env;
    if (yapiAllocMem("Environment", 1, sizeof(YapiEnv), (void**)&env, &error) != YAPI_SUCCESS) {
        return YAPI_ERROR;
    }
    if (yapiCliAllocHandle(YAPI_HANDLE_ENV, NULL, &env->envHandler, &error) == YAPI_ERROR) {
        yapiFreeMem(env);
        return YAPI_ERROR;
    }

    *inst = env;
    return YAPI_SUCCESS;
}

YapiResult yapiReleaseEnv(YapiEnv* inst)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    
    YAPI_CALL(yapiCliFreeHandle(YAPI_HANDLE_ENV, inst->envHandler, &error));
    yapiFreeMem(inst);

    return YAPI_SUCCESS;
}

YapiResult yapiEnvGetAttr(YapiEnv* hEnv, YapiEnvAttr attr, void* value, int32_t bufLength, int32_t* stringLength) 
{
    YapiErrorMsg error;
    yapiInitError(&error);

    return yapiCliGetEnvAttr(hEnv->envHandler, attr, value, bufLength, stringLength, &error);
}

char* yapiGetVersion(YapiEnv* inst) 
{
    YapiErrorMsg error;
    yapiInitError(&error);

    char* version = NULL;
    yapiCliGetVersion(&version, &error);
    return version;
}