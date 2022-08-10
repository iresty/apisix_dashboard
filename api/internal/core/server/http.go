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
package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/shiningrush/droplet"

	"github.com/apache/apisix-dashboard/api/internal"
	"github.com/apache/apisix-dashboard/api/internal/conf"
	"github.com/apache/apisix-dashboard/api/internal/handler"
)

func (s *server) setupAPI() {
	// orchestrator
	droplet.Option.Orchestrator = func(mws []droplet.Middleware) []droplet.Middleware {
		var newMws []droplet.Middleware
		// default middleware order: resp_reshape, auto_input, traffic_log
		// We should put err_transform at second to catch all error
		newMws = append(newMws, mws[0], &handler.ErrorTransformMiddleware{})
		newMws = append(newMws, mws[1:]...)
		return newMws
	}

	// routes
	r := internal.SetUpRouter()

	// HTTP
	addr := net.JoinHostPort(conf.ServerHost, strconv.Itoa(conf.ServerPort))
	s.server = &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Duration(1000) * time.Millisecond,
		WriteTimeout: time.Duration(5000) * time.Millisecond,
	}

	// HTTPS
	if conf.SSLCert != "" && conf.SSLKey != "" {
		addrSSL := net.JoinHostPort(conf.SSLHost, strconv.Itoa(conf.SSLPort))
		s.serverSSL = &http.Server{
			Addr:         addrSSL,
			Handler:      r,
			ReadTimeout:  time.Duration(1000) * time.Millisecond,
			WriteTimeout: time.Duration(5000) * time.Millisecond,
			TLSConfig: &tls.Config{
				// Causes servers to use Go's default ciphersuite preferences,
				// which are tuned to avoid attacks. Does nothing on clients.
				PreferServerCipherSuites: true,
			},
		}
	}
}
