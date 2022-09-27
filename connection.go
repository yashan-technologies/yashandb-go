package yasdb

/*
#cgo !noPkgConfig pkg-config: yacli
#include "yacli.go.h"
*/
import "C"
import (
    "context"
    "database/sql/driver"
    "unsafe"
)

type YasConn struct {
    Env        YacHandle
    Conn       YacHandle
    autoCommit bool
    closed     bool
}

func NewYasConn() *YasConn {
    return &YasConn{
        Env:  NewYacHandle(),
        Conn: NewYacHandle(),
    }
}

// Prepare statement for prepare exec
func (conn *YasConn) Prepare(query string) (driver.Stmt, error) {
    return conn.PrepareContext(context.Background(), query)
}

func (conn *YasConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
    if ctx.Err() != nil {
        return nil, ctx.Err()
    }

    stmt, err := NewYasStmt(conn, ctx)
    if err != nil {
        return nil, err
    }

    queryP := C.CString(query)
    defer C.free(unsafe.Pointer(queryP))
    sqlLength := C.YacInt32(0)
    if err := checkYasError(
        C.yacPrepare(
            *stmt.Stmt,
            queryP,
            sqlLength,
        )); err != nil {
        return nil, err
    }

    var sqltype C.YacUint32
    sqlSize := C.YacInt32(unsafe.Sizeof(&sqltype))
    if err := checkYasError(
        C.yacGetStmtAttr(
            *stmt.Stmt,
            C.YAC_ATTR_SQLTYPE,
            unsafe.Pointer(&sqltype),
            sqlSize,
            &sqlLength,
        )); err != nil {
        return nil, err
    }
    stmt.SqlType = (uint32)(sqltype)

    return stmt, nil
}

// Begin begin
func (conn *YasConn) Begin() (driver.Tx, error) {
    return conn.BeginTx(context.Background(), driver.TxOptions{})
}

func (conn *YasConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
    if ctx.Err() != nil {
        return nil, ctx.Err()
    }
    return &YasTx{Conn: conn}, nil
}

// Close close db YasConn
func (conn *YasConn) Close() error {
    if conn.closed {
        return nil
    }
    conn.closed = true

    if err := yasdbFreeHandle(conn.Conn, YAC_HANDLE_DBC); err != nil {
        return err
    }
    return yasdbFreeHandle(conn.Env, YAC_HANDLE_ENV)
}

func (conn *YasConn) Ping(ctx context.Context) error {
    if ctx.Err() != nil {
        return ctx.Err()
    }
    if conn.Conn == nil {
        return ErrNoConnect()
    }
    return nil
}

func (conn *YasConn) IsAutoCommit() bool {
    return conn.autoCommit
}

func (conn *YasConn) setAutoCommit(auto bool) error {
    var a C.YacInt32 = 0
    if auto {
        a = 1
    }
    size := C.YacInt32(unsafe.Sizeof(a))
    if err := conn.yacSetConnAttr(C.YAC_ATTR_AUTOCOMMIT, unsafe.Pointer(&a), size); err != nil {
        return err
    }
    conn.autoCommit = auto
    return nil
}

func (conn *YasConn) yacSetConnAttr(attr C.YacConnAttr, value unsafe.Pointer, bufLength C.YacInt32) error {
    if err := checkYasError(
        C.yacSetConnAttr(
            *conn.Conn,
            attr,
            value,
            bufLength,
        )); err != nil {
        return err
    }
    return nil
}

func (conn *YasConn) yacCommit() error {
    return checkYasError(C.yacCommit(*conn.Conn))
}

func (conn *YasConn) yacRollback() error {
    return checkYasError(C.yacRollback(*conn.Conn))
}

func (conn *YasConn) lobRead(lobLocator *C.YacLobLocator) ([]byte, error) {
    lobLen, err := conn.yacLobGetLength(lobLocator)
    if err != nil {
        return nil, err
    }
    data, err := conn.yacLobRead(lobLocator, lobLen)
    if err != nil {
        return nil, err
    }
    return data, nil
}

func (conn *YasConn) yacLobRead(lobLocator *C.YacLobLocator, lobLen uint64) ([]byte, error) {
    if lobLen == 0 {
        return []byte{}, nil
    }
    data := make([]byte, 0, lobLen)
    bytes := C.YacUint64(_LobBufLen)
    for {
        readBuffer := byteBufferPool.Get().([]byte)
        buf := (*C.YacUint8)((unsafe.Pointer)(&readBuffer[0]))
        if err := checkYasError(
            C.yacLobRead(
                *conn.Conn,
                lobLocator,
                &bytes,
                buf,
                _LobBufLen,
            )); err != nil {
            return nil, nil
        }
        data = append(data, readBuffer[:uint64(bytes)]...)
        if uint64(bytes) < _LobBufLen {
            break
        }
    }
    return data, nil
}

func (conn *YasConn) yacLobGetLength(lobLocator *C.YacLobLocator) (uint64, error) {
    var lobLen C.YacUint64
    if err := checkYasError(C.yacLobGetLength(*conn.Conn, lobLocator, &lobLen)); err != nil {
        return 0, err
    }
    return uint64(lobLen), nil
}

func (conn *YasConn) lobWrite(yacType C.YacType, data []byte) (*unsafe.Pointer, error) {
    desc, err := conn.yacLobDescAlloc(yacType)
    if err != nil {
        return nil, err
    }
    lobLocator := (**C.YacLobLocator)(unsafe.Pointer(desc))
    if err := conn.yacLobCreateTemporary(*lobLocator); err != nil {
        return nil, err
    }
    if err := conn.yacLobWrite(*lobLocator, data); err != nil {
        return nil, err
    }
    return desc, nil
}

func (conn *YasConn) yacLobDescAlloc(yacType C.YacType) (*unsafe.Pointer, error) {
    var desc = new(unsafe.Pointer)
    if err := checkYasError(C.yacLobDescAlloc(*conn.Conn, yacType, desc)); err != nil {
        return nil, err
    }
    return desc, nil
}

func (conn *YasConn) yacLobCreateTemporary(lobLocator *C.YacLobLocator) error {
    if err := checkYasError(C.yacLobCreateTemporary(*conn.Conn, lobLocator)); err != nil {
        return err
    }
    return nil
}

func (conn *YasConn) yacLobWrite(lobLocator *C.YacLobLocator, data []byte) error {
    if len(data) == 0 || data == nil {
        return nil
    }
    bufLen := uint64(_LobBufLen)
    dataLen := uint64(len(data))
    writeBuffer := byteBufferPool.Get().([]byte)
    if _LobBufLen > dataLen {
        bufLen = dataLen
        copy(writeBuffer, data)
    } else {
        copy(writeBuffer, data[0:_LobBufLen])
    }
    buf := (*C.YacUint8)((unsafe.Pointer)(&writeBuffer[0]))
    count := uint64(0)
    for {
        if err := checkYasError(
            C.yacLobWrite(
                *conn.Conn,
                lobLocator,
                nil,
                buf,
                C.YacUint64(bufLen),
            )); err != nil {
            return nil
        }
        count += bufLen
        if count >= dataLen {
            break
        }
        if count+bufLen < dataLen {
            copy(writeBuffer, data[count:count+bufLen])
        } else {
            copy(writeBuffer, data[count:])
            bufLen = dataLen - count
        }
    }
    return nil
}

func (conn *YasConn) lobFree(yacType C.YacType, lobLocator *C.YacLobLocator) {
    if yacType != C.YAC_TYPE_BLOB && yacType != C.YAC_TYPE_CLOB {
        return
    }
    C.yacLobFreeTemporary(*conn.Conn, lobLocator)
    C.yacLobDescFree(unsafe.Pointer(lobLocator), yacType)
}

func (conn *YasConn) handleYacCancel(ctx context.Context, done <-chan struct{}) {
    select {
    case <-done:
    case <-ctx.Done():
        select {
        case <-done:
        default:
            _ = conn.yacCancel()
        }
    }
}

func (conn *YasConn) yacCancel() error {
    return checkYasError(C.yacCancel(*conn.Conn))
}
