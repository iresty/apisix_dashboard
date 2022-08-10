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
package route_test

import (
	"net/http"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Route", func() {
	Context("test route export data empty", func() {
		It("Export route when data is empty", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   "{\"code\":10000,\"message\":\"Route data is empty, cannot be exported\",\"data\":null",
			})
		})
	})

	Context("test route export", func() {
		// 1.Export data as the route of URIs Hosts
		exportStrR1 := `
			"/hello_": {
				"get": {
					"operationId": "aaaaGET",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"security": [],
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["foo.com", "*.bar.com"],
					"x-apisix-id":"r1",
					"x-apisix-labels": {
						"build": "16",
						"env": "production",
						"version": "v2"
					},
					"x-apisix-plugins": {
						"limit-count": {
							"count": 2,
							"key": "remote_addr",
							"policy": "local",
							"rejected_code": 503,
							"time_window": 60
						}
					},
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"` + base.UpstreamIp + `:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"post": {
					"operationId": "aaaaPOST",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"security": [],
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["foo.com", "*.bar.com"],
					"x-apisix-id":"r1",
					"x-apisix-labels": {
						"build": "16",
						"env": "production",
						"version": "v2"
					},
					"x-apisix-plugins": {
						"limit-count": {
							"count": 2,
							"key": "remote_addr",
							"policy": "local",
							"rejected_code": 503,
							"time_window": 60
						}
					},
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"` + base.UpstreamIp + `:1980": 1
						},
						"type": "roundrobin"
					}
				}
			}`
		exportStrR1 = replaceStr(exportStrR1)
		It("hit route that not exist", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello_",
				Headers:      map[string]string{"Host": "foo.com"},
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			})
		})
		It("create route with uris and hosts to test whether the uris parsing is correct", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r1",
				Body: `{
					"name": "aaaa",
					"labels": {
						"build":"16",
						"env":"production",
						"version":"v2"
					},
					"plugins": {
						"limit-count": {
							"count": 2,
							"time_window": 60,
							"rejected_code": 503,
							"key": "remote_addr",
							"policy": "local"
						}
					},
					"status": 1,
					"uris": ["/hello_"],
					"hosts": ["foo.com", "*.bar.com"],
					"methods": ["GET", "POST"],
					"upstream": {
						"nodes": {
							"` + base.UpstreamIp + `:1980": 1
						},
						"type": "roundrobin"
					}
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})

		It("export route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   "{\"components\":{},\"info\":{\"title\":\"RoutesExport\",\"version\":\"3.0.0\"},\"openapi\":\"3.0.0\",\"paths\":{" + exportStrR1 + "}}",
			})
		})
		// 2.Export data as the route of URI host
		exportStrR2 := `
			"/hello2": {
				"get": {
					"operationId": "aaaa2GET",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"security": [],
					"x-apisix-enable_websocket": false,
					"x-apisix-host": "*.bar.com",
					"x-apisix-id":"r2",
					"x-apisix-labels": {
						"build": "16",
						"env": "production",
						"version": "v2"
					},
					"x-apisix-plugins": {
						"limit-count": {
							"count": 2,
							"key": "remote_addr",
							"policy": "local",
							"rejected_code": 503,
							"time_window": 60
						}
					},
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"` + base.UpstreamIp + `:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"post": {
					"operationId": "aaaa2POST",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"security": [],
					"x-apisix-enable_websocket": false,
					"x-apisix-host": "*.bar.com",
					"x-apisix-id":"r2",
					"x-apisix-labels": {
						"build": "16",
						"env": "production",
						"version": "v2"
					},
					"x-apisix-plugins": {
						"limit-count": {
							"count": 2,
							"key": "remote_addr",
							"policy": "local",
							"rejected_code": 503,
							"time_window": 60
						}
					},
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"` + base.UpstreamIp + `:1980": 1
						},
						"type": "roundrobin"
					}
				}
			}`
		exportStrR2 = replaceStr(exportStrR2)

		It("hit route2 that not exist", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello2",
				Headers:      map[string]string{"Host": "bar.com"},
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			})
		})
		It("create route2 with uri and host to test whether the uri parsing is correct", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r2",
				Body: `{
					"name": "aaaa2",
					"labels": {
						"build":"16",
						"env":"production",
						"version":"v2"
					},
					"plugins": {
						"limit-count": {
							"count": 2,
							"time_window": 60,
							"rejected_code": 503,
							"key": "remote_addr",
							"policy": "local"
						}
					},
					"status": 1,
					"uri": "/hello2",
					"host": "*.bar.com",
					"methods": ["GET", "POST"],
					"upstream": {
						"nodes": {
							"` + base.UpstreamIp + `:1980": 1
						},
						"type": "roundrobin"
					}
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("export route2", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r2",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   "{\"components\":{},\"info\":{\"title\":\"RoutesExport\",\"version\":\"3.0.0\"},\"openapi\":\"3.0.0\",\"paths\":{" + exportStrR2 + "}}",
			})
		})
		It("export route and route2", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r1,r2",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   "{\"components\":{},\"info\":{\"title\":\"RoutesExport\",\"version\":\"3.0.0\"},\"openapi\":\"3.0.0\",\"paths\":{" + exportStrR2 + "," + exportStrR1 + "}}",
			})
		})
		It("use the exportall interface to export all routes", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   "{\"components\":{},\"info\":{\"title\":\"RoutesExport\",\"version\":\"3.0.0\"},\"openapi\":\"3.0.0\",\"paths\":{" + exportStrR2 + "," + exportStrR1 + "}}",
			})
		})
		It("delete the route just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("hit the route just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello_",
				Headers:      map[string]string{"Host": "bar.com"},
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
				Sleep:        base.SleepTime,
			})
		})
		It("delete the route2 just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r2",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("hit the route2 just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello2",
				Headers:      map[string]string{"Host": "bar.com"},
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
				Sleep:        base.SleepTime,
			})
		})
		// 4.Create a service that contains complete data and use the service_ id create route
		serviceStrS1 := `
			"name": "testservice",
			"desc": "testservice_desc",
			"upstream": {
				"nodes": [{
					"host": "` + base.UpstreamIp + `",
					"port": 1980,
					"weight": 1
				}],
				"type": "roundrobin"
			},
			"plugins": {
				"limit-count": {
					"count": 100,
					"key": "remote_addr",
					"policy": "local",
					"rejected_code": 503,
					"time_window": 60
				}
			},
			"labels": {
				"build": "16",
				"env": "production",
				"version": "v2"
			},
			"enable_websocket": true
			`
		serviceStrS1 = replaceStr(serviceStrS1)

		exportStrR3 := `
			"components": {},
			"info": {
				"title": "RoutesExport",
				"version": "3.0.0"
			},
			"openapi": "3.0.0",
			"paths": {
				"/hello3": {
					"get": {
						"operationId": "route3GET",
						"requestBody": {},
						"responses": {
							"default": {
								"description": ""
							}
						},
						"x-apisix-enable_websocket": false,
						"x-apisix-id":"r3",
						"x-apisix-labels": {
							"build": "16",
							"env": "production",
							"version": "v2"
						},
						"x-apisix-plugins": {
							"limit-count": {
								"count": 100,
								"key": "remote_addr",
								"policy": "local",
								"rejected_code": 503,
								"time_window": 60
							}
						},
						"x-apisix-priority": 0,
						"x-apisix-status": 1,
						"x-apisix-upstream": {
							"nodes": [{
								"host": "` + base.UpstreamIp + `",
								"port": 1980,
								"weight": 1
							}],
							"type": "roundrobin"
						}
					}
				}
			}`
		exportStrR3 = replaceStr(exportStrR3)

		It("create service with all options", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:  base.ManagerApiExpect(),
				Method:  http.MethodPut,
				Path:    "/apisix/admin/services/s1",
				Headers: map[string]string{"Authorization": base.GetToken()},
				Body: `{
					 "name": "testservice",
					 "desc": "testservice_desc",
					 "labels": {
						 "build":"16",
						 "env":"production",
						 "version":"v2"
					 },
					 "enable_websocket":true,
					 "plugins": {
						 "limit-count": {
							 "count": 100,
							 "time_window": 60,
							 "policy": "local",
							 "rejected_code": 503,
							 "key": "remote_addr"
						 }
					 },
					 "upstream": {
						 "type": "roundrobin",
						 "create_time":1602883670,
						 "update_time":1602893670,
						 "nodes": [{
							 "host": "` + base.UpstreamIp + `",
							 "port": 1980,
							 "weight": 1
						 }]
					 }
				 }`,
				ExpectStatus: http.StatusOK,
			})
		})
		It("get the service s1", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:     base.ManagerApiExpect(),
				Method:     http.MethodGet,
				Path:       "/apisix/admin/services/s1",
				Headers:    map[string]string{"Authorization": base.GetToken()},
				ExpectCode: http.StatusOK,
				ExpectBody: serviceStrS1,
			})
		})
		It("create route3 using the service id just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r3",
				Body: `{
					 "name": "route3",
					 "methods": ["GET"],
					 "uri": "/hello3",
					 "service_id": "s1"
				 }`,
				Headers:    map[string]string{"Authorization": base.GetToken()},
				ExpectCode: http.StatusOK,
				Sleep:      base.SleepTime,
			})
		})
		It("export route3", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:     base.ManagerApiExpect(),
				Method:     http.MethodGet,
				Path:       "/apisix/admin/export/routes/r3",
				Headers:    map[string]string{"Authorization": base.GetToken()},
				ExpectCode: http.StatusOK,
				ExpectBody: exportStrR3,
			})
		})
		It("delete the route3 just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:     base.ManagerApiExpect(),
				Method:     http.MethodDelete,
				Path:       "/apisix/admin/routes/r3",
				Headers:    map[string]string{"Authorization": base.GetToken()},
				ExpectCode: http.StatusOK,
			})
		})
		It("hit the route3 just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello3",
				ExpectStatus: http.StatusNotFound,
				Sleep:        base.SleepTime,
			})
		})
		It("delete the service1", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/services/s1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		// 5.Create a service containing plugin and a route containing plugin to test the fusion of exported data
		// This test case tests the creation of a complete service.
		// And create a complete route, export rules as route upstream data for high priority, direct use.
		// However, if the data in the service plugins does not exist in the route, it will be fused and exported.
		// If it exists, the data in route will be used first
		serviceStrS2 := `
			"name": "testservice",
			"desc": "testservice_desc",
			"upstream": {
				"nodes": [{
					"host": "` + base.UpstreamIp + `",
					"port": 1980,
					"weight": 1
				}],
				"type": "roundrobin"
			},
			"plugins": {
				"limit-count": {
					"count": 100,
					"key": "remote_addr",
					"policy": "local",
					"rejected_code": 503,
					"time_window": 60
				}
			},
			"labels": {
				"build": "16",
				"env": "production",
				"version": "v2"
			},
			"enable_websocket": true`
		serviceStrS2 = replaceStr(serviceStrS2)

		exportStrR4 := `
			"components": {},
			"info": {
				"title": "RoutesExport",
				"version": "3.0.0"
			},
			"openapi": "3.0.0",
			"paths": {
				"/hello": {
					"get": {
						"operationId": "route4GET",
						"requestBody": {},
						"responses": {
							"default": {
								"description": ""
							}
						},
						"security": [],
						"x-apisix-enable_websocket": false,
						"x-apisix-id":"r4",
						"x-apisix-labels": {
							"build": "16",
							"env": "production",
							"version": "v2"
						},
						"x-apisix-plugins": {
							"limit-count": {
								"count": 100,
								"key": "remote_addr",
								"policy": "local",
								"rejected_code": 503,
								"time_window": 60
							},
							"prometheus": {
								"disable": false
							}
						},
						"x-apisix-priority": 0,
						"x-apisix-status": 1,
						"x-apisix-upstream": {
							"nodes": [{
								"host": "` + base.UpstreamIp + `",
								"port": 1980,
								"weight": 1
							}],
							"type": "roundrobin"
						}
					}
				}
			}`
		exportStrR4 = replaceStr(exportStrR4)

		It("create service with all options", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:  base.ManagerApiExpect(),
				Method:  http.MethodPut,
				Path:    "/apisix/admin/services/s2",
				Headers: map[string]string{"Authorization": base.GetToken()},
				Body: `{
					"name": "testservice",
					"desc": "testservice_desc",
					"labels": {
						"build":"16",
						"env":"production",
						"version":"v2"
					},
					"enable_websocket":true,
					"plugins": {
						"limit-count": {
							"count": 100,
							"time_window": 60,
							"rejected_code": 503,
							"key": "remote_addr",
							"policy": "local"
						}
					},
					"upstream": {
						"type": "roundrobin",
						"create_time":1602883670,
						"update_time":1602893670,
						"nodes": [{
							"host": "` + base.UpstreamIp + `",
							"port": 1980,
							"weight": 1
						}]
					}
				}`,
				ExpectStatus: http.StatusOK,
			})
		})
		It("get the service s2", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/services/s2",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   serviceStrS2,
			})
		})
		It("Create Route4 and test the priority merging function of upstream, label and plugin when both service and route are included", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r4",
				Body: `{
					"name": "route4",
					"methods": ["GET"],
					"uri": "/hello",
					"service_id": "s2",
					"enable_websocket":false,
					"plugins": {
						"prometheus": {
							"disable": false
						}
					},
					"upstream": {
						"type": "roundrobin",
						"nodes": [{
							"host": "` + base.UpstreamIp + `",
							"port": 1980,
							"weight": 1
						}]
					}
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("export route4", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r4",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   exportStrR4,
			})
		})
		It("delete the route4 just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r4",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("hit the route4 just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			})
		})
		It("delete the service2", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/services/s2",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		// 6.Create a service according to the upstream ID and a route according to the service ID.
		// The test export also contains the upstream.
		// Use the upstream of route.
		serviceStrS3 := `
			"name": "testservice",
			"desc": "testservice_desc",
			"upstream_id": "1",
			"plugins": {
				"limit-count": {
					"count": 100,
					"key": "remote_addr",
					"policy": "local",
					"rejected_code": 503,
					"time_window": 60
				}
			},
			"labels": {
				"build": "16",
				"env": "production",
				"version": "v2"
			},
			"enable_websocket": true`
		serviceStrS3 = replaceStr(serviceStrS3)

		exportStrR5 := `
			"components": {},
			"info": {
				"title": "RoutesExport",
				"version": "3.0.0"
			},
			"openapi": "3.0.0",
			"paths": {
				"/hello5": {
					"get": {
						"operationId": "route5GET",
						"requestBody": {},
						"responses": {
							"default": {
								"description": ""
							}
						},
						"security": [],
						"x-apisix-enable_websocket": false,
						"x-apisix-id":"r5",
						"x-apisix-labels": {
							"build": "16",
							"env": "production",
							"version": "v2"
						},
						"x-apisix-plugins": {
							"limit-count": {
								"count": 100,
								"key": "remote_addr",
								"policy": "local",
								"rejected_code": 503,
								"time_window": 60
							},
							"prometheus": {
								"disable": false
							}
						},
						"x-apisix-priority": 0,
						"x-apisix-status": 1,
						"x-apisix-upstream": {
							"nodes": [{
								"host": "` + base.UpstreamIp + `",
								"port": 1981,
								"weight": 1
							}],
							"type": "roundrobin"
						}
					}
				}
			}`
		exportStrR5 = replaceStr(exportStrR5)
		It("create upstream", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:  base.ManagerApiExpect(),
				Method:  http.MethodPut,
				Path:    "/apisix/admin/upstreams/1",
				Headers: map[string]string{"Authorization": base.GetToken()},
				Body: `{
					"nodes": [
						{
							"host": "` + base.UpstreamIp + `",
							"port": 1980,
							"weight": 1
						}
					],
					"type": "roundrobin"
				}`,
				ExpectStatus: http.StatusOK,
			})
		})
		It("create service with upstream id", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:  base.ManagerApiExpect(),
				Method:  http.MethodPut,
				Path:    "/apisix/admin/services/s3",
				Headers: map[string]string{"Authorization": base.GetToken()},
				Body: `{
					"name": "testservice",
					"desc": "testservice_desc",
					"labels": {
						"build":"16",
						"env":"production",
						"version":"v2"
					},
					"enable_websocket":true,
					"plugins": {
						"limit-count": {
							"count": 100,
							"time_window": 60,
							"rejected_code": 503,
							"key": "remote_addr",
							"policy": "local"
						}
					},
					"upstream_id": "1"
				}`,
				ExpectStatus: http.StatusOK,
			})
		})
		It("get the service s3", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/services/s3",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   serviceStrS3,
			})
		})
		It("Create a route5 with the id of the service3 created with upstream id", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r5",
				Body: `{
					"name": "route5",
					"methods": ["GET"],
					"uri": "/hello5",
					"service_id": "s3",
					"enable_websocket":false,
					"plugins": {
						"prometheus": {
							"disable": false
						}
					},
					"upstream": {
						"type": "roundrobin",
						"nodes": [{
							"host": "` + base.UpstreamIp + `",
							"port": 1981,
							"weight": 1
						}]
					}
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("export route5", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r5",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   exportStrR5,
			})
		})
		It("delete the route5 just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r5",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("hit the route5 just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:     base.APISIXExpect(),
				Method:     http.MethodGet,
				Path:       "/hello5",
				ExpectBody: "{\"error_msg\":\"404 Route Not Found\"}\n",
				Sleep:      base.SleepTime,
			})
		})
		It("delete the service3", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/services/s3",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("remove upstream", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/upstreams/1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		// 8.Create and export route according to upstream ID
		exportStrR8 := `
			"components": {},
			"info": {
				"title": "RoutesExport",
				"version": "3.0.0"
			},
			"openapi": "3.0.0",
			"paths": {
				"/hello": {
					"get": {
						"operationId": "route8GET",
						"requestBody": {},
						"responses": {
							"default": {
								"description": ""
							}
						},
						"security": [],
						"x-apisix-enable_websocket": false,
						"x-apisix-id":"r8",
						"x-apisix-plugins": {
							"prometheus": {
								"disable": false
							}
						},
						"x-apisix-priority": 0,
						"x-apisix-status": 1,
						"x-apisix-upstream": {
							"id": "3"
				`
		exportStrR8 = replaceStr(exportStrR8)

		It("create upstream3", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/upstreams/3",
				Body: `{
					"nodes": [
						{
						"host": "` + base.UpstreamIp + `",
						"port": 1980,
						"weight": 1
						}
					],
					"type": "roundrobin"
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("Create a route8 using upstream id", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r8",
				Body: `{
					"name": "route8",
					"methods": ["GET"],
					"uri": "/hello",
					"enable_websocket":false,
					"plugins": {
						"prometheus": {
							"disable": false
						}
					},
					"upstream_id": "3"
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("export route8", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r8",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   exportStrR8,
			})
		})
		It("delete the route8 just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r8",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("hit the route8 just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
				Sleep:        base.SleepTime,
			})
		})
		It("remove upstream3", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/upstreams/3",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})

		// 9.Create service according to upstream1 ID
		// Create route according to upstream2 ID and service ID
		// Export route
		serviceStrS4 := `
			"name": "testservice",
			"desc": "testservice_desc",
			"upstream_id": "4",
			"enable_websocket": true`
		serviceStrS4 = replaceStr(serviceStrS4)

		exportStrR9 := `
			"components": {},
			"info": {
				"title": "RoutesExport",
				"version": "3.0.0"
			},
			"openapi": "3.0.0",
			"paths": {
				"/hello": {
					"get": {
						"operationId": "route_allGET",
						"requestBody": {},
						"responses": {
							"default": {
								"description": ""
							}
						},
						"security": [],
						"summary": "所有",
						"x-apisix-enable_websocket": false,
						"x-apisix-id":"r9",
						"x-apisix-labels": {
							"API_VERSION": "v1",
							"test": "1"
						},
						"x-apisix-plugins": {
							"prometheus": {
								"disable": false
							}
						},
						"x-apisix-priority": 0,
						"x-apisix-status": 1,
						"x-apisix-upstream": {
							"id": "5"
						`
		exportStrR9 = replaceStr(exportStrR9)

		It("create upstream4", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/upstreams/4",
				Body: `{
				"nodes": [
					{
						"host": "` + base.UpstreamIp + `",
						"port": 1980,
						"weight": 1
							}
						],
						"type": "roundrobin"
					}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("create upstream5", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/upstreams/5",
				Body: `{
					"name": "upstream5",
					"nodes": [
						{
							"host": "` + base.UpstreamIp + `",
							"port": 1981,
							"weight": 1
						}
					],
					"type": "roundrobin"
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("create service", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/services/s4",
				Body: `{
					"name": "testservice",
					"desc": "testservice_desc",
					"enable_websocket":true,
					"upstream_id": "4"
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("get the service s4", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/services/s4",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   serviceStrS4,
			})
		})
		It("Create a route9 using upstream id and service id", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r9",
				Body: `{
					"name": "route_all",
					"desc": "所有",
					"status": 1,
					"methods": ["GET"],
					"priority": 0,
					"service_id": "s4",
					"labels": {
						"test": "1",
						"API_VERSION": "v1"
					},
					"vars": [
						["arg_name", "==", "test"]
					],
					"uri": "/hello",
					"enable_websocket":false,
					"plugins": {
						"prometheus": {
							"disable": false
						}
					},
					"upstream_id": "5"
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("export route9", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r9",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   exportStrR9,
			})
		})
		It("delete the route9 just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r9",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("hit the route9 just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
				Sleep:        base.SleepTime,
			})
		})
		It("delete the service4", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/services/s4",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("remove upstream4", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/upstreams/4",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("remove upstream5", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/upstreams/5",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})

		// 10.Creating route10 using service ID does not contain upstream data
		serviceStrS5 := `
			"name": "testservice",
			"desc": "testservice_desc",
			"upstream_id": "6",
			"enable_websocket": true`
		serviceStrS5 = replaceStr(serviceStrS5)

		exportStrR10 := `
			"components": {},
			"info": {
				"title": "RoutesExport",
				"version": "3.0.0"
			},
			"openapi": "3.0.0",
			"paths": {
				"/hello": {
					"get": {
						"operationId": "route_allGET",
						"requestBody": {},
						"responses": {
							"default": {
								"description": ""
							}
						},
						"security": [],
						"summary": "所有",
						"x-apisix-enable_websocket": false,
						"x-apisix-id":"r10",
						"x-apisix-labels": {
							"API_VERSION": "v1",
							"test": "1"
						},
						"x-apisix-plugins": {
							"prometheus": {
								"disable": false
							}
						},
						"x-apisix-priority": 0,
						"x-apisix-status": 1,
						"x-apisix-upstream": {
							"id": "6"
						`
		exportStrR10 = replaceStr(exportStrR10)
		It("create upstream6", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/upstreams/6",
				Body: `{
					"nodes": [
						{
							"host": "` + base.UpstreamIp + `",
							"port": 1980,
							"weight": 1
						}
					],
					"type": "roundrobin"
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("create service", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/services/s5",
				Body: `{
					"name": "testservice",
					"desc": "testservice_desc",
					"enable_websocket":true,
					"upstream_id": "6"
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("get the service s5", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/services/s5",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   serviceStrS5,
			})
		})
		It("Creating route10 using service ID does not contain upstream data", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r10",
				Body: `{
					"name": "route_all",
					"desc": "所有",
					"status": 1,
					"methods": ["GET"],
					"priority": 0,
					"service_id": "s5",
					"labels": {
						"test": "1",
						"API_VERSION": "v1"
					},
					"vars": [
						["arg_name", "==", "test"]
					],
					"uri": "/hello",
					"enable_websocket":false,
					"plugins": {
						"prometheus": {
							"disable": false
						}
					}
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("export route10", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r10",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   exportStrR10,
			})
		})
		It("delete the route10 just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r10",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("hit the route10 just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
				Sleep:        base.SleepTime,
			})
		})
		It("delete the service5", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/services/s5",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("remove upstream6", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/upstreams/6",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
	})
	Context("test export route with jwt plugin", func() {
		It("make sure the route is not created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			})
		})
		It("create route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r1",
				Body: `{
					"name": "route1",
					"uri": "/hello",
					"plugins": {
						"jwt-auth": {}
					},
					"upstream": {
						"type": "roundrobin",
						"nodes": [{
							"host": "` + base.UpstreamIp + `",
							"port": 1980,
							"weight": 1
						}]
					}
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   `"code":0`,
			})
		})
		It("make sure the consumer is not created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/consumers/jack",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusNotFound,
			})
		})
		It("create consumer", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/consumers",
				Body: `{
					"username": "jack",
					"plugins": {
						"jwt-auth": {
							"key": "user-key",
							"secret": "my-secret-key",
							"algorithm": "HS256"
						}
					},
					"desc": "test description"
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})

		jwtToken := ""
		exportStrJWT := ""
		It("sign jwt token", func() {
			time.Sleep(base.SleepTime)
			// sign jwt token
			t := GinkgoT()
			body, status, err := base.HttpGet("http://127.0.0.1:9080/apisix/plugin/jwt/sign?key=user-key", nil)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, status)
			jwtToken = string(body)

			// sign jwt token with not exists key
			_, status, err = base.HttpGet("http://127.0.0.1:9080/apisix/plugin/jwt/sign?key=not-exist-key", nil)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusNotFound, status)

			exportStrJWT = `
				"components": {
					"securitySchemes": {
						"bearerAuth": {
							"bearerFormat": "JWT",
							"scheme": "bearer",
							"type": "http"
						}
					}
				},
				"info": {
					"title": "RoutesExport",
					"version": "3.0.0"
				},
				"openapi": "3.0.0",
				`
			exportStrJWT = replaceStr(exportStrJWT)
			// verify token and clean test data
		})
		It("verify route with correct jwt token", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				Headers:      map[string]string{"Authorization": jwtToken},
				ExpectStatus: http.StatusOK,
				ExpectBody:   "hello world",
				Sleep:        base.SleepTime,
			})
		})
		It("export route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   exportStrJWT,
			})
		})
		It("delete consumer", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/consumers/jack",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("verify route with the jwt token from just deleted consumer", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				Headers:      map[string]string{"Authorization": jwtToken},
				ExpectStatus: http.StatusUnauthorized,
				ExpectBody:   `{"message":"Missing related consumer"}`,
				Sleep:        base.SleepTime,
			})
		})
		It("delete route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("verify the deleted route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
				ExpectStatus: http.StatusNotFound,
				Sleep:        base.SleepTime,
			})
		})
		exportStrJWTNoAlgorithm := `
			"components": {
				"securitySchemes": {
					"bearerAuth": {
						"bearerFormat": "JWT",
						"scheme": "bearer",
						"type": "http"
					}
				}
			},
			"info": {
				"title": "RoutesExport",
				"version": "3.0.0"
			},
			"openapi": "3.0.0",
			`
		exportStrJWTNoAlgorithm = replaceStr(exportStrJWTNoAlgorithm)

		It("create consumer with jwt (no algorithm)", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/consumers",
				Body: `{
					"username":"consumer_1",
					"desc": "test description",
					"plugins":{
						"jwt-auth":{
							"exp":86400,
							"key":"user-key",
							"secret":"my-secret-key"
						}
					}
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectBody:   "\"code\":0",
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("get the consumer", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/consumers/consumer_1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   "\"username\":\"consumer_1\"",
				Sleep:        base.SleepTime,
			})
		})
		It("create the route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r1",
				Body: `{
					"name": "route1",
					"uri": "/hello",
					"plugins": {
						"jwt-auth": {}
					},
					"upstream": {
						"type": "roundrobin",
						"nodes": [{
							"host": "` + base.UpstreamIp + `",
							"port": 1980,
							"weight": 1
						}]
					}
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})

		jwttoken := ""
		It("sign jwt token", func() {
			// sign jwt token
			t := GinkgoT()
			body, status, err := base.HttpGet("http://127.0.0.1:9080/apisix/plugin/jwt/sign?key=user-key", nil)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, status)
			jwttoken = string(body)
		})
		It("hit route with jwt token", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				Headers:      map[string]string{"Authorization": jwttoken},
				ExpectStatus: http.StatusOK,
				ExpectBody:   "hello world",
				Sleep:        base.SleepTime,
			})
		})
		It("export route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   exportStrJWTNoAlgorithm,
			})
		})
		It("delete consumer", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/consumers/consumer_1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   "\"code\":0",
			})
		})
		It("after delete consumer verify it again", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/consumers/jack",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusNotFound,
				Sleep:        base.SleepTime,
			})
		})
		It("delete the route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
	})

	Context("test export route with auth plugin", func() {
		It("make sure the route is not created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
				ExpectStatus: http.StatusNotFound,
			})
		})
		It("create route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r1",
				Body: `{
					"name": "route1",
					"uri": "/hello",
					"plugins": {
						"key-auth": {},
						"basic-auth": {}
					},
					"upstream": {
						"type": "roundrobin",
						"nodes": [{
							"host": "` + base.UpstreamIp + `",
							"port": 1980,
							"weight": 1
						}]
					}
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectBody:   `"code":0`,
				ExpectStatus: http.StatusOK,
			})
		})
		It("make sure the consumer is not created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/consumers/jack",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusNotFound,
			})
		})
		It("create consumer", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/consumers",
				Body: `{
					"username": "jack",
					"plugins": {
						"key-auth": {
							"key": "auth-one"
						},
						"basic-auth": {
							"username": "jack",
							"password": "123456"
						}
					},
					"desc": "test description"
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})

		time.Sleep(base.SleepTime)

		exportStrAuth := `
			"components": {
				"securitySchemes": {
					"api_key": {
						"in": "header",
						"name": "X-XSRF-TOKEN",
						"type": "apiKey"
					},
					"basicAuth": {
						"in": "header",
						"name": "basicAuth",
						"type": "basicAuth"
					}
					}
				},
			"info": {
				"title": "RoutesExport",
				"version": "3.0.0"
			},
			"openapi": "3.0.0"`
		time.Sleep(base.SleepTime)

		exportStrAuth = replaceStr(exportStrAuth)

		It("verify route with correct basic-auth and key-auth token", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				Headers:      map[string]string{"Authorization": "Basic amFjazoxMjM0NTYKIA==", "apikey": "auth-one"},
				ExpectStatus: http.StatusOK,
				ExpectBody:   "hello world",
				Sleep:        base.SleepTime,
			})
		})
		It("export route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   exportStrAuth,
			})
		})
		It("delete consumer", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/consumers/jack",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("verify route with the basic-auth and key-auth token from just deleted consumer", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				Headers:      map[string]string{"Authorization": "Basic amFjazoxMjM0NTYKIA==", "apikey": "auth-one"},
				ExpectStatus: http.StatusUnauthorized,
				ExpectBody:   `{"message":"Missing related consumer"}`,
				Sleep:        base.SleepTime,
			})
		})
		It("delete route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("verify the deleted route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
				Sleep:        base.SleepTime,
			})
		})
	})
	Context("test route export label", func() {
		// 10.Create a service with label data and a route with label data, and export the route.
		// Label is the original data of the route
		serviceStrS1 := `
			"name": "testservice",
			"desc": "testservice_desc",
			"labels": {
				"build": "10"
			},
			"enable_websocket": true`
		serviceStrS1 = replaceStr(serviceStrS1)

		exportStrR1 := `
			"components": {},
			"info": {
				"title": "RoutesExport",
				"version": "3.0.0"
			},
			"openapi": "3.0.0",
			"paths": {
				"/hello": {
					"get": {
						"operationId": "route_allGET",
						"requestBody": {},
						"responses": {
							"default": {
								"description": ""
							}
						},
						"summary": "所有",
						"x-apisix-enable_websocket": false,
						"x-apisix-id":"r1",
						"x-apisix-labels": {
							"build": "16",
							"env": "production",
							"version": "v2"
						},
						"x-apisix-priority": 0,
						"x-apisix-status": 1,
						"x-apisix-vars":[
							[
								"arg_name",
								"==",
								"test"
								]
							]
						}`
		exportStrR1 = replaceStr(exportStrR1)

		It("create service", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:  base.ManagerApiExpect(),
				Method:  http.MethodPut,
				Path:    "/apisix/admin/services/s1",
				Headers: map[string]string{"Authorization": base.GetToken()},
				Body: `{
					"name": "testservice",
					"desc": "testservice_desc",
					"enable_websocket":true,
					"labels": {
						"build": "10"
					}
				}`,
				ExpectStatus: http.StatusOK,
			})
		})
		It("get the service s1", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:     base.ManagerApiExpect(),
				Method:     http.MethodGet,
				Path:       "/apisix/admin/services/s1",
				Headers:    map[string]string{"Authorization": base.GetToken()},
				ExpectCode: http.StatusOK,
				ExpectBody: serviceStrS1,
			})
		})
		It("Create a service with label data and a route with label data", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r1",
				Body: `{
					"name": "route_all",
					"desc": "所有",
					"status": 1,
					"methods": ["GET"],
					"service_id": "s1",
					"labels": {
						"build": "16",
						"env": "production",
						"version": "v2"
					},
					"vars": [
						["arg_name", "==", "test"]
					],
					"uri": "/hello"
				}`,
				Headers:    map[string]string{"Authorization": base.GetToken()},
				ExpectCode: http.StatusOK,
				Sleep:      base.SleepTime,
			})
		})
		It("export route1", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:     base.ManagerApiExpect(),
				Method:     http.MethodGet,
				Path:       "/apisix/admin/export/routes/r1",
				Headers:    map[string]string{"Authorization": base.GetToken()},
				ExpectCode: http.StatusOK,
				ExpectBody: exportStrR1,
			})
		})
		It("delete the route1 just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:     base.ManagerApiExpect(),
				Method:     http.MethodDelete,
				Path:       "/apisix/admin/routes/r1",
				Headers:    map[string]string{"Authorization": base.GetToken()},
				ExpectCode: http.StatusOK,
			})
		})
		It("hit the route1 just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
				Sleep:        base.SleepTime,
			})
		})
		It("delete the service1", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/services/s1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})

		// 11.Create a service with label data and a route without label data, and export the route.
		//  Label is the data of the service
		serviceStrS2 := `
			"name": "testservice",
			"desc": "testservice_desc",
			"labels": {
				"build": "16",
				"env": "production",
				"version": "v2"
			},
			"enable_websocket": true`
		serviceStrS2 = replaceStr(serviceStrS2)

		exportStrR2 := `
			"components": {},
			"info": {
				"title": "RoutesExport",
				"version": "3.0.0"
			},
			"openapi": "3.0.0",
			"paths": {
				"/hello": {
					"get": {
						"operationId": "route_allGET",
						"requestBody": {},
						"responses": {
							"default": {
								"description": ""
							}
						},
						"summary": "所有",
						"x-apisix-enable_websocket": false,
						"x-apisix-id":"r2",
						"x-apisix-labels": {
							"build": "16",
							"env": "production",
							"version": "v2"
						},
						"x-apisix-priority": 0,
						"x-apisix-status": 1,
						"x-apisix-vars":[
							[
								"arg_name",
								"==",
								"test"
								]
							]
						}`
		exportStrR2 = replaceStr(exportStrR2)

		It("create service", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:  base.ManagerApiExpect(),
				Method:  http.MethodPut,
				Path:    "/apisix/admin/services/s2",
				Headers: map[string]string{"Authorization": base.GetToken()},
				Body: `{
					"name": "testservice",
					"desc": "testservice_desc",
					"enable_websocket":true,
					"labels": {
						"build": "16",
						"env": "production",
						"version": "v2"
					}
				}`,
				ExpectStatus: http.StatusOK,
			})
		})
		It("get the service s2", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/services/s2",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectBody:   serviceStrS2,
				ExpectStatus: http.StatusOK,
			})
		})
		It("Create a service with label data and a route without label data", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r2",
				Body: `{
					"name": "route_all",
					"desc": "所有",
					"status": 1,
					"methods": ["GET"],
					"service_id": "s2",
					"vars": [
						["arg_name", "==", "test"]
					],
					"uri": "/hello"
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("export route2", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r2",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectBody:   exportStrR2,
				ExpectStatus: http.StatusOK,
			})
		})
		It("delete the route2 just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r2",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("hit the route2 just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
				Sleep:        base.SleepTime,
			})
		})
		It("delete the service2", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/services/s2",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
	})
	Context("test route export request validation", func() {
		// 12.Test export route request_ validation data correctness
		exportStrR1 := `
			"components": {},
			"info": {
				"title": "RoutesExport",
				"version": "3.0.0"
			},
			"openapi": "3.0.0",
			"paths": {
				"/test-test": {
					"get": {
						"operationId": "route_allGET",
						"parameters": [{
							"in": "header",
							"name": "test",
							"schema": {
								"type": "string"
							}
						}],
						"requestBody": {
							"content": {
								"*/*": {
									"schema": {
										"properties": {
											"boolean_payload": {
												"type": "boolean"
											},
											"required_payload": {
												"type": "string"
											}
										},
										"required": ["required_payload"],
										"type": "object"
									}
								}
							}
						},
						"responses": {
							"default": {
								"description": ""
							}
						},
						"security": [],
						"summary": "所有",
						"x-apisix-enable_websocket": false,
						"x-apisix-hosts": ["test.com"],
						"x-apisix-id":"r1",
						"x-apisix-priority": 0,
						"x-apisix-status": 1
					}
				}
			}`
		exportStrR1 = replaceStr(exportStrR1)

		It("Create a route containing request_ validation data", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r1",
				Body: `{
					"uris": ["/test-test"],
					"name": "route_all",
					"desc": "所有",
					"methods": ["GET"],
					"hosts": ["test.com"],
					"plugins": {
						"request-validation": {
							"body_schema": {
								"properties": {
									"boolean_payload": {
										"type": "boolean"
									},
									"required_payload": {
										"type": "string"
									}
								},
								"required": ["required_payload"],
								"type": "object"
							},
							"disable": false,
							"header_schema": {
								"properties": {
									"test": {
										"enum": "test-enum",
										"type": "string"
									}
								},
								"type": "string"
							}
						}
					},
					"status": 1
			 	}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("export route1", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   exportStrR1,
			})
		})
		It("delete the route1 just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("hit the route1 just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
				Sleep:        base.SleepTime,
			})
		})
	})
	Context("test route export equal uri", func() {
		// 13.Add suffix when testing the same URI export
		exportStrAll := `
			"components": {},
			"info": {
				"title": "RoutesExport",
				"version": "3.0.0"
			},
			"openapi": "3.0.0",
			"paths": {
				"/test-test": {
					"get": {
						"operationId": "route_allGET",
						"requestBody": {},
						"responses": {
							"default": {
								"description": ""
							}
						},
						"summary": "所有",
						"x-apisix-enable_websocket": false,
						"x-apisix-hosts": ["test.com"],
						"x-apisix-id":"r1",
						"x-apisix-priority": 0,
						"x-apisix-status": 1,
						"x-apisix-upstream": {
							"nodes": {
								"` + base.UpstreamIp + `:1980": 1
							},
							"type": "roundrobin"
						}
					}
				},
				"/test-test-APISIX-REPEAT-URI-2": {
					"get": {
						"operationId": "route_all2GET",
						"requestBody": {},
						"responses": {
							"default": {
								"description": ""
							}
						},
						"summary": "所有1",
						"x-apisix-enable_websocket": false,
						"x-apisix-hosts": ["test.com"],
						"x-apisix-id":"r2",
						"x-apisix-priority": 0,
						"x-apisix-status": 1,
						"x-apisix-upstream": {
							"nodes": {
								"` + base.UpstreamIp + `:1980": 1
							},
							"type": "roundrobin"
						}
					}
				},
				"/test-test-APISIX-REPEAT-URI-3": {
					"get": {
						"operationId": "route_all3GET",
						"requestBody": {},
						"responses": {
							"default": {
								"description": ""
							}
						},
						"summary": "所有2",
						"x-apisix-enable_websocket": false,
						"x-apisix-hosts": ["test.com"],
						"x-apisix-id":"r3",
						"x-apisix-priority": 0,
						"x-apisix-status": 1,
						"x-apisix-upstream": {
							"nodes": {
								"` + base.UpstreamIp + `:1981": 1
							},
							"type": "roundrobin"
						}
					}
				}
			}`
		exportStrAll = replaceStr(exportStrAll)

		It("Create a route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r1",
				Body: `{
				"uris": ["/test-test"],
				"name": "route_all",
				"desc": "所有",
				"methods": ["GET"],
				"hosts": ["test.com"],
				"status": 1,
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
						},
					"type": "roundrobin"
				 	}
		 		}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("Create a route2", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r2",
				Body: `{
					"uris": ["/test-test"],
					"name": "route_all2",
					"desc": "所有1",
					"methods": ["GET"],
					"hosts": ["test.com"],
					"status": 1,
					"upstream": {
						"nodes": {
							"` + base.UpstreamIp + `:1980": 1
						},
						"type": "roundrobin"
					}
			 	}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("Create a route3", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r3",
				Body: `{
					"uris": ["/test-test"],
					"name": "route_all3",
					"desc": "所有2",
					"methods": ["GET"],
					"hosts": ["test.com"],
					"status": 1,
					"upstream": {
						"nodes": {
							"` + base.UpstreamIp + `:1981": 1
						},
						"type": "roundrobin"
					}
			 	}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("use the exportall interface to export all routes", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/export/routes/r1,r2,r3",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   exportStrAll,
			})
		})
		It("delete the route1 just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("hit the route1 just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:     base.APISIXExpect(),
				Method:     http.MethodGet,
				Path:       "/hello",
				Headers:    map[string]string{"Authorization": base.GetToken()},
				ExpectBody: "{\"error_msg\":\"404 Route Not Found\"}\n",
				Sleep:      base.SleepTime,
			})
		})
		It("delete the route2 just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r2",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("hit the route2 just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
				Sleep:        base.SleepTime,
			})
		})
	})
})

func replaceStr(str string) string {
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, " ", "", -1)
	return str
}
