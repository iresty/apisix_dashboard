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

var _ = DescribeTable("Consumer With key-auth",
	func(tc base.HttpTestCase) {
		base.RunTestCase(tc)
	},
	Entry("create route", base.HttpTestCase{
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
	Entry("hit route without apikey", base.HttpTestCase{
		Object:       base.APISIXExpect(),
		Method:       http.MethodGet,
		Path:         "/hello",
		ExpectStatus: http.StatusUnauthorized,
		ExpectBody:   "Missing API key found in request",
		Sleep:        base.SleepTime * 2,
	}),
	Entry("create consumer", base.HttpTestCase{
		Object: base.ManagerApiExpect(),
		Method: http.MethodPut,
		Path:   "/apisix/admin/consumers",
		Body: `{
			"username": "jack",
			"plugins": {
				"key-auth": {
					"key": "auth-one"
				}
			},
			"desc": "test description"
		}`,
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
	}),
	Entry("hit route with correct apikey", base.HttpTestCase{
		Object:       base.APISIXExpect(),
		Method:       http.MethodGet,
		Path:         "/hello",
		Headers:      map[string]string{"apikey": "auth-one"},
		ExpectStatus: http.StatusOK,
		ExpectBody:   "hello world",
		Sleep:        base.SleepTime * 2,
	}),
	Entry("hit route with incorrect apikey", base.HttpTestCase{
		Object:       base.APISIXExpect(),
		Method:       http.MethodGet,
		Path:         "/hello",
		Headers:      map[string]string{"apikey": "auth-new"},
		ExpectStatus: http.StatusUnauthorized,
		ExpectBody:   "Invalid API key in request",
		Sleep:        base.SleepTime * 2,
	}),
	Entry("delete consumer", base.HttpTestCase{
		Object:       base.ManagerApiExpect(),
		Method:       http.MethodDelete,
		Path:         "/apisix/admin/consumers/jack",
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
	}),
	Entry("hit route (consumer deleted)", base.HttpTestCase{
		Object:       base.APISIXExpect(),
		Method:       http.MethodGet,
		Path:         "/hello",
		Headers:      map[string]string{"apikey": "auth-one"},
		ExpectStatus: http.StatusUnauthorized,
		ExpectBody:   "Missing related consumer",
		Sleep:        base.SleepTime * 2,
	}),
	Entry("delete route", base.HttpTestCase{
		Object:       base.ManagerApiExpect(),
		Method:       http.MethodDelete,
		Path:         "/apisix/admin/routes/r1",
		Headers:      map[string]string{"Authorization": base.GetToken()},
		ExpectStatus: http.StatusOK,
	}),
)
