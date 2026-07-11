package yasdb

import (
	"database/sql"
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// A. OUT NUMBER tests
//
// Each scenario creates a dedicated procedure that assigns a known value
// to an OUT NUMBER parameter, then verifies the Go-side destination
// receives the expected result.
// ---------------------------------------------------------------------------

func TestNumberOut_Float64_Integer(t *testing.T) {
	runSqlTest(t, testNumberOutFloat64Integer)
}

func testNumberOutFloat64Integer(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_out_f64_int(p_out OUT NUMBER) AS
BEGIN
    p_out := 42;
END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_out_f64_int")

	var f64 float64
	value, err := NewOutputBindValue(&f64, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	st.mustExec(
		`BEGIN "P_NUMBER_OUT_F64_INT"(:1); END;`,
		sql.Out{Dest: value},
	)
	if f64 != 42.0 {
		st.Fatalf("got %v, want 42.0", f64)
	}
}

func TestNumberOut_Int64_LargeInteger(t *testing.T) {
	runSqlTest(t, testNumberOutInt64LargeInteger)
}

func testNumberOutInt64LargeInteger(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_out_i64_large(p_out OUT NUMBER) AS
BEGIN
    p_out := 9223372036854775807;
END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_out_i64_large")

	var i64 int64
	value, err := NewOutputBindValue(&i64, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	st.mustExec(
		`BEGIN "P_NUMBER_OUT_I64_LARGE"(:1); END;`,
		sql.Out{Dest: value},
	)
	if i64 != math.MaxInt64 {
		st.Fatalf("got %d, want %d", i64, int64(math.MaxInt64))
	}
}

func TestNumberOut_String_SuperLargeInteger(t *testing.T) {
	runSqlTest(t, testNumberOutStringSuperLargeInteger)
}

func testNumberOutStringSuperLargeInteger(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_out_str_super(p_out OUT NUMBER) AS
BEGIN
    p_out := 123456789012345678901234567890;
END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_out_str_super")

	var s string
	value, err := NewOutputBindValue(&s, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	st.mustExec(
		`BEGIN "P_NUMBER_OUT_STR_SUPER"(:1); END;`,
		sql.Out{Dest: value},
	)
	expected := "123456789012345678901234567890"
	if s != expected {
		st.Fatalf("got %q, want %q", s, expected)
	}
}

func TestNumberOut_Number_LargeInteger(t *testing.T) {
	runSqlTest(t, testNumberOutNumberLargeInteger)
}

func testNumberOutNumberLargeInteger(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_out_num_large(p_out OUT NUMBER) AS
BEGIN
    p_out := 12345678901234567890;
END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_out_num_large")

	var n Number
	value, err := NewOutputBindValue(&n, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	st.mustExec(
		`BEGIN "P_NUMBER_OUT_NUM_LARGE"(:1); END;`,
		sql.Out{Dest: value},
	)
	expected := "12345678901234567890"
	if !n.Valid() {
		st.Fatal("expected valid Number, got NULL")
	}
	if n.String() != expected {
		st.Fatalf("got %q, want %q", n.String(), expected)
	}
}

func TestNumberOut_Int64_DecimalReturnsError(t *testing.T) {
	runSqlTest(t, testNumberOutInt64DecimalReturnsError)
}

func testNumberOutInt64DecimalReturnsError(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_out_i64_dec(p_out OUT NUMBER) AS
BEGIN
    p_out := 123.4;
END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_out_i64_dec")

	var i64 int64
	value, err := NewOutputBindValue(&i64, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	_, err = st.Exec(
		`BEGIN "P_NUMBER_OUT_I64_DEC"(:1); END;`,
		sql.Out{Dest: value},
	)
	if err == nil {
		st.Fatal("expected error when OUT NUMBER decimal 123.4 is read into *int64, got nil")
	}
}

func TestNumberOut_NULL_Float64(t *testing.T) {
	runSqlTest(t, testNumberOutNullFloat64)
}

func testNumberOutNullFloat64(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_out_null_f64(p_out OUT NUMBER) AS
BEGIN
    p_out := NULL;
END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_out_null_f64")

	f64 := float64(99.9) // pre-set; should become 0 after NULL OUT
	value, err := NewOutputBindValue(&f64, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	st.mustExec(
		`BEGIN "P_NUMBER_OUT_NULL_F64"(:1); END;`,
		sql.Out{Dest: value},
	)
	if f64 != 0 {
		st.Fatalf("got %v, want 0 (NULL → float64(0))", f64)
	}
}

func TestNumberOut_NULL_Int64(t *testing.T) {
	runSqlTest(t, testNumberOutNullInt64)
}

func testNumberOutNullInt64(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_out_null_i64(p_out OUT NUMBER) AS
BEGIN
    p_out := NULL;
END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_out_null_i64")

	i64 := int64(999) // pre-set; should become 0 after NULL OUT
	value, err := NewOutputBindValue(&i64, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	st.mustExec(
		`BEGIN "P_NUMBER_OUT_NULL_I64"(:1); END;`,
		sql.Out{Dest: value},
	)
	if i64 != 0 {
		st.Fatalf("got %d, want 0 (NULL → int64(0))", i64)
	}
}

func TestNumberOut_NULL_String(t *testing.T) {
	runSqlTest(t, testNumberOutNullString)
}

func testNumberOutNullString(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_out_null_str(p_out OUT NUMBER) AS
BEGIN
    p_out := NULL;
END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_out_null_str")

	s := "preset" // pre-set; should become "" after NULL OUT
	value, err := NewOutputBindValue(&s, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	st.mustExec(
		`BEGIN "P_NUMBER_OUT_NULL_STR"(:1); END;`,
		sql.Out{Dest: value},
	)
	if s != "" {
		st.Fatalf("got %q, want empty string (NULL → \"\")", s)
	}
}

func TestNumberOut_NULL_Number(t *testing.T) {
	runSqlTest(t, testNumberOutNullNumber)
}

func testNumberOutNullNumber(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_out_null_num(p_out OUT NUMBER) AS
BEGIN
    p_out := NULL;
END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_out_null_num")

	n := NewNumber("999") // pre-set; should become Number{} after NULL OUT
	value, err := NewOutputBindValue(&n, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	st.mustExec(
		`BEGIN "P_NUMBER_OUT_NULL_NUM"(:1); END;`,
		sql.Out{Dest: value},
	)
	if n.Valid() {
		st.Fatalf("expected NULL Number (Valid()=false), got Valid()=true, String()=%q", n.String())
	}
	if !n.IsNull() {
		st.Fatal("expected IsNull()=true after NULL OUT")
	}
}

// ---------------------------------------------------------------------------
// B. IN OUT NUMBER tests
//
// Each scenario creates a procedure with an IN OUT NUMBER parameter,
// passes an initial value, and verifies the output after PL/SQL processing.
// ---------------------------------------------------------------------------

func TestNumberInOut_Int64(t *testing.T) {
	runSqlTest(t, testNumberInOutInt64)
}

func testNumberInOutInt64(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_inout_i64(p_inout IN OUT NUMBER) AS
BEGIN
    p_inout := p_inout + 1;
END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_inout_i64")

	// Use a value beyond float64 exact range (>2^53) to verify text path precision.
	// 9223372036854775806 (MaxInt64-1) + 1 = 9223372036854775807 (MaxInt64).
	// Through float64, this would lose precision and not produce the exact result.
	i64 := int64(9223372036854775806)
	value, err := NewOutputBindValue(&i64, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	st.mustExec(
		`BEGIN "P_NUMBER_INOUT_I64"(:1); END;`,
		sql.Out{Dest: value, In: true},
	)
	if i64 != 9223372036854775807 {
		st.Fatalf("got %d, want 9223372036854775807", i64)
	}
}

func TestNumberInOut_String_LargeInteger(t *testing.T) {
	runSqlTest(t, testNumberInOutStringLargeInteger)
}

func testNumberInOutStringLargeInteger(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_inout_str(p_inout IN OUT NUMBER) AS
BEGIN
    p_inout := p_inout + 1;
END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_inout_str")

	// IN OUT input uses VARCHAR bind path to avoid float64 precision loss.
	// This 20-digit integer exceeds float64 exact range (>2^53 = 9007199254740992).
	// The database handles VARCHAR ↔ NUMBER conversion; the output is also lossless for *string.
	s := "12345678901234567890"
	value, err := NewOutputBindValue(&s, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	st.mustExec(
		`BEGIN "P_NUMBER_INOUT_STR"(:1); END;`,
		sql.Out{Dest: value, In: true},
	)
	expected := "12345678901234567891"
	if s != expected {
		st.Fatalf("got %q, want %q", s, expected)
	}
}

func TestNumberInOut_Number_Reassign(t *testing.T) {
	runSqlTest(t, testNumberInOutNumberReassign)
}

func testNumberInOutNumberReassign(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_inout_num(p_inout IN OUT NUMBER) AS
BEGIN
    p_inout := 12345678901234567890;
END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_inout_num")

	n := NewNumber("1")
	value, err := NewOutputBindValue(&n, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	st.mustExec(
		`BEGIN "P_NUMBER_INOUT_NUM"(:1); END;`,
		sql.Out{Dest: value, In: true},
	)
	expected := "12345678901234567890"
	if !n.Valid() {
		st.Fatal("expected valid Number, got NULL")
	}
	if n.String() != expected {
		st.Fatalf("got %q, want %q", n.String(), expected)
	}
}

func TestNumberInOut_Number_NullInput(t *testing.T) {
	runSqlTest(t, testNumberInOutNumberNullInput)
}

func testNumberInOutNumberNullInput(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	// The procedure assigns a concrete value even when input is NULL.
	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_inout_null(p_inout IN OUT NUMBER) AS
BEGIN
    IF p_inout IS NULL THEN
        p_inout := 0;
    ELSE
        p_inout := p_inout + 1;
    END IF;
END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_inout_null")

	// NULL Number (zero value) → procedure detects NULL and assigns 0.
	n := Number{} // NULL
	value, err := NewOutputBindValue(&n, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	st.mustExec(
		`BEGIN "P_NUMBER_INOUT_NULL"(:1); END;`,
		sql.Out{Dest: value, In: true},
	)

	// After procedure processing, NULL input → procedure assigns 0.
	if !n.Valid() {
		st.Fatal("expected valid Number after procedure assigns 0, got NULL")
	}
	if n.String() != "0" {
		st.Fatalf("got %q, want \"0\" after NULL input + procedure assignment", n.String())
	}
}

// ---------------------------------------------------------------------------
// H. Regression test: NUMBER OUT followed by other OUT parameters
//
// This catches a bug where getBindValueDest returned early on NUMBER OUT,
// preventing subsequent OUT parameters from being processed.
// ---------------------------------------------------------------------------

func TestNumberOut_MultipleOutParams(t *testing.T) {
	runSqlTest(t, testNumberOutMultipleOutParams)
}

func testNumberOutMultipleOutParams(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	// Procedure with NUMBER OUT followed by VARCHAR2 OUT
	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_out_multi(
		p_num OUT NUMBER,
		p_str OUT VARCHAR2
	) AS
	BEGIN
		p_num := 42;
		p_str := 'hello';
	END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_out_multi")

	var num float64
	var str string

	numValue, err := NewOutputBindValue(&num, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	strValue, err := NewOutputBindValue(&str, WithTypeVarchar(), WithBindSize(100))
	if err != nil {
		st.Fatal(err)
	}

	st.mustExec(
		`BEGIN "P_NUMBER_OUT_MULTI"(:1, :2); END;`,
		sql.Out{Dest: numValue},
		sql.Out{Dest: strValue},
	)

	if num != 42 {
		st.Fatalf("got num=%v, want 42", num)
	}
	if str != "hello" {
		st.Fatalf("got str=%q, want \"hello\"", str)
	}
}

func TestNumberOut_NumberInMiddle(t *testing.T) {
	runSqlTest(t, testNumberOutNumberInMiddle)
}

func testNumberOutNumberInMiddle(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{}

	// Procedure with NUMBER OUT in the middle
	st.mustExec(`CREATE OR REPLACE PROCEDURE p_number_out_middle(
		p_str1 OUT VARCHAR2,
		p_num OUT NUMBER,
		p_str2 OUT VARCHAR2
	) AS
	BEGIN
		p_str1 := 'before';
		p_num := 99;
		p_str2 := 'after';
	END;`)
	defer st.mustExec("DROP PROCEDURE IF EXISTS p_number_out_middle")

	var str1, str2 string
	var num float64

	str1Value, err := NewOutputBindValue(&str1, WithTypeVarchar(), WithBindSize(100))
	if err != nil {
		st.Fatal(err)
	}
	numValue, err := NewOutputBindValue(&num, WithTypeNumber())
	if err != nil {
		st.Fatal(err)
	}
	str2Value, err := NewOutputBindValue(&str2, WithTypeVarchar(), WithBindSize(100))
	if err != nil {
		st.Fatal(err)
	}

	st.mustExec(
		`BEGIN "P_NUMBER_OUT_MIDDLE"(:1, :2, :3); END;`,
		sql.Out{Dest: str1Value},
		sql.Out{Dest: numValue},
		sql.Out{Dest: str2Value},
	)

	if str1 != "before" {
		st.Fatalf("got str1=%q, want \"before\"", str1)
	}
	if num != 99 {
		st.Fatalf("got num=%v, want 99", num)
	}
	if str2 != "after" {
		st.Fatalf("got str2=%q, want \"after\"", str2)
	}
}
