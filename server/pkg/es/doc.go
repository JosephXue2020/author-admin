package es

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/pkg/util"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type DocMeta struct {
	IndexName  string
	ID         string
	FlatMap    map[string]interface{}
	Reader     io.Reader
	ReadSeeker io.ReadSeeker
}

// StructToMapForES convert struct to flat map with given depth.
// A better choice is always to use this function to make such a conversion.
func StructToMapForES(x interface{}, depth int) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := util.StructToFlatMapWithJSONKey(x, m, depth)
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

func NewDocMeta(indexName string, x interface{}, depth int) (*DocMeta, error) {
	m, err := StructToMapForES(x, depth)
	if err != nil {
		return nil, err
	}

	id, ok := m["id"]
	if !ok {
		err = errors.New("db table lacks id field.")
		return nil, err
	}

	idStr, err := util.ItfToStr(id)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		return nil, err
	}

	s := string(buf.Bytes())
	readSeeker := strings.NewReader(s)

	meta := &DocMeta{
		IndexName:  indexName,
		ID:         idStr,
		FlatMap:    m,
		Reader:     &buf,
		ReadSeeker: readSeeker,
	}
	return meta, nil
}

func DocExists(indexName string, id string) bool {
	_, err := dao.ES.Exists(indexName, id)
	if err != nil {
		return false
	}
	return true
}

func CreateDoc(indexName string, x interface{}, depth int) error {
	meta, err := NewDocMeta(indexName, x, depth)
	if err != nil {
		return err
	}
	// log.Println("meta:", meta.FlatMap)

	resp, err := dao.ES.Create(indexName, meta.ID, meta.Reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// err := fmt.Errorf("Wrong response code: %v\n", resp.StatusCode)
		err := RespError(resp)
		return err
	}

	return nil
}

func UpdateDoc(indexName string, x interface{}, depth int) error {
	meta, err := NewDocMeta(indexName, x, depth)
	if err != nil {
		return err
	}

	resp, err := dao.ES.Update(indexName, meta.ID, meta.Reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Wrong response code: %v\n", resp.StatusCode)
		return err
	}

	return nil
}

// Create if not exist; update if exist.
func IndexDoc(indexName string, x interface{}, depth int) error {
	meta, err := NewDocMeta(indexName, x, depth)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: meta.ID,
		Body:       meta.Reader,
		Refresh:    "true",
	}
	resp, err := req.Do(context.Background(), dao.ES)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Create if not exist; update if exist.
// Bulk operation.
func IndexDocBulk(indexName string, xs interface{}, depth int) error {
	var xSlc []interface{}
	t := reflect.TypeOf(xs)
	v := reflect.ValueOf(xs)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	for i := 0; i < v.Len(); i++ {
		xSlc = append(xSlc, v.Index(i).Interface())
	}

	var metas []*DocMeta
	for _, x := range xSlc {
		meta, err := NewDocMeta(indexName, x, depth)
		if err != nil {
			return err
		}
		metas = append(metas, meta)
	}

	// esutil.BulkIndexer()
	cfg := esutil.BulkIndexerConfig{
		NumWorkers: 10,
		Client:     dao.ES,
		FlushBytes: 1,
	}
	bi, err := esutil.NewBulkIndexer(cfg)
	if err != nil {
		return err
	}

	for _, meta := range metas {
		err = bi.Add(context.Background(), esutil.BulkIndexerItem{
			Index:      indexName,
			DocumentID: meta.ID,
			Body:       meta.ReadSeeker,
		})
		// log.Println(meta.FlatMap)
	}
	if err != nil {
		return err
	}

	if err := bi.Close(context.Background()); err != nil {
		return err
	}

	stats := bi.Stats()
	// log.Printf("%#v\n", stats)
	if stats.NumFailed != 0 {
		err = fmt.Errorf("Number failed: %v / %v", stats.NumFailed, len(metas))
		return err
	}
	return nil
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
