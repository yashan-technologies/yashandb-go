/*
Copyright  2022, YashanDB and/or its affiliates. All rights reserved.
YashanDB Driver for golang is licensed under the terms of the mulan PSL v2.0

License: 	http://license.coscl.org.cn/MulanPSL2
Home page: 	https://www.yashandb.com/
*/

package yasdb

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	dsnRegExpr       = `^(.*?)/(.*?)@(.*?)(\?(.*?))?$`
	ipv4UrlRegExpr   = `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d+$`
	ipv6UrlRegExpr   = `^\[[:|\d|A-Z|a-z|%]+\]:\d+$`
	mappedUrlRegExpr = `^\[[:|\d|A-Z|a-z|\.]+\]:\d+$`
	udsRegExpr       = `^(.*?)(\?(.*?))?$`

	_UkeyName   = `ukey_name`
	_UkeyPin    = `ukey_pin`
	_Autocommit = "autocommit"
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
	DataPath     string
	IsAutoCommit bool
	ukeyName     string
	ukeyPin      string
}

// ParseDSN parses a DataSourceName used to connect to YashanDB
//
// It expects to receive a string in the form:
//
// [username/[password]@]host[:port][?param1=value1&...&paramN=valueN]
// OR
// YASDB_DATA_PATH[?param1=value1&...&paramN=valueN]
// Supported parameters are:
//
// autocommit - When it is true, the transaction will be automatically committed every time an SQL statement is executed. Default is false
func ParseDSN(dsnStr string) (*DataSourceName, error) {
	if dsnStr == "" {
		return nil, ErrDsnNoSet()
	}
	if isDsn(dsnStr) {
		return parseDSN(dsnStr)
	}
	return parseUDS(dsnStr)
}

func parseDSN(dsnStr string) (*DataSourceName, error) {
	dsnStr = replaceSpecialChars(dsnStr)
	dsnReg, _ := regexp.Compile(dsnRegExpr)

	if !dsnReg.MatchString(dsnStr) {
		return nil, ErrDsnNoStandard(dsnStr)
	}
	matchStrs := dsnReg.FindStringSubmatch(dsnStr)
	dsn := &DataSourceName{
		User:     recoverySpecialChars(matchStrs[1]),
		Password: recoverySpecialChars(matchStrs[2]),
		Url:      matchStrs[3],
		DataPath: "",
	}

	if err := checkUrl(dsn.Url); err != nil {
		return nil, err
	}

	if err := parseParams(dsn, matchStrs[4]); err != nil {
		return nil, err
	}
	genUkeyUrl(dsn)

	return dsn, nil
}

func parseUDS(dsnStr string) (*DataSourceName, error) {
	udsReg, _ := regexp.Compile(udsRegExpr)
	if !udsReg.MatchString(dsnStr) {
		return nil, ErrDsnNoStandard(dsnStr)
	}
	matchStrs := udsReg.FindStringSubmatch(dsnStr)
	dsn := &DataSourceName{
		User:     "sys",
		Password: "",
		Url:      "",
		DataPath: matchStrs[1],
	}
	_, err := os.Stat(dsn.DataPath)
	if err != nil && !os.IsExist(err) {
		return nil, ErrDataPathNoExist(dsnStr)
	}
	if err := parseParams(dsn, matchStrs[2]); err != nil {
		return nil, err
	}
	return dsn, nil
}

func parseParams(dsn *DataSourceName, argStr string) error {
	if argStr == "" {
		return nil
	}
	paramStr := argStr
	if argStr[0] == '?' || argStr[0] == '&' {
		paramStr = argStr[1:]
	}
	connParams := strings.Split(paramStr, "&")
	for _, param := range connParams {
		strs := strings.Split(param, "=")
		if len(strs) < 2 {
			return ErrDsnNoStandard(argStr)
		}
		switch strings.ToLower(strs[0]) {
		case _Autocommit:
			value := strings.ToLower(strs[1])
			if value == "1" || value == "true" {
				dsn.IsAutoCommit = true
			}
		case _UkeyName:
			dsn.ukeyName = strs[1]
		case _UkeyPin:
			dsn.ukeyPin = strs[1]
		default:
			return fmt.Errorf("unknown param %s", strs[0])
		}
	}

	return nil
}

func genUkeyUrl(dsn *DataSourceName) {
	if dsn.ukeyName != "" && dsn.ukeyPin != "" {
		dsn.Url += fmt.Sprintf("?%s=%s&%s=%s", _UkeyName, dsn.ukeyName, _UkeyPin, dsn.ukeyPin)
	} else if dsn.ukeyName != "" {
		dsn.Url += fmt.Sprintf("?%s=%s", _UkeyName, dsn.ukeyName)
	} else if dsn.ukeyPin != "" {
		dsn.Url += fmt.Sprintf("?%s=%s", _UkeyPin, dsn.ukeyPin)
	}
}

func isDsn(dsnStr string) bool {
	dsnReg, _ := regexp.Compile(dsnRegExpr)
	return dsnReg.MatchString(dsnStr)
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

func checkUrl(url string) error {
	ipv4UrlReg, _ := regexp.Compile(ipv4UrlRegExpr)
	ipv6UrlReg, _ := regexp.Compile(ipv6UrlRegExpr)
	mappedUrlReg, _ := regexp.Compile(mappedUrlRegExpr)
	primaryStr := "primary:"
	strs := strings.Split(strings.ToLower(url), ",")
	for i, str := range strs {
		if i == 0 && strings.HasPrefix(str, primaryStr) {
			str = strings.TrimPrefix(str, primaryStr)
		}
		if !ipv4UrlReg.MatchString(str) && !ipv6UrlReg.MatchString(str) && !mappedUrlReg.MatchString(str) {
			return ErrDsnNoStandard(str)
		}
	}
	return nil
}
