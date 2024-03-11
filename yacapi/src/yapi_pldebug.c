#include "yapi_inc.h"
#include "stdlib.h"

YapiResult yapiPdbgStart(YapiStmt* hStmt, char* procName, uint32_t procNameLen)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgStart(hStmt->stmtHandler, procName, procNameLen, &error);
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

YapiResult yapiPdbgShowSource(YapiStmt* hStmt)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgShowSource(hStmt->stmtHandler, &error);
}

YapiResult yapiPdbgDeleteAllBreakpoints(YapiStmt* hStmt)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgDeleteAllBreakpoints(hStmt->stmtHandler, &error);
}

YapiResult yapiPdbgAddBreakpoint(YapiStmt* hStmt, int lineNum, uint32_t* bpID)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgAddBreakpoint(hStmt->stmtHandler, lineNum, bpID, &error);
}

YapiResult yapiPdbgDeleteBreakpoint(YapiStmt* hStmt, uint32_t bpID)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgDeleteBreakpoint(hStmt->stmtHandler, bpID, &error);
}

YapiResult yapiPdbgShowBreakpoints(YapiStmt* hStmt)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgShowBreakpoints(hStmt->stmtHandler, &error);
}

YapiResult yapiPdbgShowFrameVariables(YapiStmt* hStmt)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgShowFrameVariables(hStmt->stmtHandler, &error);
}

YapiResult yapiPdbgShowFrames(YapiStmt* hStmt)
{
    YapiErrorMsg error;
    yapiInitError(&error);
    return yapiCiPdbgShowFrames(hStmt->stmtHandler, &error);
}