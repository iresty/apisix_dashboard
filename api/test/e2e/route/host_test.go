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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Route", func() {
	DescribeTable("test route with host",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("invalid host", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Path:   "/apisix/admin/routes/r1",
			Method: http.MethodPut,
			Body: `{
				"name": "route1",
				"uri": "/hello_",
				"host": "$%$foo.com",
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		Entry("invalid hosts", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello_",
				"hosts": ["$%$foo.com", "*.bar.com"],
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		Entry("create route with host and hosts together at the same time", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello_",
				"host": "github.com",
				"hosts": ["foo.com", "*.bar.com"],
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		Entry("hit route not created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			Headers:      map[string]string{"Host": "foo.com"},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		}),
		Entry("hit route not created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			Headers:      map[string]string{"Host": "$%$foo.com"},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		}),
	)

	DescribeTable("test route with hosts",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("make sure route not created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			Headers:      map[string]string{"Host": "foo.com"},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		}),
		Entry("create route", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello_",
				"hosts": ["foo.com", "*.bar.com"],
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("create route with int uri", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route2",
				"uri": 123456
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		Entry("hit the route just created - wildcard domain name", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			Headers:      map[string]string{"Host": "test.bar.com"},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world\n",
			Sleep:        base.SleepTime,
		}),
		Entry("hit the route just created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			Headers:      map[string]string{"Host": "foo.com"},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world\n",
		}),
		Entry("hit the route not exists", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_111",
			Headers:      map[string]string{"Host": "foo.com"},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		}),
		Entry("delete the route just created", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("hit the route just deleted", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			Headers:      map[string]string{"Host": "bar.com"},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		}),
	)

	DescribeTable("update routes with hosts",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("hit route that not exist", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Host": "foo.com"},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		}),
		Entry("create route with host foo.com", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
			"name": "route1",
			"uri": "/hello",
			"methods": ["GET"],
			"hosts": ["foo.com"],
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
			ExpectBody:   []string{"\"id\":\"r1\"", "\"hosts\":[\"foo.com\"]"},
		}),
		Entry("hit the route just create", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Host": "foo.com"},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		}),
		Entry("update route with host bar.com", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"hosts": ["bar.com"],
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"id\":\"r1\"", "\"hosts\":[\"bar.com\"]"},
		}),
		Entry("hit the route with host foo.com", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Host": "foo.com"},
			ExpectStatus: http.StatusNotFound,
			Sleep:        base.SleepTime,
		}),
		Entry("hit the route just updated", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Host": "bar.com"},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world\n",
		}),
		Entry("delete route", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("hit the route just deleted", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Host": "bar.com"},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		}),
	)

	DescribeTable("test route with empty array",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("create route with empty hosts and host", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Path:   "/apisix/admin/routes/r1",
			Method: http.MethodPut,
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"hosts": [],
				"host": "test.com",
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   `{"code":10000,"message":"schema validate failed: (root): Must validate one and only one schema (oneOf)\n(root): Must validate all the schemas (allOf)\nhosts: Array must have at least 1 items"}`,
		}),
		Entry("make sure the route not created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Host": "test.com"},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
		}),
		Entry("create route with empty hosts", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Path:   "/apisix/admin/routes/r1",
			Method: http.MethodPut,
			Body: `{
				"name": "route2",
				"uri": "/hello",
				"hosts": [],
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   `{"code":10000,"message":"schema validate failed: hosts: Array must have at least 1 items"}`,
		}),
		Entry("make sure the route not created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
		}),
		Entry("create route with empty uris and uri", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Path:   "/apisix/admin/routes/r1",
			Method: http.MethodPut,
			Body: `{
				"name": "route2",
				"uri": "/hello",
				"uris": [],
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   `{"code":10000,"message":"schema validate failed: (root): Must validate one and only one schema (oneOf)\n(root): Must validate all the schemas (allOf)\nuris: Array must have at least 1 items"}`,
		}),
		Entry("create route with empty remote_addrs and remote_addr", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Path:   "/apisix/admin/routes/r1",
			Method: http.MethodPut,
			Body: `{
				"name": "route2",
				"uri": "/hello",
				"remote_addrs": [],
				"remote_addr": "0.0.0.0",
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   `{"code":10000,"message":"schema validate failed: (root): Must validate one and only one schema (oneOf)\n(root): Must validate all the schemas (allOf)\nremote_addrs: Array must have at least 1 items"}`,
		}),
		Entry("make sure the route not created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
		}),
	)

})
