/*
Copyright  2022, YashanDB and/or its affiliates. All rights reserved.
YashanDB Driver for golang is licensed under the terms of the mulan PSL v2.0

License: 	http://license.coscl.org.cn/MulanPSL2
Home page: 	https://www.yashandb.com/
*/

package yasdb

/*
#cgo CFLAGS: -I./yacapi/include -I./yacapi/src
#cgo !windows LDFLAGS: -ldl

#include "yacapi.h"
#include "yapi_inc.h"
#include "yacapi.go.h"
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"strings"
	"sync"
	"unsafe"
)

const (
	_LobBufLen      = 8192
	_OutputBindSize = 8192
	_DefaultSize    = 32*1024 + 1
)

type valueFreeType int8

var (
	mutex = &sync.Mutex{}

	yapiErr C.YapiErrorInfo

	byteBufferPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, _LobBufLen)
		},
	}

	keySqls = []string{
		"create or replace procedure",
		"create procedure",
		"create or replace trigger",
		"create trigger",
		"create or replace function",
		"create function",
		"create or replace package",
		"create or replace editionable package",
		"create package",
		"create or replace package body",
		"create package body",
		"create or replace type body",
		"create type body",
		"begin",
		"declare",
	}

	notFree    valueFreeType = 0
	normalFree valueFreeType = 1
	lobFree    valueFreeType = 2

	commentRegStr_1 = `\/\*([^*]|\*+[^*/])*\*+\/`
	commentRegStr_2 = `^--.*`
)

type bindStruct struct {
	direction C.YapiParamDirection
	yacType   C.YapiType
	value     C.YapiPointer
	bindSize  C.int32_t
	bufLength C.int32_t
	indicator *C.int32_t
	out       sql.Out
	freeType  valueFreeType
}

func stringToYasChar(str string) *C.char {
	p := C.malloc(C.size_t(len(str) + 1))
	pp := (*[1 << 30]byte)(p)
	copy(pp[:], str)
	pp[len(str)] = 0
	return (*C.char)(p)
}

func intToYacInt16(n int) C.int16_t {
	return C.int16_t(n)
}

func intToYacUint16(n int) C.uint16_t {
	return C.uint16_t(n)
}

func intToYacInt32(n int) C.int32_t {
	return C.int32_t(n)
}

func intToYacUint32(n int) C.uint32_t {
	return C.uint32_t(n)
}

func intToYacInt(n int) C.int {
	return C.int(n)
}

func yacPointerToInt64(p C.YapiPointer) int64 {
	return int64(*(*C.int64_t)(p))
}

func yacPointerToUint64(p C.YapiPointer) uint64 {
	return uint64(*(*C.uint64_t)(p))
}

func yacPointerToFloat64(p C.YapiPointer) float64 {
	return float64(*(*C.double)(p))
}

func yacPointerToBool(p C.YapiPointer) bool {
	return bool(*(*C.bool)(p))
}

func mallocBytes(size uint32) unsafe.Pointer {
	p := C.malloc(C.size_t(size))
	pp := (*[1 << 30]byte)(p)
	return unsafe.Pointer(pp)
}

func sizeToAlign4(size uint32) uint32 {
	margin := uint32(size % 4)
	if margin == 0 {
		return size
	}
	return size + (4 - margin)
}

func freeFetchRows(rows []*yasRow) {
	if len(rows) == 0 || rows == nil {
		return
	}
	for i := 0; i < len(rows); i++ {
		if rows[i] == nil {
			continue
		}
		freeFetchRow(rows[i])
	}
}

func freeFetchRow(row *yasRow) {
	if row == nil {
		return
	}
	switch row.freeType {
	case lobFree:
		lobLocator := (**C.YapiLobLocator)(unsafe.Pointer(row.Data))
		C.yapiLobDescFree(unsafe.Pointer(*lobLocator), row.yacType)
	case normalFree:
		C.free(row.Data)
	}

	if row.Indicator != nil {
		C.free(unsafe.Pointer(row.Indicator))
	}
	row.Data = nil
	row.Indicator = nil
}

func checkYasError(ret C.YapiResult) error {
	if int(ret) == 0 {
		return nil
	}
	mutex.Lock()
	defer func() {
		yapiErr.errCode = -1
		yapiErr.pos.line = 0
		yapiErr.pos.column = 0
		yapiErr.message = nil
		yapiErr.sqlState = nil
		mutex.Unlock()
	}()

	C.yapiGetLastError(&yapiErr)
	err := &YasDBError{
		Code:     int(yapiErr.errCode),
		Msg:      C.GoString(yapiErr.message),
		SqlState: C.GoString(yapiErr.sqlState),
		Line:     int(yapiErr.pos.line),
		Column:   int(yapiErr.pos.column),
	}
	return err
}

func tryRmSemicolon(query string) string {
	if isKeySql(query) {
		return query
	}
	return strings.TrimSuffix(strings.TrimSpace(query), ";")
}

func isKeySql(query string) bool {
	query = rmComment(query)
	strs := strings.Fields(strings.TrimSpace(query))
	sqlStr := strings.ToLower(strings.Join(strs, " "))
	for _, v := range keySqls {
		if strings.HasPrefix(sqlStr, v) {
			return true
		}
	}
	return false
}

func rmComment(query string) string {
	reg1, _ := regexp.Compile(commentRegStr_1)
	if reg1.MatchString(query) {
		query = reg1.ReplaceAllString(query, "")
	}

	nQuery := ""
	reg2, _ := regexp.Compile(commentRegStr_2)
	for _, line := range strings.Split(query, "\n") {
		nline := strings.TrimSpace(line)
		if nline == "" {
			continue
		}
		if reg2.MatchString(nline) {
			continue
		}
		nQuery += nline + "\n"
	}
	return nQuery
}

func isDisconnetionErr(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	if strings.Contains(errStr, "YAS-08012") || strings.Contains(errStr, "YAS-00406") {
		return true
	}
	return false
}

func releaseConn(conn *C.YapiConnect) error {
	if conn == nil {
		return nil
	}
	if err := checkYasError(C.yapiReleaseConn(conn)); err != nil {
		return err
	}
	conn = nil
	return nil
}

func releaseEnv(env *C.YapiEnv) error {
	if env == nil {
		return nil
	}
	if err := checkYasError(C.yapiReleaseEnv(env)); err != nil {
		return err
	}
	env = nil
	return nil
}

func ConvertToNameValue(args ...any) ([]driver.NamedValue, error) {
	nargs := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		var (
			nargValue driver.Value
			err       error
		)
		outValue, isOut := arg.(sql.Out)
		if isOut {
			nargValue = outValue
		} else {
			nargValue, err = driver.DefaultParameterConverter.ConvertValue(arg)
			if err != nil {
				return nil, err
			}
		}
		nargs[i].Ordinal = i + 1
		nargs[i].Value = nargValue
	}
	return nargs, nil
}

func GetDatabaseTypeName(yapiType uint32) string {
	switch C.YapiType(yapiType) {
	case C.YAPI_TYPE_BOOL:
		return "BOOLEAN"
	case C.YAPI_TYPE_TINYINT:
		return "TINYINT"
	case C.YAPI_TYPE_SMALLINT:
		return "SMALLINT"
	case C.YAPI_TYPE_INTEGER:
		return "INTEGER"
	case C.YAPI_TYPE_BIGINT:
		return "BIGINT"
	case C.YAPI_TYPE_FLOAT:
		return "FLOAT"
	case C.YAPI_TYPE_DOUBLE:
		return "DOUBLE"
	case C.YAPI_TYPE_NUMBER:
		return "NUMBER"
	case C.YAPI_TYPE_DATE:
		return "DATE"
	case C.YAPI_TYPE_SHORTTIME:
		return "TIME"
	case C.YAPI_TYPE_TIMESTAMP:
		return "TIMESTAMP"
	case C.YAPI_TYPE_CHAR:
		return "CHAR"
	case C.YAPI_TYPE_NCHAR:
		return "NCHAR"
	case C.YAPI_TYPE_VARCHAR:
		return "VARCHAR"
	case C.YAPI_TYPE_NVARCHAR:
		return "NVARCHAR"
	case C.YAPI_TYPE_CLOB:
		return "CLOB"
	case C.YAPI_TYPE_BLOB:
		return "BLOB"
	case C.YAPI_TYPE_BINARY:
		return "RAW"
	case C.YAPI_TYPE_ROWID:
		return "ROWID"
	case C.YAPI_TYPE_BIT:
		return "BIT"
	case C.YAPI_TYPE_NCLOB:
		return "NCLOB"
	case C.YAPI_TYPE_JSON:
		return "JSON"
	case C.YAPI_TYPE_YM_INTERVAL:
		return "INTERVAL YEAR TO MONTH"
	case C.YAPI_TYPE_DS_INTERVAL:
		return "INTERVAL DAY TO SECOND"
	default:
		return ""
	}
}

func GetDatabaseTypeSize(yType C.YapiType) int32 {
	switch yType {
	case C.YAPI_TYPE_BOOL, C.YAPI_TYPE_TINYINT, C.YAPI_TYPE_UTINYINT:
		return 1
	case C.YAPI_TYPE_SMALLINT, C.YAPI_TYPE_USMALLINT:
		return 2
	case C.YAPI_TYPE_INTEGER, C.YAPI_TYPE_UINTEGER, C.YAPI_TYPE_FLOAT:
		return 4
	case C.YAPI_TYPE_BIGINT, C.YAPI_TYPE_DOUBLE, C.YAPI_TYPE_UBIGINT:
		return 8
	case C.YAPI_TYPE_NUMBER:
		return 22
	case C.YAPI_TYPE_DATE, C.YAPI_TYPE_SHORTDATE, C.YAPI_TYPE_SHORTTIME, C.YAPI_TYPE_TIMESTAMP, C.YAPI_TYPE_TIMESTAMP_TZ, C.YAPI_TYPE_TIMESTAMP_LTZ, C.YAPI_TYPE_YM_INTERVAL, C.YAPI_TYPE_DS_INTERVAL:
		return 12
	case C.YAPI_TYPE_CHAR, C.YAPI_TYPE_NCHAR, C.YAPI_TYPE_VARCHAR, C.YAPI_TYPE_NVARCHAR, C.YAPI_TYPE_BINARY, C.YAPI_TYPE_CLOB, C.YAPI_TYPE_BLOB, C.YAPI_TYPE_BIT:
		return -1
	case C.YAPI_TYPE_ROWID:
		return 44
	default:
		return _DefaultSize
	}
}
