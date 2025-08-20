#include "yapi_inc.h"

YapiResult yapiGetDateStruct(YapiDate date, YapiDateStruct* ds)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliGetDateStruct(date, ds, &error);
}

YapiResult yapiDateGetDate(const YapiDate date, int16_t* year, uint8_t* month, uint8_t* day)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliDateGetDate(date, year, month, day, &error);
}

YapiResult yapiShortTimeGetShortTime(const YapiShortTime time, uint8_t* hour, uint8_t* minute, uint8_t* second,
                                     uint32_t* fraction)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliShortTimeGetShortTime(time, hour, minute, second, fraction, &error);
}

YapiResult yapiTimestampGetTimestamp(const YapiTimestamp timestamp, int16_t* year, uint8_t* month, uint8_t* day,
                                     uint8_t* hour, uint8_t* minute, uint8_t* second, uint32_t* fraction)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliTimestampGetTimestamp(timestamp, year, month, day, hour, minute, second, fraction, &error);
}

YapiResult yapiYMIntervalGetYearMonth(const YapiYMInterval ymInterval, int32_t* year, int32_t* month)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliYMIntervalGetYearMonth(ymInterval, year, month, &error);
}

YapiResult yapiDSIntervalGetDaySecond(const YapiDSInterval dsInterval, int32_t* day, int32_t* hour, int32_t* minute,
                                      int32_t* second, int32_t* fraction)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliDSIntervalGetDaySecond(dsInterval, day, hour, minute, second, fraction, &error);
}

YapiResult yapiDateSetDate(YapiDate* date, int16_t year, uint8_t month, uint8_t day)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliDateSetDate(date, year, month, day, &error);
}

YapiResult yapiShortTimeSetShortTime(YapiShortTime* time, uint8_t hour, uint8_t minute, uint8_t second,
                                     uint32_t fraction)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliShortTimeSetShortTime(time, hour, minute, second, fraction, &error);
}

YapiResult yapiTimestampSetTimestamp(YapiTimestamp* timestamp, int16_t year, uint8_t month, uint8_t day, uint8_t hour,
                                     uint8_t minute, uint8_t second, uint32_t fraction)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliTimestampSetTimestamp(timestamp, year, month, day, hour, minute, second, fraction, &error);
}

YapiResult yapiDateTimeGetTimeZoneOffset(YapiEnv* hEnv, YapiTimestamp timestamp, int8_t* hr, int8_t* mm)

{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliDateTimeGetTimeZoneOffset(hEnv->envHandler, timestamp, hr, mm, &error);
}

YapiResult yapiYMIntervalSetYearMonth(YapiYMInterval* ymInterval, int32_t year, int32_t month)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliYMIntervalSetYearMonth(ymInterval, year, month, &error);
}

YapiResult yapiDSIntervalSetDaySecond(YapiDSInterval* dsInterval, int32_t day, int32_t hour, int32_t minute,
                                      int32_t second, int32_t fraction)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliDSIntervalSetDaySecond(dsInterval, day, hour, minute, second, fraction, &error);
}

YapiResult yapiDSIntervalFromText(YapiEnv* hEnv, YapiDSInterval* dsInterval, const char* str, uint32_t strLen)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliDSIntervalFromText(hEnv->envHandler, dsInterval, str, strLen, &error);
}

YapiResult yapiYMIntervalFromText(YapiEnv* hEnv, YapiYMInterval* ymInterval, const char* str, uint32_t strLen)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliYMIntervalFromText(hEnv->envHandler, ymInterval, str, strLen, &error);
}

YapiResult yapiNumberRound(YapiNumber* n, int32_t precision, int32_t scale)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliNumberRound(n, precision, scale, &error);
}

YapiResult yapiNumberToText(const YapiNumber* number, const char* fmt, uint32_t fmtLength, const char* nlsParam,
                            uint32_t nlsParamLength, char* str, int32_t bufLength, int32_t* length)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliNumberToText(number, fmt, fmtLength, nlsParam, nlsParamLength, str, bufLength, length, &error);
}

YapiResult yapiNumberFromText(const char* str, uint32_t strLength, const char* fmt, uint32_t fmtLength,
                              const char* nlsParam, uint32_t nlsParamLength, YapiNumber* number)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliNumberFromText(str, strLength, fmt, fmtLength, nlsParam, nlsParamLength, number, &error);
}

YapiResult yapiNumberFromReal(const YapiPointer rnum, uint32_t length, YapiNumber* number)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliNumberFromReal(rnum, length, number, &error);
}

YapiResult yapiNumberToReal(const YapiNumber* number, uint32_t length, YapiPointer rsl)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliNumberToReal(number, length, rsl, &error);
}