/*
Copyright  2022, YashanDB and/or its affiliates. All rights reserved.
YashanDB Driver for golang is licensed under the terms of the mulan PSL v2.0

License: 	http://license.coscl.org.cn/MulanPSL2
Home page: 	https://www.yashandb.com/
*/

package yasdb

/*
#cgo CFLAGS: -I./yacapi/include -I./yacapi/src
#cgo LDFLAGS: -ldl

#include "yacapi.h"
#include "yapi_inc.h"
#include "yacapi.go.h"
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
    "database/sql"
    "strings"
    "sync"
    "unsafe"
)

const (
    _LobBufLen      = 8192
    _OutputBindSize = 8192
)

var (
    mutex = &sync.Mutex{}

    byteBufferPool = sync.Pool{
        New: func() interface{} {
            return make([]byte, _LobBufLen)
        },
    }
)

type bindStruct struct {
    direction C.YapiParamDirection
    yacType   C.YapiType
    value     C.YapiPointer
    bindSize  C.int32_t
    bufLength C.int32_t
    indicator *C.int32_t
    out       sql.Out
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

func intToYacUint32(n int) C.uint32_t {
    return C.uint32_t(n)
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
        switch rows[i].yacType {
        case C.YAPI_TYPE_CLOB, C.YAPI_TYPE_BLOB:
            lobLocator := (**C.YapiLobLocator)(unsafe.Pointer(rows[i].Data))
            C.yapiLobDescFree(unsafe.Pointer(*lobLocator), rows[i].yacType)
        default:
            C.free(rows[i].Data)
        }

        if rows[i].Indicator != nil {
            C.free(unsafe.Pointer(rows[i].Indicator))
        }
        rows[i].Data = nil
        rows[i].Indicator = nil
    }
}

func checkYasError(ret C.YapiResult) error {
    if int(ret) == 0 {
        return nil
    }
    mutex.Lock()
    defer mutex.Unlock()
    var yapErr C.YapiErrorInfo
    C.yapiGetLastError(&yapErr)
    err := &YasDBError{
        Code:     int(yapErr.errCode),
        Msg:      C.GoString(yapErr.message),
        SqlState: C.GoString(yapErr.sqlState),
        Line:     int(yapErr.pos.line),
        Column:   int(yapErr.pos.column),
    }
    return err
}

func rmSqlSemicolon(query string) string {
    return strings.TrimSuffix(strings.TrimSpace(query), ";")
}
