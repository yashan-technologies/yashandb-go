package yasdb

import (
	"strings"
	"testing"
)

func TestBlob(t *testing.T) {
	t.Parallel()
	runSqlTest(t, testBlob)
}

func testBlob(t *sqlTest) {
	si := sqlGenInfo{}
	t.sqlGenInfo = &si

	str := "你好, YashanDB!"
	blob1 := []byte("hello, YashanDB!")
	blob2 := []byte(str)
	blob3 := getBlobTestValue()                // Over 8192bytes
	blob4 := []byte(strings.Repeat(str, 8000)) // Over 8192bytes
	si = sqlGenInfo{
		tableName: "test_blob",
		columnNameType: [][2]string{
			{"id", "int"},
			{"blob1", "blob"},
		},
		execArgs: [][]interface{}{
			{1, blob1},
			{2, blob2},
			{3, blob3},
			{4, blob4},
		},
		queryResult: [][]interface{}{
			{int32(1), blob1},
			{int32(2), blob2},
			{int32(3), blob3},
			{int32(4), blob4},
		},
	}
	t.genTableTest()
	t.runInsertTest()
	t.runSelectTest()
}

func getBlobTestValue() []byte {
	return []byte(getClobTestValue())
}
