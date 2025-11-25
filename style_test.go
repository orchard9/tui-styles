package tuistyles

import (
	"reflect"
	"testing"
)

// TestNewStyle_ZeroValue verifies that NewStyle() returns a struct with all nil fields.
func TestNewStyle_ZeroValue(t *testing.T) {
	s := NewStyle()

	// Use reflection to verify all fields are nil pointers
	v := reflect.ValueOf(s)
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := typ.Field(i).Name

		if !field.IsNil() {
			t.Errorf("Field %s should be nil, got %v", fieldName, field.Interface())
		}
	}
}

// TestStyle_FieldCount verifies the struct has the expected number of fields.
// This test prevents accidental field removal during refactoring.
func TestStyle_FieldCount(t *testing.T) {
	s := NewStyle()
	v := reflect.ValueOf(s)

	expectedFields := 30 // 7 text attrs + 2 colors + 4 layout + 2 align + 8 spacing + 7 border (incl 2 border colors)
	actualFields := v.NumField()

	if actualFields != expectedFields {
		t.Errorf("Expected %d fields, got %d. Fields may have been added/removed.", expectedFields, actualFields)
	}
}

// TestStyle_PointerFields verifies all fields are pointers (for optionality).
func TestStyle_PointerFields(t *testing.T) {
	s := NewStyle()
	v := reflect.ValueOf(s)
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := typ.Field(i).Name
		fieldKind := field.Kind()

		if fieldKind != reflect.Ptr {
			t.Errorf("Field %s should be a pointer (for optionality), got %v", fieldName, fieldKind)
		}
	}
}

// TestStyle_ImmutableByDefault verifies that NewStyle() returns a value, not a pointer.
// This ensures the copy-on-write pattern works correctly.
func TestStyle_ImmutableByDefault(t *testing.T) {
	s1 := NewStyle()
	s2 := NewStyle()

	// Verify they are separate values
	v1 := reflect.ValueOf(s1)
	v2 := reflect.ValueOf(s2)

	if v1.Kind() != reflect.Struct || v2.Kind() != reflect.Struct {
		t.Errorf("NewStyle() should return a struct value, not a pointer")
	}
}
