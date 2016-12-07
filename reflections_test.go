// Copyright (c) 2013 Théo Crevon
//
// See the file LICENSE for copying permission.

package reflectme

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestDummyOnlyStruct struct {
	unexported uint64
	Dummy      string `test:"dummytag"`
}

type TestStruct struct {
	unexported uint64
	Dummy      string `test:"dummytag"`
	Yummy      int    `test:"yummytag"`
}

type TestNestedPointerStruct struct {
	unexported uint64
	Dummy      string `test:"dummytag"`
	Yummy      int    `test:"yummytag"`
	Nested     *NestedStruct
}

type TestNestedStruct struct {
	unexported uint64
	Dummy      string `test:"dummytag"`
	Yummy      int    `test:"yummytag"`
	Nested     NestedStruct
}

type TestInnerStruct struct {
	unexported uint64
	Dummy      string `test:"dummytag"`
	Yummy      int    `test:"yummytag"`
	Nested     TestNestedStruct
}

type NestedStruct struct {
	Dummy string `test:"dummytag"`
	Yummy int    `test:"yummytag"`
}

func TestGetField_on_struct(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
	}
	value, err := GetField(dummyStruct, "Dummy")
	assert.NoError(t, err)
	assert.Equal(t, value, "test")
}

func TestGetField_on_nested_struct(t *testing.T) {
	dummyStruct := TestNestedStruct{
		Dummy: "test",
		Nested: NestedStruct{
			Dummy: "nested",
		},
	}
	value, err := GetField(dummyStruct, "Nested.Dummy")
	assert.NoError(t, err)
	assert.Equal(t, value, "nested")
}

func TestGetField_on_nested_pointer_struct(t *testing.T) {
	dummyStruct := TestNestedPointerStruct{
		Dummy: "test",
		Nested: &NestedStruct{
			Dummy: "nested",
		},
	}
	value, err := GetField(dummyStruct, "Nested.Dummy")
	assert.NoError(t, err)
	assert.Equal(t, value, "nested")
}

func TestGetField_on_inner_pointer_struct(t *testing.T) {
	dummyStruct := TestInnerStruct{
		Nested: TestNestedStruct{
			Dummy: "test",
			Yummy: 123,
		},
	}
	dummyStruct.Nested.Nested.Dummy = "dummy"
	value, err := GetField(dummyStruct, "Nested.Nested.Dummy")
	assert.NoError(t, err)
	assert.Equal(t, value, "dummy")
}

func TestGetField_on_struct_pointer(t *testing.T) {
	dummyStruct := &TestStruct{
		Dummy: "test",
	}

	value, err := GetField(dummyStruct, "Dummy")
	assert.NoError(t, err)
	assert.Equal(t, value, "test")
}

func TestGetField_on_nested_struct_pointer(t *testing.T) {
	dummyStruct := &TestNestedStruct{
		Dummy: "test",
		Nested: NestedStruct{
			Dummy: "nested",
		},
	}
	value, err := GetField(dummyStruct, "Nested.Dummy")
	assert.NoError(t, err)
	assert.Equal(t, value, "nested")
}

func TestGetField_on_nested_pointer_struct_pointer(t *testing.T) {
	dummyStruct := &TestNestedPointerStruct{
		Dummy: "test",
		Nested: &NestedStruct{
			Dummy: "nested",
		},
	}
	value, err := GetField(dummyStruct, "Nested.Dummy")
	assert.NoError(t, err)
	assert.Equal(t, value, "nested")
}

func TestGetField_on_non_struct(t *testing.T) {
	dummy := "abc 123"

	_, err := GetField(dummy, "Dummy")
	assert.Error(t, err)
}

func TestGetField_non_existing_field(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
	}

	_, err := GetField(dummyStruct, "obladioblada")
	assert.Error(t, err)
}

func TestGetField_non_existing_nested_field(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
	}

	_, err := GetField(dummyStruct, "Dummy1.Bla")
	assert.Error(t, err)
}

func TestGetField_on_non_struct_nested_field(t *testing.T) {
	dummyStruct := TestNestedPointerStruct{
		Dummy: "test",
	}

	_, err := GetField(dummyStruct, "Nested.Bla")
	assert.Error(t, err)
}

func TestGetField_unexported_field(t *testing.T) {
	dummyStruct := TestStruct{
		unexported: 12345,
		Dummy:      "test",
	}

	assert.Panics(t, func() {
		GetField(dummyStruct, "unexported")
	})
}

func TestGetFieldKind_on_struct(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	kind, err := GetFieldKind(dummyStruct, "Dummy")
	assert.NoError(t, err)
	assert.Equal(t, kind, reflect.String)

	kind, err = GetFieldKind(dummyStruct, "Yummy")
	assert.NoError(t, err)
	assert.Equal(t, kind, reflect.Int)
}

func TestGetFieldKind_on_nested_struct(t *testing.T) {
	dummyStruct := TestNestedStruct{
		Nested: NestedStruct{
			Dummy: "test",
			Yummy: 123,
		},
	}

	kind, err := GetFieldKind(dummyStruct, "Nested.Dummy")
	assert.NoError(t, err)
	assert.Equal(t, kind, reflect.String)

	kind, err = GetFieldKind(dummyStruct, "Nested.Yummy")
	assert.NoError(t, err)
	assert.Equal(t, kind, reflect.Int)
}

func TestGetFieldKind_on_nested_pointer_struct(t *testing.T) {
	dummyStruct := TestNestedPointerStruct{
		Nested: &NestedStruct{
			Dummy: "test",
			Yummy: 123,
		},
	}

	kind, err := GetFieldKind(dummyStruct, "Nested.Dummy")
	assert.NoError(t, err)
	assert.Equal(t, kind, reflect.String)

	kind, err = GetFieldKind(dummyStruct, "Nested.Yummy")
	assert.NoError(t, err)
	assert.Equal(t, kind, reflect.Int)
}

func TestGetFieldKind_on_struct_pointer(t *testing.T) {
	dummyStruct := &TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	kind, err := GetFieldKind(dummyStruct, "Dummy")
	assert.NoError(t, err)
	assert.Equal(t, kind, reflect.String)

	kind, err = GetFieldKind(dummyStruct, "Yummy")
	assert.NoError(t, err)
	assert.Equal(t, kind, reflect.Int)
}

func TestGetFieldKind_on_nested_struct_pointer(t *testing.T) {
	dummyStruct := &TestNestedStruct{
		Nested: NestedStruct{
			Dummy: "test",
			Yummy: 123,
		},
	}

	kind, err := GetFieldKind(dummyStruct, "Nested.Dummy")
	assert.NoError(t, err)
	assert.Equal(t, kind, reflect.String)

	kind, err = GetFieldKind(dummyStruct, "Nested.Yummy")
	assert.NoError(t, err)
	assert.Equal(t, kind, reflect.Int)
}

func TestGetFieldKind_on_nested_pointer_struct_pointer(t *testing.T) {
	dummyStruct := &TestNestedPointerStruct{
		Nested: &NestedStruct{
			Dummy: "test",
			Yummy: 123,
		},
	}

	kind, err := GetFieldKind(dummyStruct, "Nested.Dummy")
	assert.NoError(t, err)
	assert.Equal(t, kind, reflect.String)

	kind, err = GetFieldKind(dummyStruct, "Nested.Yummy")
	assert.NoError(t, err)
	assert.Equal(t, kind, reflect.Int)
}

func TestGetFieldKind_on_non_struct(t *testing.T) {
	dummy := "abc 123"

	_, err := GetFieldKind(dummy, "Dummy")
	assert.Error(t, err)
}

func TestGetFieldKind_non_existing_field(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	_, err := GetFieldKind(dummyStruct, "obladioblada")
	assert.Error(t, err)
}

func TestGetFieldTag_on_struct(t *testing.T) {
	dummyStruct := TestStruct{}

	tag, err := GetFieldTag(dummyStruct, "Dummy", "test")
	assert.NoError(t, err)
	assert.Equal(t, tag, "dummytag")

	tag, err = GetFieldTag(dummyStruct, "Yummy", "test")
	assert.NoError(t, err)
	assert.Equal(t, tag, "yummytag")
}

func TestGetFieldTag_on_nested_struct(t *testing.T) {
	dummyStruct := TestNestedStruct{}

	tag, err := GetFieldTag(dummyStruct, "Nested.Dummy", "test")
	assert.NoError(t, err)
	assert.Equal(t, tag, "dummytag")

	tag, err = GetFieldTag(dummyStruct, "Nested.Yummy", "test")
	assert.NoError(t, err)
	assert.Equal(t, tag, "yummytag")
}

func TestGetFieldTag_on_nested_pointer_struct(t *testing.T) {
	dummyStruct := TestNestedPointerStruct{
		Nested: &NestedStruct{},
	}

	tag, err := GetFieldTag(dummyStruct, "Nested.Dummy", "test")
	assert.NoError(t, err)
	assert.Equal(t, tag, "dummytag")

	tag, err = GetFieldTag(dummyStruct, "Nested.Yummy", "test")
	assert.NoError(t, err)
	assert.Equal(t, tag, "yummytag")
}

func TestGetFieldTag_on_nested_nil_pointer_struct(t *testing.T) {
	dummyStruct := TestNestedPointerStruct{}

	_, err := GetFieldTag(dummyStruct, "Nested.Dummy", "test")
	assert.Error(t, err)

	_, err = GetFieldTag(dummyStruct, "Nested.Yummy", "test")
	assert.Error(t, err)
}

func TestGetFieldTag_on_struct_pointer(t *testing.T) {
	dummyStruct := &TestStruct{}

	tag, err := GetFieldTag(dummyStruct, "Dummy", "test")
	assert.NoError(t, err)
	assert.Equal(t, tag, "dummytag")

	tag, err = GetFieldTag(dummyStruct, "Yummy", "test")
	assert.NoError(t, err)
	assert.Equal(t, tag, "yummytag")
}

func TestGetFieldTag_on_nested_struct_pointer(t *testing.T) {
	dummyStruct := &TestNestedStruct{}

	tag, err := GetFieldTag(dummyStruct, "Nested.Dummy", "test")
	assert.NoError(t, err)
	assert.Equal(t, tag, "dummytag")

	tag, err = GetFieldTag(dummyStruct, "Nested.Yummy", "test")
	assert.NoError(t, err)
	assert.Equal(t, tag, "yummytag")
}

func TestGetFieldTag_on_nested_pointer_struct_pointer(t *testing.T) {
	dummyStruct := &TestNestedPointerStruct{
		Nested: &NestedStruct{},
	}

	tag, err := GetFieldTag(dummyStruct, "Nested.Dummy", "test")
	assert.NoError(t, err)
	assert.Equal(t, tag, "dummytag")

	tag, err = GetFieldTag(dummyStruct, "Nested.Yummy", "test")
	assert.NoError(t, err)
	assert.Equal(t, tag, "yummytag")
}

func TestGetFieldTag_on_nested_nil_pointer_struct_pointer(t *testing.T) {
	dummyStruct := &TestNestedPointerStruct{}

	_, err := GetFieldTag(dummyStruct, "Nested.Dummy", "test")
	assert.Error(t, err)

	_, err = GetFieldTag(dummyStruct, "Nested.Yummy", "test")
	assert.Error(t, err)
}

func TestGetFieldTag_on_non_struct(t *testing.T) {
	dummy := "abc 123"

	_, err := GetFieldTag(dummy, "Dummy", "test")
	assert.Error(t, err)
}

func TestGetFieldTag_non_existing_field(t *testing.T) {
	dummyStruct := TestStruct{}

	_, err := GetFieldTag(dummyStruct, "obladioblada", "test")
	assert.Error(t, err)
}

func TestGetFieldTag_non_existing_nested_field(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
	}

	_, err := GetFieldTag(dummyStruct, "Dummy1.Bla", "test")
	assert.Error(t, err)
}

func TestGetFieldTag_unexported_field(t *testing.T) {
	dummyStruct := TestStruct{
		unexported: 12345,
		Dummy:      "test",
	}

	_, err := GetFieldTag(dummyStruct, "unexported", "test")
	assert.Error(t, err)
}

func TestSetField_on_struct_with_valid_value_type(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
	}

	err := SetField(&dummyStruct, "Dummy", "abc")
	assert.NoError(t, err)
	assert.Equal(t, dummyStruct.Dummy, "abc")
}

func TestSetField_on_nested_struct_pointer_with_valid_value_type(t *testing.T) {
	dummyStruct := TestNestedPointerStruct{
		Dummy: "test",
		Nested: &NestedStruct{
			Dummy: "nested",
		},
	}
	err := SetField(&dummyStruct, "Nested.Dummy", "abc")
	assert.NoError(t, err)
	assert.Equal(t, "abc", dummyStruct.Nested.Dummy)
}

func TestSetField_on_nested_struct_with_valid_value_type(t *testing.T) {
	dummyStruct := TestNestedStruct{
		Dummy: "test",
		Nested: NestedStruct{
			Dummy: "nested",
		},
	}
	err := SetField(&dummyStruct, "Nested.Dummy", "abc")
	assert.NoError(t, err)
	assert.Equal(t, dummyStruct.Nested.Dummy, "abc")
}

func TestSetField_on_inner_struct_with_valid_value_type(t *testing.T) {
	dummyStruct := TestInnerStruct{
		Dummy: "test",
	}
	dummyStruct.Nested.Nested.Dummy = "cba"
	err := SetField(&dummyStruct, "Nested.Nested.Dummy", "abc")
	assert.NoError(t, err)
	assert.Equal(t, dummyStruct.Nested.Nested.Dummy, "abc")
}

func TestSetField_non_existing_field(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
	}

	err := SetField(&dummyStruct, "obladioblada", "life goes on")
	assert.Error(t, err)
}

func TestSetField_invalid_value_type(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
	}

	err := SetField(&dummyStruct, "Yummy", "123")
	assert.Error(t, err)
}

func TestSetField_non_exported_field(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
	}

	assert.Error(t, SetField(&dummyStruct, "unexported", "fail, bitch"))
}

func TestSetField_non_pointer(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
	}

	assert.Error(t, SetField(dummyStruct, "Dummy", "abc"))
}

func TestCopyField_from_structs(t *testing.T) {
	fromStruct := TestStruct{
		Dummy: "test",
	}

	toStruct := TestStruct{}

	err := CopyField(fromStruct, &toStruct, "Dummy")
	assert.NoError(t, err)
	assert.Equal(t, fromStruct.Dummy, toStruct.Dummy)
}

func TestCopyField_with_error(t *testing.T) {
	fromStruct := "test"

	toStruct := TestStruct{}

	err := CopyField(fromStruct, &toStruct, "unexported")
	assert.Error(t, err)
}

func TestFieldsNames_on_struct(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	fields, err := FieldsNames(dummyStruct)
	assert.NoError(t, err)
	assert.Equal(t, fields, []string{"Dummy", "Yummy"})
}

func TestFieldsNames_on_nested_struct(t *testing.T) {
	dummyStruct := TestNestedStruct{
		Dummy: "test",
		Yummy: 123,
	}

	fields, err := FieldsNames(dummyStruct)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Dummy", "Yummy", "Nested", "Nested.Dummy", "Nested.Yummy"}, fields)
}

func TestFieldsNames_on_nested_pointer_struct(t *testing.T) {
	dummyStruct := TestNestedPointerStruct{
		Dummy:  "test",
		Yummy:  123,
		Nested: &NestedStruct{},
	}

	fields, err := FieldsNames(dummyStruct)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Dummy", "Yummy", "Nested", "Nested.Dummy", "Nested.Yummy"}, fields)
}

func TestFieldsNames_on_struct_pointer(t *testing.T) {
	dummyStruct := &TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	fields, err := FieldsNames(dummyStruct)
	assert.NoError(t, err)
	assert.Equal(t, fields, []string{"Dummy", "Yummy"})
}

func TestFieldsNames_on_nested_struct_pointer(t *testing.T) {
	dummyStruct := &TestNestedStruct{
		Dummy: "test",
		Yummy: 123,
	}

	fields, err := FieldsNames(dummyStruct)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Dummy", "Yummy", "Nested", "Nested.Dummy", "Nested.Yummy"}, fields)
}

func TestFieldsNames_on_nested_pointer_struct_pointer(t *testing.T) {
	dummyStruct := &TestNestedPointerStruct{
		Dummy:  "test",
		Yummy:  123,
		Nested: &NestedStruct{},
	}

	fields, err := FieldsNames(dummyStruct)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Dummy", "Yummy", "Nested", "Nested.Dummy", "Nested.Yummy"}, fields)
}

func TestFieldsNames_on_non_struct(t *testing.T) {
	dummy := "abc 123"

	_, err := FieldsNames(dummy)
	assert.Error(t, err)
}

func TestFieldsNames_with_non_exported_fields(t *testing.T) {
	dummyStruct := TestStruct{
		unexported: 6789,
		Dummy:      "test",
		Yummy:      123,
	}

	fields, err := FieldsNames(dummyStruct)
	assert.NoError(t, err)
	assert.Equal(t, fields, []string{"Dummy", "Yummy"})
}

func TestFields_on_struct(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	fields, err := Fields(dummyStruct)
	assert.NoError(t, err)
	expFields := []string{"Dummy", "Yummy"}
	for i, field := range fields {
		assert.Equal(t, expFields[i], field.Name)
	}
}

func TestFields_on_struct_pointer(t *testing.T) {
	dummyStruct := &TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	fields, err := Fields(dummyStruct)
	assert.NoError(t, err)
	expFields := []string{"Dummy", "Yummy"}
	assert.Equal(t, len(expFields), len(fields))
	for i, field := range fields {
		assert.Equal(t, expFields[i], field.Name)
	}
}

func TestFields_on_non_struct(t *testing.T) {
	dummy := "abc 123"

	_, err := Fields(dummy)
	assert.Error(t, err)
}

func TestFields_with_non_exported_fields(t *testing.T) {
	dummyStruct := TestStruct{
		unexported: 6789,
		Dummy:      "test",
		Yummy:      123,
	}

	fields, err := Fields(dummyStruct)
	assert.NoError(t, err)
	expFields := []string{"Dummy", "Yummy"}
	assert.Equal(t, len(expFields), len(fields))
	for i, field := range fields {
		assert.Equal(t, expFields[i], field.Name)
	}
}

func TestFields_on_nested_struct(t *testing.T) {
	dummyStruct := TestNestedStruct{
		Dummy: "test",
		Yummy: 123,
	}

	fields, err := Fields(dummyStruct)
	assert.NoError(t, err)
	expFields := []string{"Dummy", "Yummy", "Nested", "Dummy", "Yummy"}
	assert.Equal(t, len(expFields), len(fields))
	for i, field := range fields {
		assert.Equal(t, expFields[i], field.Name)
	}
}

func TestFields_on_nested_pointer_struct(t *testing.T) {
	dummyStruct := TestNestedPointerStruct{
		Dummy:  "test",
		Yummy:  123,
		Nested: &NestedStruct{},
	}

	fields, err := Fields(dummyStruct)
	assert.NoError(t, err)
	expFields := []string{"Dummy", "Yummy", "Nested", "Dummy", "Yummy"}
	assert.Equal(t, len(expFields), len(fields))
	for i, field := range fields {
		assert.Equal(t, expFields[i], field.Name)
	}
}

func TestFields_on_nested_struct_pointer(t *testing.T) {
	dummyStruct := &TestNestedStruct{
		Dummy: "test",
		Yummy: 123,
	}

	fields, err := Fields(dummyStruct)
	assert.NoError(t, err)
	expFields := []string{"Dummy", "Yummy", "Nested", "Dummy", "Yummy"}
	assert.Equal(t, len(expFields), len(fields))
	for i, field := range fields {
		assert.Equal(t, expFields[i], field.Name)
	}
}

func TestFields_on_nested_pointer_struct_pointer(t *testing.T) {
	dummyStruct := &TestNestedPointerStruct{
		Dummy:  "test",
		Yummy:  123,
		Nested: &NestedStruct{},
	}

	fields, err := Fields(dummyStruct)
	assert.NoError(t, err)
	expFields := []string{"Dummy", "Yummy", "Nested", "Dummy", "Yummy"}
	assert.Equal(t, len(expFields), len(fields))
	for i, field := range fields {
		assert.Equal(t, expFields[i], field.Name)
	}
}

func TestHasField_on_struct_with_existing_field(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	has, err := HasField(dummyStruct, "Dummy")
	assert.NoError(t, err)
	assert.True(t, has)
}

func TestHasField_on_struct_pointer_with_existing_field(t *testing.T) {
	dummyStruct := &TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	has, err := HasField(dummyStruct, "Dummy")
	assert.NoError(t, err)
	assert.True(t, has)
}

func TestHasField_non_existing_field(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	has, err := HasField(dummyStruct, "Test")
	assert.NoError(t, err)
	assert.False(t, has)
}

func TestHasField_on_non_struct(t *testing.T) {
	dummy := "abc 123"

	_, err := HasField(dummy, "Test")
	assert.Error(t, err)
}

func TestHasField_unexported_field(t *testing.T) {
	dummyStruct := TestStruct{
		unexported: 7890,
		Dummy:      "test",
		Yummy:      123,
	}

	has, err := HasField(dummyStruct, "unexported")
	assert.NoError(t, err)
	assert.False(t, has)
}

func TestTags_on_struct(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	tags, err := Tags(dummyStruct, "test")
	assert.NoError(t, err)
	assert.Equal(t, tags, map[string]string{
		"Dummy": "dummytag",
		"Yummy": "yummytag",
	})
}

func TestTags_on_struct_pointer(t *testing.T) {
	dummyStruct := &TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	tags, err := Tags(dummyStruct, "test")
	assert.NoError(t, err)
	assert.Equal(t, tags, map[string]string{
		"Dummy": "dummytag",
		"Yummy": "yummytag",
	})
}

func TestTags_on_non_struct(t *testing.T) {
	dummy := "abc 123"

	_, err := Tags(dummy, "test")
	assert.Error(t, err)
}

func TestItems_on_struct(t *testing.T) {
	dummyStruct := TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	tags, err := Items(dummyStruct)
	assert.NoError(t, err)
	assert.Equal(t, tags, map[string]interface{}{
		"Dummy": "test",
		"Yummy": 123,
	})
}

func TestItems_on_struct_pointer(t *testing.T) {
	dummyStruct := &TestStruct{
		Dummy: "test",
		Yummy: 123,
	}

	tags, err := Items(dummyStruct)
	assert.NoError(t, err)
	assert.Equal(t, tags, map[string]interface{}{
		"Dummy": "test",
		"Yummy": 123,
	})
}

func TestItems_on_non_struct(t *testing.T) {
	dummy := "abc 123"

	_, err := Items(dummy)
	assert.Error(t, err)
}

func TestCopy_with_default_values(t *testing.T) {
	dummyStruct1 := TestStruct{
		Dummy: "test",
		Yummy: 0,
	}
	dummyStruct2 := TestStruct{
		Dummy: "test1",
		Yummy: 1,
	}
	err := Copy(dummyStruct1, &dummyStruct2)
	assert.NoError(t, err)
	assert.Equal(t, dummyStruct1.Dummy, dummyStruct2.Dummy)
	assert.Equal(t, dummyStruct1.Yummy, dummyStruct2.Yummy)
}

func TestCopy_not_copying_zero_values(t *testing.T) {
	DefaultCopyOptions.CopyZeroValues = false
	dummyStruct1 := TestStruct{
		Dummy: "test",
		Yummy: 0,
	}
	dummyStruct2 := TestStruct{
		Dummy: "test1",
		Yummy: 1,
	}
	err := Copy(dummyStruct1, &dummyStruct2)
	assert.NoError(t, err)
	assert.Equal(t, dummyStruct1.Dummy, dummyStruct2.Dummy)
	assert.Equal(t, 1, dummyStruct2.Yummy)
}

func TestCopy_ignoring_not_founded_fields(t *testing.T) {
	DefaultCopyOptions.IgnoreNotFoundedFields = true
	dummyStruct1 := TestStruct{
		Dummy: "test",
	}

	dummyStruct2 := TestDummyOnlyStruct{
		Dummy: "test1",
	}

	err := Copy(dummyStruct1, &dummyStruct2)
	assert.NoError(t, err)
	assert.Equal(t, dummyStruct1.Dummy, dummyStruct2.Dummy)
}

func TestCopy_not_ignoring_not_founded_fields(t *testing.T) {
	DefaultCopyOptions.IgnoreNotFoundedFields = false
	dummyStruct1 := TestStruct{
		Dummy: "test",
		Yummy: 1,
	}

	dummyStruct2 := TestDummyOnlyStruct{
		Dummy: "test1",
	}

	err := Copy(dummyStruct1, &dummyStruct2)
	assert.Error(t, err)
}

func TestCopy_on_nested_struct(t *testing.T) {
	DefaultCopyOptions.IgnoreNotFoundedFields = true
	dummyStruct1 := TestNestedStruct{
		Dummy: "test",
		Nested: NestedStruct{
			Dummy: "nested",
		},
	}

	dummyStruct2 := TestNestedStruct{
		Dummy: "test1",
		Nested: NestedStruct{
			Dummy: "nested1",
		},
	}

	err := Copy(dummyStruct1, &dummyStruct2)
	assert.NoError(t, err)
	assert.Equal(t, dummyStruct1.Dummy, dummyStruct2.Dummy)
	assert.Equal(t, dummyStruct1.Nested.Dummy, dummyStruct2.Nested.Dummy)
}

func TestCopy_on_non_pointer(t *testing.T) {
	DefaultCopyOptions.IgnoreNotFoundedFields = true
	dummyStruct1 := TestStruct{
		Dummy: "test",
		Yummy: 1,
	}

	dummyStruct2 := TestStruct{
		Dummy: "test1",
	}

	err := Copy(dummyStruct1, dummyStruct2)
	assert.Error(t, err)
}

func TestCopy_on_non_struct(t *testing.T) {
	DefaultCopyOptions.IgnoreNotFoundedFields = true
	dummyStruct1 := TestStruct{
		Dummy: "test",
		Yummy: 1,
	}

	err := Copy(dummyStruct1, "test1")
	assert.Error(t, err)
}

func TestIsZeroValue(t *testing.T) {
	var i int
	var f float64
	var b bool
	var s string
	st := TestStruct{}
	assert.True(t, IsZeroValue(i))
	assert.True(t, IsZeroValue(f))
	assert.True(t, IsZeroValue(b))
	assert.True(t, IsZeroValue(s))
	assert.True(t, IsZeroValue(st))
	assert.True(t, IsZeroValue(nil))
	i = 1
	f = 1.0
	b = true
	s = "1"
	pSt := &TestStruct{}
	assert.False(t, IsZeroValue(i))
	assert.False(t, IsZeroValue(f))
	assert.False(t, IsZeroValue(b))
	assert.False(t, IsZeroValue(s))
	assert.False(t, IsZeroValue(pSt))
}

func TestIsStruct(t *testing.T) {
	var i int
	var f float64
	var b bool
	var s string
	st := TestStruct{}
	pSt := &TestStruct{}
	assert.False(t, IsStruct(i))
	assert.False(t, IsStruct(f))
	assert.False(t, IsStruct(b))
	assert.False(t, IsStruct(s))
	assert.False(t, IsStruct(nil))
	assert.False(t, IsStruct(pSt))
	assert.True(t, IsStruct(st))
}

func items(obj interface{}) map[string]interface{} {
	m, _ := Items(obj)
	return m
}

func fields(obj interface{}) []string {
	f, _ := FieldsNames(obj)
	return f
}
