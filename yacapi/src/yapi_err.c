#include "yacapi.h"
#include "yapi_inc.h"
#include <stdlib.h>
#include <stdarg.h>
#include <stdio.h>

#ifdef _WIN32
#include <windows.h>
#include <tchar.h>
#else
#include <errno.h>
#include <pthread.h>
#include <sys/time.h>
#include <dlfcn.h>
#endif
#ifdef __linux
#include <unistd.h>
#include <sys/syscall.h>
#endif

#ifdef _MSC_VER
#define __thread __declspec(thread)  // Thread Local Storage
#endif

__thread YapiErrorBuffer gErrorBuf;

void yapiSetError(YapiErrorMsg* error, yapiErrorNum errorNum, const char* format, ...)
{
    if(error == NULL) {
        return;
    }
    va_list args;
    error->buf->code = errorNum;

    va_start(args, format);
    error->buf->messageLen = (uint32_t)vsnprintf(error->buf->message,
                            sizeof(error->buf->message),
                            format, args);
    va_end(args);
}

void yapiInitError(YapiErrorMsg *error)
{
    error->buf = &gErrorBuf;
}

void yapiGetCliError(YapiErrorMsg* error)
{
    yapiCliGetLastError(error);
}

void yapiGetErrorInfo(YapiErrorMsg *error, YapiErrorInfo *info)
{
    info->errCode = error->buf->code;
    info->message = error->buf->message;
    info->pos = &error->buf->pos;
    info->sqlState = error->buf->sqlState;
}
