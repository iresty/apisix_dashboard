/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package conf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

const (
	EnvPROD  = "prod"
	EnvBETA  = "beta"
	EnvDEV   = "dev"
	EnvLOCAL = "local"
	EnvTEST  = "test"

	WebDir = "html/"
)

var (
	ENV              string
	Schema           gjson.Result
	WorkDir          = "."
	ServerHost       = "0.0.0.0"
	ServerPort       = 80
	SSLHost          = "0.0.0.0"
	SSLPort          = 443
	SSLCert          string
	SSLKey           string
	ETCDConfig       *Etcd
	ErrorLogLevel    = "warn"
	ErrorLogPath     = "logs/error.log"
	AccessLogPath    = "logs/access.log"
	AuthConf         Authentication
	SSLDefaultStatus = 1 //enable ssl by default
	ImportSizeLimit  = 10 * 1024 * 1024
	PIDPath          = "/tmp/manager-api.pid"
	AllowList        []string
	Plugins          = map[string]bool{}
)

type MTLS struct {
	CaFile   string `yaml:"ca_file"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type Etcd struct {
	Endpoints []string
	Username  string
	Password  string
	MTLS      *MTLS
	Prefix    string
}

type SSL struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

type Listen struct {
	Host string
	Port int
}

type ErrorLog struct {
	Level    string
	FilePath string `yaml:"file_path"`
}

type AccessLog struct {
	FilePath string `yaml:"file_path"`
}

type Log struct {
	ErrorLog  ErrorLog  `yaml:"error_log"`
	AccessLog AccessLog `yaml:"access_log"`
}

type Conf struct {
	Etcd      Etcd
	Listen    Listen
	SSL       SSL
	Log       Log
	AllowList []string `yaml:"allow_list"`
	MaxCpu    int      `yaml:"max_cpu"`
}

type Config struct {
	Conf           Conf
	Authentication Authentication
	Plugins        []string
}

// TODO: we should no longer use init() function after remove all handler's integration tests
// ENV=test is for integration tests only, other ENV should call "InitConf" explicitly
func init() {
	if env := os.Getenv("ENV"); env == EnvTEST {
		InitConf()
	}
}

func InitConf() {
	//go test
	if workDir := os.Getenv("APISIX_API_WORKDIR"); workDir != "" {
		WorkDir = workDir
	}

	setConf()
	setEnvironment()
	initSchema()
}

func setConf() {
	filePath := WorkDir + "/conf/conf.yaml"
	if configurationContent, err := ioutil.ReadFile(filePath); err != nil {
		panic(fmt.Sprintf("fail to read configuration: %s", filePath))
	} else {
		// configuration := gjson.ParseBytes(configurationContent)
		config := Config{}
		err := yaml.Unmarshal(configurationContent, &config)
		if err != nil {
			log.Printf("conf: %s, error: %v", configurationContent, err)
		}

		// listen
		if config.Conf.Listen.Port != 0 {
			ServerPort = config.Conf.Listen.Port
		}
		if config.Conf.Listen.Host != "" {
			ServerHost = config.Conf.Listen.Host
		}

		// SSL
		if config.Conf.SSL.Port != 0 {
			SSLPort = config.Conf.SSL.Port
		}
		if config.Conf.SSL.Cert != "" {
			SSLCert = config.Conf.SSL.Cert
		}
		if config.Conf.SSL.Key != "" {
			SSLKey = config.Conf.SSL.Key
		}

		// for etcd
		if len(config.Conf.Etcd.Endpoints) > 0 {
			initEtcdConfig(config.Conf.Etcd)
		}

		// error log
		if config.Conf.Log.ErrorLog.Level != "" {
			ErrorLogLevel = config.Conf.Log.ErrorLog.Level
		}
		if config.Conf.Log.ErrorLog.FilePath != "" {
			ErrorLogPath = config.Conf.Log.ErrorLog.FilePath
		}

		// access log
		if config.Conf.Log.AccessLog.FilePath != "" {
			AccessLogPath = config.Conf.Log.AccessLog.FilePath
		}

		if !filepath.IsAbs(ErrorLogPath) {
			if strings.HasPrefix(ErrorLogPath, "winfile") {
				return
			}
			ErrorLogPath, err = filepath.Abs(filepath.Join(WorkDir, ErrorLogPath))
			if err != nil {
				panic(err)
			}
			if runtime.GOOS == "windows" {
				ErrorLogPath = `winfile:///` + ErrorLogPath
			}
		}
		if !filepath.IsAbs(AccessLogPath) {
			if strings.HasPrefix(AccessLogPath, "winfile") {
				return
			}
			AccessLogPath, err = filepath.Abs(filepath.Join(WorkDir, AccessLogPath))
			if err != nil {
				panic(err)
			}
			if runtime.GOOS == "windows" {
				AccessLogPath = `winfile:///` + AccessLogPath
			}
		}

		AllowList = config.Conf.AllowList

		// set degree of parallelism
		initParallelism(config.Conf.MaxCpu)

		// auth
		initAuthentication(config.Authentication)

		initPlugins(config.Plugins)
	}
}

func setEnvironment() {
	ENV = EnvPROD
	if env := os.Getenv("ENV"); env != "" {
		ENV = env
	}
}

func initPlugins(plugins []string) {
	for _, pluginName := range plugins {
		Plugins[pluginName] = true
	}
}

func initSchema() {
	filePath := WorkDir + "/conf/schema.json"
	if schemaContent, err := ioutil.ReadFile(filePath); err != nil {
		panic(fmt.Sprintf("fail to read configuration: %s", filePath))
	} else {
		Schema = gjson.ParseBytes(schemaContent)
	}
}

// initialize etcd config
func initEtcdConfig(conf Etcd) {
	var endpoints = []string{"127.0.0.1:2379"}
	if len(conf.Endpoints) > 0 {
		endpoints = conf.Endpoints
	}

	prefix := "/apisix"
	if len(conf.Prefix) > 0 {
		prefix = conf.Prefix
	}

	ETCDConfig = &Etcd{
		Endpoints: endpoints,
		Username:  conf.Username,
		Password:  conf.Password,
		MTLS:      conf.MTLS,
		Prefix:    prefix,
	}
}

// initialize parallelism settings
func initParallelism(choiceCores int) {
	if choiceCores < 1 {
		return
	}
	maxSupportedCores := runtime.NumCPU()

	if choiceCores > maxSupportedCores {
		choiceCores = maxSupportedCores
	}
	runtime.GOMAXPROCS(choiceCores)
}
