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
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	"github.com/apisix/manager-api/internal/conf"
)

func genToken(username string, issueAt, expireAt int64) string {
	claims := jwt.StandardClaims{
		Subject:   username,
		IssuedAt:  issueAt,
		ExpiresAt: expireAt,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(conf.AuthConf.Secret))

	return signedToken
}

func TestAuthenticationMiddleware_Handle(t *testing.T) {
	r := gin.New()
	r.Use(Authentication())
	r.GET("/*path", func(c *gin.Context) {
	})

	w := performRequest(r, "GET", "/apisix/admin/user/login", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w = performRequest(r, "GET", "/apisix/admin/routes", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// test with token expire
	expireToken := genToken("admin", time.Now().Unix(), time.Now().Unix()-60*3600)
	w = performRequest(r, "GET", "/apisix/admin/routes", map[string]string{"Authorization": expireToken})
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// test with empty subject
	emptySubjectToken := genToken("", time.Now().Unix(), time.Now().Unix()+60*3600)
	w = performRequest(r, "GET", "/apisix/admin/routes", map[string]string{"Authorization": emptySubjectToken})
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// test token with nonexistent username
	nonexistentUserToken := genToken("user1", time.Now().Unix(), time.Now().Unix()+60*3600)
	w = performRequest(r, "GET", "/apisix/admin/routes", map[string]string{"Authorization": nonexistentUserToken})
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// test auth success
	validToken := genToken("admin", time.Now().Unix(), time.Now().Unix()+60*3600)
	w = performRequest(r, "GET", "/apisix/admin/routes", map[string]string{"Authorization": validToken})
	assert.Equal(t, http.StatusOK, w.Code)
}
