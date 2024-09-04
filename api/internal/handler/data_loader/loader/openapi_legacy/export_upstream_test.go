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

// Create and export route according to upstream ID
func TestExportRoutesCreateByUpstreamId(t *testing.T) {
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

	r := `{
		"methods": ["GET"],
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
					"operationId": "GET",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"security": [],
					"x-apisix-enable_websocket": false,
					"x-apisix-plugins": {
						"prometheus": {
							"disable": false
						}
					},
					"x-apisix-priority": 0,
					"x-apisix-status": 0,
					"x-apisix-upstream": {
						"id": "u1",
						"nodes": [{
							"host": "172.16.238.20",
							"port": 1980,
							"weight": 1
						}],
						"type": "roundrobin"
					}
				}
			}
		}
	}`
	var route *entity.Route
	var upstream *entity.Upstream
	err := json.Unmarshal([]byte(r), &route)
	err = json.Unmarshal([]byte(us), &upstream)

	l := &Loader{}
	exported, err := l.Export(loader.DataSets{
		Routes:    []entity.Route{*route},
		Upstreams: []entity.Upstream{*upstream},
	})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.Equal(t, replaceStr(exportR), string(ret))
	assert.NotNil(t, ret)
}
