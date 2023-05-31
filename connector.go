/*
Copyright  2022, YashanDB and/or its affiliates. All rights reserved.
YashanDB Driver for golang is licensed under the terms of the mulan PSL v2.0

License: 	http://license.coscl.org.cn/MulanPSL2
Home page: 	https://www.yashandb.com/
*/

package yasdb

import (
	"context"
	"database/sql/driver"
)

type YasConnector struct {
}

// Driver returns the underlying Driver of the Connector,
// mainly to maintain compatibility with the Driver method
// on sql.DB.
func (connector *YasConnector) Driver() driver.Driver {
	return &YasdbDriver{}
}

// Connect returns a connection to the database.
func (connectot *YasConnector) Connect(ctx context.Context) (driver.Conn, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	yasConn := &YasConn{}
	return yasConn, nil
}
