package dao

import (
	"goweb/author-admin/server/pkg/setting"
	"log"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
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

	resp, _ = ES.Cat.Health()
	log.Println("ES cluster health status: ", resp)

	return nil
}
