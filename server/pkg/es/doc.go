package es

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/pkg/util"
	"io/ioutil"
	"log"
	"net/http"
)

func CreateDoc(indexName string, x interface{}, depth int) error {
	m := make(map[string]interface{})
	err := util.StructToFlattenMapWithJSONKey(x, m, depth)
	if err != nil {
		return err
	}
	log.Println("map from doc is: ", m)

	id, ok := m["id"]
	if !ok {
		return errors.New("db table lacks id field.")
	}
	idStr, err := util.ItfToStr(id)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		return err
	}

	log.Println("index name:", indexName)
	log.Println("idstr:", idStr)
	resp, err := dao.ES.Create(indexName, idStr, &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	r, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(r))

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

	return nil
}

func DeleteDocByID(indexName string, id int) error {

	return nil
}

func DeleteDocByIDBatch(indexName string, ids []int) error {

	return nil
}
