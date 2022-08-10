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
package server_info_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Server Info Test", func() {
	DescribeTable("Get server info",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("Get server info (apisix-server1)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Path:         "/apisix/admin/server_info/apisix-server1",
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"hostname":"apisix_server1"`,
		}),
		Entry("Get server info (apisix-server2)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Path:         "/apisix/admin/server_info/apisix-server2",
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"hostname":"apisix_server2"`,
		}),
	)

	DescribeTable("List server info",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("List all server info", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Path:         "/apisix/admin/server_info",
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":2`,
		}),
		Entry("List server info (With condition, hostname = apisix_)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Path:         "/apisix/admin/server_info",
			Query:        "hostname=apisix_",
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":2`,
		}),
		Entry("List server info (With condition, hostname = apisix_server2)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Path:         "/apisix/admin/server_info",
			Query:        "hostname=apisix_server2",
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":1`,
		}),
	)

	DescribeTable("Get server info (Omit Empty Value)",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("get server info", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Path:         "/apisix/admin/server_info/apisix-server1",
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			UnexpectBody: []string{`"create_time":`, `"update_time":`},
		}),
	)

	DescribeTable("List server info (Omit Empty Value)",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("List all server info", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Path:         "/apisix/admin/server_info",
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":2`,
			UnexpectBody: []string{`"create_time":`, `"update_time":`},
		}),
		Entry("List server info (With condition, hostname = apisix_)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Path:         "/apisix/admin/server_info",
			Query:        "hostname=apisix_",
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":2`,
			UnexpectBody: []string{`"create_time":`, `"update_time":`},
		}),
	)
})
