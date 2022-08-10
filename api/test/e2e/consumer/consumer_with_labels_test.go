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
package consumer_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = DescribeTable("Consumer With Labels",
	func(tc base.HttpTestCase) {
		base.RunTestCase(tc)
	},
	Entry("create the consumer", base.HttpTestCase{
		Object: base.ManagerApiExpect(),
		Method: http.MethodPut,
		Path:   "/apisix/admin/consumers",
		Body: `{
			"username": "jack",
			"labels": {
				"build":"16",
				"env":"production",
				"version":"v2"
			},
			"plugins": {
				"key-auth": {
					"key": "auth-two"
				}
			},
			"desc": "test description"
		}`,
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
	}),
	Entry("verify the consumer", base.HttpTestCase{
		Object:       base.ManagerApiExpect(),
		Method:       http.MethodGet,
		Path:         "/apisix/admin/consumers/jack",
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
		ExpectBody:   "\"username\":\"jack\",\"desc\":\"test description\",\"plugins\":{\"key-auth\":{\"key\":\"auth-two\"}},\"labels\":{\"build\":\"16\",\"env\":\"production\",\"version\":\"v2\"}",
	}),
	Entry("create the route", base.HttpTestCase{
		Object: base.ManagerApiExpect(),
		Method: http.MethodPut,
		Path:   "/apisix/admin/routes/r1",
		Body: `{
				"name": "route1",
				"uri": "/hello",
				"plugins": {
					"key-auth": {}
				},
				"upstream": {
					"type": "roundrobin",
					"nodes": [{
						"host": "` + base.UpstreamIp + `",
						"port": 1980,
						"weight": 1
					}]
				}
			}`,
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
	}),
	Entry("hit route", base.HttpTestCase{
		Object:       base.APISIXExpect(),
		Method:       http.MethodGet,
		Path:         "/hello",
		Headers:      map[string]string{"apikey": "auth-two"},
		ExpectStatus: http.StatusOK,
	}),
	Entry("delete consumer", base.HttpTestCase{
		Object:       base.ManagerApiExpect(),
		Method:       http.MethodDelete,
		Path:         "/apisix/admin/consumers/jack",
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
	}),
	Entry("delete route", base.HttpTestCase{
		Object:       base.ManagerApiExpect(),
		Method:       http.MethodDelete,
		Path:         "/apisix/admin/routes/r1",
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
	}),
)
