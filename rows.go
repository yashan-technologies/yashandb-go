package yasdb

/*
#cgo !noPkgConfig pkg-config: yacli
#include "yacli.go.h"
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
    Indicator  int32
    yacType    C.YacType
    name       string
}

func NewYasRow(stmt *YasStmt, size uint32, yacType C.YacType) *yasRow {
    row := &yasRow{
        Elements: 1,
        Size:     size,
        yacType:  yacType,
    }
    row.Data = mallocBytes(size)
    row.Indicator = 0
    return row
}

type YasRows struct {
    stmt      *YasStmt
    fetchRows []*yasRow
    isClosed  bool
}

func (r *YasRows) Columns() []string {
    names := make([]string, 0, len(r.fetchRows))
    for _, row := range r.fetchRows {
        names = append(names, row.name)
    }
    return names
}

func (r *YasRows) Close() error {
    if r.isClosed {
        return nil
    }
    freeFetchRows(r.fetchRows)
    r.isClosed = true
    return nil
}

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

func (r *YasRows) ColumnTypeScanType(index int) reflect.Type {
    if len(r.fetchRows) < index+1 {
        return reflect.TypeOf(nil)
    }
    switch r.fetchRows[index].yacType {
    case C.YAC_TYPE_BOOL:
        return reflect.TypeOf(false)
    case C.YAC_TYPE_TINYINT:
        return reflect.TypeOf(int8(0))
    case C.YAC_TYPE_SMALLINT:
        return reflect.TypeOf(int16(0))
    case C.YAC_TYPE_INTEGER:
        return reflect.TypeOf(int(0))
    case C.YAC_TYPE_BIGINT:
        return reflect.TypeOf(int64(0))
    case C.YAC_TYPE_FLOAT:
        return reflect.TypeOf(float32(0))
    case C.YAC_TYPE_DOUBLE, C.YAC_TYPE_NUMBER:
        return reflect.TypeOf(float64(0))
    case C.YAC_TYPE_DATE, C.YAC_TYPE_TIMESTAMP:
        return reflect.TypeOf(time.Time{})
    case C.YAC_TYPE_CHAR, C.YAC_TYPE_NCHAR, C.YAC_TYPE_VARCHAR, C.YAC_TYPE_NVARCHAR, C.YAC_TYPE_CLOB:
        return reflect.TypeOf("")
    case C.YAC_TYPE_BLOB:
        return reflect.TypeOf([]byte(nil))
    default:
        return reflect.TypeOf(nil)
    }
}

func (r *YasRows) ColumnTypeDatabaseTypeName(index int) string {
    if len(r.fetchRows) < index+1 {
        return ""
    }
    switch r.fetchRows[index].yacType {
    case C.YAC_TYPE_BOOL:
        return "BOOLEAN"
    case C.YAC_TYPE_TINYINT:
        return "TINYINT"
    case C.YAC_TYPE_SMALLINT:
        return "SMALLINT"
    case C.YAC_TYPE_INTEGER:
        return "INTEGER"
    case C.YAC_TYPE_BIGINT:
        return "BIGINT"
    case C.YAC_TYPE_FLOAT:
        return "FLOAT"
    case C.YAC_TYPE_DOUBLE:
        return "DOUBLE"
    case C.YAC_TYPE_NUMBER:
        return "NUMBER"
    case C.YAC_TYPE_DATE:
        return "DATE"
    case C.YAC_TYPE_TIMESTAMP:
        return "TIMESTAMP"
    case C.YAC_TYPE_CHAR:
        return "CHAR"
    case C.YAC_TYPE_NCHAR:
        return "NCHAR"
    case C.YAC_TYPE_VARCHAR:
        return "VARCHAR"
    case C.YAC_TYPE_NVARCHAR:
        return "NVARCHAR"
    case C.YAC_TYPE_CLOB:
        return "CLOB"
    case C.YAC_TYPE_BLOB:
        return "BLOB"
    default:
        return ""
    }
}

func (r *YasRows) ColumnTypeLength(index int) (length int64, ok bool) {
    if len(r.fetchRows) < index+1 {
        return 0, false
    }
    switch r.fetchRows[index].yacType {
    case C.YAC_TYPE_CHAR, C.YAC_TYPE_NCHAR, C.YAC_TYPE_VARCHAR, C.YAC_TYPE_NVARCHAR:
        return int64(r.fetchRows[index].Size), true
    case C.YAC_TYPE_BLOB, C.YAC_TYPE_CLOB:
        return math.MaxInt64, true
    default:
        return 0, false
    }
}

func (r *YasRows) getValues() (*[]driver.Value, error) {
    var err error
    unsafeRows := (unsafe.Pointer)(new(uint32))
    rows := (*C.YacUint32)(unsafeRows)
    if err = checkYasError(C.yacFetch(*r.stmt.Stmt, rows)); err != nil {
        return nil, err
    }
    if *rows == 0 {
        return nil, nil
    }
    columns := len(r.fetchRows)
    dest := make([]driver.Value, columns)
    for i := 0; i < columns; i++ {
        var value driver.Value
        row := r.fetchRows[i]
        if row == nil {
            return &dest, nil
        }
        switch row.yacType {
        case C.YAC_TYPE_BOOL:
            value = (*(*bool)(row.Data))
        case C.YAC_TYPE_TINYINT:
            value = (*(*int8)(row.Data))
        case C.YAC_TYPE_SMALLINT:
            value = (*(*int16)(row.Data))
        case C.YAC_TYPE_INTEGER:
            value = (*(*int32)(row.Data))
        case C.YAC_TYPE_BIGINT:
            value = (*(*int64)(row.Data))
        case C.YAC_TYPE_FLOAT:
            value = (*(*float32)(row.Data))
        case C.YAC_TYPE_DOUBLE:
            value = (*(*float64)(row.Data))
        case C.YAC_TYPE_DATE, C.YAC_TYPE_TIMESTAMP:
            value = time.Unix(0, (*(*int64)(row.Data))*1e3)
        case C.YAC_TYPE_CHAR, C.YAC_TYPE_NCHAR, C.YAC_TYPE_VARCHAR, C.YAC_TYPE_NVARCHAR:
            value = (C.GoString((*C.char)(row.Data)))
        case C.YAC_TYPE_NUMBER:
            value, err = strconv.ParseFloat(C.GoString((*C.char)(row.Data)), 64)
            if err != nil {
                return nil, err
            }
        case C.YAC_TYPE_CLOB, C.YAC_TYPE_BLOB:
            lobLocator := (**C.YacLobLocator)(row.Data)
            data, err := r.stmt.Conn.lobRead(*lobLocator)
            if err != nil {
                return nil, err
            }
            value = data
            if C.YacType(row.yacType) == C.YAC_TYPE_CLOB {
                value = string(data)
            }
        }
        dest[i] = value
    }
    return &dest, nil
}
