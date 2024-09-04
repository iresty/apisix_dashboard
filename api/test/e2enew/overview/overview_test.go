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
package overview_test

import (
	"net/http"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"

	"github.com/apisix/manager-api/test/e2enew/base"
)

var _ = ginkgo.Describe("Overview", func() {
	table.DescribeTable("test dashboard overview",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		table.Entry("create route1 success(online)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				 "name": "route1",
				 "uri": "/hello",
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
		table.Entry("create route2 success(offline)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r2",
			Body: `{
				 "name": "route2",
				 "uri": "/hello1",
				 "status": 0,
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
		table.Entry("create upstream", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/upstreams",
			Body: `{
				 "id": "u1",
				 "name":"upstream1",
				 "nodes":[
					 {
						 "host":"` + base.UpstreamIp + `",
						 "port":1980,
						 "weight":1
					 },
					 {
						 "host":"` + base.UpstreamIp + `",
						 "port":1981,
						 "weight":1
					 }
				 ],
				 "type":"roundrobin"
			 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"name":"upstream1"`,
		}),
		table.Entry("create service success", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/services/s1",
			Body: `{
				 "name": "testservice",
				 "upstream": {
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
				 }
			 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		table.Entry("get dashboard overview", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Path:         "/apisix/admin/overview",
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
			ExpectBody:   []string{`"router_cnt":2`, `"service_cnt":1`, `"upstream_cnt":1`, `gateway_info`, `plugins`},
		}),
		table.Entry("delete service", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Path:         "/apisix/admin/services/s1",
			ExpectStatus: http.StatusOK,
		}),
		table.Entry("delete upstream", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Path:         "/apisix/admin/upstreams/u1",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"code":0`,
		}),
		table.Entry("delete routes", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Path:         "/apisix/admin/routes/r1,r2",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"code":0`,
		}),
	)
})
