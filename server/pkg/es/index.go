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
**/

// Scanner scan DB rows which have been updated or deleted.
type Scanner interface {
	ScanUpdate(int, int) []Scanner
	ScanDelete(int, int) []Scanner
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

func DefaultSettings() map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

// Each Index instance defines an index in ES server.
type Index struct {
	IndexScanner Scanner
	Name         string
	// Depth indicate the iterative depth of Scanner instance.
	Depth    int
	Mappings map[string]map[string]map[string]string
	Err      error
}

func (idx *Index) GetName() string {
	if idx.IndexScanner == nil {
		idx.Err = errors.New("Failed to get the index Name from Scanner.")
		return ""
	}

	n := strings.ToLower(util.GetStructName(idx.IndexScanner))
	idx.Name = n
	return n
}

func (idx *Index) GetMappings() map[string]map[string]map[string]string {
	m := make(map[string]map[string]string)
	err := FlatMappings(idx.IndexScanner, m, idx.Depth)
	if err != nil {
		idx.Err = errors.New("Failed to get the index Mappings from Scanner.")
		return nil
	}

	r := map[string]map[string]map[string]string{
		"properties": m,
	}
	idx.Mappings = r
	return r
}

func (idx *Index) IndexExist() bool {
	return IndexExist(idx.Name)
}

func (idx *Index) CreateIndex() error {
	exist := idx.IndexExist()
	m := make(map[string]interface{})
	if !exist {
		mappings := idx.GetMappings()
		if mappings == nil {
			return idx.Err
		}
		m["mappings"] = mappings
		m["settings"] = DefaultSettings()

		byteData, err := json.Marshal(m)
		if err != nil {
			return err
		}
		s := string(byteData)
		resp, err := dao.ES.Indices.Create(idx.Name, dao.ES.Indices.Create.WithBody(strings.NewReader(s)))
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

func (idx *Index) AutoMigrate() error {
	return idx.AutoMigrate()
}

// Set default iteration depth to 10.
func NewDefaultIndex(sc Scanner) *Index {
	idx := &Index{
		IndexScanner: sc,
		Depth:        10,
	}

	idx.GetName()
	if idx.Err != nil {
		log.Println("Failed to create Index instance.")
		return nil
	}
	return idx
}

func CreateIndex(idx *Index) error {
	idx.CreateIndex()
	if idx.Err != nil {
		return idx.Err
	}
	return nil
}

func CreateIndices(ids ...Index) error {
	for _, idx := range ids {
		err := idx.CreateIndex()
		if err != nil {
			return err
		}
	}
	return nil
}
