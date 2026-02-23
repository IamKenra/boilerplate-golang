package slices

import (
	"reflect"
	"strings"
	"unsafe"
)

func RemoveIDFields(data any) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	removeIDFieldsRecursive(v)
}

func removeIDFieldsRecursive(v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if !v.IsNil() {
			removeIDFieldsRecursive(v.Elem())
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := v.Type().Field(i)

			jsonTag := fieldType.Tag.Get("json")
			tagParts := strings.Split(jsonTag, ",")

			if len(tagParts) > 0 && (strings.HasPrefix(tagParts[0], "id_") || tagParts[0] == "id") {
				if field.CanSet() {
					field.Set(reflect.Zero(field.Type()))
				} else if field.CanAddr() {
					setUnexportedField(field, reflect.Zero(field.Type()))
				}
			} else {
				removeIDFieldsRecursive(field)
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			removeIDFieldsRecursive(v.Index(i))
		}
	}
}

func setUnexportedField(field reflect.Value, value reflect.Value) {
	ptr := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	ptr.Set(value)
}
