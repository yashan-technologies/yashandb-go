package yasdb

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// 1. Construction & state
// ---------------------------------------------------------------------------

func TestNumberNewNumber(t *testing.T) {
	n := NewNumber("123")
	if n.String() != "123" {
		t.Errorf("String() = %q, want %q", n.String(), "123")
	}
	if !n.Valid() {
		t.Error("Valid() = false, want true")
	}
	if n.IsNull() {
		t.Error("IsNull() = true, want false")
	}
}

func TestNumberNewNumberEmpty(t *testing.T) {
	n := NewNumber("")
	if !n.Valid() {
		t.Error("Valid() = false, want true for empty string")
	}
	if n.IsNull() {
		t.Error("IsNull() = true, want false for empty string")
	}
	if n.String() != "" {
		t.Errorf("String() = %q, want empty string", n.String())
	}
	// Int64 / Float64 should fail on empty text
	if _, err := n.Int64(); err == nil {
		t.Error("Int64() on empty string should return error")
	}
	if _, err := n.Float64(); err == nil {
		t.Error("Float64() on empty string should return error")
	}
}

func TestNumberZeroValue(t *testing.T) {
	var n Number
	if n.Valid() {
		t.Error("zero value Valid() = true, want false")
	}
	if !n.IsNull() {
		t.Error("zero value IsNull() = false, want true")
	}
	if n.String() != "" {
		t.Errorf("zero value String() = %q, want empty", n.String())
	}
}

func TestNumberText(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		n := NewNumber("99.9")
		text, valid := n.Text()
		if text != "99.9" || !valid {
			t.Errorf("Text() = (%q, %v), want (%q, true)", text, valid, "99.9")
		}
	})
	t.Run("null", func(t *testing.T) {
		var n Number
		text, valid := n.Text()
		if text != "" || valid {
			t.Errorf("Text() = (%q, %v), want (%q, false)", text, valid, "")
		}
	})
}

// ---------------------------------------------------------------------------
// 2. Scan
// ---------------------------------------------------------------------------

func TestNumberScan(t *testing.T) {
	tests := []struct {
		name    string
		src     interface{}
		wantStr string
		wantErr bool
		valid   bool // expected Valid()
	}{
		{"nil", nil, "", false, false},
		{"string", "12345678901234567890", "12345678901234567890", false, true},
		{"byte_slice", []byte("123.45"), "123.45", false, true},
		{"int64", int64(123), "123", false, true},
		{"int", int(42), "42", false, true},
		{"int32", int32(99), "99", false, true},
		{"float64", float64(123.5), "123.5", false, true},
		{"float32", float32(1.5), "1.5", false, true},
		{"unsupported", struct{}{}, "", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n Number
			err := n.Scan(tt.src)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Scan(%v) error = %v, wantErr = %v", tt.src, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if n.String() != tt.wantStr {
				t.Errorf("Scan(%v) String() = %q, want %q", tt.src, n.String(), tt.wantStr)
			}
			if n.Valid() != tt.valid {
				t.Errorf("Scan(%v) Valid() = %v, want %v", tt.src, n.Valid(), tt.valid)
			}
		})
	}
}

func TestNumberScanNumberValue(t *testing.T) {
	orig := NewNumber("999.888")
	var n Number
	if err := n.Scan(orig); err != nil {
		t.Fatalf("Scan(Number) error: %v", err)
	}
	if n.String() != "999.888" || !n.Valid() {
		t.Errorf("Scan(Number) = (%q, valid=%v), want (%q, true)", n.String(), n.Valid(), "999.888")
	}
}

func TestNumberScanNumberPointer(t *testing.T) {
	t.Run("non-nil pointer", func(t *testing.T) {
		orig := NewNumber("42.0")
		var n Number
		if err := n.Scan(&orig); err != nil {
			t.Fatalf("Scan(*Number) error: %v", err)
		}
		if n.String() != "42.0" || !n.Valid() {
			t.Errorf("Scan(*Number) = (%q, valid=%v), want (%q, true)", n.String(), n.Valid(), "42.0")
		}
	})
	t.Run("nil pointer", func(t *testing.T) {
		var n Number
		if err := n.Scan((*Number)(nil)); err != nil {
			t.Fatalf("Scan(*Number(nil)) error: %v", err)
		}
		if n.Valid() || !n.IsNull() {
			t.Error("Scan(*Number(nil)) should produce NULL Number")
		}
	})
}

func TestNumberScanUnsupportedError(t *testing.T) {
	var n Number
	err := n.Scan(struct{}{})
	if err == nil {
		t.Fatal("Scan(struct{}{}) should return error")
	}
	if !strings.Contains(err.Error(), "cannot scan") {
		t.Errorf("error message = %q, want it to contain 'cannot scan'", err.Error())
	}
}

// ---------------------------------------------------------------------------
// 3. Value
// ---------------------------------------------------------------------------

func TestNumberValue(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		n := NewNumber("123.45")
		v, err := n.Value()
		if err != nil {
			t.Fatalf("Value() error: %v", err)
		}
		if v != "123.45" {
			t.Errorf("Value() = %v, want %q", v, "123.45")
		}
	})
	t.Run("null", func(t *testing.T) {
		var n Number
		v, err := n.Value()
		if err != nil {
			t.Fatalf("Value() error: %v", err)
		}
		if v != nil {
			t.Errorf("Value() = %v, want nil", v)
		}
	})
}

// ---------------------------------------------------------------------------
// 4. Int64 conversion
// ---------------------------------------------------------------------------

func TestNumberInt64(t *testing.T) {
	tests := []struct {
		name    string
		input   Number
		want    int64
		wantErr bool
		errSub  string // expected substring in error
	}{
		{"MaxInt64", NewNumber("9223372036854775807"), math.MaxInt64, false, ""},
		{"MinInt64", NewNumber("-9223372036854775808"), math.MinInt64, false, ""},
		{"zero", NewNumber("0"), 0, false, ""},
		{"negative", NewNumber("-42"), -42, false, ""},
		{"overflow", NewNumber("9223372036854775808"), 0, true, ""},
		{"decimal", NewNumber("123.4"), 0, true, ""},
		{"non_numeric", NewNumber("abc"), 0, true, ""},
		{"null", Number{}, 0, true, "NULL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Int64()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Int64() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil && tt.errSub != "" && !strings.Contains(err.Error(), tt.errSub) {
				t.Errorf("error = %q, want substring %q", err, tt.errSub)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Int64() = %d, want %d", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 5. Float64 conversion
// ---------------------------------------------------------------------------

func TestNumberFloat64(t *testing.T) {
	tests := []struct {
		name    string
		input   Number
		want    float64
		wantErr bool
		errSub  string
	}{
		{"decimal", NewNumber("123.4"), 123.4, false, ""},
		{"zero", NewNumber("0"), 0, false, ""},
		{"negative", NewNumber("-42.5"), -42.5, false, ""},
		{"integer_text", NewNumber("100"), 100.0, false, ""},
		{"null", Number{}, 0, true, "NULL"},
		{"non_numeric", NewNumber("abc"), 0, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Float64()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Float64() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil && tt.errSub != "" && !strings.Contains(err.Error(), tt.errSub) {
				t.Errorf("error = %q, want substring %q", err, tt.errSub)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 6. MarshalJSON / UnmarshalJSON
// ---------------------------------------------------------------------------

func TestNumberMarshalJSON(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		n := NewNumber("123.45")
		data, err := json.Marshal(n)
		if err != nil {
			t.Fatalf("MarshalJSON error: %v", err)
		}
		// Should be a JSON string "123.45", not a bare number
		if string(data) != `"123.45"` {
			t.Errorf("MarshalJSON = %s, want %q", data, `"123.45"`)
		}
	})
	t.Run("null", func(t *testing.T) {
		var n Number
		data, err := json.Marshal(n)
		if err != nil {
			t.Fatalf("MarshalJSON error: %v", err)
		}
		if string(data) != "null" {
			t.Errorf("MarshalJSON = %s, want null", data)
		}
	})
}

func TestNumberUnmarshalJSON(t *testing.T) {
	t.Run("null", func(t *testing.T) {
		var n Number
		if err := json.Unmarshal([]byte("null"), &n); err != nil {
			t.Fatalf("UnmarshalJSON(null) error: %v", err)
		}
		if n.Valid() || !n.IsNull() {
			t.Error("UnmarshalJSON(null) should produce NULL Number")
		}
	})
	t.Run("json_string", func(t *testing.T) {
		var n Number
		if err := json.Unmarshal([]byte(`"123.45"`), &n); err != nil {
			t.Fatalf("UnmarshalJSON(string) error: %v", err)
		}
		if n.String() != "123.45" || !n.Valid() {
			t.Errorf("UnmarshalJSON(string) = (%q, valid=%v), want (%q, true)", n.String(), n.Valid(), "123.45")
		}
	})
	t.Run("json_number", func(t *testing.T) {
		var n Number
		if err := json.Unmarshal([]byte("123"), &n); err != nil {
			t.Fatalf("UnmarshalJSON(number) error: %v", err)
		}
		if !n.Valid() {
			t.Error("UnmarshalJSON(number) should produce valid Number")
		}
		// 123 as float64 -> FormatFloat -> "123"
		if n.String() != "123" {
			t.Errorf("UnmarshalJSON(number) String() = %q, want %q", n.String(), "123")
		}
	})
	t.Run("invalid_string", func(t *testing.T) {
		var n Number
		if err := json.Unmarshal([]byte(`"invalid"`), &n); err != nil {
			t.Fatalf("UnmarshalJSON(invalid string) error: %v", err)
		}
		// It should still be a valid Number with text "invalid"
		if !n.Valid() || n.String() != "invalid" {
			t.Errorf("UnmarshalJSON(invalid) = (%q, valid=%v), want (%q, true)", n.String(), n.Valid(), "invalid")
		}
		// But conversion methods should fail
		if _, err := n.Int64(); err == nil {
			t.Error("Int64() on 'invalid' should fail")
		}
		if _, err := n.Float64(); err == nil {
			t.Error("Float64() on 'invalid' should fail")
		}
	})
}

// ---------------------------------------------------------------------------
// 7. MarshalText / UnmarshalText
// ---------------------------------------------------------------------------

func TestNumberMarshalText(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		n := NewNumber("123.45")
		data, err := n.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText error: %v", err)
		}
		if string(data) != "123.45" {
			t.Errorf("MarshalText = %q, want %q", data, "123.45")
		}
	})
	t.Run("null", func(t *testing.T) {
		var n Number
		data, err := n.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText error: %v", err)
		}
		if data != nil {
			t.Errorf("MarshalText = %v, want nil", data)
		}
	})
}

func TestNumberUnmarshalText(t *testing.T) {
	var n Number
	if err := n.UnmarshalText([]byte("123.45")); err != nil {
		t.Fatalf("UnmarshalText error: %v", err)
	}
	if !n.Valid() {
		t.Error("UnmarshalText should set Valid() = true")
	}
	if n.String() != "123.45" {
		t.Errorf("String() = %q, want %q", n.String(), "123.45")
	}
}

// ---------------------------------------------------------------------------
// 8. Edge cases
// ---------------------------------------------------------------------------

func TestNumberEdgeCases(t *testing.T) {
	t.Run("large_integer", func(t *testing.T) {
		big := "123456789012345678901234567890"
		n := NewNumber(big)
		if n.String() != big {
			t.Errorf("String() = %q, want %q", n.String(), big)
		}
		// Int64 should overflow
		if _, err := n.Int64(); err == nil {
			t.Error("Int64() on large integer should overflow")
		}
	})
	t.Run("negative_decimal", func(t *testing.T) {
		n := NewNumber("-42.5")
		if n.String() != "-42.5" {
			t.Errorf("String() = %q, want %q", n.String(), "-42.5")
		}
		f, err := n.Float64()
		if err != nil {
			t.Fatalf("Float64() error: %v", err)
		}
		if f != -42.5 {
			t.Errorf("Float64() = %v, want -42.5", f)
		}
	})
	t.Run("leading_zeros", func(t *testing.T) {
		n := NewNumber("007")
		if n.String() != "007" {
			t.Errorf("String() = %q, want %q", n.String(), "007")
		}
		v, err := n.Int64()
		if err != nil {
			t.Fatalf("Int64() error: %v", err)
		}
		if v != 7 {
			t.Errorf("Int64() = %d, want 7", v)
		}
	})
	t.Run("scientific_notation", func(t *testing.T) {
		n := NewNumber("1.23e+10")
		if n.String() != "1.23e+10" {
			t.Errorf("String() = %q, want %q", n.String(), "1.23e+10")
		}
		// Float64 should parse scientific notation
		f, err := n.Float64()
		if err != nil {
			t.Fatalf("Float64() error: %v", err)
		}
		if f != 1.23e+10 {
			t.Errorf("Float64() = %v, want 1.23e+10", f)
		}
	})
}

// ---------------------------------------------------------------------------
// JSON roundtrip integration
// ---------------------------------------------------------------------------

func TestNumberJSONRoundtrip(t *testing.T) {
	type payload struct {
		Val Number `json:"val"`
	}

	t.Run("valid roundtrip", func(t *testing.T) {
		orig := payload{Val: NewNumber("99999.99999")}
		data, err := json.Marshal(orig)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		var decoded payload
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if decoded.Val.String() != orig.Val.String() {
			t.Errorf("roundtrip String() = %q, want %q", decoded.Val.String(), orig.Val.String())
		}
		if !decoded.Val.Valid() {
			t.Error("roundtrip Valid() = false, want true")
		}
	})

	t.Run("null roundtrip", func(t *testing.T) {
		orig := payload{} // zero value -> NULL
		data, err := json.Marshal(orig)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		var decoded payload
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if decoded.Val.Valid() {
			t.Error("roundtrip Valid() = true, want false for NULL")
		}
	})
}

// ---------------------------------------------------------------------------
// 9. NUMBER CRUD (table integration)
// ---------------------------------------------------------------------------

func TestNumberCRUD(t *testing.T) {
	runSqlTestWithParams(t, "number_as_string=true", testNumberCRUD)
}

func testNumberCRUD(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{
		tableName: "test_number_crud",
		columnNameType: [][2]string{
			{"id", "INT"},
			{"n_float", "NUMBER"},
			{"n_int", "NUMBER"},
			{"n_str", "NUMBER"},
			{"n_number", "NUMBER"},
		},
	}
	st.genTableTest()
	defer st.dropTable()

	// -----------------------------------------------------------------
	// CREATE – INSERT rows using float64, int64, string
	// -----------------------------------------------------------------
	insertSQL := fmt.Sprintf(
		"INSERT INTO %s (id, n_float, n_int, n_str, n_number) VALUES (?, ?, ?, ?, ?)",
		st.tableName,
	)

	// Row 1: float64=3.14, int64=100, string="999"
	st.mustExec(insertSQL, 1, float64(3.14), int64(100), "999", NewNumber("987654321.123456789"))
	// Row 2: float64=2.718, int64=200, string="12345678901234567890"
	st.mustExec(insertSQL, 2, float64(2.718), int64(200), "12345678901234567890", NewNumber("987654321.123456789"))
	// Row 3: float64=0.0, int64=0, string="0"
	st.mustExec(insertSQL, 3, float64(0.0), int64(0), "0", NewNumber(""))

	// -----------------------------------------------------------------
	// READ – SELECT and verify values
	// -----------------------------------------------------------------
	selectSQL := fmt.Sprintf(
		"SELECT id, n_float, n_int, n_str, n_number FROM %s ORDER BY id",
		st.tableName,
	)
	rows := st.mustQuery(selectSQL)
	defer rows.Close()

	type row struct {
		id     int
		nFloat float64
		nInt   int64
		nStr   string
		nNum   Number
	}
	var got []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.id, &r.nFloat, &r.nInt, &r.nStr, &r.nNum); err != nil {
			st.Fatalf("Scan error: %v", err)
		}
		got = append(got, r)
	}
	if err := rows.Err(); err != nil {
		st.Fatalf("Rows error: %v", err)
	}

	want := []row{
		{1, 3.14, 100, "999", NewNumber("987654321.123456789")},
		{2, 2.718, 200, "12345678901234567890", NewNumber("987654321.123456789")},
		{3, 0.0, 0, "0", NewNumber("")},
	}
	if len(got) != len(want) {
		st.Fatalf("row count = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].id != want[i].id {
			st.Errorf("row %d: id = %d, want %d", i, got[i].id, want[i].id)
		}
		if got[i].nFloat != want[i].nFloat {
			st.Errorf("row %d: n_float = %v, want %v", i, got[i].nFloat, want[i].nFloat)
		}
		if got[i].nInt != want[i].nInt {
			st.Errorf("row %d: n_int = %d, want %d", i, got[i].nInt, want[i].nInt)
		}
		if got[i].nStr != want[i].nStr {
			st.Errorf("row %d: n_str = %q, want %q", i, got[i].nStr, want[i].nStr)
		}
		if got[i].nNum.String() != want[i].nNum.String() {
			st.Errorf("row %d: n_number = %q, want %q", i, got[i].nStr, want[i].nStr)
		}
	}

	// -----------------------------------------------------------------
	// UPDATE – modify row id=1: n_float→9.99, n_int→999, n_str→"3.14159"
	// -----------------------------------------------------------------
	updateSQL := fmt.Sprintf(
		"UPDATE %s SET n_float = ?, n_int = ?, n_str = ?, n_number = ? WHERE id = ?",
		st.tableName,
	)
	res := st.mustExec(updateSQL, float64(9.99), int64(999), "3.14159", NewNumber("3.14159"), 1)
	if n, _ := res.RowsAffected(); n != 1 {
		st.Fatalf("UPDATE rows affected = %d, want 1", n)
	}

	// Verify update
	var updFloat float64
	var updInt int64
	var updStr string
	var upNum Number
	err := st.QueryRow(
		fmt.Sprintf("SELECT n_float, n_int, n_str, n_number FROM %s WHERE id = ?", st.tableName),
		1,
	).Scan(&updFloat, &updInt, &updStr, &upNum)
	if err != nil {
		st.Fatalf("QueryRow after UPDATE error: %v", err)
	}
	if updFloat != 9.99 {
		st.Errorf("after UPDATE: n_float = %v, want 9.99", updFloat)
	}
	if updInt != 999 {
		st.Errorf("after UPDATE: n_int = %d, want 999", updInt)
	}
	if updStr != "3.14159" {
		st.Errorf("after UPDATE: n_str = %q, want %q", updStr, "3.14159")
	}
	if upNum.String() != "3.14159" {
		st.Errorf("after UPDATE: n_str = %q, want %q", updStr, "3.14159")
	}
	// -----------------------------------------------------------------
	// DELETE – remove row id=3, verify remaining rows
	// -----------------------------------------------------------------
	deleteSQL := fmt.Sprintf("DELETE FROM %s WHERE id = ?", st.tableName)
	res = st.mustExec(deleteSQL, 3)
	if n, _ := res.RowsAffected(); n != 1 {
		st.Fatalf("DELETE rows affected = %d, want 1", n)
	}

	// Verify deletion – should have 2 rows left
	var count int
	err = st.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", st.tableName)).Scan(&count)
	if err != nil {
		st.Fatalf("COUNT after DELETE error: %v", err)
	}
	if count != 2 {
		st.Errorf("after DELETE: count = %d, want 2", count)
	}

	// Verify id=3 is gone
	var exist int
	err = st.QueryRow(
		fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = ?", st.tableName),
		3,
	).Scan(&exist)
	if err != nil {
		st.Fatalf("check deleted row error: %v", err)
	}
	if exist != 0 {
		st.Errorf("id=3 still exists after DELETE")
	}
}

// ---------------------------------------------------------------------------
// 9.5 INSERT NULL Number
// ---------------------------------------------------------------------------

func TestNumberInsertNull(t *testing.T) {
	runSqlTest(t, testNumberInsertNull)
}

func testNumberInsertNull(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{
		tableName: "test_number_insert_null",
		columnNameType: [][2]string{
			{"id", "INT"},
			{"n_val", "NUMBER"},
		},
	}
	st.genTableTest()
	defer st.dropTable()

	insertSQL := fmt.Sprintf("INSERT INTO %s (id, n_val) VALUES (?, ?)", st.tableName)

	// Row 1: NULL Number (zero value)
	st.mustExec(insertSQL, 1, Number{})

	// Row 2: *Number pointing to NULL Number
	nullNum := Number{}
	st.mustExec(insertSQL, 2, &nullNum)

	// Row 3: nil *Number
	var nilNum *Number
	st.mustExec(insertSQL, 3, nilNum)

	// Row 4: valid Number for comparison
	st.mustExec(insertSQL, 4, NewNumber("123.45"))

	// Verify all values
	selectSQL := fmt.Sprintf("SELECT id, n_val FROM %s ORDER BY id", st.tableName)
	rows := st.mustQuery(selectSQL)
	defer rows.Close()

	type row struct {
		id    int
		nVal  Number
		valid bool
	}
	var got []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.id, &r.nVal); err != nil {
			st.Fatalf("Scan error: %v", err)
		}
		r.valid = r.nVal.Valid()
		got = append(got, r)
	}
	if err := rows.Err(); err != nil {
		st.Fatalf("Rows error: %v", err)
	}

	if len(got) != 4 {
		st.Fatalf("row count = %d, want 4", len(got))
	}

	// Rows 1, 2, 3 should be NULL
	for i := 0; i < 3; i++ {
		if got[i].valid {
			st.Errorf("row %d (id=%d): expected NULL Number, got Valid()=true, String()=%q",
				i, got[i].id, got[i].nVal.String())
		}
	}

	// Row 4 should be valid
	if !got[3].valid {
		st.Errorf("row 3 (id=%d): expected valid Number, got NULL", got[3].id)
	}
	if got[3].nVal.String() != "123.45" {
		st.Errorf("row 3 (id=%d): n_val = %q, want %q", got[3].id, got[3].nVal.String(), "123.45")
	}
}

func TestNumberInsertEdgeCases(t *testing.T) {
	runSqlTest(t, testNumberInsertEdgeCases)
}

func testNumberInsertEdgeCases(st *sqlTest) {
	st.sqlGenInfo = &sqlGenInfo{
		tableName: "test_number_edge_cases",
		columnNameType: [][2]string{
			{"id", "INT"},
			{"n_val", "NUMBER"},
		},
	}
	st.genTableTest()
	defer st.dropTable()

	insertSQL := fmt.Sprintf("INSERT INTO %s (id, n_val) VALUES (?, ?)", st.tableName)

	tests := []struct {
		name   string
		id     int
		number Number
		want   string // expected string representation after DB roundtrip
	}{
		{"zero", 1, NewNumber("0"), "0"},
		{"negative_integer", 2, NewNumber("-42"), "-42"},
		{"negative_decimal", 3, NewNumber("-3.14159"), "-3.14159"},
		{"small_decimal", 4, NewNumber("0.000000000001"), "0.000000000001"},
		{"reasonable_decimals", 5, NewNumber("3.14159265358979"), "3.14159265358979"},
		{"simple_integer", 6, NewNumber("123456789"), "123456789"},
		{"moderate_large", 7, NewNumber("999999999999999"), "999999999999999"},
		{"moderate_negative", 8, NewNumber("-999999999999999"), "-999999999999999"},
	}

	for _, tt := range tests {
		st.mustExec(insertSQL, tt.id, tt.number)
	}

	// Verify all values
	selectSQL := fmt.Sprintf("SELECT id, n_val FROM %s ORDER BY id", st.tableName)
	rows := st.mustQuery(selectSQL)
	defer rows.Close()

	var got []struct {
		id   int
		nVal Number
	}
	for rows.Next() {
		var r struct {
			id   int
			nVal Number
		}
		if err := rows.Scan(&r.id, &r.nVal); err != nil {
			st.Fatalf("Scan error: %v", err)
		}
		got = append(got, r)
	}
	if err := rows.Err(); err != nil {
		st.Fatalf("Rows error: %v", err)
	}

	if len(got) != len(tests) {
		st.Fatalf("row count = %d, want %d", len(got), len(tests))
	}

	// Verify each value matches expected
	for i, tt := range tests {
		if got[i].id != tt.id {
			st.Errorf("row %d: id = %d, want %d", i, got[i].id, tt.id)
		}
		gotStr := got[i].nVal.String()
		if gotStr != tt.want {
			st.Errorf("row %d (%s): n_val = %q, want %q", i, tt.name, gotStr, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// 10. DatabaseTypeName
// ---------------------------------------------------------------------------

func TestDatabaseTypeName(t *testing.T) {
	runSqlTest(t, testNumberFloatDatabaseTypeName)
}

func testNumberFloatDatabaseTypeName(t *sqlTest) {
	si := sqlGenInfo{}
	t.sqlGenInfo = &si

	columnNameType := [][2]string{
		{"c1", "TINYINT"},
		{"c2", "SMALLINT"},
		{"c3", "INT"},
		{"c4", "BIGINT"},
		{"c5", "FLOAT"},
		{"c6", "DOUBLE"},
		{"c7", "NUMBER"},
		{"c8", "BIT(64)"},
		{"c9", "CHAR(126)"},
		{"c10", "VARCHAR(126)"},
		{"c11", "NCHAR(126)"},
		{"c12", "NVARCHAR(126)"},
		{"c13", "BOOLEAN"},
		{"c14", "TIME"},
		{"c15", "TIMESTAMP"},
		{"c18", "INTERVAL YEAR TO MONTH"},
		{"c19", "INTERVAL DAY TO SECOND"},
		{"c20", "BLOB"},
		{"c21", "CLOB"},
		{"c22", "NCLOB"},
		{"c23", "RAW(10)"},
		{"c24", "JSON"},
		{"c26", "ROWID"},
		{"c27", "UROWID(20)"},
	}

	columnTypes := []string{
		"TINYINT",
		"SMALLINT",
		"INTEGER",
		"BIGINT",
		"FLOAT",
		"DOUBLE",
		"NUMBER",
		"BIT",
		"CHAR",
		"VARCHAR",
		"NCHAR",
		"NVARCHAR",
		"BOOLEAN",
		"TIME",
		"TIMESTAMP",
		"INTERVAL YEAR TO MONTH",
		"INTERVAL DAY TO SECOND",
		"BLOB",
		"CLOB",
		"NCLOB",
		"RAW",
		"JSON",
		"ROWID",
		"RAW",
	}

	if t.isBfileSupport() {
		columnNameType = append(columnNameType, [2]string{"c25", "BFILE"})
		columnTypes = append(columnTypes, "BFILE")
	}
	if t.isToTimestampTzSupport() {
		columnNameType = append(columnNameType, [][2]string{
			{"c16", "TIMESTAMP WITH LOCAL TIME ZONE"},
			{"c17", "TIMESTAMP WITH TIME ZONE"}}...)
		columnTypes = append(columnTypes,
			"TIMESTAMP WITH LOCAL TIME ZONE",
			"TIMESTAMP WITH TIME ZONE")
	}

	si = sqlGenInfo{
		tableName:      "database_type_name",
		columnNameType: columnNameType,
	}
	t.genTableTest()
	t.runSelectTest()

	rows, err := t.getRowsColumnTypes()
	if err != nil {
		t.Fatalf(err.Error())
	}
	for i, row := range rows {
		if row.DatabaseTypeName() != columnTypes[i] {
			t.Fatalf("column %d database type expected: %s actual: %s", i, columnTypes[i], row.DatabaseTypeName())
		}
	}
}
