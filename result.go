package yasdb

import "errors"

var ErrInsertIdUnsupport = errors.New("last insert id is unsupport")

type YasResult struct {
    rowsAffected    int64
    rowsAffectedErr error
    lastInsertId    int64
    lastInsertIdErr error
}

func (result *YasResult) LastInsertId() (int64, error) {
    return result.lastInsertId, result.lastInsertIdErr
}
func (result *YasResult) RowsAffected() (int64, error) {
    return result.rowsAffected, result.rowsAffectedErr
}
