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
package route_online_debug

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/core/store"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/wrapper"
	wgin "github.com/shiningrush/droplet/wrapper/gin"

	"github.com/apisix/manager-api/internal/handler"
)

type Handler struct {
	routeStore    store.Interface
	svcStore      store.Interface
	upstreamStore store.Interface
}

func NewHandler() (handler.RouteRegister, error) {
	return &Handler{
		routeStore:    store.GetStore(store.HubKeyRoute),
		svcStore:      store.GetStore(store.HubKeyService),
		upstreamStore: store.GetStore(store.HubKeyUpstream),
	}, nil
}

type ProtocolSupport interface {
	RequestForwarding(c droplet.Context) (interface{}, error)
}

func (h *Handler) ApplyRoute(r *gin.Engine) {
	r.POST("/apisix/admin/debug-request-forwarding", wgin.Wraps(h.DebugRequestForwarding,
		wrapper.InputType(reflect.TypeOf(DebugOnlineInput{}))))
}

type DebugOnlineInput struct {
	ID              string `auto_read:"online_debug_route_id,header"`
	RequestPath     string `auto_read:"online_debug_path,header"`
	RequestProtocol string `auto_read:"online_debug_request_protocol,header"`
	Method          string `auto_read:"online_debug_method,header"`
	HeaderParams    string `auto_read:"online_debug_header_params,header"`
	ContentType     string `auto_read:"Content-Type,header"`
	Body            []byte `auto_read:"@body"`
}

type Result struct {
	Code    int         `json:"code,omitempty"`
	Header  interface{} `json:"header,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *Handler) DebugRequestForwarding(c droplet.Context) (interface{}, error) {
	//TODO: other Protocols, e.g: grpc, websocket
	input := c.Input().(*DebugOnlineInput)
	protocol := input.RequestProtocol
	if protocol == "" {
		protocol = "http"
	}

	protocolMap := make(map[string]ProtocolSupport)
	protocolMap["http"] = &HTTPProtocolSupport{handler: h}
	protocolMap["https"] = &HTTPProtocolSupport{handler: h}

	if v, ok := protocolMap[protocol]; ok {
		ret, err := v.RequestForwarding(c)
		return ret, err
	}

	return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, fmt.Errorf("protocol unsupported %s", protocol)
}

type HTTPProtocolSupport struct {
	handler *Handler
}

func (h *HTTPProtocolSupport) RequestForwarding(c droplet.Context) (interface{}, error) {
	input := c.Input().(*DebugOnlineInput)
	requestPath := input.RequestPath
	method := input.Method
	body := input.Body
	contentType := input.ContentType
	var upstreamAddress string

	r, err := h.handler.routeStore.Get(c.Context(), input.ID)
	if err != nil {
		return nil, err
	}
	route := r.(*entity.Route)

	host := route.Host

	if route.Plugins != nil {
		if uri, ok := route.Plugins["proxy-rewrite"].(map[string]interface{})["uri"]; ok {
			requestPath = uri.(string)
		}
		if rewriteHost, ok := route.Plugins["proxy-rewrite"].(map[string]interface{})["host"]; ok {
			host = rewriteHost.(string)
		}
	}

	if route.Upstream != nil {
		upstreamAddress = route.Upstream.Nodes.([]*entity.Node)[0].Host + ":" + strconv.Itoa(route.Upstream.Nodes.([]*entity.Node)[0].Port)
	} else if route.UpstreamID != nil {
		up, err := h.handler.upstreamStore.Get(c.Context(), route.UpstreamID.(string))
		if err != nil {
			return nil, err
		}
		upstream := up.(*entity.Upstream)
		upstreamAddress = upstream.UpstreamDef.Nodes.([]*entity.Node)[0].Host + ":" + strconv.Itoa(upstream.UpstreamDef.Nodes.([]*entity.Node)[0].Port)
	}

	url := input.RequestProtocol + "://" + upstreamAddress + requestPath

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableCompression = true

	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}

	var tempMap map[string][]string
	err = json.Unmarshal([]byte(input.HeaderParams), &tempMap)

	if err != nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, fmt.Errorf("can not get header")
	}

	req, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewReader(body))
	if err != nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusBadRequest}, err
	}

	req.Header.Add("Content-Type", contentType)
	for k, v := range tempMap {
		for _, v1 := range v {
			if !strings.EqualFold(k, "Content-Type") {
				req.Header.Add(k, v1)
			} else {
				req.Header.Set(k, v1)
			}
		}
	}

	if req.Header.Get("Host") != "" {
		req.Header.Set("Host", host)
	} else {
		req.Header.Add("Host", host)
	}

	resp, err := client.Do(req)
	if err != nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusInternalServerError}, err
	}

	defer resp.Body.Close()

	// handle gzip content encoding
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return &data.SpecCodeResponse{StatusCode: http.StatusInternalServerError}, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	_body, err := ioutil.ReadAll(reader)
	if err != nil {
		return &data.SpecCodeResponse{StatusCode: http.StatusInternalServerError}, err
	}
	returnData := make(map[string]interface{})
	result := &Result{}
	err = json.Unmarshal(_body, &returnData)
	if err != nil {
		result.Code = resp.StatusCode
		result.Header = resp.Header
		result.Message = resp.Status
		result.Data = string(_body)
	} else {
		result.Code = resp.StatusCode
		result.Header = resp.Header
		result.Message = resp.Status
		result.Data = returnData
	}
	return result, nil
}
