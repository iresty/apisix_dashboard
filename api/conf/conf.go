package conf

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

const ServerPort = 8080
const PROD = "prod"
const BETA = "beta"
const DEV = "dev"
const LOCAL = "local"
const TimeLayout = "2006-01-02 15:04:05"
const TimeTLayout = "2006-01-02T15:04:05"
const DateLayout = "2006-01-02"
const confPath = "/root/api7-manager-api/conf.json"
const RequestId = "requestId"

var (
	ENV      string
	basePath string
	ApiKey   = "edd1c9f034335f136f87ad84b625c8f1"
	BaseUrl  = "http://127.0.0.1:9080/apisix/admin"
)

func init() {
	setEnvironment()
	initMysql()
	initApisix()
}

func setEnvironment() {
	if env := os.Getenv("ENV"); env == "" {
		ENV = LOCAL
	} else {
		ENV = env
	}
	_, basePath, _, _ = runtime.Caller(1)
}

func configurationPath() string {
	if ENV == LOCAL {
		return filepath.Join(filepath.Dir(basePath), "conf.json")
	} else {
		return confPath
	}
}

type mysqlConfig struct {
	Address  string
	User     string
	Password string

	MaxConns     int
	MaxIdleConns int
	MaxLifeTime  int
}

var MysqlConfig mysqlConfig

func initMysql() {
	filePath := configurationPath()
	if configurationContent, err := ioutil.ReadFile(filePath); err != nil {
		panic(fmt.Sprintf("fail to read configuration: %s", filePath))
	} else {
		configuration := gjson.ParseBytes(configurationContent)
		mysqlConf := configuration.Get("conf.mysql")
		MysqlConfig.Address = mysqlConf.Get("address").String()
		MysqlConfig.User = mysqlConf.Get("user").String()
		MysqlConfig.Password = mysqlConf.Get("password").String()
		MysqlConfig.MaxConns = int(mysqlConf.Get("maxConns").Int())
		MysqlConfig.MaxIdleConns = int(mysqlConf.Get("maxIdleConns").Int())
		MysqlConfig.MaxLifeTime = int(mysqlConf.Get("maxLifeTime").Int())
	}
}

func initApisix() {
	filePath := configurationPath()
	if configurationContent, err := ioutil.ReadFile(filePath); err != nil {
		panic(fmt.Sprintf("fail to read configuration: %s", filePath))
	} else {
		configuration := gjson.ParseBytes(configurationContent)
		apisixConf := configuration.Get("conf.apisix")
		BaseUrl = apisixConf.Get("base_url").String()
		ApiKey = apisixConf.Get("api_key").String()
	}
}
