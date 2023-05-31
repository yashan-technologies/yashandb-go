package yasdb

import "testing"

func TestBool(t *testing.T) {
	runSqlTest(t, testBool)
}

func testBool(t *sqlTest) {
	si := sqlGenInfo{}
	t.sqlGenInfo = &si

	si = sqlGenInfo{
		tableName: "test_bool",
		columnNameType: [][2]string{
			{"id", "int"},
			{"b1", "boolean"},
		},
		execArgs: [][]interface{}{
			{1, nil},
			{2, "t"},
			{3, "yes"},
			{4, "y"},
			{5, "on"},
			{6, "1"},
			{7, true},
			{8, "false"},
			{9, "f"},
			{10, "no"},
			{11, "n"},
			{12, "off"},
			{13, "0"},
			{14, false},
			{15, 0},
			{16, 12},
			{17, -2},
		},
		queryResult: [][]interface{}{
			{int32(1), nil},
			{int32(2), true},
			{int32(3), true},
			{int32(4), true},
			{int32(5), true},
			{int32(6), true},
			{int32(7), true},
			{int32(8), false},
			{int32(9), false},
			{int32(10), false},
			{int32(11), false},
			{int32(12), false},
			{int32(13), false},
			{int32(14), false},
			{int32(15), false},
			{int32(16), true},
			{int32(17), true},
		},
	}
	t.genTableTest()
	t.runInsertTest()
	t.runSelectTest()
}
