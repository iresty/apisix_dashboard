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
package base

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

var (
	token string

	UpstreamIp             = "172.16.238.20"
	UpstreamGrpcIp         = "172.16.238.21"
	APISIXHost             = "http://127.0.0.1:9080"
	APISIXInternalUrl      = "http://172.16.238.30:9080"
	APISIXSingleWorkerHost = "http://127.0.0.1:9081"
	ManagerAPIHost         = "http://127.0.0.1:9000"
	PrometheusExporter     = "http://127.0.0.1:9091"
)

func GetToken() string {
	if token != "" {
		return token
	}

	requestBody := `{
		"username": "admin",
		"password": "admin"
	}`

	url := ManagerAPIHost + "/apisix/admin/user/login"
	body, _, err := HttpPost(url, nil, requestBody)
	if err != nil {
		panic(err)
	}

	respond := gjson.ParseBytes(body)
	token = respond.Get("data.token").String()

	return token
}

func getTestingHandle() httpexpect.LoggerReporter {
	return ginkgo.GinkgoT()
}

func ManagerApiExpect() *httpexpect.Expect {
	t := getTestingHandle()
	return httpexpect.New(t, ManagerAPIHost)
}

func APISIXExpect() *httpexpect.Expect {
	t := getTestingHandle()
	return httpexpect.New(t, APISIXHost)
}

func APISIXStreamProxyExpect(port uint16, sni string) *httpexpect.Expect {
	if port == 0 {
		port = 10090
	}
	t := getTestingHandle()

	if sni != "" {
		addr := net.JoinHostPort(sni, strconv.Itoa(int(port)))
		return httpexpect.WithConfig(httpexpect.Config{
			BaseURL:  "https://" + addr,
			Reporter: httpexpect.NewAssertReporter(t),
			Client: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						// accept any certificate; for testing only!
						InsecureSkipVerify: true,
					},
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						addr = net.JoinHostPort("127.0.0.1", strconv.Itoa(int(port)))
						dialer := &net.Dialer{}
						return dialer.DialContext(ctx, network, addr)
					},
				},
			},
		})
	} else {
		return httpexpect.New(t, "http://"+net.JoinHostPort("127.0.0.1", strconv.Itoa(int(port))))
	}
}

func PrometheusExporterExpect() *httpexpect.Expect {
	t := getTestingHandle()
	return httpexpect.New(t, PrometheusExporter)
}

func APISIXHTTPSExpect() *httpexpect.Expect {
	t := getTestingHandle()
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "https://www.test2.com:9443",
		Reporter: httpexpect.NewAssertReporter(t),
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// accept any certificate; for testing only!
					InsecureSkipVerify: true,
				},
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					if addr == "www.test2.com:9443" {
						addr = "127.0.0.1:9443"
					}
					dialer := &net.Dialer{}
					return dialer.DialContext(ctx, network, addr)
				},
			},
		},
	})

	return e
}

var SleepTime = time.Duration(300) * time.Millisecond

type HttpTestCase struct {
	Desc          string
	Object        *httpexpect.Expect
	Method        string
	Path          string
	Query         string
	Body          string
	Headers       map[string]string
	Headers_test  map[string]interface{}
	ExpectStatus  int
	ExpectCode    int
	ExpectMessage string
	ExpectBody    interface{}
	UnexpectBody  interface{}
	ExpectHeaders map[string]string
	Sleep         time.Duration //ms
}

func RunTestCase(tc HttpTestCase) {
	//init
	expectObj := tc.Object
	var req *httpexpect.Request
	switch tc.Method {
	case http.MethodGet:
		req = expectObj.GET(tc.Path)
	case http.MethodPut:
		req = expectObj.PUT(tc.Path)
	case http.MethodPost:
		req = expectObj.POST(tc.Path)
	case http.MethodDelete:
		req = expectObj.DELETE(tc.Path)
	case http.MethodPatch:
		req = expectObj.PATCH(tc.Path)
	case http.MethodOptions:
		req = expectObj.OPTIONS(tc.Path)
	default:
	}

	if req == nil {
		panic("fail to init request")
	}

	if tc.Sleep != 0 {
		time.Sleep(tc.Sleep)
	} else {
		time.Sleep(time.Duration(50) * time.Millisecond)
	}

	if tc.Query != "" {
		req.WithQueryString(tc.Query)
	}

	// set header
	setContentType := false
	for key, val := range tc.Headers {
		req.WithHeader(key, val)
		if strings.ToLower(key) == "content-type" {
			setContentType = true
		}
	}

	// set default content-type
	if !setContentType {
		req.WithHeader("Content-Type", "application/json")
	}

	// set body
	if tc.Body != "" {
		req.WithText(tc.Body)
	}

	// respond check
	resp := req.Expect()

	// match http status
	if tc.ExpectStatus != 0 {
		resp.Status(tc.ExpectStatus)
	}

	// match headers
	if tc.ExpectHeaders != nil {
		for key, val := range tc.ExpectHeaders {
			resp.Header(key).Equal(val)
		}
	}

	// match body
	if tc.ExpectBody != nil {
		//assert.Contains(t, []string{"string", "[]string"}, reflect.TypeOf(tc.ExpectBody).String())
		if body, ok := tc.ExpectBody.(string); ok {
			if body == "" {
				// "" indicates the body is expected to be empty
				resp.Body().Empty()
			} else {
				resp.Body().Contains(body)
			}
		} else if bodies, ok := tc.ExpectBody.([]string); ok && len(bodies) != 0 {
			for _, b := range bodies {
				resp.Body().Contains(b)
			}
		}
	}

	// match UnexpectBody
	if tc.UnexpectBody != nil {
		//assert.Contains(t, []string{"string", "[]string"}, reflect.TypeOf(tc.UnexpectBody).String())
		if body, ok := tc.UnexpectBody.(string); ok {
			// "" indicates the body is expected to be non empty
			if body == "" {
				resp.Body().NotEmpty()
			} else {
				resp.Body().NotContains(body)
			}
		} else if bodies, ok := tc.UnexpectBody.([]string); ok && len(bodies) != 0 {
			for _, b := range bodies {
				resp.Body().NotContains(b)
			}
		}
	}
}

func ReadAPISIXErrorLog() string {
	t := getTestingHandle()
	cmd := exec.Command("pwd")
	pwdByte, err := cmd.CombinedOutput()
	pwd := string(pwdByte)
	pwd = strings.Replace(pwd, "\n", "", 1)
	pwd = pwd[:strings.Index(pwd, "/e2e")]
	bytes, err := ioutil.ReadFile(pwd + "/docker/apisix_logs/error.log")
	assert.Nil(t, err)
	logContent := string(bytes)

	return logContent
}

func CleanAPISIXErrorLog() {
	t := getTestingHandle()
	cmd := exec.Command("pwd")
	pwdByte, err := cmd.CombinedOutput()
	pwd := string(pwdByte)
	pwd = strings.Replace(pwd, "\n", "", 1)
	pwd = pwd[:strings.Index(pwd, "/e2e")]
	cmdStr := "echo | sudo tee " + pwd + "/docker/apisix_logs/error.log"
	cmd = exec.Command("bash", "-c", cmdStr)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println("cmd error:", err.Error())
	}
	assert.Nil(t, err)
}

func GetResourceList(resource string) string {
	t := getTestingHandle()
	body, _, err := HttpGet(ManagerAPIHost+"/apisix/admin/"+resource, map[string]string{"Authorization": GetToken()})
	assert.Nil(t, err)
	return string(body)
}

func CleanResource(resource string) {
	resources := GetResourceList(resource)
	list := gjson.Get(resources, "data.rows").Value().([]interface{})
	for _, item := range list {
		resourceObj := item.(map[string]interface{})
		idTag := "id"
		if resource == "consumers" {
			idTag = "username"
		}
		tc := HttpTestCase{
			Desc:    "delete " + resource + "/" + resourceObj[idTag].(string),
			Object:  ManagerApiExpect(),
			Method:  http.MethodDelete,
			Path:    "/apisix/admin/" + resource + "/" + resourceObj[idTag].(string),
			Headers: map[string]string{"Authorization": GetToken()},
		}
		RunTestCase(tc)
	}
	time.Sleep(SleepTime)
}

var jwtToken string

func GetJwtToken(userKey string) string {
	if jwtToken != "" {
		return jwtToken
	}
	time.Sleep(SleepTime)

	body, status, err := HttpGet(APISIXHost+"/apisix/plugin/jwt/sign?key="+userKey, nil)
	assert.Nil(ginkgo.GinkgoT(), err)
	assert.Equal(ginkgo.GinkgoT(), http.StatusOK, status)
	jwtToken = string(body)

	return jwtToken
}
