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

package server

import (
	"github.com/pkg/errors"

	"github.com/apache/apisix-dashboard/api/internal/core/storage"
	"github.com/apache/apisix-dashboard/api/internal/core/store"
	"github.com/apache/apisix-dashboard/api/internal/log"
)

func (s *server) setupStore() error {
	dataSourceConfig := s.options.Config.DataSource

	if len(dataSourceConfig) <= 0 {
		return errors.New("no data source is configured")
	}

	etcdConfig := dataSourceConfig[0].ETCD
	if err := storage.InitETCDClient(etcdConfig); err != nil {
		log.Errorf("init etcd client fail: %w", err)
		return err
	}

	if err := store.InitStores(etcdConfig.Prefix); err != nil {
		log.Errorf("init stores fail: %w", err)
		return err
	}

	return nil
}
