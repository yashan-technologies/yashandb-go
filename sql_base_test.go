package yasdb

import (
    "database/sql"
    "fmt"
    "reflect"
    "testing"
    "time"
)

var (
    testDsn     = "sys/yasdb_123@192.168.6.177:1688?autocommit=true"
    tablePrefix = "gosqltest_"
)

type sqlTest struct {
    *sqlGenInfo
    *testing.T
    *sql.DB
}

type sqlGenInfo struct {
    tableName      string
    columnNameType [][2]string
    orderByColumn  string
    order          string
    execArgs       [][]interface{}
    queryArgs      []interface{}
    queryResult    [][]interface{}
    execResult     [][]interface{}
    query          string
}

func newSqlTest(t *testing.T) *sqlTest {
    db, err := sql.Open("yasdb", testDsn)
    if err != nil {
        t.Fatalf("error connecting: %v", err)
    }
    t.Parallel()
    return &sqlTest{T: t, DB: db}
}

func (st *sqlTest) dropTable() {
    dropTableSql := fmt.Sprintf("drop table if exists %s", st.tableName)
    st.mustExec(dropTableSql)
}

func (st *sqlTest) createTable() {
    columnDefineStr := ""
    flag := false
    for _, column := range st.columnNameType {
        if flag {
            columnDefineStr += fmt.Sprintf(", %s %s", column[0], column[1])
        } else {
            columnDefineStr += fmt.Sprintf("%s %s", column[0], column[1])
        }
        flag = true
    }
    createTableSql := fmt.Sprintf(`create table %s( %s )`, st.tableName, columnDefineStr)
    st.mustExec(createTableSql)
}

func (st *sqlTest) genTableTest() {
    st.dropTable()
    st.createTable()
}

func (st *sqlTest) mustExec(query string, args ...interface{}) sql.Result {
    res, err := st.Exec(query, args...)
    if err != nil {
        st.Fatalf("Error running %s [%v]: %v\n", query, args, err)
    }
    return res
}

func (st *sqlTest) mustQuery(query string, args ...interface{}) *sql.Rows {
    res, err := st.Query(query, args...)
    if err != nil {
        st.Fatalf("Error running %s [%v]: %v\n", query, args, err)
    }
    return res
}

func (st *sqlTest) runInsertTest() {
    columnName := ""
    bindValue := ""
    flag := false
    for _, column := range st.columnNameType {
        if flag {
            columnName += ", " + column[0]
            bindValue += ", ?"
        } else {
            columnName += column[0]
            bindValue += "?"
        }
        flag = true
    }
    st.query = fmt.Sprintf("insert into %s (%s) values(%s)", st.tableName, columnName, bindValue)
    st.runExecTest()
}

func (st *sqlTest) runExecTest() {
    for i := range st.execArgs {
        st.mustExec(st.query, st.execArgs[i]...)
    }
}

func (st *sqlTest) runSelectTest() {
    columnName := ""
    flag := false
    for _, column := range st.columnNameType {
        if flag {
            columnName += ", " + column[0]
        } else {
            columnName += column[0]
        }
        flag = true
    }
    st.query = fmt.Sprintf("select %s from %s %s", columnName, st.tableName, st.genOrderBy())
    st.runQueryTest()
}

func (st *sqlTest) runQueryTest() {
    actualResult, err := st.getRowsTest()
    if err != nil {
        st.Errorf("get rows error: %v - query: %s", err, st.query)
        return
    }
    if err = st.resultComparison(actualResult, st.queryResult); err != nil {
        st.Errorf("result comparison error: %v - query: %s", err, st.query)
    }
}

func (st *sqlTest) genOrderBy() string {
    orderBy := "order by "
    if st.orderByColumn == "" {
        orderBy += "1"
    } else {
        orderBy += fmt.Sprintf("%s %s", st.orderByColumn, st.order)
    }
    return orderBy
}

func (st *sqlTest) getRowsTest() ([][]interface{}, error) {
    rows := st.mustQuery(st.query, st.queryArgs...)
    defer rows.Close()
    columns, err := rows.Columns()
    if err != nil {
        rows.Close()
        return nil, fmt.Errorf("columns error: %v", err)
    }

    result := make([][]interface{}, 0)
    columnCount := len(columns)
    rowResult := make([]interface{}, columnCount)
    for rows.Next() {
        rowInterface := make([]interface{}, columnCount)
        for i := 0; i < columnCount; i++ {
            rowResult[i] = &rowInterface[i]
        }
        err = rows.Scan(rowResult...)
        if err != nil {
            return nil, fmt.Errorf("scan error: %v", err)
        }
        result = append(result, rowInterface)
    }

    return result, nil
}

func (st *sqlTest) resultComparison(actualResults [][]interface{}, expectedResults [][]interface{}) error {
    if actualResults == nil && expectedResults != nil {
        return fmt.Errorf("result is nil")
    }
    if len(actualResults) != len(expectedResults) {
        return fmt.Errorf("result rows len %v not equal to expected results len %v",
            len(actualResults), len(expectedResults))
    }

    testTypeTime := reflect.TypeOf(time.Time{})
    for i := 0; i < len(actualResults); i++ {
        for j := 0; j < len(actualResults[i]); j++ {
            actualResult := actualResults[i][j]
            expectedResult := expectedResults[i][j]
            bad := false
            type1 := reflect.TypeOf(actualResult)
            type2 := reflect.TypeOf(expectedResult)
            switch {
            case type1 == nil || type2 == nil:
                if type1 != type2 {
                    bad = true
                }
            case type1 == testTypeTime || type2 == testTypeTime:
                if type1 != type2 {
                    bad = true
                    break
                }
                time1 := actualResult.(time.Time)
                time2 := expectedResult.(time.Time)
                if !time1.Equal(time2) {
                    bad = true
                }
            case type1.Kind() == reflect.Slice || type2.Kind() == reflect.Slice:
                if !reflect.DeepEqual(actualResult, expectedResult) {
                    bad = true
                }
            default:
                if actualResult != expectedResult {
                    bad = true
                }
            }
            if bad {
                return fmt.Errorf("result - row %v, %v - received: %T, %v - expected: %T, %v",
                    i, j, actualResult, actualResult, expectedResult, expectedResult)
            }

        }
    }

    return nil
}

func runSqlTest(t *testing.T, fn func(st *sqlTest)) {
    db, err := sql.Open("yasdb", testDsn)
    if err != nil {
        t.Fatalf("error connecting: %v", err)
    }
    defer db.Close()
    fn(&sqlTest{T: t, DB: db})
}
