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
package stream_route_test

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Stream Route", func() {
	DescribeTable("Test stream route Integrate (HTTP upstream)",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("Create stream route", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/stream_routes/1",
			Body: `{
				"server_port": 10090,
				"upstream": {
        			"nodes": {
            			"` + base.UpstreamIp + `:1980": 1
        			},
        			"type": "roundrobin"
    			}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("Hit stream route", base.HttpTestCase{
			Object:       base.APISIXStreamProxyExpect(10090, ""),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
		}),
		Entry("Delete stream route", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/stream_routes/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
	)

	// prepare ssl certificate
	apisixCert, err := ioutil.ReadFile("../../certs/apisix.crt")
	Expect(err).To(BeNil())
	apisixKey, err := ioutil.ReadFile("../../certs/apisix.key")
	Expect(err).To(BeNil())
	apisixSSLBody, err := json.Marshal(map[string]string{"cert": string(apisixCert), "key": string(apisixKey)})
	Expect(err).To(BeNil())
	DescribeTable("Test stream route Integrate (HTTPS upstream)",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("Create SSL certificate", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPost,
			Path:         "/apisix/admin/ssl",
			Body:         string(apisixSSLBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("Create stream route", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/stream_routes/1",
			Body: `{
				"server_port": 10093,
				"sni": "test.com",
				"upstream": {
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("Hit stream route through https", base.HttpTestCase{
			Object:       base.APISIXStreamProxyExpect(10093, "test.com"),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
		}),
	)

	Describe("Test stream route (TCP upstream)", func() {
		It("Create stream route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/stream_routes/sr1tcp",
				Body: `{
					"server_port": 10090,
					"upstream": {
						"nodes": {
							"` + base.UpstreamIp + `:1991": 1
						},
						"type": "roundrobin"
					}
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("Hit stream route through tcp", func() {
			time.Sleep(base.SleepTime)
			conn, err := net.Dial("tcp", "127.0.0.1:1991")
			Expect(err).To(BeNil())

			_ = conn.SetDeadline(time.Now().Add(time.Second * 3))

			_, err = conn.Write([]byte("world"))
			Expect(err).To(BeNil())

			result := make([]byte, 11)
			n, err := conn.Read(result)
			Expect(n).Should(BeNumerically("==", 11))
			Expect(err).To(BeNil())
			Expect(string(result)).To(ContainSubstring("hello world"))

			err = conn.Close()
			Expect(err).To(BeNil())
		})
	})

	Describe("Test stream route (UDP upstream)", func() {
		It("Create stream route", func() {
			base.RunTestCase(base.HttpTestCase{
				Object: base.ManagerApiExpect(),
				Method: http.MethodPut,
				Path:   "/apisix/admin/stream_routes/sr1udp",
				Body: `{
					"id": "sr1udp",
					"server_port": 10095,
					"upstream": {
						"nodes": {
							"` + base.UpstreamIp + `:1992": 1
						},
						"type": "roundrobin"
					}
				}`,
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		})
		It("Hit stream route through udp", func() {
			time.Sleep(base.SleepTime)
			conn, err := net.Dial("udp", "127.0.0.1:10095")
			Expect(err).To(BeNil())

			_ = conn.SetDeadline(time.Now().Add(time.Second * 3))

			_, err = conn.Write([]byte("world"))
			Expect(err).To(BeNil())

			result := make([]byte, 11)
			n, err := conn.Read(result)
			Expect(n).Should(BeNumerically("==", 11))
			Expect(err).To(BeNil())
			Expect(string(result)).To(ContainSubstring("hello world"))

			err = conn.Close()
			Expect(err).To(BeNil())
		})
	})
})
