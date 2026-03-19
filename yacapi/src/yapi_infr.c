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

//-----------------------------------------------------------------------------
// Column Desc Function
//-----------------------------------------------------------------------------
uint8_t yapiColumnDescGetPrecision(const YapiColumnDesc* desc)
{
    return desc->precision;
}

void yapiColumnDescSetPrecision(YapiColumnDesc* desc, uint8_t precision)
{
    desc->precision = precision;
}

int8_t yapiColumnDescGetScale(const YapiColumnDesc* desc)
{
    return desc->scale;
}

void yapiColumnDescSetScale(YapiColumnDesc* desc, int8_t scale)
{
    desc->scale = scale;
}

uint8_t yapiColumnDescGetVectorFormat(const YapiColumnDesc* desc)
{
    return desc->vector.format;
}

void yapiColumnDescSetVectorFormat(YapiColumnDesc* desc, uint8_t format)
{
    desc->vector.format = format;
}
