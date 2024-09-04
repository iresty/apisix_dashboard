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
package overview

import (
	"reflect"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/wrapper"
	wgin "github.com/shiningrush/droplet/wrapper/gin"

	"github.com/apisix/manager-api/internal/conf"
	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/core/store"
	"github.com/apisix/manager-api/internal/handler"
	"github.com/apisix/manager-api/internal/utils"
)

type Handler struct {
	routeStore      store.Interface
	upstreamStore   store.Interface
	serverInfoStore store.Interface
	serviceStore    store.Interface
	sslStore        store.Interface
}

func NewHandler() (handler.RouteRegister, error) {
	return &Handler{
		routeStore:      store.GetStore(store.HubKeyRoute),
		upstreamStore:   store.GetStore(store.HubKeyUpstream),
		serverInfoStore: store.GetStore(store.HubKeyServerInfo),
		serviceStore:    store.GetStore(store.HubKeyService),
		sslStore:        store.GetStore(store.HubKeySsl),
	}, nil
}

func (h *Handler) ApplyRoute(r *gin.Engine) {
	r.GET("/apisix/admin/overview", wgin.Wraps(h.List,
		wrapper.InputType(reflect.TypeOf(ListInput{}))))
}

type ListInput struct {
}

func (h *Handler) List(c droplet.Context) (interface{}, error) {
	routerCnt, onlineRouterCnt, err := h.GetRouteCnt(c)
	if err != nil {
		return nil, err
	}

	upstreamCnt, err := h.GetUpstreamCnt(c)
	if err != nil {
		return nil, err
	}

	serviceCnt, err := h.GetServiceCnt(c)
	if err != nil {
		return nil, err
	}

	certificateCnt, err := h.GetCertificateCnt(c)
	if err != nil {
		return nil, err
	}

	plugins, err := h.GetPlugins(c)
	if err != nil {
		return nil, err
	}

	dashboardVersion, err := h.GetApisixVersion()
	if err != nil {
		return nil, err
	}

	serverInfo, err := h.GetServerInfo(c)
	if err != nil {
		return nil, err
	}

	res := &entity.Overview{
		RouteCnt:         routerCnt,
		OnlineRouterCnt:  onlineRouterCnt,
		UpstreamCnt:      upstreamCnt,
		ServiceCnt:       serviceCnt,
		CertificateCnt:   certificateCnt,
		DashboardVersion: dashboardVersion,
		GatewayInfo:      serverInfo,
		Plugins:          plugins,
	}

	return res, nil
}

func (h *Handler) GetCertificateCnt(c droplet.Context) (int64, error) {
	var certificateCnt int64 = 0

	certificateRet, err := h.sslStore.List(c.Context(), store.ListInput{
		Predicate: func(obj interface{}) bool {
			return true
		},
		PageSize:   0,
		PageNumber: 0,
	})

	if err == nil {
		certificateCnt = int64(certificateRet.TotalSize)
	}

	return certificateCnt, err
}

func (h *Handler) GetPlugins(_ droplet.Context) ([]string, error) {

	plugins := conf.Schema.Get("plugins")

	var res []string
	list := plugins.Value().(map[string]interface{})
	for name := range list {
		res = append(res, name)
	}
	sort.Strings(res)

	return res, nil
}

func (h *Handler) GetServerInfo(c droplet.Context) ([]*entity.ServerInfo, error) {
	ret, err := h.serverInfoStore.List(c.Context(), store.ListInput{
		Predicate: func(obj interface{}) bool {
			return true
		},
		PageSize:   0,
		PageNumber: 0,
	})

	if err != nil {
		return nil, err
	}
	var res [](*entity.ServerInfo)

	for _, item := range ret.Rows {
		res = append(res, item.(*entity.ServerInfo))
	}

	return res, nil
}

func (h *Handler) GetApisixVersion() (string, error) {
	_, version := utils.GetHashAndVersion()
	return version, nil
}

func (h *Handler) GetRouteCnt(c droplet.Context) (int64, int64, error) {
	var routesCnt, onlineRoutesCnt int64 = 0, 0

	routesRet, err := h.routeStore.List(c.Context(), store.ListInput{
		Predicate: func(obj interface{}) bool {
			// get all routes
			return true
		},
		PageSize:   0,
		PageNumber: 0,
	})

	if err == nil {
		routesCnt = int64(routesRet.TotalSize)
	}

	onlineRoutesRet, err := h.routeStore.List(c.Context(), store.ListInput{
		Predicate: func(obj interface{}) bool {
			// get all the online routes
			if int(obj.(*entity.Route).Status) == 1 {
				return true
			} else {
				return false
			}
		},
		PageSize:   0,
		PageNumber: 0,
	})

	if err == nil {
		onlineRoutesCnt = int64(onlineRoutesRet.TotalSize)
	}

	return routesCnt, onlineRoutesCnt, err
}

func (h *Handler) GetUpstreamCnt(c droplet.Context) (int64, error) {
	var upstreamCnt int64 = 0
	upstreamRet, err := h.upstreamStore.List(c.Context(), store.ListInput{
		Predicate: func(obj interface{}) bool {
			return true
		},
		Format: func(obj interface{}) interface{} {
			upstream := obj.(*entity.Upstream)
			upstream.Nodes = entity.NodesFormat(upstream.Nodes)
			return upstream
		},
		PageSize:   0,
		PageNumber: 0,
	})

	if err == nil {
		upstreamCnt = int64(upstreamRet.TotalSize)
	}

	return upstreamCnt, err
}

func (h *Handler) GetServiceCnt(c droplet.Context) (int64, error) {
	var serviceCnt int64 = 0

	servicesRet, err := h.serviceStore.List(c.Context(), store.ListInput{
		Predicate: func(obj interface{}) bool {
			return true
		},
		Format: func(obj interface{}) interface{} {
			service := obj.(*entity.Service)
			if service.Upstream != nil && service.Upstream.Nodes != nil {
				service.Upstream.Nodes = entity.NodesFormat(service.Upstream.Nodes)
			}
			return service
		},
		PageSize:   0,
		PageNumber: 0,
	})

	if err == nil {
		serviceCnt = int64(servicesRet.TotalSize)
	}

	return serviceCnt, err
}
