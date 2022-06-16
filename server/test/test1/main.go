package main

import (
	"encoding/json"
	"fmt"
	"goweb/author-admin/server/pkg/e"
	"log"
	"reflect"

	"github.com/google/uuid"
)

func main() {
	resp := NewResponse()
	fmt.Println(resp)

	m := resp.ToMap(0)
	fmt.Printf("%+v\n", m)
	fmt.Printf("%T\n", m["meta"])

	js := resp.ToJSON(0)
	fmt.Printf("%+v\n", string(js))

	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxx")

	m = resp.ToMap(3)
	fmt.Printf("%+v\n", m)

	js = resp.ToJSON(3)
	fmt.Printf("%+v\n", string(js))
}

type Page struct {
	PageNum  int    `form:"pageNum" json:"pageNum"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Keyword  string `form:"keyword" json:"keyword"`
	Desc     bool   `form:"desc" json:"desc"`
}

type Meta struct {
	Page      `json:"page"`
	RequestID string `json:"requestid"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`

	// Meta存储如请求ID、分页等信息
	Meta Meta `json:"meta"`
}

func (resp *Response) ToMap(depth int) map[string]interface{} {
	m := make(map[string]interface{})
	err := StructToMapWithTagKey(*resp, m, depth)
	if err != nil {
		log.Println(err)
		return nil
	}

	return m
}

func (resp *Response) ToJSON(depth int) []byte {
	m := resp.ToMap(depth)

	r, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return nil
	}

	return r
}

func NewResponse() *Response {
	uuid1, _ := uuid.NewUUID()
	uuidStr := uuid1.String()
	return &Response{
		Code:    e.SUCCESS,
		Message: e.GetMsg(e.SUCCESS),
		Meta: Meta{
			RequestID: uuidStr,
		},
	}
}

func NeWReponseWithCode(code int) *Response {
	resp := NewResponse()
	resp.Code = code
	resp.Message = e.GetMsg(code)
	return resp
}

func FailedResponseMap(code int) map[string]interface{} {
	resp := NeWReponseWithCode(code)
	return resp.ToMap(0)
}

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

func GetStructTagName(structName interface{}) []string {
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
func GetStructTagNameIter(x interface{}, names []string, depth int) error {
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
			err = GetStructTagNameIter(fv, names, depth-1)
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
func StructToMapWithTagKey(x interface{}, m map[string]interface{}, depth int) error {
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
			err = StructToMapWithTagKey(fv, subm, depth-1)
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
