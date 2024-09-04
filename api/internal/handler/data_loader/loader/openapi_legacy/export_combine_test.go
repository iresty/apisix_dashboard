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
package openapi_legacy

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/handler/data_loader/loader"
)

// Create route according to upstream ID and service ID
func TestExportRoutesCreateByUpstreamIdandServiceId(t *testing.T) {
	us := `{
		"id": "u1",
		"nodes": [
			{
				"host": "172.16.238.20",
				"port": 1980,
				"weight": 1
			}
		],
		"type": "roundrobin"
	}`

	s := `{
		"id": "s1",
		"name": "testservice",
		"desc": "testservice_desc",
		"labels": {
			"build":"16",
			"env":"production",
			"version":"v2"
		},
		"enable_websocket":true,
		"plugins": {
			"limit-count": {
				"count": 100,
				"time_window": 60,
				"rejected_code": 503,
				"key": "remote_addr"
			}
		},
		"upstream_id": "u1"
	}`

	r := `{
		"name": "route_all",
		"desc": "所有",
		"status": 1,
		"methods": ["GET"],
		"priority": 0,
		"service_id": "s1",
		"labels": {
			"test": "1",
			"API_VERSION": "v1"
		},
		"vars": [
			["arg_name", "==", "test"]
		],
		"uri": "/hello",
		"enable_websocket":false,
		"plugins": {
			"prometheus": {
				"disable": false
			}
		},
		"upstream_id": "u1"
	}`

	exportR := `{
		"components": {},
		"info": {
			"title": "RoutesExport",
			"version": "3.0.0"
		},
		"openapi": "3.0.0",
		"paths": {
			"/hello": {
				"get": {
					"operationId": "route_allGET",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"security": [],
					"summary": "所有",
					"x-apisix-enable_websocket": false,
					"x-apisix-labels": {
						"API_VERSION": "v1",
						"test": "1"
					},
					"x-apisix-plugins": {
						"limit-count": {
							"count": 100,
							"key": "remote_addr",
							"rejected_code": 503,
							"time_window": 60
						},
						"prometheus": {
							"disable": false
						}
					},
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"id": "u1",
						"nodes": [{
							"host": "172.16.238.20",
							"port": 1980,
							"weight": 1
						}],
						"type": "roundrobin"
					},
					"x-apisix-vars": [
						["arg_name", "==", "test"]
					]
				}
			}
		}
	}`
	var route *entity.Route
	var service *entity.Service
	var upstream *entity.Upstream
	err := json.Unmarshal([]byte(r), &route)
	err = json.Unmarshal([]byte(s), &service)
	err = json.Unmarshal([]byte(us), &upstream)

	l := &Loader{}
	exported, err := l.Export(loader.DataSets{
		Routes:    []entity.Route{*route},
		Services:  []entity.Service{*service},
		Upstreams: []entity.Upstream{*upstream},
	})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.Equal(t, replaceStr(exportR), string(ret))
	assert.NotNil(t, ret)
}

// Creating route using service ID does not contain upstream data
func TestExportRoutesCreateByServiceIdNoUpstream(t *testing.T) {
	us := `{
		"id": "u1",
		"nodes": [
			{
				"host": "172.16.238.20",
				"port": 1980,
				"weight": 1
			}
		],
		"type": "roundrobin"
	}`

	s := `{
		"id": "s5",
		"name": "testservice",
		"desc": "testservice_desc",
		"enable_websocket":true,
		"upstream_id": "u1"
	}`

	r := `{
		"name": "route_all",
		"desc": "所有",
		"status": 1,
		"methods": ["GET"],
		"priority": 0,
		"service_id": "s5",
		"labels": {
			"test": "1",
			"API_VERSION": "v1"
		},
		"vars": [
			["arg_name", "==", "test"]
		],
		"uri": "/hello",
		"enable_websocket":false,
		"plugins": {
			"prometheus": {
				"disable": false
			}
		}
	}`

	exportR := `{
		"components": {},
		"info": {
			"title": "RoutesExport",
			"version": "3.0.0"
		},
		"openapi": "3.0.0",
		"paths": {
			"/hello": {
				"get": {
					"operationId": "route_allGET",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"security": [],
					"summary": "所有",
					"x-apisix-enable_websocket": false,
					"x-apisix-labels": {
						"API_VERSION": "v1",
						"test": "1"
					},
					"x-apisix-plugins": {
						"prometheus": {
							"disable": false
						}
					},
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"id": "u1",
						"nodes": [{
							"host": "172.16.238.20",
							"port": 1980,
							"weight": 1
						}],
						"type": "roundrobin"
					},
					"x-apisix-vars": [
						["arg_name", "==", "test"]
					]
				}
			}
		}
	}`
	var route *entity.Route
	var service *entity.Service
	var upstream *entity.Upstream
	err := json.Unmarshal([]byte(r), &route)
	err = json.Unmarshal([]byte(s), &service)
	err = json.Unmarshal([]byte(us), &upstream)

	l := &Loader{}
	exported, err := l.Export(loader.DataSets{
		Routes:    []entity.Route{*route},
		Services:  []entity.Service{*service},
		Upstreams: []entity.Upstream{*upstream},
	})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.Equal(t, replaceStr(exportR), string(ret))
	assert.NotNil(t, ret)
}

// Create service according to upstream1 ID
// Create route according to upstream2 ID and service ID
func TestExportRoutesCreateByUpstreamIDAndServiceID2(t *testing.T) {
	us := `{
		"id": "u1",
		"nodes": [
			{
				"host": "172.16.238.20",
				"port": 1980,
				"weight": 1
			}
		],
		"type": "roundrobin"
	}`

	us2 := `{
		"id": "u2",
		"nodes": [
			{
				"host": "172.16.238.20",
				"port": 1981,
				"weight": 1
			}
		],
		"type": "roundrobin"
	}`

	s := `{
		"id": "s1",
		"name": "testservice",
		"desc": "testservice_desc",
		"enable_websocket":true,
		"upstream_id": "u2"
	}`

	r := `{
		"name": "route_all",
		"desc": "所有",
		"status": 1,
		"methods": ["GET"],
		"priority": 0,
		"service_id": "s1",
		"labels": {
			"test": "1",
			"API_VERSION": "v1"
		},
		"vars": [
			["arg_name", "==", "test"]
		],
		"uri": "/hello",
		"enable_websocket":false,
		"plugins": {
			"prometheus": {
				"disable": false
			}
		},
		"upstream_id": "u1"
	}`

	exportR := `{
		"components": {},
		"info": {
			"title": "RoutesExport",
			"version": "3.0.0"
		},
		"openapi": "3.0.0",
		"paths": {
			"/hello": {
				"get": {
					"operationId": "route_allGET",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"security": [],
					"summary": "所有",
					"x-apisix-enable_websocket": false,
					"x-apisix-labels": {
						"API_VERSION": "v1",
						"test": "1"
					},
					"x-apisix-plugins": {
						"prometheus": {
							"disable": false
						}
					},
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"id": "u1",
						"nodes": [{
							"host": "172.16.238.20",
							"port": 1980,
							"weight": 1
						}],
						"type": "roundrobin"
					},
					"x-apisix-vars": [
						["arg_name", "==", "test"]
					]
				}
			}
		}
	}`
	var route *entity.Route
	var service *entity.Service
	var upstream *entity.Upstream
	var upstream2 *entity.Upstream

	err := json.Unmarshal([]byte(r), &route)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(s), &service)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(us), &upstream)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(us2), &upstream2)
	assert.NoError(t, err)

	l := &Loader{}
	exported, err := l.Export(loader.DataSets{
		Routes:    []entity.Route{*route},
		Services:  []entity.Service{*service},
		Upstreams: []entity.Upstream{*upstream, *upstream2},
	})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.Equal(t, replaceStr(exportR), string(ret))
	assert.NotNil(t, ret)
}
