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
	"github.com/gin-gonic/gin"

	"github.com/apache/apisix-dashboard/api/internal/conf"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if conf.SecurityConf.AllowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", conf.SecurityConf.AllowOrigin)
		}

		if conf.SecurityConf.AllowHeaders != "" {
			c.Writer.Header().Set("Access-Control-Allow-Headers", conf.SecurityConf.AllowHeaders)
		}

		if conf.SecurityConf.AllowMethods != "" {
			c.Writer.Header().Set("Access-Control-Allow-Methods", conf.SecurityConf.AllowMethods)
		}

		if conf.SecurityConf.AllowCredentials != "" {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", conf.SecurityConf.AllowCredentials)
		}

		if conf.SecurityConf.XFrameOptions != "" {
			c.Writer.Header().Set("X-Frame-Options", conf.SecurityConf.XFrameOptions)
		}

		if conf.SecurityConf.ContentSecurityPolicy != "" {
			c.Writer.Header().Set("Content-Security-Policy", conf.SecurityConf.ContentSecurityPolicy)
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
