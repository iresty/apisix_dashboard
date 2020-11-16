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

func TestRoute_Invalid_Service_And_Service(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "use service that not exist",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/routes/r1",
			Body: `{
				"uri": "/hello_",
				"service_id": "not-exists"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusBadRequest,
		},
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
			caseDesc: "create service and upstream together at the same time",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/routes/r1",
			Body: `{
				"uri": "/hello_",
				"service_id": "not-exists-service",
				"upstream_id": "not-exists-upstream",
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusBadRequest,
		},
	}
	for _, tc := range tests {
		testCaseCheck(tc)
	}
}

func TestRoute_Create_Upstream(t *testing.T) {
	tests := []HttpTestCase{
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
			Path:     "/apisix/admin/routes/r1",
			Body: `{
				"uri": "/server_port",
				"upstream_id": "1"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
			Sleep:        sleepTime,
		},
		{
			caseDesc:     "hit the route just created",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/server_port",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "1980",
			Sleep:        sleepTime,
		},
	}
	for _, tc := range tests {
		testCaseCheck(tc)
	}
}

func TestRoute_Create_Service(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "create service",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/services/200",
			Body: `{
    			"upstream": {
        		"type": "roundrobin",
                "nodes": [{
                    "host": "172.16.238.20",
                    "port": 1980,
                    "weight": 1
                	}]
                }
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc: "create route using the service just created",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/routes/r2",
			Body: `{
				"uri": "/hello",
				"service_id": "200"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
			Sleep:        sleepTime,
		},
		{
			caseDesc:     "hit the route just created",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world\n",
			Sleep:        sleepTime,
		},
	}
	for _, tc := range tests {
		testCaseCheck(tc)
	}
}

//func TestRoute_Delete_Upstream(t *testing.T) {
//	tests := []HttpTestCase{
//		{
//			caseDesc:     "remove upstream before deleting route",
//			Object:       MangerApiExpect(t),
//			Method:       http.MethodDelete,
//			Path:         "/apisix/admin/upstreams/1",
//			Headers:      map[string]string{"Authorization": token},
//			ExpectStatus: http.StatusBadRequest,
//		},
//		{
//			caseDesc:     "delete route",
//			Object:       MangerApiExpect(t),
//			Method:       http.MethodDelete,
//			Path:         "/apisix/admin/routes/r1",
//			Headers:      map[string]string{"Authorization": token},
//			ExpectStatus: http.StatusOK,
//		},
//		{
//			caseDesc:     "remove upstream",
//			Object:       MangerApiExpect(t),
//			Method:       http.MethodDelete,
//			Path:         "/apisix/admin/upstreams/1",
//			Headers:      map[string]string{"Authorization": token},
//			ExpectStatus: http.StatusBadRequest,
//		},
//	}
//	for _, tc := range tests {
//		testCaseCheck(tc)
//	}
//}
