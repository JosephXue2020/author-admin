package es

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func RespBodyStr(resp *esapi.Response) string {
	byteData, _ := ioutil.ReadAll(resp.Body)
	return string(byteData)
}

func RespBodyMap(resp *esapi.Response) map[string]interface{} {
	var r map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return nil
	}

	return r
}

func RespError(resp *esapi.Response) error {
	m := RespBodyMap(resp)

	v, ok := m["error"]
	if !ok {
		return nil
	}

	byteData, _ := json.Marshal(v)
	s := string(byteData)
	if s == "" {
		return nil
	}

	return errors.New(s)
}
