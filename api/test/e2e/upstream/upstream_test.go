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
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Upstream", func() {
	DescribeTable("test upstream create and update",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("create upstream1 by id", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/upstreams/1",
			Body: `{
					"nodes": {
						"` + base.UpstreamIp + `:1980": 1
					},
					"type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get upstream1", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/upstreams/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"key\":\"/apisix/upstreams/1", "\"nodes\":{\"127.0.0.1:1980\":1}"},
		}),
		Entry("create upstream2 by random id", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/upstreams",
			Body: `{
					"nodes": {
						"` + base.UpstreamIp + `:1981": 1
					},
					"type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get upstreams", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/upstreams",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"nodes\":{\"127.0.0.1:1981\":1}"},
		}),
		Entry("patch upstream1", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPatch,
			Path:   "/apisix/admin/upstreams/1",
			Body: `{
					"nodes": {
						"` + base.UpstreamIp + `:1982": 1
					}
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get upstream1", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/upstreams/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"nodes\":{\"127.0.0.1:1980\":1,\"127.0.0.1:1982\":1}"},
		}),
		Entry("subpath patch upstream1", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPatch,
			Path:   "/apisix/admin/upstreams/1/upstream/nodes",
			Body: `{
				"` + base.UpstreamIp + `:1983": 1
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get route1", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/upstreams/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"nodes\":{\"127.0.0.1:1983\":1}"},
		}),
	)
})
