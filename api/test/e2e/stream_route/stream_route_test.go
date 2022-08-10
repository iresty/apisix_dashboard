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
package stream_route_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Stream Route", func() {
	DescribeTable("Test stream route CURD",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("Get stream route (Not Exist)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/stream_routes/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		}),
		Entry("List stream route (Empty)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/stream_routes",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":0`,
		}),
		Entry("Create stream route #1", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/stream_routes/1",
			Body: `{
				"remote_addr": "127.0.0.1",
				"server_addr": "127.0.0.1",
				"server_port": 10090,
				"sni": "test.com",
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"remote_addr":"127.0.0.1"`,
		}),
		Entry("Get stream route (Exist)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/stream_routes/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"server_port":10090`,
		}),
		Entry("List stream route (1 item)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/stream_routes",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":1`,
		}),
		Entry("Create stream route #2", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/stream_routes",
			Body: `{
				"id": "2",
				"remote_addr": "127.0.0.1",
				"server_addr": "127.0.0.1",
				"server_port": 10091,
				"sni": "test.com",
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"server_port":10091`,
		}),
		Entry("List stream route (2 items)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/stream_routes",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":2`,
		}),
		Entry("List stream route (Paginate)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/stream_routes",
			Query:        "page=2&page_size=1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"id":"2"`,
		}),
		Entry("Update stream route (Full)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/stream_routes/1",
			Body: `{
				"remote_addr": "127.0.0.1",
				"server_addr": "127.0.0.1",
				"server_port": 10091,
				"sni": "test.com",
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"server_port":10091`,
		}),
		Entry("Batch Delete stream route", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/stream_routes/1,2",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
	)

	DescribeTable("Test stream route CURD exception",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("Create stream route with upstream id not found", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/stream_routes",
			Body: `{
				"id": "sr1",
				"server_port": 10090,
				"upstream_id": "u1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
	)

	It("Clean all resources", func() {
		base.CleanResource("stream_routes")
		base.CleanResource("ssl")
	})
})
