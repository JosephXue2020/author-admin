package dao

import (
	"context"
	"encoding/json"
	"goweb/author-admin/server/pkg/setting"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var ES *elasticsearch.Client

func InitES() error {
	cfg := elasticsearch.Config{
		Addresses: setting.ESHosts,
		Username:  setting.ESUser,
		Password:  setting.ESPassword,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
		},
	}
	var err error
	ES, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Println("Failed to init ES: ", err)
		return err
	}

	resp, err := ES.Info()
	if err != nil {
		log.Println("Failed to connect ES: ", err)
		return err
	}
	defer resp.Body.Close()
	// resp, _ = ES.Cat.Health()
	// log.Println(resp)

	return nil
}

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
		_, err = ES.Indices.Create(indexName, ES.Indices.Create.WithBody(strings.NewReader(s)))
		if err != nil {
			return err
		}
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

func CreateToDBES(value interface{}) {
	DB.Create(value)
	// ES.Create()
}
