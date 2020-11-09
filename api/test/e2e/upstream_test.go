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
package e2e

import (
	"net/http"
	"testing"
)

func TestUpstream_Create(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "use upstream that not exist",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/routes/r1",
			Body: `{
				"uri": "/hello_",
				"upstream_id": "not-exists"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusBadRequest,
		},
		{
			caseDesc: "create upstream",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/upstreams/1",
			Body: `{
                "nodes": [{
                    "host": "172.16.238.20",
                    "port": 1980,
                    "weight": 1
                }],
                "type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc: "create route using the upstream just created",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/routes/1",
			Body: `{
				"uri": "/hello_",
				"upstream_id": "1"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc:     "hit the route just created",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world\n",
			Sleep:        sleepTime,
		},
	}

	for _, tc := range tests {
		testCaseCheck(tc)
	}
}

func TestUpstream_Update(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "update upstream with domain",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/upstreams/1",
			Body: `{
               "nodes": [{
                   "host": "httpbin.org",
                   "port": 80,
                   "weight": 1
               }],
               "type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc:     "hit the route using upstream 1",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/hello_",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "/hello_\n",
			Sleep:        sleepTime,
		},
		{
			caseDesc: "update route using the upstream just created",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/routes/1",
			Body: `{
				"uri": "/get",
				"upstream_id": "1"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc:     "hit the route just created",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/get",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "/get\n",
			Sleep:        sleepTime,
		},
	}

	for _, tc := range tests {
		testCaseCheck(tc)
	}
}

//func TestRoute_Websocket(t *testing.T) {
//   tests := []HttpTestCase{
//       {
//           caseDesc:     "delete route",
//           Object:       MangerApiExpect(t),
//           Method:       http.MethodDelete,
//           Path:         "/apisix/admin/routes/r1",
//           Headers:      map[string]string{"Authorization": token},
//           ExpectStatus: http.StatusOK,
//       },
//       {
//           caseDesc:     "delete not exist route",
//           Object:       MangerApiExpect(t),
//           Method:       http.MethodDelete,
//           Path:         "/apisix/admin/routes/not-exist",
//           Headers:      map[string]string{"Authorization": token},
//           ExpectStatus: http.StatusNotFound,
//       },
//       {
//           caseDesc:     "hit the route just deleted",
//           Object:       APISIXExpect(t),
//           Method:       http.MethodGet,
//           Path:         "/hello1",
//           Headers:      map[string]string{"Host": "bar.com"},
//           ExpectStatus: http.StatusNotFound,
//           ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
//           Sleep:        sleepTime,
//       },
//   }
//
//   for _, tc := range tests {
//       testCaseCheck(tc)
//   }
//}

func TestRoute_Delete(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc:     "delete not exist upstream",
			Object:       MangerApiExpect(t),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/not-exist",
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusNotFound,
		},
		{
			caseDesc:     "delete upstream - being used by route 1",
			Object:       MangerApiExpect(t),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/1",
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusBadRequest,
		},
		{
			caseDesc:     "delete route",
			Object:       MangerApiExpect(t),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/1",
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc:     "delete upstream",
			Object:       MangerApiExpect(t),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/1",
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc:     "hit the route just deleted",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/hello1",
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   "{\"error_msg\":\"404 Route Not Found\"}\n",
			Sleep:        sleepTime,
		},
	}

	for _, tc := range tests {
		testCaseCheck(tc)
	}
}
