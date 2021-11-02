package yasdb

import (
    "database/sql/driver"
    "log"
    "regexp"
)

type YasDriver struct{}

func init() {
    log.Println("driver is call ")
}

func (driver *YasDriver) Open(dsn string) (driver.Conn, error) {
    log.Println("exec open driver")
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
