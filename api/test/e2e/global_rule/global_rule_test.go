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
package global_rule_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Global Rule", func() {
	DescribeTable("Test global rule CURD",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("Get global rule (Not Exist)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/global_rules/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		}),
		Entry("List global rule (Empty)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/global_rules",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":0`,
		}),
		Entry("Create global rule #1", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/global_rules/1",
			Body: `{
				"plugins": {
					"response-rewrite": {
						"headers": {
							"X-VERSION":"1.0"
						}
					}
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"X-VERSION":"1.0"`,
		}),
		Entry("Get global rule (Exist)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/global_rules/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"X-VERSION":"1.0"`,
		}),
		Entry("List global rule (1 item)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/global_rules",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":1`,
		}),
		Entry("Create global rule #2", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/global_rules/2",
			Body: `{
				"plugins": {
					"response-rewrite": {
						"headers": {
							"X-VERSION":"2.0"
						}
					}
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"X-VERSION":"2.0"`,
		}),
		Entry("List global rule (2 item)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/global_rules",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":2`,
		}),
		Entry("List global rule (Paginate)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/global_rules",
			Query:        "page=2&page_size=1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"id":"2"`,
		}),
		Entry("Update global rule (Full)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/global_rules/1",
			Body: `{
				"plugins": {
					"response-rewrite": {
						"headers": {
							"X-TEST":"1.0"
						}
					}
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "\"X-TEST\":\"1.0\"",
		}),
		Entry("Update global rule (Partial)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPatch,
			Path:   "/apisix/admin/global_rules/1/plugins",
			Body: `{
				"response-rewrite": {
					"headers": {
						"X-VERSION":"1.0"
					}
				},
				"key-auth": {}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "\"key-auth\":{}",
		}),
		Entry("Delete global rule", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/global_rules/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
	)

	DescribeTable("Test global rule Integration",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("Create route", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello",
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
		Entry("Create global rule", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/global_rules/1",
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
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "\"X-VERSION\":\"1.0\"",
		}),
		Entry("Verify route with header", base.HttpTestCase{
			Object:        base.APISIXExpect(),
			Method:        http.MethodGet,
			Path:          "/hello",
			ExpectStatus:  http.StatusOK,
			ExpectBody:    "hello world",
			ExpectHeaders: map[string]string{"X-VERSION": "1.0"},
		}),
		Entry("Verify route that should be blocked", base.HttpTestCase{
			Object:        base.APISIXExpect(),
			Method:        http.MethodGet,
			Path:          "/hello",
			Query:         "name=%3Bselect%20from%20sys",
			ExpectStatus:  http.StatusForbidden,
			ExpectHeaders: map[string]string{"X-VERSION": "1.0"},
		}),
		Entry("Update route with same plugin response-rewrite", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"plugins": {
					"response-rewrite": {
						"headers": {
							"X-VERSION":"2.0"
						}
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
			ExpectBody:   "\"X-VERSION\":\"2.0\"",
		}),
		Entry("Verify route that header should override", base.HttpTestCase{
			Object:        base.APISIXExpect(),
			Method:        http.MethodGet,
			Path:          "/hello",
			ExpectStatus:  http.StatusOK,
			ExpectBody:    "hello world",
			ExpectHeaders: map[string]string{"X-VERSION": "2.0"},
		}),
		Entry("Delete global rule", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/global_rules/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("Delete route", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
	)
})
