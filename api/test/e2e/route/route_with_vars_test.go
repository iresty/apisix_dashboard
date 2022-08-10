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
package route_test

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var upstream map[string]interface{} = map[string]interface{}{
	"type": "roundrobin",
	"nodes": []map[string]interface{}{
		{
			"host":   base.UpstreamIp,
			"port":   1980,
			"weight": 1,
		},
	},
}

var _ = Describe("test route with vars (args)", func() {
	It("add route with vars (args)", func() {
		var createRouteBody map[string]interface{} = map[string]interface{}{
			"name": "route1",
			"uri":  "/hello",
			"vars": [][]string{
				{"arg_name", "==", "aaa"},
			},
			"upstream": upstream,
		}
		_createRouteBody, err := json.Marshal(createRouteBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/routes/r1",
			Body:         string(_createRouteBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("hit the route with right args", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         `/hello`,
			Query:        "name=aaa",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		})
	})
	It("hit the route with wrong args", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         `/hello`,
			Query:        "name=bbb",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		})
	})
	It("hit the route with no args", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         `/hello`,
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		})
	})
	It("update route with vars (header)", func() {
		var createRouteBody map[string]interface{} = map[string]interface{}{
			"name": "route1",
			"uri":  "/hello",
			"vars": [][]string{
				{"http_k", "==", "header"},
			},
			"upstream": upstream,
		}
		_createRouteBody, err := json.Marshal(createRouteBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/routes/r1",
			Body:         string(_createRouteBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("hit the route with right header", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"k": "header"},
			Path:         `/hello`,
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		})
	})
	It("hit the route with wrong header", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"k": "jack"},
			Path:         `/hello`,
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		})
	})
	It("hit the route with no header", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         `/hello`,
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		})
	})
	It("update route with vars (cookie)", func() {
		var createRouteBody map[string]interface{} = map[string]interface{}{
			"name": "route1",
			"uri":  "/hello",
			"vars": [][]string{
				{"http_cookie", "==", "_octo=GH1.1.572248189.1598928545; _device_id=2c1a1a52074e66a3a008e4b73c690500; logged_in=yes;"},
			},
			"upstream": upstream,
		}
		_createRouteBody, err := json.Marshal(createRouteBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/routes/r1",
			Body:         string(_createRouteBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("hit the route with right Cookie", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Cookie": "_octo=GH1.1.572248189.1598928545; _device_id=2c1a1a52074e66a3a008e4b73c690500; logged_in=yes;"},
			Path:         `/hello`,
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		})
	})
	It("hit the route with wrong Cookie", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Cookie": "jack"},
			Path:         `/hello`,
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		})
	})
	It("hit the route with no Cookie", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         `/hello`,
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		})
	})
	It("delete route", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	It("hit the route just delete", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Cookie": "_octo=GH1.1.572248189.1598928545; _device_id=2c1a1a52074e66a3a008e4b73c690500; logged_in=yes;"},
			Path:         `/hello`,
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		})
	})
})

var _ = Describe("test route with multiple vars (args, cookie and header)", func() {
	It("add route with multiple vars (args, cookie and header)", func() {
		var createRouteBody map[string]interface{} = map[string]interface{}{
			"name": "route1",
			"uri":  "/hello",
			"vars": [][]string{
				{"http_cookie", "==", "_octo=GH1.1.572248189.1598928545; _device_id=2c1a1a52074e66a3a008e4b73c690500; logged_in=yes;"},
				{"http_k", "==", "header"},
				{"arg_name", "==", "aaa"},
			},
			"upstream": upstream,
		}
		_createRouteBody, err := json.Marshal(createRouteBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/routes/r1",
			Body:         string(_createRouteBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("hit the route with right parameters", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"k": "header", "Cookie": "_octo=GH1.1.572248189.1598928545; _device_id=2c1a1a52074e66a3a008e4b73c690500; logged_in=yes;"},
			Path:         `/hello`,
			Query:        "name=aaa",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		})
	})
	It("hit the route with wrong arg", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"k": "header", "Cookie": "_octo=GH1.1.572248189.1598928545; _device_id=2c1a1a52074e66a3a008e4b73c690500; logged_in=yes;"},
			Path:         `/hello`,
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		})
	})
	It("hit the route with wrong arg", func() {
		base.RunTestCase(base.HttpTestCase{
			Desc:         "hit the route with wrong header",
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"k": "test", "Cookie": "_octo=GH1.1.572248189.1598928545; _device_id=2c1a1a52074e66a3a008e4b73c690500; logged_in=yes;"},
			Path:         `/hello`,
			Query:        "name=aaa",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		})
	})
	It("hit the route with wrong cookie", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"k": "header", "Cookie": "_octo=GH1.1.572248189.1598928545; logged_in=yes;"},
			Path:         `/hello`,
			Query:        "name=aaa",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		})
	})
	It("delete route", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	It("hit the route just delete", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"k": "header", "Cookie": "_octo=GH1.1.572248189.1598928545; _device_id=2c1a1a52074e66a3a008e4b73c690500; logged_in=yes;"},
			Path:         `/hello`,
			Query:        "name=aaa",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		})
	})
})

var _ = Describe("test route with vars (args is digital)", func() {
	It("add route with vars (args is digital)", func() {
		var createRouteBody map[string]interface{} = map[string]interface{}{
			"name": "route1",
			"uri":  "/hello",
			"vars": [][]string{
				{"arg_name", "==", "111"},
			},
			"upstream": upstream,
		}
		_createRouteBody, err := json.Marshal(createRouteBody)
		Expect(err).To(BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/routes/r1",
			Body:         string(_createRouteBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("verify route", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         `/hello`,
			Query:        "name=111",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		})
	})
	It("delete the route with vars (args is digital)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	It("hit route just delete", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         `/hello`,
			Query:        "name=111",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		})
	})
})
