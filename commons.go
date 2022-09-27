package yasdb

/*
#cgo !noPkgConfig pkg-config: yacli
#include "yacli.go.h"
*/
import "C"
import (
    "database/sql"
    "sync"
    "unsafe"
)

const (
    _LobBufLen      = 8192
    _OutputBindSize = 8192
)

var (
    YAC_HANDLE_UNKNOWN = 0
    YAC_HANDLE_ENV     = 1
    YAC_HANDLE_DBC     = 2
    YAC_HANDLE_STMT    = 3
    YAC_HANDLE_DESC    = 4

    Mutex        = &sync.Mutex{}
    LastErrMsg   = (*C.YacChar)(C.malloc(512))
    LastSqlState = (*C.YacChar)(C.malloc(512))

    byteBufferPool = sync.Pool{
        New: func() interface{} {
            return make([]byte, _LobBufLen)
        },
    }
)

type YacHandle *C.YacHandle

type bindStruct struct {
    direction C.YacParamDirection
    yacType   C.YacType
    value     C.YacPointer
    bindSize  C.YacUint32
    bufLength C.YacInt32
    indicator *C.YacInt32
    out       sql.Out
}

func NewYacHandle() YacHandle {
    return (YacHandle)((unsafe.Pointer)(new([]byte)))
}

// stringToYasChar converts golang string to C *YacChar
func stringToYasChar(str string) *C.YacChar {
    p := C.malloc(C.size_t(len(str) + 1))
    pp := (*[1 << 30]byte)(p)
    copy(pp[:], str)
    pp[len(str)] = 0
    return (*C.YacChar)(p) //C.CString(str)
}

// stringToYasUint8 converts golang string to C *YacUint8
func stringToYasUint8(str string) *C.YacUint8 {
    p := C.malloc(C.size_t(len(str) + 1))
    pp := (*[1 << 30]byte)(p)
    copy(pp[:], str)
    pp[len(str)] = 0
    return (*C.YacUint8)(p) //C.CString(str)
}

// intToYacInt32 converts golang int to C YacInt16
func intToYacInt16(n int) C.YacInt16 {
    return C.YacInt16(n)
}

// intToYacInt32 converts golang int to C YacInt32
func intToYacInt32(n int) C.YacInt32 {
    return C.YacInt32(n)
}

// intToYacUint16 converts golang int to C YacUint16
func intToYacUint16(n int) C.YacUint16 {
    return C.YacUint16(n)
}

// intToYacUint32 converts golang int to C YacUint32
func intToYacUint32(n int) C.YacUint32 {
    return C.YacUint32(n)
}

func yacPointerToInt64(p C.YacPointer) int64 {
    return int64(*(*C.YacInt64)(p))
}

func yacPointerToUint64(p C.YacPointer) uint64 {
    return uint64(*(*C.YacUint64)(p))
}

func yacPointerToFloat64(p C.YacPointer) float64 {
    return float64(*(*C.YacDouble)(p))
}

func yacPointerToBool(p C.YacPointer) bool {
    return bool(*(*C.YacBool)(p))
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

func yasdbFreeHandle(a YacHandle, t int) error {
    if a == nil {
        return nil
    }
    if err := checkYasError(C.yacFreeHandle(C.YacHandleType(t), *a)); err != nil {
        return err
    }
    a = nil
    return nil
}

func freeFetchRows(rows []*yasRow) {
    if len(rows) == 0 || rows == nil {
        return
    }
    for i := range rows {
        row := rows[i]
        if row == nil {
            continue
        }
        switch row.yacType {
        case C.YAC_TYPE_CLOB, C.YAC_TYPE_BLOB:
            lobLocator := (**C.YacLobLocator)(unsafe.Pointer(row.Data))
            C.yacLobDescFree(unsafe.Pointer(*lobLocator), row.yacType)
        default:
            C.free(row.Data)
        }
        rows[i].Data = nil
    }
}

func checkYasError(ret C.YacResult) error {
    if int(ret) == 0 {
        return nil
    }
    Mutex.Lock()
    defer Mutex.Unlock()
    var errCode C.YacInt32
    pos := &C.struct_StYacTextPos{}
    C.yacGetLastError(&errCode, &LastErrMsg, &LastSqlState, pos)
    err := &YasDBError{
        Code:     int(errCode),
        Msg:      C.GoString(LastErrMsg),
        SqlState: C.GoString(LastSqlState),
        Line:     int(pos.line),
        Column:   int(pos.column),
    }
    return err
}
