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
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

func DocExist(indexName string, id string) bool {
	_, err := dao.ES.Exists(indexName, id)
	if err != nil {
		return false
	}
	return true
}

type Doc struct {
	Idx        *Index
	ID         string
	FlatMap    map[string]interface{}
	Reader     io.Reader
	ReadSeeker io.ReadSeeker
	Err        error
}

func (doc *Doc) GetFlatMap() map[string]interface{} {
	m := make(map[string]interface{})
	err := util.StructToFlatMapWithJSONKey(doc.Idx.IndexScanner, m, doc.Idx.Depth)
	if err != nil {
		doc.Err = err
		return nil
	}

	doc.FlatMap = m
	return m
}

func (doc *Doc) GetID() string {
	idItf, ok := doc.FlatMap["id"]
	if !ok {
		doc.Err = errors.New("db table lacks id field.")
		return ""
	}

	id, err := util.ItfToStr(idItf)
	if err != nil {
		doc.Err = err
		return ""
	}

	return id
}

func (doc *Doc) GetReaderAndReadSeeker() {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc.FlatMap); err != nil {
		doc.Err = err
		return
	}
	doc.Reader = &buf

	s := string(buf.Bytes())
	doc.ReadSeeker = strings.NewReader(s)
	return
}

func (doc *Doc) DocExist() bool {
	return DocExist(doc.Idx.Name, doc.ID)
}

func (doc *Doc) CreateDoc() error {
	resp, err := dao.ES.Create(doc.Idx.Name, doc.ID, doc.Reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("Failed to create doc, response code: %v", resp.StatusCode)
	}

	return nil
}

func (doc *Doc) UpdateDoc() error {
	resp, err := dao.ES.Update(doc.Idx.Name, doc.ID, doc.Reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("Failed to update doc, response code: %v", resp.StatusCode)
	}

	return nil
}

func (doc *Doc) IndexDoc() error {
	req := esapi.IndexRequest{
		Index:      doc.Idx.Name,
		DocumentID: doc.ID,
		Body:       doc.Reader,
		Refresh:    "true",
	}
	resp, err := req.Do(context.Background(), dao.ES)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("Failed to Index doc, response code: %v", resp.StatusCode)
	}
	return nil
}

func NewDoc(idx *Index) (doc *Doc) {
	doc = &Doc{
		Idx: idx,
	}

	doc.GetFlatMap()
	if doc.Err != nil {
		return
	}

	doc.GetID()
	if doc.Err != nil {
		return
	}

	doc.GetReaderAndReadSeeker()
	if doc.Err != nil {
		return
	}

	return
}

func NewDocFromScanner(sc Scanner) *Doc {
	idx := NewDefaultIndex(sc)
	return NewDoc(idx)
}

// Create if not exist; update if exist.
func IndexDoc(doc *Doc) error {
	return doc.IndexDoc()
}

// Create if not exist; update if exist.
// Bulk operation.
func IndexDocBulk(docs []*Doc, workers int) error {
	// esutil.BulkIndexer()
	cfg := esutil.BulkIndexerConfig{
		NumWorkers: workers,
		Client:     dao.ES,
		FlushBytes: 1,
	}
	bi, err := esutil.NewBulkIndexer(cfg)
	if err != nil {
		return err
	}

	for _, doc := range docs {
		err = bi.Add(context.Background(), esutil.BulkIndexerItem{
			Index:      doc.Idx.Name,
			DocumentID: doc.ID,
			Body:       doc.ReadSeeker,
		})
		if err != nil {
			return err
		}
	}

	if err := bi.Close(context.Background()); err != nil {
		return err
	}

	stats := bi.Stats()
	// log.Printf("%#v\n", stats)
	if stats.NumFailed != 0 {
		err = fmt.Errorf("Number failed: %v / %v", stats.NumFailed, len(docs))
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
