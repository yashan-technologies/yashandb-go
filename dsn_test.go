package yasdb

import (
    "fmt"
    "os"
    "reflect"
    "testing"
)

func TestParseDsn(t *testing.T) {
    t.Parallel()

    os.Mkdir("./test", 0755)

    var dsnTests = []struct {
        dsnStr      string
        expectedDSN *DataSourceName
    }{
        {`sys/yasdb_123@127.0.0.1:1688?autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: true, DataPath: ""}},
        {`sys/yasdb\@_123@127.0.0.1:1688?autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb@_123", Url: "127.0.0.1:1688", IsAutoCommit: true, DataPath: ""}},
        {`sys\//yasdb\@_123@127.0.0.1:1688?autocommit=true`, &DataSourceName{User: "sys/", Password: "yasdb@_123", Url: "127.0.0.1:1688", IsAutoCommit: true, DataPath: ""}},
        {"sys/yasdb_123@127.0.0.1:1688", &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: false, DataPath: ""}},
        {`sys\/\\/yasdb_123@127.0.0.1:1688?autocommit=true`, &DataSourceName{User: `sys/\`, Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: true, DataPath: ""}},
        {"sys/yasdb_123@1X7.0.0.1:1688", nil},
        {"sysyasdb_123@127.0.0.1:1688", nil},
        {"sys/yasdb_123127.0.0.1:1688", nil},
        {`./test?autocommit=true`, &DataSourceName{User: "sys", Password: "", Url: "", IsAutoCommit: true, DataPath: "./test"}},
        {`./test`, &DataSourceName{User: "sys", Password: "", Url: "", IsAutoCommit: false, DataPath: "./test"}},
        {`/a/b/c?autocommit=true`, nil},
    }

    for index, dt := range dsnTests {
        fmt.Println(dt.dsnStr)
        dsn, _ := ParseDSN(dt.dsnStr)
        if !reflect.DeepEqual(dsn, dt.expectedDSN) {
            t.Errorf("test case:%d. failed to parse dsn:%s  expected %+v, actual %+v", index, dt.dsnStr, dt.expectedDSN, dsn)
        }
    }
}
