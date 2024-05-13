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
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"
	"time"
	"unsafe"
)

type YasStmt struct {
	Conn    *YasConn
	Stmt    *C.YapiStmt
	closed  bool
	SqlType uint32
	ctx     context.Context
	binds   []*bindStruct
	sync.Mutex
}

// Query executes a query that may return rows, such as a SELECT.
//
// Deprecated: Drivers should implement StmtQueryContext instead (or additionally).
func (stmt *YasStmt) Query(args []driver.Value) (driver.Rows, error) {
	nargs := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		nargs[i].Ordinal = i + 1
		nargs[i].Value = arg
	}
	return stmt.QueryContext(context.Background(), nargs)
}

// QueryContext executes a query that may return rows, such as a SELECT.
//
// QueryContext must honor the context timeout and return when it is canceled.
func (stmt *YasStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	stmt.Lock()
	defer stmt.Unlock()
	stmt.ctx = ctx

	defer stmt.freeBindValues()
	if err := stmt.bindValues(args); err != nil {
		return nil, err
	}
	return stmt.query()
}

// Exec executes a query that doesn't return rows, such as an INSERT or UPDATE.
//
// Deprecated: Drivers should implement StmtExecContext instead (or additionally).
func (stmt *YasStmt) Exec(args []driver.Value) (driver.Result, error) {
	nargs := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		nargs[i].Ordinal = i + 1
		nargs[i].Value = arg
	}

	return stmt.ExecContext(context.Background(), nargs)
}

// ExecContext executes a query that doesn't return rows, such as an INSERT or UPDATE.
//
// ExecContext must honor the context timeout and return when it is canceled.
func (stmt *YasStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	stmt.Lock()
	defer stmt.Unlock()
	stmt.ctx = ctx

	defer stmt.freeBindValues()
	if err := stmt.bindValues(args); err != nil {
		return nil, err
	}

	return stmt.exec()
}

// NumInput returns the number of placeholder parameters.
func (stmt *YasStmt) NumInput() int {
	return -1
}

// Close closes the statement.
func (stmt *YasStmt) Close() error {
	if stmt.closed {
		return nil
	}
	stmt.closed = true
	return stmt.yapiReleaseStmt()
}

// CheckNamedValue is called before passing arguments to the driver
// and is called in place of any ColumnConverter. CheckNamedValue must do type
// validation and conversion as appropriate for the driver.
func (stmt *YasStmt) CheckNamedValue(namedValue *driver.NamedValue) error {
	switch namedValue.Value.(type) {
	case sql.Out:
		return nil
	}
	return driver.ErrSkip
}

func (stmt *YasStmt) query() (driver.Rows, error) {
	if stmt.ctx.Err() != nil {
		return nil, stmt.ctx.Err()
	}

	done := make(chan struct{})
	go stmt.Conn.handleYacCancel(stmt.ctx, done)
	err := stmt.yacExecute()
	close(done)
	if err != nil {
		return nil, err
	}

	fetchRows, err := stmt.getFetchRows()
	if err != nil {
		return nil, err
	}
	rows := YasRows{
		stmt:      stmt,
		fetchRows: fetchRows,
	}
	return &rows, nil
}

func (stmt *YasStmt) exec() (driver.Result, error) {
	if stmt.ctx.Err() != nil {
		return nil, stmt.ctx.Err()
	}

	done := make(chan struct{})
	go stmt.Conn.handleYacCancel(stmt.ctx, done)
	err := stmt.yacExecute()
	close(done)
	if err != nil {
		return nil, err
	}

	rowsAffected, rowsAffectedErr := stmt.getRowsAffected()
	result := YasResult{
		rowsAffected:    rowsAffected,
		rowsAffectedErr: rowsAffectedErr,
	}

	if err := stmt.getBindValueDest(); err != nil {
		return nil, err
	}
	return &result, nil
}

func (stmt *YasStmt) yacExecute() error {
	return checkYasError(C.yapiExecute(stmt.Stmt))
}

func (stmt *YasStmt) yapiReleaseStmt() error {
	if stmt.Stmt == nil {
		return nil
	}
	if err := checkYasError(C.yapiReleaseStmt(stmt.Stmt)); err != nil {
		return err
	}
	stmt.Stmt = nil
	return nil
}

func (stmt *YasStmt) getFetchRows() ([]*yasRow, error) {
	columns := C.int16_t(0)
	if err := checkYasError(C.yapiNumResultCols(stmt.Stmt, &columns)); err != nil {
		return nil, err
	}
	columnCount := int(columns)
	yasRows := make([]*yasRow, 0, columnCount)
	for i := 0; i < columnCount; i++ {
		row, err := stmt.getFetchRow(i)
		if err != nil {
			freeFetchRows(yasRows)
			return nil, err
		}
		yasRows = append(yasRows, row)
	}
	return yasRows, nil
}

func (stmt *YasStmt) getFetchRow(pos int) (*yasRow, error) {
	item := C.YapiColumnDesc{}
	if err := checkYasError(C.yapiDescribeCol2(stmt.Stmt, C.uint16_t(pos), &item)); err != nil {
		return nil, err
	}
	yacType := C.YapiType(item._type)
	size, indicator := uint32(item.size), (*C.int32_t)(C.malloc(4))
	row := NewYasRow(size, yacType)
	bufLen := int32(size)
	freeType := notFree

	switch yacType {
	case C.YAPI_TYPE_NCHAR, C.YAPI_TYPE_NVARCHAR:
		yacType = C.YAPI_TYPE_VARCHAR
		bufLen = int32(sizeToAlign4(size)) + 1
		row.Data = mallocBytes(uint32(bufLen))
		freeType = normalFree
	case C.YAPI_TYPE_CHAR, C.YAPI_TYPE_VARCHAR:
		bufLen = int32(sizeToAlign4(size)) + 1
		row.Data = mallocBytes(uint32(bufLen))
		freeType = normalFree
	case C.YAPI_TYPE_NUMBER, C.YAPI_TYPE_YM_INTERVAL, C.YAPI_TYPE_DS_INTERVAL: // number to string
		yacType = C.YAPI_TYPE_VARCHAR
		bufLen = int32(sizeToAlign4(uint32(item.precision) + 8))
		row.Data = mallocBytes(uint32(bufLen))
		freeType = normalFree
	case C.YAPI_TYPE_DATE, C.YAPI_TYPE_TIMESTAMP, C.YAPI_TYPE_SHORTDATE, C.YAPI_TYPE_SHORTTIME:
		bufLen = 12
		row.Data = mallocBytes(uint32(bufLen))
		freeType = normalFree
	case C.YAPI_TYPE_CLOB, C.YAPI_TYPE_BLOB:
		desc := new(unsafe.Pointer)
		if err := checkYasError(C.yapiLobDescAlloc(stmt.Conn.Conn, yacType, desc)); err != nil {
			return nil, err
		}
		bufLen = -1
		row.Data = unsafe.Pointer(desc)
		freeType = lobFree
	case C.YAPI_TYPE_BOOL, C.YAPI_TYPE_TINYINT, C.YAPI_TYPE_SMALLINT, C.YAPI_TYPE_INTEGER, C.YAPI_TYPE_BIGINT, C.YAPI_TYPE_FLOAT, C.YAPI_TYPE_DOUBLE, C.YAPI_TYPE_BINARY, C.YAPI_TYPE_BIT:
		row.Data = mallocBytes(size)
		freeType = normalFree
	case C.YAPI_TYPE_ROWID:
		yacType = C.YAPI_TYPE_VARCHAR
		bufLen = 44
		row.Data = mallocBytes(uint32(bufLen))
		freeType = normalFree
	default:
		yacType = C.YAPI_TYPE_VARCHAR
		bufLen = _DefaultSize
		row.Data = mallocBytes(uint32(bufLen))
		freeType = normalFree
	}
	row.Indicator = indicator
	row.freeType = freeType
	if err := checkYasError(
		C.yapiBindColumn(
			stmt.Stmt,
			C.uint16_t(pos),
			yacType,
			C.YapiPointer(row.Data),
			C.int32_t(bufLen),
			indicator,
		),
	); err != nil {
		freeFetchRow(row)
		return nil, err
	}
	row.name = C.GoString(item.name)
	return row, nil
}

func (stmt *YasStmt) getRowsAffected() (int64, error) {
	var rowsCount C.uint32_t
	size := C.int32_t(unsafe.Sizeof(rowsCount))
	s_length := C.int32_t(0)
	err := checkYasError(
		C.yapiGetStmtAttr(
			stmt.Stmt,
			C.YAPI_ATTR_ROWS_AFFECTED,
			unsafe.Pointer(&rowsCount),
			size,
			&s_length,
		),
	)
	return int64(rowsCount), err
}

func (stmt *YasStmt) bindValues(args []driver.NamedValue) error {
	if len(args) == 0 {
		return nil
	}
	stmt.binds = make([]*bindStruct, 0, len(args))
	var err error
	for index, narg := range args {
		arg := narg.Value
		sqlOut, isOut := arg.(sql.Out)
		bind := &bindStruct{}

		if isOut {
			bind, err = stmt.getOutputBindValue(sqlOut)
			bind.out = sqlOut
		} else {
			bind, err = stmt.getInputBindValue(arg)
		}
		if err != nil {
			return err
		}

		if len(narg.Name) == 0 {
			err = stmt.yacBindParameter(bind, intToYacUint16(index+1))
		} else {
			err = stmt.yacBindParameterByName(bind, narg.Name)
		}
		if err != nil {
			return err
		}
		stmt.binds = append(stmt.binds, bind)
	}

	return nil
}

func (stmt *YasStmt) yacBindParameter(b *bindStruct, pos C.uint16_t) error {
	if err := checkYasError(
		C.yapiBindParameter(
			stmt.Stmt,
			pos,
			b.direction,
			b.yacType,
			b.value,
			b.bindSize,
			C.int32_t(0),
			b.indicator,
		),
	); err != nil {
		return err
	}
	return nil
}

func (stmt *YasStmt) yacBindParameterByName(b *bindStruct, name string) error {
	charName := stringToYasChar(name)
	defer C.free(unsafe.Pointer(charName))
	if err := checkYasError(
		C.yapiBindParameterByName(
			stmt.Stmt,
			charName,
			b.direction,
			b.yacType,
			b.value,
			b.bindSize,
			C.int32_t(0),
			b.indicator,
		),
	); err != nil {
		return err
	}
	return nil
}

func (stmt *YasStmt) getInputBindValue(arg driver.Value) (*bindStruct, error) {
	bind := &bindStruct{}
	var (
		yacType   C.YapiType
		bindSize  C.int32_t
		value     C.YapiPointer
		indicator *C.int32_t
		bufLength C.int32_t
		freeType  valueFreeType
	)

	bindSize = C.int32_t(unsafe.Sizeof(arg)) + 1
	bufLength = C.int32_t(bindSize - 1)
	indicator = new(C.int32_t)
	*indicator = C.int32_t(bindSize - 1)
	freeType = notFree

	switch v := arg.(type) {
	case int64:
		yacType = C.YAPI_TYPE_BIGINT
		value = C.YapiPointer(unsafe.Pointer(&v))
	case float64:
		yacType = C.YAPI_TYPE_DOUBLE
		value = C.YapiPointer(unsafe.Pointer(&v))
	case bool:
		yacType = C.YAPI_TYPE_BOOL
		value = C.YapiPointer(unsafe.Pointer(&v))
	case string:
		yacType = C.YAPI_TYPE_VARCHAR
		bindSize = C.int32_t(len(v)) + 1
		bufLength = C.int32_t(bindSize - 1)
		indicator = nil
		value = C.YapiPointer(unsafe.Pointer(stringToYasChar(v)))
		freeType = normalFree
	case []byte:
		desc, err := stmt.Conn.lobWrite(C.YAPI_TYPE_BLOB, v)
		if err != nil {
			return nil, err
		}
		yacType = C.YAPI_TYPE_BLOB
		bindSize = -1
		bufLength = -1
		indicator = nil
		value = C.YapiPointer(desc)
		freeType = lobFree
	case time.Time:
		yacType = C.YAPI_TYPE_TIMESTAMP
		t := v.UnixNano() / 1e3
		value = C.YapiPointer(unsafe.Pointer(&t))
	case nil:
		yacType = C.YAPI_TYPE_CHAR
		bindSize = 0
		*indicator = C.YAPI_NULL_DATA
		value = C.YapiPointer(unsafe.Pointer(&v))
	default:
		return nil, ErrUnknowType(arg)
	}

	bind.yacType = yacType
	bind.value = value
	bind.bindSize = bindSize
	bind.bufLength = bufLength
	bind.indicator = indicator
	bind.direction = C.YAPI_PARAM_INPUT
	bind.freeType = freeType
	return bind, nil
}

func (stmt *YasStmt) getOutputBindValue(sqlOut sql.Out) (*bindStruct, error) {
	if obi, ok := sqlOut.Dest.(*outputBindInfo); ok {
		return stmt.getOutputBindValueByInfo(obi)
	} else {
		return stmt.getOutputBindValueByDest(sqlOut.Dest)
	}
}

func (stmt *YasStmt) getOutputBindValueByDest(dest interface{}) (*bindStruct, error) {
	bind := &bindStruct{}
	var (
		yacType   C.YapiType
		bindSize  C.int32_t
		value     C.YapiPointer
		indicator *C.int32_t
		bufLength C.int32_t
		arg       driver.Value
		err       error

		freeType = notFree
	)

	arg, err = driver.DefaultParameterConverter.ConvertValue(dest)
	if err != nil {
		return bind, err
	}

	switch arg.(type) {
	case nil:
		arg = dest
		switch arg.(type) {
		case *sql.NullBool:
			arg = false
		case *sql.NullFloat64:
			arg = float64(0)
		case *sql.NullInt64:
			arg = int64(0)
		case *sql.NullString:
			arg = ""
		}
	}

	bindSize = C.int32_t(unsafe.Sizeof(arg)) + 1
	bufLength = C.int32_t(bindSize)
	indicator = new(C.int32_t)
	*indicator = C.int32_t(bindSize - 1)

	switch v := arg.(type) {
	case int64:
		yacType = C.YAPI_TYPE_INTEGER
		value = C.YapiPointer(unsafe.Pointer(&v))
	case float64:
		yacType = C.YAPI_TYPE_DOUBLE
		value = C.YapiPointer(unsafe.Pointer(&v))
	case bool:
		yacType = C.YAPI_TYPE_BOOL
		value = C.YapiPointer(unsafe.Pointer(&v))
	case string:
		yacType = C.YAPI_TYPE_VARCHAR
		bindSize = _OutputBindSize
		bufLength = C.int32_t(bindSize - 1)
		value = C.YapiPointer(unsafe.Pointer(stringToYasChar(v)))
		freeType = normalFree
	case []byte:
		desc, err := stmt.Conn.lobWrite(C.YAPI_TYPE_BLOB, v)
		if err != nil {
			return bind, err
		}
		yacType = C.YAPI_TYPE_BLOB
		bindSize = -1
		bufLength = -1
		indicator = nil
		value = C.YapiPointer(desc)
		freeType = lobFree
	case time.Time:
		yacType = C.YAPI_TYPE_TIMESTAMP
		t := int64(0)
		value = C.YapiPointer(unsafe.Pointer(&t))
	case nil:
		yacType = C.YAPI_TYPE_CHAR
		bindSize = 0
		*indicator = C.YAPI_NULL_DATA
		value = C.YapiPointer(unsafe.Pointer(&v))
	default:
		return bind, ErrUnknowType(v)
	}

	bind.yacType = yacType
	bind.value = value
	bind.bindSize = bindSize
	bind.bufLength = bufLength
	bind.indicator = indicator
	bind.direction = C.YAPI_PARAM_OUTPUT
	bind.freeType = freeType
	return bind, nil
}

func (stmt *YasStmt) getOutputBindValueByInfo(obi *outputBindInfo) (*bindStruct, error) {
	bind := &bindStruct{}
	var (
		yacType   C.YapiType = obi.yacType
		bindSize  C.int32_t
		value     C.YapiPointer
		indicator *C.int32_t
		bufLength C.int32_t

		freeType = notFree
	)

	if obi.bindSize == 0 {
		bindSize = _OutputBindSize
	} else {
		bindSize = obi.bindSize
	}
	bufLength = C.int32_t(bindSize)
	indicator = new(C.int32_t)
	*indicator = C.int32_t(bindSize - 1)

	switch yacType {
	case C.YAPI_TYPE_BLOB:
		v, err := obi.getBlobBindDest()
		if err != nil {
			return bind, err
		}
		desc, err := stmt.Conn.lobWrite(C.YAPI_TYPE_BLOB, *v)
		if err != nil {
			return bind, err
		}
		bindSize = -1
		bufLength = -1
		indicator = nil
		value = C.YapiPointer(desc)
		freeType = lobFree
	case C.YAPI_TYPE_CLOB:
		v, err := obi.getClobBindDest()
		if err != nil {
			return bind, err
		}
		desc, err := stmt.Conn.lobWrite(C.YAPI_TYPE_CLOB, []byte(*v))
		if err != nil {
			return bind, err
		}
		bindSize = -1
		bufLength = -1
		indicator = nil
		value = C.YapiPointer(desc)
		freeType = lobFree
	case C.YAPI_TYPE_CHAR:
		v, err := obi.getCharBindDest()
		if err != nil {
			return bind, err
		}
		bufLength = C.int32_t(bindSize - 1)
		value = C.YapiPointer(unsafe.Pointer(stringToYasChar(*v)))
		freeType = normalFree
	case C.YAPI_TYPE_VARCHAR:
		v, err := obi.getCharBindDest()
		if err != nil {
			return bind, err
		}
		bufLength = C.int32_t(bindSize - 1)
		value = C.YapiPointer(unsafe.Pointer(stringToYasChar(*v)))
		freeType = normalFree
	default:
		return bind, ErrUnknowType(yacType)
	}

	bind.yacType = yacType
	bind.value = value
	bind.bindSize = bindSize
	bind.bufLength = bufLength
	bind.indicator = indicator
	bind.direction = C.YAPI_PARAM_OUTPUT
	bind.freeType = freeType
	return bind, nil
}

func (stmt *YasStmt) getBindValueDest() error {
	var err error
	for index, bind := range stmt.binds {
		if bind.value == nil || bind.out.Dest == nil {
			continue
		}
		switch dest := bind.out.Dest.(type) {
		case *int8:
			*dest = int8(yacPointerToInt64(bind.value))
		case *int16:
			*dest = int16(yacPointerToInt64(bind.value))
		case *int32:
			*dest = int32(yacPointerToInt64(bind.value))
		case *int:
			*dest = int(yacPointerToInt64(bind.value))
		case *int64:
			*dest = yacPointerToInt64(bind.value)
		case *uint8:
			*dest = uint8(yacPointerToUint64(bind.value))
		case *uint16:
			*dest = uint16(yacPointerToUint64(bind.value))
		case *uint32:
			*dest = uint32(yacPointerToUint64(bind.value))
		case *uint:
			*dest = uint(yacPointerToUint64(bind.value))
		case *uint64:
			*dest = yacPointerToUint64(bind.value)
		case *uintptr:
			*dest = uintptr(yacPointerToUint64(bind.value))
		case *float32:
			*dest = float32(yacPointerToFloat64(bind.value))
		case *float64:
			*dest = yacPointerToFloat64(bind.value)
		case *string:
			*dest = C.GoString((*C.char)(bind.value))
		case *time.Time:
			*dest = time.Unix(0, yacPointerToInt64(bind.value)*1e3)
		case *bool:
			*dest = yacPointerToBool(bind.value)
		case *[]byte:
			lobLocator := (**C.YapiLobLocator)(bind.value)
			*dest, err = stmt.Conn.lobRead(*lobLocator)
			if err != nil {
				return err
			}
		case *outputBindInfo:
			switch dest.yacType {
			case C.YAPI_TYPE_BLOB:
				bindDest, _ := dest.getBlobBindDest()
				lobLocator := (**C.YapiLobLocator)(bind.value)
				*bindDest, err = stmt.Conn.lobRead(*lobLocator)
				if err != nil {
					return err
				}
			case C.YAPI_TYPE_CLOB:
				bindDest, _ := dest.getClobBindDest()
				lobLocator := (**C.YapiLobLocator)(bind.value)
				byteDest, err := stmt.Conn.lobRead(*lobLocator)
				if err != nil {
					return err
				}
				*bindDest = string(byteDest)
			case C.YAPI_TYPE_VARCHAR:
				bindDest, _ := dest.getVarcharBindDest()
				*bindDest = C.GoString((*C.char)(bind.value))
			case C.YAPI_TYPE_CHAR:
				bindDest, _ := dest.getVarcharBindDest()
				*bindDest = C.GoString((*C.char)(bind.value))
			}
		default:
			return fmt.Errorf("unknown column %v", index)
		}
	}
	return nil
}

func (stmt *YasStmt) freeBindValues() {
	for _, bind := range stmt.binds {
		stmt.freeBIndValue(bind)
	}
	stmt.binds = []*bindStruct{}
}

func (stmt *YasStmt) freeBIndValue(bind *bindStruct) {
	if bind.value == nil {
		return
	}
	switch bind.freeType {
	case lobFree:
		lobLocator := (**C.YapiLobLocator)(unsafe.Pointer(bind.value))
		stmt.Conn.lobFree(bind.yacType, *lobLocator)
	case normalFree:
		C.free(unsafe.Pointer(bind.value))
	}
	bind.value = nil

}

type outputBindInfo struct {
	yacType  C.YapiType
	dest     interface{}
	bindSize C.int32_t
}
type outputBindOpt func(*outputBindInfo)

func WithTypeClob() outputBindOpt {
	return func(obi *outputBindInfo) { obi.yacType = C.YAPI_TYPE_CLOB }
}

func WithTypeBlob() outputBindOpt {
	return func(obi *outputBindInfo) { obi.yacType = C.YAPI_TYPE_BLOB }
}

func WithTypeVarchar() outputBindOpt {
	return func(obi *outputBindInfo) { obi.yacType = C.YAPI_TYPE_VARCHAR }
}

func WithTypeChar() outputBindOpt {
	return func(obi *outputBindInfo) { obi.yacType = C.YAPI_TYPE_CHAR }
}

func WithBindSize(bindSize uint32) outputBindOpt {
	return func(obi *outputBindInfo) { obi.bindSize = C.int32_t(bindSize) }
}

func NewOutputBindValue(dest interface{}, opts ...outputBindOpt) (*outputBindInfo, error) {
	out := &outputBindInfo{
		dest:     dest,
		bindSize: C.int32_t(0),
		yacType:  C.YapiType(0),
	}
	if err := out.setBindOpt(opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (obi *outputBindInfo) setBindOpt(opts ...outputBindOpt) error {
	for _, opt := range opts {
		opt(obi)
	}
	return obi.checkBindOptParams()
}

func (obi *outputBindInfo) checkBindOptParams() (err error) {
	switch obi.yacType {
	case C.YAPI_TYPE_BLOB:
		_, err = obi.getBlobBindDest()
	case C.YAPI_TYPE_CLOB:
		_, err = obi.getClobBindDest()
	case C.YAPI_TYPE_VARCHAR:
		_, err = obi.getVarcharBindDest()
	case C.YAPI_TYPE_CHAR:
		_, err = obi.getCharBindDest()
	default:
		return ErrUnknowType(obi.yacType)
	}
	return err
}

func (obi *outputBindInfo) getClobBindDest() (*string, error) {
	if value, ok := obi.dest.(*string); ok {
		return value, nil
	}
	return nil, fmt.Errorf("the dest parameter type must be *string")
}

func (obi *outputBindInfo) getBlobBindDest() (*[]byte, error) {
	if value, ok := obi.dest.(*[]byte); ok {
		return value, nil
	}
	return nil, fmt.Errorf("the dest parameter type must be *[]byte")
}

func (obi *outputBindInfo) getCharBindDest() (*string, error) {
	if value, ok := obi.dest.(*string); ok {
		return value, nil
	}
	return nil, fmt.Errorf("the dest parameter type must be *string")
}

func (obi *outputBindInfo) getVarcharBindDest() (*string, error) {
	if value, ok := obi.dest.(*string); ok {
		return value, nil
	}
	return nil, fmt.Errorf("the dest parameter type must be *string")
}
