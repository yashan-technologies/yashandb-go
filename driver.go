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

var _ClientDriverName = "YashanDB GO Driver"

func SetClientDriverName(name string) {
	if name == "" {
		return
	}
	_ClientDriverName = name
}

type YasdbDriver struct{}

// Open returns a new connection to the database.
func (yasDriver *YasdbDriver) Open(dsnStr string) (driver.Conn, error) {
	return GenYasconn(dsnStr)
}

func GenYasconn(dsnStr string) (*YasConn, error) {
	dsn, err := ParseDSN(dsnStr)
	if err != nil {
		return nil, err
	}

	var env *C.YapiEnv
	if err := checkYasError(C.yapiAllocEnv(&env)); err != nil {
		return nil, err
	}
	if dsn.DataPath != "" {
		dataPath := stringToYasChar(dsn.DataPath)
		defer C.free(unsafe.Pointer(dataPath))

		dpLen := intToYacInt32(len(dsn.DataPath))
		if err := checkYasError(C.yapiSetEnvAttr(env, C.YAPI_ATTR_DATA_PATH, unsafe.Pointer(dataPath), dpLen)); err != nil {
			_ = releaseEnv(env)
			return nil, err
		}
	}

	charset := C.YAPI_CHARSET_UTF8
	if err := checkYasError(C.yapiSetEnvAttr(env, C.YAPI_ATTR_CHARSET_CODE, unsafe.Pointer(&charset), 4)); err != nil {
		_ = releaseEnv(env)
		return nil, err
	}

	driverName := stringToYasChar(_ClientDriverName)
	driverNameLen := intToYacInt32(len(_ClientDriverName))
	defer C.free(unsafe.Pointer(driverName))
	C.yapiSetEnvAttr(env, C.YAPI_ATTR_CLIENT_DRIVER, unsafe.Pointer(driverName), driverNameLen)

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
		_ = releaseEnv(env)
		return nil, err
	}
	yasConn := &YasConn{
		Env:        env,
		Conn:       conn,
		autoCommit: dsn.IsAutoCommit,
	}
	if err := yasConn.setAutoCommit(dsn.IsAutoCommit); err != nil {
		_ = yasConn.Close()
		return nil, err
	}

	if err := yasConn.getConnAttr(); err != nil {
		_ = yasConn.Close()
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
