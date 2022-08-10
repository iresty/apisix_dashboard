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
package stream_route

import (
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
	"github.com/apache/apisix-dashboard/api/internal/utils"
)

type Handler struct {
	streamRouteStore store.Interface
	upstreamStore    store.Interface
}

func NewHandler() (handler.RouteRegister, error) {
	return &Handler{
		streamRouteStore: store.GetStore(store.HubKeyStreamRoute),
		upstreamStore:    store.GetStore(store.HubKeyUpstream),
	}, nil
}

func (h *Handler) ApplyRoute(r *gin.Engine) {
	r.GET("/apisix/admin/stream_routes/:id", wgin.Wraps(h.Get,
		wrapper.InputType(reflect.TypeOf(GetInput{}))))
	r.GET("/apisix/admin/stream_routes", wgin.Wraps(h.List,
		wrapper.InputType(reflect.TypeOf(ListInput{}))))
	r.POST("/apisix/admin/stream_routes", wgin.Wraps(h.Create,
		wrapper.InputType(reflect.TypeOf(entity.StreamRoute{}))))
	r.PUT("/apisix/admin/stream_routes", wgin.Wraps(h.Update,
		wrapper.InputType(reflect.TypeOf(UpdateInput{}))))
	r.PUT("/apisix/admin/stream_routes/:id", wgin.Wraps(h.Update,
		wrapper.InputType(reflect.TypeOf(UpdateInput{}))))
	r.DELETE("/apisix/admin/stream_routes/:ids", wgin.Wraps(h.BatchDelete,
		wrapper.InputType(reflect.TypeOf(BatchDelete{}))))
}

type GetInput struct {
	ID string `auto_read:"id,path" validate:"required"`
}

func (h *Handler) Get(c droplet.Context) (interface{}, error) {
	input := c.Input().(*GetInput)
	streamRoute, err := h.streamRouteStore.Get(c.Context(), input.ID)
	if err != nil {
		return handler.SpecCodeResponse(err), err
	}
	return streamRoute, nil
}

type ListInput struct {
	RemoteAddr string `auto_read:"remote_addr,query"`
	ServerAddr string `auto_read:"server_addr,query"`
	ServerPort int    `auto_read:"server_port,query"`
	SNI        string `auto_read:"sni,query"`
	store.Pagination
}

func (h *Handler) List(c droplet.Context) (interface{}, error) {
	input := c.Input().(*ListInput)
	ret, err := h.streamRouteStore.List(c.Context(), store.ListInput{
		Predicate: func(obj interface{}) bool {
			if input.RemoteAddr != "" && !strings.Contains(obj.(*entity.StreamRoute).RemoteAddr, input.RemoteAddr) {
				return false
			}

			if input.ServerAddr != "" && !strings.Contains(obj.(*entity.StreamRoute).ServerAddr, input.ServerAddr) {
				return false
			}

			if input.ServerPort != 0 && obj.(*entity.StreamRoute).ServerPort != input.ServerPort {
				return false
			}

			if input.SNI != "" && !strings.Contains(obj.(*entity.StreamRoute).SNI, input.SNI) {
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
	streamRoute := c.Input().(*entity.StreamRoute)
	if streamRoute.UpstreamID != nil {
		upstreamID := utils.InterfaceToString(streamRoute.UpstreamID)
		_, err := h.upstreamStore.Get(c.Context(), upstreamID)
		if err != nil {
			if err == data.ErrNotFound {
				return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
					fmt.Errorf("upstream id: %s not found", streamRoute.UpstreamID)
			}
			return handler.SpecCodeResponse(err), err
		}
	}
	create, err := h.streamRouteStore.Create(c.Context(), streamRoute)
	if err != nil {
		return handler.SpecCodeResponse(err), err
	}
	return create, nil
}

type UpdateInput struct {
	ID string `auto_read:"id,path"`
	entity.StreamRoute
}

func (h *Handler) Update(c droplet.Context) (interface{}, error) {
	input := c.Input().(*UpdateInput)

	// check if ID in body is equal ID in path
	if err := handler.IDCompare(input.ID, input.StreamRoute.ID); err != nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
	}

	// if has id in path, use it
	if input.ID != "" {
		input.StreamRoute.ID = input.ID
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
	res, err := h.streamRouteStore.Update(c.Context(), &input.StreamRoute, true)
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

	if err := h.streamRouteStore.BatchDelete(c.Context(), strings.Split(input.IDs, ",")); err != nil {
		return handler.SpecCodeResponse(err), err
	}

	return nil, nil
}
