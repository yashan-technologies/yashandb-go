package yasdb

import (
    "database/sql/driver"
)

// Conn for db open
type Connection struct {
    Env        AncHandle
    Conn       AncHandle
    Stmt       AncHandle
    Dsn        string
    Username   string
    autoCommit bool
}

func NewConnection() *Connection {
    return &Connection{
        Env:  NewAncHandle(),
        Conn: NewAncHandle(),
        Stmt: NewAncHandle(),
    }
}

// Prepare statement for prepare exec
func (c *Connection) Prepare(query string) (driver.Stmt, error) {
    stmt, err := NewYasStmt(c)
    if err != nil {
        return nil, err
    }
    if err := yasdbPrepare(stmt, query); err != nil {
        return nil, err
    }
    return stmt, nil
}

// Close close db connection
func (c *Connection) Close() error {
    if err := yasdbFreeHandle(c.Stmt, ANC_HANDLE_STMT); err != nil {
        return err
    }
    if err := yasdbFreeHandle(c.Conn, ANC_HANDLE_DBC); err != nil {
        return err
    }
    return yasdbFreeHandle(c.Env, ANC_HANDLE_ENV)
}

// Begin begin
func (c *Connection) Begin() (driver.Tx, error) {
    return &YasTx{Conn: c}, nil
}

func (c *Connection) AutoCommit(auto bool) error {
    if auto == c.autoCommit {
        return nil
    }
    if err := setAutoCommit(c, auto); err != nil {
        return err
    }
    c.autoCommit = auto
    return nil
}
