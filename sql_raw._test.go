package yasdb

import (
	"testing"
)

func TestRaw(t *testing.T) {
	t.Parallel()
	runSqlTest(t, testRaw)
}

func testRaw(t *sqlTest) {
	si := sqlGenInfo{}
	t.sqlGenInfo = &si

	raw1 := []byte("hello, YashanDB!")
	si = sqlGenInfo{
		tableName: "test_raw",
		columnNameType: [][2]string{
			{"id", "int"},
			{"raw1", "raw(16)"},
		},
		execArgs: [][]interface{}{
			{1, raw1},
		},
		queryResult: [][]interface{}{
			{int32(1), raw1},
		},
	}
	t.genTableTest()
	t.runInsertTest()
	t.runSelectTest()
}
