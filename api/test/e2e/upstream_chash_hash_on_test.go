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
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpstream_chash_hash_on_custom_header(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "create chash upstream with hash_on (custom_header)",
			Object:   ManagerApiExpect(t),
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
					}],
					"type": "chash",
					"key": "custom_header",
					"hash_on": "header"
				}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc: "create route using the upstream just created",
			Object:   ManagerApiExpect(t),
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
	time.Sleep(time.Duration(100) * time.Millisecond)
	basepath := "http://127.0.0.1:9080"
	var req *http.Request
	var err error
	var url string
	var resp *http.Response
	var respBody []byte
	res := map[string]int{}
	for i := 0; i <= 3; i++ {
		url = basepath + "/server_port?var=2&var2=" + strconv.Itoa(i)
		req, err = http.NewRequest("GET", url, nil)
		req.Header.Add("custom_header", `custom-one`)
		resp, err = http.DefaultClient.Do(req)
		assert.Nil(t, err)
		respBody, err = ioutil.ReadAll(resp.Body)
		body := string(respBody)
		if _, ok := res[body]; !ok {
			res[body] = 1
		} else {
			res[body] += 1
		}
	}
	assert.Equal(t, true, res["1980"] == 4 || res["1981"] == 4)
	resp.Body.Close()
}

func TestUpstream_chash_hash_on_cookie(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "create cHash upstream with hash_on (cookie)",
			Object:   ManagerApiExpect(t),
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
					}],
					"type": "chash",
					"key": "custom-cookie",
					"hash_on": "cookie"
				}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc: "create route using the upstream just created",
			Object:   ManagerApiExpect(t),
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
	time.Sleep(time.Duration(100) * time.Millisecond)
	basepath := "http://127.0.0.1:9080"
	var req *http.Request
	var err error
	var url string
	var resp *http.Response
	var respBody []byte
	res := map[string]int{}
	for i := 0; i <= 3; i++ {
		url = basepath + "/server_port"
		req, err = http.NewRequest("GET", url, nil)
		req.Header.Add("Cookie", `custom-cookie=cuscookie`)
		resp, err = http.DefaultClient.Do(req)
		assert.Nil(t, err)
		respBody, err = ioutil.ReadAll(resp.Body)
		body := string(respBody)
		if _, ok := res[body]; !ok {
			res[body] = 1
		} else {
			res[body] += 1
		}
	}
	assert.Equal(t, true, res["1980"] == 4 || res["1981"] == 4)
	resp.Body.Close()

	//hit routes with miss cookie
	res = map[string]int{}
	for i := 0; i <= 3; i++ {
		url = basepath + "/server_port"
		req, err = http.NewRequest("GET", url, nil)
		req.Header.Add("Cookie", `miss-custom-cookie=cuscookie`)
		resp, err = http.DefaultClient.Do(req)
		assert.Nil(t, err)
		respBody, err = ioutil.ReadAll(resp.Body)
		body := string(respBody)
		if _, ok := res[body]; !ok {
			res[body] = 1
		} else {
			res[body] += 1
		}
	}
	assert.Equal(t, true, res["1980"] == 4 || res["1981"] == 4)
	resp.Body.Close()
}

func TestUpstream_key_contains_uppercase_letters_and_hyphen(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "create cHash upstream with key contains uppercase letters and hyphen",
			Object:   ManagerApiExpect(t),
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
					}],
					"type": "chash",
					"key": "X-Sessionid",
					"hash_on": "header"
				}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc: "create route using the upstream just created",
			Object:   ManagerApiExpect(t),
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
	time.Sleep(time.Duration(100) * time.Millisecond)
	basepath := "http://127.0.0.1:9080"
	var req *http.Request
	var err error
	var url string
	var resp *http.Response
	var respBody []byte
	res := map[string]int{}
	for i := 0; i <= 15; i++ {
		url = basepath + "/server_port"
		req, err = http.NewRequest("GET", url, nil)
		req.Header.Add("X-Sessionid", `chash_val_`+strconv.Itoa(i))
		resp, err = http.DefaultClient.Do(req)
		assert.Nil(t, err)
		respBody, err = ioutil.ReadAll(resp.Body)
		body := string(respBody)
		if _, ok := res[body]; !ok {
			res[body] = 1
		} else {
			res[body] += 1
		}
	}
	assert.Equal(t, true, res["1980"] == 8 && res["1981"] == 8)
	resp.Body.Close()
}

func TestUpstream_chash_hash_on_consumer(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "create consumer with key-auth",
			Object:   ManagerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/consumers",
			Body: `{
					"username": "jack",
					"plugins": {
						"key-auth": {
							"key": "auth-jack"
						}
					}
				}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc: "create route with key-auth",
			Object:   ManagerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/routes/1",
			Body: `{
					"uri": "/server_port",
					"plugins": {
						"key-auth": {}
					},
					"upstream": {
						"nodes": [{
							"host": "172.16.238.20",
							"port": 1980,
							"weight": 1
						},
						{
							"host": "172.16.238.20",
							"port": 1981,
							"weight": 1
						}],
						"type": "chash",
						"hash_on": "consumer"
					}
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
	time.Sleep(time.Duration(100) * time.Millisecond)
	basepath := "http://127.0.0.1:9080"
	var req *http.Request
	var err error
	var url string
	var resp *http.Response
	var respBody []byte
	res := map[string]int{}
	for i := 0; i <= 3; i++ {
		url = basepath + "/server_port"
		req, err = http.NewRequest("GET", url, nil)
		req.Header.Add("apikey", `auth-jack`)
		resp, err = http.DefaultClient.Do(req)
		assert.Nil(t, err)
		respBody, err = ioutil.ReadAll(resp.Body)
		body := string(respBody)
		if _, ok := res[body]; !ok {
			res[body] = 1
		} else {
			res[body] += 1
		}
	}
	assert.Equal(t, true, res["1980"] == 4 || res["1981"] == 4)
	resp.Body.Close()
}

func TestUpstream_Delete_hash_on(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc:     "delete consumer",
			Object:       ManagerApiExpect(t),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc:     "delete route",
			Object:       ManagerApiExpect(t),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/1",
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc:     "delete upstream",
			Object:       ManagerApiExpect(t),
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
