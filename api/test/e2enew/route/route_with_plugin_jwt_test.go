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
package route

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/gomega"

	"github.com/apisix/manager-api/test/e2enew/base"
)

var _ = ginkgo.Describe("route with jwt plugin", func() {
	var jwtToken string

	table.DescribeTable("create route and consumer with jwt plugin",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		table.Entry("make sure the route is not created ", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
		}),
		table.Entry("create route", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				 "name": "route1",
				 "uri": "/hello",
				 "plugins": {
					 "jwt-auth": {}
				 },
				 "upstream": {
					"type": "roundrobin",
					"nodes": {
						"` + base.UpstreamIp + `:1981": 1
					}
				 }
			 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{`"code":0`, `"id":"r1"`, `"uri":"/hello"`, `"name":"route1"`, `"jwt-auth":{}`},
		}),
		table.Entry("make sure the consumer is not created", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		}),
		table.Entry("create consumer", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Path:   "/apisix/admin/consumers",
			Method: http.MethodPut,
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
			ExpectBody:   []string{`"code":0`, `"username":"jack"`, `"key":"user-key"`, `"secret":"my-secret-key"`},
		}),
	)
	ginkgo.It("sign jwt token", func() {
		time.Sleep(base.SleepTime)

		// sign jwt token
		body, status, err := base.HttpGet(base.APISIXHost+"/apisix/plugin/jwt/sign?key=user-key", nil)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(status).To(gomega.Equal(http.StatusOK))
		jwtToken = string(body)
		// sign jwt token with not exists key
		body, status, err = base.HttpGet(base.APISIXHost+"/apisix/plugin/jwt/sign?key=not-exist-key", nil)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(status).To(gomega.Equal(http.StatusNotFound))
	})

	ginkgo.It("verify route with correct jwt token", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": jwtToken},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
		})
	})

	table.DescribeTable("verify token and clean consumer",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		table.Entry("verify route without jwt token", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusUnauthorized,
			ExpectBody:   `{"message":"Missing JWT token in request"}`,
			Sleep:        base.SleepTime,
		}),
		table.Entry("verify route with incorrect jwt token", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": "invalid-token"},
			ExpectStatus: http.StatusUnauthorized,
			ExpectBody:   `{"message":"invalid jwt string"}`,
		}),
		table.Entry("delete consumer", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
	)

	ginkgo.It("verify route with the jwt token from just deleted consumer", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": jwtToken},
			ExpectStatus: http.StatusUnauthorized,
			ExpectBody:   `{"message":"Missing related consumer"}`,
			Sleep:        base.SleepTime,
		})
	})

	table.DescribeTable("cleanup route and verify",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		table.Entry("delete route", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		table.Entry("verify the deleted route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		}),
	)

	table.DescribeTable("create route and consumer with jwt (no algorithm)",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		table.Entry("create consumer with jwt (no algorithm)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Path:   "/apisix/admin/consumers",
			Method: http.MethodPut,
			Body: `{
				"username":"consumer_1",
				"desc": "test description",
				"plugins":{
					"jwt-auth":{
						"exp":86400,
						"key":"user-key",
						"secret":"my-secret-key"
					}
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody: []string{`"code":0`, `"username":"consumer_1"`,
				`"jwt-auth":{"algorithm":"HS256","exp":86400,"key":"user-key","secret":"my-secret-key"}`},
		}),
		table.Entry("get the consumer", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Path:         "/apisix/admin/consumers/consumer_1",
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"username":"consumer_1"`,
			Sleep:        base.SleepTime,
		}),
		table.Entry("create the route", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"plugins": {
					"jwt-auth": {}
				},
				"upstream": {
					"type": "roundrobin",
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					}
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectBody:   []string{`"code":0`, `"id":"r1"`, `"uri":"/hello"`, `"name":"route1"`, `"jwt-auth":{}`},
			ExpectStatus: http.StatusOK,
		}),
	)

	ginkgo.It("get the jwt token", func() {
		time.Sleep(base.SleepTime)

		request, _ := http.NewRequest("GET", base.APISIXHost+"/apisix/plugin/jwt/sign?key=user-key", nil)
		resp, err := http.DefaultClient.Do(request)
		gomega.Expect(err).To(gomega.BeNil())
		defer resp.Body.Close()
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))
		jwtTokenBytes, _ := ioutil.ReadAll(resp.Body)
		jwtToken = string(jwtTokenBytes)
	})

	ginkgo.It("hit route with jwt token", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": jwtToken},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		})
	})

	table.DescribeTable("cleanup consumer and route",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		table.Entry("delete consumer", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Path:         "/apisix/admin/consumers/consumer_1",
			Method:       http.MethodDelete,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"code":0`,
		}),
		table.Entry("after delete consumer verify it again", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `"message":"data not found"`,
			Sleep:        base.SleepTime,
		}),
		table.Entry("delete the route", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		table.Entry("verify the deleted route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		}),
	)
})
