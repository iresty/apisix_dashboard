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
	"compress/gzip"
	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/core/store"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

var TestResponse = "test"

func mockServer() *httptest.Server {
	l, _ := net.Listen("tcp", "127.0.0.1:9000")
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "plain/text")
		w.Header().Set("Content-Encoding", "gzip")
		writer, _ := gzip.NewWriterLevel(w, gzip.BestCompression)
		defer writer.Close()
		_, _ = writer.Write([]byte(TestResponse))
	}
	ts := httptest.NewUnstartedServer(http.HandlerFunc(f))
	ts.Listener.Close()
	ts.Listener = l
	ts.Start()
	return ts
}

func TestDebugRequestForwardingANDRequestForwarding(t *testing.T) {
	server := mockServer()
	defer server.Close()

	var (
		cases = []struct {
			Desc       string
			Input      *DebugOnlineInput
			Result     interface{}
			MockRoute  *entity.Route
			MockStream *entity.Upstream
		}{
			{
				Desc: "correct request with no UpstreamID",
				Input: &DebugOnlineInput{
					ID:              "test1",
					RequestPath:     "test/abc",
					RequestProtocol: "http",
					Method:          "Get",
					HeaderParams:    `{"Accept-Encoding": ["gzip"]}`,
					ContentType:     "application/json",
					Body:            nil,
				},
				MockRoute: &entity.Route{
					BaseInfo: entity.BaseInfo{
						ID: "test1",
					},
					Host: "1.2.3.4",
					Plugins: map[string]interface{}{
						"proxy-rewrite": map[string]interface{}{
							"uri":  "/test/1",
							"host": "1.2.3.4",
						},
					},
					Upstream: &entity.UpstreamDef{
						Nodes: []*entity.Node{
							{
								Host: "127.0.0.1",
								Port: 9000,
							},
						},
					},
					UpstreamID: nil,
				},
				MockStream: &entity.Upstream{
					UpstreamDef: entity.UpstreamDef{
						Nodes: nil,
					},
				},
				Result: TestResponse,
			},
			{
				Desc: "correct request with no Upstream",
				Input: &DebugOnlineInput{
					ID:              "test2",
					RequestPath:     "test/abc",
					RequestProtocol: "http",
					Method:          "Get",
					HeaderParams:    `{"Accept-Encoding": ["gzip"]}`,
					ContentType:     "application/json",
					Body:            nil,
				},
				MockRoute: &entity.Route{
					BaseInfo: entity.BaseInfo{
						ID: "test1",
					},
					Host: "1.2.3.4",
					Plugins: map[string]interface{}{
						"proxy-rewrite": map[string]interface{}{
							"uri":  "/test/1",
							"host": "1.2.3.4",
						},
					},
					Upstream:   nil,
					UpstreamID: "test2",
				},
				MockStream: &entity.Upstream{
					BaseInfo: entity.BaseInfo{
						ID: "test2",
					},
					UpstreamDef: entity.UpstreamDef{
						Nodes: []*entity.Node{
							{
								Host: "127.0.0.1",
								Port: 9000,
							},
						},
					},
				},
				Result: TestResponse,
			},
			{
				Desc: "unsupported protocol",
				Input: &DebugOnlineInput{
					ID:              "test1",
					RequestPath:     "test/abc",
					RequestProtocol: "grpc",
					Method:          "Get",
					HeaderParams:    `{"Accept-Encoding": ["gzip"]}`,
					ContentType:     "application/json",
					Body:            nil,
				},
				MockRoute: &entity.Route{
					BaseInfo: entity.BaseInfo{
						ID: "test1",
					},
					Host: "1.2.3.4",
					Plugins: map[string]interface{}{
						"proxy-rewrite": map[string]interface{}{
							"uri":  "/test/1",
							"host": "1.2.3.4",
						},
					},
					Upstream: &entity.UpstreamDef{
						Nodes: []*entity.Node{
							{
								Host: "127.0.0.1",
								Port: 9000,
							},
						},
					},
					UpstreamID: nil,
				},
				MockStream: &entity.Upstream{
					UpstreamDef: entity.UpstreamDef{
						Nodes: nil,
					},
				},
				Result: &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
			},
			{
				Desc: "unsupported HeaderParams",
				Input: &DebugOnlineInput{
					ID:              "test1",
					RequestPath:     "test/abc",
					RequestProtocol: "http",
					Method:          "Get",
					HeaderParams:    "",
					ContentType:     "application/json",
					Body:            nil,
				},
				MockRoute: &entity.Route{
					BaseInfo: entity.BaseInfo{
						ID: "test1",
					},
					Host: "1.2.3.4",
					Plugins: map[string]interface{}{
						"proxy-rewrite": map[string]interface{}{
							"uri":  "/test/1",
							"host": "1.2.3.4",
						},
					},
					Upstream: &entity.UpstreamDef{
						Nodes: []*entity.Node{
							{
								Host: "127.0.0.1",
								Port: 9000,
							},
						},
					},
					UpstreamID: nil,
				},
				MockStream: &entity.Upstream{
					UpstreamDef: entity.UpstreamDef{
						Nodes: nil,
					},
				},
				Result: &data.SpecCodeResponse{StatusCode: http.StatusBadRequest},
			},
			{
				Desc: "wrong url",
				Input: &DebugOnlineInput{
					ID:              "test1",
					RequestPath:     "test/abc",
					RequestProtocol: "http",
					Method:          "Get",
					HeaderParams:    `{"Accept-Encoding": ["gzip"]}`,
					ContentType:     "application/json",
					Body:            nil,
				},
				MockRoute: &entity.Route{
					BaseInfo: entity.BaseInfo{
						ID: "test1",
					},
					Host: "1.2.3.4",
					Plugins: map[string]interface{}{
						"proxy-rewrite": map[string]interface{}{
							"uri":  "/test/1",
							"host": "1.2.3.4",
						},
					},
					Upstream: &entity.UpstreamDef{
						Nodes: []*entity.Node{
							{
								Host: "127.0.0.2",
								Port: 9000,
							},
						},
					},
					UpstreamID: nil,
				},
				MockStream: &entity.Upstream{
					UpstreamDef: entity.UpstreamDef{
						Nodes: nil,
					},
				},
				Result: &data.SpecCodeResponse{StatusCode: http.StatusInternalServerError},
			},
		}
	)

	for _, c := range cases {
		routeStore := &store.MockInterface{}
		routeStore.On("Get", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		}).Return(c.MockRoute, nil)
		upstreamStore := &store.MockInterface{}
		upstreamStore.On("Get", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		}).Return(c.MockStream, nil)
		t.Run(c.Desc, func(t *testing.T) {
			handler := Handler{
				routeStore:    routeStore,
				svcStore:      nil,
				upstreamStore: upstreamStore,
			}

			context := droplet.NewContext()
			context.SetInput(c.Input)

			result, _ := handler.DebugRequestForwarding(context)

			switch result.(type) {
			case *Result:
				assert.Equal(t, c.Result, result.(*Result).Data)
			case *data.SpecCodeResponse:
				assert.Equal(t, c.Result, result)
			}
		})
	}
}
