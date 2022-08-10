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
package label

import (
	"net/http"

	. "github.com/onsi/ginkgo"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Test label", func() {
	It("Create route", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Path:   "/apisix/admin/routes/r1",
			Method: http.MethodPut,
			Body: `{
					 "name": "route1",
					 "uri": "/hello",
					 "labels": {
						 "build":"16",
						 "env":"production",
						 "version":"v2"
					 },
					 "upstream": {
						 "type": "roundrobin",
						 "nodes": [{
							 "host": "` + base.UpstreamIp + `",
							 "port": 1980,
							 "weight": 1
						 }]
					 }
				 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("Create consumer", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Path:   "/apisix/admin/consumers/c1",
			Method: http.MethodPut,
			Body: `{
					 "username": "c1",
					 "plugins": {
						 "key-auth": {
							 "key": "auth-one"
						 }
					 },
					 "labels": {
						 "build":"16",
						 "env":"production",
						 "version":"v3"
					 },
					 "desc": "test description"
				 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("Create upstream", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/upstreams/u1",
			Body: `{
					 "nodes": [{
						 "host": "` + base.UpstreamIp + `",
						 "port": 1980,
						 "weight": 1
					 }],
					 "labels": {
						 "build":"17",
						 "env":"production",
						 "version":"v2"
					 },
					 "type": "roundrobin"
				 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("Create service", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/services",
			Body: `{
					 "id": "s1",
					 "plugins": {
						 "limit-count": {
							 "count": 2,
							 "time_window": 60,
							 "rejected_code": 503,
							 "key": "remote_addr",
							 "policy": "local"
						 }
					 },
					 "upstream": {
						 "type": "roundrobin",
						 "nodes": [{
							 "host": "39.97.63.215",
							 "port": 80,
							 "weight": 1
						 }]
					 },
					 "labels": {
						 "build":"16",
						 "env":"production",
						 "version":"v2",
						 "extra": "test"
					 }
				 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("Create plugin config", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/plugin_configs/1",
			Body: `{
					 "plugins": {
						 "response-rewrite": {
							 "headers": {
								 "X-VERSION":"22.0"
							 }
						 }
					 },
					 "labels": {
						 "version": "v2",
						 "build":   "17",
						 "extra":   "test"
					 }
				 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("Get route label", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Path:         "/apisix/admin/labels/route",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"16"},{"env":"production"},{"version":"v2"}`,
			Sleep:        base.SleepTime,
		})
	})
	It("Get consumer label", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Path:         "/apisix/admin/labels/consumer",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"16"},{"env":"production"},{"version":"v3"}`,
		})
	})
	It("Get upstream label", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Path:         "/apisix/admin/labels/upstream",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"17"},{"env":"production"},{"version":"v2"}`,
		})
	})
	It("Get service label", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Path:         "/apisix/admin/labels/service",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"16"},{"env":"production"},{"extra":"test"},{"version":"v2"}`,
		})
	})
	It("Get plugin config label", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Path:         "/apisix/admin/labels/plugin_config",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"17"},{"extra":"test"},{"version":"v2"}`,
		})
	})
	It("Update plugin config", func() {
		base.RunTestCase(base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/plugin_configs/1",
			Body: `{
					 "plugins": {
						 "response-rewrite": {
							 "headers": {
								 "X-VERSION":"22.0"
							 }
						 }
					 },
					 "labels": {
						 "version": "v3",
						 "build":   "16",
						 "extra":   "test"
					 }
				 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("Get plugin config (Updated)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Path:         "/apisix/admin/labels/plugin_config",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"16"},{"extra":"test"},{"version":"v3"}`,
			Sleep:        base.SleepTime,
		})
	})
	It("List label", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Path:         "/apisix/admin/labels/all",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"16"},{"build":"17"},{"env":"production"},{"extra":"test"},{"version":"v2"},{"version":"v3"}`,
		})
	})
	It("List label (Paginate) #1", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Query:        "page=1&page_size=1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Path:         "/apisix/admin/labels/all",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"16"}`,
		})
	})
	It("List label (Paginate) #2", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Query:        "page=3&page_size=1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Path:         "/apisix/admin/labels/all",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"env":"production"}`,
		})
	})
	It("Get labels (With condition, key = build)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Query:        "label=build",
			Path:         "/apisix/admin/labels/all",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"16"},{"build":"17"}`,
		})
	})
	It("Get labels with the same key (With condition, key = build)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Query:        "label=build:16,build:17",
			Path:         "/apisix/admin/labels/all",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"16"},{"build":"17"}`,
		})
	})
	It("Get labels (With condition, key = build; Paginate)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Query:        "label=build&page=2&page_size=1",
			Path:         "/apisix/admin/labels/all",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"17"}`,
		})
	})
	It("Get labels with same key (With condition, key = build; Paginate)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Query:        "label=build:16,build:17&page=1&page_size=2",
			Path:         "/apisix/admin/labels/all",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"16"},{"build":"17"}`,
		})
	})
	It("Get labels with same key (With condition, key = build; Paginate)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Query:        "label=build:16,build:17&page=2&page_size=1",
			Path:         "/apisix/admin/labels/all",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"17"}`,
		})
	})
	It("Get labels (With condition, key = build && env = production)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Query:        "label=build,env:production",
			Path:         "/apisix/admin/labels/all",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"16"},{"build":"17"},{"env":"production"}`,
		})
	})
	It("Get labels (With condition, build=16|17 && env = production)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Query:        "label=build:16,build:17,env:production",
			Path:         "/apisix/admin/labels/all",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"build":"16"},{"build":"17"},{"env":"production"}`,
		})
	})
	It("Get labels (With condition, key = build && env = production; Paginate)", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			Query:        "label=build,env:production&page=3&page_size=1",
			Path:         "/apisix/admin/labels/all",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"env":"production"}`,
		})
	})
	It("Delete route", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("Delete consumer", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/consumers/c1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("Delete service", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/services/s1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("Delete upstream", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/u1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	It("Delete plugin config", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/plugin_configs/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
})
