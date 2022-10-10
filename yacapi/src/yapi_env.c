#include "yapi_inc.h"
#include "stdlib.h"

#ifdef _WIN32
#define YACLI_LIB_NAME "yascli.dll"
#else
#define YACLI_LIB_NAME "libyascli.so"
#endif

YapiResult yapiAllocEnv(YapiEnv** inst)
{
    YapiErrorMsg error;

    void* handle;
    if (yapiOpenDynamicLib(YACLI_LIB_NAME, &handle, &error) == YAPI_ERROR) {
        return YAPI_ERROR;
    }

    YapiEnv* env = malloc(sizeof(YapiEnv));
    if (env == NULL) {
        return YAPI_ERROR;
    }
    if (yapiCliAllocHandle(YAPI_HANDLE_ENV, NULL, &env->envHandler) == YAPI_ERROR) {
        return YAPI_ERROR;
    }

    *inst = env;
    return YAPI_SUCCESS;
}

YapiResult yapiReleaseEnv(YapiEnv* inst)
{
    return yapiCliFreeHandle(YAPI_HANDLE_ENV, inst->envHandler);
}