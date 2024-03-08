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
	"unsafe"
)

type PlsqlDebug struct {
	Stmt     *YasStmt
	procName string
}

func NewPlsqlDebug(dsn string, sql string, procName string, args ...any) (*PlsqlDebug, error) {
	stmt, err := GenPdbgStatement(dsn, sql, args...)
	if err != nil {
		return nil, err
	}
	return &PlsqlDebug{Stmt: stmt, procName: procName}, nil
}

func (p *PlsqlDebug) Start() error {
	return PdbgStart(p.Stmt, p.procName)
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

func (p *PlsqlDebug) AddBreakpoint(lineNum int) (uint32, error) {
	return PdbgAddBreakpoint(p.Stmt, lineNum)
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

func PdbgStart(stmt *YasStmt, procName string) error {
	charProcName := stringToYasChar(procName)
	defer C.free(unsafe.Pointer(charProcName))
	proceNameLen := intToYacUint32(len(procName))

	return checkYasError(
		C.yapiPdbgStart(stmt.Stmt, charProcName, proceNameLen),
	)
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

func PdbgAddBreakpoint(stmt *YasStmt, lineNum int) (uint32, error) {
	// bpID := new(C.uint32_t)
	bpID := (*C.uint32_t)(C.malloc(32))
	defer C.free(unsafe.Pointer(bpID))
	err := checkYasError(C.yapiPdbgAddBreakpoint(stmt.Stmt, intToYacInt(lineNum), bpID))
	if err != nil {
		return 0, err
	}
	return uint32(*bpID), nil
}
func PdbgStatementClose(stmt *YasStmt) error {
	stmt.freeBindValues()
	if err := stmt.Close(); err != nil {
		return err
	}
	return stmt.Conn.Close()
}

// C接口暂未实现
/*
func (p *PlsqlDebug) ShowSorce() error {
	return nil
}

func (p *PlsqlDebug) DeleteAllBreadpoints() error {
	return nil
}

func (p *PlsqlDebug) DeleteBreakpoint(bpID uint32) error {
	return nil
}

func (p *PlsqlDebug) ShowBreakpoints() error {
	return nil
}

func (p *PlsqlDebug) ShowFrameVariables() error {
	return nil
}

func (p *PlsqlDebug) ShowFrames() error {
	return nil
}
*/
