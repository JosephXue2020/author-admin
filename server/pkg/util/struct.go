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
	t := reflect.TypeOf(structName)
	if !IsStructOrStructPtr(t) {
		log.Println("Parameter should be struct or struct pointer type.")
		return nil
	}

	var names []string
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		names = append(names, name)
	}
	return names
}

func GetStructJSONName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if !IsStructOrStructPtr(t) {
		log.Println("Parameter should be struct or struct pointer type.")
		return nil
	}

	var names []string
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Tag.Get("json")
		names = append(names, name)
	}
	return names
}

// Iterative
// depth<0 for no limit; depth>0 for finite iterative depth; depth=0 for no iteration.
func GetStructFieldNameIter(x interface{}, names []string, depth int) error {
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
		vField := v.Field(i)
		if !vField.CanInterface() {
			continue
		}
		fv := vField.Interface()
		if IsStructOrStructPtr(ft) && depth != 0 {
			err = GetStructFieldNameIter(fv, names, depth-1)
			continue
		}
		names = append(names, structField.Name)
	}

	return err
}

// Iterative
// depth<0 for no limit; depth>0 for finite iterative depth; depth=0 for no iteration.
func GetStructJSONNameIter(x interface{}, names []string, depth int) error {
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
		vField := v.Field(i)
		if !vField.CanInterface() {
			continue
		}
		fv := vField.Interface()
		if IsStructOrStructPtr(ft) && depth != 0 {
			err = GetStructJSONNameIter(fv, names, depth-1)
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
// depth<0 for no limit; depth>0 for finite iterative depth; depth=0 for no iteration.
func StructToMapWithFieldKey(x interface{}, m map[string]interface{}, depth int) error {
	var err error
	if m == nil {
		err = fmt.Errorf("Map parameter should not be nil before used.")
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
		key := structField.Name
		ft := structField.Type

		vField := v.Field(i)
		if !vField.CanInterface() {
			continue
		}
		fv := vField.Interface()

		if IsStructOrStructPtr(ft) && depth != 0 {
			subm := make(map[string]interface{})
			err = StructToMapWithFieldKey(fv, subm, depth-1)
			if err != nil {
				return err
			}
			m[key] = subm
		} else {
			m[key] = fv
		}
	}

	return nil
}

// Iterative
// depth<0 for no limit; depth>0 for finite iterative depth; depth=0 for no iteration.
func StructToMapWithJSONKey(x interface{}, m map[string]interface{}, depth int) error {
	var err error
	if m == nil {
		err = fmt.Errorf("Map parameter should not be nil before used.")
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
		key := structField.Tag.Get("json")
		if key == "" {
			continue
		}

		vField := v.Field(i)
		if !vField.CanInterface() {
			continue
		}
		fv := vField.Interface()

		if IsStructOrStructPtr(ft) && depth != 0 {
			subm := make(map[string]interface{})
			err = StructToMapWithJSONKey(fv, subm, depth-1)
			if err != nil {
				return err
			}
			m[key] = subm
		} else {
			m[key] = fv
		}
	}

	return nil
}

// Iterative
// depth<0 for no limit; depth>0 for finite iterative depth; depth=0 for no iteration.
func StructToFlattenMapWithFieldKey(x interface{}, m map[string]interface{}, depth int) error {
	var err error
	if m == nil {
		err = fmt.Errorf("Map parameter should not be nil before used.")
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
		key := structField.Name
		ft := structField.Type

		vField := v.Field(i)
		if !vField.CanInterface() {
			continue
		}
		fv := vField.Interface()

		if IsStructOrStructPtr(ft) && depth != 0 {
			subm := make(map[string]interface{})
			err = StructToFlattenMapWithFieldKey(fv, subm, depth-1)
			if err != nil {
				return err
			}

			mKeys := MapKeysWithStr(m)
			for k := range subm {
				if ContainStr(mKeys, k) {
					err := fmt.Errorf("Failed to convert struct to map iterablely of conflict keys.")
					return err
				}
			}

			for subk, subv := range subm {
				m[subk] = subv
			}
		} else {
			m[key] = fv
		}
	}

	return nil
}

// Iterative
// depth<0 for no limit; depth>0 for finite iterative depth; depth=0 for no iteration.
func StructToFlattenMapWithJSONKey(x interface{}, m map[string]interface{}, depth int) error {
	var err error
	if m == nil {
		err = fmt.Errorf("Map parameter should not be nil before used.")
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
		key := structField.Tag.Get("json")
		if key == "" {
			continue
		}

		vField := v.Field(i)
		if !vField.CanInterface() {
			continue
		}
		fv := vField.Interface()

		if IsStructOrStructPtr(ft) && depth != 0 {
			subm := make(map[string]interface{})
			err = StructToFlattenMapWithJSONKey(fv, subm, depth-1)
			if err != nil {
				return err
			}

			mKeys := MapKeysWithStr(m)
			for k := range subm {
				if ContainStr(mKeys, k) {
					err := fmt.Errorf("Failed to convert struct to map iterablely of conflict keys.")
					return err
				}
			}

			for subk, subv := range subm {
				m[subk] = subv
			}
		} else {
			m[key] = fv
		}
	}

	return nil
}
