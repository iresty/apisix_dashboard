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

// Export data as the route of URIs Hosts
func TestExportRoutes1(t *testing.T) {
	//*entity.Route
	r1 := `{
		"name": "aaaa",
		"labels": {
			"build":"16",
			"env":"production",
			"version":"v2"
		},
		"plugins": {
			"limit-count": {
				"count": 2,
				"time_window": 60,
				"rejected_code": 503,
				"key": "remote_addr"
			}
		},
		"status": 1,
		"uris": ["/hello_"],
		"hosts": ["foo.com", "*.bar.com"],
		"methods": ["GET", "POST"],
		"upstream": {
			"nodes": {
				"172.16.238.20:1980": 1
			},
			"type": "roundrobin"
		}
	}`

	exportR1 := `{
		"components": {},
		"info": {
			"title": "RoutesExport",
			"version": "3.0.0"
		},
		"openapi": "3.0.0",
		"paths": {
			"/hello_": {
				"get": {
					"operationId": "aaaaGET",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"security": [],
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["foo.com", "*.bar.com"],
					"x-apisix-labels": {
						"build": "16",
						"env": "production",
						"version": "v2"
					},
					"x-apisix-plugins": {
						"limit-count": {
							"count": 2,
							"key": "remote_addr",
							"rejected_code": 503,
							"time_window": 60
						}
					},
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"post": {
					"operationId": "aaaaPOST",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"security": [],
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["foo.com", "*.bar.com"],
					"x-apisix-labels": {
						"build": "16",
						"env": "production",
						"version": "v2"
					},
					"x-apisix-plugins": {
						"limit-count": {
							"count": 2,
							"key": "remote_addr",
							"rejected_code": 503,
							"time_window": 60
						}
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
	err := json.Unmarshal([]byte(r1), &route)

	l := &Loader{}
	exported, err := l.Export(loader.DataSets{Routes: []entity.Route{*route}})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.Equal(t, replaceStr(exportR1), string(ret))
	assert.NotNil(t, ret)
}

// Export data as the route of URI host
func TestExportRoutes2(t *testing.T) {
	//*entity.Route
	r2 := `{
		"name": "aaaa2",
		"labels": {
			"build":"16",
			"env":"production",
			"version":"v2"
		},
		"plugins": {
			"limit-count": {
				"count": 2,
				"time_window": 60,
				"rejected_code": 503,
				"key": "remote_addr"
			}
		},
		"status": 1,
		"uri": "/hello2",
		"host": "*.bar.com",
		"methods": ["GET", "POST"],
		"upstream": {
			"nodes": {
				"172.16.238.20:1980": 1
			},
			"type": "roundrobin"
		}
	}`

	exportR2 := `{
		"components": {},
		"info": {
			"title": "RoutesExport",
			"version": "3.0.0"
		},
		"openapi": "3.0.0",
		"paths": {
			"/hello2": {
				"get": {
					"operationId": "aaaa2GET",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"security": [],
					"x-apisix-enable_websocket": false,
					"x-apisix-host": "*.bar.com",
					"x-apisix-labels": {
						"build": "16",
						"env": "production",
						"version": "v2"
					},
					"x-apisix-plugins": {
						"limit-count": {
							"count": 2,
							"key": "remote_addr",
							"rejected_code": 503,
							"time_window": 60
						}
					},
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"post": {
					"operationId": "aaaa2POST",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"security": [],
					"x-apisix-enable_websocket": false,
					"x-apisix-host": "*.bar.com",
					"x-apisix-labels": {
						"build": "16",
						"env": "production",
						"version": "v2"
					},
					"x-apisix-plugins": {
						"limit-count": {
							"count": 2,
							"key": "remote_addr",
							"rejected_code": 503,
							"time_window": 60
						}
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
	err := json.Unmarshal([]byte(r2), &route)

	l := &Loader{}
	exported, err := l.Export(loader.DataSets{Routes: []entity.Route{*route}})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.Equal(t, replaceStr(exportR2), string(ret))
	assert.NotNil(t, ret)
}
