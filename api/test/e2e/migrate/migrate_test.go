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

package migrate_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"hash/crc32"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

const (
	checksumLength = 4 // 4 bytes (uint32)
)

type AllData struct {
	Consumers     []interface{}
	Routes        []interface{}
	Services      []interface{}
	SSLs          []interface{}
	Upstreams     []interface{}
	Scripts       []interface{}
	GlobalPlugins []interface{}
	PluginConfigs []interface{}
}

type response struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
	Data    struct {
		ConflictItems AllData `json:"ConflictItems"`
	} `json:"Data"`
}

var _ = Describe("Migrate", func() {
	var exportData []byte

	DescribeTable("Prepare config data",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("create test route", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route",
				"uri": "/hello_",
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
		Entry("create test upstream", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/upstreams/u1",
			Body: `{
				"name": "upstream",
				"nodes": [
					{
						"host": "` + base.UpstreamIp + `",
						"port": 1980,
						"weight": 1
					}
				],
				"type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("create test service", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/services/s1",
			Body: `{
  				"name": "service",
  				"upstream": {
    				"nodes": [
      					{
        					"host": "` + base.UpstreamIp + `",
        					"port": 1980,
        					"weight": 1
      					},
      					{
							"host": "` + base.UpstreamIp + `",
        					"port": 1981,
        					"weight": 2
      					},
      					{
        					"host": "` + base.UpstreamIp + `",
        					"port": 1982,
        					"weight": 3
      					}
    				],
    				"type": "roundrobin"
  				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        time.Second * 1,
		}),
		Entry("migrate export auth test", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPost,
			Path:         "/apisix/admin/migrate/export",
			ExpectStatus: http.StatusUnauthorized,
			ExpectBody:   "request unauthorized",
			Sleep:        base.SleepTime,
		}),
		Entry("migrate import auth test", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPost,
			Path:         "/apisix/admin/migrate/import",
			ExpectStatus: http.StatusUnauthorized,
			ExpectBody:   "request unauthorized",
			Sleep:        base.SleepTime,
		}),
	)

	It("export config success", func() {
		req := base.ManagerApiExpect().GET("/apisix/admin/migrate/export")
		req.WithHeader("Authorization", base.GetToken())
		resp := req.Expect()
		resp.Status(http.StatusOK)
		exportData = []byte(resp.Body().Raw())
		data := exportData[:len(exportData)-checksumLength]
		checksum := binary.BigEndian.Uint32(exportData[len(exportData)-checksumLength:])
		Expect(checksum).Should(Equal(crc32.ChecksumIEEE(data)))
	})

	It("import config conflict and return", func() {
		req := base.ManagerApiExpect().POST("/apisix/admin/migrate/import")
		buffer := bytes.NewBuffer(exportData)
		req.WithMultipart().WithForm(map[string]string{"mode": "return"})
		req.WithMultipart().WithFile("file", "apisix-config.bak", buffer)
		req.WithHeader("Authorization", base.GetToken())
		resp := req.Expect()
		resp.Status(http.StatusOK)
		rsp := &response{}
		err := json.Unmarshal([]byte(resp.Body().Raw()), rsp)
		Expect(err).Should(BeNil())
		Expect(rsp.Code).Should(Equal(20001))
		Expect(len(rsp.Data.ConflictItems.Routes)).Should(Equal(1))
		Expect(len(rsp.Data.ConflictItems.Upstreams)).Should(Equal(1))
		Expect(len(rsp.Data.ConflictItems.Services)).Should(Equal(1))
	})

	It("import config conflict and skip", func() {
		req := base.ManagerApiExpect().POST("/apisix/admin/migrate/import")
		buffer := bytes.NewBuffer(exportData)
		req.WithMultipart().WithForm(map[string]string{"mode": "skip"})
		req.WithMultipart().WithFile("file", "apisix-config.bak", buffer)
		req.WithHeader("Authorization", base.GetToken())
		resp := req.Expect()
		resp.Status(http.StatusOK)
		rsp := &response{}
		err := json.Unmarshal([]byte(resp.Body().Raw()), rsp)
		Expect(err).Should(BeNil())
		Expect(rsp.Code).Should(Equal(0))
	})

	It("import config conflict and overwrite", func() {
		req := base.ManagerApiExpect().POST("/apisix/admin/migrate/import")
		buffer := bytes.NewBuffer(exportData)
		req.WithMultipart().WithForm(map[string]string{"mode": "overwrite"})
		req.WithMultipart().WithFile("file", "apisix-config.bak", buffer)
		req.WithHeader("Authorization", base.GetToken())
		resp := req.Expect()
		resp.Status(http.StatusOK)
		rsp := &response{}
		err := json.Unmarshal([]byte(resp.Body().Raw()), rsp)
		Expect(err).Should(BeNil())
		Expect(rsp.Code).Should(Equal(0))
	})

	It("request hit route r1", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		})
	})

	DescribeTable("Cleanup config data",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("remove test route", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("remove test upstream", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/u1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("remove test service", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/services/s1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
	)

	It("delete imported route failed", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		})
	})

	It("request route r1 not found", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusNotFound,
			Sleep:        base.SleepTime,
		})
	})

	It("import config success", func() {
		req := base.ManagerApiExpect().POST("/apisix/admin/migrate/import")
		buffer := bytes.NewBuffer(exportData)
		req.WithMultipart().WithForm(map[string]string{"mode": "return"})
		req.WithMultipart().WithFile("file", "apisix-config.bak", buffer)
		req.WithHeader("Authorization", base.GetToken())
		resp := req.Expect()
		resp.Status(http.StatusOK)
		rsp := &response{}
		err := json.Unmarshal([]byte(resp.Body().Raw()), rsp)
		Expect(err).Should(BeNil())
		Expect(rsp.Code).Should(Equal(0))
	})

	It("Wait for APISIX re-sync", func() {
		time.Sleep(5 * time.Second)
	})

	It("request hit route r1", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		})
	})

	It("delete imported route success", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})

	It("delete imported upstream success", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/u1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})

	It("delete imported service success", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/services/s1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})

})
