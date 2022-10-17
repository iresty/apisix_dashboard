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

package handler

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/apache/apisix-dashboard/api/internal/config"
)

type RegisterFactory func() (RouteRegister, error)

type RouteRegister interface {
	ApplyRoute(r *gin.Engine, cfg config.Config)
}

type Func = func(c *gin.Context, input interface{}) Response
type Response struct {
	StatusCode int
	Success    bool
	Data       interface{}
	ErrMsg     string
}

func Wrap(f Func, inputType reflect.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input interface{}
		if inputType != nil {
			input = reflect.New(inputType).Interface()
			if err := c.ShouldBind(&input); err != nil {
				if _, ok := err.(validator.ValidationErrors); ok {
					c.AbortWithStatusJSON(http.StatusBadRequest,
						buildResponse(false, errors.Wrap(err, "input validate error").Error(), nil))
					return
				}
				c.AbortWithStatusJSON(http.StatusBadRequest,
					buildResponse(false, errors.Wrap(err, "input parse error").Error(), nil))
				return
			}
		}

		resp := f(c, input)

		c.JSON(resp.StatusCode, buildResponse(resp.Success, resp.ErrMsg, resp.Data))
	}
}

// buildResponse builds the response body for the Dashboard Resource API output.
func buildResponse(success bool, message string, data interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"success": success,
		"message": message,
	}
	if data != nil {
		resp["data"] = data
	}
	return resp
}
