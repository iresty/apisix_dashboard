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
package data_loader

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/middleware"
	wgin "github.com/shiningrush/droplet/wrapper/gin"

	"github.com/apisix/manager-api/internal/conf"
	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/core/store"
	"github.com/apisix/manager-api/internal/handler"
	routeHandler "github.com/apisix/manager-api/internal/handler/route"
	"github.com/apisix/manager-api/internal/log"
	"github.com/apisix/manager-api/internal/utils"
	"github.com/apisix/manager-api/internal/utils/consts"
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
	r.POST("/apisix/admin/import", wgin.Wraps(h.Import))
}

type ImportInput struct {
	Force bool `auto_read:"force,query"`
}

func (h *Handler) Import(c droplet.Context) (interface{}, error) {
	httpReq := c.Get(middleware.KeyHttpRequest)
	if httpReq == nil {
		return nil, errors.New("input middleware cannot get http request")
	}
	req := httpReq.(*http.Request)
	req.Body = http.MaxBytesReader(nil, req.Body, int64(conf.ImportSizeLimit))
	if err := req.ParseMultipartForm(int64(conf.ImportSizeLimit)); err != nil {
		log.Warnf("upload file size exceeds limit: %s", err)
		return nil, fmt.Errorf("the file size exceeds the limit; limit %d", conf.ImportSizeLimit)
	}

	Force := req.URL.Query().Get("force")

	_, fileHeader, err := req.FormFile("file")
	if err != nil {
		return nil, err
	}

	// file check
	suffix := path.Ext(fileHeader.Filename)
	if suffix != ".json" && suffix != ".yaml" && suffix != ".yml" {
		return nil, fmt.Errorf("the file type error: %s", suffix)
	}

	// read file and parse
	handle, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = handle.Close()
	}()

	reader := bufio.NewReader(handle)
	bytes := make([]byte, fileHeader.Size)
	_, err = reader.Read(bytes)
	if err != nil {
		return nil, err
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(bytes)
	if err != nil {
		return nil, err
	}

	if len(swagger.Paths) < 1 {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
			errors.New("empty or invalid imported file")
	}

	routes, err := OpenAPI3ToRoute(swagger)
	if err != nil {
		return nil, err
	}

	routeStore := store.GetStore(store.HubKeyRoute)
	upstreamStore := store.GetStore(store.HubKeyUpstream)
	scriptStore := store.GetStore(store.HubKeyScript)

	//input := c.Input().(*ImportInput)

	// check route
	for _, route := range routes {
		err := checkRouteExist(c.Context(), route)
		if err != nil && Force != "1" {
			log.Warnf("import duplicate: %s, route: %#v", err, route)
			return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
		}
		if route.ServiceID != nil {
			_, err := routeStore.Get(c.Context(), utils.InterfaceToString(route.ServiceID))
			if err != nil {
				if err == data.ErrNotFound {
					return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
						fmt.Errorf("service id: %s not found", route.ServiceID)
				}
				return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
			}
		}
		if route.UpstreamID != nil {
			_, err := upstreamStore.Get(c.Context(), utils.InterfaceToString(route.UpstreamID))
			if err != nil {
				if err == data.ErrNotFound {
					return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
						fmt.Errorf("upstream id: %s not found", route.UpstreamID)
				}
				return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
			}
		}
		if route.Script != nil {
			if route.ID == "" {
				route.ID = utils.GetFlakeUidStr()
			}
			script := &entity.Script{
				ID:     utils.InterfaceToString(route.ID),
				Script: route.Script,
			}
			// to lua
			var err error
			route.Script, err = routeHandler.GenerateLuaCode(route.Script.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			// save original conf
			if _, err = scriptStore.Create(c.Context(), script); err != nil {
				return nil, err
			}
		}

		if _, err := routeStore.CreateCheck(route); err != nil {
			return handler.SpecCodeResponse(err), err
		}
	}

	// create route
	for _, route := range routes {
		if _, err := routeStore.Create(c.Context(), route); err != nil {
			return handler.SpecCodeResponse(err), err
		}
	}

	return nil, nil
}

func checkRouteExist(ctx context.Context, route *entity.Route) error {
	routeStore := store.GetStore(store.HubKeyRoute)
	ret, err := routeStore.List(ctx, store.ListInput{
		Predicate: func(obj interface{}) bool {
			id := utils.InterfaceToString(route.ID)
			item := obj.(*entity.Route)
			if id != "" && id != utils.InterfaceToString(item.ID) {
				return false
			}

			if !(item.Host == route.Host && item.URI == route.URI && utils.StringSliceEqual(item.Uris, route.Uris) &&
				utils.StringSliceEqual(item.RemoteAddrs, route.RemoteAddrs) && item.RemoteAddr == route.RemoteAddr &&
				utils.StringSliceEqual(item.Hosts, route.Hosts) && item.Priority == route.Priority &&
				utils.ValueEqual(item.Vars, route.Vars) && item.FilterFunc == route.FilterFunc) {
				return false
			}
			return true
		},
		PageSize:   0,
		PageNumber: 0,
	})
	if err != nil {
		return err
	}
	if len(ret.Rows) > 0 {
		return consts.InvalidParam("route is duplicate")
	}
	return nil
}

func parseExtension(val *openapi3.Operation) (*entity.Route, error) {
	routeMap := map[string]interface{}{}
	for key, val := range val.Extensions {
		if strings.HasPrefix(key, "x-apisix-") {
			routeMap[strings.TrimPrefix(key, "x-apisix-")] = val
		}
	}

	route := new(entity.Route)
	routeJson, err := json.Marshal(routeMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(routeJson, &route)
	if err != nil {
		return nil, err
	}

	return route, nil
}

type PathValue struct {
	Method string
	Value  *openapi3.Operation
}

func mergePathValue(key string, values []PathValue, swagger *openapi3.Swagger) (map[string]*entity.Route, error) {
	var parsed []PathValue
	var routes = map[string]*entity.Route{}
	for _, value := range values {
		value.Value.OperationID = strings.Replace(value.Value.OperationID, value.Method, "", 1)
		var eq = false
		for _, v := range parsed {
			if utils.ValueEqual(v.Value, value.Value) {
				eq = true
				if routes[v.Method].Methods == nil {
					routes[v.Method].Methods = []string{}
				}
				routes[v.Method].Methods = append(routes[v.Method].Methods, value.Method)
			}
		}
		// not equal to the previous ones
		if !eq {
			route, err := getRouteFromPaths(value.Method, key, value.Value, swagger)
			if err != nil {
				return nil, err
			}
			routes[value.Method] = route
			parsed = append(parsed, value)
		}
	}

	return routes, nil
}

func OpenAPI3ToRoute(swagger *openapi3.Swagger) ([]*entity.Route, error) {
	var routes []*entity.Route
	paths := swagger.Paths
	var upstream *entity.UpstreamDef
	var err error
	for k, v := range paths {
		upstream = &entity.UpstreamDef{}
		if up, ok := v.Extensions["x-apisix-upstream"]; ok {
			err = json.Unmarshal(up.(json.RawMessage), upstream)
			if err != nil {
				return nil, err
			}
		}

		var values []PathValue
		if v.Get != nil {
			value := PathValue{
				Method: http.MethodGet,
				Value:  v.Get,
			}
			values = append(values, value)
		}
		if v.Post != nil {
			value := PathValue{
				Method: http.MethodPost,
				Value:  v.Post,
			}
			values = append(values, value)
		}
		if v.Head != nil {
			value := PathValue{
				Method: http.MethodHead,
				Value:  v.Head,
			}
			values = append(values, value)
		}
		if v.Put != nil {
			value := PathValue{
				Method: http.MethodPut,
				Value:  v.Put,
			}
			values = append(values, value)
		}
		if v.Patch != nil {
			value := PathValue{
				Method: http.MethodPatch,
				Value:  v.Patch,
			}
			values = append(values, value)
		}
		if v.Delete != nil {
			value := PathValue{
				Method: http.MethodDelete,
				Value:  v.Delete,
			}
			values = append(values, value)
		}

		// merge same route
		tmp, err := mergePathValue(k, values, swagger)
		if err != nil {
			return nil, err
		}

		for _, route := range tmp {
			routes = append(routes, route)
		}
	}

	return routes, nil
}

func parseParameters(parameters openapi3.Parameters, plugins map[string]interface{}) {
	props := make(map[string]interface{})
	var required []string
	for _, v := range parameters {
		if v.Value.Schema != nil {
			v.Value.Schema.Value.Format = ""
			v.Value.Schema.Value.XML = nil
		}

		switch v.Value.In {
		case "header":
			if v.Value.Schema != nil && v.Value.Schema.Value != nil {
				props[v.Value.Name] = v.Value.Schema.Value
			}
			if v.Value.Required {
				required = append(required, v.Value.Name)
			}
		}
	}

	requestValidation := make(map[string]interface{})
	if rv, ok := plugins["request-validation"]; ok {
		requestValidation = rv.(map[string]interface{})
	}
	requestValidation["header_schema"] = &entity.RequestValidation{
		Type:       "object",
		Required:   required,
		Properties: props,
	}
	plugins["request-validation"] = requestValidation
}

func parseRequestBody(requestBody *openapi3.RequestBodyRef, swagger *openapi3.Swagger, plugins map[string]interface{}) {
	schema := requestBody.Value.Content
	requestValidation := make(map[string]interface{})
	if rv, ok := plugins["request-validation"]; ok {
		requestValidation = rv.(map[string]interface{})
	}
	for _, v := range schema {
		if v.Schema.Ref != "" {
			s := getParameters(v.Schema.Ref, &swagger.Components).Value
			requestValidation["body_schema"] = &entity.RequestValidation{
				Type:       s.Type,
				Required:   s.Required,
				Properties: s.Properties,
			}
			plugins["request-validation"] = requestValidation
		} else if v.Schema.Value != nil {
			if v.Schema.Value.Properties != nil {
				for k1, v1 := range v.Schema.Value.Properties {
					if v1.Ref != "" {
						s := getParameters(v1.Ref, &swagger.Components)
						v.Schema.Value.Properties[k1] = s
					}
					v1.Value.Format = ""
				}
				requestValidation["body_schema"] = &entity.RequestValidation{
					Type:       v.Schema.Value.Type,
					Required:   v.Schema.Value.Required,
					Properties: v.Schema.Value.Properties,
				}
				plugins["request-validation"] = requestValidation
			} else if v.Schema.Value.Items != nil {
				if v.Schema.Value.Items.Ref != "" {
					s := getParameters(v.Schema.Value.Items.Ref, &swagger.Components).Value
					requestValidation["body_schema"] = &entity.RequestValidation{
						Type:       s.Type,
						Required:   s.Required,
						Properties: s.Properties,
					}
					plugins["request-validation"] = requestValidation
				}
			} else {
				requestValidation["body_schema"] = &entity.RequestValidation{
					Type:       "object",
					Required:   []string{},
					Properties: v.Schema.Value.Properties,
				}
			}
		}
		plugins["request-validation"] = requestValidation
	}
}

func parseSecurity(security openapi3.SecurityRequirements, securitySchemes openapi3.SecuritySchemes, plugins map[string]interface{}) {
	// todo: import consumers
	for _, securities := range security {
		for name := range securities {
			if schema, ok := securitySchemes[name]; ok {
				value := schema.Value
				if value == nil {
					continue
				}

				// basic auth
				if value.Type == "http" && value.Scheme == "basic" {
					plugins["basic-auth"] = map[string]interface{}{}
					//username, ok := value.Extensions["username"]
					//if !ok {
					//	continue
					//}
					//password, ok := value.Extensions["password"]
					//if !ok {
					//	continue
					//}
					//plugins["basic-auth"] = map[string]interface{}{
					//	"username": username,
					//	"password": password,
					//}
					// jwt auth
				} else if value.Type == "http" && value.Scheme == "bearer" && value.BearerFormat == "JWT" {
					plugins["jwt-auth"] = map[string]interface{}{}
					//key, ok := value.Extensions["key"]
					//if !ok {
					//	continue
					//}
					//secret, ok := value.Extensions["secret"]
					//if !ok {
					//	continue
					//}
					//plugins["jwt-auth"] = map[string]interface{}{
					//	"key":    key,
					//	"secret": secret,
					//}
					// key auth
				} else if value.Type == "apiKey" {
					plugins["key-auth"] = map[string]interface{}{}
					//key, ok := value.Extensions["key"]
					//if !ok {
					//	continue
					//}
					//plugins["key-auth"] = map[string]interface{}{
					//	"key": key,
					//}
				}
			}
		}
	}
}

func getRouteFromPaths(method, key string, value *openapi3.Operation, swagger *openapi3.Swagger) (*entity.Route, error) {
	// transform /path/{var} to  /path/*
	reg := regexp.MustCompile(`{[\w.]*}`)
	foundStr := reg.FindString(key)
	if foundStr != "" {
		key = strings.Split(key, foundStr)[0] + "*"
	}

	route, err := parseExtension(value)
	if err != nil {
		return nil, err
	}

	route.URI = key
	route.Name = value.OperationID
	route.Desc = value.Summary
	route.Methods = []string{method}

	if route.Plugins == nil {
		route.Plugins = make(map[string]interface{})
	}

	if value.Parameters != nil {
		parseParameters(value.Parameters, route.Plugins)
	}

	if value.RequestBody != nil {
		parseRequestBody(value.RequestBody, swagger, route.Plugins)
	}

	if value.Security != nil && swagger.Components.SecuritySchemes != nil {
		parseSecurity(*value.Security, swagger.Components.SecuritySchemes, route.Plugins)
	}

	return route, nil
}

func getParameters(ref string, components *openapi3.Components) *openapi3.SchemaRef {
	schemaRef := &openapi3.SchemaRef{}
	arr := strings.Split(ref, "/")
	if arr[0] == "#" && arr[1] == "components" && arr[2] == "schemas" {
		schemaRef = components.Schemas[arr[3]]
		schemaRef.Value.XML = nil
		// traverse properties to find another ref
		for k, v := range schemaRef.Value.Properties {
			if v.Value != nil {
				v.Value.XML = nil
				v.Value.Format = ""
			}
			if v.Ref != "" {
				schemaRef.Value.Properties[k] = getParameters(v.Ref, components)
			} else if v.Value.Items != nil && v.Value.Items.Ref != "" {
				v.Value.Items = getParameters(v.Value.Items.Ref, components)
			} else if v.Value.Items != nil && v.Value.Items.Value != nil {
				v.Value.Items.Value.XML = nil
				v.Value.Items.Value.Format = ""
			}
		}
	}
	return schemaRef
}
