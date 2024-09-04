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

// Create a service with label data and a route with label data, and export the route.
// Label is the original data of the route
func TestExportRoutesCreateByLabel(t *testing.T) {
	s := `{
		"name": "testservice",
		"desc": "testservice_desc",
		"enable_websocket":true,
		"labels": {
			"build": "10"
		}
	}`

	r := `{
		"name": "route_all",
		"desc": "所有",
		"status": 1,
		"methods": ["GET"],
		"service_id": "s1",
		"labels": {
			"test": "1",
			"API_VERSION": "v1"
		},
		"uri": "/hello",
		"enable_websocket":false,
		"upstream": {
			"nodes": {
				"172.16.238.20:1980": 1
			},
			"type": "roundrobin"
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
					"summary": "所有",
					"x-apisix-enable_websocket": false,
					"x-apisix-labels": {
						"API_VERSION": "v1",
						"test": "1"
					},
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				}
			}
		}
	}`
	var route *entity.Route
	var service *entity.Service
	err := json.Unmarshal([]byte(r), &route)
	err = json.Unmarshal([]byte(s), &service)

	l := &Loader{}
	exported, err := l.Export(loader.DataSets{
		Routes:   []entity.Route{*route},
		Services: []entity.Service{*service},
	})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.Equal(t, replaceStr(exportR), string(ret))
	assert.NotNil(t, ret)
}

// Create a service with label data and a route without label data, and export the route.
// Label is the data of the service
func TestExportRoutesCreateByLabel2(t *testing.T) {
	s := `{
		"id": "s2",
		"name": "testservice",
		"desc": "testservice_desc",
		"enable_websocket":true,
		"labels": {
			"build": "16",
			"env": "production",
			"version": "v2"
		}
	}`

	r := `{
		"name": "route_all",
		"desc": "所有",
		"status": 1,
		"methods": ["GET"],
		"service_id": "s2",
		"vars": [
			["arg_name", "==", "test"]
		],
		"uri": "/hello"
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
					"summary": "所有",
					"x-apisix-enable_websocket": false,
					"x-apisix-labels": {
						"build": "16",
						"env": "production",
						"version": "v2"
					},
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-vars": [
						["arg_name", "==", "test"]
					]
				}
			}
		}
	}`
	var route *entity.Route
	var service *entity.Service
	err := json.Unmarshal([]byte(r), &route)
	err = json.Unmarshal([]byte(s), &service)

	l := &Loader{}
	exported, err := l.Export(loader.DataSets{
		Routes:   []entity.Route{*route},
		Services: []entity.Service{*service},
	})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.Equal(t, replaceStr(exportR), string(ret))
	assert.NotNil(t, ret)
}
