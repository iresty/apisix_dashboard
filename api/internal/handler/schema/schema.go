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
package schema

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/wrapper"
	wgin "github.com/shiningrush/droplet/wrapper/gin"

	"github.com/apache/apisix-dashboard/api/internal/conf"
	"github.com/apache/apisix-dashboard/api/internal/handler"
)

type SchemaHandler struct {
}

func NewSchemaHandler() (handler.RouteRegister, error) {
	return &SchemaHandler{}, nil
}

func (h *SchemaHandler) ApplyRoute(r *gin.Engine) {
	r.GET("/apisix/admin/schema/plugins/:name", wgin.Wraps(h.PluginSchema,
		wrapper.InputType(reflect.TypeOf(PluginSchemaInput{}))))

	r.GET("/apisix/admin/schemas/:resource", wgin.Wraps(h.Schema,
		wrapper.InputType(reflect.TypeOf(SchemaInput{}))))
}

type SchemaInput struct {
	Resource string `auto_read:"resource,path" validate:"required"`
}

func (h *SchemaHandler) Schema(c droplet.Context) (interface{}, error) {
	input := c.Input().(*SchemaInput)

	ret := conf.Schema.Get("main." + input.Resource).Value()

	if ret == nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusNotFound},
			fmt.Errorf("schema of %s not found", input.Resource)
	}

	return ret, nil
}

type PluginSchemaInput struct {
	Name       string `auto_read:"name,path" validate:"required"`
	SchemaType string `auto_read:"schema_type,query"`
}

func (h *SchemaHandler) PluginSchema(c droplet.Context) (interface{}, error) {
	input := c.Input().(*PluginSchemaInput)

	var ret interface{}
	if input.SchemaType == "consumer" {
		ret = conf.Schema.Get("plugins." + input.Name + ".consumer_schema").Value()
	}

	if ret == nil {
		ret = conf.Schema.Get("plugins." + input.Name + ".schema").Value()
	}

	if ret == nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusNotFound},
			fmt.Errorf("schema of plugins %s not found", input.Name)
	}

	return ret, nil
}
