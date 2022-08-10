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
package route_online_debug

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/gomega"
	"github.com/tidwall/gjson"

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

var _ = ginkgo.Describe("Route_Online_Debug_Route_Not_Exist", func() {
	table.DescribeTable("Route_Online_Debug_Route_Not_Exist",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		table.Entry("hit route that not exist", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		}),
		table.Entry("online debug route with query params", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/debug-request-forwarding",
			Headers: map[string]string{
				"Authorization":                 base.GetToken(),
				"online_debug_url":              base.APISIXInternalUrl + `/hello_`,
				"online_debug_request_protocol": "http",
				"online_debug_method":           "GET",
				"Content-Type":                  "multipart/form-data",
				"online_debug_header_params":    `{"test":["test1"]}`,
			},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"code":0,"message":"","data":{"code":404,"header":{"Connection":["keep-alive"],"Content-Type":["text/plain; charset=utf-8"]`,
			Sleep:        base.SleepTime,
		}),
	)
})

var _ = ginkgo.Describe("Route_Online_Debug_Route_With_Query_Params", func() {
	ginkgo.It("hit route that not exist", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		})
	})
	ginkgo.It("create route with query params", func() {
		var routeBody map[string]interface{} = map[string]interface{}{
			"name":    "route1",
			"uri":     "/hello",
			"methods": []string{"GET"},
			"vars": []interface{}{
				[]string{"arg_name", "==", "aaa"},
			},
			"upstream": upstream,
		}
		_routeBody, err := json.Marshal(routeBody)
		gomega.Expect(err).To(gomega.BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/routes/r1",
			Body:         string(_routeBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("online debug route with query params", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/debug-request-forwarding",
			Headers: map[string]string{
				"Authorization":                 base.GetToken(),
				"online_debug_url":              base.APISIXInternalUrl + `/hello?name=aaa`,
				"online_debug_request_protocol": "http",
				"online_debug_method":           "GET",
				"Content-Type":                  "multipart/form-data",
				"online_debug_header_params":    `{"test":["test1"]}`,
			},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"code":0,"message":"","data":{"code":200,"header":{"Connection":["keep-alive"],"Content-Type":["application/octet-stream"],`,
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("delete the route just created", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("hit the route just deleted", func() {
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

var _ = ginkgo.Describe("Route_Online_Debug_Route_With_Header_Params", func() {
	ginkgo.It("hit route that not exist", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		})
	})
	ginkgo.It("create route with header params", func() {
		var routeBody map[string]interface{} = map[string]interface{}{
			"name":    "route1",
			"uri":     "/hello",
			"methods": []string{"GET"},
			"vars": []interface{}{
				[]string{"http_version", "==", "v2"},
			},
			"upstream": upstream,
		}
		_reqRouteBody, err := json.Marshal(routeBody)
		gomega.Expect(err).To(gomega.BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/routes/r1",
			Body:         string(_reqRouteBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("online debug route with header params", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/debug-request-forwarding",
			Headers: map[string]string{
				"Authorization":                 base.GetToken(),
				"online_debug_url":              base.APISIXInternalUrl + `/hello`,
				"online_debug_request_protocol": "http",
				"online_debug_method":           "GET",
				"Content-Type":                  "multipart/form-data",
				"online_debug_header_params":    `{"test":["test1"],"version":["v2"]}`,
			},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"code":0,"message":"","data":{"code":200,"header":{"Connection":["keep-alive"],"Content-Type":["application/octet-stream"],`,
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("online debug route with header params(add Content-type to header params to create route)", func() {
		var routeBody map[string]interface{} = map[string]interface{}{
			"name":     "route2",
			"status":   1,
			"uri":      "/hello_",
			"methods":  []string{"GET"},
			"upstream": upstream,
		}
		_reqRouteBody, err := json.Marshal(routeBody)
		gomega.Expect(err).To(gomega.BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/debug-request-forwarding",
			Body:   string(_reqRouteBody),
			Headers: map[string]string{
				"Authorization":                 base.GetToken(),
				"online_debug_url":              base.ManagerAPIHost + `/apisix/admin/routes/r2`,
				"online_debug_request_protocol": "http",
				"online_debug_method":           http.MethodPut,
				"Content-Type":                  "text/plain;charset=UTF-8",
				"online_debug_header_params":    `{"Content-type":["application/json"],"Authorization":["` + base.GetToken() + `"]}`,
			},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"code":0,"message":"","data":{"code":200,"header":{"Access-Control-Allow-Credentials":["true"],"Access-Control-Allow-Headers":["Authorization"],"Access-Control-Allow-Methods":["*"],"Access-Control-Allow-Origin":["*"],"Content-Length":["296"],"Content-Type":["application/json"]`,
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("hit the route (r2)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world\n",
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("delete the route just created (r1)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("delete the route just created (r2)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r2",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("hit the route just deleted (r1)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("hit the route just deleted  (r2)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		})
	})
})

var _ = ginkgo.Describe("Route_Online_Debug_Route_With_Body_Params", func() {
	ginkgo.It("hit route that not exist", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		})
	})
	ginkgo.It("create route with method POST", func() {
		var routeBody map[string]interface{} = map[string]interface{}{
			"name":     "route1",
			"uri":      "/hello",
			"methods":  []string{"POST"},
			"upstream": upstream,
		}
		_reqRouteBody, err := json.Marshal(routeBody)
		gomega.Expect(err).To(gomega.BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/routes/r1",
			Body:         string(_reqRouteBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("online debug route with body params", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/debug-request-forwarding",
			Body: `{
				"name": "test",
				"desc": "online debug route with body params"
			}`,
			Headers: map[string]string{
				"Authorization":                 base.GetToken(),
				"online_debug_url":              base.APISIXInternalUrl + `/hello`,
				"online_debug_request_protocol": "http",
				"online_debug_method":           http.MethodPost,
				"Content-Type":                  "application/json",
				"online_debug_header_params":    `{"test":["test1"]}`,
			},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"code":0,"message":"","data":{"code":200,"header":{"Connection":["keep-alive"],"Content-Type":["application/octet-stream"]`,
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("delete the route just created", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("hit the route just deleted", func() {
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
var _ = ginkgo.Describe("Route_Online_Debug_Route_With_Basic_Auth", func() {
	ginkgo.It("hit route that not exist", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		})
	})
	ginkgo.It("create route enable basic-auth plugin", func() {
		var routeBody map[string]interface{} = map[string]interface{}{
			"name": "route1",
			"uri":  "/hello",
			"plugins": map[string]interface{}{
				"basic-auth": map[string]interface{}{},
			},
			"methods":  []string{"GET"},
			"upstream": upstream,
		}
		_reqRouteBody, err := json.Marshal(routeBody)
		gomega.Expect(err).To(gomega.BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/routes/r1",
			Body:         string(_reqRouteBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("make sure the consumer is not created", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		})
	})
	ginkgo.It("create consumer", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/consumers",
			Body: `{
				"username": "jack",
				"plugins": {
					"basic-auth": {
						"disable": false,
						"username": "jack",
						"password": "123456"
					}
				},
				"desc": "test description"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("online debug with basic-auth", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/debug-request-forwarding",
			Headers: map[string]string{
				"Authorization":                 base.GetToken(),
				"online_debug_url":              base.APISIXInternalUrl + `/hello`,
				"online_debug_request_protocol": "http",
				"online_debug_method":           "GET",
				"Content-Type":                  "multipart/form-data",
				"online_debug_header_params":    `{"test":["test1"],"Authorization": ["Basic amFjazoxMjM0NTYKIA=="]}`,
			},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"code":0,"message":"","data":{"code":200,"header":{"Connection":["keep-alive"],"Content-Type":["application/octet-stream"]`,
		})
	})
	ginkgo.It("online debug without basic-auth", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/debug-request-forwarding",
			Headers: map[string]string{
				"Authorization":                 base.GetToken(),
				"online_debug_url":              base.APISIXInternalUrl + `/hello`,
				"online_debug_request_protocol": "http",
				"online_debug_method":           "GET",
				"Content-Type":                  "multipart/form-data",
				"online_debug_header_params":    `{"test":["test1"]}`,
			},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"code":0,"message":"","data":{"code":401,"header":{"Connection":["keep-alive"],"Content-Type":["text/plain; charset=utf-8"],`,
		})
	})
	ginkgo.It("delete the route just created", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("hit the route just deleted", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("delete consumer", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
})

var _ = ginkgo.Describe("Route_Online_Debug_Route_With_Key_Auth", func() {

	ginkgo.It("hit route that not exist", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		})
	})
	ginkgo.It("create route enable key-auth plugin", func() {
		var routeBody map[string]interface{} = map[string]interface{}{
			"name": "route1",
			"uri":  "/hello",
			"plugins": map[string]interface{}{
				"key-auth": map[string]interface{}{},
			},
			"methods":  []string{"GET"},
			"upstream": upstream,
		}
		_reqRouteBody, err := json.Marshal(routeBody)
		gomega.Expect(err).To(gomega.BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/routes/r1",
			Body:         string(_reqRouteBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("make sure the consumer is not created", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		})
	})
	ginkgo.It("create consumer", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/consumers",
			Body: `{
				"username": "jack",
				"plugins": {
					"key-auth": {
						"key": "user-key"
					}
				},
				"desc": "test description"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("online debug with key-auth", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/debug-request-forwarding",
			Headers: map[string]string{
				"Authorization":                 base.GetToken(),
				"online_debug_url":              base.APISIXInternalUrl + `/hello`,
				"online_debug_request_protocol": "http",
				"online_debug_method":           "GET",
				"Content-Type":                  "multipart/form-data",
				"online_debug_header_params":    `{"test":["test1"],"apikey":["user-key"]}`,
			},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"code":0,"message":"","data":{"code":200,"header":{"Connection":["keep-alive"],"Content-Type":["application/octet-stream"],`,
		})
	})
	ginkgo.It("online debug without key-auth", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/debug-request-forwarding",
			Headers: map[string]string{
				"Authorization":                 base.GetToken(),
				"online_debug_url":              base.APISIXInternalUrl + `/hello`,
				"online_debug_request_protocol": "http",
				"online_debug_method":           "GET",
				"Content-Type":                  "multipart/form-data",
				"online_debug_header_params":    `{"test":["test1"]}`,
			},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"code":0,"message":"","data":{"code":401,"header":{"Connection":["keep-alive"],"Content-Type":["text/plain; charset=utf-8"]`,
		})
	})
	ginkgo.It("delete the route just created", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("hit the route just deleted", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("delete consumer", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
})

var _ = ginkgo.Describe("Route_Online_Debug_Route_With_JWT_Auth", func() {
	ginkgo.It("hit route that not exist", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		})
	})
	ginkgo.It("create route enable jwt-auth plugin", func() {
		var routeBody map[string]interface{} = map[string]interface{}{
			"name": "route1",
			"uri":  "/hello",
			"plugins": map[string]interface{}{
				"jwt-auth": map[string]interface{}{},
			},
			"methods":  []string{"GET"},
			"upstream": upstream,
		}
		_reqRouteBody, err := json.Marshal(routeBody)
		gomega.Expect(err).To(gomega.BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/routes/r1",
			Body:         string(_reqRouteBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("make sure the consumer is not created", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		})
	})
	ginkgo.It("create consumer", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/consumers",
			Body: `{
				"username": "jack",
				"plugins": {
					"jwt-auth": {
						"key": "user-key",
						"secret": "my-secret-key",
						"algorithm": "HS256"
					}
				},
				"desc": "test description"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("online debug with JWT-auth", func() {
		jsonStr := `{"test":["test1"]}`
		var _headerParams map[string]interface{}
		err := json.Unmarshal([]byte(jsonStr), &_headerParams)
		gomega.Expect(err).To(gomega.BeNil())
		jwtToken := base.GetJwtToken("user-key")
		l := []string{jwtToken}
		_headerParams["Authorization"] = l
		headerParams, err := json.Marshal(_headerParams)
		gomega.Expect(err).To(gomega.BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/debug-request-forwarding",
			Headers: map[string]string{
				"Authorization":                 base.GetToken(),
				"online_debug_url":              base.APISIXInternalUrl + `/hello`,
				"online_debug_request_protocol": "http",
				"online_debug_method":           "GET",
				"Content-Type":                  "multipart/form-data",
				"online_debug_header_params":    string(headerParams),
			},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"code":0,"message":"","data":{"code":200,"header":{"Connection":["keep-alive"],"Content-Type":["application/octet-stream"],`,
		})
	})
	ginkgo.It("online debug without JWT-auth", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/debug-request-forwarding",
			Headers: map[string]string{
				"Authorization":                 base.GetToken(),
				"online_debug_url":              base.APISIXInternalUrl + `/hello`,
				"online_debug_request_protocol": "http",
				"online_debug_method":           "GET",
				"Content-Type":                  "multipart/form-data",
				"online_debug_header_params":    `{"test":["test1"]}`,
			},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"code":0,"message":"","data":{"code":401,"header":{"Connection":["keep-alive"],"Content-Type":["text/plain; charset=utf-8"],`,
		})
	})
	ginkgo.It("delete the route just created", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("hit the route just deleted", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("delete consumer", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
})

var _ = ginkgo.Describe("Route_Online_Debug_Route_With_Files", func() {
	ginkgo.It("hit route that not exist", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
		})
	})
	ginkgo.It("create route enable basic-auth plugin", func() {
		var routeBody map[string]interface{} = map[string]interface{}{
			"name":     "route1",
			"uri":      "/hello_",
			"methods":  []string{"POST"},
			"upstream": upstream,
		}
		_reqRouteBody, err := json.Marshal(routeBody)
		gomega.Expect(err).To(gomega.BeNil())
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/routes/r2",
			Body:         string(_reqRouteBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		})
	})
	ginkgo.It("online debug with file", func() {
		path, err := filepath.Abs("../../testdata/import/default.yaml")
		gomega.Expect(err).To(gomega.BeNil())
		files := []base.UploadFile{
			{Name: "file", Filepath: path},
		}

		headers := map[string]string{}

		jsonStr := `{"test":["test1"]}`
		var _headerParams map[string]interface{}
		err = json.Unmarshal([]byte(jsonStr), &_headerParams)
		gomega.Expect(err).To(gomega.BeNil())
		l := []string{base.GetToken()}
		_headerParams["Authorization"] = l
		headerParams, err := json.Marshal(_headerParams)
		gomega.Expect(err).To(gomega.BeNil())

		basePath := base.ManagerAPIHost + "/apisix/admin/debug-request-forwarding"
		requestBody, requestContentType, err := base.GetReader(headers, "multipart/form-data", files)
		gomega.Expect(err).To(gomega.BeNil())
		httpRequest, err := http.NewRequest(http.MethodPost, basePath, requestBody)
		gomega.Expect(err).To(gomega.BeNil())
		httpRequest.Header.Add("Content-Type", requestContentType)
		httpRequest.Header.Add("Authorization", base.GetToken())
		httpRequest.Header.Add("online_debug_request_protocol", "http")
		httpRequest.Header.Add("online_debug_url", base.ManagerAPIHost+`/apisix/admin/import/routes`)
		httpRequest.Header.Add("online_debug_method", http.MethodPost)
		httpRequest.Header.Add("online_debug_header_params", string(headerParams))
		client := &http.Client{}
		resp, err := client.Do(httpRequest)
		gomega.Expect(err).To(gomega.BeNil())

		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		gomega.Expect(err).To(gomega.BeNil())
		realBody := gjson.Get(string(respBody), "data")
		// todo get successful result and compare
		gomega.Expect(realBody.String()).Should(gomega.ContainSubstring(`"data":{"paths":1,"routes":1}`))
	})

	ginkgo.It("verify the route just imported and delete data", func() {
		time.Sleep(time.Duration(500) * time.Millisecond)
		request, _ := http.NewRequest("GET", base.ManagerAPIHost+"/apisix/admin/routes", nil)
		request.Header.Add("Authorization", base.GetToken())
		resp, err := http.DefaultClient.Do(request)
		gomega.Expect(err).To(gomega.BeNil())
		defer resp.Body.Close()
		respBody, _ := ioutil.ReadAll(resp.Body)
		list := gjson.Get(string(respBody), "data.rows").Value().([]interface{})

		var tests []base.HttpTestCase
		for _, item := range list {
			route := item.(map[string]interface{})
			tc := base.HttpTestCase{
				Desc:         "route patch for update status(online)",
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodPatch,
				Path:         "/apisix/admin/routes/" + route["id"].(string),
				Body:         `{"status":1}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				Sleep:        base.SleepTime,
			}
			tests = append(tests, tc)
		}

		// verify route
		tests = append(tests, base.HttpTestCase{
			Desc:         "verify the route just imported",
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		})

		// delete test data
		for _, item := range list {
			route := item.(map[string]interface{})
			tc := base.HttpTestCase{
				Desc:         "delete route",
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodDelete,
				Path:         "/apisix/admin/routes/" + route["id"].(string),
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			}
			tests = append(tests, tc)
		}

		for _, tc := range tests {
			base.RunTestCase(tc)
		}

	})

	ginkgo.It("hit the route just deleted", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		})
	})

	ginkgo.It("hit the route just deleted", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        base.SleepTime,
		})
	})
})
