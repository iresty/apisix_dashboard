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
	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

func TestUpstream_Create(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "use upstream that not exist",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/routes/r1",
			Body: `{
				"uri": "/hello",
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
				"uri": "/hello",
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
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
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
					"host": "172.16.238.20",
					"port": 1981,
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
			Path:         "/hello",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world",
			Sleep:        sleepTime,
		},
	}

	for _, tc := range tests {
		testCaseCheck(tc)
	}
}

func TestRoute_Node_Host(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "update upstream - pass host: node",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/upstreams/1",
			Body: `{
				"nodes": [{
					"host": "httpbin.org",
					"port": 80,
					"weight": 1
				}],
				"type": "roundrobin",
				"pass_host": "node"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc: "update path for route",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/routes/1",
			Body: `{
				"uri": "/*",
				"upstream_id": "1"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc:     "hit the route ",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/get",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "\"Host\": \"httpbin.org\"",
			Sleep:        sleepTime,
		},
		{
			caseDesc: "update upstream - pass host: rewrite",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/upstreams/1",
			Body: `{
				"nodes": [{
					"host": "172.16.238.20",
					"port": 1980,
					"weight": 1
				}],
				"type": "roundrobin",
				"pass_host": "rewrite",
				"upstream_host": "httpbin.org"  
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc:     "hit the route ",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/uri",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "x-forwarded-host: 127.0.0.1",
			Sleep:        sleepTime,
		},
	}

	for _, tc := range tests {
		testCaseCheck(tc)
	}
}

func TestUpstream_cHash(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "create cHash upstream with key (remote_addr)",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/upstreams/1",
			Body: `{
                "nodes": [{
                    "host": "172.16.238.20",
                    "port": 1980,
                    "weight": 1
				},
				{
                    "host": "172.16.238.20",
                    "port": 1981,
                    "weight": 1
				},
				{
                    "host": "172.16.238.20",
                    "port": 1982,
                    "weight": 1
                }],
				"type": "chash",
				"hash_on":"header",
				"key": "remote_addr"
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
				"uri": "/server_port",
				"upstream_id": "1"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
			Sleep:        sleepTime,
		},
		
	}

	for _, tc := range tests {
		testCaseCheck(tc)
	}

	//hit routes
	basepath := "http://127.0.0.1:9080/"
	request, err := http.NewRequest("GET", basepath+"/server_port", nil)
	request.Header.Add("Authorization", token)
	var resp *http.Response
	var respBody []byte
	var count int
	for i := 0; i <= 17; i++ {
		resp, err = http.DefaultClient.Do(request)
		assert.Nil(t, err)
		respBody, err = ioutil.ReadAll(resp.Body)
		if string(respBody) == "1982"{
			count ++
		}
	}
	assert.Equal(t, count, 18)
	defer resp.Body.Close()

}

func TestUpstream_Delete(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc:     "delete not exist upstream",
			Object:       MangerApiExpect(t),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/upstreams/not-exist",
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusNotFound,
		},
		// TODO delete upstream - being used by route 1
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
