#include "yapi_inc.h"
#include "stdlib.h"

YapiResult yapiPdbgStart(YapiStmt* hStmt)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgStart(hStmt->stmtHandler, &error);
}

YapiResult yapiPdbgCheckVersion(YapiStmt* hStmt, uint64_t objId, uint16_t subId, uint32_t version)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgCheckVersion(hStmt->stmtHandler, objId, subId, version, &error);
}

YapiResult yapiPdbgAbort(YapiStmt* hStmt)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgAbort(hStmt->stmtHandler, &error);
}

YapiResult yapiPdbgContinue(YapiStmt* hStmt)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgContinue(hStmt->stmtHandler, &error);
}

YapiResult yapiPdbgStepInto(YapiStmt* hStmt)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgStepInto(hStmt->stmtHandler, &error);
}

YapiResult yapiPdbgStepOut(YapiStmt* hStmt)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgStepOut(hStmt->stmtHandler, &error);
}

YapiResult yapiPdbgStepNext(YapiStmt* hStmt)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgStepNext(hStmt->stmtHandler, &error);
}

YapiResult yapiPdbgDeleteAllBreakpoints(YapiStmt* hStmt)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgDeleteAllBreakpoints(hStmt->stmtHandler, &error);
}

YapiResult yapiPdbgAddBreakpoint(YapiStmt* hStmt, uint64_t objId, uint16_t subId, uint32_t lineNo, uint32_t* bpId)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgAddBreakpoint(hStmt->stmtHandler, objId, subId, lineNo, bpId, &error);
}

YapiResult yapiPdbgDeleteBreakpoint(YapiStmt* hStmt, uint64_t objId, uint16_t subId, uint32_t lineNo)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgDeleteBreakpoint(hStmt->stmtHandler, objId, subId, lineNo, &error);
}

YapiResult yapiPdbgGetBreakpointsCount(YapiStmt* hStmt, uint32_t* bpCount)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgGetBreakpointsCount(hStmt->stmtHandler, bpCount, &error);
}

YapiResult yapiPdbgGetAllVars(YapiStmt* hStmt, uint32_t* bpCount)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgGetAllVars(hStmt->stmtHandler, bpCount, &error);
}

YapiResult yapiPdbgGetAllFrames(YapiStmt* hStmt, uint32_t* frameCount)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgGetAllFrames(hStmt->stmtHandler, frameCount, &error);
}

YapiResult yapiPdbgGetRunningData(YapiStmt* hStmt, YapiDebugRunningAttr attr, void* value, int32_t bufLen)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgGetRunningData(hStmt->stmtHandler, attr, value, bufLen, &error);
}

YapiResult yapiPdbgGetFrameData(YapiStmt* hStmt, uint32_t id, YapiDebugFrameAttr attr, void* value, int32_t bufLen)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgGetFrameData(hStmt->stmtHandler, id, attr, value, bufLen, &error);
}

YapiResult yapiPdbgGetVarData(YapiStmt* hStmt, uint32_t id, YapiDebugVarAttr attr, void* value, int32_t bufLen)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgGetVarData(hStmt->stmtHandler, id, attr, value, bufLen, &error);
}

YapiResult yapiPdbgGetVarValue(YapiStmt* hStmt, uint32_t id, uint32_t valueType, void* value, int32_t bufLen,
                               int32_t* indicator)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgGetVarValue(hStmt->stmtHandler, id, valueType, value, bufLen, indicator, &error);
}

YapiResult yapiPdbgGetBreakpointData(YapiStmt* hStmt, uint32_t id, YapiDebugBpAttr attr, void* value, int32_t bufLen)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgGetBreakpointData(hStmt->stmtHandler, id, attr, value, bufLen, &error);
}