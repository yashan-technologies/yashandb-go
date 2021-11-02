package yasdb

import (
    "database/sql/driver"
    "regexp"
)

type YasDriver struct{}

func (driver *YasDriver) Open(dsn string) (driver.Conn, error) {
    conn := NewConnection()
    re := regexp.MustCompile(`^(.*?)/(.*?)@(.*)$`)
    if !re.MatchString(dsn) {
        return nil, ErrDsnNoStandard(dsn)
    }
    conn.Dsn = dsn
    if err := yasdbConnect(conn); err != nil {
        return nil, err
    }
    return conn, nil
}

type YasTx struct {
    Conn *Connection
}

func (tx *YasTx) Commit() error {
    return yasdbCommit(tx.Conn)
}
func (tx *YasTx) Rollback() error {
    return yasdbRollback(tx.Conn)
}
