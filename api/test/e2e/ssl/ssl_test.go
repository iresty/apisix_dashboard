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
package ssl_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("SSL", func() {
	testCertByte, err := ioutil.ReadFile("../../certs/test2.crt")
	Expect(err).To(BeNil())
	testCert := strings.ReplaceAll(string(testCertByte), "\n", "\\n")

	testKeyByte, err := ioutil.ReadFile("../../certs/test2.key")
	Expect(err).To(BeNil())
	testKey := strings.ReplaceAll(string(testKeyByte), "\n", "\\n")

	DescribeTable("Test SSL CURD",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("Get SSL (Not Exist)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/ssl/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		}),
		Entry("List SSL (Empty)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/ssl",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":0`,
		}),
		Entry("Create SSL #1", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/ssl/1",
			Body: `{
				"cert": "` + testCert + `",
				"key": "` + testKey + `",
				"labels": {
					"build":   "16",
					"env":     "production",
					"version": "v3"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"version":"v3"`,
		}),
		Entry("Get SSL (Exist)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/ssl/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"build":"16"`,
		}),
		Entry("List SSL (1 item)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/ssl",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":1`,
		}),
		Entry("Create SSL #2", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/ssl",
			Body: `{
				"id": "2",
				"cert": "` + testCert + `",
				"key": "` + testKey + `",
				"labels": {
					"build":   "16",
					"env":     "production",
					"version": "v3"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"version":"v3"`,
		}),
		Entry("List SSL (2 items)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/ssl",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"total_size":2`,
		}),
		Entry("List SSL (Paginate)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/ssl",
			Query:        "page=2&page_size=1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"id":"2"`,
		}),
		Entry("Check Exist SSL", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPost,
			Path:         "/apisix/admin/check_ssl_exists",
			Body:         `{"hosts": ["www.test2.com"]}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("Update SSL (Full)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/ssl/1",
			Body: `{
				"cert": "` + testCert + `",
				"key": "` + testKey + `",
				"labels": {
					"build":   "16",
					"env":     "production",
					"version": "v2"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"version":"v2"`,
		}),
		Entry("Update SSL (Partial)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPatch,
			Path:   "/apisix/admin/ssl/2/labels",
			Body: `{
				"build":   "16",
				"env":     "production",
				"version": "v1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `"version":"v1"`,
		}),
		Entry("Batch Delete SSL", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/ssl/1,2",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
	)

	DescribeTable("Test SSL Integrate",
		func(f func() *base.HttpTestCase) {
			if c := f(); c != nil {
				base.RunTestCase(*c)
			}
		},
		Entry("Ensure SSL cleaned", func() *base.HttpTestCase {
			base.CleanResource("ssl")
			time.Sleep(time.Second)
			return &base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodGet,
				Path:         "/apisix/admin/ssl",
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   `"total_size":0`,
			}
		}),
		Entry("Request HTTPS (without certificate)", func() *base.HttpTestCase {
			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				if addr == "www.test2.com:9443" {
					addr = "127.0.0.1:9443"
				}
				dialer := &net.Dialer{}
				return dialer.DialContext(ctx, network, addr)
			}

			_, err := http.Get("https://www.test2.com:9443")
			Expect(fmt.Sprintf("%v", err)).Should(Equal(`Get "https://www.test2.com:9443": remote error: tls: internal error`))
			return nil
		}),
		Entry("Create SSL", func() *base.HttpTestCase {
			return &base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/ssl/1",
				Body: `{
					"cert": "` + testCert + `",
					"key": "` + testKey + `"
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   `"cert":"-----BEGIN CERTIFICATE-----`,
			}
		}),
		Entry("Create route", func() *base.HttpTestCase {
			return &base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/routes/r1",
				Body: `{
					"name": "route1",
					"uri": "/hello_",
					"hosts": ["test2.com", "*.test2.com"],
					"upstream": {
						"nodes": {
							"` + base.UpstreamIp + `:1980": 1
						},
						"type": "roundrobin"
					}
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
				ExpectBody:   []string{`"id":"r1"`, `"name":"route1"`, `"uri":"/hello_"`},
			}
		}),
		Entry("Hit route (with HTTPS)", func() *base.HttpTestCase {
			return &base.HttpTestCase{
				Object:       base.APISIXHTTPSExpect(),
				Method:       http.MethodGet,
				Path:         "/hello_",
				Headers:      map[string]string{"Host": "www.test2.com"},
				ExpectStatus: http.StatusOK,
				ExpectBody:   "hello world",
			}
		}),
	)
})
