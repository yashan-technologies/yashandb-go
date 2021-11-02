package yasdb

import "errors"

var ErrInsertIdUnsupport = errors.New("last insert id is unsupport")

type YasResult struct {
    stmt *YasStmt
}

func (y *YasResult) LastInsertId() (int64, error) {
    return 0, ErrInsertIdUnsupport
}
func (y *YasResult) RowsAffected() (int64, error) {
    return yasdbRowAffected(y.stmt)
}
