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
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Upstream", func() {
	It("create route failed, using non-existent upstream_id", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				 "uri": "/hello",
				 "upstream_id": "not-exists"
			 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		})
	})
	It("create upstream success", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["name"] = "upstream1"
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":   base.UpstreamIp,
				"port":   1980,
				"weight": 1,
			},
		}
		createUpstreamBody["type"] = "roundrobin"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/1",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("create upstream2 success", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["name"] = "upstream2"
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":   base.UpstreamIp,
				"port":   1980,
				"weight": 1,
			},
		}
		createUpstreamBody["type"] = "roundrobin"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/2",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("create upstream3 success when pass host is 'node' and nodes without port", func() {
		By("create upstream3", func() {
			createUpstreamBody := make(map[string]interface{})
			createUpstreamBody["name"] = "upstream3"
			createUpstreamBody["nodes"] = map[string]float64{base.UpstreamIp: 100}
			createUpstreamBody["type"] = "roundrobin"
			createUpstreamBody["pass_host"] = "node"

			_createUpstreamBody, err := json.Marshal(createUpstreamBody)
			Expect(err).To(BeNil())
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodPut,
				Path:         "/apisix/admin/upstreams/3",
				Body:         string(_createUpstreamBody),
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})

		By("create route using the upstream3", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/1",
				Body: `{
				 "name": "route1",
				 "uri": "/hello",
				 "upstream_id": "3"
			 }`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			})
		})

		By("hit the route just created", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.APISIXExpect(),
				Method:       http.MethodGet,
				Path:         "/hello",
				ExpectStatus: http.StatusOK,
				ExpectBody:   "hello",
				Sleep:        base.SleepTime,
			})
		})

		By("delete route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/1",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
	})
	It("create upstream failed, name existed", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["name"] = "upstream2"
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":   base.UpstreamIp,
				"port":   1980,
				"weight": 1,
			},
		}
		createUpstreamBody["type"] = "roundrobin"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPost,
			Path:         "/apisix/admin/upstreams",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   `upstream name exists`,
			Sleep:        base.SleepTime,
		})
	})
	It("update upstream failed, name existed", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["name"] = "upstream1"
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":   base.UpstreamIp,
				"port":   1980,
				"weight": 1,
			},
		}
		createUpstreamBody["type"] = "roundrobin"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/2",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   `upstream name exists`,
		})
	})
	It("update upstream success", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["name"] = "upstream22"
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":   base.UpstreamIp,
				"port":   1980,
				"weight": 1,
			},
		}
		createUpstreamBody["type"] = "roundrobin"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/2",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("check upstream exists by name", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/notexist/upstreams",
			Query:        "name=upstream1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   "Upstream name is reduplicate",
			Sleep:        base.SleepTime,
		})
	})
	It("upstream name list", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/names/upstreams",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"name":"upstream1"`,
			Sleep:        base.SleepTime,
		})
	})
	It("check upstream exists by name (exclude it self)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/notexist/upstreams",
			Query:        "name=upstream1&exclude=1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	It("create route using the upstream just created", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/1",
			Body: `{
				 "name": "route1",
				 "uri": "/hello",
				 "upstream_id": "1"
			 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	It("hit the route just created", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		})
	})
	It("delete not exist upstream", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/not-exist",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		})
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
	It("delete upstream2", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/2",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("delete upstream3", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/3",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("hit the route just deleted", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		})
	})
})

var _ = Describe("Upstream update with domain", func() {
	It("create upstream success", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["name"] = "upstream1"
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":     base.UpstreamIp,
				"port":     1980,
				"weight":   1,
				"priority": 10,
			},
		}
		createUpstreamBody["type"] = "roundrobin"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/1",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"code":0`,
		})
	})
	It("create route using the upstream(use proxy rewriteproxy rewrite plugin)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/1",
			Body: `{
				"name": "route1",
				 "uri": "/*",
				 "upstream_id": "1",
				 "plugins": {
					"proxy-rewrite": {
						"uri": "/",
						"scheme": "https"
					}
				}
			 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	It("update upstream with domain", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":     "www.google.com",
				"port":     443,
				"weight":   1,
				"priority": 10,
			},
		}
		createUpstreamBody["type"] = "roundrobin"
		createUpstreamBody["pass_host"] = "node"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/1",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("hit the route using upstream 1", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "google",
			Sleep:        base.SleepTime,
		})
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
	It("hit the route just deleted", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/get",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		})
	})
})

var _ = Describe("Upstream chash remote addr", func() {
	It("create chash upstream with key (remote_addr)", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":   base.UpstreamIp,
				"port":   1980,
				"weight": 1,
			},
			{
				"host":   base.UpstreamIp,
				"port":   1981,
				"weight": 1,
			},
			{
				"host":   base.UpstreamIp,
				"port":   1982,
				"weight": 1,
			},
		}
		createUpstreamBody["type"] = "chash"
		createUpstreamBody["hash_on"] = "header"
		createUpstreamBody["key"] = "remote_addr"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/1",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})

	It("create route using the upstream just created", func() {
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

	It("hit routes(upstream weight 1)", func() {
		time.Sleep(time.Duration(500) * time.Millisecond)
		basepath := base.APISIXHost
		request, err := http.NewRequest("GET", basepath+"/server_port", nil)
		Expect(err).To(BeNil())
		request.Header.Add("Authorization", base.GetToken())
		res := map[string]int{}
		for i := 0; i < 18; i++ {
			resp, err := http.DefaultClient.Do(request)
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			respBody, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			body := string(respBody)
			if _, ok := res[body]; !ok {
				res[body] = 1
			} else {
				res[body]++
			}
		}
		Expect(res["1982"]).Should(Equal(18))
	})

	It("create chash upstream with key (remote_addr, weight equal 0 or 1)", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":   base.UpstreamIp,
				"port":   1980,
				"weight": 1,
			},
			{
				"host":   base.UpstreamIp,
				"port":   1981,
				"weight": 0,
			},
			{
				"host":   base.UpstreamIp,
				"port":   1982,
				"weight": 0,
			},
		}
		createUpstreamBody["type"] = "chash"
		createUpstreamBody["hash_on"] = "header"
		createUpstreamBody["key"] = "remote_addr"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/1",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("create route using the upstream just created", func() {
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
	It("hit routes(remote_addr, weight equal 0 or 1)", func() {
		time.Sleep(time.Duration(500) * time.Millisecond)
		basepath := base.APISIXHost
		request, err := http.NewRequest("GET", basepath+"/server_port", nil)
		Expect(err).To(BeNil())
		request.Header.Add("Authorization", base.GetToken())
		count := 0
		for i := 0; i <= 17; i++ {
			resp, err := http.DefaultClient.Do(request)
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			respBody, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			if string(respBody) == "1980" {
				count++
			}
		}
		Expect(count).Should(Equal(18))
	})
	It("create chash upstream with key (remote_addr, all weight equal 0)", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":   base.UpstreamIp,
				"port":   1980,
				"weight": 0,
			},
			{
				"host":   base.UpstreamIp,
				"port":   1982,
				"weight": 0,
			},
		}
		createUpstreamBody["type"] = "chash"
		createUpstreamBody["hash_on"] = "header"
		createUpstreamBody["key"] = "remote_addr"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())

		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/1",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("create route using the upstream just created", func() {
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
	It("hit the route(remote_addr, all weight equal 0)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/server_port",
			ExpectStatus: http.StatusBadGateway,
			ExpectBody:   "<head><title>502 Bad Gateway</title></head>",
			Sleep:        base.SleepTime,
		})
	})
	It("create chash upstream u2", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/upstreams/u2",
			Body: `{
				"retries": 1,
				"timeout": {
					"connect":15,
					"send":15,
					"read":15
				},
				"type":"roundrobin",
				"scheme": "https",
				"service_name": "USER-SERVICE",
				"discovery_type": "eureka",
				"name": "upstream-for-test",
				"desc": "hello world"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("get the upstream to verify config", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/upstreams/u2",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"type":"roundrobin","scheme":"https","discovery_type":"eureka"`,
			Sleep:        base.SleepTime,
		})
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
	It("delete upstream u2", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/u2",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("hit the route just deleted", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/server_port",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		})
	})
})

var _ = Describe("Upstream create via post", func() {
	It("create upstream via POST", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["id"] = "u1"
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":   base.UpstreamIp,
				"port":   1981,
				"weight": 1,
			},
		}
		createUpstreamBody["type"] = "roundrobin"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPost,
			Path:         "/apisix/admin/upstreams",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			// should return id and other request body content
			ExpectBody: []string{"\"id\":\"u1\"", "\"type\":\"roundrobin\""},
		})
	})
	It("create route using the upstream just created", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/1",
			Body: `{
				 "name": "route1",
				 "uri": "/hello",
				 "upstream_id": "u1"
			 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	It("hit the route just created", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		})
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
			Path:         "/apisix/admin/upstreams/u1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("hit the route just deleted", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		})
	})
})

var _ = Describe("Upstream update use patch method", func() {
	It("create upstream via POST", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["id"] = "u1"
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":   base.UpstreamIp,
				"port":   1981,
				"weight": 1,
			},
		}
		createUpstreamBody["type"] = "roundrobin"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPost,
			Path:         "/apisix/admin/upstreams",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			// should return id and other request body content
			ExpectBody: []string{"\"id\":\"u1\"", "\"type\":\"roundrobin\""},
		})
	})
	It("update upstream use patch method", func() {
		createUpstreamBody := make(map[string]interface{})
		createUpstreamBody["nodes"] = []map[string]interface{}{
			{
				"host":     base.UpstreamIp,
				"port":     1981,
				"weight":   1,
				"priority": 10,
			},
		}
		createUpstreamBody["type"] = "roundrobin"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPatch,
			Path:         "/apisix/admin/upstreams/u1",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("get upstream data", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/upstreams/u1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "nodes\":[{\"host\":\"" + base.UpstreamIp + "\",\"port\":1981,\"weight\":1,\"priority\":10}],\"type\":\"roundrobin\"}",
		})
	})
	It("Upstream update use patch method", func() {
		var nodes []map[string]interface{} = []map[string]interface{}{
			{
				"host":     base.UpstreamIp,
				"port":     1980,
				"weight":   1,
				"priority": 10,
			},
		}
		_nodes, err := json.Marshal(nodes)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPatch,
			Path:   "/apisix/admin/upstreams/u1/nodes",
			Body:   string(_nodes),
			Headers: map[string]string{
				"Authorization": base.GetToken(),
				"Content-Type":  "text/plain",
			},
			ExpectStatus: http.StatusOK,
		})
	})
	It("get upstream data", func() {

		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/upstreams/u1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "nodes\":[{\"host\":\"" + base.UpstreamIp + "\",\"port\":1980,\"weight\":1,\"priority\":10}],\"type\":\"roundrobin\"}",
		})
	})
	It("delete upstream", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/u1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
})

var _ = Describe("test upstream delete (route is in use)", func() {
	DescribeTable("test upstream delete",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("create upstream without plugin", base.HttpTestCase{
			Desc:    "create upstream without plugin",
			Object:  base.ManagerApiExpect(),
			Method:  http.MethodPut,
			Path:    "/apisix/admin/upstreams/u1",
			Headers: map[string]string{"Authorization": base.GetToken()},
			Body: `{
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}`,
			ExpectStatus: http.StatusOK,
			ExpectBody:   "\"id\":\"u1\"",
		}),
		Entry("create route use upstream r1", base.HttpTestCase{
			Desc:   "create route use upstream u1",
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"id": "r1",
				"uri": "/hello",
				"upstream_id": "u1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "\"upstream_id\":\"u1\"",
		}),
		Entry("delete upstream failed", base.HttpTestCase{
			Desc:         "delete upstream failed",
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/u1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   "route: route1 is using this upstream",
		}),
		Entry("delete route first", base.HttpTestCase{
			Desc:         "delete route first",
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("check route exist", base.HttpTestCase{
			Desc:         "check route exist",
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		}),
		Entry("delete upstream success", base.HttpTestCase{
			Desc:         "delete upstream success",
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/u1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("check upstream exist", base.HttpTestCase{
			Desc:         "check upstream exist",
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/upstreams/u1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		}))
})

var _ = Describe("test upstream delete (service is in use)", func() {
	DescribeTable("test upstream delete",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("create upstream without plugin", base.HttpTestCase{
			Desc:    "create upstream without plugin",
			Object:  base.ManagerApiExpect(),
			Method:  http.MethodPut,
			Path:    "/apisix/admin/upstreams/u1",
			Headers: map[string]string{"Authorization": base.GetToken()},
			Body: `{
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}`,
			ExpectStatus: http.StatusOK,
			ExpectBody:   "\"id\":\"u1\"",
		}),
		Entry("create service use upstream r1", base.HttpTestCase{
			Desc:   "create service use upstream r1",
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/services/s1",
			Body: `{
				"id": "s1",
				"name": "service1",
				"upstream_id": "u1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "\"upstream_id\":\"u1\"",
		}),
		Entry("delete upstream failed", base.HttpTestCase{
			Desc:         "delete upstream failed",
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/u1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   "service: service1 is using this upstream",
		}),
		Entry("delete service first", base.HttpTestCase{
			Desc:         "delete service first",
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/services/s1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("check service exist", base.HttpTestCase{
			Desc:         "check service exist",
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/services/s1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		}),
		Entry("delete upstream success", base.HttpTestCase{
			Desc:         "delete upstream success",
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/u1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("check upstream exist", base.HttpTestCase{
			Desc:         "check upstream exist",
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/upstreams/u1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		}))
})
