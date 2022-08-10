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
package label

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/wrapper"
	wgin "github.com/shiningrush/droplet/wrapper/gin"

	"github.com/apache/apisix-dashboard/api/internal/core/entity"
	"github.com/apache/apisix-dashboard/api/internal/core/store"
	"github.com/apache/apisix-dashboard/api/internal/handler"
	"github.com/apache/apisix-dashboard/api/internal/utils"
)

type Handler struct {
	routeStore        store.Interface
	serviceStore      store.Interface
	upstreamStore     store.Interface
	sslStore          store.Interface
	consumerStore     store.Interface
	pluginConfigStore store.Interface
}

var _ json.Marshaler = Pair{}

type Pair struct {
	Key string
	Val string
}

func (p Pair) MarshalJSON() ([]byte, error) {
	res := fmt.Sprintf("{%s:%s}", strconv.Quote(p.Key), strconv.Quote(p.Val))
	return []byte(res), nil
}

func NewHandler() (handler.RouteRegister, error) {
	return &Handler{
		routeStore:        store.GetStore(store.HubKeyRoute),
		serviceStore:      store.GetStore(store.HubKeyService),
		upstreamStore:     store.GetStore(store.HubKeyUpstream),
		sslStore:          store.GetStore(store.HubKeySsl),
		consumerStore:     store.GetStore(store.HubKeyConsumer),
		pluginConfigStore: store.GetStore(store.HubKeyPluginConfig),
	}, nil
}

func (h *Handler) ApplyRoute(r *gin.Engine) {
	r.GET("/apisix/admin/labels/:type", wgin.Wraps(h.List,
		wrapper.InputType(reflect.TypeOf(ListInput{}))))
}

type ListInput struct {
	store.Pagination
	Type  string `auto_read:"type,path" validate:"required"`
	Label string `auto_read:"label,query"`
}

func subsetOf(reqLabels map[string]struct{}, labels map[string]string) map[string]string {
	if len(reqLabels) == 0 {
		return labels
	}

	var res = make(map[string]string)
	for k, v := range labels {
		if _, exist := reqLabels[k]; exist {
			res[k] = v
		}

		if _, exist := reqLabels[k+":"+v]; exist {
			res[k] = v
		}
	}

	return res
}

// swagger:operation GET /api/labels getLabelsList
//
// Return the labels list among `route,ssl,consumer,upstream,service`
// according to the specified page number and page size, and can search labels by label.
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
// - name: label
//   in: query
//   description: label filter of labels
//   required: false
//   type: string
// responses:
//   '0':
//     description: list response
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/service"
//   default:
//     description: unexpected error
//     schema:
//       "$ref": "#/definitions/ApiError"
func (h *Handler) List(c droplet.Context) (interface{}, error) {
	input := c.Input().(*ListInput)

	typ := input.Type
	reqLabels, err := utils.GenLabelMap(input.Label)
	if err != nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
			fmt.Errorf("%s: \"%s\"", err.Error(), input.Label)
	}

	var items []interface{}
	switch typ {
	case "route":
		items = append(items, h.routeStore)
	case "service":
		items = append(items, h.serviceStore)
	case "consumer":
		items = append(items, h.consumerStore)
	case "ssl":
		items = append(items, h.sslStore)
	case "upstream":
		items = append(items, h.upstreamStore)
	case "plugin_config":
		items = append(items, h.pluginConfigStore)
	case "all":
		items = append(items, h.routeStore, h.serviceStore, h.upstreamStore,
			h.sslStore, h.consumerStore, h.pluginConfigStore)
	}

	predicate := func(obj interface{}) bool {
		var ls map[string]string

		switch obj := obj.(type) {
		case *entity.Route:
			ls = obj.Labels
		case *entity.Consumer:
			ls = obj.Labels
		case *entity.SSL:
			ls = obj.Labels
		case *entity.Service:
			ls = obj.Labels
		case *entity.Upstream:
			ls = obj.Labels
		case *entity.PluginConfig:
			ls = obj.Labels
		default:
			return false
		}

		return utils.LabelContains(ls, reqLabels)
	}

	format := func(obj interface{}) interface{} {
		val := reflect.ValueOf(obj).Elem()
		l := val.FieldByName("Labels")
		if l.IsNil() {
			return nil
		}

		ls := l.Interface().(map[string]string)
		return subsetOf(reqLabels, ls)
	}

	var totalRet = store.NewListOutput()
	var existMap = make(map[string]struct{})
	for _, item := range items {
		ret, err := item.(store.Interface).List(c.Context(),
			store.ListInput{
				Predicate: predicate,
				Format:    format,
				// Sort it later.
				PageSize:   0,
				PageNumber: 0,
				Less: func(i, j interface{}) bool {
					return true
				},
			},
		)

		if err != nil {
			return nil, err
		}

		for _, r := range ret.Rows {
			if r == nil {
				continue
			}

			for k, v := range r.(map[string]string) {
				key := fmt.Sprintf("%s:%s", k, v)
				if _, exist := existMap[key]; exist {
					continue
				}

				existMap[key] = struct{}{}
				p := Pair{Key: k, Val: v}
				totalRet.Rows = append(totalRet.Rows, p)
			}
		}
	}
	totalRet.TotalSize = len(totalRet.Rows)

	sort.Slice(totalRet.Rows, func(i, j int) bool {
		p1 := totalRet.Rows[i].(Pair)
		p2 := totalRet.Rows[j].(Pair)

		if strings.Compare(p1.Key, p2.Key) == 0 {
			return strings.Compare(p1.Val, p2.Val) < 0
		}

		return strings.Compare(p1.Key, p2.Key) < 0
	})

	/* There are more than one store items,
	   So we need sort after getting all of labels.
	*/
	if input.PageSize > 0 && input.PageNumber > 0 {
		skipCount := (input.PageNumber - 1) * input.PageSize
		if skipCount > totalRet.TotalSize {
			totalRet.Rows = []interface{}{}
			return totalRet, nil
		}

		endIdx := skipCount + input.PageSize
		if endIdx >= totalRet.TotalSize {
			totalRet.Rows = totalRet.Rows[skipCount:]
			return totalRet, nil
		}

		totalRet.Rows = totalRet.Rows[skipCount:endIdx]
	}

	return totalRet, nil
}
