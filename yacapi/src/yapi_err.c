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

void yapiSetError(YapiErrorMsg* error, yapiErrorNum errorNum, const char* format, ...)
{
    if(error == NULL) {
        return;
    }
    va_list args;
    error->code = errorNum;

    va_start(args, format);
    error->messageLen = (uint32_t)vsnprintf(error->message,
                            sizeof(error->message),
                            format, args);
    va_end(args);
}

#ifdef _WIN32
static YapiResult yapiGetWindowsError(DWORD errNum, YapiErrorMsg* error)
{
    TCHAR *errBuf = NULL;
    DWORD status = FormatMessage(FORMAT_MESSAGE_FROM_SYSTEM |
            FORMAT_MESSAGE_IGNORE_INSERTS | FORMAT_MESSAGE_ALLOCATE_BUFFER,
            NULL, errNum, MAKELANGID(LANG_ENGLISH, SUBLANG_ENGLISH_US),
            (LPTSTR) &errBuf, 0, NULL);
    if (errBuf == NULL) {
        return YAPI_SUCCESS;
    }

    size_t len = _tcslen(errBuf);
    if (len >= T2S_BUFFER_SIZE) {
        len = T2S_BUFFER_SIZE - 1;
    }
    memcpy(error->message, errBuf, len);
    LocalFree(errBuf);
    return YAPI_ERROR;
}
#endif