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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("route with limit plugin", func() {
	DescribeTable("test route with limit plugin",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("make sure the route is not created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
		}),
		Entry("create route with limit plugin", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"plugins": {
					"limit-count": {
						"count": 2,
						"time_window": 2,
						"rejected_code": 503,
						"key": "remote_addr",
						"policy": "local"
					}
				},
				"upstream": {
					"type": "roundrobin",
					"nodes": {
						"` + base.UpstreamIp + `:1981": 1
					}
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"code":0`,
		}),
		Entry("verify route that should not be limited", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		}),
		Entry("verify route that should not be limited 2", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
		}),
		Entry("verify route that should be limited", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusServiceUnavailable,
			ExpectBody:   "503 Service Temporarily Unavailable",
		}),
		Entry("verify route that should not be limited since time window pass", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        3 * time.Second,
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
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		}),
	)
})

var _ = Describe("route with limit plugin by consumer", func() {
	DescribeTable("test route with limit plugin by consumer",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("make sure the route is not created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
		}),
		Entry("create route with limit plugin", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
			   "name": "route1",
			   "uri": "/hello",
			   "plugins": {
				"key-auth": {},
				   "limit-count": {
					   "count": 2,
					   "time_window": 2,
					   "rejected_code": 503,
					   "key": "consumer_name",
					   "policy": "local"
				   }
			   },
			   "upstream": {
				   "type": "roundrobin",
				   "nodes": {
					   "` + base.UpstreamIp + `:1981": 1
				   }
			   }
		   }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"code":0`,
		}),
		Entry("make sure the consumer is not created", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		}),
		Entry("create consumer", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Path:   "/apisix/admin/consumers",
			Method: http.MethodPut,
			Body: `{
			"username": "jack",
			"plugins": {
				"key-auth": {
					"key": "auth-jack"
				}
			}
		}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("create consumer 2", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Path:   "/apisix/admin/consumers",
			Method: http.MethodPut,
			Body: `{
			"username": "pony",
			"plugins": {
				"key-auth": {
					"key": "auth-pony"
				}
			}
		}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("verify route that should not be limited", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"apikey": "auth-jack"},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		}),
		Entry("verify route that should not be limited 2", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"apikey": "auth-jack"},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
		}),
		Entry("verify route that should be limited", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"apikey": "auth-jack"},
			ExpectStatus: http.StatusServiceUnavailable,
			ExpectBody:   "503 Service Temporarily Unavailable",
		}),
		Entry("verify route that should not be limited (other consumer)", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"apikey": "auth-pony"},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
		}),
		Entry("verify route that should not be limited since time window pass", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"apikey": "auth-jack"},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        2 * time.Second,
		}),
		Entry("delete consumer pony", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/consumers/pony",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("delete consumer jack", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("make sure pony has been deleted", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"apikey": "auth-pony"},
			ExpectStatus: http.StatusUnauthorized,
			ExpectBody:   `{"message":"Missing related consumer"}`,
			Sleep:        base.SleepTime,
		}),
		Entry("make sure jack has been deleted", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"apikey": "auth-jack"},
			ExpectStatus: http.StatusUnauthorized,
			ExpectBody:   `{"message":"Missing related consumer"}`,
			Sleep:        base.SleepTime,
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
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		}),
	)
})

var _ = Describe("route with limit count and disable", func() {
	DescribeTable("test route with limit count and disable",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("make sure the route is not created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
		}),
		Entry("create route with limit plugin", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
			   "name": "route1",
			   "uri": "/hello",
			   "plugins": {
				   "limit-count": {
					   "count": 2,
					   "time_window": 2,
					   "rejected_code": 503,
					   "key": "remote_addr",
					   "disable": false,
					   "policy": "local"
				   }
			   },
			   "upstream": {
				   "type": "roundrobin",
				   "nodes": {
					   "` + base.UpstreamIp + `:1981": 1
				   }
			   }
		   }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"code":0`,
		}),
		Entry("verify route that should not be limited", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		}),
		Entry("verify route that should not be limited 2", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
		}),
		Entry("verify route that should be limited", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusServiceUnavailable,
			ExpectBody:   "503 Service Temporarily Unavailable",
		}),
		Entry("verify route that should not be limited since time window pass", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        2 * time.Second,
		}),
		Entry("update route to disable plugin limit-count", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
			   "name": "route1",
			   "uri": "/hello",
			   "plugins": {
				   "limit-count": {
					   "count": 1,
					   "time_window": 2,
					   "rejected_code": 503,
					   "key": "remote_addr",
					   "disable": true,
					   "policy": "local"
				   }
			   },
			   "upstream": {
				   "type": "roundrobin",
				   "nodes": {
					   "` + base.UpstreamIp + `:1981": 1
				   }
			   }
		   }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"code":0`,
		}),
		Entry("verify route that should not be limited", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		}),
		Entry("verify route that should not be limited (exceed config count)", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
		}),
		Entry("verify route that should not be limited (exceed config count again)", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
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
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		}),
	)
})
