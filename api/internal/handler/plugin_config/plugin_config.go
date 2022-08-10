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
package plugin_config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/wrapper"
	wgin "github.com/shiningrush/droplet/wrapper/gin"

	"github.com/apache/apisix-dashboard/api/internal/core/entity"
	"github.com/apache/apisix-dashboard/api/internal/core/store"
	"github.com/apache/apisix-dashboard/api/internal/handler"
	"github.com/apache/apisix-dashboard/api/internal/log"
	"github.com/apache/apisix-dashboard/api/internal/utils"
)

type Handler struct {
	pluginConfigStore store.Interface
	routeStore        store.Interface
}

func NewHandler() (handler.RouteRegister, error) {
	return &Handler{
		pluginConfigStore: store.GetStore(store.HubKeyPluginConfig),
		routeStore:        store.GetStore(store.HubKeyRoute),
	}, nil
}

func (h *Handler) ApplyRoute(r *gin.Engine) {
	r.GET("/apisix/admin/plugin_configs/:id", wgin.Wraps(h.Get,
		wrapper.InputType(reflect.TypeOf(GetInput{}))))
	r.GET("/apisix/admin/plugin_configs", wgin.Wraps(h.List,
		wrapper.InputType(reflect.TypeOf(ListInput{}))))
	r.POST("/apisix/admin/plugin_configs", wgin.Wraps(h.Create,
		wrapper.InputType(reflect.TypeOf(entity.PluginConfig{}))))
	r.PUT("/apisix/admin/plugin_configs", wgin.Wraps(h.Update,
		wrapper.InputType(reflect.TypeOf(UpdateInput{}))))
	r.PUT("/apisix/admin/plugin_configs/:id", wgin.Wraps(h.Update,
		wrapper.InputType(reflect.TypeOf(UpdateInput{}))))
	r.PATCH("/apisix/admin/plugin_configs/:id", wgin.Wraps(h.Patch,
		wrapper.InputType(reflect.TypeOf(PatchInput{}))))
	r.PATCH("/apisix/admin/plugin_configs/:id/*path", wgin.Wraps(h.Patch,
		wrapper.InputType(reflect.TypeOf(PatchInput{}))))
	r.DELETE("/apisix/admin/plugin_configs/:ids", wgin.Wraps(h.BatchDelete,
		wrapper.InputType(reflect.TypeOf(BatchDelete{}))))
}

type GetInput struct {
	ID string `auto_read:"id,path" validate:"required"`
}

func (h *Handler) Get(c droplet.Context) (interface{}, error) {
	input := c.Input().(*GetInput)

	pluginConfig, err := h.pluginConfigStore.Get(c.Context(), input.ID)
	if err != nil {
		return handler.SpecCodeResponse(err), err
	}

	return pluginConfig, nil
}

type ListInput struct {
	Search string `auto_read:"search,query"`
	Label  string `auto_read:"label,query"`
	store.Pagination
}

// swagger:operation GET /apisix/admin/plugin_configs getPluginConfigList
//
// Return the plugin_config list according to the specified page number and page size, and support search.
//
// ---
// produces:
// - application/json
// parameters:
// - name: page
//   in: query
//   description: page number
//   required: false
//   type: integer
// - name: page_size
//   in: query
//   description: page size
//   required: false
//   type: integer
// - name: search
//   in: query
//   description: search keyword
//   required: false
//   type: string
// responses:
//   '0':
//     description: list response
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/pluginConfig"
//   default:
//     description: unexpected error
//     schema:
//       "$ref": "#/definitions/ApiError"
func (h *Handler) List(c droplet.Context) (interface{}, error) {
	input := c.Input().(*ListInput)
	labelMap, err := utils.GenLabelMap(input.Label)
	if err != nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
			fmt.Errorf("%s: \"%s\"", err.Error(), input.Label)
	}

	ret, err := h.pluginConfigStore.List(c.Context(), store.ListInput{
		Predicate: func(obj interface{}) bool {
			if input.Search != "" {
				return strings.Contains(obj.(*entity.PluginConfig).Desc, input.Search)
			}

			if input.Label != "" && !utils.LabelContains(obj.(*entity.PluginConfig).Labels, labelMap) {
				return false
			}

			return true
		},
		PageSize:   input.PageSize,
		PageNumber: input.PageNumber,
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (h *Handler) Create(c droplet.Context) (interface{}, error) {
	input := c.Input().(*entity.PluginConfig)

	ret, err := h.pluginConfigStore.Create(c.Context(), input)
	if err != nil {
		return handler.SpecCodeResponse(err), err
	}

	return ret, nil
}

type UpdateInput struct {
	ID string `auto_read:"id,path"`
	entity.PluginConfig
}

func (h *Handler) Update(c droplet.Context) (interface{}, error) {
	input := c.Input().(*UpdateInput)

	// check if ID in body is equal ID in path
	if err := handler.IDCompare(input.ID, input.PluginConfig.ID); err != nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
	}

	if input.ID != "" {
		input.PluginConfig.ID = input.ID
	}

	ret, err := h.pluginConfigStore.Update(c.Context(), &input.PluginConfig, true)
	if err != nil {
		return handler.SpecCodeResponse(err), err
	}

	return ret, nil
}

type BatchDelete struct {
	IDs string `auto_read:"ids,path"`
}

func (h *Handler) BatchDelete(c droplet.Context) (interface{}, error) {
	input := c.Input().(*BatchDelete)

	IDs := strings.Split(input.IDs, ",")
	IDMap := map[string]bool{}
	for _, id := range IDs {
		IDMap[id] = true
	}
	ret, err := h.routeStore.List(c.Context(), store.ListInput{
		Predicate: func(obj interface{}) bool {
			id := utils.InterfaceToString(obj.(*entity.Route).PluginConfigID)
			if _, ok := IDMap[id]; ok {
				return true
			}
			return false
		},
	})

	if err != nil {
		return nil, err
	}

	if len(ret.Rows) > 0 {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
			fmt.Errorf("please disconnect the route (ID: %s) with this plugin config first",
				ret.Rows[0].(*entity.Route).ID)
	}

	if err := h.pluginConfigStore.BatchDelete(c.Context(), strings.Split(input.IDs, ",")); err != nil {
		return handler.SpecCodeResponse(err), err
	}

	return nil, nil
}

type PatchInput struct {
	ID      string `auto_read:"id,path"`
	SubPath string `auto_read:"path,path"`
	Body    []byte `auto_read:"@body"`
}

func (h *Handler) Patch(c droplet.Context) (interface{}, error) {
	input := c.Input().(*PatchInput)
	reqBody := input.Body
	id := input.ID
	subPath := input.SubPath

	stored, err := h.pluginConfigStore.Get(c.Context(), id)
	if err != nil {
		log.Warnf("get stored data from etcd failed: %s", err)
		return handler.SpecCodeResponse(err), err
	}

	res, err := utils.MergePatch(stored, subPath, reqBody)
	if err != nil {
		log.Warnf("merge failed: %s", err)
		return handler.SpecCodeResponse(err), err
	}

	var pluginConfig entity.PluginConfig
	if err := json.Unmarshal(res, &pluginConfig); err != nil {
		log.Warnf("unmarshal to pluginConfig failed: %s", err)
		return handler.SpecCodeResponse(err), err
	}

	ret, err := h.pluginConfigStore.Update(c.Context(), &pluginConfig, false)
	if err != nil {
		log.Warnf("update failed: %s", err)
		return handler.SpecCodeResponse(err), err
	}

	return ret, nil
}
