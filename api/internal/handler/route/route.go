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
package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/wrapper"
	wgin "github.com/shiningrush/droplet/wrapper/gin"
	lua "github.com/yuin/gopher-lua"

	"github.com/apache/apisix-dashboard/api/internal/core/entity"
	"github.com/apache/apisix-dashboard/api/internal/core/store"
	"github.com/apache/apisix-dashboard/api/internal/handler"
	"github.com/apache/apisix-dashboard/api/internal/log"
	"github.com/apache/apisix-dashboard/api/internal/utils"
	"github.com/apache/apisix-dashboard/api/internal/utils/consts"
)

type Handler struct {
	routeStore    store.Interface
	svcStore      store.Interface
	upstreamStore store.Interface
	scriptStore   store.Interface
}

func NewHandler() (handler.RouteRegister, error) {
	return &Handler{
		routeStore:    store.GetStore(store.HubKeyRoute),
		svcStore:      store.GetStore(store.HubKeyService),
		upstreamStore: store.GetStore(store.HubKeyUpstream),
		scriptStore:   store.GetStore(store.HubKeyScript),
	}, nil
}

func (h *Handler) ApplyRoute(r *gin.Engine) {
	r.GET("/apisix/admin/routes/:id", wgin.Wraps(h.Get,
		wrapper.InputType(reflect.TypeOf(GetInput{}))))
	r.GET("/apisix/admin/routes", wgin.Wraps(h.List,
		wrapper.InputType(reflect.TypeOf(ListInput{}))))
	r.POST("/apisix/admin/routes", wgin.Wraps(h.Create,
		wrapper.InputType(reflect.TypeOf(entity.Route{}))))
	r.PUT("/apisix/admin/routes", wgin.Wraps(h.Update,
		wrapper.InputType(reflect.TypeOf(UpdateInput{}))))
	r.PUT("/apisix/admin/routes/:id", wgin.Wraps(h.Update,
		wrapper.InputType(reflect.TypeOf(UpdateInput{}))))

	r.DELETE("/apisix/admin/routes/:ids", wgin.Wraps(h.BatchDelete,
		wrapper.InputType(reflect.TypeOf(BatchDelete{}))))

	r.PATCH("/apisix/admin/routes/:id", wgin.Wraps(h.Patch,
		wrapper.InputType(reflect.TypeOf(PatchInput{}))))
	r.PATCH("/apisix/admin/routes/:id/*path", wgin.Wraps(h.Patch,
		wrapper.InputType(reflect.TypeOf(PatchInput{}))))

	r.GET("/apisix/admin/notexist/routes", wgin.Wraps(h.Exist,
		wrapper.InputType(reflect.TypeOf(ExistCheckInput{}))))
}

type PatchInput struct {
	ID      string `auto_read:"id,path"`
	SubPath string `auto_read:"path,path"`
	Body    []byte `auto_read:"@body"`
}

func (h *Handler) Patch(c droplet.Context) (interface{}, error) {
	input := c.Input().(*PatchInput)
	reqBody := input.Body
	ID := input.ID
	subPath := input.SubPath

	stored, err := h.routeStore.Get(c.Context(), ID)
	if err != nil {
		return handler.SpecCodeResponse(err), err
	}

	res, err := utils.MergePatch(stored, subPath, reqBody)
	if err != nil {
		return handler.SpecCodeResponse(err), err
	}

	var route entity.Route
	err = json.Unmarshal(res, &route)
	if err != nil {
		return handler.SpecCodeResponse(err), err
	}

	ret, err := h.routeStore.Update(c.Context(), &route, false)
	if err != nil {
		return handler.SpecCodeResponse(err), err
	}

	return ret, nil
}

type GetInput struct {
	ID string `auto_read:"id,path" validate:"required"`
}

// swagger:operation GET /apisix/admin/routes getRouteList
//
// Return the route list according to the specified page number and page size, and can search routes by name and uri.
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
// - name: name
//   in: query
//   description: name of route
//   required: false
//   type: string
// - name: uri
//   in: query
//   description: uri of route
//   required: false
//   type: string
// - name: label
//   in: query
//   description: label of route
//   required: false
//   type: string
// responses:
//   '0':
//     description: list response
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/route"
//   default:
//     description: unexpected error
//     schema:
//       "$ref": "#/definitions/ApiError"
func (h *Handler) Get(c droplet.Context) (interface{}, error) {
	input := c.Input().(*GetInput)

	r, err := h.routeStore.Get(c.Context(), input.ID)
	if err != nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusNotFound}, err
	}

	//format respond
	route := r.(*entity.Route)
	script, _ := h.scriptStore.Get(c.Context(), input.ID)
	if script != nil {
		route.Script = script.(*entity.Script).Script
	}

	//format
	if route.Upstream != nil && route.Upstream.Nodes != nil {
		route.Upstream.Nodes = entity.NodesFormat(route.Upstream.Nodes)
	}

	return route, nil
}

type ListInput struct {
	Name   string `auto_read:"name,query"`
	URI    string `auto_read:"uri,query"`
	Label  string `auto_read:"label,query"`
	Status string `auto_read:"status,query"`
	store.Pagination
}

func uriContains(obj *entity.Route, uri string) bool {
	if strings.Contains(obj.URI, uri) {
		return true
	}
	for _, str := range obj.Uris {
		result := strings.Contains(str, uri)
		if result {
			return true
		}
	}
	return false
}

func (h *Handler) List(c droplet.Context) (interface{}, error) {
	input := c.Input().(*ListInput)
	labelMap, err := utils.GenLabelMap(input.Label)
	if err != nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
			fmt.Errorf("%s: \"%s\"", err.Error(), input.Label)
	}

	ret, err := h.routeStore.List(c.Context(), store.ListInput{
		Predicate: func(obj interface{}) bool {
			if input.Name != "" && !strings.Contains(obj.(*entity.Route).Name, input.Name) {
				return false
			}

			if input.URI != "" && !uriContains(obj.(*entity.Route), input.URI) {
				return false
			}

			if input.Label != "" && !utils.LabelContains(obj.(*entity.Route).Labels, labelMap) {
				return false
			}

			if input.Status != "" && strconv.Itoa(int(obj.(*entity.Route).Status)) != input.Status {
				return false
			}

			return true
		},
		Format: func(obj interface{}) interface{} {
			route := obj.(*entity.Route)
			if route.Upstream != nil && route.Upstream.Nodes != nil {
				route.Upstream.Nodes = entity.NodesFormat(route.Upstream.Nodes)
			}
			return route
		},
		PageSize:   input.PageSize,
		PageNumber: input.PageNumber,
	})

	if err != nil {
		return nil, err
	}

	//format respond
	var route *entity.Route
	for i, item := range ret.Rows {
		route = item.(*entity.Route)
		id := utils.InterfaceToString(route.ID)
		script, _ := h.scriptStore.Get(c.Context(), id)
		if script != nil {
			route.Script = script.(*entity.Script).Script
		}
		ret.Rows[i] = route
	}

	return ret, nil
}

func generateLuaCode(script map[string]interface{}) (string, error) {
	scriptString, err := json.Marshal(script)
	if err != nil {
		return "", err
	}
	workDir, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}
	libDir := filepath.Join(".", "dag-to-lua/")
	if err := os.Chdir(libDir); err != nil {
		log.Errorf("Chdir to libDir failed: %s", err)
		return "", err
	}

	defer func() {
		if err := os.Chdir(workDir); err != nil {
			log.Errorf("Chdir to workDir failed: %s", err)
		}
	}()

	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(`
	        local dag_to_lua = require 'dag-to-lua'
		local conf = [==[` + string(scriptString) + `]==]
	        code = dag_to_lua.generate(conf)
        `); err != nil {
		return "", err
	}

	code := L.GetGlobal("code")

	return code.String(), nil
}

func (h *Handler) Create(c droplet.Context) (interface{}, error) {
	input := c.Input().(*entity.Route)
	//check depend
	if input.ServiceID != nil {
		serviceID := utils.InterfaceToString(input.ServiceID)
		_, err := h.svcStore.Get(c.Context(), serviceID)
		if err != nil {
			if err == data.ErrNotFound {
				return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
					fmt.Errorf(consts.IDNotFound, "service", input.ServiceID)
			}
			return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
		}
	}
	if input.UpstreamID != nil {
		upstreamID := utils.InterfaceToString(input.UpstreamID)
		_, err := h.upstreamStore.Get(c.Context(), upstreamID)
		if err != nil {
			if err == data.ErrNotFound {
				return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
					fmt.Errorf("upstream id: %s not found", input.UpstreamID)
			}
			return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
		}
	}

	// If route's script_id is set, it must be equals to the route's id.
	if input.ScriptID != nil && (utils.InterfaceToString(input.ID) != utils.InterfaceToString(input.ScriptID)) {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
			fmt.Errorf("script_id must be the same as id")
	}

	if input.Script != nil {
		if utils.InterfaceToString(input.ID) == "" {
			input.ID = utils.GetFlakeUidStr()
		}
		script := &entity.Script{}
		script.ID = utils.InterfaceToString(input.ID)
		script.Script = input.Script

		var err error
		// Explicitly to lua if input script is of the map type, otherwise
		// it will always represent a piece of lua code of the string type.
		if scriptConf, ok := input.Script.(map[string]interface{}); ok {
			// For lua code of map type, syntax validation is done by
			// the generateLuaCode function
			input.Script, err = generateLuaCode(scriptConf)
			if err != nil {
				return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
			}
		} else {
			// For lua code of string type, use utility func to syntax validation
			err = utils.ValidateLuaCode(input.Script.(string))
			if err != nil {
				return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
			}
		}

		//save original conf
		if _, err = h.scriptStore.Create(c.Context(), script); err != nil {
			return nil, err
		}

		// After saving the Script entity, always set route's script_id
		// the same as route's id.
		input.ScriptID = input.ID
	} else {
		// If script is unset, script_id must be unset neither.
		if input.ScriptID != nil {
			return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
				fmt.Errorf("script_id cannot be set if script is unset")
		}
	}

	// check name existed
	ret, err := handler.NameExistCheck(c.Context(), h.routeStore, "route", input.Name, nil)
	if err != nil {
		return ret, err
	}

	// create
	res, err := h.routeStore.Create(c.Context(), input)
	if err != nil {
		return handler.SpecCodeResponse(err), err
	}

	return res, nil
}

type UpdateInput struct {
	ID string `auto_read:"id,path"`
	entity.Route
}

func (h *Handler) Update(c droplet.Context) (interface{}, error) {
	input := c.Input().(*UpdateInput)

	// check if ID in body is equal ID in path
	if err := handler.IDCompare(input.ID, input.Route.ID); err != nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
	}

	// if has id in path, use it
	if input.ID != "" {
		input.Route.ID = input.ID
	}

	//check depend
	if input.ServiceID != nil {
		serviceID := utils.InterfaceToString(input.ServiceID)
		_, err := h.svcStore.Get(c.Context(), serviceID)
		if err != nil {
			if err == data.ErrNotFound {
				return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
					fmt.Errorf(consts.IDNotFound, "service", input.ServiceID)
			}
			return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
		}
	}
	if input.UpstreamID != nil {
		upstreamID := utils.InterfaceToString(input.UpstreamID)
		_, err := h.upstreamStore.Get(c.Context(), upstreamID)
		if err != nil {
			if err == data.ErrNotFound {
				return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
					fmt.Errorf("upstream id: %s not found", input.UpstreamID)
			}
			return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
		}
	}

	// If route's script_id is set, it must be equals to the route's id.
	if input.Route.ScriptID != nil && (utils.InterfaceToString(input.ID) != utils.InterfaceToString(input.Route.ScriptID)) {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
			fmt.Errorf("script_id must be the same as id")
	}

	if input.Script != nil {
		script := &entity.Script{}
		script.ID = input.ID
		script.Script = input.Script

		var err error
		// Explicitly to lua if input script is of the map type, otherwise
		// it will always represent a piece of lua code of the string type.
		if scriptConf, ok := input.Script.(map[string]interface{}); ok {
			// For lua code of map type, syntax validation is done by
			// the generateLuaCode function
			input.Route.Script, err = generateLuaCode(scriptConf)
			if err != nil {
				return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
			}
		} else {
			// For lua code of string type, use utility func to syntax validation
			err = utils.ValidateLuaCode(input.Script.(string))
			if err != nil {
				return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
			}
		}

		//save original conf
		if _, err = h.scriptStore.Update(c.Context(), script, true); err != nil {
			//if not exists, create
			if err.Error() == fmt.Sprintf("key: %s is not found", script.ID) {
				if _, err := h.scriptStore.Create(c.Context(), script); err != nil {
					return handler.SpecCodeResponse(err), err
				}
			} else {
				return handler.SpecCodeResponse(err), err
			}
		}

		// After updating the Script entity, always set route's script_id
		// the same as route's id.
		input.Route.ScriptID = input.ID
	} else {
		// If script is unset, script_id must be unset neither.
		if input.Route.ScriptID != nil {
			return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
				fmt.Errorf("script_id cannot be set if script is unset")
		}
		//remove exists script
		id := utils.InterfaceToString(input.Route.ID)
		script, _ := h.scriptStore.Get(c.Context(), id)
		if script != nil {
			if err := h.scriptStore.BatchDelete(c.Context(), strings.Split(id, ",")); err != nil {
				log.Warnf("delete script %s failed", input.Route.ID)
			}
		}
	}

	// check name existed
	ret, err := handler.NameExistCheck(c.Context(), h.routeStore, "route", input.Name, input.ID)
	if err != nil {
		return ret, err
	}

	// create
	res, err := h.routeStore.Update(c.Context(), &input.Route, true)
	if err != nil {
		return handler.SpecCodeResponse(err), err
	}

	return res, nil
}

type BatchDelete struct {
	IDs string `auto_read:"ids,path"`
}

func (h *Handler) BatchDelete(c droplet.Context) (interface{}, error) {
	input := c.Input().(*BatchDelete)

	//delete route
	if err := h.routeStore.BatchDelete(c.Context(), strings.Split(input.IDs, ",")); err != nil {
		return handler.SpecCodeResponse(err), err
	}

	//delete stored script
	if err := h.scriptStore.BatchDelete(c.Context(), strings.Split(input.IDs, ",")); err != nil {
		//try again
		log.Warn("try to delete script %s again", input.IDs)
		if err := h.scriptStore.BatchDelete(c.Context(), strings.Split(input.IDs, ",")); err != nil {
			return nil, nil
		}
	}

	return nil, nil
}

type ExistCheckInput struct {
	Name    string `auto_read:"name,query"`
	Exclude string `auto_read:"exclude,query"`
}

// swagger:operation GET /apisix/admin/notexist/routes checkRouteExist
//
// Return result of route exists checking by name and exclude id.
//
// ---
// produces:
// - application/json
// parameters:
// - name: name
//   in: query
//   description: name of route
//   required: false
//   type: string
// - name: exclude
//   in: query
//   description: id of route that exclude checking
//   required: false
//   type: string
// responses:
//   '0':
//     description: route not exists
//     schema:
//       "$ref": "#/definitions/ApiError"
//   default:
//     description: unexpected error
//     schema:
//       "$ref": "#/definitions/ApiError"
func (h *Handler) Exist(c droplet.Context) (interface{}, error) {
	input := c.Input().(*ExistCheckInput)
	name := input.Name
	exclude := input.Exclude

	ret, err := h.routeStore.List(c.Context(), store.ListInput{
		Predicate: func(obj interface{}) bool {
			r := obj.(*entity.Route)
			if r.Name == name && r.ID != exclude {
				return true
			}

			return false
		},
		PageSize:   0,
		PageNumber: 0,
	})
	if err != nil {
		return nil, err
	}

	if ret.TotalSize > 0 {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
			consts.InvalidParam("Route name is reduplicate")
	}

	return nil, nil
}
