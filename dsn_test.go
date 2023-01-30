package yasdb

import (
    "reflect"
    "testing"
)

func TestParseDsn(t *testing.T) {
    t.Parallel()

    var dsnTests = []struct {
        dsnStr      string
        expectedDSN *DataSourceName
    }{
        {`sys/yasdb_123@127.0.0.1:1688?autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: true}},
        {`sys/yasdb\@_123@127.0.0.1:1688?autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb@_123", Url: "127.0.0.1:1688", IsAutoCommit: true}},
        {`sys\//yasdb\@_123@127.0.0.1:1688?autocommit=true`, &DataSourceName{User: "sys/", Password: "yasdb@_123", Url: "127.0.0.1:1688", IsAutoCommit: true}},
        {"sys/yasdb_123@127.0.0.1:1688", &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: false}},
        {`sys\/\\/yasdb_123@127.0.0.1:1688?autocommit=true`, &DataSourceName{User: `sys/\`, Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: true}},
        {"sys/yasdb_123@1X7.0.0.1:1688", nil},
        {"sysyasdb_123@127.0.0.1:1688", nil},
        {"sys/yasdb_123127.0.0.1:1688", nil},
    }

    for index, dt := range dsnTests {
        dsn, _ := ParseDSN(dt.dsnStr)
        if !reflect.DeepEqual(dsn, dt.expectedDSN) {
            t.Errorf("test case:%d. failed to parse dsn:%s  expected %+v, actual %+v", index, dt.dsnStr, dt.expectedDSN, dsn)
        }
    }
}
