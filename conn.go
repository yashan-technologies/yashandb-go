package yasdb

import (
    "database/sql/driver"
)

// Conn for db open
type Connection struct {
    Env        YacHandle
    Conn       YacHandle
    Stmt       YacHandle
    Dsn        string
    Username   string
    AutoCommit bool
}

func NewConnection() *Connection {
    return &Connection{
        Env:  NewYacHandle(),
        Conn: NewYacHandle(),
        Stmt: NewYacHandle(),
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
    if err := yasdbFreeHandle(c.Stmt, YAC_HANDLE_STMT); err != nil {
        return err
    }
    if err := yasdbFreeHandle(c.Conn, YAC_HANDLE_DBC); err != nil {
        return err
    }
    return yasdbFreeHandle(c.Env, YAC_HANDLE_ENV)
}

// Begin begin
func (c *Connection) Begin() (driver.Tx, error) {
    return &YasTx{Conn: c}, nil
}

func (c *Connection) SetAutoCommit(auto bool) error {
    if auto == c.AutoCommit {
        return nil
    }
    if err := setAutoCommit(c, auto); err != nil {
        return err
    }
    c.AutoCommit = auto
    return nil
}

func (c *Connection) IsAutoCommit() bool {
    return c.AutoCommit
}
