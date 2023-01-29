/*
Copyright  2022, YashanDB and/or its affiliates. All rights reserved.
YashanDB Driver for golang is licensed under the terms of the mulan PSL v2.0

License: 	http://license.coscl.org.cn/MulanPSL2
Home page: 	https://www.yashandb.com/
*/

package yasdb

import (
    "regexp"
    "strings"
)

const (
    dsnRegExpr = `^(.*?)/(.*?)@(.*?)(\?(.*?))?$`
    urlRegExpr = `^\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}:\d+$`
)

var (
    specialChars = map[string]string{
        `\`: "{dsn_placeholder_1}",
        `/`: "{dsn_placeholder_2}",
        `@`: "{dsn_placeholder_3}",
    }
    escapeChar = `\`
)

type DataSourceName struct {
    User         string
    Password     string
    Url          string
    IsAutoCommit bool
}

// ParseDSN parses a DataSourceName used to connect to YashanDB
//
// It expects to receive a string in the form:
//
// [username/[password]@]host[:port][?param1=value1&...&paramN=valueN]
//
// Supported parameters are:
//
// autocommit - When it is true, the transaction will be automatically committed every time an SQL statement is executed. Default is false
func ParseDSN(dsnStr string) (*DataSourceName, error) {
    if dsnStr == "" {
        return nil, ErrDsnNoSet()
    }
    return parseDSN(dsnStr)
}

func parseDSN(dsnStr string) (*DataSourceName, error) {
    dsnStr = replaceSpecialChars(dsnStr)
    dsnReg, _ := regexp.Compile(dsnRegExpr)
    urlReg, _ := regexp.Compile(urlRegExpr)
    if !dsnReg.MatchString(dsnStr) {
        return nil, ErrDsnNoStandard(dsnStr)
    }
    matchStrs := dsnReg.FindStringSubmatch(dsnStr)
    dsn := &DataSourceName{
        User:     recoverySpecialChars(matchStrs[1]),
        Password: recoverySpecialChars(matchStrs[2]),
        Url:      matchStrs[3],
    }
    if !urlReg.MatchString(dsn.Url) {
        return nil, ErrDsnNoStandard(dsnStr)
    }
    if len(matchStrs[4]) > 1 {
        paramStr := strings.ToLower(matchStrs[4][1:])
        connParams := strings.Split(paramStr, "&")
        for _, param := range connParams {
            if param == "autocommit=1" || param == "autocommit=true" {
                dsn.IsAutoCommit = true
            }
        }
    }
    return dsn, nil
}

func replaceSpecialChars(dsnStr string) string {
    for k, v := range specialChars {
        dsnStr = strings.ReplaceAll(dsnStr, escapeChar+k, v)
    }
    return dsnStr
}

func recoverySpecialChars(str string) string {
    for k, v := range specialChars {
        str = strings.ReplaceAll(str, v, k)
    }
    return str
}
