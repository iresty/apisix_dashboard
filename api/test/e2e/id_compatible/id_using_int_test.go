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
package id_compatible_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = DescribeTable("Id Using Int",
	func(tc base.HttpTestCase) {
		base.RunTestCase(tc)
	},
	Entry("create upstream", base.HttpTestCase{
		Object: base.ManagerApiExpect(),
		Method: http.MethodPut,
		Path:   "/apisix/admin/upstreams",
		Body: `{
				"id": 1,
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
	Entry("create route using the upstream just created", base.HttpTestCase{
		Object: base.ManagerApiExpect(),
		Method: http.MethodPut,
		Path:   "/apisix/admin/routes/1",
		Body: `{
				"name": "route1",
				"uri": "/hello",
				"upstream_id": 1
			}`,
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
		Sleep:        base.SleepTime,
	}),
	Entry("hit the route just created", base.HttpTestCase{
		Object:       base.APISIXExpect(),
		Method:       http.MethodGet,
		Path:         "/hello",
		ExpectStatus: http.StatusOK,
		ExpectBody:   "hello world",
		Sleep:        base.SleepTime,
	}),
	Entry("create service", base.HttpTestCase{
		Object: base.ManagerApiExpect(),
		Method: http.MethodPut,
		Path:   "/apisix/admin/services",
		Body: `{
				"id": 1,
				"upstream_id": 1
			}`,
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
	}),
	Entry("update route to use the service just created", base.HttpTestCase{
		Object: base.ManagerApiExpect(),
		Method: http.MethodPut,
		Path:   "/apisix/admin/routes/1",
		Body: `{
				"name": "route1",
				"uri": "/hello",
				"service_id": 1
			}`,
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
		Sleep:        base.SleepTime,
	}),
	Entry("hit the route just updated", base.HttpTestCase{
		Object:       base.APISIXExpect(),
		Method:       http.MethodGet,
		Path:         "/hello",
		ExpectStatus: http.StatusOK,
		ExpectBody:   "hello world",
		Sleep:        base.SleepTime,
	}),
	Entry("delete the route", base.HttpTestCase{
		Object:       base.ManagerApiExpect(),
		Method:       http.MethodDelete,
		Path:         "/apisix/admin/routes/1",
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
	}),
	Entry("delete the service", base.HttpTestCase{
		Object:       base.ManagerApiExpect(),
		Method:       http.MethodDelete,
		Path:         "/apisix/admin/services/1",
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
		Sleep:        base.SleepTime,
	}),
	Entry("make sure the service has been deleted", base.HttpTestCase{
		Object:       base.ManagerApiExpect(),
		Method:       http.MethodGet,
		Path:         "/apisix/admin/services/1",
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusNotFound,
	}),
	Entry("delete the upstream", base.HttpTestCase{
		Object:       base.ManagerApiExpect(),
		Method:       http.MethodDelete,
		Path:         "/apisix/admin/upstreams/1",
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
	}),
	Entry("make sure the upstream has been deleted", base.HttpTestCase{
		Object:       base.ManagerApiExpect(),
		Method:       http.MethodGet,
		Path:         "/apisix/admin/upstreams/1",
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusNotFound,
		Sleep:        base.SleepTime,
	}),
	Entry("hit deleted route", base.HttpTestCase{
		Object:       base.APISIXExpect(),
		Method:       http.MethodGet,
		Path:         "/hello",
		ExpectStatus: http.StatusNotFound,
		Sleep:        base.SleepTime,
	}),
)
