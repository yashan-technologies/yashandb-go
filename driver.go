/*
Copyright  2022, YashanDB and/or its affiliates. All rights reserved.
YashanDB Driver for golang is licensed under the terms of the mulan PSL v2.0

License: 	http://license.coscl.org.cn/MulanPSL2
Home page: 	https://www.yashandb.com/
*/

package yasdb

/*
#cgo CFLAGS: -I./yacapi/include

#include "yacapi.h"
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
    "database/sql/driver"
    "unsafe"
)

type YasdbDriver struct{}

// Open returns a new connection to the database.
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
    var env *C.YapiEnv
    if err := checkYasError(C.yapiAllocEnv(&env)); err != nil {
        return nil, err
    }
    if dsn.DataPath != "" {
        dataPath := stringToYasChar(dsn.DataPath)
        defer C.free(unsafe.Pointer(dataPath))

        dpLen := intToYacInt32(len(dsn.DataPath))
        if err := checkYasError(C.yapiSetEnvAttr(env, C.YAPI_ATTR_DATA_PATH, unsafe.Pointer(dataPath), dpLen)); err != nil {
            return nil, err
        }
    }
    var conn *C.YapiConnect

    url := stringToYasChar(dsn.Url)
    defer C.free(unsafe.Pointer(url))
    user := stringToYasChar(dsn.User)
    defer C.free(unsafe.Pointer(user))
    password := stringToYasChar(dsn.Password)
    defer C.free(unsafe.Pointer(password))
    urlLen := intToYacInt16(len(dsn.Url))
    userLen := intToYacInt16(len(dsn.User))
    pwLen := intToYacInt16(len(dsn.Password))
    if err := checkYasError(C.yapiConnect(env, url, urlLen, user, userLen, password, pwLen, &conn)); err != nil {
        return nil, err
    }
    yasConn := &YasConn{
        Env:        env,
        Conn:       conn,
        autoCommit: dsn.IsAutoCommit,
    }
    if err := yasConn.setAutoCommit(dsn.IsAutoCommit); err != nil {
        return nil, err
    }
    return yasConn, nil
}

type YasTx struct {
    Conn *YasConn
}

// Commit transaction commit
func (tx *YasTx) Commit() error {
    return tx.Conn.yacCommit()
}

// Rollback transaction rollback
func (tx *YasTx) Rollback() error {
    return tx.Conn.yacRollback()
}
