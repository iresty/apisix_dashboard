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

var _ = Describe("create route that not exists service or upstream", func() {
	DescribeTable("test create route that not exists service or upstream",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("make sure the route has not created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		}),
		Entry("create route that not exists service", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello_",
				"service_id": "not-exists"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		Entry("verify not-exist route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		}),
		Entry("create route that not exists upstream", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello_",
				"upstream_id": "not-exists"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		Entry("verify not-exist route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		}),
		Entry("create route that not exists service and upstream", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello_",
				"service_id": "not-exists-service",
				"upstream_id": "not-exists-upstream"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		Entry("verify not-exist route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		}),
		Entry("create service with not-exist upstream", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/services/100",
			Body: `{
				"upstream_id": "not-exists-upstream"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		Entry("create route with service(service with not exist upstream)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello_",
				"service_id": "not-exists"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		Entry("verify not-exist route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		}),
	)
})

var _ = Describe("route create with service", func() {
	DescribeTable("test route create with service",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("make sure the route has not created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/server_port",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		}),
		Entry("create service", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/services/200",
			Body: `{
				"upstream": {
				"type": "roundrobin",
				"nodes": [
						{
							"host": "` + base.UpstreamIp + `",
							"port": 1980,
							"weight": 1
						},
						{
							"host": "` + base.UpstreamIp + `",
							"port": 1981,
							"weight": 1
						},
						{
							"host": "` + base.UpstreamIp + `",
							"port": 1982,
							"weight": 1
						}
					]
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("create route using the service just created", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/server_port",
				"service_id": "200"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		}),
	)
	It("batch test /server_port api", func() {
		// sleep for etcd sync
		time.Sleep(time.Duration(300) * time.Millisecond)

		// batch test /server_port api
		res := base.BatchTestServerPort(18, nil, "")

		Expect(res["1980"]).Should(Equal(6))
		Expect(res["1981"]).Should(Equal(6))
		Expect(res["1982"]).Should(Equal(6))
	})
	DescribeTable("delete route and service",
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
		Entry("remove service", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/services/200",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("make sure the route deleted", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/server_port",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		}),
	)
})

var _ = Describe("route create upstream", func() {
	DescribeTable("test route create upstream",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("create upstream", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/upstreams/1",
			Body: `{
				"nodes":[
					{
						"host": "` + base.UpstreamIp + `",
						"port": 1980,
						"weight": 1
					},
					{
						"host": "` + base.UpstreamIp + `",
						"port": 1981,
						"weight": 1
					},
					{
						"host": "` + base.UpstreamIp + `",
						"port": 1982,
						"weight": 1
					}
				],
				"type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("make sure the route has not created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/server_port",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		}),
		Entry("create route using the upstream just created", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/server_port",
				"upstream_id": "1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
	)
	It("batch test /server_port api", func() {
		// sleep for etcd sync
		time.Sleep(time.Duration(300) * time.Millisecond)

		// batch test /server_port api
		res := base.BatchTestServerPort(12, nil, "")

		Expect(res["1980"]).Should(Equal(4))
		Expect(res["1981"]).Should(Equal(4))
		Expect(res["1982"]).Should(Equal(4))
	})
	DescribeTable("delete route and upstream",
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
		Entry("remove upstream", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("make sure the route deleted", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/server_port",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		}),
	)
})

var _ = Describe("route create with service that contains upstream", func() {
	DescribeTable("test route create with service that contains upstream",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("make sure the route has not created", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/server_port",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		}),
		Entry("create upstream", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/upstreams/1",
			Body: `{
				"nodes":[
					{
						"host": "` + base.UpstreamIp + `",
						"port": 1980,
						"weight": 1
					},
					{
						"host": "` + base.UpstreamIp + `",
						"port": 1981,
						"weight": 1
					},
					{
						"host": "` + base.UpstreamIp + `",
						"port": 1982,
						"weight": 1
					}
				],
				"type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("create service", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/services/200",
			Body: `{
				"upstream_id": "1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("create route using the service just created", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/server_port",
				"service_id": "200"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		}),
	)
	It("batch test /server_port api", func() {
		// sleep for etcd sync
		time.Sleep(time.Duration(300) * time.Millisecond)

		// batch test /server_port api
		res := base.BatchTestServerPort(18, nil, "")

		Expect(res["1980"]).Should(Equal(6))
		Expect(res["1981"]).Should(Equal(6))
		Expect(res["1982"]).Should(Equal(6))
	})
	DescribeTable("delete route and service",
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
		Entry("remove service", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/services/200",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("make sure the route deleted", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/server_port",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		}),
	)
})
