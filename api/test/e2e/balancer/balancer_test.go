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
package balancer

import (
	"net/http"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	"github.com/stretchr/testify/assert"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = ginkgo.Describe("Balancer", func() {
	table.DescribeTable("test create upstream and route",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		table.Entry("create upstream (roundrobin with same weight)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/upstreams/1",
			Body: `{
				"nodes": [{
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
				}],
				"type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		table.Entry("create route using the upstream just created", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/1",
			Body: `{
				"name": "route1",
				"uri": "/server_port",
				"upstream_id": "1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		}),
	)
	ginkgo.It("verify balancer by access count(same weight)", func() {
		time.Sleep(base.SleepTime)
		// batch test /server_port api
		res := base.BatchTestServerPort(18, nil, "")
		assert.Equal(ginkgo.GinkgoT(), 6, res["1980"])
		assert.Equal(ginkgo.GinkgoT(), 6, res["1981"])
		assert.Equal(ginkgo.GinkgoT(), 6, res["1982"])
	})

	table.DescribeTable("test update upstream",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		table.Entry("update upstream (roundrobin with different weight)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/upstreams/1",
			Body: `{
				"nodes": [{
					"host": "` + base.UpstreamIp + `",
					"port": 1980,
					"weight": 1
				},
				{
					"host": "` + base.UpstreamIp + `",
					"port": 1981,
					"weight": 2
				},
				{
					"host": "` + base.UpstreamIp + `",
					"port": 1982,
					"weight": 3
				}],
				"type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
	)
	ginkgo.It("verify balancer by access count(different weight)", func() {
		time.Sleep(base.SleepTime)
		// batch test /server_port api
		res := base.BatchTestServerPort(18, nil, "")
		assert.Equal(ginkgo.GinkgoT(), 3, res["1980"])
		assert.Equal(ginkgo.GinkgoT(), 6, res["1981"])
		assert.Equal(ginkgo.GinkgoT(), 9, res["1982"])
	})

	table.DescribeTable("update upstream",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		table.Entry("update upstream (roundrobin with weight 1 and 0)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/upstreams/1",
			Body: `{
				"nodes": [{
					"host": "` + base.UpstreamIp + `",
					"port": 1980,
					"weight": 1
				},
				{
					"host": "` + base.UpstreamIp + `",
					"port": 1981,
					"weight": 0
				}],
				"type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
	)
	ginkgo.It("verify balancer by access count(weight 1 and 0)", func() {
		time.Sleep(base.SleepTime)
		// batch test /server_port api
		res := base.BatchTestServerPort(18, nil, "")
		assert.Equal(ginkgo.GinkgoT(), 18, res["1980"])
	})

	table.DescribeTable("update upstream",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		table.Entry("update upstream (roundrobin with weight only 1)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/upstreams/1",
			Body: `{
				"nodes": [{
					"host": "` + base.UpstreamIp + `",
					"port": 1980,
					"weight": 1
				}],
				"type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
	)
	ginkgo.It("verify balancer by access count(weight only 1)", func() {
		time.Sleep(base.SleepTime)
		// batch test /server_port api
		res := base.BatchTestServerPort(18, nil, "")
		assert.Equal(ginkgo.GinkgoT(), 18, res["1980"])
	})

	ginkgo.Context("test balancer delete", func() {
		ginkgo.It("delete route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		ginkgo.It("delete upstream", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/upstreams/1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		ginkgo.It("hit the route just deleted", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/server_port",
				ExpectStatus: http.StatusNotFound,
				ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
				Sleep:        base.SleepTime,
			})
		})
	})

})
