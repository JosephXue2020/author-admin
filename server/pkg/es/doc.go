package es

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/pkg/util"
	"net/http"
)

// StructToMapForES convert struct to flat map with given depth.
// A better choice is always to use this function to make such a conversion.
func StructToMapForES(x interface{}, depth int) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := util.StructToFlattenMapWithJSONKey(x, m, depth)
	if err != nil {
		return nil, err
	}

	id, ok := m["id"]
	if !ok {
		err = errors.New("db table lacks id field.")
		return nil, err
	}

	_, err = util.ItfToStr(id)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func DocExists(indexName string, id string) bool {
	_, err := dao.ES.Exists(indexName, id)
	if err != nil {
		return false
	}
	return true
}

func CreateDoc(indexName string, x interface{}, depth int) error {
	m, err := StructToMapForES(x, depth)
	if err != nil {
		return err
	}

	idStr, _ := util.ItfToStr(m["id"])

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		return err
	}

	resp, err := dao.ES.Create(indexName, idStr, &buf)
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

func CreateDocBatch(indexName string, x interface{}, depth int) error {

	return nil
}

func UpdateDoc(indexName string, x interface{}, depth int) error {
	m, err := StructToMapForES(x, depth)
	if err != nil {
		return err
	}

	idStr, _ := util.ItfToStr(m["id"])

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		return err
	}

	resp, err := dao.ES.Update(indexName, idStr, &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func CreateUpdate(indexName string, x interface{}, depth int) error {
	m, err := StructToMapForES(x, depth)
	if err != nil {
		return err
	}

	idStr, _ := util.ItfToStr(m["id"])

	exist := DocExists(indexName, idStr)
	if !exist {
		return CreateDoc(indexName, x, depth)
	} else {
		return UpdateDoc(indexName, x, depth)
	}
}

func DeleteDocByID(indexName string, id int) error {

	return nil
}

func DeleteDoc(indexName string, x interface{}, depth int) error {

	return nil
}

func DeleteDocBatch(indexName string, xs []interface{}) error {

	return nil
}
