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
	"bytes"
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/juliangruber/go-intersect"
	"github.com/pkg/errors"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/wrapper"
	wgin "github.com/shiningrush/droplet/wrapper/gin"

	"github.com/apisix/manager-api/internal/conf"
	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/core/store"
	"github.com/apisix/manager-api/internal/handler"
	loader "github.com/apisix/manager-api/internal/handler/data_loader/loader"
	"github.com/apisix/manager-api/internal/handler/data_loader/loader/openapi3"
)

type ImportHandler struct {
	routeStore        store.Interface
	upstreamStore     store.Interface
	serviceStore      store.Interface
	consumerStore     store.Interface
	sslStore          store.Interface
	streamRouteStore  store.Interface
	globalPluginStore store.Interface
	pluginConfigStore store.Interface
	protoStore        store.Interface
}

func NewImportHandler() (handler.RouteRegister, error) {
	return &ImportHandler{
		routeStore:        store.GetStore(store.HubKeyRoute),
		upstreamStore:     store.GetStore(store.HubKeyUpstream),
		serviceStore:      store.GetStore(store.HubKeyService),
		consumerStore:     store.GetStore(store.HubKeyConsumer),
		sslStore:          store.GetStore(store.HubKeySsl),
		streamRouteStore:  store.GetStore(store.HubKeyStreamRoute),
		globalPluginStore: store.GetStore(store.HubKeyGlobalRule),
		pluginConfigStore: store.GetStore(store.HubKeyPluginConfig),
		protoStore:        store.GetStore(store.HubKeyProto),
	}, nil
}

func (h *ImportHandler) ApplyRoute(r *gin.Engine) {
	r.POST("/apisix/admin/import/routes", wgin.Wraps(h.Import,
		wrapper.InputType(reflect.TypeOf(ImportInput{}))))
}

type ImportResult struct {
	Total  int      `json:"total"`
	Failed int      `json:"failed"`
	Errors []string `json:"errors"`
}

type LoaderType string

const (
	LoaderTypeOpenAPI3 LoaderType = "openapi3"
)

var (
	storeList = []store.HubKey{store.HubKeyRoute, store.HubKeyUpstream, store.HubKeyService,
		store.HubKeyConsumer, store.HubKeySsl, store.HubKeyStreamRoute, store.HubKeyGlobalRule,
		store.HubKeyPluginConfig, store.HubKeyProto}
)

type ImportInput struct {
	Type        string `auto_read:"type"`
	TaskName    string `auto_read:"task_name"`
	FileName    string `auto_read:"_file"`
	FileContent []byte `auto_read:"file"`

	MergeMethod string `auto_read:"merge_method"`
}

func (h *ImportHandler) Import(c droplet.Context) (interface{}, error) {
	input := c.Input().(*ImportInput)

	// input file content check
	contentLen := bytes.Count(input.FileContent, nil) - 1
	if contentLen <= 0 {
		return nil, errors.New("uploaded file is empty")
	}
	if contentLen > conf.ImportSizeLimit {
		return nil, errors.Errorf("uploaded file size exceeds the limit, limit is %d", conf.ImportSizeLimit)
	}

	var l loader.Loader
	switch LoaderType(input.Type) {
	case LoaderTypeOpenAPI3:
		l = &openapi3.Loader{
			MergeMethod: input.MergeMethod == "true",
			TaskName:    input.TaskName,
		}
		break
	default:
		return nil, fmt.Errorf("unsupported data loader type: %s", input.Type)
	}

	dataSets, err := l.Import(input.FileContent)
	if err != nil {
		return nil, err
	}

	// Pre-checking for route duplication
	preCheckErrs := h.preCheck(c.Context(), dataSets)
	if len(preCheckErrs) > 0 {
		var errStr string
		for _, err := range preCheckErrs {
			errStr += err.Error() + "\n"
		}
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, errors.New(errStr)
	}

	// Create APISIX resources
	createErrs := h.createEntities(c.Context(), dataSets)
	return h.convertToImportResult(dataSets, createErrs), nil
}

// Pre-check imported data for duplicates
// The main problem facing duplication is routing, so here
// we mainly check the duplication of routes, based on
// domain name and uri.
func (h *ImportHandler) preCheck(ctx context.Context, data *loader.DataSets) []error {
	var errs = make([]error, 0)
	for _, route := range data.Routes {
		o, err := h.routeStore.List(ctx, store.ListInput{
			// The check logic here is that if when a duplicate HOST or URI
			// has been found, the HTTP method is checked for overlap, and
			// if there is overlap it is determined to be a duplicate route
			// and the import is rejected.
			Predicate: func(obj interface{}) bool {
				r := obj.(*entity.Route)
				methodCheck := func() bool {
					return len(intersect.Hash(r.Methods, route.Methods)) > 0
				}

				// Check URI duplication
				if r.URI == route.URI || len(intersect.Hash(r.Uris, route.Uris)) > 0 {
					return methodCheck()
				}
				// Check host duplication
				if r.Host == route.Host || len(intersect.Hash(r.Hosts, route.Hosts)) > 0 {
					return methodCheck()
				}
				return false
			},
			PageSize:   0,
			PageNumber: 0,
		})
		if err != nil {
			// When a special storage layer error occurs, return directly.
			return []error{err}
		}

		// Duplicate routes found
		if o.TotalSize > 0 {
			for _, row := range o.Rows {
				route, ok := row.(*entity.Route)
				if ok {
					errs = append(errs, fmt.Errorf("%s uri is duplicated with route %s", route.Uris[0], route.Name))
				}
			}
		}
	}

	return []error{}
}

// Create parsed resources
func (h *ImportHandler) createEntities(ctx context.Context, data *loader.DataSets) map[store.HubKey][]string {
	errs := make(map[store.HubKey][]string)

	for _, route := range data.Routes {
		_, err := h.routeStore.Create(ctx, &route)
		if err != nil {
			if errs[store.HubKeyRoute] == nil {
				errs[store.HubKeyRoute] = []string{}
			}
			errs[store.HubKeyRoute] = append(errs[store.HubKeyRoute], err.Error())
		}
	}
	for _, upstream := range data.Upstreams {
		_, err := h.upstreamStore.Create(ctx, &upstream)
		if err != nil {
			if errs[store.HubKeyUpstream] == nil {
				errs[store.HubKeyUpstream] = []string{}
			}
			errs[store.HubKeyUpstream] = append(errs[store.HubKeyUpstream], err.Error())
		}
	}
	for _, service := range data.Services {
		_, err := h.serviceStore.Create(ctx, &service)
		if err != nil {
			if errs[store.HubKeyService] == nil {
				errs[store.HubKeyService] = []string{}
			}
			errs[store.HubKeyService] = append(errs[store.HubKeyService], err.Error())
		}
	}
	for _, consumer := range data.Consumers {
		_, err := h.consumerStore.Create(ctx, &consumer)
		if err != nil {
			if errs[store.HubKeyConsumer] == nil {
				errs[store.HubKeyConsumer] = []string{}
			}
			errs[store.HubKeyConsumer] = append(errs[store.HubKeyConsumer], err.Error())
		}
	}
	for _, ssl := range data.SSLs {
		_, err := h.sslStore.Create(ctx, &ssl)
		if err != nil {
			if errs[store.HubKeySsl] == nil {
				errs[store.HubKeySsl] = []string{}
			}
			errs[store.HubKeySsl] = append(errs[store.HubKeySsl], err.Error())
		}
	}
	for _, route := range data.StreamRoutes {
		_, err := h.streamRouteStore.Create(ctx, &route)
		if err != nil {
			if errs[store.HubKeyStreamRoute] == nil {
				errs[store.HubKeyStreamRoute] = []string{}
			}
			errs[store.HubKeyStreamRoute] = append(errs[store.HubKeyStreamRoute], err.Error())
		}
	}
	for _, plugin := range data.GlobalPlugins {
		_, err := h.globalPluginStore.Create(ctx, &plugin)
		if err != nil {
			if errs[store.HubKeyGlobalRule] == nil {
				errs[store.HubKeyGlobalRule] = []string{}
			}
			errs[store.HubKeyGlobalRule] = append(errs[store.HubKeyGlobalRule], err.Error())
		}
	}
	for _, config := range data.PluginConfigs {
		_, err := h.pluginConfigStore.Create(ctx, &config)
		if err != nil {
			if errs[store.HubKeyPluginConfig] == nil {
				errs[store.HubKeyPluginConfig] = []string{}
			}
			errs[store.HubKeyPluginConfig] = append(errs[store.HubKeyPluginConfig], err.Error())
		}
	}
	for _, proto := range data.Protos {
		_, err := h.protoStore.Create(ctx, &proto)
		if err != nil {
			if errs[store.HubKeyProto] == nil {
				errs[store.HubKeyProto] = []string{}
			}
			errs[store.HubKeyProto] = append(errs[store.HubKeyProto], err.Error())
		}
	}

	return errs
}

// Convert import errors to response result
func (ImportHandler) convertToImportResult(data *loader.DataSets, errs map[store.HubKey][]string) map[store.HubKey]ImportResult {
	return map[store.HubKey]ImportResult{
		store.HubKeyRoute: {
			Total:  len(data.Routes),
			Failed: len(errs[store.HubKeyRoute]),
			Errors: errs[store.HubKeyRoute],
		},
		store.HubKeyUpstream: {
			Total:  len(data.Upstreams),
			Failed: len(errs[store.HubKeyUpstream]),
			Errors: errs[store.HubKeyUpstream],
		},
		store.HubKeyService: {
			Total:  len(data.Services),
			Failed: len(errs[store.HubKeyService]),
			Errors: errs[store.HubKeyService],
		},
		store.HubKeyConsumer: {
			Total:  len(data.Consumers),
			Failed: len(errs[store.HubKeyConsumer]),
			Errors: errs[store.HubKeyConsumer],
		},
		store.HubKeySsl: {
			Total:  len(data.SSLs),
			Failed: len(errs[store.HubKeySsl]),
			Errors: errs[store.HubKeySsl],
		},
		store.HubKeyStreamRoute: {
			Total:  len(data.StreamRoutes),
			Failed: len(errs[store.HubKeyStreamRoute]),
			Errors: errs[store.HubKeyStreamRoute],
		},
		store.HubKeyGlobalRule: {
			Total:  len(data.GlobalPlugins),
			Failed: len(errs[store.HubKeyGlobalRule]),
			Errors: errs[store.HubKeyGlobalRule],
		},
		store.HubKeyPluginConfig: {
			Total:  len(data.PluginConfigs),
			Failed: len(errs[store.HubKeyPluginConfig]),
			Errors: errs[store.HubKeyPluginConfig],
		},
		store.HubKeyProto: {
			Total:  len(data.Protos),
			Failed: len(errs[store.HubKeyProto]),
			Errors: errs[store.HubKeyProto],
		},
	}
}
