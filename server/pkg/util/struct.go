package util

import (
	"fmt"
	"log"
	"reflect"
)

func IsStructOrStructPtr(t reflect.Type) bool {
	if t.Kind() == reflect.Struct {
		return true
	}
	if t.Kind() == reflect.Ptr {
		elem := t.Elem()
		if elem.Kind() == reflect.Struct {
			return true
		}
	}
	return false
}

func GetStructFieldName(structName interface{}) []string {
	typ := reflect.TypeOf(structName)
	if !IsStructOrStructPtr(typ) {
		log.Println("Parameter should be struct or struct pointer type.")
		return nil
	}

	var names []string
	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Name
		names = append(names, name)
	}
	return names
}

func GetStructTagName(structName interface{}) []string {
	typ := reflect.TypeOf(structName)
	if !IsStructOrStructPtr(typ) {
		log.Println("Parameter should be struct or struct pointer type.")
		return nil
	}

	var names []string
	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Tag.Get("json")
		names = append(names, name)
	}
	return names
}

func GetStructFieldNameIter(x interface{}, names []string) error {
	var err error
	if names == nil {
		err = fmt.Errorf("Slice parameter should not be nil before used.")
		return err
	}

	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		structField := t.Field(i)
		ft := structField.Type
		fv := v.Field(i).Interface()
		if IsStructOrStructPtr(ft) {
			err = GetStructFieldNameIter(fv, names)
			continue
		}
		names = append(names, structField.Name)
	}

	return err
}

func GetStructTagNameIter(x interface{}, names []string) error {
	var err error
	if names == nil {
		err = fmt.Errorf("Slice parameter should not be nil before used.")
		return err
	}

	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		structField := t.Field(i)
		ft := structField.Type
		fv := v.Field(i).Interface()
		if IsStructOrStructPtr(ft) {
			err = GetStructTagNameIter(fv, names)
			continue
		}
		name := structField.Tag.Get("json")
		if name != "" {
			names = append(names, structField.Name)
		}
	}

	return err
}

// Iterative
func StructToMapWithFieldKey(x interface{}, m map[string]interface{}) error {
	if m == nil {
		return fmt.Errorf("Map parameter should not be nil before used.")
	}

	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		structField := t.Field(i)
		ft := structField.Type
		fv := v.Field(i).Interface()
		if IsStructOrStructPtr(ft) {
			StructToMapWithFieldKey(fv, m)
			continue
		}
		key := structField.Name
		m[key] = fv
	}

	return nil
}

// Iterative
func StructToMapWithTagKey(x interface{}, m map[string]interface{}) error {
	if m == nil {
		return fmt.Errorf("Map parameter should not be nil before used.")
	}

	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		structField := t.Field(i)
		ft := structField.Type
		fv := v.Field(i).Interface()
		if IsStructOrStructPtr(ft) {
			StructToMapWithTagKey(fv, m)
			continue
		}

		key := structField.Tag.Get("json")
		if key == "" {
			continue
		}
		m[key] = fv
	}

	return nil
}
