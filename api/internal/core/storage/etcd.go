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
	"github.com/apache/apisix-dashboard/api/pkg/storage"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"

	"github.com/apache/apisix-dashboard/api/internal/conf"
	"github.com/apache/apisix-dashboard/api/internal/log"
	"github.com/apache/apisix-dashboard/api/internal/utils"
	"github.com/apache/apisix-dashboard/api/internal/utils/runtime"
)

const (
	// SkippedValueEtcdInitDir indicates the init_dir
	// etcd event will be skipped.
	SkippedValueEtcdInitDir = "init_dir"

	// SkippedValueEtcdEmptyObject indicates the data with an
	// empty JSON value {}, which may be set by APISIX,
	// should be also skipped.
	//
	// Important: at present, {} is considered as invalid,
	// but may be changed in the future.
	SkippedValueEtcdEmptyObject = "{}"
)

var (
	Client *clientv3.Client
)

type EtcdV3Storage struct {
	client *clientv3.Client
}

func InitETCDClient(etcdConf *conf.Etcd) error {
	config := clientv3.Config{
		Endpoints:   etcdConf.Endpoints,
		DialTimeout: 5 * time.Second,
		Username:    etcdConf.Username,
		Password:    etcdConf.Password,
	}
	// mTLS
	if etcdConf.MTLS != nil && etcdConf.MTLS.CaFile != "" &&
		etcdConf.MTLS.CertFile != "" && etcdConf.MTLS.KeyFile != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      etcdConf.MTLS.CertFile,
			KeyFile:       etcdConf.MTLS.KeyFile,
			TrustedCAFile: etcdConf.MTLS.CaFile,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return err
		}
		config.TLS = tlsConfig
	}

	cli, err := clientv3.New(config)
	if err != nil {
		log.Errorf("init etcd failed: %s", err)
		return fmt.Errorf("init etcd failed: %s", err)
	}

	Client = cli
	utils.AppendToClosers(Close)
	return nil
}

func GenEtcdStorage() *EtcdV3Storage {
	return &EtcdV3Storage{
		client: Client,
	}
}

func NewETCDStorage(etcdConf *conf.Etcd) (*EtcdV3Storage, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdConf.Endpoints,
		DialTimeout: 5 * time.Second,
		Username:    etcdConf.Username,
		Password:    etcdConf.Password,
	})
	if err != nil {
		log.Errorf("init etcd failed: %s", err)
		return nil, fmt.Errorf("init etcd failed: %s", err)
	}

	s := &EtcdV3Storage{
		client: cli,
	}

	utils.AppendToClosers(s.Close)
	return s, nil
}

func Close() error {
	if err := Client.Close(); err != nil {
		log.Errorf("etcd client close failed: %s", err)
		return err
	}
	return nil
}

func (s *EtcdV3Storage) Close() error {
	if err := s.client.Close(); err != nil {
		log.Errorf("etcd client close failed: %s", err)
		return err
	}
	return nil
}

func (s *EtcdV3Storage) Get(ctx context.Context, key string) (string, error) {
	resp, err := s.client.Get(ctx, key)
	if err != nil {
		log.Errorf("etcd get failed: %s", err)
		return "", fmt.Errorf("etcd get failed: %s", err)
	}
	if resp.Count == 0 {
		log.Warnf("key: %s is not found", key)
		return "", fmt.Errorf("key: %s is not found", key)
	}

	return string(resp.Kvs[0].Value), nil
}

func (s *EtcdV3Storage) List(ctx context.Context, key string) ([]storage_api.Keypair, error) {
	resp, err := s.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.Errorf("etcd get failed: %s", err)
		return nil, fmt.Errorf("etcd get failed: %s", err)
	}
	var ret []storage_api.Keypair
	for i := range resp.Kvs {
		key := string(resp.Kvs[i].Key)
		value := string(resp.Kvs[i].Value)

		// Skip the data if its value is init_dir or {}
		// during fetching-all phase.
		//
		// For more complex cases, an explicit function to determine if
		// skippable would be better.
		if value == SkippedValueEtcdInitDir || value == SkippedValueEtcdEmptyObject {
			continue
		}

		data := storage_api.Keypair{
			Key:   key,
			Value: value,
		}
		ret = append(ret, data)
	}

	return ret, nil
}

func (s *EtcdV3Storage) Create(ctx context.Context, key, val string) error {
	_, err := s.client.Put(ctx, key, val)
	if err != nil {
		log.Errorf("etcd put failed: %s", err)
		return fmt.Errorf("etcd put failed: %s", err)
	}
	return nil
}

func (s *EtcdV3Storage) Update(ctx context.Context, key, val string) error {
	_, err := s.client.Put(ctx, key, val)
	if err != nil {
		log.Errorf("etcd put failed: %s", err)
		return fmt.Errorf("etcd put failed: %s", err)
	}
	return nil
}

func (s *EtcdV3Storage) BatchDelete(ctx context.Context, keys []string) error {
	for i := range keys {
		resp, err := s.client.Delete(ctx, keys[i])
		if err != nil {
			log.Errorf("delete etcd key[%s] failed: %s", keys[i], err)
			return fmt.Errorf("delete etcd key[%s] failed: %s", keys[i], err)
		}
		if resp.Deleted == 0 {
			log.Warnf("key: %s is not found", keys[i])
			return fmt.Errorf("key: %s is not found", keys[i])
		}
	}
	return nil
}

func (s *EtcdV3Storage) Watch(ctx context.Context, key string) <-chan storage_api.WatchResponse {
	eventChan := s.client.Watch(ctx, key, clientv3.WithPrefix())
	ch := make(chan storage_api.WatchResponse, 1)
	go func() {
		defer runtime.HandlePanic()
		for event := range eventChan {
			output := storage_api.WatchResponse{
				Canceled: event.Canceled,
			}

			for i := range event.Events {
				key := string(event.Events[i].Kv.Key)
				value := string(event.Events[i].Kv.Value)

				// Skip the data if its value is init_dir or {}
				// during watching phase.
				//
				// For more complex cases, an explicit function to determine if
				// skippable would be better.
				if value == SkippedValueEtcdInitDir || value == SkippedValueEtcdEmptyObject {
					continue
				}

				e := storage_api.Event{
					Keypair: storage_api.Keypair{
						Key:   key,
						Value: value,
					},
				}
				switch event.Events[i].Type {
				case clientv3.EventTypePut:
					e.Type = storage_api.EventTypePut
				case clientv3.EventTypeDelete:
					e.Type = storage_api.EventTypeDelete
				}
				output.Events = append(output.Events, e)
			}
			if output.Canceled {
				log.Error("channel canceled")
				output.Error = fmt.Errorf("channel canceled")
			}
			ch <- output
		}

		close(ch)
	}()

	return ch
}
