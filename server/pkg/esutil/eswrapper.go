package esutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/pkg/util"
	"net/http"
	"reflect"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

type ESWrapper elasticsearch.Client

func (es *ESWrapper) CreateDoc(x interface{}) error {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	flag := util.IsStructOrStructPtr(t)
	if !flag {
		err := fmt.Errorf("Param should be struct or struct pointer type.")
		return err
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	m := make(map[string]interface{})
	err := util.StructToMapWithJSONKey(x, m, 0)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		return err
	}

	indexName := strings.ToLower(t.Name())
	// Model是被嵌套的
	id, ok := m["id"]
	if !ok {
		err = fmt.Errorf("db table lacks id field.")
		return err
	}
	idStr, err := util.ItfToStr(id)
	if err != nil {
		return err
	}

	resp, err := es.Create(indexName, idStr, &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		err := fmt.Errorf("Wrong response code: %v\n", resp.StatusCode)
		return err
	}

	return nil
}

func (es *ESWrapper) DeleteDocByID(x interface{}, id int) error {
	t := reflect.TypeOf(x)
	flag := util.IsStructOrStructPtr(t)
	if !flag {
		err := fmt.Errorf("Param should be struct or struct pointer type.")
		return err
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	m := make(map[string]interface{})
	err := util.StructToMapWithJSONKey(x, m, 0)
	if err != nil {
		return err
	}

	indexName := strings.ToLower(t.Name())
	idStr, err := util.ItfToStr(id)
	if err != nil {
		return err
	}

	resp, err := es.Delete(indexName, idStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusCreated {
	// 	err := fmt.Errorf("Wrong response code: %v\n", resp.StatusCode)
	// 	return err
	// }

	return nil
}

func (es *ESWrapper) UpdateDoc(x interface{}) error {
	t := reflect.TypeOf(x)
	flag := util.IsStructOrStructPtr(t)
	if !flag {
		err := fmt.Errorf("Param should be struct or struct pointer type.")
		return err
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	m := make(map[string]interface{})
	err := util.StructToMapWithJSONKey(x, m, 0)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		return err
	}

	indexName := strings.ToLower(t.Name())
	id, ok := m["id"]
	if !ok {
		err = fmt.Errorf("db table lacks id field.")
		return err
	}
	idStr, err := util.ItfToStr(id)
	if err != nil {
		return err
	}

	resp, err := dao.ES.Update(indexName, idStr, &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
