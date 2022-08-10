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
// API doc of Manager API.
//
// Manager API directly operates ETCD and provides data management for Apache APISIX, provides APIs for Front-end or other clients.
//
// Terms Of Service:
//     Schemes: http, https
//     Host: 127.0.0.1
//     License: Apache License 2.0 http://www.apache.org/licenses/LICENSE-2.0
//
//     Consumes:
//     - application/json
//     - application/xml
//
//     Produces:
//     - application/json
//     - application/xml
//
// swagger:meta
package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/middleware"

	"github.com/apache/apisix-dashboard/api/internal/core/entity"
	"github.com/apache/apisix-dashboard/api/internal/core/store"
	"github.com/apache/apisix-dashboard/api/internal/utils"
)

type RegisterFactory func() (RouteRegister, error)

type RouteRegister interface {
	ApplyRoute(r *gin.Engine)
}

func SpecCodeResponse(err error) *data.SpecCodeResponse {
	errMsg := err.Error()
	if strings.Contains(errMsg, "required") ||
		strings.Contains(errMsg, "conflicted") ||
		strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "missing") ||
		strings.Contains(errMsg, "schema validate failed") {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}
	}

	if strings.Contains(errMsg, "not found") {
		return &data.SpecCodeResponse{StatusCode: http.StatusNotFound}
	}

	return &data.SpecCodeResponse{StatusCode: http.StatusInternalServerError}
}

type ErrorTransformMiddleware struct {
	middleware.BaseMiddleware
}

func (mw *ErrorTransformMiddleware) Handle(ctx droplet.Context) error {
	if err := mw.BaseMiddleware.Handle(ctx); err != nil {
		bErr, ok := err.(*data.BaseError)
		if !ok {
			return err
		}
		switch bErr.Code {
		case data.ErrCodeValidate, data.ErrCodeFormat:
			ctx.SetOutput(&data.SpecCodeResponse{StatusCode: http.StatusBadRequest})
		case data.ErrCodeInternal:
			ctx.SetOutput(&data.SpecCodeResponse{StatusCode: http.StatusInternalServerError})
		}
		return err
	}
	return nil
}

func IDCompare(idOnPath string, idOnBody interface{}) error {
	idOnBodyStr, ok := idOnBody.(string)
	if !ok {
		idOnBodyStr = utils.InterfaceToString(idOnBody)
	}

	// check if id on path is == to id on body ONLY if both ids are valid
	if idOnPath != "" && idOnBodyStr != "" && idOnBodyStr != idOnPath {
		return fmt.Errorf("ID on path (%s) doesn't match ID on body (%s)", idOnPath, idOnBodyStr)
	}

	return nil
}

func NameExistCheck(ctx context.Context, stg store.Interface, resource, name string, excludeID interface{}) (interface{}, error) {
	ret, err := stg.List(ctx, store.ListInput{
		Predicate: func(obj interface{}) bool {
			var objName string
			var objID interface{}
			switch resource {
			case "route":
				objID = obj.(*entity.Route).ID
				objName = obj.(*entity.Route).Name
			case "service":
				objID = obj.(*entity.Service).ID
				objName = obj.(*entity.Service).Name
			case "upstream":
				objID = obj.(*entity.Upstream).ID
				objName = obj.(*entity.Upstream).Name
			default:
				panic("bad resource")
			}

			if excludeID != nil && objID == excludeID {
				return false
			}
			if objName == name {
				return true
			}

			return false
		},
	})
	if err != nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusInternalServerError}, err
	}
	if ret.TotalSize > 0 {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
			fmt.Errorf("%s name exists", resource)
	}

	return nil, nil
}
