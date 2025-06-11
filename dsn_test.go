package yasdb

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestParseDsn(t *testing.T) {

	os.Mkdir("./test", 0o755)

	dsnTests := []struct {
		dsnStr      string
		expectedDSN *DataSourceName
	}{
		{`sys/yasdb_123@[::ffff:127.0.0.1]:1688?autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "[::ffff:127.0.0.1]:1688", IsAutoCommit: true, DataPath: "", directInsert: true}},
		{`sys/yasdb_123@127.0.0.1:1688?Autocommit=TRUE`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: true, DataPath: "", directInsert: true}},
		{`sys/yasdb\@_123@127.0.0.1:1688?autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb@_123", Url: "127.0.0.1:1688", IsAutoCommit: true, DataPath: "", directInsert: true}},
		{`sys\//yasdb\@_123@127.0.0.1:1688?autocommit=true`, &DataSourceName{User: "sys/", Password: "yasdb@_123", Url: "127.0.0.1:1688", IsAutoCommit: true, DataPath: "", directInsert: true}},
		{"sys/yasdb_123@127.0.0.1:1688", &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: false, DataPath: "", directInsert: true}},
		{`sys\/\\/yasdb_123@127.0.0.1:1688?autocommit=true`, &DataSourceName{User: `sys/\`, Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: true, DataPath: "", directInsert: true}},
		{`sys\/\\/yasdb_123@localhost:1688?autocommit=true`, &DataSourceName{User: `sys/\`, Password: "yasdb_123", Url: "localhost:1688", IsAutoCommit: true, DataPath: "", directInsert: true}},
		{"sysyasdb_123@127.0.0.1:1688", nil},
		{"sys/yasdb_123127.0.0.1:1688", nil},
		{`./test?autocommit=true`, &DataSourceName{User: "sys", Password: "", Url: "", IsAutoCommit: true, DataPath: "./test", directInsert: true}},
		{`./test`, &DataSourceName{User: "sys", Password: "", Url: "", IsAutoCommit: false, DataPath: "./test", directInsert: true}},
		{`/a/b/c?autocommit=true`, nil},
		{`sys/yasdb_123@[::1]:1688?autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "[::1]:1688", IsAutoCommit: true, DataPath: "", directInsert: true}},
		{`sys/yasdb_123@[fe80::8535:da9e:9517:1775]:7000?autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "[fe80::8535:da9e:9517:1775]:7000", IsAutoCommit: true, DataPath: "", directInsert: true}},
		{`sys/yasdb_123@[fe80::d5a0:6043:483c:4bfd%ens192]:1688?autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "[fe80::d5a0:6043:483c:4bfd%ens192]:1688", IsAutoCommit: true, DataPath: "", directInsert: true}},
		{"sys/Cod-2022@primary:192.168.6.177:2300,192.168.6.177:2302,192.168.6.177:2304", &DataSourceName{User: "sys", Password: "Cod-2022", Url: "primary:192.168.6.177:2300,192.168.6.177:2302,192.168.6.177:2304", directInsert: true}},
		{"sys/Cod-2022@192.168.6.177:2300,192.168.6.177:2302,192.168.6.177:2304", &DataSourceName{User: "sys", Password: "Cod-2022", Url: "192.168.6.177:2300,192.168.6.177:2302,192.168.6.177:2304", directInsert: true}},
		{`sys/yasdb_123@127.0.0.1:1688?ukey_name=dba&ukey_pin=Cod-2022`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688?ukey_name=dba&ukey_pin=Cod-2022", IsAutoCommit: false, DataPath: "", ukeyName: "dba", ukeyPin: "Cod-2022", directInsert: true}},
		{`sys/yasdb_123@127.0.0.1:1688?ukey_name=1&ukey_pin=123&autocommit=true`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688?ukey_name=1&ukey_pin=123", IsAutoCommit: true, DataPath: "", ukeyName: "1", ukeyPin: "123", directInsert: true}},
		{`sys/yasdb_123@127.0.0.1:1688?ukey_name=1`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688?ukey_name=1", IsAutoCommit: false, DataPath: "", ukeyName: "1", directInsert: true}},
		{`sys/yasdb_123@127.0.0.1:1688?ukey_pin=123`, &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688?ukey_pin=123", IsAutoCommit: false, DataPath: "", ukeyPin: "123", directInsert: true}},
		{"sys/Cod-2022@LOADBALANCE:192.168.6.177:2300,192.168.6.177:2302,192.168.6.177:2304", &DataSourceName{User: "sys", Password: "Cod-2022", Url: "LOADBALANCE:192.168.6.177:2300,192.168.6.177:2302,192.168.6.177:2304", directInsert: true}},
		{"sys/Cod-2022@loadbalance:192.168.6.177:2300,192.168.6.177:2302,192.168.6.177:2304", &DataSourceName{User: "sys", Password: "Cod-2022", Url: "loadbalance:192.168.6.177:2300,192.168.6.177:2302,192.168.6.177:2304", directInsert: true}},
		{"sys/yasdb_123@127.0.0.1:1688?heartbeat_enable=true", &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: false, heartbeatEnable: true, directInsert: true}},
		{"sys/yasdb_123@127.0.0.1:1688?compat_vector=yashan", &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: false, compatVector: "yashan", directInsert: true}},
		{"sys/yasdb_123@127.0.0.1:1688?compat_vector=yashan", &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: false, compatVector: "yashan", directInsert: true}},
		{"sys/yasdb_123@127.0.0.1:1688?compat_vector=mysql", &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: false, compatVector: "mysql", directInsert: true}},
		{"sys/yasdb_123@127.0.0.1:1688?compat_vector=null", &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: false, compatVector: "null", directInsert: true}},
		{"sys/yasdb_123@127.0.0.1:1688?compat_vector=yashan&autocommit=true&number_as_string=true", &DataSourceName{User: "sys", Password: "yasdb_123", Url: "127.0.0.1:1688", IsAutoCommit: true, numberAsString: true, compatVector: "yashan", directInsert: true}},
	}

	for index, dt := range dsnTests {
		fmt.Println(dt.dsnStr)
		dsn, _ := ParseDSN(dt.dsnStr)
		if !reflect.DeepEqual(dsn, dt.expectedDSN) {
			t.Errorf("test case:%d. failed to parse dsn:%s  expected %+v, actual %+v", index, dt.dsnStr, dt.expectedDSN, dsn)
		}
	}
}
