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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/handler/data_loader/loader"
)

// Test export route request_ validation data correctness
func TestExportRoutesCreateByRequestValidation(t *testing.T) {
	r := `{
		"uris": ["/test-test"],
		"name": "route_all",
		"desc": "所有",
		"methods": ["GET"],
		"hosts": ["test.com"],
		"plugins": {
			"request-validation": {
				"body_schema": {
					"properties": {
						"boolean_payload": {
							"type": "boolean"
						},
						"required_payload": {
							"type": "string"
						}
					},
					"required": ["required_payload"],
					"type": "object"
				},
				"disable": false,
				"header_schema": {
					"properties": {
						"test": {
							"enum": "test-enum",
							"type": "string"
						}
					},
					"type": "string"
				}
			}
		},
		"status": 1
	}`

	exportR := `{
		"components": {},
		"info": {
			"title": "RoutesExport",
			"version": "3.0.0"
		},
		"openapi": "3.0.0",
		"paths": {
			"/test-test": {
				"get": {
					"operationId": "route_allGET",
					"parameters": [{
						"in": "header",
						"name": "test",
						"schema": {
							"type": "string"
						}
					}],
					"requestBody": {
						"content": {
							"*/*": {
								"schema": {
									"properties": {
										"boolean_payload": {
											"type": "boolean"
										},
										"required_payload": {
											"type": "string"
										}
									},
									"required": ["required_payload"],
									"type": "object"
								}
							}
						}
					},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"security": [],
					"summary": "所有",
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1
				}
			}
		}
	}`
	var route *entity.Route
	err := json.Unmarshal([]byte(r), &route)

	l := &Loader{}
	exported, err := l.Export(loader.DataSets{
		Routes: []entity.Route{*route},
	})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.Equal(t, replaceStr(exportR), string(ret))
	assert.NotNil(t, ret)
}

// Export route create by jwt-auth plugin
func TestExportRoutesCreateByJWTAuth(t *testing.T) {
	r := `{
		"uri": "/hello",
		"methods": ["Get"],
		"plugins": {
			"jwt-auth": {}
		},
		"upstream": {
			"type": "roundrobin",
		   "nodes": [{
			   "host": "172.16.238.20",
			   "port": 1980,
			   "weight": 1
		   }]
		}
	}`

	c := `{
		"username": "jack",
		"plugins": {
			"jwt-auth": {
				"key": "user-key",
				"secret": "my-secret-key",
				"algorithm": "HS256"
			}
		},
		"desc": "test description"
	}`

	exportR := `{
		"components": {
			"securitySchemes": {
				"bearerAuth": {
					"bearerFormat": "JWT",
					"scheme": "bearer",
					"type": "http"
				}
			}
		},
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
					"security": [{
						"bearerAuth": [""]
					}],
					"x-apisix-enable_websocket": false,
					"x-apisix-priority": 0,
					"x-apisix-status": 0,
					"x-apisix-upstream": {
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
	var consumer *entity.Consumer
	err := json.Unmarshal([]byte(r), &route)
	err = json.Unmarshal([]byte(c), &consumer)

	l := &Loader{}
	exported, err := l.Export(loader.DataSets{
		Routes:    []entity.Route{*route},
		Consumers: []entity.Consumer{*consumer},
	})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.Equal(t, replaceStr(exportR), strings.Replace(string(ret), " ", "", -1))
	assert.NotNil(t, ret)
}

// 13.Export route create by apikey-auth plugin  basic-auth plugin
func TestExportRoutesCreateByKeyAuthAndBasicAuth(t *testing.T) {
	r := `{
		"uri": "/hello",
		"methods": ["Get"],
		"plugins": {
			"key-auth": {},
			"basic-auth": {}
		},
		"upstream": {
			"type": "roundrobin",
		   "nodes": [{
			   "host": "172.16.238.20",
			   "port": 1980,
			   "weight": 1
		   }]
		}
	}`

	c := `{
		"username": "jack",
		"plugins": {
			"key-auth": {
				"key": "auth-one"
			},
			"basic-auth": {
				"username": "jack",
				"password": "123456"
			}
		},
		"desc": "test description"
	}`

	exportR := `{
		"components": {
			"securitySchemes": {
				"api_key": {
					"in": "header",
					"name": "X-XSRF-TOKEN",
					"type": "apiKey"
				},
				"basicAuth": {
					"in": "header",
					"name": "basicAuth",
					"type": "basicAuth"
				}
			}
		},
		"info": {
			"title": "RoutesExport",
			"version": "3.0.0"
		},
		"openapi": "3.0.0",
		"paths": {
			"/hello": {
				"get": {`

	var route *entity.Route
	var consumer *entity.Consumer
	err := json.Unmarshal([]byte(r), &route)
	err = json.Unmarshal([]byte(c), &consumer)

	l := &Loader{}
	exported, err := l.Export(loader.DataSets{
		Routes:    []entity.Route{*route},
		Consumers: []entity.Consumer{*consumer},
	})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.NotNil(t, ret)
	assert.Contains(t, strings.Replace(string(ret), " ", "", -1), replaceStr(exportR))
}
