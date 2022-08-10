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
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("route with valid remote_addr remote_addrs", func() {
	DescribeTable("test route with valid remote_addr remote_addrs",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("add route with valid remote_addr", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				 "name": "route1",
				 "uri": "/hello",
				 "remote_addr": "172.16.238.1",
				 "upstream": {
					 "type": "roundrobin",
					 "nodes": {
						 "` + base.UpstreamIp + `:1980": 1
					 }
				 }
			 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("verify route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		}),
		Entry("update route with valid remote_addr (CIDR)", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				 "name": "route1",
				 "uri": "/hello",
				 "remote_addr": "172.16.238.1/24",
				 "upstream": {
					 "type": "roundrobin",
					 "nodes": {
						 "` + base.UpstreamIp + `:1980": 1
					 }
				 }
			 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("verify route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		}),
		Entry("update route with valid remote_addrs", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				 "name": "route1",
				 "uri": "/hello",
				 "remote_addrs": ["172.16.238.1","192.168.0.2/24"],
				 "upstream": {
					 "type": "roundrobin",
					 "nodes": {
						 "` + base.UpstreamIp + `:1980": 1
					 }
				 }
			 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("verify route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        base.SleepTime,
		}),
		Entry("update remote_addr to not be hit", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				 "name": "route1",
				 "uri": "/hello",
				 "remote_addr": "10.10.10.10",
				 "upstream": {
					 "type": "roundrobin",
					 "nodes": {
						 "` + base.UpstreamIp + `:1980": 1
					 }
				 }
			 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("verify route not found", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		}),
		Entry("update remote_addrs to not be hit", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				 "name": "route1",
				 "uri": "/hello",
				 "remote_addrs": ["10.10.10.10","11.11.11.1/24"],
				 "upstream": {
					 "type": "roundrobin",
					 "nodes": {
						 "` + base.UpstreamIp + `:1980": 1
					 }
				 }
			 }`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("verify route not found", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		}),
		Entry("delete route", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		}),
		Entry("verify it again after deleting the route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		}),
	)
})

var _ = Describe("route with invalid remote_addr", func() {
	DescribeTable("route with remote_addr",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("config route with invalid remote_addr", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				 "name": "route1",
				 "uri": "/hello",
				 "remote_addr": "127.0.0.",
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
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   "{\"code\":10000,\"message\":\"schema validate failed: remote_addr: Must validate at least one schema (anyOf)\\nremote_addr: Does not match format 'ipv4'\"}",
		}),
		Entry("verify route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
			Sleep:        base.SleepTime,
		}),
		Entry("config route with invalid remote_addr", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				 "name": "route1",
				 "uri": "/hello",
				 "remote_addr": "127.0.0.aa",
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
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   "{\"code\":10000,\"message\":\"schema validate failed: remote_addr: Must validate at least one schema (anyOf)\\nremote_addr: Does not match format 'ipv4'\"}",
		}),
		Entry("verify route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
			Sleep:        base.SleepTime,
		}),
		Entry("config route with invalid remote_addrs", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				 "name": "route1",
				 "uri": "/hello",
				 "remote_addrs": ["127.0.0.1","192.168.0."],
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
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   "{\"code\":10000,\"message\":\"schema validate failed: remote_addrs.1: Must validate at least one schema (anyOf)\\nremote_addrs.1: Does not match format 'ipv4'\"}",
		}),
		Entry("verify route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
			Sleep:        base.SleepTime,
		}),
	)
})
