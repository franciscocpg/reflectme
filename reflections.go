// Copyright (c) 2013 Théo Crevon
//
// See the file LICENSE for copying permission.

/*
Package reflections provides high level abstractions above the
reflect library.

Reflect library is very low-level and as can be quite complex when it comes to do simple things like accessing a structure field value, a field tag...

The purpose of reflections package is to make developers life easier when it comes to introspect structures at runtime.
It's API is freely inspired from python language (getattr, setattr, hasattr...) and provides a simplified access to structure fields and tags.
*/
package reflectme

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type (
	// CopyOptions are options for copy function
	CopyOptions struct {
		CopyZeroValues         bool
		IgnoreNotFoundedFields bool
	}
)

var (
	// DefaultCopyOptions are the default options for copy function
	DefaultCopyOptions = CopyOptions{
		CopyZeroValues:         true,
		IgnoreNotFoundedFields: true,
	}
)

// GetField returns the value of the provided obj field. obj can whether
// be a structure or pointer to structure.
func GetField(obj interface{}, name string) (interface{}, error) {
	field, err := getInnerField(obj, name)
	if err != nil {
		return nil, err
	}

	return field.Interface(), nil
}

// GetFieldKind returns the kind of the provided obj field. obj can whether
// be a structure or pointer to structure.
func GetFieldKind(obj interface{}, name string) (reflect.Kind, error) {
	field, err := getInnerField(obj, name)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Kind(), nil
}

// GetFieldTag returns the provided obj field tag value. obj can whether
// be a structure or pointer to structure.
func GetFieldTag(obj interface{}, fieldName, tagKey string) (string, error) {
	field, err := getInnerFieldType(obj, fieldName, fieldName)
	if err != nil {
		return "", err
	}

	if !isExportableField(field) {
		return "", errors.New("Cannot GetFieldTag on a non-exported struct field")
	}

	return field.Tag.Get(tagKey), nil
}

// SetField sets the provided obj field with provided value. obj param has
// to be a pointer to a struct, otherwise it will soundly fail. Provided
// value type should match with the struct field you're trying to set.
func SetField(s interface{}, name string, value interface{}) error {
	return setField(reflect.ValueOf(s), name, name, value)
}

func setField(v reflect.Value, name string, currName string, value interface{}) error {

	if v.Kind() != reflect.Ptr {
		return errors.New("Not a pointer value")
	}
	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Struct, reflect.Ptr:
		currName, nextFieldName := getCurrAndNextFieldName(currName)
		if v.Kind() == reflect.Struct {
			v = v.FieldByName(currName)
		} else {
			v = v.Elem().FieldByName(currName)
		}
		if !v.IsValid() {
			return fmt.Errorf("No such field: %s in obj", name)
		}
		err := setField(v.Addr(), name, nextFieldName, value)
		if err != nil {
			return err
		}
	default:
		valueOf := reflect.ValueOf(value)
		if v.Type() != valueOf.Type() {
			return fmt.Errorf("Provided value type (%v) didn't match obj field type (%v)\n", valueOf.Type(), v.Type())
		}
		v.Set(valueOf)

	}

	return nil
}

// HasField checks if the provided field name is part of a struct. obj can whether
// be a structure or pointer to structure.
func HasField(obj interface{}, name string) (bool, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return false, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()
	field, ok := objType.FieldByName(name)
	if !ok || !isExportableField(field) {
		return false, nil
	}

	return true, nil
}

// FieldsNames returns the struct fields names list. obj can whether
// be a structure or pointer to structure.
func FieldsNames(obj interface{}) ([]string, error) {
	return fieldsNames(obj, "")
}

func fieldsNames(obj interface{}, parent string) ([]string, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return nil, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()
	fieldsCount := objType.NumField()

	var fields []string
	for i := 0; i < fieldsCount; i++ {
		field := objType.Field(i)
		var fieldName string
		if isExportableField(field) {
			fieldName = field.Name
			if len(parent) > 0 {
				fieldName = parent + "." + fieldName
			}
			fields = append(fields, fieldName)
		}
		if k := objValue.Field(i).Kind(); k == reflect.Struct || k == reflect.Ptr {
			nestedFields, err := fieldsNames(objValue.Field(i).Interface(), fieldName)
			if err == nil {
				fields = append(fields, nestedFields...)
			} else {
				return fields, err
			}
		}
	}

	return fields, nil
}

// Fields returns the struct fields list. obj can whether
// be a structure or pointer to structure.
func Fields(obj interface{}) ([]reflect.StructField, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return nil, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()
	fieldsCount := objType.NumField()

	var fields []reflect.StructField
	for i := 0; i < fieldsCount; i++ {
		field := objType.Field(i)
		if isExportableField(field) {
			fields = append(fields, field)
		}
		if k := objValue.Field(i).Kind(); k == reflect.Struct || k == reflect.Ptr {
			nestedFields, err := Fields(objValue.Field(i).Interface())
			if err == nil {
				fields = append(fields, nestedFields...)
			} else {
				return fields, err
			}
		}
	}

	return fields, nil
}

// Items returns the field - value struct pairs as a map. obj can whether
// be a structure or pointer to structure.
func Items(obj interface{}) (map[string]interface{}, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return nil, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()
	fieldsCount := objType.NumField()

	items := make(map[string]interface{})

	for i := 0; i < fieldsCount; i++ {
		field := objType.Field(i)
		fieldValue := objValue.Field(i)

		// Make sure only exportable and addressable fields are
		// returned by Items
		if isExportableField(field) {
			items[field.Name] = fieldValue.Interface()
		}
	}

	return items, nil
}

// Tags lists the struct tag fields. obj can whether
// be a structure or pointer to structure.
func Tags(obj interface{}, key string) (map[string]string, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return nil, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()
	fieldsCount := objType.NumField()

	tags := make(map[string]string)

	for i := 0; i < fieldsCount; i++ {
		structField := objType.Field(i)

		if isExportableField(structField) {
			tags[structField.Name] = structField.Tag.Get(key)
		}
	}

	return tags, nil
}

// Copy copies all values from "from" to "to" with
// DefaultCopyOptions
func Copy(from interface{}, to interface{}) error {
	return CopyWithOptions(from, to, DefaultCopyOptions)
}

// CopyWithOptions copies all values from "from" to "to" with
// CopyOptions
func CopyWithOptions(from interface{}, to interface{}, options CopyOptions) error {
	if !isPointer(to) {
		return errors.New("To must be a pointer")
	}
	fromFields, err := FieldsNames(from)
	if err != nil {
		return err
	}
	_, err = FieldsNames(to)
	if err != nil {
		return err
	}
	for _, field := range fromFields {
		v, err := GetField(from, field)
		if err != nil {
			return err
		}
		if !options.CopyZeroValues && IsZeroValue(v) {
			continue
		}
		err = SetField(to, field, v)
		if !options.IgnoreNotFoundedFields && err != nil {
			return err
		}
	}
	return nil
}

// IsZeroValue indicates if the interface has value
// according to golang spec: https://golang.org/ref/spec#The_zero_value
func IsZeroValue(i interface{}) bool {
	return i == nil || reflect.DeepEqual(i, reflect.Zero(reflect.TypeOf(i)).Interface())
}

func reflectValue(obj interface{}) reflect.Value {
	var val reflect.Value

	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		val = reflect.ValueOf(obj).Elem()
	} else {
		val = reflect.ValueOf(obj)
	}

	return val
}

func isExportableField(field reflect.StructField) bool {
	// PkgPath is empty for exported fields.
	return field.PkgPath == ""
}

func hasValidType(obj interface{}, types []reflect.Kind) bool {
	for _, t := range types {
		if reflect.TypeOf(obj).Kind() == t {
			return true
		}
	}

	return false
}

func isStruct(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Struct
}

func isPointer(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Ptr
}

func getCurrAndNextFieldName(name string) (string, string) {
	currName := name
	nextFieldName := ""
	if i := strings.Index(name, "."); i > -1 {
		currName = name[0:i]
		nextFieldName = name[i+1 : len(name)]
	}
	return currName, nextFieldName
}

func getInnerField(obj interface{}, name string) (reflect.Value, error) {
	field := reflect.Value{}
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return field, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	if i := strings.Index(name, "."); i > -1 {
		currFieldName := name[0:i]
		field = objValue.FieldByName(currFieldName)
		if !field.IsValid() {
			return field, fmt.Errorf("No such field: %s in1 obj", name)
		}
		if !isStruct(field) {
			return field, fmt.Errorf("Field %s expected to be an struct", currFieldName)
		}
		nextFieldName := name[i+1 : len(name)]
		return getInnerField(field.Interface(), nextFieldName)
	}
	field = objValue.FieldByName(name)
	if !field.IsValid() {
		return field, fmt.Errorf("No such field: %s in obj", name)
	}
	return field, nil
}

func getInnerFieldType(obj interface{}, fullName, name string) (reflect.StructField, error) {
	field := reflect.StructField{}
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return field, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	if i := strings.Index(name, "."); i > -1 {
		currFieldName := name[0:i]
		fieldRv := objValue.FieldByName(currFieldName)
		if !fieldRv.IsValid() {
			return field, fmt.Errorf("No such field: %s in obj", name)
		}
		if !isStruct(fieldRv) {
			return field, fmt.Errorf("Field %s expected to be an struct", currFieldName)
		}
		nextFieldName := name[i+1 : len(name)]
		return getInnerFieldType(fieldRv.Interface(), fullName, nextFieldName)
	}
	if !objValue.IsValid() {
		return field, fmt.Errorf("Nil pointer: %s in obj", fullName)
	}
	objType := objValue.Type()
	field, ok := objType.FieldByName(name)
	if !ok {
		return field, fmt.Errorf("No such field: %s in obj", name)
	}
	return field, nil
}
