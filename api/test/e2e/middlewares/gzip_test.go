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
package middlewares_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("Gzip enable", func() {
	It("get index.html", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:        base.ManagerApiExpect(),
			Method:        http.MethodGet,
			Path:          "/",
			Headers:       map[string]string{"Accept-Encoding": "gzip, deflate, br"},
			ExpectHeaders: map[string]string{"Content-Encoding": "gzip"},
		})
	})
})
