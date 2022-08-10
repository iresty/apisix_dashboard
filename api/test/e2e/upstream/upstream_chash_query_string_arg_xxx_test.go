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
	"sort"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var createUpstreamBody map[string]interface{} = map[string]interface{}{
	"nodes": []map[string]interface{}{
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
	},
	"type": "chash",
}

var _ = Describe("Upstream chash query string", func() {
	It("create chash upstream with key (query_string)", func() {
		createUpstreamBody["key"] = "query_string"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/1",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"id\":\"1\"", "\"key\":\"query_string\""},
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
	It("hit routes(upstream query_string)", func() {
		time.Sleep(time.Duration(500) * time.Millisecond)
		basepath := base.APISIXHost
		res := map[string]int{}
		for i := 0; i < 180; i++ {
			url := basepath + "/server_port?var=2&var2=" + strconv.Itoa(i)
			req, err := http.NewRequest("GET", url, nil)
			Expect(err).To(BeNil())
			resp, err := http.DefaultClient.Do(req)
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
		var counts []int
		for _, value := range res {
			counts = append(counts, value)
		}
		sort.Ints(counts)
		Expect(float64(counts[2]-counts[0]) / float64(counts[1])).Should(BeNumerically("<", 0.4))
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
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		})
	})
})

var _ = Describe("Upstream chash query string", func() {
	It("create chash upstream with key (arg_xxx)", func() {
		createUpstreamBody["key"] = "arg_device_id"
		_createUpstreamBody, err := json.Marshal(createUpstreamBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/upstreams/1",
			Body:         string(_createUpstreamBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"id\":\"1\"", "\"key\":\"arg_device_id\""},
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
	It("hit routes(upstream arg_device_id)", func() {
		time.Sleep(time.Duration(500) * time.Millisecond)
		basepath := base.APISIXHost
		res := map[string]int{}
		for i := 0; i <= 17; i++ {
			url := basepath + "/server_port?device_id=" + strconv.Itoa(i)
			req, err := http.NewRequest("GET", url, nil)
			Expect(err).To(BeNil())
			resp, err := http.DefaultClient.Do(req)
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
		var counts []int
		for _, value := range res {
			counts = append(counts, value)
		}
		sort.Ints(counts)
		Expect(float64(counts[2]-counts[0]) / float64(counts[1])).Should(BeNumerically("<", 0.4))
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
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		})
	})
})
