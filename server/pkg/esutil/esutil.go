package esutil

import (
	"context"
	"encoding/json"
	"fmt"
	"goweb/author-admin/server/dao"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

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

func IndexExist(indexName string) bool {
	ctx := context.Background()
	req := esapi.IndicesExistsRequest{}
	req.Index = append(req.Index, indexName)
	resp, err := req.Do(ctx, dao.ES.Transport)
	if err != nil {
		log.Println("ES requests failed.")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func CreateIndex(x interface{}) error {
	t := reflect.TypeOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	indexName := t.Name()
	indexName = strings.ToLower(indexName) // 必须是lowecase

	exist := IndexExist(indexName)
	if !exist {
		m := make(map[string]interface{})
		m["settings"] = DefaultSettings()
		m["mappings"] = ESMappings(x)
		byteData, err := json.Marshal(m)
		if err != nil {
			return err
		}
		s := string(byteData)
		resp, err := dao.ES.Indices.Create(indexName, dao.ES.Indices.Create.WithBody(strings.NewReader(s)))
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
