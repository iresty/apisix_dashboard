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
package conf

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mergeSchema(t *testing.T) {
	type args struct {
		apisixSchema    []byte
		customizeSchema []byte
	}
	tests := []struct {
		name           string
		args           args
		wantRes        []byte
		wantErr        bool
		wantErrMessage string
	}{
		{
			name: "should failed when have duplicates key",
			args: args{
				apisixSchema:    []byte(`{"main":{"a":1,"b":2},"plugins":{"a":1}}`),
				customizeSchema: []byte(`{"main":{"b":1}}`),
			},
			wantErr:        true,
			wantErrMessage: "duplicates key: main.b between schema.json and customize_schema.json",
		},
		{
			name: "should success",
			args: args{
				apisixSchema:    []byte(`{"main":{"a":1,"b":2},"plugins":{"a":1}}`),
				customizeSchema: []byte(`{"main":{"c":3}}`),
			},
			wantErr: false,
			wantRes: []byte(`{"main":{"a":1,"b":2,"c":3},"plugins":{"a":1}}`),
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			var (
				wantMap map[string]interface{}
				gotMap  map[string]interface{}
			)

			got, err := mergeSchema(tt.args.apisixSchema, tt.args.customizeSchema)
			if tt.wantErr {
				assert.Equal(t, tt.wantErrMessage, err.Error())
				return
			}

			assert.NoError(t, err)
			err = json.Unmarshal(got, &gotMap)
			assert.NoError(t, err)
			err = json.Unmarshal(tt.wantRes, &wantMap)
			assert.NoError(t, err)
			assert.Equal(t, wantMap, gotMap)
		})
	}
}
