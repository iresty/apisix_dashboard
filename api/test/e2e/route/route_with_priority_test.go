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

var _ = Describe("route with priority test", func() {
	DescribeTable("test route with priority",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("add another route with no priority (default 0)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
					"name": "route1",
					"uri": "/server_port",
					"methods": ["GET"],
					"upstream": {
						"type": "roundrobin",
						"nodes": {
							"` + base.UpstreamIp + `:1981": 1
						}
					}
				}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("access the route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/server_port",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "1981",
			Sleep:        base.SleepTime,
		}),
		Entry("add another route with valid priority (1), upstream is different from the others", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r2",
			Body: `{
					"name": "route2",
					"uri": "/server_port",
					"methods": ["GET"],
					"priority": 1,
					"upstream": {
						"type": "roundrobin",
						"nodes": {
							"` + base.UpstreamIp + `:1982": 1
						}
					}
				}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("access the route to determine whether it meets the priority (compare 1 and default)", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/server_port",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "1982",
			Sleep:        base.SleepTime,
		}),
		Entry("delete route (r1)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		}),
		Entry("delete route (r2)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r2",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		}),
		Entry("hit the route just delete", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/server_port",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		}),
	)
})
