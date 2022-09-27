package yasdb

/*
#cgo !noPkgConfig pkg-config: yacli
#include "yacli.go.h"
*/
import "C"
import (
    "database/sql/driver"
    "unsafe"
)

type YasdbDriver struct{}

func (yasDriver *YasdbDriver) Open(dsnStr string) (driver.Conn, error) {
    dsn, err := ParseDSN(dsnStr)
    if err != nil {
        return nil, err
    }
    conn, err := yasDriver.getYasConn(dsn)
    if err != nil {
        return nil, err
    }
    return conn, nil
}

func (yasDriver *YasdbDriver) getYasConn(dsn *DataSourceName) (driver.Conn, error) {
    var isAllocHandleEnv, isAllocHandleDbc bool
    var err error
    conn := NewYasConn()

    defer func(errHandle *error) {
        if *errHandle == nil {
            return
        }
        if isAllocHandleEnv {
            yasdbFreeHandle(conn.Env, C.YAC_HANDLE_ENV)
        }
        if isAllocHandleDbc {
            yasdbFreeHandle(conn.Conn, C.YAC_HANDLE_DBC)
        }

    }(&err)

    if err = checkYasError(C.yacAllocHandle(C.YAC_HANDLE_ENV, nil, conn.Env)); err != nil {
        return nil, err
    }
    isAllocHandleEnv = true

    if err = checkYasError(C.yacAllocHandle(C.YAC_HANDLE_DBC, *conn.Env, conn.Conn)); err != nil {
        return nil, err
    }
    isAllocHandleDbc = true

    url := stringToYasChar(dsn.Url)
    defer C.free(unsafe.Pointer(url))
    user := stringToYasChar(dsn.User)
    defer C.free(unsafe.Pointer(user))
    password := stringToYasChar(dsn.Password)
    defer C.free(unsafe.Pointer(password))
    urlLen := intToYacInt16(len(dsn.Url))
    userLen := intToYacInt16(len(dsn.User))
    pwLen := intToYacInt16(len(dsn.Password))
    if err = checkYasError(C.yacConnect(*conn.Conn, url, urlLen, user, userLen, password, pwLen)); err != nil {
        return nil, err
    }

    if err := conn.setAutoCommit(dsn.IsAutoCommit); err != nil {
        return nil, err
    }
    return conn, nil
}

type YasTx struct {
    Conn *YasConn
}

func (tx *YasTx) Commit() error {
    return tx.Conn.yacCommit()
}
func (tx *YasTx) Rollback() error {
    return tx.Conn.yacRollback()
}
