package yasdb

import (
    "regexp"
    "strings"
)

const (
    dsnRegExpr = `^(.*?)/(.*?)@(.*?)(\?(.*?))?$`
    urlRegExpr = `^\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}:\d+$`
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
    dsnReg, _ := regexp.Compile(dsnRegExpr)
    urlReg, _ := regexp.Compile(urlRegExpr)
    if !dsnReg.MatchString(dsnStr) {
        return nil, ErrDsnNoStandard(dsnStr)
    }
    matchStrs := dsnReg.FindStringSubmatch(dsnStr)
    dsn := &DataSourceName{
        User:     matchStrs[1],
        Password: matchStrs[2],
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
