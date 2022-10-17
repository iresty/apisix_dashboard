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
package ssl_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/gomega"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var _ = Describe("SSL", func() {
	var (
		testCert   []byte
		testKey    []byte
		validBody  []byte
		validBody2 []byte
		patchBody  []byte
	)

	var err error
	testCert, err = ioutil.ReadFile("../../certs/test2.crt")
	gomega.Expect(err).To(gomega.BeNil())
	testKey, err = ioutil.ReadFile("../../certs/test2.key")
	gomega.Expect(err).To(gomega.BeNil())

	validBody, err = json.Marshal(map[string]interface{}{
		"cert": string(testCert),
		"key":  string(testKey),
		"sni":  "test.com",
		"labels": map[string]string{
			"build":   "16",
			"env":     "production",
			"version": "v3",
		},
	})
	gomega.Expect(err).To(gomega.BeNil())

	validBody2, err = json.Marshal(map[string]interface{}{
		"cert": string(testCert),
		"key":  string(testKey),
		"sni":  "test1.com",
		"labels": map[string]string{
			"build":   "17",
			"env":     "production",
			"version": "v3",
		},
	})
	gomega.Expect(err).To(gomega.BeNil())

	patchBody, err = json.Marshal(map[string]interface{}{
		"sni": "test2.com",
	})
	gomega.Expect(err).To(gomega.BeNil())

	DescribeTable("test ssl create and update",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		Entry("create ssl1 by id", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/ssls/1",
			Body:         string(validBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get ssl1", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/ssls/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"key\":\"/apisix/ssls/1\"", "\"build\":\"16\""},
		}),
		Entry("create ssl2 by random id", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPost,
			Path:         "/apisix/admin/ssls",
			Body:         string(validBody2),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get ssls", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/ssls",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"key\":\"/apisix/ssls/1\"", "\"build\":\"17\""},
		}),
		Entry("patch ssl1", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPatch,
			Path:         "/apisix/admin/ssls/1",
			Body:         string(patchBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		Entry("get ssl1", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/ssls/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   []string{"\"key\":\"/apisix/ssls/1\"", "\"sni\":\"test2.com\""},
		}),
	)
})
