package yasdb

import (
    "database/sql/driver"
    "io"
    "unsafe"
)

type YasRow struct {
    Elements   uint32
    Size       uint32
    dataSize   uint32
    IsValueSet bool
    IsArray    bool
    Data       unsafe.Pointer
    Indicator  int32
    DbType     int
    TransType  int
}

func NewYasRow(stmt *YasStmt, size uint32, dbType int) *YasRow {
    row := &YasRow{
        Elements: 1,
        Size:     size,
        DbType:   dbType,
    }
    row.dataSize = size * row.Elements
    row.Data = mallocBytes(size)
    row.Indicator = 0
    return row
}

type YasRows struct {
    stmt *YasStmt
}

func (r *YasRows) Columns() []string {
    return *r.stmt.Columns
}

func (r *YasRows) Close() error {
    freeFetchRows(r.stmt)
    return nil
}

func (r *YasRows) Next(dest []driver.Value) error {
    results, err := yasdbFetch(r.stmt)
    if err != nil {
        return err
    }
    if results == nil {
        freeFetchRows(r.stmt)
        return io.EOF
    }
    for i, d := range *results {
        dest[i] = d
    }
    return nil
}
