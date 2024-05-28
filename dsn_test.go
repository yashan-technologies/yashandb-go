package yasdb

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestParseDsn(t *testing.T) {
	t.Parallel()

	os.Mkdir("./test", 0o755)

	dsnTests := []struct {
		dsnStr      string
		expectedDSN *DataSourceName
	}{
		{`sys/yasdb_123@[::ffff:127.0.0.1]:1688?autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "[::ffff:127.0.0.1]:1688", IsAutoCommit: true, DataPath: ""}},
		{`sys/yasdb_123@127.0.0.1:1688?Autocommit=TRUE`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: true, DataPath: ""}},
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
		{`sys/yasdb_123@[::1]:1688?autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "[::1]:1688", IsAutoCommit: true, DataPath: ""}},
		{`sys/yasdb_123@[fe80::8535:da9e:9517:1775]:7000?autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "[fe80::8535:da9e:9517:1775]:7000", IsAutoCommit: true, DataPath: ""}},
		{`sys/yasdb_123@[fe80::d5a0:6043:483c:4bfd%ens192]:1688?autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "[fe80::d5a0:6043:483c:4bfd%ens192]:1688", IsAutoCommit: true, DataPath: ""}},
		{"sys/Cod-2022@primary:192.168.6.177:2300,192.168.6.177:2302,192.168.6.177:2304", &DataSourceName{User: "sys", Password: "Cod-2022", Url: "primary:192.168.6.177:2300,192.168.6.177:2302,192.168.6.177:2304"}},
		{"sys/Cod-2022@192.168.6.177:2300,192.168.6.177:2302,192.168.6.177:2304", &DataSourceName{User: "sys", Password: "Cod-2022", Url: "192.168.6.177:2300,192.168.6.177:2302,192.168.6.177:2304"}},
		{`sys/yasdb_123@127.0.0.1:1688?ukey_name=dba&ukey_pin=Cod-2022`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688?ukey_name=dba&ukey_pin=Cod-2022", IsAutoCommit: false, DataPath: "", ukeyName: "dba", ukeyPin: "Cod-2022"}},
		{`sys/yasdb_123@127.0.0.1:1688?ukey_name=1&ukey_pin=123&autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688?ukey_name=1&ukey_pin=123", IsAutoCommit: true, DataPath: "", ukeyName: "1", ukeyPin: "123"}},
		{`sys/yasdb_123@127.0.0.1:1688?ukey_name=1`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688?ukey_name=1", IsAutoCommit: false, DataPath: "", ukeyName: "1"}},
		{`sys/yasdb_123@127.0.0.1:1688?ukey_pin=123`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688?ukey_pin=123", IsAutoCommit: false, DataPath: "", ukeyPin: "123"}},
	}

	for index, dt := range dsnTests {
		fmt.Println(dt.dsnStr)
		dsn, _ := ParseDSN(dt.dsnStr)
		if !reflect.DeepEqual(dsn, dt.expectedDSN) {
			t.Errorf("test case:%d. failed to parse dsn:%s  expected %+v, actual %+v", index, dt.dsnStr, dt.expectedDSN, dsn)
		}
	}
}
