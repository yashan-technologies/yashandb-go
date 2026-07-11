package yasdb

import (
	"database/sql"
	"database/sql/driver"
	"testing"
)

func TestCheckNamedValue_Number(t *testing.T) {
	stmt := &YasStmt{}

	tests := []struct {
		name    string
		value   interface{}
		wantErr error
	}{
		{"Number_valid", NewNumber("123.45"), nil},
		{"Number_null", Number{}, nil},
		{"*Number_valid", func() interface{} { n := NewNumber("99.99"); return &n }(), nil},
		{"*Number_null", func() interface{} { var n Number; return &n }(), nil},
		{"*Number_nil", (*Number)(nil), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nv := &driver.NamedValue{Value: tt.value}
			err := stmt.CheckNamedValue(nv)
			if err != tt.wantErr {
				t.Errorf("CheckNamedValue(%v) = %v, want %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestCheckNamedValue_SqlNull(t *testing.T) {
	stmt := &YasStmt{}

	tests := []struct {
		name    string
		value   interface{}
		wantErr error
	}{
		{"*sql.NullBool_valid", &sql.NullBool{Bool: true, Valid: true}, nil},
		{"*sql.NullBool_null", &sql.NullBool{Bool: false, Valid: false}, nil},
		{"*sql.NullFloat64_valid", &sql.NullFloat64{Float64: 3.14, Valid: true}, nil},
		{"*sql.NullFloat64_null", &sql.NullFloat64{Float64: 0, Valid: false}, nil},
		{"*sql.NullInt64_valid", &sql.NullInt64{Int64: 42, Valid: true}, nil},
		{"*sql.NullInt64_null", &sql.NullInt64{Int64: 0, Valid: false}, nil},
		{"*sql.NullString_valid", &sql.NullString{String: "test", Valid: true}, nil},
		{"*sql.NullString_null", &sql.NullString{String: "", Valid: false}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nv := &driver.NamedValue{Value: tt.value}
			err := stmt.CheckNamedValue(nv)
			if err != tt.wantErr {
				t.Errorf("CheckNamedValue(%v) = %v, want %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestCheckNamedValue_OtherTypes(t *testing.T) {
	stmt := &YasStmt{}

	tests := []struct {
		name    string
		value   interface{}
		wantErr error
	}{
		{"sql.Out", sql.Out{}, nil},
		{"DSInterval", DSInterval{}, nil},
		{"YMInterval", YMInterval{}, nil},
		{"Vector", Vector{}, nil},
		{"*Vector", &Vector{}, nil},
		{"[]Vector", []Vector{}, nil},
		{"[]*Vector", []*Vector{}, nil},
		{"unsupported_type", "string_value", driver.ErrSkip},
		{"unsupported_int", int64(123), driver.ErrSkip},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nv := &driver.NamedValue{Value: tt.value}
			err := stmt.CheckNamedValue(nv)
			if err != tt.wantErr {
				t.Errorf("CheckNamedValue(%v) = %v, want %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
