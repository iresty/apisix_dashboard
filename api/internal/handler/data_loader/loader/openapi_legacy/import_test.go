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
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	TestFilePrefix = "../../../../../"
)

func TestImport_invalid_content(t *testing.T) {
	l := &Loader{}
	_, err := l.Import([]byte(`{"test": "a"}`))
	assert.EqualError(t, err, "empty or invalid imported file")
}

func TestImport_with_service_id(t *testing.T) {
	fileContent, err := ioutil.ReadFile(TestFilePrefix + "test/testdata/import/with-service-id.yaml")
	assert.NoError(t, err)

	l := &Loader{}
	data, err := l.Import(fileContent)
	assert.NoError(t, err)

	assert.Equal(t, "service1", data.Routes[0].ServiceID)
}

func TestImport_with_upstream_id(t *testing.T) {
	fileContent, err := ioutil.ReadFile(TestFilePrefix + "test/testdata/import/with-upstream-id.yaml")
	assert.NoError(t, err)

	l := &Loader{}
	data, err := l.Import(fileContent)
	assert.NoError(t, err)

	assert.Equal(t, "upstream1", data.Routes[0].UpstreamID)
}
