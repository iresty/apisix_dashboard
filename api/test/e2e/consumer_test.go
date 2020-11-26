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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestConsumer_without_username(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "create consumer without username",
			Object:   MangerApiExpect(t),
			Path:     "/apisix/admin/consumers",
			Method:   http.MethodPut,
			Body: `{
				 "plugins": {
					 "key-auth": {
						 "key": "auth-new"
					 }
				 },
				 "desc": "test description"
			 }`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   "scheme validate fail",
		},
	}

	for _, tc := range tests {
		testCaseCheck(tc)
	}
}

func TestConsumer_delete_notexit_consumer(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc:     "delete notexit consumer",
			Object:       MangerApiExpect(t),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/consumers/notexit",
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		testCaseCheck(tc)
	}
}

func TestConsumer_with_error_key(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "create consumer with error key",
			Object:   MangerApiExpect(t),
			Path:     "/apisix/admin/consumers",
			Method:   http.MethodPut,
			Body: `{
				 "username": "jack",
				 "plugins": {
					 "key-authaa": {
						 "key": "auth-one"
					 }
				 },
				 "desc": "test description"
			 }`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusBadRequest,
			ExpectBody:   "scheme validate failed",
		},
		{
			caseDesc:     "verify consumer",
			Object:       MangerApiExpect(t),
			Path:         "/apisix/admin/consumers/jack",
			Method:       http.MethodGet,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusNotFound,
			Sleep:        sleepTime,
		},
	}

	for _, tc := range tests {
		testCaseCheck(tc)
	}
}

func TestConsumer_add_consumer_with_labels(t *testing.T) {
	tests := []HttpTestCase{
		{
			caseDesc: "create consumer",
			Object:   MangerApiExpect(t),
			Path:     "/apisix/admin/consumers",
			Method:   http.MethodPut,
			Body: `{
				"username": "jack",
				"labels": {
					"build":"16",
					"env":"production",
					"version":"v2"
				},
				"plugins": {
					"key-auth": {
						"key": "auth-two"
					}
				},
			    "desc": "test description"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc:     "verify consumer",
			Object:       MangerApiExpect(t),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "\"username\":\"jack\",\"desc\":\"test description\",\"plugins\":{\"key-auth\":{\"key\":\"auth-two\"}},\"labels\":{\"build\":\"16\",\"env\":\"production\",\"version\":\"v2\"}",
			Sleep:        sleepTime,
		},
		{
			caseDesc: "create route",
			Object:   MangerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/routes/r1",
			Body: `{
				"uri": "/hello",
				"plugins": {
					"key-auth": {}
				},
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
			caseDesc:     "verify route",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"apikey": "auth-two"},
			ExpectStatus: http.StatusOK,
			Sleep:        sleepTime,
		},
		{
			caseDesc:     "delete consumer",
			Object:       MangerApiExpect(t),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/consumers/jack",
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc:     "delete route",
			Object:       MangerApiExpect(t),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		testCaseCheck(tc)
	}
}

func TestConsumer_with_createtime_updatetime(t *testing.T) {
	//create consumer
	basepath := "http://127.0.0.1:8080/apisix/admin/consumers"
	data := `{
		"username":"jack",
		"desc": "new consumer"
    }`
	request, _ := http.NewRequest("PUT", basepath, strings.NewReader(data))
	request.Header.Add("Authorization", token)
	resp, _ := http.DefaultClient.Do(request)
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	code := gjson.Get(string(respBody), "code")
	assert.Equal(t, code.String(), "0")

	time.Sleep(1 * time.Second)

	request, _ = http.NewRequest("GET", basepath+"/jack", nil)
	request.Header.Add("Authorization", token)
	resp, _ = http.DefaultClient.Do(request)
	defer resp.Body.Close()
	respBody, _ = ioutil.ReadAll(resp.Body)
	createtime := gjson.Get(string(respBody), "data.create_time")
	updatetime := gjson.Get(string(respBody), "data.update_time")

	//create consumer again, compare create_time and update_time
	data = `{
		"username":"jack",
		"desc": "new consumer haha"
    }`
	request, _ = http.NewRequest("PUT", basepath, strings.NewReader(data))
	request.Header.Add("Authorization", token)
	resp, _ = http.DefaultClient.Do(request)
	defer resp.Body.Close()
	respBody, _ = ioutil.ReadAll(resp.Body)
	code = gjson.Get(string(respBody), "code")
	assert.Equal(t, code.String(), "0")

	time.Sleep(1 * time.Second)

	request, _ = http.NewRequest("GET", basepath+"/jack", nil)
	request.Header.Add("Authorization", token)
	resp, _ = http.DefaultClient.Do(request)
	defer resp.Body.Close()
	respBody, _ = ioutil.ReadAll(resp.Body)
	createtime2 := gjson.Get(string(respBody), "data.create_time")
	updatetime2 := gjson.Get(string(respBody), "data.update_time")

	assert.Equal(t, createtime.String(), createtime2.String())
	assert.NotEqual(t, updatetime.String(), updatetime2.String())

	//deletea consumer
	request, _ = http.NewRequest("DELETE", basepath+"/jack", nil)
	request.Header.Add("Authorization", token)
	_, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)
}
