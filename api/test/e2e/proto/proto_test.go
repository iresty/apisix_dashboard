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
package proto

import (
	"encoding/json"
	"net/http"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"github.com/apache/apisix-dashboard/api/test/e2e/base"
)

var correctProtobuf = `syntax = "proto3";
    package helloworld;
    service Greeter {
        rpc SayHello (HelloRequest) returns (HelloReply) {}
    }
    message HelloRequest {
        string name = 1;
    }
    message HelloReply {
        string message = 1;
    }`

var _ = ginkgo.Describe("Proto", func() {
	ginkgo.It("create proto success", func() {
		createProtoBody := make(map[string]interface{})
		createProtoBody["id"] = 1
		createProtoBody["desc"] = "test_proto1"
		createProtoBody["content"] = correctProtobuf

		_createProtoBody, err := json.Marshal(createProtoBody)
		gomega.Expect(err).To(gomega.BeNil())

		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPost,
			Path:         "/apisix/admin/proto",
			Body:         string(_createProtoBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("create proto failed, id existed", func() {
		createProtoBody := make(map[string]interface{})
		createProtoBody["id"] = 1
		createProtoBody["desc"] = "test_proto1"
		createProtoBody["content"] = correctProtobuf

		_createProtoBody, err := json.Marshal(createProtoBody)
		gomega.Expect(err).To(gomega.BeNil())

		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPost,
			Path:         "/apisix/admin/proto",
			Body:         string(_createProtoBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectBody:   "proto id exists",
			ExpectStatus: http.StatusBadRequest,
		})
	})
	ginkgo.It("update proto success", func() {
		updateProtoBody := make(map[string]interface{})
		updateProtoBody["id"] = 1
		updateProtoBody["desc"] = "test_proto1_modify"
		updateProtoBody["content"] = correctProtobuf

		_updateProtoBody, err := json.Marshal(updateProtoBody)
		gomega.Expect(err).To(gomega.BeNil())

		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPut,
			Path:         "/apisix/admin/proto",
			Body:         string(_updateProtoBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectBody:   "test_proto1_modify",
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("list proto", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/proto",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectBody:   "test_proto1_modify",
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("get proto", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/proto/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectBody:   "test_proto1_modify",
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("delete not existed proto", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/proto/not-exist",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
		})
	})
	ginkgo.It("delete proto", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/proto/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
})

var _ = ginkgo.Describe("Proto with grpc-transcode plugin", func() {
	ginkgo.It("create proto success", func() {
		createProtoBody := make(map[string]interface{})
		createProtoBody["id"] = 1
		createProtoBody["desc"] = "test_proto1"
		createProtoBody["content"] = correctProtobuf

		_createProtoBody, err := json.Marshal(createProtoBody)
		gomega.Expect(err).To(gomega.BeNil())

		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPost,
			Path:         "/apisix/admin/proto",
			Body:         string(_createProtoBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("create route with grpc-transcode", func() {
		createRouteBody := make(map[string]interface{})
		createRouteBody["id"] = 1
		createRouteBody["name"] = "test_route"
		createRouteBody["uri"] = "/grpc_test"
		createRouteBody["methods"] = []string{"GET", "POST"}
		createRouteBody["upstream"] = map[string]interface{}{
			"nodes": []map[string]interface{}{
				{
					"host":   base.UpstreamGrpcIp,
					"port":   50051,
					"weight": 1,
				},
			},
			"type":   "roundrobin",
			"scheme": "grpc",
		}
		createRouteBody["plugins"] = map[string]interface{}{
			"grpc-transcode": map[string]interface{}{
				"method":   "SayHello",
				"proto_id": "1",
				"service":  "helloworld.Greeter",
			},
		}

		_createRouteBody, err := json.Marshal(createRouteBody)
		gomega.Expect(err).To(gomega.BeNil())

		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodPost,
			Path:         "/apisix/admin/routes",
			Body:         string(_createRouteBody),
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("hit GET route for grpc-transcode test", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/grpc_test",
			Query:        "name=world",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectBody:   "{\"message\":\"Hello world\"}",
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("hit POST route for grpc-transcode test", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodPost,
			Path:         "/grpc_test",
			Body:         "name=world",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectBody:   "{\"message\":\"Hello world\"}",
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("hit JSON POST route for grpc-transcode test", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodPost,
			Path:         "/grpc_test",
			Body:         "{\"name\": \"world\"}",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectBody:   "{\"message\":\"Hello world\"}",
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("delete route used proto", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/proto/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectBody:   "proto used check invalid: route",
			ExpectStatus: http.StatusBadRequest,
		})
	})
	ginkgo.It("delete conflict route", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
	ginkgo.It("delete proto again", func() {
		base.RunTestCase(base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/proto/1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		})
	})
})
