package util

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/google/uuid"
)

func PressAnyKeyToExit() {
	var x string
	fmt.Println("Press any key to exit.")
	fmt.Scan(&x)
}

func ContainStr(set []string, v string) bool {
	for _, item := range set {
		if item == v {
			return true
		}
	}
	return false
}

func ContainInt(set []int, v int) bool {
	for _, item := range set {
		if item == v {
			return true
		}
	}
	return false
}

func ContainRefl(set interface{}, v interface{}) bool {
	setT := reflect.TypeOf(set)
	if setT.Kind() != reflect.Slice {
		log.Fatal("First parameter should be slice of some type.")
	}

	vT := reflect.TypeOf(v)
	if !vT.Comparable() {
		log.Fatal("Element of set and 2rd parameter should be comparable.")
	}

	setV := reflect.ValueOf(set)
	firstV := setV.Index(0)
	firstT := firstV.Type().Name()
	if firstT != vT.Name() {
		log.Fatal("The element in set must have the same type as value to inspect.")
	}

	for i := 0; i < setV.Len(); i++ {
		if setV.Index(i).Interface() == v {
			return true
		}
	}
	return false
}

func CurrentTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GenUUID() string {
	u, _ := uuid.NewUUID() // v1 uuid
	s := u.String()
	return s
}

func IsStruct(x interface{}) bool {
	typ := reflect.TypeOf(x)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return false
	}

	return true
}

func GetStructFieldName(structName interface{}) []string {
	if !IsStruct(structName) {
		log.Println("Parameter should be struct or struct pointer type.")
		return nil
	}

	var names []string
	typ := reflect.TypeOf(structName)
	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Name
		names = append(names, name)
	}
	return names
}

func GetStructTagName(structName interface{}) []string {
	if !IsStruct(structName) {
		log.Println("Parameter should be struct or struct pointer type.")
		return nil
	}

	var names []string
	typ := reflect.TypeOf(structName)
	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Tag.Get("json")
		names = append(names, name)
	}
	return names
}

func StructToMap(structName interface{}) map[string]interface{} {
	if !IsStruct(structName) {
		log.Println("Parameter should be struct or struct pointer type.")
		return nil
	}

	m := make(map[string]interface{})
	typ := reflect.TypeOf(structName)
	val := reflect.ValueOf(structName)
	for i := 0; i < val.NumField(); i++ {
		key := typ.Field(i).Name
		v := val.Field(i).Interface()
		m[key] = v
	}

	return m
}

func StructToMapWithTagKey(structName interface{}) map[string]interface{} {
	if !IsStruct(structName) {
		log.Println("Parameter should be struct or struct pointer type.")
		return nil
	}

	m := make(map[string]interface{})
	typ := reflect.TypeOf(structName)
	val := reflect.ValueOf(structName)
	for i := 0; i < val.NumField(); i++ {
		key := typ.Field(i).Tag.Get("json")
		if key != "" {
			v := val.Field(i).Interface()
			m[key] = v
		}
	}

	return m
}
