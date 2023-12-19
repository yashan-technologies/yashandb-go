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
	"context"
	"database/sql/driver"
	"fmt"
	"unsafe"
)

type YasConn struct {
	Env        *C.YapiEnv
	Conn       *C.YapiConnect
	autoCommit bool
	closed     bool
}

func (conn *YasConn) Prepare(query string) (driver.Stmt, error) {
	return conn.PrepareContext(context.Background(), query)
}

func (conn *YasConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var stmt *C.YapiStmt
	queryP := C.CString(tryRmSqlSemicolon(query))
	defer C.free(unsafe.Pointer(queryP))
	sqlLength := C.int32_t(len(query))
	if err := checkYasError(
		C.yapiPrepare(
			conn.Conn,
			queryP,
			sqlLength,
			&stmt,
		)); err != nil {
		return nil, err
	}

	var sqltype C.uint32_t
	sqlSize := C.int32_t(unsafe.Sizeof(sqltype))
	if err := checkYasError(
		C.yapiGetStmtAttr(
			stmt,
			C.YAPI_ATTR_SQLTYPE,
			unsafe.Pointer(&sqltype),
			sqlSize,
			&sqlLength,
		)); err != nil {
		return nil, err
	}

	yasStmt := &YasStmt{
		Conn:    conn,
		Stmt:    stmt,
		SqlType: (uint32)(sqltype),
	}

	return yasStmt, nil
}

func (conn *YasConn) Begin() (driver.Tx, error) {
	return conn.BeginTx(context.Background(), driver.TxOptions{})
}

func (conn *YasConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return &YasTx{Conn: conn}, nil
}

func (conn *YasConn) Close() error {
	if conn.closed {
		return nil
	}
	conn.closed = true

	connErr := releaseConn(conn.Conn)
	envErr := releaseEnv(conn.Env)
	if envErr != nil && connErr != nil {
		return fmt.Errorf("release env err: %s ; release conn err: %s", envErr, connErr)
	} else if envErr != nil {
		return envErr
	} else if connErr != nil {
		return connErr
	}
	return nil
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

func (conn *YasConn) setAutoCommit(auto bool) error {
	var a C.int32_t = 0
	if auto {
		a = 1
	}
	size := C.int32_t(unsafe.Sizeof(a))
	if err := conn.yapiSetConnAttr(C.YAPI_ATTR_AUTOCOMMIT, unsafe.Pointer(&a), size); err != nil {
		return err
	}
	conn.autoCommit = auto
	return nil
}

func (conn *YasConn) yapiSetConnAttr(attr C.YapiConnAttr, value unsafe.Pointer, bufLength C.int32_t) error {
	if err := checkYasError(
		C.yapiSetConnAttr(
			conn.Conn,
			attr,
			value,
			bufLength,
		)); err != nil {
		return err
	}
	return nil
}

func (conn *YasConn) yacCommit() error {
	return checkYasError(C.yapiCommit(conn.Conn))
}

func (conn *YasConn) yacRollback() error {
	return checkYasError(C.yapiRollback(conn.Conn))
}

func (conn *YasConn) lobRead(lobLocator *C.YapiLobLocator) ([]byte, error) {
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

func (conn *YasConn) yacLobRead(lobLocator *C.YapiLobLocator, lobLen uint64) ([]byte, error) {
	if lobLen == 0 {
		return []byte{}, nil
	}
	data := make([]byte, 0, lobLen)
	bytes := C.uint64_t(_LobBufLen)
	for {
		readBuffer := byteBufferPool.Get().([]byte)
		buf := (*C.uint8_t)((unsafe.Pointer)(&readBuffer[0]))
		if err := checkYasError(
			C.yapiLobRead(
				conn.Conn,
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

func (conn *YasConn) yacLobGetLength(lobLocator *C.YapiLobLocator) (uint64, error) {
	var lobLen C.uint64_t
	if err := checkYasError(C.yapiLobGetLength(conn.Conn, lobLocator, &lobLen)); err != nil {
		return 0, err
	}
	return uint64(lobLen), nil
}

func (conn *YasConn) lobWrite(yacType C.YapiType, data []byte) (*unsafe.Pointer, error) {
	desc, err := conn.yacLobDescAlloc(yacType)
	if err != nil {
		return nil, err
	}
	lobLocator := (**C.YapiLobLocator)(unsafe.Pointer(desc))
	if err := conn.yacLobCreateTemporary(*lobLocator); err != nil {
		return nil, err
	}
	if err := conn.yacLobWrite(*lobLocator, data); err != nil {
		return nil, err
	}
	return desc, nil
}

func (conn *YasConn) yacLobDescAlloc(yacType C.YapiType) (*unsafe.Pointer, error) {
	desc := new(unsafe.Pointer)
	if err := checkYasError(C.yapiLobDescAlloc(conn.Conn, yacType, desc)); err != nil {
		return nil, err
	}
	return desc, nil
}

func (conn *YasConn) yacLobCreateTemporary(lobLocator *C.YapiLobLocator) error {
	if err := checkYasError(C.yapiLobCreateTemporary(conn.Conn, lobLocator)); err != nil {
		return err
	}
	return nil
}

func (conn *YasConn) yacLobWrite(lobLocator *C.YapiLobLocator, data []byte) error {
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
	buf := (*C.uint8_t)((unsafe.Pointer)(&writeBuffer[0]))
	count := uint64(0)
	for {
		if err := checkYasError(
			C.yapiLobWrite(
				conn.Conn,
				lobLocator,
				nil,
				buf,
				C.uint64_t(bufLen),
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

func (conn *YasConn) lobFree(yacType C.YapiType, lobLocator *C.YapiLobLocator) {
	if yacType != C.YAPI_TYPE_BLOB && yacType != C.YAPI_TYPE_CLOB {
		return
	}
	C.yapiLobFreeTemporary(conn.Conn, lobLocator)
	C.yapiLobDescFree(unsafe.Pointer(lobLocator), yacType)
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
	return checkYasError(C.yapiCancel(conn.Conn))
}

func (conn *YasConn) ResetSession(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	stmt, err := conn.PrepareContext(ctx, "select 1 from dual")
	if err != nil {
		return conn.handleRestSessionErr(err)
	}
	defer stmt.Close()
	return nil
}

func (conn *YasConn) handleRestSessionErr(err error) error {
	if err == nil {
		return nil
	}
	if isDisconnetionErr(err) {
		return driver.ErrBadConn
	}
	return nil
}
