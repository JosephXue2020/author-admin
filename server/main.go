package main

import (
	"fmt"
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/models"
	"goweb/author-admin/server/pkg/setting"
	"goweb/author-admin/server/pkg/util"
	"goweb/author-admin/server/router"
	"log"
	"net/http"
)

func main() {
	defer util.PressAnyKeyToExit()

	// 初始化数据库连接
	err := dao.InitMySQL()
	if err != nil {
		log.Panic("Can not connect sql database.")
	} else {
		log.Println("Success to connect sql database.")
	}
	defer dao.Close()

	// 初始化模型
	models.AutoMigrate()

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
