package yasdb

import (
    "database/sql/driver"
    "unsafe"
)

type YasStmt struct {
    Conn      *Connection
    Stmt      AncHandle
    Columns   *[]string
    SqlType   uint32
    ArrSize   uint32
    RowCount  uint64
    IsOpen    bool
    fetchRows []*YasRow
    bindVals  []unsafe.Pointer
}

func NewYasStmt(conn *Connection) (*YasStmt, error) {
    stmt := &YasStmt{
        Conn: conn,
        Stmt: NewAncHandle(),
    }
    if err := yasdbStmtInit(stmt); err != nil {
        return nil, err
    }
    return stmt, nil
}

// Close  implement for stmt
func (stmt *YasStmt) Close() error {
    return yasdbFreeHandle(stmt.Stmt, ANC_HANDLE_STMT)
}

// Query  implement for Query
func (stmt *YasStmt) Query(args []driver.Value) (driver.Rows, error) {
    freeBindVals(stmt)
    if err := yasdbBindParams(stmt, args); err != nil {
        return nil, err
    }
    if err := yasdbExecute(stmt); err != nil {
        return nil, err
    }
    return &YasRows{stmt: stmt}, nil
}

// NumInput row numbers
func (stmt *YasStmt) NumInput() int {
    // don't know how many row numbers
    return -1
}

// Exec exec  implement
func (stmt *YasStmt) Exec(args []driver.Value) (driver.Result, error) {
    freeBindVals(stmt)
    if err := yasdbBindParams(stmt, args); err != nil {
        return nil, err
    }
    if err := yasdbExecute(stmt); err != nil {
        return nil, err
    }
    return &YasResult{stmt: stmt}, nil
}
