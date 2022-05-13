package dao

import (
	"fmt"
	"goweb/author-admin/server/pkg/setting"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitMySQL() error {
	user := setting.Mysql.User
	password := setting.Mysql.Password
	host := setting.Mysql.Host
	name := setting.Mysql.Name
	tablePrefix := setting.Mysql.TablePrefix
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, name)
	var err error
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		return err
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	DB.SingularTable(true) //表名不自动转为复数
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)

	// Ping方法判断是否能够连接
	err = DB.DB().Ping()

	// gorm调试模式：显示sql语句
	// DB.LogMode(true)

	return err
}

func Close() error {
	err := DB.Close()
	return err
}
