/*
Copyright  2022, YashanDB and/or its affiliates. All rights reserved.
YashanDB Driver for golang is licensed under the terms of the mulan PSL v2.0

License: 	http://license.coscl.org.cn/MulanPSL2
Home page: 	https://www.yashandb.com/
*/
package yasdb

/*
#cgo CFLAGS: -I./yacapi/include

#include "yacapi.h"
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"fmt"
	"unsafe"
)

type DebugRunningAttr int8
type DebugFrameAttr int8
type DebugVarAttr int8
type DebugBpAttr int8
type DebuggerStatus int8

const (
	DBG_RUNNING_ATTR_STATUS          DebugRunningAttr = 0
	DBG_RUNNING_ATTR_OBJ_ID          DebugRunningAttr = 1
	DBG_RUNNING_ATTR_SUB_ID          DebugRunningAttr = 2
	DBG_RUNNING_ATTR_LINE_NO         DebugRunningAttr = 3
	DBG_RUNNING_ATTR_CLASS_NAME_LEN  DebugRunningAttr = 4
	DBG_RUNNING_ATTR_METHOD_NAME_LEN DebugRunningAttr = 5
	DBG_RUNNING_ATTR_CLASS_NAME      DebugRunningAttr = 6
	DBG_RUNNING_ATTR_METHOD_NAME     DebugRunningAttr = 7

	DBG_FRAME_ATTR_OBJ_ID          DebugFrameAttr = 0
	DBG_FRAME_ATTR_SUB_ID          DebugFrameAttr = 1
	DBG_FRAME_ATTR_LINE_NO         DebugFrameAttr = 2
	DBG_FRAME_ATTR_STACK_NO        DebugFrameAttr = 3
	DBG_FRAME_ATTR_CLASS_NAME_LEN  DebugFrameAttr = 4
	DBG_FRAME_ATTR_METHOD_NAME_LEN DebugFrameAttr = 5
	DBG_FRAME_ATTR_CLASS_NAME      DebugFrameAttr = 6
	DBG_FRAME_ATTR_METHOD_NAME     DebugFrameAttr = 7

	DBG_VAR_ATTR_BLOCK_NO   DebugVarAttr = 0
	DBG_VAR_ATTR_TYPE       DebugVarAttr = 1
	DBG_VAR_ATTR_IS_GLOBAL  DebugVarAttr = 2
	DBG_VAR_ATTR_NAME_LEN   DebugVarAttr = 3
	DBG_VAR_ATTR_NAME       DebugVarAttr = 4
	DBG_VAR_ATTR_VALUE_SIZE DebugVarAttr = 5

	DBG_BP_ATTR_OBJ_ID  DebugBpAttr = 0
	DBG_BP_ATTR_SUB_ID  DebugBpAttr = 1
	DBG_BP_ATTR_LINE_NO DebugBpAttr = 2
	DBG_BP_ATTR_BP_ID   DebugBpAttr = 3

	DBG_STATUS_OFF DebuggerStatus = 0
	DBG_STATUS_ON  DebuggerStatus = 1
)

type PlsqlDebug struct {
	Stmt *YasStmt
}

func NewPlsqlDebug(dsn string, sql string, args ...any) (*PlsqlDebug, error) {
	stmt, err := GenPdbgStatement(dsn, sql, args...)
	if err != nil {
		return nil, err
	}
	return &PlsqlDebug{Stmt: stmt}, nil
}

func (p *PlsqlDebug) Start() error {
	return PdbgStart(p.Stmt)
}

func (p *PlsqlDebug) CheckVersion(objId uint64, subId uint16, version uint32) error {
	return PdbgCheckVersion(p.Stmt, objId, subId, version)
}

func (p *PlsqlDebug) Abort() error {
	return PdbgAbort(p.Stmt)
}

func (p *PlsqlDebug) Continue() error {
	return PdbgContinue(p.Stmt)
}

func (p *PlsqlDebug) StepInto() error {
	return PdbgStepInto(p.Stmt)
}

func (p *PlsqlDebug) StepOut() error {
	return PdbgStepOut(p.Stmt)
}

func (p *PlsqlDebug) StepNext() error {
	return PdbgStepNext(p.Stmt)
}

func (p *PlsqlDebug) DeleteAllBreakpoint() error {
	return PdbgDeleteAllBreakpoints(p.Stmt)
}

func (p *PlsqlDebug) AddBreakpoint(objId uint64, subId uint16, lineNo uint32) (uint32, error) {
	return PdbgAddBreakpoint(p.Stmt, objId, subId, lineNo)
}

func (p *PlsqlDebug) DeleteBreakpoint(objId uint64, subId uint16, lineNo uint32) error {
	return PdbgDeleteBreakpoint(p.Stmt, objId, subId, lineNo)
}

func (p *PlsqlDebug) GetBreakpointsCount() (uint32, error) {
	return PdbgGetBreakpointsCount(p.Stmt)
}

func (p *PlsqlDebug) GetAllVars() (uint32, error) {
	return PdbgGetAllVars(p.Stmt)
}

func (p *PlsqlDebug) GetAllFrames() (uint32, error) {
	return PdbgGetAllFrames(p.Stmt)
}

func (p *PlsqlDebug) GetRunningData(attr DebugRunningAttr, value interface{}) error {
	return PdbgGetRunningData(p.Stmt, attr, value)
}

func (p *PlsqlDebug) GetFrameData(id uint32, attr DebugFrameAttr, value interface{}) error {
	return PdbgGetFrameData(p.Stmt, id, attr, value)
}

func (p *PlsqlDebug) GetVarData(id uint32, attr DebugVarAttr, value interface{}) error {
	return PdbgGetVarData(p.Stmt, id, attr, value)
}

func (p *PlsqlDebug) GetVarValue(id uint32) (string, error) {
	return PdbgGetVarValue(p.Stmt, id)
}

func (p *PlsqlDebug) GetBreakpointData(id uint32, attr DebugBpAttr, value interface{}) error {
	return PdbgGetBreakpointData(p.Stmt, id, attr, value)
}

func (p *PlsqlDebug) GetAllRunningData() (*PdbgRunningData, error) {
	return PdbgGetAllRunningData(p.Stmt)
}

func (p *PlsqlDebug) GetAllBreakpointData() ([]*PdbgBreakpointData, error) {
	return PdbgGetAllBreakpointData(p.Stmt)
}

func (p *PlsqlDebug) GetAllVarData() ([]*PdbgVarData, error) {
	return PdbgGetAllVarData(p.Stmt)
}

func (p *PlsqlDebug) GetAllFrameData() ([]*PdbgFrameData, error) {
	return PdbgGetAllFrameData(p.Stmt)
}

func (p *PlsqlDebug) Close() error {
	return PdbgStatementClose(p.Stmt)
}

func (p *PlsqlDebug) GetBindOutValue() error {
	return PbdgGetBindOutValue(p.Stmt)
}

func GenPdbgStatement(dsn string, sql string, args ...any) (*YasStmt, error) {
	conn, err := GenYasconn(dsn)
	if err != nil {
		return nil, err
	}
	stmt, err := PrepareContext(conn, context.Background(), sql)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	if len(args) == 0 {
		return stmt, nil
	}

	nArgs, err := ConvertToNameValue(args...)
	if err != nil {
		return nil, err
	}
	if err := stmt.bindValues(nArgs); err != nil {
		_ = conn.Close()
		_ = stmt.Close()
		return nil, err
	}
	return stmt, nil
}

func PdbgStart(stmt *YasStmt) error {
	return checkYasError(C.yapiPdbgStart(stmt.Stmt))
}

func PdbgCheckVersion(stmt *YasStmt, objId uint64, subId uint16, version uint32) error {
	return checkYasError(C.yapiPdbgCheckVersion(stmt.Stmt, C.uint64_t(objId), C.uint16_t(subId), C.uint32_t(version)))
}

func PdbgAbort(stmt *YasStmt) error {
	return checkYasError(C.yapiPdbgAbort(stmt.Stmt))
}

func PdbgContinue(stmt *YasStmt) error {
	return checkYasError(C.yapiPdbgContinue(stmt.Stmt))
}

func PdbgStepInto(stmt *YasStmt) error {
	return checkYasError(C.yapiPdbgStepInto(stmt.Stmt))
}

func PdbgStepOut(stmt *YasStmt) error {
	return checkYasError(C.yapiPdbgStepOut(stmt.Stmt))
}

func PdbgStepNext(stmt *YasStmt) error {
	return checkYasError(C.yapiPdbgStepNext(stmt.Stmt))
}

func PbdgGetBindOutValue(stmt *YasStmt) error {
	if err := stmt.getBindValueDest(); err != nil {
		return err
	}
	return nil
}

func PdbgStatementClose(stmt *YasStmt) error {
	stmt.freeBindValues()
	if err := stmt.Close(); err != nil {
		return err
	}
	return stmt.Conn.Close()
}

func PdbgDeleteAllBreakpoints(stmt *YasStmt) error {
	return checkYasError(C.yapiPdbgDeleteAllBreakpoints(stmt.Stmt))
}

func PdbgAddBreakpoint(stmt *YasStmt, objId uint64, subId uint16, lineNo uint32) (uint32, error) {
	bpID := (*C.uint32_t)(C.malloc(32))
	defer C.free(unsafe.Pointer(bpID))
	err := checkYasError(C.yapiPdbgAddBreakpoint(stmt.Stmt, C.uint64_t(objId), C.uint16_t(subId), C.uint32_t(lineNo), bpID))
	if err != nil {
		return 0, err
	}
	return uint32(*bpID), nil
}

func PdbgDeleteBreakpoint(stmt *YasStmt, objId uint64, subId uint16, lineNo uint32) error {
	return checkYasError(C.yapiPdbgDeleteBreakpoint(stmt.Stmt, C.uint64_t(objId), C.uint16_t(subId), C.uint32_t(lineNo)))
}

func PdbgGetBreakpointsCount(stmt *YasStmt) (uint32, error) {
	bpCount := new(C.uint32_t)
	err := checkYasError(C.yapiPdbgGetBreakpointsCount(stmt.Stmt, bpCount))
	if err != nil {
		return 0, err
	}
	return uint32(*bpCount), nil
}

func PdbgGetAllVars(stmt *YasStmt) (uint32, error) {
	varCount := new(C.uint32_t)
	err := checkYasError(C.yapiPdbgGetAllVars(stmt.Stmt, varCount))
	if err != nil {
		return 0, err
	}
	return uint32(*varCount), nil
}

func PdbgGetAllFrames(stmt *YasStmt) (uint32, error) {
	frameCount := new(C.uint32_t)
	err := checkYasError(C.yapiPdbgGetAllFrames(stmt.Stmt, frameCount))
	if err != nil {
		return 0, err
	}
	return uint32(*frameCount), nil
}

func PdbgGetRunningData(stmt *YasStmt, attr DebugRunningAttr, value interface{}) error {
	switch attr {
	case DBG_RUNNING_ATTR_STATUS:
		data, ok := value.(*DebuggerStatus)
		if !ok {
			return fmt.Errorf("the value parameter type must be *DebuggerStatus")
		}
		status := new(C.uint32_t)
		bufferSize := C.int32_t(32)
		err := checkYasError(C.yapiPdbgGetRunningData(stmt.Stmt, C.YAPI_DBG_RUNNING_ATTR_STATUS, C.YapiPointer(status), bufferSize))
		if err != nil {
			return err
		}
		if *status == C.uint32_t(C.YAPI_DBG_STATUS_ON) {
			*data = DBG_STATUS_ON
		} else {
			*data = DBG_STATUS_OFF
		}

	case DBG_RUNNING_ATTR_OBJ_ID:
		data, ok := value.(*uint64)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint64")
		}

		objId := new(C.uint64_t)
		bufferSize := C.int32_t(64)
		err := checkYasError(C.yapiPdbgGetRunningData(stmt.Stmt, C.YAPI_DBG_RUNNING_ATTR_OBJ_ID, C.YapiPointer(objId), bufferSize))
		if err != nil {
			return err
		}
		*data = uint64(*objId)

	case DBG_RUNNING_ATTR_SUB_ID:
		data, ok := value.(*uint16)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint16")
		}

		subId := new(C.uint16_t)
		bufferSize := C.int32_t(64)
		err := checkYasError(C.yapiPdbgGetRunningData(stmt.Stmt, C.YAPI_DBG_RUNNING_ATTR_SUB_ID, C.YapiPointer(subId), bufferSize))
		if err != nil {
			return err
		}
		*data = uint16(*subId)

	case DBG_RUNNING_ATTR_LINE_NO:
		data, ok := value.(*uint32)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint32")
		}

		lineNo := new(C.uint32_t)
		bufferSize := C.int32_t(32)
		err := checkYasError(C.yapiPdbgGetRunningData(stmt.Stmt, C.YAPI_DBG_RUNNING_ATTR_LINE_NO, C.YapiPointer(lineNo), bufferSize))
		if err != nil {
			return err
		}
		*data = uint32(*lineNo)

	case DBG_RUNNING_ATTR_CLASS_NAME_LEN:
		data, ok := value.(*uint32)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint32")
		}

		nameLen := new(C.uint32_t)
		bufferSize := C.int32_t(32)
		err := checkYasError(C.yapiPdbgGetRunningData(stmt.Stmt, C.YAPI_DBG_RUNNING_ATTR_CLASS_NAME_LEN, C.YapiPointer(nameLen), bufferSize))
		if err != nil {
			return err
		}
		*data = uint32(*nameLen)

	case DBG_RUNNING_ATTR_CLASS_NAME:
		data, ok := value.(*string)
		if !ok {
			return fmt.Errorf("the value parameter type must be *string")
		}
		className := new(C.char)
		var nameLen uint32
		if err := PdbgGetRunningData(stmt, DBG_RUNNING_ATTR_CLASS_NAME_LEN, &nameLen); err != nil {
			return err
		}
		bufferSize := C.int32_t(nameLen)
		err := checkYasError(C.yapiPdbgGetRunningData(stmt.Stmt, C.YAPI_DBG_RUNNING_ATTR_CLASS_NAME, C.YapiPointer(className), bufferSize))
		if err != nil {
			return err
		}
		*data = C.GoString(className)

	case DBG_RUNNING_ATTR_METHOD_NAME_LEN:
		data, ok := value.(*uint32)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint32")
		}

		nameLen := new(C.uint32_t)
		bufferSize := C.int32_t(32)
		err := checkYasError(C.yapiPdbgGetRunningData(stmt.Stmt, C.YAPI_DBG_RUNNING_ATTR_METHOD_NAME_LEN, C.YapiPointer(nameLen), bufferSize))
		if err != nil {
			return err
		}
		*data = uint32(*nameLen)

	case DBG_RUNNING_ATTR_METHOD_NAME:
		data, ok := value.(*string)
		if !ok {
			return fmt.Errorf("the value parameter type must be *string")
		}
		methodName := new(C.char)
		var nameLen uint32
		if err := PdbgGetRunningData(stmt, DBG_RUNNING_ATTR_METHOD_NAME_LEN, &nameLen); err != nil {
			return err
		}
		bufferSize := C.int32_t(nameLen)
		err := checkYasError(C.yapiPdbgGetRunningData(stmt.Stmt, C.YAPI_DBG_RUNNING_ATTR_METHOD_NAME, C.YapiPointer(methodName), bufferSize))
		if err != nil {
			return err
		}
		*data = C.GoString(methodName)

	default:
		return fmt.Errorf("unsupport debug running attr %v", attr)
	}
	return nil
}

func PdbgGetFrameData(stmt *YasStmt, id uint32, attr DebugFrameAttr, value interface{}) error {
	switch attr {
	case DBG_FRAME_ATTR_OBJ_ID:
		data, ok := value.(*uint64)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint64")
		}

		outValue := new(C.uint64_t)
		bufferSize := C.int32_t(64)
		err := checkYasError(C.yapiPdbgGetFrameData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_FRAME_ATTR_OBJ_ID, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = uint64(*outValue)

	case DBG_FRAME_ATTR_SUB_ID:
		data, ok := value.(*uint16)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint16")
		}

		outValue := new(C.uint16_t)
		bufferSize := C.int32_t(16)
		err := checkYasError(C.yapiPdbgGetFrameData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_FRAME_ATTR_SUB_ID, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = uint16(*outValue)

	case DBG_FRAME_ATTR_LINE_NO:
		data, ok := value.(*uint32)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint32")
		}

		outValue := new(C.uint32_t)
		bufferSize := C.int32_t(32)
		err := checkYasError(C.yapiPdbgGetFrameData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_FRAME_ATTR_LINE_NO, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = uint32(*outValue)

	case DBG_FRAME_ATTR_STACK_NO:
		data, ok := value.(*uint32)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint32")
		}

		outValue := new(C.uint32_t)
		bufferSize := C.int32_t(32)
		err := checkYasError(C.yapiPdbgGetFrameData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_FRAME_ATTR_STACK_NO, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = uint32(*outValue)

	case DBG_FRAME_ATTR_CLASS_NAME_LEN:
		data, ok := value.(*uint32)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint32")
		}

		nameLen := new(C.uint32_t)
		bufferSize := C.int32_t(32)
		err := checkYasError(C.yapiPdbgGetFrameData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_FRAME_ATTR_CLASS_NAME_LEN, C.YapiPointer(nameLen), bufferSize))
		if err != nil {
			return err
		}
		*data = uint32(*nameLen)

	case DBG_FRAME_ATTR_CLASS_NAME:
		data, ok := value.(*string)
		if !ok {
			return fmt.Errorf("the value parameter type must be *string")
		}
		outValue := new(C.char)
		var nameLen uint32
		if err := PdbgGetFrameData(stmt, id, DBG_FRAME_ATTR_CLASS_NAME_LEN, &nameLen); err != nil {
			return err
		}
		bufferSize := C.int32_t(nameLen)
		err := checkYasError(C.yapiPdbgGetFrameData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_FRAME_ATTR_CLASS_NAME, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = C.GoString(outValue)
	case DBG_FRAME_ATTR_METHOD_NAME_LEN:
		data, ok := value.(*uint32)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint32")
		}

		nameLen := new(C.uint32_t)
		bufferSize := C.int32_t(32)
		err := checkYasError(C.yapiPdbgGetFrameData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_FRAME_ATTR_METHOD_NAME_LEN, C.YapiPointer(nameLen), bufferSize))
		if err != nil {
			return err
		}
		*data = uint32(*nameLen)

	case DBG_FRAME_ATTR_METHOD_NAME:
		data, ok := value.(*string)
		if !ok {
			return fmt.Errorf("the value parameter type must be *string")
		}
		outValue := new(C.char)
		var nameLen uint32
		if err := PdbgGetFrameData(stmt, id, DBG_FRAME_ATTR_METHOD_NAME_LEN, &nameLen); err != nil {
			return err
		}
		bufferSize := C.int32_t(nameLen)
		err := checkYasError(C.yapiPdbgGetFrameData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_FRAME_ATTR_METHOD_NAME, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = C.GoString(outValue)

	default:
		return fmt.Errorf("unsupport debug frame attr %v", attr)
	}
	return nil
}

func PdbgGetVarData(stmt *YasStmt, id uint32, attr DebugVarAttr, value interface{}) error {
	switch attr {
	case DBG_VAR_ATTR_BLOCK_NO:
		data, ok := value.(*uint32)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint32")
		}

		outValue := new(C.uint32_t)
		bufferSize := C.int32_t(unsafe.Sizeof(data))
		err := checkYasError(C.yapiPdbgGetVarData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_VAR_ATTR_BLOCK_NO, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = uint32(*outValue)

	case DBG_VAR_ATTR_TYPE:
		data, ok := value.(*uint8)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint8")
		}

		outValue := new(C.uint8_t)
		bufferSize := C.int32_t(unsafe.Sizeof(data))
		err := checkYasError(C.yapiPdbgGetVarData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_VAR_ATTR_TYPE, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = uint8(*outValue)

	case DBG_VAR_ATTR_IS_GLOBAL:
		data, ok := value.(*bool)
		if !ok {
			return fmt.Errorf("the value parameter type must be *bool")
		}

		outValue := new(C.bool)
		bufferSize := C.int32_t(unsafe.Sizeof(data))
		err := checkYasError(C.yapiPdbgGetVarData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_VAR_ATTR_IS_GLOBAL, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = bool(*outValue)
	case DBG_VAR_ATTR_NAME_LEN:
		data, ok := value.(*uint32)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint32")
		}

		nameLen := new(C.uint32_t)
		bufferSize := C.int32_t(unsafe.Sizeof(data))
		err := checkYasError(C.yapiPdbgGetVarData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_VAR_ATTR_NAME_LEN, C.YapiPointer(nameLen), bufferSize))
		if err != nil {
			return err
		}
		*data = uint32(*nameLen)

	case DBG_VAR_ATTR_NAME:
		data, ok := value.(*string)
		if !ok {
			return fmt.Errorf("the value parameter type must be *string")
		}
		outValue := new(C.char)
		var nameLen uint32
		if err := PdbgGetVarData(stmt, id, DBG_VAR_ATTR_NAME_LEN, &nameLen); err != nil {
			return err
		}
		bufferSize := C.int32_t(nameLen)
		err := checkYasError(C.yapiPdbgGetVarData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_VAR_ATTR_NAME, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = C.GoString(outValue)

	case DBG_VAR_ATTR_VALUE_SIZE:
		data, ok := value.(*uint32)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint32")
		}
		outValue := new(C.uint32_t)
		bufferSize := C.int32_t(unsafe.Sizeof(data))
		err := checkYasError(C.yapiPdbgGetVarData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_VAR_ATTR_VALUE_SIZE, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = uint32(*outValue)

	default:
		return fmt.Errorf("unsupport debug var attr %v", attr)
	}

	return nil
}

func PdbgGetVarValue(stmt *YasStmt, id uint32) (string, error) {
	indicator := new(C.int32_t)
	value := C.YapiPointer(unsafe.Pointer(stringToYasChar("")))
	err := checkYasError(C.yapiPdbgGetVarValue(stmt.Stmt, C.uint32_t(id), C.uint32_t(C.YAPI_TYPE_VARCHAR), value, 300, indicator))
	if err != nil {
		return "", err
	}
	defer C.free(unsafe.Pointer(value))
	if *indicator == C.YAPI_NULL_DATA {
		return "", nil
	}
	return C.GoString((*C.char)(value)), nil
}

func PdbgGetBreakpointData(stmt *YasStmt, id uint32, attr DebugBpAttr, value interface{}) error {
	switch attr {
	case DBG_BP_ATTR_OBJ_ID:
		data, ok := value.(*uint64)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint64")
		}

		outValue := new(C.uint64_t)
		bufferSize := C.int32_t(64)
		err := checkYasError(C.yapiPdbgGetBreakpointData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_BP_ATTR_OBJ_ID, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = uint64(*outValue)

	case DBG_BP_ATTR_SUB_ID:
		data, ok := value.(*uint16)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint16")
		}

		outValue := new(C.uint16_t)
		bufferSize := C.int32_t(16)
		err := checkYasError(C.yapiPdbgGetBreakpointData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_BP_ATTR_SUB_ID, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = uint16(*outValue)

	case DBG_BP_ATTR_LINE_NO:
		data, ok := value.(*uint32)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint32")
		}

		outValue := new(C.uint32_t)
		bufferSize := C.int32_t(32)
		err := checkYasError(C.yapiPdbgGetBreakpointData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_BP_ATTR_LINE_NO, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = uint32(*outValue)

	case DBG_BP_ATTR_BP_ID:
		data, ok := value.(*uint32)
		if !ok {
			return fmt.Errorf("the value parameter type must be *uint32")
		}

		outValue := new(C.uint32_t)
		bufferSize := C.int32_t(32)
		err := checkYasError(C.yapiPdbgGetBreakpointData(stmt.Stmt, C.uint32_t(id), C.YAPI_DBG_BP_ATTR_BP_ID, C.YapiPointer(outValue), bufferSize))
		if err != nil {
			return err
		}
		*data = uint32(*outValue)

	default:
		return fmt.Errorf("unsupport debug breakpoint attr %v", attr)
	}

	return nil
}

type PdbgRunningData struct {
	Status     DebuggerStatus
	ObjId      uint64
	SubId      uint16
	LineNo     uint32
	ClassName  string
	MethodName string
}

func PdbgGetAllRunningData(stmt *YasStmt) (*PdbgRunningData, error) {
	data := PdbgRunningData{}
	if err := PdbgGetRunningData(stmt, DBG_RUNNING_ATTR_STATUS, &data.Status); err != nil {
		return nil, err
	}

	if err := PdbgGetRunningData(stmt, DBG_RUNNING_ATTR_OBJ_ID, &data.ObjId); err != nil {
		return nil, err
	}

	if err := PdbgGetRunningData(stmt, DBG_RUNNING_ATTR_SUB_ID, &data.SubId); err != nil {
		return nil, err
	}

	if err := PdbgGetRunningData(stmt, DBG_RUNNING_ATTR_CLASS_NAME, &data.ClassName); err != nil {
		return nil, err
	}

	if err := PdbgGetRunningData(stmt, DBG_RUNNING_ATTR_METHOD_NAME, &data.MethodName); err != nil {
		return nil, err
	}

	if err := PdbgGetRunningData(stmt, DBG_RUNNING_ATTR_LINE_NO, &data.LineNo); err != nil {
		return nil, err
	}

	return &data, nil
}

type PdbgBreakpointData struct {
	Id     uint32
	ObjId  uint64
	SubId  uint16
	LineNo uint32
	BpId   uint32
}

func PdbgGetAllBreakpointData(stmt *YasStmt) ([]*PdbgBreakpointData, error) {
	bpCount, err := PdbgGetBreakpointsCount(stmt)
	if err != nil {
		return nil, err
	}

	bpDatas := make([]*PdbgBreakpointData, 0, bpCount)
	for i := uint32(0); i < bpCount; i++ {
		data := PdbgBreakpointData{Id: i}
		if err := PdbgGetBreakpointData(stmt, i, DBG_BP_ATTR_OBJ_ID, &data.ObjId); err != nil {
			return nil, err
		}

		if err := PdbgGetBreakpointData(stmt, i, DBG_BP_ATTR_SUB_ID, &data.SubId); err != nil {
			return nil, err
		}

		if err := PdbgGetBreakpointData(stmt, i, DBG_BP_ATTR_LINE_NO, &data.LineNo); err != nil {
			return nil, err
		}

		if err := PdbgGetBreakpointData(stmt, i, DBG_BP_ATTR_BP_ID, &data.BpId); err != nil {
			return nil, err
		}
		bpDatas = append(bpDatas, &data)
	}

	return bpDatas, nil
}

type PdbgVarData struct {
	Id           uint32
	BlockNo      uint32
	DataType     uint8
	DataTypeName string
	IsGlobal     bool
	Name         string
	Size         uint32
	Value        string
}

func PdbgGetAllVarData(stmt *YasStmt) ([]*PdbgVarData, error) {
	varCount, err := PdbgGetAllVars(stmt)
	if err != nil {
		return nil, err
	}
	varDatas := make([]*PdbgVarData, 0, varCount)
	for i := uint32(0); i < varCount; i++ {
		data := PdbgVarData{Id: i}
		if err := PdbgGetVarData(stmt, i, DBG_VAR_ATTR_BLOCK_NO, &data.BlockNo); err != nil {
			return nil, err
		}

		if err := PdbgGetVarData(stmt, i, DBG_VAR_ATTR_TYPE, &data.DataType); err != nil {
			return nil, err
		}

		data.DataTypeName = GetDatabaseTypeName(uint32(data.DataType))

		if err := PdbgGetVarData(stmt, i, DBG_VAR_ATTR_IS_GLOBAL, &data.IsGlobal); err != nil {
			return nil, err
		}

		if err := PdbgGetVarData(stmt, i, DBG_VAR_ATTR_NAME, &data.Name); err != nil {
			return nil, err
		}

		if err := PdbgGetVarData(stmt, i, DBG_VAR_ATTR_VALUE_SIZE, &data.Size); err != nil {
			return nil, err
		}

		value, err := PdbgGetVarValue(stmt, i)
		if err != nil {
			return nil, err
		}
		data.Value = value
		varDatas = append(varDatas, &data)
	}
	return varDatas, nil
}

type PdbgFrameData struct {
	Id         uint32
	ObjId      uint64
	SubId      uint16
	LineNo     uint32
	StackNo    uint32
	ClassName  string
	MethodName string
}

func PdbgGetAllFrameData(stmt *YasStmt) ([]*PdbgFrameData, error) {
	frameCount, err := PdbgGetAllFrames(stmt)
	if err != nil {
		return nil, err
	}
	frameDatas := make([]*PdbgFrameData, 0, frameCount)
	for i := uint32(0); i < frameCount; i++ {
		data := PdbgFrameData{Id: i}
		if err := PdbgGetFrameData(stmt, i, DBG_FRAME_ATTR_OBJ_ID, &data.ObjId); err != nil {
			return nil, err
		}

		if err := PdbgGetFrameData(stmt, i, DBG_FRAME_ATTR_SUB_ID, &data.SubId); err != nil {
			return nil, err
		}

		if err := PdbgGetFrameData(stmt, i, DBG_FRAME_ATTR_LINE_NO, &data.LineNo); err != nil {
			return nil, err
		}

		if err := PdbgGetFrameData(stmt, i, DBG_FRAME_ATTR_STACK_NO, &data.StackNo); err != nil {
			return nil, err
		}

		if err := PdbgGetFrameData(stmt, i, DBG_FRAME_ATTR_CLASS_NAME, &data.ClassName); err != nil {
			return nil, err
		}

		if err := PdbgGetFrameData(stmt, i, DBG_FRAME_ATTR_METHOD_NAME, &data.MethodName); err != nil {
			return nil, err
		}

		frameDatas = append(frameDatas, &data)
	}
	return frameDatas, nil
}
