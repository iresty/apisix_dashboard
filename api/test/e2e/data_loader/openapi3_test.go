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
package data_loader_test

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/savsgio/gotils/bytes"
	"github.com/tidwall/gjson"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("OpenAPI 3", func() {
	DescribeTable("Import cases",
		func(f func()) {
			f()
		},
		Entry("default.yaml", func() {
			path, err := filepath.Abs("../../testdata/import/default.yaml")
			Expect(err).To(BeNil())

			req := base.ManagerApiExpect().POST("/apisix/admin/import/routes")
			req.WithMultipart().WithForm(map[string]string{
				"type":      "openapi3",
				"task_name": "test_default_yaml",
				"_file":     "default.yaml",
			})
			req.WithMultipart().WithFile("file", path)
			req.WithHeader("Authorization", base.GetToken())
			resp := req.Expect()
			resp.Status(http.StatusOK)
			r := gjson.ParseBytes([]byte(resp.Body().Raw()))

			Expect(r.Get("code").Uint()).To(Equal(uint64(0)))

			r = r.Get("data")
			for s, result := range r.Map() {
				if s == "route" {
					Expect(result.Get("total").Uint()).To(Equal(uint64(1)))
					Expect(result.Get("failed").Uint()).To(Equal(uint64(0)))
				}
			}
		}),
		Entry("default.json", func() {
			path, err := filepath.Abs("../../testdata/import/default.json")
			Expect(err).To(BeNil())

			req := base.ManagerApiExpect().POST("/apisix/admin/import/routes")
			req.WithMultipart().WithForm(map[string]string{
				"type":      "openapi3",
				"task_name": "test_default_json",
				"_file":     "default.json",
			})
			req.WithMultipart().WithFile("file", path)
			req.WithHeader("Authorization", base.GetToken())
			resp := req.Expect()
			resp.Status(http.StatusOK)
			r := gjson.ParseBytes([]byte(resp.Body().Raw()))

			Expect(r.Get("code").Uint()).To(Equal(uint64(0)))

			r = r.Get("data")
			for s, result := range r.Map() {
				if s == "route" {
					Expect(result.Get("total").Uint()).To(Equal(uint64(1)))
					Expect(result.Get("failed").Uint()).To(Equal(uint64(0)))
				}
			}
		}),
		Entry("Postman-API101.yaml merge method", func() {
			path, err := filepath.Abs("../../testdata/import/Postman-API101.yaml")
			Expect(err).To(BeNil())

			req := base.ManagerApiExpect().POST("/apisix/admin/import/routes")
			req.WithMultipart().WithForm(map[string]string{
				"type":         "openapi3",
				"task_name":    "test_postman_api101_yaml_mm",
				"_file":        "Postman-API101.yaml",
				"merge_method": "true",
			})
			req.WithMultipart().WithFile("file", path)
			req.WithHeader("Authorization", base.GetToken())
			resp := req.Expect()
			resp.Status(http.StatusOK)
			r := gjson.ParseBytes([]byte(resp.Body().Raw()))

			Expect(r.Get("code").Uint()).To(Equal(uint64(0)))

			r = r.Get("data")
			for s, result := range r.Map() {
				if s == "route" {
					Expect(result.Get("total").Uint()).To(Equal(uint64(3)))
					Expect(result.Get("failed").Uint()).To(Equal(uint64(0)))
				}
				if s == "upstream" {
					Expect(result.Get("total").Uint()).To(Equal(uint64(1)))
					Expect(result.Get("failed").Uint()).To(Equal(uint64(0)))
				}
			}
		}),
		Entry("Postman-API101.yaml non-merge method", func() {
			// clean routes
			base.CleanResource("routes")
			path, err := filepath.Abs("../../testdata/import/Postman-API101.yaml")
			Expect(err).To(BeNil())

			req := base.ManagerApiExpect().POST("/apisix/admin/import/routes")
			req.WithMultipart().WithForm(map[string]string{
				"type":         "openapi3",
				"task_name":    "test_postman_api101_yaml_nmm",
				"_file":        "Postman-API101.yaml",
				"merge_method": "false",
			})
			req.WithMultipart().WithFile("file", path)
			req.WithHeader("Authorization", base.GetToken())
			resp := req.Expect()
			resp.Status(http.StatusOK)
			r := gjson.ParseBytes([]byte(resp.Body().Raw()))

			Expect(r.Get("code").Uint()).To(Equal(uint64(0)))

			r = r.Get("data")
			for s, result := range r.Map() {
				if s == "route" {
					Expect(result.Get("total").Uint()).To(Equal(uint64(5)))
					Expect(result.Get("failed").Uint()).To(Equal(uint64(0)))
				}
				if s == "upstream" {
					Expect(result.Get("total").Uint()).To(Equal(uint64(1)))
					Expect(result.Get("failed").Uint()).To(Equal(uint64(0)))
				}
			}
		}),
		Entry("Clean resources", func() {
			base.CleanResource("routes")
			base.CleanResource("upstreams")
			base.CleanResource("services")
		}),
	)
	DescribeTable("Exception cases",
		func(f func()) {
			f()
		},
		Entry("Empty upload file", func() {
			req := base.ManagerApiExpect().POST("/apisix/admin/import/routes")
			req.WithMultipart().WithForm(map[string]string{
				"type":      "openapi3",
				"task_name": "empty_upload",
				"_file":     "default.yaml",
			})
			req.WithHeader("Authorization", base.GetToken())
			resp := req.Expect()
			resp.Status(http.StatusOK)
			r := gjson.ParseBytes([]byte(resp.Body().Raw()))

			Expect(r.Get("code").Uint()).To(Equal(uint64(10000)))
			Expect(r.Get("message").String()).To(Equal("uploaded file is empty"))
		}),
		Entry("Exceed limit upload file", func() {
			req := base.ManagerApiExpect().POST("/apisix/admin/import/routes")
			req.WithMultipart().WithForm(map[string]string{
				"type":      "openapi3",
				"task_name": "exceed_limit_upload",
				"_file":     "default.yaml",
			})

			req.WithMultipart().WithFileBytes("file", "default.yaml", bytes.Rand(make([]byte, 10*1024*1024+1)))
			req.WithHeader("Authorization", base.GetToken())
			resp := req.Expect()
			resp.Status(http.StatusOK)
			r := gjson.ParseBytes([]byte(resp.Body().Raw()))

			Expect(r.Get("code").Uint()).To(Equal(uint64(10000)))
			Expect(r.Get("message").String()).To(Equal("uploaded file size exceeds the limit, limit is 10485760"))
		}),
		Entry("Routes duplicate #1", func() {
			path, err := filepath.Abs("../../testdata/import/Postman-API101.yaml")
			Expect(err).To(BeNil())
			req := base.ManagerApiExpect().POST("/apisix/admin/import/routes")
			req.WithMultipart().WithForm(map[string]string{
				"type":         "openapi3",
				"task_name":    "duplicate",
				"_file":        "Postman-API101.yaml",
				"merge_method": "true",
			})
			req.WithMultipart().WithFile("file", path)
			req.WithHeader("Authorization", base.GetToken())
			resp := req.Expect()
			resp.Status(http.StatusOK)
			r := gjson.ParseBytes([]byte(resp.Body().Raw()))
			Expect(r.Get("code").Uint()).To(Equal(uint64(0)))
		}),
		Entry("Route duplicate #2", func() {
			path, err := filepath.Abs("../../testdata/import/Postman-API101.yaml")
			Expect(err).To(BeNil())
			req := base.ManagerApiExpect().POST("/apisix/admin/import/routes")
			req.WithMultipart().WithForm(map[string]string{
				"type":         "openapi3",
				"task_name":    "duplicate",
				"_file":        "Postman-API101.yaml",
				"merge_method": "true",
			})
			req.WithMultipart().WithFile("file", path)
			req.WithHeader("Authorization", base.GetToken())
			resp := req.Expect()
			resp.Status(http.StatusOK)
			r := gjson.ParseBytes([]byte(resp.Body().Raw()))
			Expect(r.Get("code").Uint()).To(Equal(uint64(0)))
			Expect(r.Get("data").Map()["route"].Get("failed").Uint()).To(Equal(uint64(1)))
			Expect(r.Get("data").Map()["route"].Get("errors").Array()[0].String()).
				To(ContainSubstring("is duplicated with route duplicate_"))

		}),
		Entry("Clean resources", func() {
			base.CleanResource("routes")
			base.CleanResource("upstreams")
			base.CleanResource("services")
		}),
	)
	DescribeTable("Real API cases",
		func(f func()) {
			f()
		},
		Entry("Import httpbin.org YAML", func() {
			path, err := filepath.Abs("../../testdata/import/httpbin.yaml")
			Expect(err).To(BeNil())

			req := base.ManagerApiExpect().POST("/apisix/admin/import/routes")
			req.WithMultipart().WithForm(map[string]string{
				"type":      "openapi3",
				"task_name": "httpbin",
				"_file":     "httpbin.yaml",
			})
			req.WithMultipart().WithFile("file", path)
			req.WithHeader("Authorization", base.GetToken())
			resp := req.Expect()
			resp.Status(http.StatusOK)
			r := gjson.ParseBytes([]byte(resp.Body().Raw()))

			Expect(r.Get("code").Uint()).To(Equal(uint64(0)))

			r = r.Get("data")
			for s, result := range r.Map() {
				if s == "route" {
					Expect(result.Get("total").Uint()).To(Equal(uint64(1)))
					Expect(result.Get("failed").Uint()).To(Equal(uint64(0)))
				}
			}
		}),
		Entry("Modify upstream", func() {
			body := make(map[string]interface{})
			body["nodes"] = []map[string]interface{}{
				{
					"host":   "httpbin.org",
					"port":   80,
					"weight": 1,
				},
			}
			body["type"] = "roundrobin"
			_body, err := json.Marshal(body)
			Expect(err).To(BeNil())
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodPatch,
				Path:         "/apisix/admin/upstreams/httpbin",
				Body:         string(_body),
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		}),
		Entry("Enable route", func() {
			// get route id
			req := base.ManagerApiExpect().GET("/apisix/admin/routes")
			req.WithHeader("Authorization", base.GetToken())
			resp := req.Expect()
			resp.Status(http.StatusOK)
			r := gjson.ParseBytes([]byte(resp.Body().Raw()))
			Expect(r.Get("code").Uint()).To(Equal(uint64(0)))
			id := r.Get("data").Get("rows").Array()[0].Get("id").String()

			body := make(map[string]interface{})
			body["status"] = 1
			_body, err := json.Marshal(body)
			Expect(err).To(BeNil())
			base.RunTestCase(base.HttpTestCase{
				Object:       base.ManagerApiExpect(),
				Method:       http.MethodPatch,
				Path:         "/apisix/admin/routes/" + id,
				Body:         string(_body),
				Headers:      map[string]string{"Authorization": base.GetToken()},
				ExpectStatus: http.StatusOK,
			})
		}),
		Entry("Request API", func() {
			req := base.APISIXExpect().GET("/get")
			resp := req.Expect()
			resp.Status(http.StatusOK)
			r := gjson.ParseBytes([]byte(resp.Body().Raw()))
			Expect(r.Get("headers").Get("User-Agent").String()).To(Equal("Go-http-client/1.1"))
		}),
	)
})
