/*
Copyright  2022, YashanDB and/or its affiliates. All rights reserved.
YashanDB Driver for golang is licensed under the terms of the mulan PSL v2.0

License: 	http://license.coscl.org.cn/MulanPSL2
Home page: 	https://www.yashandb.com/
*/

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
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

// Number 是 YashanDB NUMBER 类型的无损文本包装。
// 零值 Number{} 表示 NULL。
type Number struct {
	str   string
	valid bool
}

// NewNumber 创建一个有效的 Number，s 为 NUMBER 的十进制文本表示。
func NewNumber(s string) Number {
	return Number{str: s, valid: true}
}

// String 返回 Number 的十进制文本表示。NULL 时返回空字符串。
func (n Number) String() string {
	return n.str
}

// Valid 返回 Number 是否有效（非 NULL）。
func (n Number) Valid() bool {
	return n.valid
}

// IsNull 返回 Number 是否为 NULL。
func (n Number) IsNull() bool {
	return !n.valid
}

// Text 返回文本和有效性。
func (n Number) Text() (string, bool) {
	return n.str, n.valid
}

// Int64 将 Number 转换为 int64。NULL 或非整数文本时返回错误。
func (n Number) Int64() (int64, error) {
	if !n.valid {
		return 0, errors.New("yasdb: cannot convert NULL Number to int64")
	}
	return strconv.ParseInt(n.str, 10, 64)
}

// Float64 将 Number 转换为 float64。NULL 或非法文本时返回错误。
func (n Number) Float64() (float64, error) {
	if !n.valid {
		return 0, errors.New("yasdb: cannot convert NULL Number to float64")
	}
	return strconv.ParseFloat(n.str, 64)
}

// Scan 实现 sql.Scanner，从数据库或 Go 值中读取 NUMBER。
// 支持的源类型：nil、string、[]byte、Number、*Number、int64、int、int32、float64、float32。
func (n *Number) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		n.str = ""
		n.valid = false
	case string:
		n.str = v
		n.valid = true
	case []byte:
		n.str = string(v)
		n.valid = true
	case Number:
		*n = v
	case *Number:
		if v == nil {
			n.str = ""
			n.valid = false
		} else {
			*n = *v
		}
	case int64:
		n.str = strconv.FormatInt(v, 10)
		n.valid = true
	case int:
		n.str = strconv.FormatInt(int64(v), 10)
		n.valid = true
	case int32:
		n.str = strconv.FormatInt(int64(v), 10)
		n.valid = true
	case float64:
		n.str = strconv.FormatFloat(v, 'f', -1, 64)
		n.valid = true
	case float32:
		n.str = strconv.FormatFloat(float64(v), 'f', -1, 32)
		n.valid = true
	default:
		return fmt.Errorf("yasdb: cannot scan %T into Number", src)
	}
	return nil
}

// Value 实现 driver.Valuer。NULL 返回 nil；有效值返回十进制文本字符串。
func (n Number) Value() (driver.Value, error) {
	if !n.valid {
		return nil, nil
	}
	return n.str, nil
}

// MarshalJSON 实现 json.Marshaler。NULL 序列化为 null，有效值序列化为 JSON 字符串（如 "123.45"），避免 JSON 数字消费者丢失精度。
func (n Number) MarshalJSON() ([]byte, error) {
	if !n.valid {
		return json.Marshal(nil)
	}
	return json.Marshal(n.str)
}

// UnmarshalJSON 实现 json.Unmarshaler。支持 null、JSON 字符串和 JSON 数字。
func (n *Number) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.str = ""
		n.valid = false
		return nil
	}
	// 尝试 JSON 字符串
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		n.str = s
		n.valid = true
		return nil
	}
	// 尝试 JSON 数字：直接保存原始文本，不经过 float64 以避免精度丢失。
	// json.Number 验证确保是合法数字格式。
	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		n.str = num.String()
		n.valid = true
		return nil
	}
	return fmt.Errorf("yasdb: cannot unmarshal %s into Number", string(data))
}

// MarshalText 实现 encoding.TextMarshaler。NULL 时返回 nil。
func (n Number) MarshalText() ([]byte, error) {
	if !n.valid {
		return nil, nil
	}
	return []byte(n.str), nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler。
func (n *Number) UnmarshalText(text []byte) error {
	n.str = string(text)
	n.valid = true
	return nil
}

// ---------------------------------------------------------------------------
// Number ↔ VARCHAR bind buffer conversion helpers
// ---------------------------------------------------------------------------

// _numberFmt is the format model for yapiNumberToText and yapiNumberFromText.
// 38 nines + dot + 25 nines = 64 chars. C API limit: int digits + dec digits ≤ 63.
const _numberFmt = "99999999999999999999999999999999999999.9999999999999999999999999"

// yapiNumberToString converts a C.YapiNumber to Go string using yapiNumberToText.
// Used for pure OUT NUMBER binding where the database writes YapiNumber binary format.
func yapiNumberToString(number *C.YapiNumber) (string, error) {
	const bufLen = 128
	buf := (*C.char)(mallocBytes(uint32(bufLen)))
	defer C.free(unsafe.Pointer(buf))

	fmtStr := C.CString(_numberFmt)
	defer C.free(unsafe.Pointer(fmtStr))

	var outLen C.int32_t
	if err := yapiNumberToText(number, fmtStr, C.uint32_t(len(_numberFmt)), nil, 0, buf, C.int32_t(bufLen), &outLen); err != nil {
		return "", err
	}
	result := strings.TrimSpace(C.GoStringN(buf, C.int(outLen)))
	if idx := strings.IndexByte(result, '.'); idx >= 0 {
		result = strings.TrimRight(result, "0")
		result = strings.TrimRight(result, ".")
	}
	if result == "" || result == "-" {
		result = "0"
	}
	return result, nil
}

// stringToYapiNumber converts a Go string to a C.YapiNumber using yapiNumberFromText.
// fmtStr must be non-NULL (pass &fmtStr) — passing nil causes YAS-08045.
func stringToYapiNumber(s string, fmtStr *string) (*C.YapiNumber, error) {
	var number C.YapiNumber
	p := C.malloc(C.size_t(unsafe.Sizeof(number)))
	b := unsafe.Slice((*byte)(p), unsafe.Sizeof(number))
	for i := range b {
		b[i] = 0
	}
	np := (*C.YapiNumber)(p)

	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))

	var cfmt *C.char
	var cfmtLen C.uint32_t
	if fmtStr != nil {
		cfmt = C.CString(*fmtStr)
		defer C.free(unsafe.Pointer(cfmt))
		cfmtLen = C.uint32_t(len(*fmtStr))
	}

	if err := yapiNumberFromText(cstr, C.uint32_t(len(s)), cfmt, cfmtLen, nil, 0, np); err != nil {
		C.free(p)
		return nil, err
	}
	return np, nil
}

// numberDestToYapiNumber converts a Go value to a C.YapiNumber for IN OUT bind.
// Uses text path (yapiNumberFromText) to avoid float64 precision loss.
func numberDestToYapiNumber(dest interface{}, in bool) (number *C.YapiNumber, indicator C.int32_t, err error) {
	var num C.YapiNumber
	bufSize := C.int32_t(unsafe.Sizeof(num))

	allocBuf := func() *C.YapiNumber {
		p := C.malloc(C.size_t(unsafe.Sizeof(num)))
		b := unsafe.Slice((*byte)(p), unsafe.Sizeof(num))
		for i := range b {
			b[i] = 0
		}
		return (*C.YapiNumber)(p)
	}

	if !in {
		return allocBuf(), bufSize, nil
	}

	// IN OUT: convert Go value to string, then to YapiNumber via text path.
	var s string
	isNull := false
	switch v := dest.(type) {
	case *float64:
		var number C.YapiNumber
		p := C.malloc(C.size_t(unsafe.Sizeof(number)))
		np := (*C.YapiNumber)(p)
		yp := C.YapiPointer(unsafe.Pointer(v))
		length := C.uint32_t(unsafe.Sizeof(*v))
		if err := yapiNumberFromReal(yp, length, np); err != nil {
			C.free(p)
			return nil, 0, err
		}
		return np, bufSize, nil

	case *int64:
		s = strconv.FormatInt(*v, 10)
	case *string:
		if *v == "" {
			isNull = true
		} else {
			s = *v
		}
	case *Number:
		if !v.Valid() {
			isNull = true
		} else {
			s = v.String()
		}
	default:
		return nil, 0, NewBindOutDestTypeErr("*float64, *int64, *string, or *yasdb.Number")
	}

	if isNull {
		return allocBuf(), C.YAPI_NULL_DATA, nil
	}

	fmtStr := _numberFmt
	num2, err2 := stringToYapiNumber(s, &fmtStr)
	if err2 != nil {
		return nil, 0, fmt.Errorf("yasdb: cannot convert %q to NUMBER: %w", s, err2)
	}
	return num2, bufSize, nil
}

// numberToInBindValue prepares bind parameters for a Number IN parameter.
// Binds as VARCHAR — database handles VARCHAR → NUMBER conversion.
func numberToInBindValue(n Number) *bindInfo {
	if !n.Valid() {
		return &bindInfo{
			yacType:   C.YAPI_TYPE_CHAR,
			bindSize:  0,
			indicator: func() *C.int32_t { v := C.int32_t(C.YAPI_NULL_DATA); return &v }(),
			freeType:  normalFree,
		}
	}
	s := n.String()
	return &bindInfo{
		yacType:   C.YAPI_TYPE_VARCHAR,
		bindSize:  C.int32_t(len(s)) + 1,
		bufLength: C.int32_t(len(s)),
		value:     C.YapiPointer(unsafe.Pointer(stringToYasChar(s))),
		freeType:  normalFree,
	}
}

// setNumberDestFromString parses a number text string and writes the converted
// value into the Go destination pointer (*float64, *int64, *string, or *Number).
func setNumberDestFromString(dest interface{}, text string, isNull bool) error {
	if isNull {
		switch v := dest.(type) {
		case *float64:
			*v = 0
		case *int64:
			*v = 0
		case *string:
			*v = ""
		case *Number:
			*v = Number{}
		default:
			return NewBindOutDestTypeErr("*float64, *int64, *string, or *yasdb.Number")
		}
		return nil
	}

	switch v := dest.(type) {
	case *float64:
		res, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return fmt.Errorf("yasdb: cannot convert NUMBER %q to float64: %w", text, err)
		}
		*v = res
		return nil

	case *int64:
		n, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return fmt.Errorf("yasdb: cannot convert NUMBER %q to int64: %w", text, err)
		}
		*v = n
		return nil

	case *string:
		*v = text
		return nil

	case *Number:
		*v = NewNumber(text)
		return nil

	default:
		return NewBindOutDestTypeErr("*float64, *int64, *string, or *yasdb.Number")
	}
}
