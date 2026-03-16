#include "yapi_inc.h"

YapiResult yapiVectorFromText(YapiVector* vector, YapiVectorFormat format, uint16_t dim, char* text, uint32_t textlen, uint32_t mode)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliVectorFromText(vector, format, dim, text, textlen, mode, &error);
}

YapiResult yapiVectorFromArray(YapiVector* vector, YapiVectorFormat format, uint16_t dim, uint8_t* array, uint32_t arrayLen, uint32_t mode)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliVectorFromArray(vector, format, dim, array, arrayLen, mode, &error);
}

YapiResult yapiVectorToText(YapiVector* vector, char* text, uint32_t* textlen, uint32_t mode)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliVectorToText(vector, text, textlen, mode, &error);
}

YapiResult yapiVectorToArray(YapiVector* vector, YapiVectorFormat format, uint16_t* dim, uint8_t* array, uint32_t* arrayLen, uint32_t mode)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliVectorToArray(vector, format, dim, array, arrayLen, mode, &error);
}

YapiResult yapiVectorGetFormat(YapiVector* vector, YapiVectorFormat* format)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliVectorGetFormat(vector, format, &error);
}

YapiResult yapiVectorGetDimension(YapiVector* vector, uint16_t* dim)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliVectorGetDimension(vector, dim, &error);
}
