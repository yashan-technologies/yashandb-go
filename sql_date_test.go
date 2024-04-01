package yasdb

import (
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	t.Parallel()
	runSqlTest(t, testDate)
}

func testDate(t *sqlTest) {
	si := sqlGenInfo{}
	t.sqlGenInfo = &si

	// test case1:
	nowTime := time.Now()
	t1c1 := time.Unix(nowTime.Unix(), 0)
	t1c2 := t1c1.AddDate(-1, -1, -1).Add(1000 * time.Microsecond)
	si = sqlGenInfo{
		tableName: "test_date1",
		columnNameType: [][2]string{
			{"c1", "DATE"},
			{"c2", "TIMESTAMP"},
		},
		execArgs: [][]interface{}{
			{t1c1, t1c2},
		},
		queryResult: [][]interface{}{
			{t1c1, t1c2.UnixMicro()},
		},
	}
	t.genTableTest()
	t.runInsertTest()
	t.runSelectTest()

	// test case2:
	t2c1_1 := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	t2c1_2 := time.Date(2200, 12, 31, 23, 59, 59, 0, time.UTC)
	t2c1_3 := time.Date(2022, 9, 27, 11, 54, 23, 0, time.UTC)
	si = sqlGenInfo{
		tableName: "test_date2",
		columnNameType: [][2]string{
			{"id", "int"},
			{"c1", "DATE"},
		},
		execArgs: [][]interface{}{
			{1, t2c1_1},
			{2, t2c1_2},
			{3, t2c1_3},
		},
		queryResult: [][]interface{}{
			{int32(1), t2c1_1},
			{int32(2), t2c1_2},
			{int32(3), t2c1_3},
		},
	}
	t.genTableTest()
	t.runInsertTest()
	t.runSelectTest()
}
