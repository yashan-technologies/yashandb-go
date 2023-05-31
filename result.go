/*
Copyright  2022, YashanDB and/or its affiliates. All rights reserved.
YashanDB Driver for golang is licensed under the terms of the mulan PSL v2.0

License: 	http://license.coscl.org.cn/MulanPSL2
Home page: 	https://www.yashandb.com/
*/

package yasdb

import "errors"

var ErrInsertIdUnsupport = errors.New("last insert id is unsupport")

type YasResult struct {
	rowsAffected    int64
	rowsAffectedErr error
	lastInsertId    int64
	lastInsertIdErr error
}

// LastInsertId returns the database's auto-generated ID
// after, for example, an INSERT into a table with primary
// key.
func (result *YasResult) LastInsertId() (int64, error) {
	return result.lastInsertId, result.lastInsertIdErr
}

// RowsAffected returns the number of rows affected by the
// query.
func (result *YasResult) RowsAffected() (int64, error) {
	return result.rowsAffected, result.rowsAffectedErr
}
