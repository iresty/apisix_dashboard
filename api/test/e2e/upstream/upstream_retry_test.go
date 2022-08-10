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
package upstream_test

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

// just test for schema check
var _ = Describe("Upstream keepalive pool", func() {
	It("create upstream with keepalive pool", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":   base.UpstreamIp,
				"port":   1980,
				"weight": 1,
			},
		}
		createUpstreamBody["type"] = "roundrobin"
		createUpstreamBody["retries"] = 5
		createUpstreamBody["retry_timeout"] = 5.5
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/retry",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("delete upstream", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/retry",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("zero retry field", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":   base.UpstreamIp,
				"port":   1980,
				"weight": 1,
			}}
		createUpstreamBody["type"] = "roundrobin"
		createUpstreamBody["retries"] = 0
		createUpstreamBody["retry_timeout"] = 5.5
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/zero-retry",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/upstreams/zero-retry",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"retries":0`,
		})
	})
	It("nil retry field", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":   base.UpstreamIp,
				"port":   1980,
				"weight": 1,
			}}
		createUpstreamBody["type"] = "roundrobin"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/zero-retry",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/upstreams/zero-retry",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			UnexpectBody: `"retries"`,
		})
	})
	It("delete upstream", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/zero-retry",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})

})
