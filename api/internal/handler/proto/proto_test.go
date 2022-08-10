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
package proto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/apisix-dashboard/api/internal/core/entity"
)

func TestStructUnmarshal(t *testing.T) {
	// define and parse data
	jsonStr := `{
    "id": 1,
    "create_time": 1700000000,
    "update_time": 1700000000,
    "desc": "desc",
    "content": "content"
}`
	proto := entity.Proto{}
	err := json.Unmarshal([]byte(jsonStr), &proto)

	// asserts
	assert.Nil(t, err)
	assert.Equal(t, proto.ID, float64(1))
	assert.Equal(t, proto.CreateTime, int64(1700000000))
	assert.Equal(t, proto.UpdateTime, int64(1700000000))
	assert.Equal(t, proto.Desc, "desc")
	assert.Equal(t, proto.Content, "content")
}
