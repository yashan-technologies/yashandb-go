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
	"database/sql/driver"
	"io"
	"math"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

type yasRow struct {
	Elements   uint32
	Size       uint32
	IsValueSet bool
	IsArray    bool
	Data       unsafe.Pointer
	Indicator  *C.int32_t
	yacType    C.YapiType
	name       string
	freeType   valueFreeType
}

func NewYasRow(size uint32, yacType C.YapiType) *yasRow {
	row := &yasRow{
		Elements: 1,
		yacType:  yacType,
		Size:     size,
	}
	return row
}

type YasRows struct {
	stmt      *YasStmt
	fetchRows []*yasRow
	isClosed  bool
}

// Columns returns the names of the columns.
func (r *YasRows) Columns() []string {
	names := make([]string, 0, len(r.fetchRows))
	for _, row := range r.fetchRows {
		names = append(names, row.name)
	}
	return names
}

// Close closes the rows iterator.
func (r *YasRows) Close() error {
	if r.isClosed {
		return nil
	}
	freeFetchRows(r.fetchRows)
	r.isClosed = true
	return nil
}

// Next is called to populate the next row of data into
// the provided slice. The provided slice will be the same
// size as the Columns() are wide.
//
// Next should return io.EOF when there are no more rows.
//
// The dest should not be written to outside of Next. Care
// should be taken when closing Rows not to modify
// a buffer held in dest.
func (r *YasRows) Next(dest []driver.Value) error {
	if r.isClosed {
		return nil
	}
	if r.stmt.ctx.Err() != nil {
		return r.stmt.ctx.Err()
	}
	r.stmt.Lock()
	defer r.stmt.Unlock()

	done := make(chan struct{})
	defer close(done)
	go r.stmt.Conn.handleYacCancel(r.stmt.ctx, done)

	results, err := r.getValues()
	if err != nil {
		return err
	}
	if results == nil {
		return io.EOF
	}
	for i, d := range *results {
		dest[i] = d
	}
	return nil
}

// ColumnTypeScanType return the value type that can be used to scan types into.
// For example, the database column type "bigint" this should return "reflect.TypeOf(int64(0))".
func (r *YasRows) ColumnTypeScanType(index int) reflect.Type {
	if len(r.fetchRows) < index+1 {
		return reflect.TypeOf(nil)
	}
	switch r.fetchRows[index].yacType {
	case C.YAPI_TYPE_BOOL:
		return reflect.TypeOf(false)
	case C.YAPI_TYPE_TINYINT:
		return reflect.TypeOf(int8(0))
	case C.YAPI_TYPE_SMALLINT:
		return reflect.TypeOf(int16(0))
	case C.YAPI_TYPE_INTEGER:
		return reflect.TypeOf(int32(0))
	case C.YAPI_TYPE_BIGINT:
		return reflect.TypeOf(int64(0))
	case C.YAPI_TYPE_FLOAT:
		return reflect.TypeOf(float32(0))
	case C.YAPI_TYPE_DOUBLE, C.YAPI_TYPE_NUMBER:
		return reflect.TypeOf(float64(0))
	case C.YAPI_TYPE_DATE, C.YAPI_TYPE_TIMESTAMP:
		return reflect.TypeOf(time.Time{})
	case C.YAPI_TYPE_CHAR, C.YAPI_TYPE_NCHAR, C.YAPI_TYPE_VARCHAR, C.YAPI_TYPE_NVARCHAR, C.YAPI_TYPE_CLOB:
		return reflect.TypeOf("")
	case C.YAPI_TYPE_BLOB:
		return reflect.TypeOf([]byte(nil))
	default:
		return reflect.TypeOf(nil)
	}
}

// RowsColumnTypeDatabaseTypeName return the database system type name without the length. Type names should be uppercase.
func (r *YasRows) ColumnTypeDatabaseTypeName(index int) string {
	if len(r.fetchRows) < index+1 {
		return ""
	}
	switch r.fetchRows[index].yacType {
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
	default:
		return ""
	}
}

// RowsColumnTypeLength return the length of the column type if the column is a variable length type.
// If the column is not a variable length type ok should return false.
// If length is not limited other than system limits, it should return math.MaxInt64.
func (r *YasRows) ColumnTypeLength(index int) (length int64, ok bool) {
	if len(r.fetchRows) < index+1 {
		return 0, false
	}
	switch r.fetchRows[index].yacType {
	case C.YAPI_TYPE_CHAR, C.YAPI_TYPE_NCHAR, C.YAPI_TYPE_VARCHAR, C.YAPI_TYPE_NVARCHAR:
		return int64(r.fetchRows[index].Size), true
	case C.YAPI_TYPE_BLOB, C.YAPI_TYPE_CLOB:
		return math.MaxInt64, true
	default:
		return 0, false
	}
}

func (r *YasRows) getValues() (*[]driver.Value, error) {
	var err error
	unsafeRows := (unsafe.Pointer)(new(uint32))
	rows := (*C.uint32_t)(unsafeRows)
	if err = checkYasError(C.yapiFetch(r.stmt.Stmt, rows)); err != nil {
		return nil, err
	}
	if *rows == 0 {
		return nil, nil
	}
	columns := len(r.fetchRows)
	dest := make([]driver.Value, columns)
	for i := 0; i < columns; i++ {
		row := r.fetchRows[i]
		if *row.Indicator == C.YAPI_NULL_DATA {
			dest[i] = nil
			continue
		}
		var value driver.Value
		switch row.yacType {
		case C.YAPI_TYPE_BOOL:
			value = (*(*bool)(row.Data))
		case C.YAPI_TYPE_TINYINT:
			value = (*(*int8)(row.Data))
		case C.YAPI_TYPE_SMALLINT:
			value = (*(*int16)(row.Data))
		case C.YAPI_TYPE_INTEGER:
			value = (*(*int32)(row.Data))
		case C.YAPI_TYPE_BIGINT:
			value = (*(*int64)(row.Data))
		case C.YAPI_TYPE_FLOAT:
			value = (*(*float32)(row.Data))
		case C.YAPI_TYPE_DOUBLE:
			value = (*(*float64)(row.Data))
		case C.YAPI_TYPE_DATE, C.YAPI_TYPE_TIMESTAMP, C.YAPI_TYPE_SHORTDATE, C.YAPI_TYPE_SHORTTIME:
			value = time.Unix(0, (*(*int64)(row.Data))*1e3)
		case C.YAPI_TYPE_CHAR, C.YAPI_TYPE_NCHAR, C.YAPI_TYPE_VARCHAR, C.YAPI_TYPE_NVARCHAR, C.YAPI_TYPE_YM_INTERVAL, C.YAPI_TYPE_DS_INTERVAL:
			value = (C.GoString((*C.char)(row.Data)))
		case C.YAPI_TYPE_NUMBER:
			value, err = strconv.ParseFloat(C.GoString((*C.char)(row.Data)), 64)
			if err != nil {
				return nil, err
			}
		case C.YAPI_TYPE_CLOB, C.YAPI_TYPE_BLOB:
			lobLocator := (**C.YapiLobLocator)(row.Data)
			data, err := r.stmt.Conn.lobRead(*lobLocator)
			if err != nil {
				return nil, err
			}
			if row.yacType == C.YAPI_TYPE_CLOB {
				value = string(data)
			} else {
				value = data
			}
		default:
			value = (C.GoString((*C.char)(row.Data)))
		}
		dest[i] = value
	}
	return &dest, nil
}
