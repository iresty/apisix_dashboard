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
	"testing"

	"github.com/shiningrush/droplet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/core/store"
	"github.com/apisix/manager-api/internal/utils"
)

type mockRet struct {
	routeData    []*entity.Route
	upstreamData []*entity.Upstream
	serviceData  []*entity.Service
	sslData      []*entity.SSL
	serverData   []*entity.ServerInfo
}

func TestDashboard(t *testing.T) {
	_, version := utils.GetHashAndVersion()
	var (
		tests = []struct {
			caseDesc   string
			giveInput  *ListInput
			giveErr    error
			giveRet    mockRet
			wantErr    error
			wantGetKey string
			wantRet    interface{}
		}{
			{
				caseDesc:  "get server_info",
				giveInput: &ListInput{},
				giveRet: mockRet{
					routeData: []*entity.Route{
						{
							BaseInfo: entity.BaseInfo{
								ID:         "s1",
								CreateTime: 1609746531,
							},
							Name:       "s1",
							Desc:       "test_route",
							UpstreamID: "u1",
							ServiceID:  "s1",
							ScriptID:   "s1",
							Script:     "",
							Labels: map[string]string{
								"version": "v1",
							},
							Status: 1,
						},
						{
							BaseInfo: entity.BaseInfo{
								ID:         "s2",
								CreateTime: 1609746661,
							},
							Name:       "s2",
							Desc:       "test_route",
							UpstreamID: "u2",
							ServiceID:  "s2",
							ScriptID:   "s2",
							Script:     "",
							Labels: map[string]string{
								"version": "v1",
							},
							Status: 0,
						},
					},
					upstreamData: []*entity.Upstream{
						{
							BaseInfo: entity.BaseInfo{
								ID:         "u1",
								CreateTime: 1609340491,
								UpdateTime: 1609340491,
							},
							UpstreamDef: entity.UpstreamDef{
								Name: "upstream1",
								Key:  "server_addr",
								Nodes: []map[string]interface{}{
									{
										"host":   "39.97.63.215",
										"port":   80,
										"weight": 1,
									},
								},
							},
						},
						{
							BaseInfo: entity.BaseInfo{
								ID:         "u2",
								CreateTime: 1609340491,
								UpdateTime: 1609340491,
							},
							UpstreamDef: entity.UpstreamDef{
								Name: "upstream2",
								Key:  "server_addr2",

								Nodes: entity.Node{
									Host:   "39.97.63.215",
									Port:   80,
									Weight: 0,
								},
							},
						},
						{
							BaseInfo: entity.BaseInfo{
								ID:         "u3",
								CreateTime: 1609340491,
								UpdateTime: 1609340491,
							},
							UpstreamDef: entity.UpstreamDef{
								Name: "upstream3",
								Key:  "server_addr3",
								Nodes: []entity.Node{
									{
										Host:   "39.97.63.215",
										Port:   80,
										Weight: 0,
									},
								},
							},
						},
					},
					serviceData: []*entity.Service{
						{Name: "s1"},
						{Name: "s2"},
						{Name: "test_service"},
						{Name: "service_test"},
					},
					sslData: []*entity.SSL{
						{
							BaseInfo: entity.BaseInfo{
								ID:         "ssl1",
								CreateTime: 1609340491,
								UpdateTime: 1609340491,
							},
							Status: 0,
							Labels: map[string]string{
								"build":   "16",
								"env":     "production",
								"version": "v2",
							},
						},
						{
							BaseInfo: entity.BaseInfo{
								ID:         "ssl2",
								CreateTime: 1609340492,
								UpdateTime: 1609340492,
							},
							Sni:    "route",
							Status: 1,
							Labels: map[string]string{
								"build":   "17",
								"env":     "production",
								"version": "v2",
							},
						},
						{
							BaseInfo: entity.BaseInfo{
								ID:         "ssl3",
								CreateTime: 1609340493,
								UpdateTime: 1609340493,
							},
							Status: 0,
							Labels: map[string]string{
								"build":   "18",
								"env":     "production",
								"version": "v2",
							},
						},
					},
					serverData: []*entity.ServerInfo{
						{
							BaseInfo:       entity.BaseInfo{ID: "server_1"},
							UpTime:         10,
							LastReportTime: 1608195454,
							BootTime:       1608195454,
							Hostname:       "gentoo",
							Version:        "v3",
						},
						{
							BaseInfo:       entity.BaseInfo{ID: "server_2"},
							UpTime:         10,
							LastReportTime: 1608195454,
							BootTime:       1608195454,
							Hostname:       "ubuntu",
							Version:        "v2",
						},
					},
				},
				wantGetKey: "dashboard",
				wantRet: &entity.Overview{
					RouteCnt:         2,
					OnlineRouterCnt:  1,
					UpstreamCnt:      3,
					ServiceCnt:       4,
					CertificateCnt:   3,
					Plugins:          []string{"api-breaker", "authz-casbin", "authz-keycloak", "aws-lambda", "azure-functions", "basic-auth", "batch-requests", "client-control", "consumer-restriction", "cors", "datadog", "dubbo-proxy", "echo", "error-log-logger", "example-plugin", "ext-plugin-post-req", "ext-plugin-pre-req", "fault-injection", "forward-auth", "google-cloud-logging", "grpc-transcode", "grpc-web", "gzip", "hmac-auth", "http-logger", "ip-restriction", "jwt-auth", "kafka-logger", "key-auth", "ldap-auth", "limit-conn", "limit-count", "limit-req", "log-rotate", "node-status", "opa", "openid-connect", "openwhisk", "prometheus", "proxy-cache", "proxy-control", "proxy-mirror", "proxy-rewrite", "real-ip", "redirect", "referer-restriction", "request-id", "request-validation", "response-rewrite", "rocketmq-logger", "server-info", "serverless-post-function", "serverless-pre-function", "skywalking", "skywalking-logger", "sls-logger", "splunk-hec-logging", "syslog", "tcp-logger", "traffic-split", "ua-restriction", "udp-logger", "uri-blocker", "wolf-rbac", "zipkin"},
					DashboardVersion: version,
					GatewayInfo: []*entity.ServerInfo{
						{
							BaseInfo:       entity.BaseInfo{ID: "server_1"},
							UpTime:         10,
							LastReportTime: 1608195454,
							BootTime:       1608195454,
							Hostname:       "gentoo",
							Version:        "v3",
						},
						{
							BaseInfo:       entity.BaseInfo{ID: "server_2"},
							UpTime:         10,
							LastReportTime: 1608195454,
							BootTime:       1608195454,
							Hostname:       "ubuntu",
							Version:        "v2",
						},
					},
				},
			},
		}
	)

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			getCalled := false

			routeStore := &store.MockInterface{}
			routeStore.On("List", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				getCalled = true
			}).Return(func(input store.ListInput) *store.ListOutput {
				var returnData []interface{}
				for _, c := range tc.giveRet.routeData {
					if input.Predicate(c) {
						if input.Format == nil {
							returnData = append(returnData, c)
							continue
						}
						returnData = append(returnData, input.Format(c))
					}
				}
				return &store.ListOutput{
					Rows:      returnData,
					TotalSize: len(returnData),
				}
			}, tc.giveErr)

			upstreamStore := &store.MockInterface{}
			upstreamStore.On("List", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				getCalled = true
			}).Return(func(input store.ListInput) *store.ListOutput {
				var returnData []interface{}
				for _, c := range tc.giveRet.upstreamData {
					if input.Predicate(c) {
						if input.Format == nil {
							returnData = append(returnData, c)
							continue
						}
						returnData = append(returnData, input.Format(c))
					}
				}
				return &store.ListOutput{
					Rows:      returnData,
					TotalSize: len(returnData),
				}
			}, tc.giveErr)

			sslStore := &store.MockInterface{}
			sslStore.On("List", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				getCalled = true
			}).Return(func(input store.ListInput) *store.ListOutput {
				var returnData []interface{}
				for _, c := range tc.giveRet.sslData {
					if input.Predicate(c) {
						if input.Format == nil {
							returnData = append(returnData, c)
							continue
						}
						returnData = append(returnData, input.Format(c))
					}
				}
				return &store.ListOutput{
					Rows:      returnData,
					TotalSize: len(returnData),
				}
			}, tc.giveErr)

			serviceStore := &store.MockInterface{}
			serviceStore.On("List", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				getCalled = true
			}).Return(func(input store.ListInput) *store.ListOutput {
				var returnData []interface{}
				for _, c := range tc.giveRet.serviceData {
					if input.Predicate(c) {
						if input.Format == nil {
							returnData = append(returnData, c)
							continue
						}
						returnData = append(returnData, input.Format(c))
					}
				}
				return &store.ListOutput{
					Rows:      returnData,
					TotalSize: len(returnData),
				}
			}, tc.giveErr)

			serverInfoStore := &store.MockInterface{}
			serverInfoStore.On("List", mock.Anything).Run(func(args mock.Arguments) {
				getCalled = true
			}).Return(func(input store.ListInput) *store.ListOutput {
				var returnData []interface{}
				for _, c := range tc.giveRet.serverData {
					if input.Predicate(c) {
						if input.Format == nil {
							returnData = append(returnData, c)
							continue
						}
						returnData = append(returnData, input.Format(c))
					}
				}
				return &store.ListOutput{
					Rows:      returnData,
					TotalSize: len(returnData),
				}
			}, tc.giveErr)

			h := Handler{
				routeStore:      routeStore,
				upstreamStore:   upstreamStore,
				serverInfoStore: serverInfoStore,
				serviceStore:    serviceStore,
				sslStore:        sslStore,
			}
			ctx := droplet.NewContext()
			ctx.SetInput(tc.giveInput)
			ret, err := h.List(ctx)
			assert.True(t, getCalled)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}
