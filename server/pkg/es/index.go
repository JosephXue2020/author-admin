package es

import (
	"context"
	"encoding/json"
	"errors"
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

func FlatMappings(x interface{}, m map[string]map[string]string, depth int) error {
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

		unpack := structField.Tag.Get("unpack")
		if util.IsStructOrStructPtr(ft) && depth != 0 && unpack == "true" {
			subm := make(map[string]map[string]string)
			err := FlatMappings(fv, subm, depth-1)
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

			// 解决date类型兼容性
			if esType == "date" {
				m[key]["format"] = "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
			}
		}
	}

	return nil
}

func ESMappings(x interface{}, depth int) (map[string]map[string]map[string]string, error) {
	m := make(map[string]map[string]string)
	err := FlatMappings(x, m, depth)
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

func CreateIndex(sc Scanner) error {
	indexName := sc.IndexName() // 已保证是lowecase

	exist := IndexExist(indexName)
	m := make(map[string]interface{})
	if !exist {
		mappings := sc.Mappings()
		if mappings == nil {
			return errors.New("Failed to get mappings from orm scanner.")
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

func CreateIndices(scs ...Scanner) error {
	for _, sc := range scs {
		err := CreateIndex(sc)
		if err != nil {
			return err
		}
	}
	return nil
}
