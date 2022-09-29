/*
Copyright  2022, YashanDB and/or its affiliates. All rights reserved.
YashanDB Driver for golang is licensed under the terms of the mulan PSL v2.0

License: 	http://license.coscl.org.cn/MulanPSL2
Home page: 	https://www.yashandb.com/
*/

package yasdb

import (
    "fmt"
    "reflect"
)

type YasBaseError struct {
    Code int
    Msg  string
}

func (e *YasBaseError) Error() string { return e.Msg }

type YasDBError struct {
    Code     int
    Msg      string
    SqlState string
    Line     int
    Column   int
}

func (e *YasDBError) Error() string {
    return fmt.Sprintf("%s:%d:%d [%d:%s]", e.SqlState, e.Line, e.Column, e.Code, e.Msg)
}

func ErrDsnNoStandard(dsn string) *YasBaseError {
    return &YasBaseError{Code: 1001, Msg: "dsn is nonstandard"}
}
func ErrDsnNoSet() *YasBaseError {
    return &YasBaseError{Code: 1002, Msg: "dsn is unset"}
}
func ErrNoConnect() *YasBaseError {
    return &YasBaseError{Code: 1003, Msg: "yasdb is not connected"}
}
func ErrStmtInit() *YasBaseError {
    return &YasBaseError{Code: 1004, Msg: "yasdb stmt init failed"}
}
func ErrStmtNoOpen() *YasBaseError {
    return &YasBaseError{Code: 1005, Msg: "yasdb stmt is not open"}
}
func ErrDbTypeUnsupport(dbType int) *YasBaseError {
    return &YasBaseError{Code: 1006, Msg: fmt.Sprintf("yasdb type %d is not support", dbType)}
}
func ErrDbFetchEOF() *YasBaseError {
    return &YasBaseError{Code: 1007, Msg: "fetch is over"}
}
func ErrUnknowType(v interface{}) *YasBaseError {
    return &YasBaseError{Code: 1008, Msg: fmt.Sprintf("Unknow: %s", reflect.TypeOf(v).String())}
}

func ErrOutputBindValue() *YasBaseError {
    return &YasBaseError{}
}
