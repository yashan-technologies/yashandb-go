package yasdb

import (
	"math"
	"testing"
)

func TestSelectNumberFromTable(t *testing.T) {
	t.Parallel()
	runSqlTest(t, testSelectNumberFromTable)
}

func testSelectNumberFromTable(t *sqlTest) {
	si := sqlGenInfo{
		tableName: "test_num",
		columnNameType: [][2]string{
			{"id", "int"},
			{"num1", "TINYINT"},
			{"num2", "SMALLINT"},
			{"num3", "INTEGER"},
			{"num4", "BIGINT"},
			{"num5", "FLOAT"},
			{"num6", "DOUBLE"},
			{"num7", "NUMBER"},
		},
		execArgs: [][]interface{}{
			{
				1,
				1,
				2,
				3,
				4,
				float32(5.1234),
				float64(6.12345),
				float64(7.123456),
			},
			{
				2,
				int8(math.MaxInt8),
				int16(math.MaxInt16),
				int32(math.MaxInt32),
				int64(math.MaxInt64),
				float32(math.MaxFloat32),
				float64(math.MaxFloat64),
				float64(7.123456),
			},
			{
				int32(3),
				int8(math.MinInt8),
				int16(math.MinInt16),
				int32(math.MinInt32),
				int64(math.MinInt64),
				float32(math.MaxFloat32 * -1),
				float64(math.MaxFloat64 * -1),
				float64(7.123456),
			},
		},
		queryResult: [][]interface{}{
			{
				int32(1),
				int8(1),
				int16(2),
				int32(3),
				int64(4),
				float32(5.1234),
				float64(6.12345),
				float64(7.123456)},
			{
				int32(2),
				int8(math.MaxInt8),
				int16(math.MaxInt16),
				int32(math.MaxInt32),
				int64(math.MaxInt64),
				float32(math.MaxFloat32),
				float64(math.MaxFloat64),
				float64(7.123456),
			},
			{
				int32(3),
				int8(math.MinInt8),
				int16(math.MinInt16),
				int32(math.MinInt32),
				int64(math.MinInt64),
				float32(math.MaxFloat32 * -1),
				float64(math.MaxFloat64 * -1),
				float64(7.123456),
			},
		},
	}
	t.sqlGenInfo = &si
	t.genTableTest()
	t.runInsertTest()
	t.runSelectTest()
}

func TestSelectNumberFromDual(t *testing.T) {
	t.Parallel()
	runSqlTest(t, testSelectNumberFromDual)
}

func testSelectNumberFromDual(t *sqlTest) {
	si := sqlGenInfo{}
	t.sqlGenInfo = &si

	// TINYINT: -128 ~ 127
	si.query = "select cast (:1 as TINYINT) from dual"

	si.queryArgs = []interface{}{-128}
	si.queryResult = [][]interface{}{{int8(-128)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{0}
	si.queryResult = [][]interface{}{{int8(0)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{1}
	si.queryResult = [][]interface{}{{int8(1)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{127}
	si.queryResult = [][]interface{}{{int8(127)}}
	t.runQueryTest()

	// SMALLINT: -32768 ~ 32767
	si.query = "select cast (:1 as SMALLINT) from dual"

	si.queryArgs = []interface{}{-32768}
	si.queryResult = [][]interface{}{{int16(-32768)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{0}
	si.queryResult = [][]interface{}{{int16(0)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{1}
	si.queryResult = [][]interface{}{{int16(1)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{32767}
	si.queryResult = [][]interface{}{{int16(32767)}}
	t.runQueryTest()

	// INTEGER: -2147483648 ~ 2147483647
	si.query = "select cast (:1 as INTEGER) from dual"

	si.queryArgs = []interface{}{-2147483648}
	si.queryResult = [][]interface{}{{int32(-2147483648)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{0}
	si.queryResult = [][]interface{}{{int32(0)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{1}
	si.queryResult = [][]interface{}{{int32(1)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{2147483647}
	si.queryResult = [][]interface{}{{int32(2147483647)}}
	t.runQueryTest()

	// BIGINT: -9223372036854775808 ~ -9223372036854775807
	si.query = "select cast (:1 as BIGINT) from dual"

	si.queryArgs = []interface{}{-9223372036854775808}
	si.queryResult = [][]interface{}{{int64(-9223372036854775808)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{0}
	si.queryResult = [][]interface{}{{int64(0)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{1}
	si.queryResult = [][]interface{}{{int64(1)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{9223372036854775807}
	si.queryResult = [][]interface{}{{int64(9223372036854775807)}}
	t.runQueryTest()

	// FLOAT: -3.402823E38 ~ -1.401298E-45, 0 , 1.401298E-45 ~ 3.402823E38
	si.query = "select cast (:1 as FLOAT) from dual"

	si.queryArgs = []interface{}{-3.402823e38}
	si.queryResult = [][]interface{}{{float32(-3.402823e38)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{-1.401298e-45}
	si.queryResult = [][]interface{}{{float32(-1.401298e-45)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{0}
	si.queryResult = [][]interface{}{{float32(0)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{1}
	si.queryResult = [][]interface{}{{float32(1)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{1.401298e-45}
	si.queryResult = [][]interface{}{{float32(1.401298e-45)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{3.402823e38}
	si.queryResult = [][]interface{}{{float32(3.402823e38)}}
	t.runQueryTest()

	// DOUBLE: -1.7976931348623E308 ~ -4.94065645841247E-324, 0 , 4.94065645841247E-324 ~ 1.7976931348623E308
	si.query = "select cast (:1 as DOUBLE) from dual"

	si.queryArgs = []interface{}{float64(-1.7976931348623e308)}
	si.queryResult = [][]interface{}{{float64(-1.7976931348623e308)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{float64(-4.94065645841247e-324)}
	si.queryResult = [][]interface{}{{float64(-4.94065645841247e-324)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{0}
	si.queryResult = [][]interface{}{{float64(0)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{1}
	si.queryResult = [][]interface{}{{float64(1)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{float64(4.94065645841247e-324)}
	si.queryResult = [][]interface{}{{float64(4.94065645841247e-324)}}
	t.runQueryTest()

	si.queryArgs = []interface{}{float64(1.7976931348623e308)}
	si.queryResult = [][]interface{}{{float64(1.7976931348623e308)}}
	t.runQueryTest()

	// NUMBER:
	si.query = "select cast (:1 as NUMBER(3,-2)) from dual"
	si.queryArgs = []interface{}{31401}
	si.queryResult = [][]interface{}{{float64(31400)}}
	t.runQueryTest()
}
