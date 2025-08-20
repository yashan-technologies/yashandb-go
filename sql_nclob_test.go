package yasdb

import "testing"

func TestNclob(t *testing.T) {
	runSqlTest(t, testNclob)
}

func testNclob(t *sqlTest) {
	si := sqlGenInfo{}
	t.sqlGenInfo = &si

	si = sqlGenInfo{
		tableName: "test_xml",
		columnNameType: [][2]string{
			{"id", "int"},
			{"c1", "nclob"},
		},
		execArgs: [][]interface{}{
			{1, nil},
			{2, "testNclob"},
		},
		queryResult: [][]interface{}{
			{int32(1), nil},
			{int32(2), "testNclob"},
		},
	}
	t.genTableTest()
	t.runInsertTest()
	t.runSelectTest()
}
