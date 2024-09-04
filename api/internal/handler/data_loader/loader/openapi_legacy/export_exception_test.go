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

// Add suffix when testing the same URI export "APISIX-REPEAT-URI-" + Millisecond time stamp:  "APISIX-REPEAT-URI-1257894000000"
func TestExportRoutesSameURI(t *testing.T) {
	//*entity.Route
	r1 := `{
		"uris": ["/test-test"],
		"name": "route_all",
		"desc": "所有",
		"methods": ["GET"],
		"hosts": ["test.com"],
		"status": 1,
		"upstream": {
			"nodes": {
				"172.16.238.20:1980": 1
			},
			"type": "roundrobin"
		}
}`

	r2 := `{
		"uris": ["/test-test"],
		"name": "route_all",
		"desc": "所有1",
		"methods": ["GET"],
		"hosts": ["test.com"],
		"status": 1,
		"upstream": {
			"nodes": {
				"172.16.238.20:1980": 1
			},
			"type": "roundrobin"
		}
}`

	r3 := `{
	"uris": ["/test-test"],
	"name": "route_all",
	"desc": "所有2",
	"methods": ["GET"],
	"hosts": ["test.com"],
	"status": 1,
	"upstream": {
		"nodes": {
			"172.16.238.20:1981": 1
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
			"/test-test": {
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
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				}
			},
			"/test-test-APISIX-REPEAT-URI-2": {
				"get": {
					"operationId": "route_allGET",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"summary": "所有1",
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				}
			},
			"/test-test-APISIX-REPEAT-URI-3": {
				"get": {
					"operationId": "route_allGET",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"summary": "所有2",
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1981": 1
						},
						"type": "roundrobin"
					}
				}
			}
		}
	}`
	var route *entity.Route
	var route2 *entity.Route
	var route3 *entity.Route
	err := json.Unmarshal([]byte(r1), &route)
	err = json.Unmarshal([]byte(r2), &route2)
	err = json.Unmarshal([]byte(r3), &route3)

	l := &Loader{}
	exported, err := l.Export(loader.DataSets{
		Routes: []entity.Route{*route, *route2, *route3},
	})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.Equal(t, replaceStr(exportR1), string(ret))
	assert.NotNil(t, ret)
}

func TestExportRoutesMethodsFeildEmpty(t *testing.T) {
	//*entity.Route
	r1 := `{
		"uris": ["/test-test"],
		"name": "route",
		"methods": [],
		"hosts": ["test.com"],
		"status": 1,
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
			"/test-test": {
				"connect": {
					"operationId": "routeCONNECT",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"delete": {
					"operationId": "routeDELETE",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"get": {
					"operationId": "routeGET",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"head": {
					"operationId": "routeHEAD",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"options": {
					"operationId": "routeOPTIONS",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"patch": {
					"operationId": "routePATCH",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
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
					"operationId": "routePOST",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"put": {
					"operationId": "routePUT",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"trace": {
					"operationId": "routeTRACE",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
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
	exported, err := l.Export(loader.DataSets{
		Routes: []entity.Route{*route},
	})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.Equal(t, replaceStr(exportR1), strings.Replace(string(ret), " ", "", -1))
	assert.NotNil(t, ret)
}

func TestExportRoutesMethodsFeildNil(t *testing.T) {
	//*entity.Route
	r1 := `{
		"uris": ["/test-test"],
		"name": "route",
		"hosts": ["test.com"],
		"status": 1,
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
			"/test-test": {
				"connect": {
					"operationId": "routeCONNECT",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"delete": {
					"operationId": "routeDELETE",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"get": {
					"operationId": "routeGET",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"head": {
					"operationId": "routeHEAD",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"options": {
					"operationId": "routeOPTIONS",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"patch": {
					"operationId": "routePATCH",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
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
					"operationId": "routePOST",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"put": {
					"operationId": "routePUT",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
					"x-apisix-priority": 0,
					"x-apisix-status": 1,
					"x-apisix-upstream": {
						"nodes": {
							"172.16.238.20:1980": 1
						},
						"type": "roundrobin"
					}
				},
				"trace": {
					"operationId": "routeTRACE",
					"requestBody": {},
					"responses": {
						"default": {
							"description": ""
						}
					},
					"x-apisix-enable_websocket": false,
					"x-apisix-hosts": ["test.com"],
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
	exported, err := l.Export(loader.DataSets{
		Routes: []entity.Route{*route},
	})
	assert.NoError(t, err)

	ret, err := json.Marshal(exported)
	assert.Nil(t, err)
	assert.Equal(t, replaceStr(exportR1), strings.Replace(string(ret), " ", "", -1))
	assert.NotNil(t, ret)
}
