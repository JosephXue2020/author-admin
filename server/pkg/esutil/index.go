package esutil

import (
	"context"
	"encoding/json"
	"fmt"
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/pkg/util"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

/**
	All ES indices have flatten structure.
	The structs will be transformed into flatten maps iteratively.
	Attention: The iterative depth is set to 2, which means it can parse struct with embedding 2.
**/

func DefaultSettings() map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func FlattenMappings(x interface{}, m map[string]map[string]string, depth int) error {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
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

		if util.IsStructOrStructPtr(ft) && depth != 0 {
			subm := make(map[string]map[string]string)
			err := FlattenMappings(fv, subm, depth-1)
			if err != nil {
				return err
			}

			mKeys := util.MapKeysWithStr(m)
			for k := range subm {
				if util.ContainStr(mKeys, k) {
					err := fmt.Errorf("Failed to get mappings iterablely since conflict keys.")
					return err
				}
			}

			for subK, subV := range subm {
				m[subK] = subV
			}
		} else {
			esType := structField.Tag.Get("es")
			if esType == "" {
				continue
			}

			m[key] = map[string]string{
				"type": esType,
			}
		}
	}

	return nil
}

func ESMappings(x interface{}) (map[string]map[string]map[string]string, error) {
	m := make(map[string]map[string]string)
	err := FlattenMappings(x, m, 2)
	if err != nil {
		return nil, err
	}

	r := map[string]map[string]map[string]string{
		"properties": m,
	}
	return r, nil
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
		mappings, err := ESMappings(x)
		if err != nil {
			return err
		}
		m["mappings"] = mappings
		m["settings"] = DefaultSettings()

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
