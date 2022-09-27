package yasdb

import (
    "context"
    "database/sql/driver"
)

type YasConnector struct {
}

func (connector *YasConnector) Driver() driver.Driver {
    return &YasdbDriver{}
}

func (connectot *YasConnector) Connect(ctx context.Context) (driver.Conn, error) {
    if ctx.Err() != nil {
        return nil, ctx.Err()
    }
    yasConn := NewYasConn()
    return yasConn, nil
}
