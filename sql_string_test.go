package yasdb

import (
    "strconv"
    "strings"
    "testing"
)

func TestClob(t *testing.T) {
    t.Parallel()
    runSqlTest(t, testClob)
}

func testClob(t *sqlTest) {
    si := sqlGenInfo{}
    t.sqlGenInfo = &si

    clob1 := "hello, YashanDB!"
    clob2 := "你好, YashanDB!"
    clob3 := getClobTestValue()          // Over 8192bytes
    clob4 := strings.Repeat(clob2, 8192) // Over 8192bytes
    si = sqlGenInfo{
        tableName: "test_clob",
        columnNameType: [][2]string{
            {"id", "int"},
            {"clob1", "clob"},
        },
        execArgs: [][]interface{}{
            {1, clob1},
            {2, clob2},
            {3, clob3},
            {4, clob4},
        },
        queryResult: [][]interface{}{
            {int32(1), clob1},
            {int32(2), clob2},
            {int32(3), clob3},
            {int32(4), clob4},
        },
    }
    t.genTableTest()
    t.runInsertTest()
    t.runSelectTest()
}

// Over 8192bytes
func getClobTestValue() string {
    str := ""
    for i := 0; i < 10000; i++ {
        str += strconv.Itoa(i) + ","
        if len(str) > 10000 {
            break
        }
    }
    return str
}

func TestVarchar(t *testing.T) {
    t.Parallel()
    runSqlTest(t, testVarchar)
}

func testVarchar(t *sqlTest) {
    si := sqlGenInfo{}
    t.sqlGenInfo = &si

    //test case1:
    t1c1 := "01234567899876543210"
    t1c2 := "中国深圳龙华"
    si = sqlGenInfo{
        tableName: "varchar_test1",
        columnNameType: [][2]string{
            {"id", "int"},
            {"c1", "varchar(20)"},
            {"c2", "varchar(20)"},
        },
        execArgs: [][]interface{}{
            {1, t1c1, t1c2},
        },
        queryResult: [][]interface{}{
            {int32(1), t1c1, t1c2},
        },
    }
    t.genTableTest()
    t.runInsertTest()
    t.runSelectTest()

    // test case2:
    t2c1_1 := strings.Repeat("0123456789", 8000/10)
    t2c1_2 := strings.Repeat("崖山数据库", 8000/(5*3))
    t2c1_3 := ""
    si = sqlGenInfo{
        tableName: "varchar_test2",
        columnNameType: [][2]string{
            {"id", "int"},
            {"c1", "varchar(8000)"},
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

func TestChar(t *testing.T) {
    t.Parallel()
    runSqlTest(t, testChar)
}

func testChar(t *sqlTest) {
    si := sqlGenInfo{}
    t.sqlGenInfo = &si

    // test case1:
    t1c1 := "01234567899876543210"
    t1c2 := "中国深圳龙华"
    si = sqlGenInfo{
        tableName: "char_test1",
        columnNameType: [][2]string{
            {"id", "int"},
            {"c1", "char(20)"},
            {"c2", "char(18)"},
        },
        execArgs: [][]interface{}{
            {1, t1c1, t1c2},
        },
        queryResult: [][]interface{}{
            {int32(1), t1c1, t1c2},
        },
    }
    t.genTableTest()
    t.runInsertTest()
    t.runSelectTest()

    // test case2:
    t2c1_1 := strings.Repeat("0123456789", 8000/10)
    t2c1_2 := strings.Repeat("崖山数据库", 8000/(5*3))
    t2c1_3 := ""
    si = sqlGenInfo{
        tableName: "char_test2",
        columnNameType: [][2]string{
            {"id", "int"},
            {"c1", "char(8000)"},
        },
        execArgs: [][]interface{}{
            {1, t2c1_1},
            {2, t2c1_2},
            {3, t2c1_3},
        },
        queryResult: [][]interface{}{
            {int32(1), t2c1_1},
            {int32(2), t2c1_2 + strings.Repeat(" ", 8000-len(t2c1_2))},
            {int32(3), strings.Repeat(t2c1_3, 8000)},
        },
    }
    t.genTableTest()
    t.runInsertTest()
    t.runSelectTest()
}

func TestEmojiCharacters(t *testing.T) {
    runSqlTest(t, testEmojiCharacters)
}

func testEmojiCharacters(t *sqlTest) {
    si := sqlGenInfo{}
    t.sqlGenInfo = &si

    si = sqlGenInfo{
        tableName: "test_emoji",
        columnNameType: [][2]string{
            {"id", "int"},
            {"c1", "varchar(20)"},
        },
        execArgs: [][]interface{}{
            {1, "😀"},
            {2, "😄"},
            {3, "😁"},
            {4, "😇"},
            {5, "🥰"},
            {6, "🚮"},
            {7, "🚰"},
            {8, "⚠️"},
            {9, "📵"},
            {10, "🧲"},
        },
        queryResult: [][]interface{}{
            {int32(1), "😀"},
            {int32(2), "😄"},
            {int32(3), "😁"},
            {int32(4), "😇"},
            {int32(5), "🥰"},
            {int32(6), "🚮"},
            {int32(7), "🚰"},
            {int32(8), "⚠️"},
            {int32(9), "📵"},
            {int32(10), "🧲"},
        },
    }
    t.genTableTest()
    t.runInsertTest()
    t.runSelectTest()

    si.query = "select cast (? as varchar(20)) from dual"

    si.queryArgs = []interface{}{"🙂"}
    si.queryResult = [][]interface{}{{"🙂"}}
    t.runQueryTest()

    si.queryArgs = []interface{}{"😂"}
    si.queryResult = [][]interface{}{{"😂"}}
    t.runQueryTest()

    si.queryArgs = []interface{}{"🤣"}
    si.queryResult = [][]interface{}{{"🤣"}}
    t.runQueryTest()

    si.queryArgs = []interface{}{"😶‍🌫️"}
    si.queryResult = [][]interface{}{{"😶‍🌫️"}}
    t.runQueryTest()

    si.queryArgs = []interface{}{"😮‍💨"}
    si.queryResult = [][]interface{}{{"😮‍💨"}}
    t.runQueryTest()

    si.queryArgs = []interface{}{"😼"}
    si.queryResult = [][]interface{}{{"😼"}}
    t.runQueryTest()
}
