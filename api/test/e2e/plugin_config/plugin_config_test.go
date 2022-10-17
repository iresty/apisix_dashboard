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
package plugin_config_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Plugin Config", func() {
	DescribeTable("test plugin config create and update",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("create plugin config 1 by id", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/plugin_configs/1",
			Body: `{
				"desc": "enable limit-count plugin",
				"plugins": {
					"limit-count": {
						"count": 2,
						"time_window": 60,
						"rejected_code": 503
					}
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get plugin config 1", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/plugin_configs/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"key\":\"/apisix/plugin_configs/1\"", "\"plugins\":{\"limit-count\"", "\"count\":2"},
		}),
		Entry("get plugin configs", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/plugin_configs",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"key\":\"/apisix/plugin_configs/1\"", "\"plugins\":{\"limit-count\"", "\"count\":2"}}),
		Entry("patch plugin config 1", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPatch,
			Path:   "/apisix/admin/plugin_configs/1",
			Body: `{
				"plugins": {
					"limit-count": {
						"count": 3,
						"time_window": 60,
						"rejected_code": 503
					}
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get plugin config 1", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/plugin_configs/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"key\":\"/apisix/plugin_configs/1\"", "\"plugins\":{\"limit-count\"", "\"count\":3"},
		}),
		Entry("subpath patch plugin config 1", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPatch,
			Path:   "/apisix/admin/plugin_configs/1/plugins",
			Body: `{
				"limit-count": {
					"count": 4,
					"time_window": 60,
					"rejected_code": 503
				}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get plugin config 1", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/plugin_configs/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"key\":\"/apisix/plugin_configs/1\"", "\"plugins\":{\"limit-count\"", "\"count\":4"},
		}),
	)
})
