package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PageSize    int
	PageUpbound int
	JwtIssuer   string
	JwtSecret   string
)

type MysqlConf struct {
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var Mysql MysqlConf

func init() {
	var err error
	Cfg, err = ini.Load("conf/conf.ini")
	if err != nil {
		log.Fatalf("Fail to load 'conf/conf.ini': %v", err)
	}

	LoadMode()
	LoadServer()
	LoadApp()
	LoadMysql()
}

func LoadMode() {
	RunMode = Cfg.Section("mode").Key("RUN_MODE").MustString(("debug"))
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(20005)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtIssuer = sec.Key("JWT_ISSUER").MustString("abc")
	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)@!)")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
	PageUpbound = sec.Key("PAGE_UPBOUND").MustInt(1000)
}

func LoadMysql() {
	sec, err := Cfg.GetSection("mysql")
	if err != nil {
		log.Fatalf("Fail to get section 'mysql': %v", err)
	}

	user := sec.Key("USER").MustString("root")
	password := sec.Key("PASSWORD").MustString("root")
	host := sec.Key("HOST").MustString("127.0.0.1:3306")
	name := sec.Key("NAME").MustString("authoradmin")
	tablePrefix := sec.Key("TABLE_PREFIX").MustString("authoradmin_")

	Mysql.User = user
	Mysql.Password = password
	Mysql.Host = host
	Mysql.Name = name
	Mysql.TablePrefix = tablePrefix
}
