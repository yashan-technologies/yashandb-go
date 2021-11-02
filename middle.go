package yasdb

/*
#cgo CFLAGS: -I./anchor/src/cli
#cgo LDFLAGS: -L${SRCDIR}/anchor/build/lib -lcodcommon -lyas_infra -lyascli

#include "anc.h"
*/
import "C"
import (
	"bytes"
	"database/sql/driver"
	"log"
	"strings"
	"time"
	"unsafe"
)

type AncHandle *C.AncHandle

func NewAncHandle() AncHandle {
    return (AncHandle)((unsafe.Pointer)(new([]byte)))
}

type AncPointer C.AncPointer

var (
    ANC_HANDLE_UNKNOWN = 0
    ANC_HANDLE_ENV     = 1
    ANC_HANDLE_DBC     = 2
    ANC_HANDLE_STMT    = 3
    ANC_HANDLE_DESC    = 4
)

func checkYasError(ret C.AncResult) error {
    if int(ret) == 0 {
        return nil
    }
    errCode := (*C.AncInt32)((unsafe.Pointer)(new(int)))
    message := (*C.AncChar)((unsafe.Pointer)(new(string)))
    sqlState := (*C.AncChar)((unsafe.Pointer)(new(string)))

    pos := &C.struct_StAncTextPos{}
    C.ancGetLastError(errCode, &message, &sqlState, pos)
    err := &YasDBError{
        Code:     int(*errCode),
        Msg:      C.GoString(message),
        SqlState: C.GoString(sqlState),
        Line:     int(pos.line),
        Column:   int(pos.column),
    }
    log.Println(err)
    return err
}

func yasdbConnect(conn *Connection) error {
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

    if err := checkYasError(C.ancAllocHandle(C.ANC_HANDLE_ENV, nil, conn.Env)); err != nil {
        return err
    }
    if err := checkYasError(C.ancAllocHandle(C.ANC_HANDLE_DBC, *conn.Env, conn.Conn)); err != nil {
        yasdbFreeHandle(conn.Env, C.ANC_HANDLE_ENV)
        return err
    }
    if err := checkYasError(C.ancConnect(*conn.Conn, url, username, password)); err != nil {
        yasdbFreeHandle(conn.Conn, C.ANC_HANDLE_DBC)
        yasdbFreeHandle(conn.Env, C.ANC_HANDLE_ENV)
        return err
    }
    if err := checkYasError(C.ancAllocHandle(C.ANC_HANDLE_STMT, *conn.Conn, conn.Stmt)); err != nil {
        yasdbFreeHandle(conn.Conn, C.ANC_HANDLE_DBC)
        yasdbFreeHandle(conn.Env, C.ANC_HANDLE_ENV)
        return err
    }
    return getAutoCommit(conn)
}

func yasdbFreeHandle(a AncHandle, t int) error {
    if a == nil {
        return nil
    }
    if err := checkYasError(C.ancFreeHandle(C.AncHandleType(t), *a)); err != nil {
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
    if err := checkYasError(C.ancAllocHandle(C.ANC_HANDLE_STMT, *conn.Conn, stmt.Stmt)); err != nil {
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
    if err := checkYasError(C.ancPrepare(*stmt.Stmt, q)); err != nil {
        return err
    }
    sqltype := (unsafe.Pointer)(new(uint32))
    sqlSize := C.AncInt32(unsafe.Sizeof(&sqltype))
    err := checkYasError(C.ancGetStmtAttr(*stmt.Stmt, C.ANC_ATTR_SQLTYPE, sqltype, sqlSize))
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
        ancType, v, err := valueToC(arg)
        if err != nil {
            freeBindVals(stmt)
            return err
        }
        var ret C.AncResult
        var charV *C.char
        ancValue := C.AncPointer(v)
        size := C.AncInt32(unsafe.Sizeof(&arg)) + 1
        if ancType == C.ANC_TYPE_VARCHAR {
            charV = C.CString(*(*string)(v))
            size = C.AncInt32(bytes.Count([]byte(*(*string)(v)), nil))
            ancValue = C.AncPointer(charV)
            // hold for free in execute
            stmt.bindVals = append(stmt.bindVals, unsafe.Pointer(charV))
        }
        indicator := size - 1
        log.Println(ancType, arg, size, indicator)
        ret = C.ancBindParameter(*stmt.Stmt, C.AncUint16(i), C.ANC_PARAM_INPUT, ancType, ancValue, size, &indicator)
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
    columns := C.AncInt16(0)
    if err := checkYasError(C.ancNumResultCols(*stmt.Stmt, &columns)); err != nil {
        return err
    }
    if err := checkYasError(C.ancExecute(*stmt.Stmt)); err != nil {
        return err
    }
    if columns == 0 {
        if stmt.SqlType >= uint32(C.ANC_SQLTYPE_CREATE_DATABASE) {
            stmt.RowCount = 1
            return nil
        }
        rowCount := (unsafe.Pointer)(&stmt.RowCount)
        size := C.AncInt32(unsafe.Sizeof(new(int64)))
        return checkYasError(C.ancGetStmtAttr(*stmt.Stmt, C.ANC_ATTR_ROWS_AFFECTED, rowCount, size))
    } else {
        stmt.RowCount = 0
    }
    if columns > 0 {
        if err := yasdbColumns(stmt, C.AncInt32(columns)); err != nil {
            return err
        }
    }
    return nil
}

func yasdbColumns(stmt *YasStmt, columns C.AncInt32) error {
    cols := []string{}
    pos := C.AncInt32(0)
    for pos = 0; pos < columns; pos++ {
        item := C.struct_StAncColumnDesc{}
        if err := checkYasError(C.ancDescribeCol2(*stmt.Stmt, C.AncUint16(pos), &item)); err != nil {
            return err
        }
        cols = append(cols, C.GoString(item.name))

        size := uint32(item.size)
        row := NewYasRow(stmt, size, int(item._type))
        if err := checkYasError(
            C.ancBindColumn(
                *stmt.Stmt, C.AncUint16(pos), C.AncType(item._type),
                C.AncPointer(row.Data), C.AncInt32(size), nil)); err != nil {
            return err
        }
        stmt.fetchRows = append(stmt.fetchRows, row)
    }
    stmt.Columns = &cols
    return nil
}

func yasdbFetch(stmt *YasStmt) (*[]driver.Value, error) {
    rows := (*C.AncUint32)((unsafe.Pointer)(new(uint32)))
    if err := checkYasError(C.ancFetch(*stmt.Stmt, rows)); err != nil {
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
    var a C.AncInt32
    if auto {
        a = 1
    } else {
        a = 0
    }
    size := C.AncInt32(unsafe.Sizeof(a))
    return checkYasError(C.ancSetConnAttr(*conn.Conn, C.ANC_ATTR_AUTOCOMMIT, unsafe.Pointer(&a), size))
}

func getAutoCommit(conn *Connection) error {
    var auto C.AncInt32
    size := C.AncInt32(unsafe.Sizeof(auto))
    err := checkYasError(C.ancGetConnAttr(*conn.Conn, C.ANC_ATTR_AUTOCOMMIT, unsafe.Pointer(&auto), size))
    if err != nil {
        return err
    }
    if auto == 0 {
        conn.autoCommit = false
    } else {
        conn.autoCommit = true
    }
    return nil
}

func yasdbCommit(conn *Connection) error {
    return checkYasError(C.ancCommit(*conn.Conn))
}

func yasdbRollback(conn *Connection) error {
    return checkYasError(C.ancRollback(*conn.Conn))
}

func yasdbRowAffected(stmt *YasStmt) (int64, error) {
    var rows C.AncUint32
    size := C.AncInt32(unsafe.Sizeof(&rows))
    err := checkYasError(
        C.ancGetStmtAttr(
            *stmt.Stmt,
            C.ANC_ATTR_ROWS_AFFECTED,
            unsafe.Pointer(&rows), size),
    )
    return int64(rows), err
}

func codSizeAlign4(size C.AncUint32) uint32 {
    margin := uint32(size & 0x03)
    if margin == 0 {
        return uint32(size)
    }
    return uint32(size) + (4 - margin)
}

func valueToGolang(row *YasRow) (interface{}, error) {
    p := unsafe.Pointer(row.Data)
    switch C.AncType(row.DbType) {
    case C.ANC_TYPE_BOOL:
        return (*(*bool)(p)), nil
    case C.ANC_TYPE_TINYINT:
        return (*(*int8)(p)), nil
    case C.ANC_TYPE_SMALLINT:
        return (*(*int16)(p)), nil
    case C.ANC_TYPE_INTEGER:
        return (*(*int32)(p)), nil
    case C.ANC_TYPE_BIGINT:
        return (*(*int64)(p)), nil
    case C.ANC_TYPE_FLOAT:
        return (*(*float32)(p)), nil
    case C.ANC_TYPE_DOUBLE:
        return (*(*float64)(p)), nil
    case C.ANC_TYPE_DATE, C.ANC_TYPE_TIMESTAMP:
        date := time.Unix(0, (*(*int64)(p))*1e3)
        return date, nil
    case C.ANC_TYPE_CHAR, C.ANC_TYPE_NCHAR, C.ANC_TYPE_VARCHAR, C.ANC_TYPE_NVARCHAR:
        return (C.GoString((*C.char)(p))), nil
    }
    return nil, ErrDbTypeUnsupport(row.DbType)
}

func valueToC(arg driver.Value) (C.AncType, unsafe.Pointer, error) {
    switch v := arg.(type) {
    case int64:
        return C.ANC_TYPE_INTEGER, unsafe.Pointer(&v), nil
    case float64:
        return C.ANC_TYPE_DOUBLE, unsafe.Pointer(&v), nil
    case bool:
        return C.ANC_TYPE_BOOL, unsafe.Pointer(&v), nil
    case string:
        return C.ANC_TYPE_VARCHAR, unsafe.Pointer(&v), nil
    case []byte:
        return C.ANC_TYPE_BINARY, unsafe.Pointer(&v), nil
    case time.Time:
        // YashanDB 存储的是us，需要除以1000
        d := v.UnixNano() / 1e3
        return C.ANC_TYPE_TIMESTAMP, unsafe.Pointer(&d), nil
    default:
        return 0, nil, ErrUnknowType(v)
    }
}
