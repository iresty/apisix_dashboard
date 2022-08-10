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
package openapi3

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/apisix-dashboard/api/internal/core/entity"
)

var (
	TestAPI101 = "../../../../../test/testdata/import/Postman-API101.yaml"
)

// Test API 101 on no MergeMethod mode
func TestParseAPI101NoMerge(t *testing.T) {
	fileContent, err := ioutil.ReadFile(TestAPI101)
	assert.NoError(t, err)

	l := &Loader{MergeMethod: false, TaskName: "test"}
	data, err := l.Import(fileContent)
	assert.NoError(t, err)

	assert.Len(t, data.Routes, 5)
	assert.Len(t, data.Upstreams, 1)

	// Upstream
	assert.Equal(t, "test", data.Upstreams[0].Name)
	assert.Equal(t, "roundrobin", data.Upstreams[0].Type)

	// Route
	assert.Equal(t, data.Upstreams[0].ID, data.Routes[0].UpstreamID)
	for _, route := range data.Routes {
		switch route.Name {
		case "test_customers_GET":
			assert.Contains(t, route.Uris, "/customers")
			assert.Contains(t, route.Methods, "GET")
			assert.Equal(t, "Get all customers", route.Desc)
			assert.Equal(t, entity.Status(0), route.Status)
		case "test_customer_GET":
			assert.Contains(t, route.Uris, "/customer")
			assert.Contains(t, route.Methods, "GET")
			assert.Equal(t, "Get one customer", route.Desc)
			assert.Equal(t, entity.Status(0), route.Status)
		case "test_customer_POST":
			assert.Contains(t, route.Uris, "/customer")
			assert.Contains(t, route.Methods, "POST")
			assert.Equal(t, "Add new customer", route.Desc)
			assert.Equal(t, entity.Status(0), route.Status)
		case "test_customer/{customer_id}_PUT":
			assert.Contains(t, route.Uris, "/customer/*")
			assert.Contains(t, route.Methods, "PUT")
			assert.Equal(t, "Update customer", route.Desc)
			assert.Equal(t, entity.Status(0), route.Status)
		case "test_customer/{customer_id}_DELETE":
			assert.Contains(t, route.Uris, "/customer/*")
			assert.Contains(t, route.Methods, "DELETE")
			assert.Equal(t, "Remove customer", route.Desc)
			assert.Equal(t, entity.Status(0), route.Status)
		default:
			t.Fatal("bad route name exist")
		}
	}
}

// Test API 101 on MergeMethod mode
func TestParseAPI101Merge(t *testing.T) {
	fileContent, err := ioutil.ReadFile(TestAPI101)
	assert.NoError(t, err)

	l := &Loader{MergeMethod: true, TaskName: "test"}
	data, err := l.Import(fileContent)
	assert.NoError(t, err)

	assert.Len(t, data.Routes, 3)
	assert.Len(t, data.Upstreams, 1)

	// Upstream
	assert.Equal(t, "test", data.Upstreams[0].Name)
	assert.Equal(t, "roundrobin", data.Upstreams[0].Type)

	// Route
	assert.Equal(t, data.Upstreams[0].ID, data.Routes[0].UpstreamID)
	for _, route := range data.Routes {
		switch route.Name {
		case "test_customer":
			assert.Contains(t, route.Uris, "/customer")
			assert.Contains(t, route.Methods, "GET", "GET")
			assert.Equal(t, entity.Status(0), route.Status)
		case "test_customers":
			assert.Contains(t, route.Uris, "/customers")
			assert.Contains(t, route.Methods, "GET")
			assert.Equal(t, entity.Status(0), route.Status)
		case "test_customer/{customer_id}":
			assert.Contains(t, route.Uris, "/customer/*")
			assert.Contains(t, route.Methods, "PUT", "DELETE")
			assert.Equal(t, entity.Status(0), route.Status)
		default:
			t.Fatal("bad route name exist")
		}
	}
}
