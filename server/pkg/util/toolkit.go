package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/google/uuid"
)

// PressAnyKeyToExit used to stop the exits of main process.
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

// ContainRefl returns whether an element of any type contained in a slice.
// It uses reflect and empty interface.
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

func GenUUID() string {
	u, _ := uuid.NewUUID() // v1 uuid
	s := u.String()
	return s
}

func ItfToStr(x interface{}) (string, error) {
	switch x.(type) {
	case string:
		return x.(string), nil
	case int:
		i := x.(int)
		return strconv.Itoa(i), nil
	case int8:
		i := x.(int8)
		return strconv.Itoa(int(i)), nil
	case int16:
		i := x.(int16)
		return strconv.Itoa(int(i)), nil
	case int32:
		i := x.(int32)
		return strconv.Itoa(int(i)), nil
	case int64:
		i := x.(int64)
		return strconv.Itoa(int(i)), nil
	case uint:
		i := x.(uint)
		return strconv.Itoa(int(i)), nil
	case uint8:
		i := x.(uint8)
		return strconv.Itoa(int(i)), nil
	case uint16:
		i := x.(uint16)
		return strconv.Itoa(int(i)), nil
	case uint32:
		i := x.(uint32)
		return strconv.Itoa(int(i)), nil
	case uint64:
		i := x.(uint64)
		return strconv.Itoa(int(i)), nil
	default:
		err := fmt.Errorf("Failed to convert.")
		return "", err
	}
}

func MapToJsonWithStrKey(m map[string]interface{}) ([]byte, error) {
	if m == nil {
		err := fmt.Errorf("The input map is nil.")
		return nil, err
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		return nil, err
	}

	byteData := buf.Bytes()
	return byteData, nil
}

func MapKeysWithStr(x interface{}) []string {
	var keys []string

	t := reflect.TypeOf(x)
	if t.Kind() != reflect.Map {
		log.Println("Input param is not a map.")
		return keys
	}

	keyStructs := reflect.ValueOf(x).MapKeys()
	if len(keyStructs) == 0 {
		return keys
	}

	kt := keyStructs[0].Type()
	if kt.Kind() != reflect.String {
		log.Println("Keys are not string type.")
		return keys
	}

	for _, v := range keyStructs {
		keys = append(keys, v.String())
	}
	return keys
}

func PathExist(p string) bool {
	_, err := os.Lstat(p)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func FileExist(p string) bool {
	info, err := os.Stat(p)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}
