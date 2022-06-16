package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"goweb/author-admin/server/pkg/util"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func IndexExist(indexName string) (bool, error) {
	var exist bool
	ctx := context.Background()
	req := esapi.IndicesExistsRequest{}
	req.Index = append(req.Index, indexName)
	resp, err := req.Do(ctx, ES.Transport)
	if err != nil {
		log.Println("ES requests failed.")
		return exist, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		exist = true
	}
	return exist, nil
}

func DefaultSettings() map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func ESMappings(x interface{}) map[string]map[string]map[string]string {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	m := make(map[string]map[string]string)
	for i := 0; i < t.NumField(); i++ {
		esType := t.Field(i).Tag.Get("es")
		if esType != "" {
			esField := t.Field(i).Tag.Get("json")
			if esField == "" {
				// 没有json字段，就取field name
				esField = t.Field(i).Name
			}
			m[esField] = map[string]string{"type": esType}
		}
	}
	r := map[string]map[string]map[string]string{
		"properties": m,
	}
	return r
}

func CreateIndex(x interface{}) error {
	t := reflect.TypeOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	indexName := t.Name()
	indexName = strings.ToLower(indexName) // 必须是lowecase

	exist, err := IndexExist(indexName)
	if err != nil {
		return err
	}

	if !exist {
		m := make(map[string]interface{})
		m["settings"] = DefaultSettings()
		m["mappings"] = ESMappings(x)
		byteData, err := json.Marshal(m)
		if err != nil {
			return err
		}
		s := string(byteData)
		resp, err := ES.Indices.Create(indexName, ES.Indices.Create.WithBody(strings.NewReader(s)))
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("Failed to authorize to ES: %d\n", resp.StatusCode)
		}
		return nil
	}

	return nil
}

func CreateIndices(x ...interface{}) error {
	for _, v := range x {
		err := CreateIndex(v)
		if err != nil {
			return err
		}
	}
	return nil
}

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
	// id, ok := m["id"]
	// if !ok {
	// 	err = fmt.Errorf("db table lacks id field.")
	// 	return err
	// }
	// idStr, err := util.ItfToStr(id)
	// if err != nil {
	// 	return err
	// }

	idValue := v.FieldByName("ID")
	emptyValue := reflect.Value{}
	if idValue == emptyValue {
		err = fmt.Errorf("db table lacks id field.")
		return err
	}
	idStr, err := util.ItfToStr(idValue.Interface())
	if err != nil {
		return err
	}
	log.Println("id:", idStr)

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

	resp, err := ES.Update(indexName, idStr, &buf)
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
