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
package schema

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Schema Test", func() {
	DescribeTable("test schema basic", func(testCase base.HttpTestCase) {
		base.RunTestCase(testCase)
	},
		Entry("get consumer schema of plugin", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/schema/plugins/jwt-auth",
			Query:        "schema_type=consumer",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "\"enum\":[\"HS256\",\"HS512\"]}}},{\"properties\":{\"algorithm\":{\"enum\":[\"RS256\",\"ES256\"]}",
			Sleep:        base.SleepTime,
		}),
		Entry("get schema of plugin `require-id`", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/schema/plugins/jwt-auth",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "\"$comment\":\"this is a mark for our injected plugin schema\",\"properties\":{\"_meta\":{\"properties\"",
			Sleep:        base.SleepTime,
		}),
		Entry("get schema of non-existent plugin", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/schema/plugins/no-exist",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   `failed to load plugin no-exist`,
			Sleep:        base.SleepTime,
		}),
		Entry("get schema of consumer", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/schema/route",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "\"required\":[\"uri\"]}",
			Sleep:        base.SleepTime,
		}),
		Entry("get schema of non-existent resources", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/schema/no-exist",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   `not found schema: no-exist`,
			Sleep:        base.SleepTime,
		}),
	)
})
