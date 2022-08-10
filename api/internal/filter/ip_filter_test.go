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
package filter

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/apache/apisix-dashboard/api/internal/conf"
)

func TestIPFilter_Handle(t *testing.T) {
	// empty allowed ip list --> should normal
	conf.AllowList = []string{}
	r := gin.New()
	r.Use(IPFilter())

	r.GET("/", func(c *gin.Context) {
	})

	w := performRequest(r, "GET", "/", nil)
	assert.Equal(t, 200, w.Code)

	// should forbidden
	conf.AllowList = []string{"10.0.0.0/8", "10.0.0.1"}
	r = gin.New()
	r.Use(IPFilter())
	r.GET("/fbd", func(c *gin.Context) {
	})

	w = performRequest(r, "GET", "/fbd", nil)
	assert.Equal(t, 403, w.Code)

	// should allowed
	conf.AllowList = []string{"10.0.0.0/8", "0.0.0.0/0"}
	r = gin.New()
	r.Use(IPFilter())
	r.GET("/test", func(c *gin.Context) {
	})
	w = performRequest(r, "GET", "/test", nil)
	assert.Equal(t, 200, w.Code)

	// should forbidden
	conf.AllowList = []string{"127.0.0.1"}
	r = gin.New()
	r.Use(IPFilter())
	r.GET("/test", func(c *gin.Context) {})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "127.0.0.1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)

	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Real-Ip", "127.0.0.1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
}
