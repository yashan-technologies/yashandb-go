package yasdb

/*
#cgo CFLAGS: -I./yacapi/include -I./yacapi/src

#include "yacapi.h"
#include "yapi_inc.h"
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

func yapiReleaseConn(conn *C.YapiConnect) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiReleaseConn(conn))
}

func yapiReleaseEnv(env *C.YapiEnv) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiReleaseEnv(env))
}

func yapiAllocEnv(env **C.YapiEnv) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiAllocEnv(env))
}

func yapiSetConnAttr(conn *C.YapiConnect, attr C.YapiConnAttr, value unsafe.Pointer, bufLength C.int32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiSetConnAttr(
		conn,
		attr,
		value,
		bufLength,
	))
}

func yapiGetConnAttr(conn *C.YapiConnect, attr C.YapiConnAttr, value unsafe.Pointer, bufLength C.int32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var stringLen C.int32_t
	return checkYasError(
		C.yapiGetConnAttr(
			conn,
			attr,
			value,
			bufLength,
			&stringLen,
		))
}

func yapiCommit(conn *C.YapiConnect) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiCommit(conn))
}

func yapiRollback(conn *C.YapiConnect) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiRollback(conn))
}

func yapiLobRead(conn *C.YapiConnect, lobLocator *C.YapiLobLocator, bytes *C.uint64_t, buf *C.uint8_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(
		C.yapiLobRead(
			conn,
			lobLocator,
			bytes,
			buf,
			_LobBufLen,
		))
}

func yapiLobGetLength(conn *C.YapiConnect, lobLocator *C.YapiLobLocator, lobLen *C.uint64_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiLobGetLength(conn, lobLocator, lobLen))
}

func yapiLobDescAlloc(conn *C.YapiConnect, yacType C.YapiType, desc *unsafe.Pointer) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiLobDescAlloc(conn, yacType, desc))
}

func yapiLobCreateTemporary(conn *C.YapiConnect, lobLocator *C.YapiLobLocator) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiLobCreateTemporary(conn, lobLocator))
}

func yapiLobWrite(conn *C.YapiConnect, lobLocator *C.YapiLobLocator, buf *C.uint8_t, bufLen C.uint64_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(
		C.yapiLobWrite(
			conn,
			lobLocator,
			nil,
			buf,
			bufLen,
		))
}

func yapiCancel(conn *C.YapiConnect) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiCancel(conn))
}

func yapiPrepare(conn *C.YapiConnect, queryP *C.char, sqlLength C.int32_t, stmt **C.YapiStmt) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(
		C.yapiPrepare(
			conn,
			queryP,
			sqlLength,
			stmt,
		))
}

func yapiGetStmtAttr(stmt *C.YapiStmt, stmtAttr C.YapiStmtAttr, point unsafe.Pointer, sqlSize, sqlLength C.int32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(
		C.yapiGetStmtAttr(
			stmt,
			stmtAttr,
			point,
			sqlSize,
			&sqlLength,
		))
}

func yapiExecute(stmt *C.YapiStmt) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiExecute(stmt))
}

func yapiReleaseStmt(stmt *C.YapiStmt) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiReleaseStmt(stmt))
}

func yapiNumResultCols(stmt *C.YapiStmt, columns *C.int16_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiNumResultCols(stmt, columns))
}

func yapiDescribeCol2(stmt *C.YapiStmt, pos C.uint16_t, item *C.YapiColumnDesc) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiDescribeCol2(stmt, pos, item))
}

func yapiBindColumn(stmt *C.YapiStmt, pos C.uint16_t, yacType C.YapiType, point C.YapiPointer, bufLen C.int32_t, indicator *C.int32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(
		C.yapiBindColumn(
			stmt,
			pos,
			yacType,
			point,
			bufLen,
			indicator,
		),
	)
}

func yapiBindParameter(stmt *C.YapiStmt, b *bindStruct, pos C.uint16_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(
		C.yapiBindParameter(
			stmt,
			pos,
			b.direction,
			b.yacType,
			b.value,
			b.bindSize,
			C.int32_t(0),
			b.indicator,
		),
	)
}

func yapiBindParameterByName(stmt *C.YapiStmt, charName *C.char, b *bindStruct) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(
		C.yapiBindParameterByName(
			stmt,
			charName,
			b.direction,
			b.yacType,
			b.value,
			b.bindSize,
			C.int32_t(0),
			b.indicator,
		),
	)
}

func yapiSetEnvAttr(env *C.YapiEnv, envAttr C.YapiEnvAttr, point unsafe.Pointer, dpLen C.int32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiSetEnvAttr(env, envAttr, point, dpLen))
}

func yapiConnect(env *C.YapiEnv, url *C.char, urlLen C.int16_t, user *C.char, userLen C.int16_t, password *C.char, pwLen C.int16_t, conn **C.YapiConnect) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiConnect(env, url, urlLen, user, userLen, password, pwLen, conn))
}

func yapiAllocConnect(env *C.YapiEnv, conn **C.YapiConnect) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiAllocConnect(env, conn))
}

func yapiConnect2(conn *C.YapiConnect, url *C.char, urlLen C.int16_t, user *C.char, userLen C.int16_t, password *C.char, pwLen C.int16_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiConnect2(conn, url, urlLen, user, userLen, password, pwLen))
}

func yapiPdbgStart(stmt *C.YapiStmt, objId C.uint64_t, subId C.uint16_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgStart(stmt, objId, subId))
}

func yapiPdbgCheckVersion(stmt *C.YapiStmt, objId C.uint64_t, subId C.uint16_t, version C.uint32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgCheckVersion(stmt, objId, subId, version))
}

func yapiPdbgAbort(stmt *C.YapiStmt) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgAbort(stmt))
}

func yapiPdbgContinue(stmt *C.YapiStmt) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgContinue(stmt))
}

func yapiPdbgStepInto(stmt *C.YapiStmt) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgStepInto(stmt))
}

func yapiPdbgStepOut(stmt *C.YapiStmt) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgStepOut(stmt))
}

func yapiPdbgStepNext(stmt *C.YapiStmt) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgStepNext(stmt))
}

func yapiPdbgDeleteAllBreakpoints(stmt *C.YapiStmt) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgDeleteAllBreakpoints(stmt))
}

func yapiPdbgAddBreakpoint(stmt *C.YapiStmt, objId uint64, subId uint16, lineNo uint32, bpID *C.uint32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgAddBreakpoint(stmt, C.uint64_t(objId), C.uint16_t(subId), C.uint32_t(lineNo), bpID))
}

func yapiPdbgDeleteBreakpoint(stmt *C.YapiStmt, objId uint64, subId uint16, lineNo uint32) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgDeleteBreakpoint(stmt, C.uint64_t(objId), C.uint16_t(subId), C.uint32_t(lineNo)))
}

func yapiPdbgGetBreakpointsCount(stmt *C.YapiStmt, bpCount *C.uint32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgGetBreakpointsCount(stmt, bpCount))
}

func yapiPdbgGetAllVars(stmt *C.YapiStmt, varCount *C.uint32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgGetAllVars(stmt, varCount))
}

func yapiPdbgGetAllFrames(stmt *C.YapiStmt, frameCount *C.uint32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgGetAllFrames(stmt, frameCount))
}

func yapiPdbgGetRunningAttrs(stmt *C.YapiStmt, attr C.YapiDebugRunningAttr, point C.YapiPointer, len C.int32_t, stringLen *C.int32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgGetRunningAttrs(stmt, attr, point, len, stringLen))
}

func yapiPdbgGetFrameAttrs(stmt *C.YapiStmt, id C.uint32_t, attr C.YapiDebugFrameAttr, point C.YapiPointer, len C.int32_t, stringLen *C.int32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgGetFrameAttrs(stmt, id, attr, point, len, stringLen))
}

func yapiPdbgGetVarAttrs(stmt *C.YapiStmt, id C.uint32_t, attr C.YapiDebugVarAttr, point C.YapiPointer, len C.int32_t, stringLen *C.int32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgGetVarAttrs(stmt, id, attr, point, len, stringLen))
}

func yapiPdbgGetVarValue(stmt *C.YapiStmt, id C.uint32_t, bindType C.uint32_t, point C.YapiPointer, len C.int32_t, stringLen *C.int32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgGetVarValue(stmt, id, bindType, point, len, stringLen))
}

func yapiPdbgGetBreakpointAttrs(stmt *C.YapiStmt, id C.uint32_t, attr C.YapiDebugBpAttr, point C.YapiPointer, len C.int32_t, stringLen *C.int32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiPdbgGetBreakpointAttrs(stmt, C.uint32_t(id), attr, point, len, stringLen))
}

func yapiFetch(stmt *C.YapiStmt, rows *C.uint32_t) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return checkYasError(C.yapiFetch(stmt, rows))
}
