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

package authentication

import (
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/apache/apisix-dashboard/api/internal/config"
	"github.com/apache/apisix-dashboard/api/internal/handler"
	"github.com/apache/apisix-dashboard/api/internal/utils/consts"
)

type Handler struct {
	config config.Config
}

func NewHandler() (handler.RouteRegister, error) {
	return &Handler{}, nil
}

func (h *Handler) ApplyRoute(r *gin.Engine, cfg config.Config) {
	h.config = cfg
	r.POST("/apisix/admin/user/login",
		handler.Wrap(h.userLogin, reflect.TypeOf(LoginInput{})))
}

type UserSession struct {
	Token string `json:"token"`
}

// swagger:model LoginInput
type LoginInput struct {
	// username
	Username string `json:"username" binding:"required"`
	// password
	Password string `json:"password" binding:"required"`
}

// swagger:operation POST /apisix/admin/user/login userLogin
//
// user login.
//
// ---
// produces:
// - application/json
// parameters:
//   - name: username
//     in: body
//     description: user name
//     required: true
//     type: string
//   - name: password
//     in: body
//     description: password
//     required: true
//     type: string
//
// responses:
//
//	'0':
//	  description: login success
//	  schema:
//	    "$ref": "#/definitions/ApiError"
//	default:
//	  description: unexpected error
//	  schema:
//	    "$ref": "#/definitions/ApiError"
func (h *Handler) userLogin(_ *gin.Context, input interface{}) handler.Response {
	i := input.(*LoginInput)
	username := i.Username
	password := i.Password

	authnConfig := h.config.Authentication
	var user *config.AuthenticationUser
	for i := range authnConfig.Users {
		if authnConfig.Users[i].Username == username {
			user = &authnConfig.Users[i]
			break
		}
	}
	if user == nil {
		return handler.Response{
			StatusCode: http.StatusUnauthorized,
			ErrMsg:     consts.ErrUsernamePassword.Error(),
		}
	}

	if username != user.Username || password != user.Password {
		return handler.Response{
			StatusCode: http.StatusUnauthorized,
			ErrMsg:     consts.ErrUsernamePassword.Error(),
		}
	}

	// create JWT for session
	claims := jwt.StandardClaims{
		Subject:   username,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(authnConfig.ExpireTime)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(authnConfig.Secret))

	// output token
	return handler.Response{
		Success: true,
		Data: UserSession{
			Token: signedToken,
		},
	}
}
