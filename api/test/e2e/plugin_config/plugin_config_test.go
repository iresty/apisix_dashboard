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
package plugin_config_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Plugin Config", func() {
	DescribeTable("Test plugin config CURD",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("Get plugin config (Not Exist)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/plugin_configs/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		}),
		Entry("List plugin config (Empty)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/plugin_configs",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":0`,
		}),
		Entry("Create plugin config #1", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/plugin_configs/1",
			Body: `{
				"plugins": {
					"response-rewrite": {
						"headers": {
							"X-VERSION":"1.0"
						}
					},
					"uri-blocker": {
						"block_rules": ["select.+(from|limit)", "(?:(union(.*?)select))"]
					}
				},
				"labels": {
					"version": "v1",
					"build":   "16"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"X-VERSION":"1.0"`,
		}),
		Entry("Get plugin config (Exist)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/plugin_configs/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"X-VERSION":"1.0"`,
		}),
		Entry("List plugin config (1 item)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/plugin_configs",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":1`,
		}),
		Entry("Create plugin config #2", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/plugin_configs/2",
			Body: `{
				"id": "2",
				"plugins": {
					"response-rewrite": {
						"headers": {
							"X-VERSION":"2.0"
						}
					}
				},
				"labels": {
					"version": "v2",
					"build":   "17",
					"extra":   "test"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"X-VERSION":"2.0"`,
		}),
		Entry("Search plugin config (By Label key:value)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Path:         "/apisix/admin/plugin_configs",
			Query:        "label=build:16",
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"labels":{"build":"16","version":"v1"}`,
		}),
		Entry("Search plugin config (By Label key)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Path:         "/apisix/admin/plugin_configs",
			Query:        "label=extra",
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"labels":{"build":"17","extra":"test","version":"v2"}`,
		}),
		Entry("Update plugin config (Full)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/plugin_configs/1",
			Body: `{
				"plugins": {
					"response-rewrite": {
						"headers": {
							"X-TEST":"1.0"
						}
					},
					"uri-blocker": {
						"block_rules": ["select.+(from|limit)", "(?:(union(.*?)select))"]
					}
				},
				"labels": {
					"version": "v1",
					"build":   "16"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"X-TEST":"1.0"`,
		}),
		Entry("Update plugin config (Partial)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPatch,
			Path:   "/apisix/admin/plugin_configs/1/plugins",
			Body: `{
				"response-rewrite": {
					"headers": {
						"X-TEST":"2.0"
					}
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"X-TEST":"2.0"`,
		}),
		Entry("Batch Delete plugin config", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/plugin_configs/1,2",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
	)

	DescribeTable("Test plugin config Integration",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("Create plugin config", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Path:   "/apisix/admin/plugin_configs/1",
			Method: http.MethodPut,
			Body: `{
				"plugins": {
					"response-rewrite": {
						"headers": {
							"X-VERSION":"1.0"
						}
					},
					"uri-blocker": {
						"block_rules": ["select.+(from|limit)", "(?:(union(.*?)select))"]
					}
				},
				"labels": {
					"version": "v1",
					"build":   "16"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("Create route", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				 "name": "route1",
				 "uri": "/hello",
				 "plugin_config_id": "1",
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
		}),
		Entry("Verify (Add Header)", base.HttpTestCase{
			Object:        base.APISIXExpect(),
			Method:        http.MethodGet,
			Path:          "/hello",
			ExpectStatus:  http.StatusOK,
			ExpectBody:    "hello world",
			ExpectHeaders: map[string]string{"X-VERSION": "1.0"},
		}),
		Entry("Verify (URI Block)", base.HttpTestCase{
			Object:        base.APISIXExpect(),
			Method:        http.MethodGet,
			Path:          "/hello",
			Query:         "name=%3Bselect%20from%20sys",
			ExpectStatus:  http.StatusForbidden,
			ExpectHeaders: map[string]string{"X-VERSION": "1.0"},
		}),
		Entry("Delete route", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("Delete plugin config", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/plugin_configs/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
	)
})
