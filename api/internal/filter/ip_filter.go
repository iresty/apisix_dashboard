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
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/apache/apisix-dashboard/api/internal/conf"
	"github.com/apache/apisix-dashboard/api/internal/log"
	"github.com/apache/apisix-dashboard/api/internal/utils/consts"
)

type subnet struct {
	ipStr   string
	ipNet   *net.IPNet
	allowed bool
}

func generateIPSet(ipList []string) (map[string]bool, []*subnet) {
	var ips = map[string]bool{}
	var subnets []*subnet
	for _, ipStr := range ipList {
		if ip, net, err := net.ParseCIDR(ipStr); err == nil {
			if n, total := net.Mask.Size(); n == total {
				ips[ip.String()] = true
				continue
			}

			subnets = append(subnets, &subnet{
				ipStr:   ipStr,
				ipNet:   net,
				allowed: true,
			})
			continue
		}
		if ip := net.ParseIP(ipStr); ip != nil {
			ips[ip.String()] = true
		}
	}

	return ips, subnets
}

func checkIP(ipStr string, ips map[string]bool, subnets []*subnet) bool {
	allowed, ok := ips[ipStr]
	if ok {
		return allowed
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, subnet := range subnets {
		if subnet.ipNet.Contains(ip) {
			return subnet.allowed
		}
	}

	return false
}

func IPFilter() gin.HandlerFunc {
	ips, subnets := generateIPSet(conf.AllowList)
	return func(c *gin.Context) {
		var ipStr string
		if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
			ipStr = ip
		}

		if len(conf.AllowList) < 1 {
			c.Next()
			return
		}

		if ipStr == "" {
			log.Warn("forbidden by empty IP")
			c.AbortWithStatusJSON(http.StatusForbidden, consts.ErrIPNotAllow)
			return
		}

		res := checkIP(ipStr, ips, subnets)
		if !res {
			log.Warnf("forbidden by IP: %s, allowed list: %v", ipStr, conf.AllowList)
			c.AbortWithStatusJSON(http.StatusForbidden, consts.ErrIPNotAllow)
		}

		c.Next()
	}
}
