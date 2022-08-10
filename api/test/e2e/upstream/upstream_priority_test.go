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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

// just test for schema check
var _ = Describe("Upstream priority", func() {
	It("create upstream with priority", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":     base.UpstreamIp,
				"port":     1980,
				"weight":   1,
				"priority": 10,
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
			Path:         "/apisix/admin/upstreams/priority",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("delete upstream", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/priority",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
})

//test node priority
var _ = Describe("Upstream priority", func() {
	It("create upstream with priority", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/upstreams/1",
			Body: `{
				"nodes":[
					{
						"host": "` + base.UpstreamIp + `",
						"port": 1980,
						"weight": 1,
						"priority": 1
					},
					{
						"host": "` + base.UpstreamIp + `",
						"port": 1981,
						"weight": 1,
						"priority": 2
					}
				],
				"type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("create route using the upstream", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/1",
			Body: `{
				 "name": "route1",
				  "uri": "/server_port",
				  "upstream_id": "1"
			  }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	It("batch test /server_port api", func() {
		// sleep for etcd sync
		time.Sleep(time.Duration(300) * time.Millisecond)

		// batch test /server_port api
		res := base.BatchTestServerPort(12, nil, "")

		Expect(res["1980"]).Should(Equal(0))
		Expect(res["1981"]).Should(Equal(12))
	})
	It("update upstream with priority", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/upstreams/1",
			Body: `{
				"nodes":[
					{
						"host": "` + base.UpstreamIp + `",
						"port": 1980,
						"weight": 1,
						"priority": 3
					},
					{
						"host": "` + base.UpstreamIp + `",
						"port": 1981,
						"weight": 1,
						"priority": 2
					}
				],
				"type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("batch test /server_port api", func() {
		// sleep for etcd sync
		time.Sleep(time.Duration(300) * time.Millisecond)

		// batch test /server_port api
		res := base.BatchTestServerPort(12, nil, "")

		Expect(res["1980"]).Should(Equal(12))
		Expect(res["1981"]).Should(Equal(0))
	})
	It("delete route", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("delete upstream", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
})
