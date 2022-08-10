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
package internal

import (
	"fmt"
	"path/filepath"

	// "github.com/gin-contrib/pprof"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/apache/apisix-dashboard/api/internal/conf"
	"github.com/apache/apisix-dashboard/api/internal/filter"
	"github.com/apache/apisix-dashboard/api/internal/handler"
	"github.com/apache/apisix-dashboard/api/internal/handler/authentication"
	"github.com/apache/apisix-dashboard/api/internal/handler/consumer"
	"github.com/apache/apisix-dashboard/api/internal/handler/data_loader"
	"github.com/apache/apisix-dashboard/api/internal/handler/global_rule"
	"github.com/apache/apisix-dashboard/api/internal/handler/healthz"
	"github.com/apache/apisix-dashboard/api/internal/handler/label"
	"github.com/apache/apisix-dashboard/api/internal/handler/plugin_config"
	"github.com/apache/apisix-dashboard/api/internal/handler/proto"
	"github.com/apache/apisix-dashboard/api/internal/handler/route"
	"github.com/apache/apisix-dashboard/api/internal/handler/schema"
	"github.com/apache/apisix-dashboard/api/internal/handler/server_info"
	"github.com/apache/apisix-dashboard/api/internal/handler/service"
	"github.com/apache/apisix-dashboard/api/internal/handler/ssl"
	"github.com/apache/apisix-dashboard/api/internal/handler/stream_route"
	"github.com/apache/apisix-dashboard/api/internal/handler/system_config"
	"github.com/apache/apisix-dashboard/api/internal/handler/tool"
	"github.com/apache/apisix-dashboard/api/internal/handler/upstream"
	"github.com/apache/apisix-dashboard/api/internal/log"
)

func SetUpRouter() *gin.Engine {
	if conf.ENV == conf.EnvLOCAL || conf.ENV == conf.EnvDEV {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	logger := log.GetLogger(log.AccessLog)
	// security
	r.Use(filter.RequestLogHandler(logger), filter.IPFilter(), filter.InvalidRequest(), filter.Authentication())

	// misc
	r.Use(gzip.Gzip(gzip.DefaultCompression), filter.CORS(), filter.RequestId(), filter.SchemaCheck(), filter.RecoverHandler())
	r.Use(static.Serve("/", static.LocalFile(filepath.Join(conf.WorkDir, conf.WebDir), false)))
	r.NoRoute(func(c *gin.Context) {
		c.File(fmt.Sprintf("%s/index.html", filepath.Join(conf.WorkDir, conf.WebDir)))
	})

	factories := []handler.RegisterFactory{
		route.NewHandler,
		ssl.NewHandler,
		consumer.NewHandler,
		upstream.NewHandler,
		service.NewHandler,
		schema.NewHandler,
		schema.NewSchemaHandler,
		healthz.NewHandler,
		authentication.NewHandler,
		global_rule.NewHandler,
		server_info.NewHandler,
		label.NewHandler,
		data_loader.NewHandler,
		data_loader.NewImportHandler,
		tool.NewHandler,
		plugin_config.NewHandler,
		proto.NewHandler,
		stream_route.NewHandler,
		system_config.NewHandler,
	}

	for i := range factories {
		h, err := factories[i]()
		if err != nil {
			panic(err)
		}
		h.ApplyRoute(r)
	}

	// pprof.Register(r)

	return r
}
