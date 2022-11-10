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

YapiResult yapiNumberRound(YapiNumber* n, int32_t precision, int32_t scale)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCliNumberRound(n, precision, scale, &error);
}