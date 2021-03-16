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
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/stretchr/testify/assert"
)

type UploadFile struct {
	Name     string
	Filepath string
}

func HttpGet(url string, headers map[string]string) ([]byte, int, error) {
	return httpRequest(http.MethodGet, url, headers, "")
}

func HttpDelete(url string, headers map[string]string) ([]byte, int, error) {
	return httpRequest(http.MethodDelete, url, headers, "")
}

func HttpPut(url string, headers map[string]string, reqBody string) ([]byte, int, error) {
	return httpRequest(http.MethodPut, url, headers, reqBody)
}

func HttpPost(url string, headers map[string]string, reqBody string) ([]byte, int, error) {
	return httpRequest(http.MethodPost, url, headers, reqBody)
}

func httpRequest(method, url string, headers map[string]string, reqBody string) ([]byte, int, error) {
	var requestBody = new(bytes.Buffer)
	if reqBody != "" {
		requestBody = bytes.NewBuffer([]byte(reqBody))
	}
	req, err := http.NewRequest(method, url, requestBody)

	req.Close = true
	if err != nil {
		return nil, 0, err
	}

	// set header
	for key, val := range headers {
		req.Header.Add(key, val)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

func BatchTestServerPort(times int) map[string]int {
	t := getTestingHandle()
	url := APISIXSingleWorkerHost + "/server_port"
	res := map[string]int{}

	for i := 0; i < times; i++ {
		bodyByte, status, err := HttpGet(url, nil)
		assert.Nil(t, err)
		assert.Equal(t, 200, status)

		body := string(bodyByte)
		if _, ok := res[body]; !ok {
			res[body] = 1
		} else {
			res[body] += 1
		}
	}

	return res
}

func HttpPostFile(url string, reqParams map[string]string, files []UploadFile, headers map[string]string) ([]byte, int, error) {
	requestBody, newContentType, err := GetReader(reqParams, "multipart/form-data", files)
	mustCheck(err)
	httpRequest, _ := http.NewRequest("POST", url, requestBody)
	httpRequest.Header.Add("Content-Type", newContentType)
	for k, v := range headers {
		httpRequest.Header.Add(k, v)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(httpRequest)
	mustCheck(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

func GetReader(reqParams map[string]string, contentType string, files []UploadFile) (io.Reader, string, error) {
	if strings.Index(contentType, "json") > -1 {
		bytesData, _ := json.Marshal(reqParams)
		return bytes.NewReader(bytesData), contentType, nil
	}
	if files != nil {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		for _, uploadFile := range files {
			file, err := os.Open(uploadFile.Filepath)
			if err != nil {
				return nil, "", err
			}
			defer file.Close()
			part, err := writer.CreateFormFile(uploadFile.Name, filepath.Base(uploadFile.Filepath))
			if err != nil {
				return nil, "", err
			}
			_, err = io.Copy(part, file)
		}
		for k, v := range reqParams {
			if err := writer.WriteField(k, v); err != nil {
				return nil, "", err
			}
		}
		if err := writer.Close(); err != nil {
			return nil, "", err
		}
		return body, writer.FormDataContentType(), nil
	}

	urlValues := url.Values{}
	for key, val := range reqParams {
		urlValues.Set(key, val)
	}

	reqBody := urlValues.Encode()

	return strings.NewReader(reqBody), contentType, nil
}

func mustCheck(err error) {
	if err != nil {
		panic(err)
	}
}
