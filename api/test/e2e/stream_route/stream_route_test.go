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
package stream_route_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Stream Route", func() {
	DescribeTable("test stream route create and update",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("create stream route 1 by id", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/stream_routes/1",
			Body: `{
				"remote_addr": "127.0.0.1",
				"desc": "test-desc",
				"upstream": {
					"nodes": {
						"127.0.0.1:8080": 1
					},
					"type": "roundrobin"
				},
				"desc": "new route"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get stream route 1", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/stream_routes/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"key\":\"/apisix/stream_routes/1\"", "\"remote_addr\":\"127.0.0.1\""},
		}),
		Entry("create stream route by random id", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/stream_routes",
			Body: `{
				"remote_addr": "127.0.0.2",
				"desc": "test-desc",
				"upstream": {
					"nodes": {
						"127.0.0.1:8080": 1
					},
					"type": "roundrobin"
				},
				"desc": "new route"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get stream routes", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/stream_routes",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"key\":\"/apisix/stream_routes/1\"", "\"remote_addr\":\"127.0.0.2\""},
		}),
	)
})
