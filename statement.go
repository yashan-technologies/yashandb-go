package yasdb

// #include "yacli.go.h"
import "C"
import (
    "context"
    "database/sql"
    "database/sql/driver"
    "fmt"
    "math"
    "sync"
    "time"
    "unsafe"
)

type YasStmt struct {
    Conn    *YasConn
    Stmt    YacHandle
    closed  bool
    SqlType uint32
    ctx     context.Context
    binds   []bindStruct
    sync.Mutex
}

func NewYasStmt(conn *YasConn, ctx context.Context) (*YasStmt, error) {
    if conn == nil || conn.Conn == nil {
        return nil, ErrNoConnect()
    }
    stmt := &YasStmt{
        Conn: conn,
        Stmt: NewYacHandle(),
        ctx:  ctx,
    }
    if err := checkYasError(C.yacAllocHandle(C.YAC_HANDLE_STMT, *conn.Conn, stmt.Stmt)); err != nil {
        return nil, err
    }
    return stmt, nil
}

func (stmt *YasStmt) Query(args []driver.Value) (driver.Rows, error) {
    nargs := make([]driver.NamedValue, len(args))
    for i, arg := range args {
        nargs[i].Ordinal = i + 1
        nargs[i].Value = arg
    }
    return stmt.QueryContext(context.Background(), nargs)
}

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

func (stmt *YasStmt) Exec(args []driver.Value) (driver.Result, error) {
    nargs := make([]driver.NamedValue, len(args))
    for i, arg := range args {
        nargs[i].Ordinal = i + 1
        nargs[i].Value = arg
    }

    return stmt.ExecContext(context.Background(), nargs)
}

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

func (stmt *YasStmt) NumInput() int {
    return -1
}

func (stmt *YasStmt) Close() error {
    if stmt.closed {
        return nil
    }
    stmt.closed = true
    return yasdbFreeHandle(stmt.Stmt, YAC_HANDLE_STMT)
}

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
    if err := stmt.yacExecute(); err != nil {
        return nil, err
    }
    close(done)

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
    if err := stmt.yacExecute(); err != nil {
        return nil, err
    }
    close(done)

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
    return checkYasError(C.yacExecute(*stmt.Stmt))
}

func (stmt *YasStmt) getFetchRows() ([]*yasRow, error) {
    columns := C.YacInt16(0)
    if err := checkYasError(C.yacNumResultCols(*stmt.Stmt, &columns)); err != nil {
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
    item := C.YacColumnDesc{}
    if err := checkYasError(C.yacDescribeCol2(*stmt.Stmt, C.YacUint16(pos), &item)); err != nil {
        return nil, err
    }
    yacType := C.YacType(item._type)
    size, indicator := uint32(item.size), C.YacInt32(0)
    bufLen := int32(size)
    row := NewYasRow(stmt, size, yacType)
    row.name = C.GoString(item.name)
    switch yacType {
    case C.YAC_TYPE_CHAR, C.YAC_TYPE_NCHAR, C.YAC_TYPE_VARCHAR, C.YAC_TYPE_NVARCHAR:
        bufLen = int32(sizeToAlign4(size)) + 1
    case C.YAC_TYPE_NUMBER, C.YAC_TYPE_YM_INTERVAL, C.YAC_TYPE_DS_INTERVAL: // number to string
        yacType = C.YAC_TYPE_VARCHAR
        bufLen = int32(sizeToAlign4(uint32(item.precision) + 8))
    case C.YAC_TYPE_CLOB, C.YAC_TYPE_BLOB:
        var desc = new(unsafe.Pointer)
        if err := checkYasError(C.yacLobDescAlloc(*stmt.Conn.Conn, yacType, desc)); err != nil {
            return nil, err
        }
        row.Data = unsafe.Pointer(desc)
        bufLen = -1
    }
    if err := checkYasError(
        C.yacBindColumn(
            *stmt.Stmt,
            C.YacUint16(pos),
            yacType,
            C.YacPointer(row.Data),
            C.YacInt32(bufLen),
            &indicator,
        ),
    ); err != nil {
        return nil, err
    }
    row.Indicator = int32(indicator)
    return row, nil
}

func (stmt *YasStmt) getRowsAffected() (int64, error) {
    var rowsCount C.YacUint32
    size := C.YacInt32(unsafe.Sizeof(&rowsCount))
    s_length := C.YacInt32(0)
    err := checkYasError(
        C.yacGetStmtAttr(
            *stmt.Stmt,
            C.YAC_ATTR_ROWS_AFFECTED,
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
    stmt.binds = make([]bindStruct, 0, len(args))
    var err error
    for index, narg := range args {
        arg := narg.Value
        sqlOut, isOut := arg.(sql.Out)
        bind := bindStruct{}

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

func (stmt *YasStmt) yacBindParameter(b bindStruct, pos C.YacUint16) error {
    if err := checkYasError(
        C.yacBindParameter(
            *stmt.Stmt,
            pos,
            b.direction,
            b.yacType,
            b.value,
            b.bindSize,
            C.YacInt32(0),
            b.indicator,
        ),
    ); err != nil {
        return err
    }
    return nil
}

func (stmt *YasStmt) yacBindParameterByName(b bindStruct, name string) error {
    if err := checkYasError(
        C.yacBindParameterByName(
            *stmt.Stmt,
            stringToYasChar(name),
            b.direction,
            b.yacType,
            b.value,
            b.bindSize,
            C.YacInt32(0),
            nil,
        ),
    ); err != nil {
        return err
    }
    return nil
}

func (stmt *YasStmt) getInputBindValue(arg driver.Value) (bindStruct, error) {
    bind := bindStruct{}
    var (
        yacType   C.YacType
        size      C.YacUint32
        value     C.YacPointer
        indicator *C.YacInt32
        bufLength C.YacInt32
    )

    size = C.YacUint32(unsafe.Sizeof(&arg)) + 1
    bufLength = C.YacInt32(size - 1)
    indicator = new(C.YacInt32)
    *indicator = C.YacInt32(size - 1)

    switch v := arg.(type) {
    case int64:
        yacType = C.YAC_TYPE_BIGINT
        value = C.YacPointer(unsafe.Pointer(&v))
    case float64:
        yacType = C.YAC_TYPE_DOUBLE
        value = C.YacPointer(unsafe.Pointer(&v))
    case bool:
        yacType = C.YAC_TYPE_BOOL
        value = C.YacPointer(unsafe.Pointer(&v))
    case string:
        yacType = C.YAC_TYPE_VARCHAR
        size = intToYacUint32(len(v)) + 1
        bufLength = C.YacInt32(size - 1)
        indicator = nil
        value = C.YacPointer(unsafe.Pointer(stringToYasChar(v)))
    case []byte:
        desc, err := stmt.Conn.lobWrite(C.YAC_TYPE_BLOB, v)
        if err != nil {
            return bind, err
        }
        yacType = C.YAC_TYPE_BLOB
        size = C.YacUint32(math.MaxUint32)
        bufLength = -1
        indicator = nil
        value = C.YacPointer(desc)
    case time.Time:
        yacType = C.YAC_TYPE_TIMESTAMP
        t := v.UnixNano() / 1e3
        value = C.YacPointer(unsafe.Pointer(&t))
    case nil:
        yacType = C.YAC_TYPE_VARCHAR
        value = C.YacPointer(&v)
    default:
        return bind, ErrUnknowType(arg)
    }

    bind.yacType = yacType
    bind.value = value
    bind.bindSize = size
    bind.bufLength = bufLength
    bind.indicator = indicator
    bind.direction = C.YAC_PARAM_INPUT
    return bind, nil
}

func (stmt *YasStmt) getOutputBindValue(sqlOut sql.Out) (bindStruct, error) {
    if obi, ok := sqlOut.Dest.(*outputBindInfo); ok {
        return stmt.getOutputBindValueByInfo(obi)
    } else {
        return stmt.getOutputBindValueByDest(sqlOut.Dest)
    }
}

func (stmt *YasStmt) getOutputBindValueByDest(dest interface{}) (bindStruct, error) {
    bind := bindStruct{}
    var (
        yacType   C.YacType
        bindSize  C.YacUint32
        value     C.YacPointer
        indicator *C.YacInt32
        bufLength C.YacInt32
        arg       driver.Value
        err       error
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

    bindSize = C.YacUint32(unsafe.Sizeof(&arg)) + 1
    bufLength = C.YacInt32(bindSize)
    indicator = new(C.YacInt32)
    *indicator = C.YacInt32(bindSize - 1)

    switch v := arg.(type) {
    case int64:
        yacType = C.YAC_TYPE_INTEGER
        value = C.YacPointer(unsafe.Pointer(&v))
    case float64:
        yacType = C.YAC_TYPE_DOUBLE
        value = C.YacPointer(unsafe.Pointer(&v))
    case bool:
        yacType = C.YAC_TYPE_BOOL
        value = C.YacPointer(unsafe.Pointer(&v))
    case string:
        yacType = C.YAC_TYPE_VARCHAR
        bindSize = _OutputBindSize
        bufLength = C.YacInt32(bindSize - 1)
        value = C.YacPointer(unsafe.Pointer(stringToYasChar(v)))
    case []byte:
        desc, err := stmt.Conn.lobWrite(C.YAC_TYPE_BLOB, v)
        if err != nil {
            return bind, err
        }
        yacType = C.YAC_TYPE_BLOB
        bindSize = C.YacUint32(math.MaxUint32)
        bufLength = -1
        indicator = nil
        value = C.YacPointer(desc)
    case time.Time:
        yacType = C.YAC_TYPE_TIMESTAMP
        t := int64(0)
        value = C.YacPointer(unsafe.Pointer(&t))
    case nil:
        yacType = C.YAC_TYPE_VARCHAR
        value = C.YacPointer(&v)
    default:
        return bind, ErrUnknowType(v)
    }

    bind.yacType = yacType
    bind.value = value
    bind.bindSize = bindSize
    bind.bufLength = bufLength
    bind.indicator = indicator
    bind.direction = C.YAC_PARAM_OUTPUT
    return bind, nil
}

func (stmt *YasStmt) getOutputBindValueByInfo(obi *outputBindInfo) (bindStruct, error) {
    bind := bindStruct{}
    var (
        yacType   C.YacType = obi.yacType
        bindSize  C.YacUint32
        value     C.YacPointer
        indicator *C.YacInt32
        bufLength C.YacInt32
    )

    if obi.bindSize == 0 {
        bindSize = _OutputBindSize
    } else {
        bindSize = obi.bindSize
    }
    bufLength = C.YacInt32(bindSize)
    indicator = new(C.YacInt32)
    *indicator = C.YacInt32(bindSize - 1)

    switch yacType {
    case C.YAC_TYPE_BLOB:
        v, err := obi.getBlobBindDest()
        if err != nil {
            return bind, err
        }
        desc, err := stmt.Conn.lobWrite(C.YAC_TYPE_BLOB, *v)
        if err != nil {
            return bind, err
        }
        bindSize = C.YacUint32(math.MaxUint32)
        bufLength = -1
        indicator = nil
        value = C.YacPointer(desc)
    case C.YAC_TYPE_CLOB:
        v, err := obi.getClobBindDest()
        if err != nil {
            return bind, err
        }
        desc, err := stmt.Conn.lobWrite(C.YAC_TYPE_CLOB, []byte(*v))
        if err != nil {
            return bind, err
        }
        bindSize = C.YacUint32(math.MaxUint32)
        bufLength = -1
        indicator = nil
        value = C.YacPointer(desc)
    case C.YAC_TYPE_CHAR:
        v, err := obi.getCharBindDest()
        if err != nil {
            return bind, err
        }
        bufLength = C.YacInt32(bindSize - 1)
        value = C.YacPointer(unsafe.Pointer(stringToYasChar(*v)))
    case C.YAC_TYPE_VARCHAR:
        v, err := obi.getCharBindDest()
        if err != nil {
            return bind, err
        }
        bufLength = C.YacInt32(bindSize - 1)
        value = C.YacPointer(unsafe.Pointer(stringToYasChar(*v)))
    default:
        return bind, ErrUnknowType(yacType)
    }

    bind.yacType = yacType
    bind.value = value
    bind.bindSize = bindSize
    bind.bufLength = bufLength
    bind.indicator = indicator
    bind.direction = C.YAC_PARAM_OUTPUT
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
            lobLocator := (**C.YacLobLocator)(bind.value)
            *dest, err = stmt.Conn.lobRead(*lobLocator)
            if err != nil {
                return err
            }
        case *outputBindInfo:
            switch dest.yacType {
            case C.YAC_TYPE_BLOB:
                bindDest, _ := dest.getBlobBindDest()
                lobLocator := (**C.YacLobLocator)(bind.value)
                *bindDest, err = stmt.Conn.lobRead(*lobLocator)
                if err != nil {
                    return err
                }
            case C.YAC_TYPE_CLOB:
                bindDest, _ := dest.getClobBindDest()
                lobLocator := (**C.YacLobLocator)(bind.value)
                byteDest, err := stmt.Conn.lobRead(*lobLocator)
                *bindDest = string(byteDest)
                if err != nil {
                    return err
                }
            case C.YAC_TYPE_VARCHAR:
                bindDest, _ := dest.getVarcharBindDest()
                *bindDest = C.GoString((*C.char)(bind.value))
            case C.YAC_TYPE_CHAR:
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
        if bind.value != nil {
            switch bind.yacType {
            case C.YAC_TYPE_BLOB, C.YAC_TYPE_CLOB:
                lobLocator := (**C.YacLobLocator)(unsafe.Pointer(bind.value))
                stmt.Conn.lobFree(bind.yacType, *lobLocator)
            case C.YAC_TYPE_VARCHAR:
                C.free(unsafe.Pointer(bind.value))
            }
            bind.value = nil
        }
    }
    stmt.binds = []bindStruct{}
}

type outputBindInfo struct {
    yacType  C.YacType
    dest     interface{}
    bindSize C.YacUint32
}
type outputBindOpt func(*outputBindInfo)

func WithTypeClob() outputBindOpt {
    return func(obi *outputBindInfo) { obi.yacType = C.YAC_TYPE_CLOB }
}

func WithTypeBlob() outputBindOpt {
    return func(obi *outputBindInfo) { obi.yacType = C.YAC_TYPE_BLOB }
}

func WithTypeVarchar() outputBindOpt {
    return func(obi *outputBindInfo) { obi.yacType = C.YAC_TYPE_VARCHAR }
}

func WitchTypeChar() outputBindOpt {
    return func(obi *outputBindInfo) { obi.yacType = C.YAC_TYPE_CHAR }
}

func WithBindSize(bindSize uint32) outputBindOpt {
    return func(obi *outputBindInfo) { obi.bindSize = C.YacUint32(bindSize) }
}

//
//
//
func NewOutputBindValue(dest interface{}, opts ...outputBindOpt) (*outputBindInfo, error) {
    out := &outputBindInfo{
        dest:     dest,
        bindSize: C.YacUint32(0),
        yacType:  C.YacType(0),
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
    case C.YAC_TYPE_BLOB:
        _, err = obi.getBlobBindDest()
    case C.YAC_TYPE_CLOB:
        _, err = obi.getClobBindDest()
    case C.YAC_TYPE_VARCHAR:
        _, err = obi.getVarcharBindDest()
    case C.YAC_TYPE_CHAR:
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
