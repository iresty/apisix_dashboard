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
package service_test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Service", func() {
	Describe("Test service Integrate (With plugin)", func() {
		It("Create service", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:  base.ManagerApiExpect(),
				Method:  http.MethodPut,
				Path:    "/apisix/admin/services/s1",
				Headers: map[string]string{"Authorization": base.GetToken()},
				Body: `{
					"name": "testservice",
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
						"nodes": {
							"` + base.UpstreamIp + `:1980": 1
						}
					}
				}`,
				ExpectStatus: http.StatusOK,
				ExpectBody:   []string{`"id":"s1"`, `"name":"testservice"`},
			})
		})
		It("Create route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r1",
				Body: `{
					"name": "route1",
					"uri": "/server_port",
					"service_id": "s1"
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        time.Second,
			})
		})
		It("Hit routes and check the response header", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/server_port",
				ExpectStatus: http.StatusOK,
				ExpectHeaders: map[string]string{
					"X-Ratelimit-Limit":     "100",
					"X-Ratelimit-Remaining": "99",
				},
			})
		})
		It("Delete route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("Delete service", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/services/s1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})
		It("Hit the route just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/server_port",
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			})
		})
	})
})
