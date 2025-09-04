package yasdb

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"git.yasdb.com/go/yasdb-go/assert"
)

type callParam struct {
	proName      string
	callTemplate string
	returnVal    string
	c1           string
	c2           string
	c3           string
	c4           string
	c5           string
	c6           string
	c7           string
	c8           string
	c9           string
	c10          string
}

func TestProcOutMalloc(t *testing.T) {
	function := `
	CREATE OR REPLACE function func_outparam(c1 out int,c2 out float,c3 out double,c4 out varchar,c5 out char,c6 out date,c7 out boolean,c8 out clob,c9 out rowid,c10 out json) return clob is 
	res clob; 
		v1 int := 943093745; 
		v2 float := 1506141.9; 
		v3 double := 107175737.7; 
		v4 varchar(20) := 'yasdb'; 
		v5 char(10) := 'yasql'; 
		v6 date := '2023-01-20'; 
		v7 boolean := false; 
		v8 clob := 'It gives me great pleasure to introduce our company.'; 
		v9 rowid := '1350:5:0:148:0'; 
		v10 json := '{"name":"Jack", "city":"Beijing","school":"TsingHua University"}'; 
	begin 
		c1 := v1; 
		c2 := v2; 
		c3 := v3; 
		c4 := v4; 
		c5 := v5; 
		c6 := v6; 
		c7 := v7; 
		c8 := v8; 
		c9 := v9; 
		c10 := v10; 
		res := c1||':'||c4||':'||c5||':'||c6||':'||c7||':'||c8||':'||c9; 
		return res; 
	end;
	`

	callTemplate := `
	BEGIN
		:RESULT := "SYS"."FUNC_OUTPARAM"(
			"C1" => :C1,
			"C2" => :C2,
			"C3" => :C3,
			"C4" => :C4,
			"C5" => :C5,
			"C6" => :C6,
			"C7" => :C7,
			"C8" => :C8,
			"C9" => :C9,
			"C10" => :C10);
	END;
	`

	createProcedute(t, function)

	for i := 0; i < 100; i++ {
		call := &callParam{
			proName:      "FUNC_OUTPARAM",
			callTemplate: callTemplate,
		}
		session := startOutParamDebug(t, call)
		debugContinue(t, session)
		closeDebugSession(session)

		expect := callParam{
			returnVal: "943093745:yasdb:yasql     :2023-01-20:false:It gives me great pleasure to introduce our company.:1350:5:0:148:0",
			c1:        "943093745",
			c2:        "1.50614188E+006",
			c3:        "1.071757377E+008",
			c4:        "yasdb",
			c5:        "yasql     ",
			c6:        "2023-01-20",
			c7:        "false",
			c8:        "It gives me great pleasure to introduce our company.",
			c9:        "1350:5:0:148:0",
			c10:       `{"city":"Beijing","name":"Jack","school":"TsingHua University"}`,
		}

		assert := assert.NewAssert(t)
		assert.Equal(call.returnVal, expect.returnVal)
		assert.Equal(call.c1, expect.c1)
		assert.Equal(call.c2, expect.c2)
		assert.Equal(call.c3, expect.c3)
		assert.Equal(call.c4, expect.c4)
		assert.Equal(call.c5, expect.c5)
		assert.Equal(call.c6, expect.c6)
		assert.Equal(call.c7, expect.c7)
		assert.Equal(call.c8, expect.c8)
		assert.Equal(call.c9, expect.c9)
		assert.Equal(call.c10, expect.c10)
	}
}

func startOutParamDebug(t *testing.T, callParam *callParam) *PlsqlDebug {
	oid, sid := queryObjIdAndSubId(t, callParam.proName)

	namedArgs := []any{
		sql.Named("RESULT", sql.Out{Dest: &callParam.returnVal}),
		sql.Named("C1", sql.Out{Dest: &callParam.c1}),
		sql.Named("C2", sql.Out{Dest: &callParam.c2}),
		sql.Named("C3", sql.Out{Dest: &callParam.c3}),
		sql.Named("C4", sql.Out{Dest: &callParam.c4}),
		sql.Named("C5", sql.Out{Dest: &callParam.c5}),
		sql.Named("C6", sql.Out{Dest: &callParam.c6}),
		sql.Named("C7", sql.Out{Dest: &callParam.c7}),
		sql.Named("C8", sql.Out{Dest: &callParam.c8}),
		sql.Named("C9", sql.Out{Dest: &callParam.c9}),
		sql.Named("C10", sql.Out{Dest: &callParam.c10}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callParam.callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}
	return debugsession
}

func debugContinue(t *testing.T, session *PlsqlDebug) {
	if err := session.Continue(); err != nil {
		t.Fatal(err)
	}
	if err := session.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}
}

func closeDebugSession(session *PlsqlDebug) {
	_ = session.Abort()
	_ = session.Close()
	session = nil
}

func TestBitInout(t *testing.T) {

	proc := `
	CREATE OR REPLACE PROCEDURE p_bit_inout(
		p_bit IN OUT BIT
	) AS
	BEGIN
		-- 0x02 0x10 0x40 0x20 0x80 0xc0 0x80 0x03
		p_bit := b'0000001000010000010000000010000010000000110000001000000000000011';
	END;
	`

	callTemplate := `
	BEGIN
		"P_BIT_INOUT"(
			"P_BIT" => :P_BIT);
	END;
	`

	procName := "P_BIT_INOUT"

	createProcedute(t, proc)

	oid, sid := queryObjIdAndSubId(t, procName)

	bit := []byte{1}
	value, err := NewOutputBindValue(&bit, WithTypeBit())
	if err != nil {
		t.Fatal(err)
	}

	// expect := []byte{0x03, 0x80, 0xc0, 0x80, 0x20, 0x40, 0x10, 0x02}
	namedArgs := []any{
		sql.Named("P_BIT", sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}
	closeDebugSession(debugsession)
	expected := "0000001000010000010000000010000010000000110000001000000000000011"
	actual := littleEndianBytesToBinaryString(bit)
	assert := assert.NewAssert(t)
	assert.Equal(actual, expected)

}

func TestBoolInout(t *testing.T) {

	proc := `
	CREATE OR REPLACE PROCEDURE p_boolean_inout(
		p_boolean IN OUT BOOLEAN
	) AS
	BEGIN
		p_boolean := TRUE;
	END;
	`

	callTemplate := `
	BEGIN
		"P_BOOLEAN_INOUT"(
			"P_BOOLEAN" => :P_BOOLEAN);
	END;
	`

	procName := "P_BOOLEAN_INOUT"

	createProcedute(t, proc)
	oid, sid := queryObjIdAndSubId(t, procName)

	var boolVal bool
	value, err := NewOutputBindValue(&boolVal, WithTypeBool())
	if err != nil {
		t.Fatal(err)
	}

	namedArgs := []any{
		sql.Named("P_BOOLEAN", sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}

	closeDebugSession(debugsession)

	expect := true
	assert := assert.NewAssert(t)
	assert.Equal(boolVal, expect)
}

func TestDateInout(t *testing.T) {

	proc := `
	CREATE OR REPLACE PROCEDURE p_date_inout(
		p_date IN OUT DATE
	) AS
	BEGIN
		p_date := TO_DATE('2023-10-01', 'YYYY-MM-DD');
	END;
	`

	callTemplate := `
	BEGIN
		"P_DATE_INOUT"(
			"P_DATE" => :P_DATE);
	END;
	`

	procName := "P_DATE_INOUT"

	timestamp := time.Now()

	createProcedute(t, proc)
	oid, sid := queryObjIdAndSubId(t, procName)

	value, err := NewOutputBindValue(&timestamp, WithTypeDate())
	if err != nil {
		t.Fatal(err)
	}

	namedArgs := []any{
		sql.Named("P_DATE", sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}

	closeDebugSession(debugsession)

	expectedTime, _ := time.Parse("2006-01-02", "2023-10-01")
	expect := expectedTime.UnixNano()
	actual := timestamp.UnixNano()
	assert := assert.NewAssert(t)
	assert.Equal(actual, expect)
}

func TestTimestampInout(t *testing.T) {

	proc := `
	CREATE OR REPLACE PROCEDURE p_timestamp_inout(
		p_timestamp IN OUT TIMESTAMP
	) AS
	BEGIN
		p_timestamp := TO_TIMESTAMP('2023-10-01 12:34:56', 'YYYY-MM-DD HH24:MI:SS');
	END;
	`

	callTemplate := `
	BEGIN
		"P_TIMESTAMP_INOUT"(
			"P_TIMESTAMP" => :P_TIMESTAMP);
	END;
	`

	procName := "P_TIMESTAMP_INOUT"

	timestamp := time.Now()

	createProcedute(t, proc)
	oid, sid := queryObjIdAndSubId(t, procName)

	value, err := NewOutputBindValue(&timestamp, WithTypeTimestamp())
	if err != nil {
		t.Fatal(err)
	}

	namedArgs := []any{
		sql.Named("P_TIMESTAMP", sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}

	closeDebugSession(debugsession)

	expectedTime, _ := time.Parse("2006-01-02 15:04:05", "2023-10-01 12:34:56")
	expect := expectedTime.UnixNano()
	actual := timestamp.UnixNano()
	assert := assert.NewAssert(t)
	assert.Equal(actual, expect)
}

// CI环境数据库不支持
func TestTimestampWithTimeZoneInout(t *testing.T) {

	test := newSqlTest(t)
	test.sqlGenInfo = &sqlGenInfo{}
	if !test.isToTimestampTzSupport() {
		return
	}
	proc := `
	CREATE OR REPLACE PROCEDURE p_timestamp_tz_inout(
		p_timestamp_tz IN OUT TIMESTAMP WITH TIME ZONE
	) AS
	BEGIN
		p_timestamp_tz := TO_TIMESTAMP_TZ('2023-10-01 12:34:56 +08:00','YYYY-MM-DD HH24:MI:SS TZH:TZM');
	END;
	`

	callTemplate := `
	BEGIN
		"P_TIMESTAMP_TZ_INOUT"(
			"P_TIMESTAMP_TZ" => :P_TIMESTAMP_TZ);
	END;
	`

	procName := "P_TIMESTAMP_TZ_INOUT"

	timestamp := time.Now()

	createProcedute(t, proc)
	oid, sid := queryObjIdAndSubId(t, procName)

	value, err := NewOutputBindValue(&timestamp, WithTypeTimestampTimeZone())
	if err != nil {
		t.Fatal(err)
	}

	namedArgs := []any{
		sql.Named("P_TIMESTAMP_TZ", sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}

	closeDebugSession(debugsession)

	expectedTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-10-01 12:34:56", time.Local)
	expect := expectedTime.UnixNano()
	actual := timestamp.UnixNano()
	assert := assert.NewAssert(t)
	assert.Equal(actual, expect)
}

func TestDoubleInout(t *testing.T) {

	proc := `
	CREATE OR REPLACE PROCEDURE p_double_inout(
		p_double IN OUT DOUBLE
	) AS
	BEGIN
		p_double := 1234567890.123456789;
	END;
	`

	callTemplate := `
	BEGIN
		"P_DOUBLE_INOUT"(
			"P_DOUBLE" => :P_DOUBLE);
	END;
	`

	procName := "P_DOUBLE_INOUT"

	doubleVal := float64(12.1)

	createProcedute(t, proc)
	oid, sid := queryObjIdAndSubId(t, procName)

	value, err := NewOutputBindValue(&doubleVal, WithTypeDouble())
	if err != nil {
		t.Fatal(err)
	}

	namedArgs := []any{
		sql.Named("P_DOUBLE", sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}

	closeDebugSession(debugsession)

	expected := float64(1234567890.123456789)

	assert := assert.NewAssert(t)
	assert.Equal(doubleVal, expected)
}

func TestFloatInout(t *testing.T) {

	proc := `
	CREATE OR REPLACE PROCEDURE p_float_inout(
		p_float IN OUT DOUBLE
	) AS
	BEGIN
		p_float := 1234.1234;
	END;
	`

	callTemplate := `
	BEGIN
		"P_FLOAT_INOUT"(
			"P_FLOAT" => :P_FLOAT);
	END;
	`

	procName := "P_FLOAT_INOUT"

	floatVal := float32(1234.1234)

	createProcedute(t, proc)
	oid, sid := queryObjIdAndSubId(t, procName)

	value, err := NewOutputBindValue(&floatVal, WithTypeFloat())
	if err != nil {
		t.Fatal(err)
	}

	namedArgs := []any{
		sql.Named("P_FLOAT", sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}

	closeDebugSession(debugsession)

	expected := float32(1234.1234)

	assert := assert.NewAssert(t)
	assert.Equal(floatVal, expected)
}

func TestTinyintInout(t *testing.T) {

	procFunc := `
					CREATE OR REPLACE PROCEDURE p_tinyint_inout(
						p_tinyint IN OUT TINYINT
					) AS
					BEGIN
						p_tinyint := 100;
					END;`
	procName := "P_TINYINT_INOUT"
	callTemplate := `
					BEGIN
						"P_TINYINT_INOUT"(
							"P_TINYINT" => :P_TINYINT);
					END;`
	bindName := "P_TINYINT"
	withBindType := WithTypeTinyint()
	bindVal := int8(0)
	expected := int8(100)

	createProcedute(t, procFunc)
	oid, sid := queryObjIdAndSubId(t, procName)

	value, err := NewOutputBindValue(&bindVal, withBindType)
	if err != nil {
		t.Fatal(err)
	}

	namedArgs := []any{
		sql.Named(bindName, sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}

	closeDebugSession(debugsession)

	assert := assert.NewAssert(t)
	assert.DeepEqual(bindVal, expected)

}

func TestSamllintInout(t *testing.T) {

	procFunc := `
					CREATE OR REPLACE PROCEDURE p_smallint_inout(
						p_smallint IN OUT SMALLINT
					) AS
					BEGIN
						p_smallint := 10000;
					END;`
	procName := "P_SMALLINT_INOUT"
	callTemplate := `
					BEGIN
						"P_SMALLINT_INOUT"(
							"P_SMALLINT" => :P_SMALLINT);
					END;`
	bindName := "P_SMALLINT"
	withBindType := WithTypeSmallInt()
	bindVal := int16(0)
	expected := int16(10000)

	createProcedute(t, procFunc)
	oid, sid := queryObjIdAndSubId(t, procName)

	value, err := NewOutputBindValue(&bindVal, withBindType)
	if err != nil {
		t.Fatal(err)
	}

	namedArgs := []any{
		sql.Named(bindName, sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}

	closeDebugSession(debugsession)

	assert := assert.NewAssert(t)
	assert.DeepEqual(bindVal, expected)

}

func TestIntegerInout(t *testing.T) {

	procFunc := `
				CREATE OR REPLACE PROCEDURE p_int_inout(
					p_int IN OUT INT
				) AS
				BEGIN
					p_int := 1000000;
				END;`
	procName := "P_INT_INOUT"
	callTemplate := `
					BEGIN
						"P_INT_INOUT"(
							"P_INT" => :P_INT);
					END;`
	bindName := "P_INT"
	withBindType := WithTypeInteger()
	bindVal := int32(0)
	expected := int32(1000000)

	createProcedute(t, procFunc)
	oid, sid := queryObjIdAndSubId(t, procName)

	value, err := NewOutputBindValue(&bindVal, withBindType)
	if err != nil {
		t.Fatal(err)
	}

	namedArgs := []any{
		sql.Named(bindName, sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}

	closeDebugSession(debugsession)

	assert := assert.NewAssert(t)
	assert.DeepEqual(bindVal, expected)

}

func TestBigintInout(t *testing.T) {

	procFunc := `
					CREATE OR REPLACE PROCEDURE p_bigint_inout(
						p_bigint in OUT BIGINT
					) AS
					BEGIN
						p_bigint := 9223372036854775807;
					END;`
	procName := "P_BIGINT_INOUT"
	callTemplate := `
					BEGIN
						"P_BIGINT_INOUT"(
							"P_BIGINT" => :P_BIGINT);
					END;`
	bindName := "P_BIGINT"
	withBindType := WithTypeBigInt()

	bindVal := int64(0)
	expected := int64(9223372036854775807)

	createProcedute(t, procFunc)
	oid, sid := queryObjIdAndSubId(t, procName)

	value, err := NewOutputBindValue(&bindVal, withBindType)
	if err != nil {
		t.Fatal(err)
	}

	namedArgs := []any{
		sql.Named(bindName, sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}

	closeDebugSession(debugsession)

	assert := assert.NewAssert(t)
	assert.DeepEqual(bindVal, expected)

}

func TestByteInout(t *testing.T) {

	type inoutCase struct {
		name          string
		procFunc      string
		procName      string
		callTemplate  string
		withBindTypes []outputBindOpt
		bindName      string
		bindVal       []byte
		expected      []byte
	}

	cases := []inoutCase{
		{
			name: "blob",
			procFunc: `
					CREATE OR REPLACE PROCEDURE p_blob_inout(
						p_blob IN OUT BLOB
					) AS
					BEGIN
						p_blob := HEXTORAW('1234ab');
					END;`,
			procName: "P_BLOB_INOUT",
			callTemplate: `
					BEGIN
						"P_BLOB_INOUT"(
							"P_BLOB" => :P_BLOB);
					END;`,
			bindName:      "P_BLOB",
			withBindTypes: append([]outputBindOpt{}, WithTypeBlob()),
			bindVal:       []byte{1},
			expected:      []byte{0x12, 0x34, 0xab},
		},
		{
			name: "json",
			procFunc: `
					CREATE OR REPLACE PROCEDURE p_json_inout(
						p_json IN OUT JSON
					) AS
					BEGIN
						p_json := '{"key": "value"}';
					END;`,
			procName: "P_JSON_INOUT",
			callTemplate: `
					BEGIN
						"P_JSON_INOUT"(
							"P_JSON" => :P_JSON);
					END;`,
			bindName:      "P_JSON",
			withBindTypes: append([]outputBindOpt{}, WithTypeBlob()),
			bindVal:       []byte("{}"),
			expected:      []byte("{\"key\":\"value\"}"),
		},
		{
			name: "raw",
			procFunc: `
					CREATE OR REPLACE PROCEDURE p_raw_inout(
						p_raw IN OUT RAW
					) AS
					BEGIN
						p_raw := HEXTORAW('01020304');
					END;`,
			procName: "P_RAW_INOUT",
			callTemplate: `
					BEGIN
						"P_RAW_INOUT"(
							"P_RAW" => :P_RAW);
					END;`,
			bindName:      "P_RAW",
			withBindTypes: append([]outputBindOpt{}, WithTypeBlob()),
			bindVal:       []byte("{}"),
			expected:      []byte{0x01, 0x02, 0x03, 0x04},
		},
		{
			name: "urowid",
			procFunc: `
					CREATE OR REPLACE PROCEDURE p_urowid_inout(
						p_urowid IN OUT UROWID
					) AS
					BEGIN
						p_urowid := HEXTORAW('01020304'); -- 需要实际行的 UROWID
					END;`,
			procName: "P_UROWID_INOUT",
			callTemplate: `
					BEGIN
						"P_UROWID_INOUT"(
							"P_UROWID" => :P_UROWID);
					END;`,
			bindName:      "P_UROWID",
			withBindTypes: append([]outputBindOpt{}, WithTypeBlob()),
			bindVal:       []byte{},
			expected:      []byte{0x01, 0x02, 0x03, 0x04},
		},
	}
	for i := range cases {
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			createProcedute(t, c.procFunc)
			oid, sid := queryObjIdAndSubId(t, c.procName)

			value, err := NewOutputBindValue(&c.bindVal, c.withBindTypes...)
			if err != nil {
				t.Fatal(err)
			}

			namedArgs := []any{
				sql.Named(c.bindName, sql.Out{Dest: value, In: true}),
			}

			debugsession, err := NewPlsqlDebug(testDsn, c.callTemplate, namedArgs...)
			if err != nil {
				t.Fatal(err)
			}
			if err := debugsession.Start(oid, sid); err != nil {
				t.Fatal(err)
			}

			if err := debugsession.Continue(); err != nil {
				t.Fatal(err)
			}

			if err := debugsession.GetBindOutValue(); err != nil {
				t.Fatal(err)
			}

			closeDebugSession(debugsession)

			assert := assert.NewAssert(t)
			assert.Equal(len(c.bindVal), len(c.expected))
			for i := range c.bindVal {
				assert.Equal(c.bindVal[i], c.expected[i])
			}
		})
	}
}

// nvarchar不能绑定
func TestStringInout(t *testing.T) {

	type inoutCase struct {
		name         string
		procFunc     string
		procName     string
		callTemplate string
		withBinds    []outputBindOpt
		bindName     string
		bindVal      string
		expected     string
	}

	cases := []inoutCase{
		{
			name: "char",
			procFunc: `
					CREATE OR REPLACE PROCEDURE p_char_inout(
						p_char IN OUT CHAR
					) AS
					BEGIN
						p_char := 'A';
					END;`,
			procName: "P_CHAR_INOUT",
			callTemplate: `
					BEGIN
						"P_CHAR_INOUT"(
							"P_CHAR" => :P_CHAR);
					END;`,
			bindName:  "P_CHAR",
			withBinds: append([]outputBindOpt{}, WithTypeVarchar(), WithBindSize(13)),
			expected:  "A",
		},
		{
			name: "varchar",
			procFunc: `
					CREATE OR REPLACE PROCEDURE p_varchar_inout(
						p_varchar IN OUT VARCHAR
					) AS
					BEGIN
						p_varchar := 'Hello, World!';
					END;`,
			procName: "P_VARCHAR_INOUT",
			callTemplate: `
					BEGIN
						"P_VARCHAR_INOUT"(
							"P_VARCHAR" => :P_VARCHAR);
					END;`,
			bindName:  "P_VARCHAR",
			withBinds: append([]outputBindOpt{}, WithTypeVarchar(), WithBindSize(20)),
			expected:  "Hello, World!",
		},
		{

			name: "nchar",
			procFunc: `
					CREATE OR REPLACE PROCEDURE p_nchar_inout(
						p_nchar IN OUT NCHAR
					) AS
					BEGIN
						p_nchar := 'Bbbbbbbbbffffffffffffffffffffffffffffffffffffffffffffffffff';
					END;`,
			procName: "P_NCHAR_INOUT",
			callTemplate: `
					BEGIN
						"P_NCHAR_INOUT"(
							"P_NCHAR" => :P_NCHAR);
					END;`,
			bindName:  "P_NCHAR",
			bindVal:   "BBB",
			withBinds: append([]outputBindOpt{}, WithTypeNvarchar(), WithBindSize(88)),
			expected:  "Bbbbbbbbbffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
		{
			name: "nvarchar",
			procFunc: `
					CREATE OR REPLACE PROCEDURE p_nvarchar_inout(
						p_nvarchar IN OUT NVARCHAR
					) AS
					BEGIN
						p_nvarchar := 'A';
					END;`,
			procName: "P_NVARCHAR_INOUT",
			callTemplate: `
					BEGIN
						"P_NVARCHAR_INOUT"(
							"P_NVARCHAR" => :P_NVARCHAR);
					END;`,
			bindName:  "P_NVARCHAR",
			bindVal:   "BBB",
			withBinds: append([]outputBindOpt{}, WithTypeNvarchar(), WithBindSize(50)),
			expected:  "A",
		},
	}
	for i := range cases {
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			createProcedute(t, c.procFunc)
			oid, sid := queryObjIdAndSubId(t, c.procName)

			value, err := NewOutputBindValue(&c.bindVal, c.withBinds...)
			if err != nil {
				t.Fatal(err)
			}

			namedArgs := []any{
				sql.Named(c.bindName, sql.Out{Dest: value, In: true}),
			}

			debugsession, err := NewPlsqlDebug(testDsn, c.callTemplate, namedArgs...)
			if err != nil {
				t.Fatal(err)
			}
			if err := debugsession.Start(oid, sid); err != nil {
				t.Fatal(err)
			}

			if err := debugsession.Continue(); err != nil {
				t.Fatal(err)
			}

			if err := debugsession.GetBindOutValue(); err != nil {
				t.Fatal(err)
			}

			closeDebugSession(debugsession)

			assert := assert.NewAssert(t)
			assert.DeepEqual(c.bindVal, c.expected)
		})
	}
}

func TestDsIntervalInout(t *testing.T) {

	proc := `
	CREATE OR REPLACE PROCEDURE p_interval_day_to_second_inout(
		p_interval_day_to_second IN OUT INTERVAL DAY TO SECOND
	) AS
	BEGIN
		p_interval_day_to_second := INTERVAL '1 12:34:56' DAY TO SECOND;
	END;
	`
	callTemplate := `
	BEGIN
		"P_INTERVAL_DAY_TO_SECOND_INOUT"(
			"P_INTERVAL_DAY_TO_SECOND" => :P_INTERVAL_DAY_TO_SECOND);
	END;
	`

	procName := "P_INTERVAL_DAY_TO_SECOND_INOUT"

	dsInterval := "50 10:30:59.999999"

	createProcedute(t, proc)
	oid, sid := queryObjIdAndSubId(t, procName)

	value, err := NewOutputBindValue(&dsInterval, WithTypeDSInterval())
	if err != nil {
		t.Fatal(err)
	}

	namedArgs := []any{
		sql.Named("P_INTERVAL_DAY_TO_SECOND", sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}

	closeDebugSession(debugsession)

	expected := "01 12:34:56.000000"

	assert := assert.NewAssert(t)
	assert.Equal(dsInterval, expected)
}

func TestNumberInout(t *testing.T) {

	proc := `
		CREATE OR REPLACE PROCEDURE p_number_inout(
			p_number IN OUT NUMBER
		) AS
		BEGIN
			p_number := 123.4;
		END;`
	callTemplate := `
		BEGIN
			"P_NUMBER_INOUT"(
				"P_NUMBER" => :P_NUMBER);
		END;`

	procName := "P_NUMBER_INOUT"

	numberVal := float64(1.2)

	createProcedute(t, proc)
	oid, sid := queryObjIdAndSubId(t, procName)

	value, err := NewOutputBindValue(&numberVal, WithTypeNumber())
	if err != nil {
		t.Fatal(err)
	}

	namedArgs := []any{
		sql.Named("P_NUMBER", sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}

	closeDebugSession(debugsession)

	expected := float64(123.4)

	assert := assert.NewAssert(t)
	assert.Equal(numberVal, expected)
}

func TestRowidInout(t *testing.T) {

	proc := `
		CREATE OR REPLACE PROCEDURE p_rowid_inout(
			p_rowid IN OUT ROWID
		) AS
		BEGIN
			p_rowid := '3574:4:0:156:0'; -- 需要实际行的 ROWID
		END;`
	callTemplate := `
		BEGIN
			"P_ROWID_INOUT"(
				"P_ROWID" => :P_ROWID);
		END;`

	procName := "P_ROWID_INOUT"

	rowidVal := string("2345:4:0:156:0")

	fmt.Println(len(rowidVal))

	createProcedute(t, proc)
	oid, sid := queryObjIdAndSubId(t, procName)

	value, err := NewOutputBindValue(&rowidVal, WithTypeRowid())
	if err != nil {
		t.Fatal(err)
	}

	namedArgs := []any{
		sql.Named("P_ROWID", sql.Out{Dest: value, In: true}),
	}

	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
	if err != nil {
		t.Fatal(err)
	}
	if err := debugsession.Start(oid, sid); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := debugsession.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}

	closeDebugSession(debugsession)

	expected := "3574:4:0:156:0"

	assert := assert.NewAssert(t)
	assert.Equal(rowidVal, expected)
}

// func TestYmIntervalInout(t *testing.T) {

// 	proc := `
// 	CREATE OR REPLACE PROCEDURE p_interval_year_to_month_inout(
// 		p_interval_year_to_month IN OUT INTERVAL YEAR TO MONTH
// 	) AS
// 	BEGIN
// 		p_interval_year_to_month := '1000-10';
// 	END;
// 	`
// 	callTemplate := `
// 	BEGIN
// 		"P_INTERVAL_YEAR_TO_MONTH_INOUT"(
// 			"P_INTERVAL_YEAR_TO_MONTH" => :P_INTERVAL_YEAR_TO_MONTH);
// 	END;
// 	`

// 	procName := "P_INTERVAL_YEAR_TO_MONTH_INOUT"

// 	ymInterval := "1000-10"

// 	createProcedute(t, proc)
// 	oid, sid := queryObjIdAndSubId(t, procName)

// 	value, err := NewOutputBindValue(&ymInterval, WithTypeYMInterval())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	namedArgs := []any{
// 		sql.Named("P_INTERVAL_YEAR_TO_MONTH", sql.Out{Dest: value, In: true}),
// 	}

// 	debugsession, err := NewPlsqlDebug(testDsn, callTemplate, namedArgs...)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := debugsession.Start(oid, sid); err != nil {
// 		t.Fatal(err)
// 	}

// 	if err := debugsession.Continue(); err != nil {
// 		t.Fatal(err)
// 	}

// 	if err := debugsession.GetBindOutValue(); err != nil {
// 		t.Fatal(err)
// 	}

// 	closeDebugSession(debugsession)

// 	// expected := float64(1234567890.123456789)

// 	// assert := assert.NewAssert(t)
// 	// assert.Equal(doubleVal, expected)
// }

func littleEndianBytesToBinaryString(data []byte) string {
	var result string
	for i := range data {
		// 将每个字节转换为 8 位二进制字符串（补前导零）
		result = fmt.Sprintf("%08b", data[i]) + result
	}
	return result
}
