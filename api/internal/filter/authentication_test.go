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
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/middleware"
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

type mockMiddleware struct {
	middleware.BaseMiddleware
}

func (mw *mockMiddleware) Handle(ctx droplet.Context) error {
	return errors.New("next middleware")
}

func testPanic(t *testing.T, mw AuthenticationMiddleware, ctx droplet.Context) {
	defer func() {
		panicErr := recover()
		assert.Contains(t, panicErr.(error).Error(), "input middleware cannot get http request")
	}()
	_ = mw.Handle(ctx)
}

func TestAuthenticationMiddleware_Handle(t *testing.T) {
	ctx := droplet.NewContext()
	fakeReq, _ := http.NewRequest(http.MethodGet, "", nil)
	expectOutput := &data.SpecCodeResponse{
		Response: data.Response{
			Code:    010013,
			Message: "request unauthorized",
		},
		StatusCode: http.StatusUnauthorized,
	}

	mw := AuthenticationMiddleware{}
	mockMw := mockMiddleware{}
	mw.SetNext(&mockMw)

	// test without http.Request
	testPanic(t, mw, ctx)

	ctx.Set(middleware.KeyHttpRequest, fakeReq)

	// test without token check
	fakeReq.URL = &url.URL{Path: "/apisix/admin/user/login"}
	assert.Equal(t, mw.Handle(ctx), errors.New("next middleware"))

	// test without authorization header
	fakeReq.URL = &url.URL{Path: "/apisix/admin/routes"}
	assert.Nil(t, mw.Handle(ctx))
	assert.Equal(t, expectOutput, ctx.Output().(*data.SpecCodeResponse))

	// test with token expire
	expireToken := genToken("admin", time.Now().Unix(), time.Now().Unix()-60*3600)
	fakeReq.Header.Set("Authorization", expireToken)
	assert.Nil(t, mw.Handle(ctx))
	assert.Equal(t, expectOutput, ctx.Output().(*data.SpecCodeResponse))

	// test with temp subject
	tempSubjectToken := genToken("", time.Now().Unix(), time.Now().Unix()+60*3600)
	fakeReq.Header.Set("Authorization", tempSubjectToken)
	assert.Nil(t, mw.Handle(ctx))
	assert.Equal(t, expectOutput, ctx.Output().(*data.SpecCodeResponse))

	// test username doesn't exist
	userToken := genToken("user1", time.Now().Unix(), time.Now().Unix()+60*3600)
	fakeReq.Header.Set("Authorization", userToken)
	assert.Nil(t, mw.Handle(ctx))
	assert.Equal(t, expectOutput, ctx.Output().(*data.SpecCodeResponse))

	// test auth success
	adminToken := genToken("admin", time.Now().Unix(), time.Now().Unix()+60*3600)
	fakeReq.Header.Set("Authorization", adminToken)
	ctx.SetOutput("test data")
	assert.Equal(t, errors.New("next middleware"), mw.Handle(ctx))
	assert.Equal(t, "test data", ctx.Output().(string))
}
