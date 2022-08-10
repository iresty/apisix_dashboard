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

var _ = Describe("route with plugin prometheus", func() {
	DescribeTable("test route with plugin prometheus",
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

		Entry("create route with plugin prometheus", base.HttpTestCase{
			Desc:   "create route",
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"plugins": {
					"prometheus": {}
				},
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
		Entry("fetch the prometheus metric data", base.HttpTestCase{
			Object:       base.PrometheusExporterExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/prometheus/metrics",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "apisix_etcd_reachable 1",
			Sleep:        base.SleepTime,
		}),
		Entry("request from client (200)", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
		}),
		Entry("create route that uri not exists in upstream", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello-not-exists",
				"plugins": {
					"prometheus": {}
				},
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
		Entry("request from client (404)", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello-not-exists",
			ExpectStatus: http.StatusNotFound,
			Sleep:        base.SleepTime,
		}),
		Entry("verify the prometheus metric data (apisix_http_status 200)", base.HttpTestCase{
			Object:       base.PrometheusExporterExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/prometheus/metrics",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `apisix_http_status{code="200",route="r1",matched_uri="/hello",matched_host="",service="",consumer=""`,
			Sleep:        1 * time.Second,
		}),
		Entry("verify the prometheus metric data (apisix_http_status 404)", base.HttpTestCase{
			Object:       base.PrometheusExporterExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/prometheus/metrics",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `apisix_http_status{code="404",route="r1",matched_uri="/hello-not-exists",matched_host="",service="",consumer=""`,
			Sleep:        base.SleepTime,
		}),
		Entry("delete route", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		}),
		Entry("make sure the route deleted", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		}),
	)
})
