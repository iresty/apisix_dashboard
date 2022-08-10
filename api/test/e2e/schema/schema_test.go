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
			ExpectBody:   "{\"dependencies\":{\"algorithm\":{\"oneOf\":[{\"properties\":{\"algorithm\":{\"default\":\"HS256\",\"enum\":[\"HS256\",\"HS512\"]}}},{\"properties\":{\"algorithm\":{\"enum\":[\"RS256\"]},\"private_key\":{\"type\":\"string\"},\"public_key\":{\"type\":\"string\"}},\"required\":[\"private_key\",\"public_key\"]},{\"properties\":{\"algorithm\":{\"enum\":[\"RS256\"]},\"vault\":{\"properties\":{},\"type\":\"object\"}},\"required\":[\"vault\"]}]}},\"properties\":{\"algorithm\":{\"default\":\"HS256\",\"enum\":[\"HS256\",\"HS512\",\"RS256\"],\"type\":\"string\"},\"base64_secret\":{\"default\":false,\"type\":\"boolean\"},\"exp\":{\"default\":86400,\"minimum\":1,\"type\":\"integer\"},\"key\":{\"type\":\"string\"},\"secret\":{\"type\":\"string\"},\"vault\":{\"properties\":{},\"type\":\"object\"}},\"required\":[\"key\"],\"type\":\"object\"}",
			Sleep:        base.SleepTime,
		}),
		Entry("get schema of plugin `require-id`", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/schema/plugins/jwt-auth",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"$comment":"this is a mark for our injected plugin schema","properties":{"cookie":{"default":"jwt","type":"string"},"disable":{"type":"boolean"},"header":{"default":"authorization","type":"string"},"query":{"default":"jwt","type":"string"}}`,
			Sleep:        base.SleepTime,
		}),
		Entry("get schema of non-existent plugin", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/schema/plugins/non-existent",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `schema of plugins non-existent not found`,
			Sleep:        base.SleepTime,
		}),
		Entry("get schema of consumer", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/schemas/consumer",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"properties":{"create_time":{"type":"integer"},"desc":{"maxLength":256,"type":"string"},"labels":{"description":"key/value pairs to specify attributes","patternProperties":{".*":{"description":"value of label","maxLength":64,"minLength":1,"pattern":"^\\S+$","type":"string"}},"type":"object"},"plugins":{"type":"object"},"update_time":{"type":"integer"},"username":{"maxLength":100,"minLength":1,"pattern":"^[a-zA-Z0-9_]+$","type":"string"}}`,
			Sleep:        base.SleepTime,
		}),
		Entry("get schema of non-existent resources", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/schemas/non-existent",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `schema of non-existent not found`,
			Sleep:        base.SleepTime,
		}),
	)
})
