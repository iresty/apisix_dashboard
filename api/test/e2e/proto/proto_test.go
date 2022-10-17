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
package proto

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Proto", func() {
	DescribeTable("test proto create and update",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("create proto1 by id", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/protos/1",
			Body: `{
				"content": "syntax = \"proto3\";
				package proto;
				message HelloRequest{
					string name = 1;
				}

				message HelloResponse{
					int32 code = 1;
					string msg = 2;
				}
				// The greeting service definition.
				service Hello {
					// Sends a greeting
					rpc SayHi (HelloRequest) returns (HelloResponse){}
				}"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get proto1", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/protos/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"key\":\"/apisix/protos/1\"", "string name = 1;"},
		}),
		Entry("delete proto1", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/protos/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get protos", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/protos",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			UnexpectBody: []string{"\"key\":\"/apisix/protos/1\""},
		}),
	)
})
