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

package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/apache/apisix-dashboard/api/pkg/storage"
)

const (
	adminAPIErrorWarp = "failed to request Admin API"
)

type adminAPIVersion = int

const (
	adminAPIVersion2 adminAPIVersion = 2
	adminApiVersion3 adminAPIVersion = 3
)

// Ensure that the AdminAPI Storage implementation conforms to the storage_api interface definition
var _ storage_api.Interface = AdminAPIStorage{}

// NewAdminAPIStorage will create an instance of AdminAPIStorage and complete the basic configure
func NewAdminAPIStorage() *AdminAPIStorage {
	return &AdminAPIStorage{
		client: resty.New().
			SetBaseURL("http://127.0.0.1:9080/apisix/admin/").
			SetHeader("X-API-KEY", "edd1c9f034335f136f87ad84b625c8f1").
			SetTimeout(1 * time.Second),
	}
}

// AdminAPIStorage is the storage layer implementation of the APISIX Admin API
type AdminAPIStorage struct {
	client *resty.Client
}

func (s AdminAPIStorage) Get(ctx context.Context, key string) (string, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		Get(s.removeKeyPrefix(key))
	if err != nil {
		return "", errors.Wrap(err, adminAPIErrorWarp)
	}

	json := gjson.Parse(resp.String())
	if json.Get("count").Int() < 1 {
		return "", errors.Errorf("key: %s is not found", key)
	}

	return gjson.Get(resp.String(), "node.value").Raw, nil
}

func (s AdminAPIStorage) List(ctx context.Context, key string) ([]storage_api.Keypair, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		Get(s.removeKeyPrefix(key))
	if err != nil {
		return nil, errors.Wrap(err, adminAPIErrorWarp)
	}

	switch resp.Header().Get("X-API-VERSION") {
	case "v3":
		return s.parseList(adminApiVersion3, resp.String()), nil
	case "v2":
		fallthrough
	default:
		return s.parseList(adminAPIVersion2, resp.String()), nil
	}
}

func (s AdminAPIStorage) Create(ctx context.Context, key, val string) error {
	return s.Update(ctx, key, val)
}

func (s AdminAPIStorage) Update(ctx context.Context, key, val string) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(val).
		Put(s.removeKeyPrefix(key))
	if err != nil {
		return errors.Wrap(err, adminAPIErrorWarp)
	}

	if err = s.checkError(resp.String()); err != nil {
		return err
	}

	return nil
}

func (s AdminAPIStorage) BatchDelete(ctx context.Context, keys []string) error {
	for _, key := range keys {
		resp, err := s.client.R().
			SetContext(ctx).
			Delete(s.removeKeyPrefix(key))
		if err != nil {
			return errors.Wrap(err, adminAPIErrorWarp)
		}

		if err = s.checkError(resp.String()); err != nil {
			return fmt.Errorf("delete key[%s] failed: %s", key, err)
		}

		if gjson.Get(resp.String(), "deleted").Int() < 1 {
			return errors.Errorf("key: %s is not found", key)
		}
	}

	return nil
}

// Watch watches for changes in data sources and generates event streams
// Admin API storage implementation does not support Watch capability
func (s AdminAPIStorage) Watch(_ context.Context, _ string) <-chan storage_api.WatchResponse {
	return make(chan storage_api.WatchResponse)
}

// removeKeyPrefix will remove the "/apisix/" prefix of key
// the pattern of the key is usually like "/apisix/routes"
func (AdminAPIStorage) removeKeyPrefix(s string) string {
	return strings.ReplaceAll(s, "/apisix/", "")
}

// checkError extracts errors from the response body
func (AdminAPIStorage) checkError(json string) error {
	if msg := gjson.Get(json, "message"); msg.Exists() {
		return errors.New(msg.String())
	} else if msg := gjson.Get(json, "error_msg"); msg.Exists() {
		return errors.New(msg.String())
	}
	return nil
}

// parseList helps us to extract data from APISIX Admin API v2/v3 list responses
// APISIX Admin API v2 response data structure usually contains some etcd-related attributes that are nested when performing List operations
// The v2 structure looks like this {"count":2,"node":{"key":"/apisix/routes","nodes":[{xxx}, {xxx}],"dir":true},"action":"get"}
// APISIX Admin API v3 response data structure is the new API structure introduced by APISIX in its v3 version
// The v3 structure looks like this {"count":2,"list":[{"key":"/apisix/routes/route1","value":{xxx}}, {xxx}]}
func (AdminAPIStorage) parseList(version adminAPIVersion, json string) []storage_api.Keypair {
	var (
		data   = gjson.Parse(json)
		count  = data.Get("count").Int()
		result []storage_api.Keypair
		path   string
	)

	if version == adminApiVersion3 {
		path = "list.%d"
	} else {
		path = "node.nodes.%d"
	}

	for i := int64(0); i < count; i++ {
		itemPath := fmt.Sprintf(path, i)
		if data.Get(itemPath).Exists() {
			result = append(result, storage_api.Keypair{
				Key:   data.Get(itemPath + ".key").String(),
				Value: data.Get(itemPath + ".value").Raw,
			})
		}
	}

	return result
}
