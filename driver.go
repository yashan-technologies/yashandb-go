package yasdb

import (
    "database/sql/driver"
    "regexp"
    "strings"
)

type YasDriver struct{}

func (driver *YasDriver) Open(dsn string) (driver.Conn, error) {
    conn := NewConnection()
    re1 := regexp.MustCompile(`^(.*?)/(.*?)@(.*)\?(.*?)$`)
    re2 := regexp.MustCompile(`^(.*?)/(.*?)@(.*)$`)
    autoCommit := false
    if re1.MatchString(dsn) {
        items := strings.Split(dsn, "?")
        dsn = items[0]
        options := strings.Split(items[1], "&")
        for _, opt := range options {
            o := strings.ToLower(opt)
            if o == "autocommit=1" || o == "autocommit=true" {
                autoCommit = true
            }
        }
    } else if !re2.MatchString(dsn) {
        return nil, ErrDsnNoStandard(dsn)
    }
    conn.Dsn = dsn
    if err := yasdbConnect(conn, autoCommit); err != nil {
        return nil, err
    }
    getAutoCommit(conn)
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
