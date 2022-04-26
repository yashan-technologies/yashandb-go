package yasdb

/*
#cgo CFLAGS: -I./deps/include
#cgo LDFLAGS: -L${SRCDIR}/deps/lib -lcodcommon -lyas_infra -lyascli

#include "yacli.h"
*/
import "C"
import (
    "database/sql/driver"
    "strconv"
    "strings"
    "sync"
    "time"
    "unsafe"
)

type YacHandle *C.YacHandle

func NewYacHandle() YacHandle {
    return (YacHandle)((unsafe.Pointer)(new([]byte)))
}

type YacPointer C.YacPointer

var (
    YAC_HANDLE_UNKNOWN = 0
    YAC_HANDLE_ENV     = 1
    YAC_HANDLE_DBC     = 2
    YAC_HANDLE_STMT    = 3
    YAC_HANDLE_DESC    = 4
)

var (
    Mutex        = &sync.Mutex{}
    LastErrMsg   = (*C.YacChar)(C.malloc(512))
    LastSqlState = (*C.YacChar)(C.malloc(512))
)

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

func yasdbConnect(conn *Connection, autoCommit bool) error {
    if conn.Dsn == "" {
        return ErrDsnNoSet()
    }
    items := strings.Split(conn.Dsn, "@")
    url := C.CString(items[1])
    defer C.free(unsafe.Pointer(url))

    user := strings.Split(items[0], "/")
    username, password := C.CString(user[0]), C.CString(user[1])
    defer C.free(unsafe.Pointer(username))
    defer C.free(unsafe.Pointer(password))
    conn.Username = user[0]

    if err := checkYasError(C.yacAllocHandle(C.YAC_HANDLE_ENV, nil, conn.Env)); err != nil {
        return err
    }
    if err := checkYasError(C.yacAllocHandle(C.YAC_HANDLE_DBC, *conn.Env, conn.Conn)); err != nil {
        yasdbFreeHandle(conn.Env, C.YAC_HANDLE_ENV)
        return err
    }
    if err := checkYasError(C.yacConnect(*conn.Conn, url, username, password)); err != nil {
        yasdbFreeHandle(conn.Conn, C.YAC_HANDLE_DBC)
        yasdbFreeHandle(conn.Env, C.YAC_HANDLE_ENV)
        return err
    }
    if err := checkYasError(C.yacAllocHandle(C.YAC_HANDLE_STMT, *conn.Conn, conn.Stmt)); err != nil {
        yasdbFreeHandle(conn.Conn, C.YAC_HANDLE_DBC)
        yasdbFreeHandle(conn.Env, C.YAC_HANDLE_ENV)
        return err
    }
    if err := setAutoCommit(conn, autoCommit); err != nil {
        yasdbFreeHandle(conn.Conn, C.YAC_HANDLE_DBC)
        yasdbFreeHandle(conn.Env, C.YAC_HANDLE_ENV)
        yasdbFreeHandle(conn.Stmt, C.YAC_HANDLE_STMT)
        return err
    }
    return nil
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

func yasdbStmtInit(stmt *YasStmt) error {
    conn := stmt.Conn
    if conn == nil || conn.Conn == nil {
        return ErrNoConnect()
    }
    if err := checkYasError(C.yacAllocHandle(C.YAC_HANDLE_STMT, *conn.Conn, stmt.Stmt)); err != nil {
        return err
    }
    stmt.ArrSize = 100
    stmt.IsOpen = true
    return nil
}

func yasdbPrepare(stmt *YasStmt, query string) error {
    if !stmt.IsOpen {
        return ErrStmtNoOpen()
    }
    q := C.CString(query)
    defer C.free(unsafe.Pointer(q))
    if err := checkYasError(C.yacPrepare(*stmt.Stmt, q)); err != nil {
        return err
    }
    sqltype := (unsafe.Pointer)(new(uint32))
    sqlSize := C.YacInt32(unsafe.Sizeof(&sqltype))
    err := checkYasError(C.yacGetStmtAttr(*stmt.Stmt, C.YAC_ATTR_SQLTYPE, sqltype, sqlSize))
    if err != nil {
        return err
    }
    stmt.SqlType = *(*uint32)(sqltype)
    return nil
}

func freeBindVals(stmt *YasStmt) {
    if len(stmt.bindVals) == 0 {
        return
    }
    for _, v := range stmt.bindVals {
        if v != nil {
            C.free(v)
        }
    }
    stmt.bindVals = []unsafe.Pointer{}
}

func yasdbBindParams(stmt *YasStmt, args []driver.Value) error {
    for i, arg := range args {
        yacType, v, err := valueToC(arg)
        if err != nil {
            freeBindVals(stmt)
            return err
        }
        var ret C.YacResult
        var charV *C.char
        yacValue := C.YacPointer(v)
        size := C.YacInt32(unsafe.Sizeof(&arg)) + 1
        if yacType == C.YAC_TYPE_VARCHAR {
            charV = C.CString(*(*string)(v))
            size = C.YacInt32(len([]byte(*(*string)(v))) + 1)
            yacValue = C.YacPointer(charV)
            // hold for free in execute
            stmt.bindVals = append(stmt.bindVals, unsafe.Pointer(charV))
        }
        indicator := size - 1
        index := i + 1
        ret = C.yacBindParameter(*stmt.Stmt, C.YacUint16(index), C.YAC_PARAM_INPUT, yacType, yacValue, size, &indicator)
        err = checkYasError(ret)
        if err != nil {
            freeBindVals(stmt)
            return err
        }
    }
    return nil
}

func yasdbExecute(stmt *YasStmt) error {
    defer freeBindVals(stmt)
    columns := C.YacInt16(0)
    if err := checkYasError(C.yacNumResultCols(*stmt.Stmt, &columns)); err != nil {
        return err
    }
    if err := checkYasError(C.yacExecute(*stmt.Stmt)); err != nil {
        return err
    }
    if columns == 0 {
        if stmt.SqlType >= uint32(C.YAC_SQLTYPE_CREATE_DATABASE) {
            stmt.RowCount = 1
            return nil
        }
        rowCount := (unsafe.Pointer)(new(uint64))
        size := C.YacInt32(unsafe.Sizeof(new(int64)))
        err := checkYasError(C.yacGetStmtAttr(*stmt.Stmt, C.YAC_ATTR_ROWS_AFFECTED, rowCount, size))
        if err != nil {
            return err
        }
        stmt.RowCount = *(*uint64)(rowCount)
        return nil
    } else {
        stmt.RowCount = 0
    }
    if columns > 0 {
        if err := yasdbColumns(stmt, C.YacInt32(columns)); err != nil {
            return err
        }
    }
    return nil
}

func yasdbColumns(stmt *YasStmt, columns C.YacInt32) error {
    cols := []string{}
    pos := C.YacInt32(0)
    for pos = 0; pos < columns; pos++ {
        item := C.struct_StYacColumnDesc{}
        if err := checkYasError(C.yacDescribeCol2(*stmt.Stmt, C.YacUint16(pos), &item)); err != nil {
            return err
        }
        cols = append(cols, strings.ToLower(C.GoString(item.name)))

        yacType := C.YacType(item._type)
        size, indicator := uint32(item.size), C.YacInt32(0)

        // number to string
        if C.YAC_TYPE_NUMBER == yacType {
            yacType = C.YAC_TYPE_VARCHAR
            size = size + 8
        }
        row := NewYasRow(stmt, size, int(item._type))
        if err := checkYasError(
            C.yacBindColumn(
                *stmt.Stmt, C.YacUint16(pos), yacType,
                C.YacPointer(row.Data), C.YacInt32(size), &indicator)); err != nil {
            return err
        }
        row.Indicator = int32(indicator)
        stmt.fetchRows = append(stmt.fetchRows, row)
    }
    stmt.Columns = &cols
    return nil
}

func yasdbFetch(stmt *YasStmt) (*[]driver.Value, error) {
    rows := (*C.YacUint32)((unsafe.Pointer)(new(uint32)))
    if err := checkYasError(C.yacFetch(*stmt.Stmt, rows)); err != nil {
        return nil, err
    }

    if *rows == 0 {
        return nil, nil
    }
    dest := []driver.Value{}
    stmt.RowCount++
    for i := 0; i < len(*stmt.Columns); i++ {
        row := stmt.fetchRows[i]
        if row == nil {
            return &dest, nil
        }
        v, err := valueToGolang(row)
        if err != nil {
            return nil, err
        }
        dest = append(dest, v)
    }
    return &dest, nil
}

func setAutoCommit(conn *Connection, auto bool) error {
    var a C.YacInt32
    if auto {
        a = 1
    } else {
        a = 0
    }
    size := C.YacInt32(unsafe.Sizeof(a))
    return checkYasError(C.yacSetConnAttr(*conn.Conn, C.YAC_ATTR_AUTOCOMMIT, unsafe.Pointer(&a), size))
}

func getAutoCommit(conn *Connection) error {
    var auto C.YacInt32
    size := C.YacInt32(unsafe.Sizeof(auto))
    err := checkYasError(C.yacGetConnAttr(*conn.Conn, C.YAC_ATTR_AUTOCOMMIT, unsafe.Pointer(&auto), size))
    if err != nil {
        return err
    }
    if auto == 0 {
        conn.AutoCommit = false
    } else {
        conn.AutoCommit = true
    }
    return nil
}

func yasdbCommit(conn *Connection) error {
    return checkYasError(C.yacCommit(*conn.Conn))
}

func yasdbRollback(conn *Connection) error {
    return checkYasError(C.yacRollback(*conn.Conn))
}

func yasdbRowAffected(stmt *YasStmt) (int64, error) {
    var rows C.YacUint32
    size := C.YacInt32(unsafe.Sizeof(&rows))
    err := checkYasError(
        C.yacGetStmtAttr(
            *stmt.Stmt,
            C.YAC_ATTR_ROWS_AFFECTED,
            unsafe.Pointer(&rows), size),
    )
    return int64(rows), err
}

func codSizeAlign4(size C.YacUint32) uint32 {
    margin := uint32(size & 0x03)
    if margin == 0 {
        return uint32(size)
    }
    return uint32(size) + (4 - margin)
}

func valueToGolang(row *YasRow) (interface{}, error) {
    p := unsafe.Pointer(row.Data)
    switch C.YacType(row.DbType) {
    case C.YAC_TYPE_BOOL:
        return (*(*bool)(p)), nil
    case C.YAC_TYPE_TINYINT:
        return (*(*int8)(p)), nil
    case C.YAC_TYPE_SMALLINT:
        return (*(*int16)(p)), nil
    case C.YAC_TYPE_INTEGER:
        return (*(*int32)(p)), nil
    case C.YAC_TYPE_BIGINT:
        return (*(*int64)(p)), nil
    case C.YAC_TYPE_FLOAT:
        return (*(*float32)(p)), nil
    case C.YAC_TYPE_DOUBLE:
        return (*(*float64)(p)), nil
    case C.YAC_TYPE_DATE, C.YAC_TYPE_TIMESTAMP:
        date := time.Unix(0, (*(*int64)(p))*1e3)
        return date, nil
    case C.YAC_TYPE_CHAR, C.YAC_TYPE_NCHAR, C.YAC_TYPE_VARCHAR, C.YAC_TYPE_NVARCHAR:
        return (C.GoString((*C.char)(p))), nil
    case C.YAC_TYPE_NUMBER:
        return strconv.ParseFloat(C.GoString((*C.char)(p)), 64)
    }
    return nil, ErrDbTypeUnsupport(row.DbType)
}

func valueToC(arg driver.Value) (C.YacType, unsafe.Pointer, error) {
    switch v := arg.(type) {
    case int64:
        return C.YAC_TYPE_INTEGER, unsafe.Pointer(&v), nil
    case float64:
        return C.YAC_TYPE_DOUBLE, unsafe.Pointer(&v), nil
    case bool:
        return C.YAC_TYPE_BOOL, unsafe.Pointer(&v), nil
    case string:
        return C.YAC_TYPE_VARCHAR, unsafe.Pointer(&v), nil
    case []byte:
        return C.YAC_TYPE_VARCHAR, unsafe.Pointer(&v), nil
    case time.Time:
        // YashanDB 存储的是us，需要除以1000
        d := v.UnixNano() / 1e3
        return C.YAC_TYPE_TIMESTAMP, unsafe.Pointer(&d), nil
    case nil:
        return C.YAC_TYPE_VARCHAR, unsafe.Pointer(&v), nil
    default:
        return 0, nil, ErrUnknowType(v)
    }
}
