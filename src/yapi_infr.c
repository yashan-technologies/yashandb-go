#include "yapi_inc.h"
#include <inttypes.h>

YapiResult yapiAllocMem(const char* name, size_t numMembers, size_t memberSize, void **ptr, YapiErrorMsg *error)
{
    size_t size = numMembers * memberSize;
    *ptr = malloc(size);
    if (!*ptr){
        yapiSetError(error, YAPI_ERR_ALLOC_MEM, "cannot allocate %" PRId64 " bytes for %s", size, name);
        return YAPI_ERROR;
    }
    return YAPI_SUCCESS;
}

//-----------------------------------------------------------------------------
void yapiFreeMem(void *ptr)
{
    if (ptr == NULL) {
        return;
    }
    free(ptr);
}
