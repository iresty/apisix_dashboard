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
package route

import (
	"net/http"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = ginkgo.Describe("Route", func() {
	ginkgo.Context("test route plugin skywalking", func() {
		ginkgo.It("make sure the route is not created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			})
		})
		ginkgo.It("create route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r1",
				Body: `{
					"name": "route1",
					"uri": "/hello",
					"plugins": {
						"skywalking": {
							"sample_ratio": 1
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
			})
		})
		ginkgo.It("trigger skywalking", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectStatus: http.StatusOK,
				ExpectBody:   "hello world",
				Sleep:        base.SleepTime,
			})
		})
		ginkgo.It("verify the log", func() {
			// sleep for process log
			time.Sleep(4 * time.Second)
			t := ginkgo.GinkgoT()

			// verify by checking log
			logContent := base.ReadAPISIXErrorLog()
			assert.Contains(t, logContent, "segments reported")

			// clean log
			base.CleanAPISIXErrorLog()
		})
		ginkgo.It("create route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r1",
				Body: `{
					"name": "route1",
					"uri": "/hello",
					"plugins": {
						"skywalking": {
							"sample_ratio": 0.00001
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
			})
		})
		ginkgo.It("trigger skywalking", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectStatus: http.StatusOK,
				ExpectBody:   "hello world",
				Sleep:        base.SleepTime,
			})
		})
		ginkgo.It("verify the log", func() {

			// sleep for process log
			time.Sleep(4 * time.Second)
			t := ginkgo.GinkgoT()

			// verify by checking log
			logContent := base.ReadAPISIXErrorLog()
			assert.Contains(t, logContent, "miss sampling, ignore")

			// clean log
			base.CleanAPISIXErrorLog()
		})
		ginkgo.It("delete route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/r1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		ginkgo.It("make sure the route has been deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello_",
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
				Sleep:        base.SleepTime,
			})
		})
	})
})
