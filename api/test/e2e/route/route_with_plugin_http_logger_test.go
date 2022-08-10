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
	. "github.com/onsi/gomega"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("route with plugin http logger", func() {
	It("cleanup previous error logs", func() {
		base.CleanAPISIXErrorLog()
	})

	DescribeTable("test route with http logger plugin",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("make sure the route is not created ", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
		}),
		Entry("create route", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello_",
				"plugins": {
					"http-logger": {
						"uri": "http://` + base.UpstreamIp + `:1982/hello",
						"batch_max_size": 1,
						"max_retry_count": 1,
						"retry_delay": 2,
						"buffer_duration": 2,
						"inactive_timeout": 2,
						"name": "http logger",
						"timeout": 3,
						"concat_method": "json"
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
			ExpectBody:   []string{`"code":0`, `"id":"r1"`, `"uri":"/hello_"`, `"name":"route1"`, `"name":"http logger"`},
			ExpectStatus: http.StatusOK,
		}),
		Entry("access route to trigger log", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		}),
	)

	It("verify http logger by checking log", func() {
		// sleep for process log
		time.Sleep(1500 * time.Millisecond)

		// verify http logger by checking log
		//todo: should use a fake upstream for confirming whether we got the log data.
		logContent := base.ReadAPISIXErrorLog()
		Expect(logContent).Should(ContainSubstring("Batch Processor[http logger] successfully processed the entries"))

		// clean log
		base.CleanAPISIXErrorLog()
	})

	DescribeTable("test route for unreachable logger uri",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("create route with wrong https endpoint", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r2",
			Body: `{
				"name": "route2",
				"uri": "/hello",
				"plugins": {
					"http-logger": {
						"uri": "https://127.0.0.1:8888/hello-world-http",
						"batch_max_size": 1,
						"max_retry_count": 1,
						"retry_delay": 2,
						"buffer_duration": 2,
						"inactive_timeout": 2,
						"name": "http logger",
						"timeout": 3,
						"concat_method": "json"
					}
				},
				"upstream": {
					"type": "roundrobin",
					"nodes": {
							"` + base.UpstreamIp + `:1982": 1
						}
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectBody:   []string{`"code":0`, `"id":"r2"`, `"uri":"/hello"`, `"name":"route2"`, `"name":"http logger"`},
			ExpectStatus: http.StatusOK,
		}),
		Entry("access route to trigger log", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		}),
	)

	It("verify http logger by checking log for second route", func() {
		// sleep for process log
		time.Sleep(1500 * time.Millisecond)

		// verify http logger by checking log
		//todo: should use a fake upstream for confirming whether we got the log data.
		logContent := base.ReadAPISIXErrorLog()
		Expect(logContent).Should(ContainSubstring("Batch Processor[http logger] failed to process entries: failed to connect to host[127.0.0.1] port[8888] connection refused"))

		// clean log
		base.CleanAPISIXErrorLog()
	})

	// todo: check disable http logger - Done
	DescribeTable("rechecking logger after disabling plugin",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("disable route http logger plugin", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r2",
			Body: `{
				"name": "route2",
				"uri": "/hello",
				"plugins": {},
				"upstream": {
					"type": "roundrobin",
					"nodes": {
							"` + base.UpstreamIp + `:1982": 1
						}
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectBody:   []string{`"code":0`, `"id":"r2"`, `"uri":"/hello"`, `"name":"route2"`},
			ExpectStatus: http.StatusOK,
		}),
		Entry("access route to trigger log (though should not be triggered)", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		}),
	)

	It("verify http logger has been successfully disabled", func() {
		// sleep for process log
		time.Sleep(1500 * time.Millisecond)

		// verify http logger by checking log
		logContent := base.ReadAPISIXErrorLog()
		Expect(logContent).ShouldNot(ContainSubstring("Batch Processor[http logger] successfully processed the entries"))

		// clean log
		base.CleanAPISIXErrorLog()
	})

	DescribeTable("cleanup test data",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("delete route", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("make sure the route has been deleted", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		}),
		Entry("delete route 2", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r2",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("make sure the route 2 has been deleted", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		}),
	)
})
