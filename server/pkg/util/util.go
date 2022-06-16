package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
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

func MapToJsonWithStrKey(m map[string]interface{}) ([]byte, error) {
	var byteData []byte
	if m == nil {
		err := fmt.Errorf("The input map is nil.")
		return byteData, err
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		return byteData, err
	}

	return byteData, nil
}

func ItfToStr(x interface{}) (string, error) {
	switch x.(type) {
	case string:
		return x.(string), nil
	case int, int8, int16, int32, int64:
		i := x.(int)
		return strconv.Itoa(i), nil
	case uint, uint8, uint16, uint32, uint64:
		i := x.(uint)
		return strconv.Itoa(int(i)), nil
	default:
		err := fmt.Errorf("Failed to convert.")
		return "", err
	}
}

func MapKeys(m map[string]interface{}) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
