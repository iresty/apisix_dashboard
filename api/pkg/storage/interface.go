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
package storage_api

import "context"

type Interface interface {
	Get(ctx context.Context, key string) (string, error)
	List(ctx context.Context, key string) ([]Keypair, error)
	Create(ctx context.Context, key, val string) error
	Update(ctx context.Context, key, val string) error
	BatchDelete(ctx context.Context, keys []string) error
	Watch(ctx context.Context, key string) <-chan WatchResponse
}

type WatchResponse struct {
	Events   []Event
	Error    error
	Canceled bool
}

type Keypair struct {
	Key   string
	Value string
}

type Event struct {
	Keypair
	Type EventType
}

type EventType string

var (
	EventTypePut    EventType = "put"
	EventTypeDelete EventType = "delete"
)
