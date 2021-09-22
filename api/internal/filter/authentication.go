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
	"math"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/middleware"

	"github.com/apisix/manager-api/internal/conf"
	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/core/store"
	"github.com/apisix/manager-api/internal/log"
)

type AuthenticationMiddleware struct {
	middleware.BaseMiddleware
}

func (mw *AuthenticationMiddleware) Handle(ctx droplet.Context) error {
	httpReq := ctx.Get(middleware.KeyHttpRequest)
	if httpReq == nil {
		err := errors.New("input middleware cannot get http request")

		// Wrong usage, just panic here and let recoverHandler to deal with
		panic(err)
	}

	req := httpReq.(*http.Request)

	if req.URL.Path == "/apisix/admin/tool/version" || req.URL.Path == "/apisix/admin/user/login" {
		return mw.BaseMiddleware.Handle(ctx)
	}

	if !strings.HasPrefix(req.URL.Path, "/apisix") {
		return mw.BaseMiddleware.Handle(ctx)
	}

	// Need check the auth header
	tokenStr := req.Header.Get("Authorization")

	// verify token
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.AuthConf.Secret), nil
	})

	// TODO: design the response error code
	response := data.Response{Code: 010013, Message: "request unauthorized"}

	if err != nil || token == nil || !token.Valid {
		log.Warnf("token validate failed: %s", err)
		log.Warn("please check the secret in conf.yaml")
		ctx.SetOutput(&data.SpecCodeResponse{StatusCode: http.StatusUnauthorized, Response: response})
		return nil
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		log.Warnf("token validate failed: %s, %v", err, token.Valid)
		ctx.SetOutput(&data.SpecCodeResponse{StatusCode: http.StatusUnauthorized, Response: response})
		return nil
	}

	if err := token.Claims.Valid(); err != nil {
		log.Warnf("token claims validate failed: %s", err)
		ctx.SetOutput(&data.SpecCodeResponse{StatusCode: http.StatusUnauthorized, Response: response})
		return nil
	}

	if claims.Subject == "" {
		log.Warn("token claims subject empty")
		ctx.SetOutput(&data.SpecCodeResponse{StatusCode: http.StatusUnauthorized, Response: response})
		return nil
	}

	userExist := false
	switch conf.DataSourceType(claims.Audience) {
	case conf.DataSourceTypeEtcd:
		dashboardUserStore := store.GetStore(store.HubKeyDashboardUser)
		ret, err := dashboardUserStore.List(ctx.Context(), store.ListInput{
			Predicate: func(obj interface{}) bool {
				return claims.Subject == obj.(*entity.DashboardUser).Username
			},
			Format: func(obj interface{}) interface{} {
				return obj.(*entity.DashboardUser)
			},
			PageSize:   math.MaxInt32,
			PageNumber: 1,
		})
		if err != nil {
			log.Warn("user not exist in etcd")
			ctx.SetOutput(&data.SpecCodeResponse{StatusCode: http.StatusUnauthorized, Response: response})
			return nil
		}

		if ret.TotalSize > 0 {
			userExist = true
		}
	case conf.DataSourceTypeLocal:
		fallthrough
	default:
		for _, item := range conf.UserList {
			if claims.Subject == item.Username {
				userExist = true
			}
		}
	}

	if !userExist {
		log.Warnf("user not exists by token claims subject %s", claims.Subject)
		ctx.SetOutput(&data.SpecCodeResponse{StatusCode: http.StatusUnauthorized, Response: response})
		return nil
	}

	return mw.BaseMiddleware.Handle(ctx)
}
