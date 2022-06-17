package main

import (
	"fmt"
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/indices"
	"goweb/author-admin/server/models"
	"goweb/author-admin/server/pkg/setting"
	"goweb/author-admin/server/pkg/util"
	"goweb/author-admin/server/router"
	"log"
	"net/http"
)

// @title author-admin restful API
// @version 0.0.1
// @description This is the backend of author-admin system.

// @host 127.0.0.1:20005
// @BasePath /

func main() {
	defer util.PressAnyKeyToExit()

	// 初始化数据库连接
	err := dao.InitMySQL()
	if err != nil {
		log.Panic("Can not connect sql database: ", err)
	} else {
		log.Println("Success to connect sql database.")
	}
	defer dao.CloseDB()

	// 初始化ES连接
	err = dao.InitES()
	if err != nil {
		log.Panic("Can not connect ES service: ", err)
	} else {
		log.Println("Success to connect ES service.")
	}

	// 初始化DB表格
	err = models.InitModels()
	if err != nil {
		log.Panic("Can not migrate DB tables: ", err)
	} else {
		log.Println("Success to migrate DB tables.")
	}

	// 初始化ES索引
	err = indices.InitIndices()
	if err != nil {
		log.Panic("Can not migrate ES indices: ", err)
	} else {
		log.Println("Success to migrate ES indices.")
	}

	// 初始化路由
	r := router.InitRouter()

	log.Println("Web server is serving at port: ", setting.HTTPPort)
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        r,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	err = s.ListenAndServe()
	if err != nil {
		log.Println("Fail to start server, err:", err)
	}
}
