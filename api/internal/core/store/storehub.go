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
package store

import (
	"fmt"
	"reflect"

	"github.com/apisix/manager-api/internal/conf"
	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/log"
	"github.com/apisix/manager-api/internal/utils"
)

type HubKey string

const (
	HubKeyConsumer     HubKey = "consumer"
	HubKeyRoute        HubKey = "route"
	HubKeyService      HubKey = "service"
	HubKeySsl          HubKey = "ssl"
	HubKeyUpstream     HubKey = "upstream"
	HubKeyScript       HubKey = "script"
	HubKeyGlobalRule   HubKey = "global_rule"
	HubKeyServerInfo   HubKey = "server_info"
	HubKeyPluginConfig HubKey = "plugin_config"

	HubKeyDashboardUser HubKey = "dashboard/user"
)

var (
	storeHub = map[HubKey]*GenericStore{}
)

func InitStore(key HubKey, opt GenericStoreOption) error {
	hubsNeedCheck := map[HubKey]bool{
		HubKeyConsumer:   true,
		HubKeyRoute:      true,
		HubKeySsl:        true,
		HubKeyService:    true,
		HubKeyUpstream:   true,
		HubKeyGlobalRule: true,
	}
	if _, ok := hubsNeedCheck[key]; ok {
		validator, err := NewAPISIXJsonSchemaValidator("main." + string(key))
		if err != nil {
			return err
		}
		opt.Validator = validator
	}
	s, err := NewGenericStore(opt)
	if err != nil {
		log.Errorf("NewGenericStore error: %s", err)
		return err
	}
	if err := s.Init(); err != nil {
		log.Errorf("GenericStore init error: %s", err)
		return err
	}

	utils.AppendToClosers(s.Close)
	storeHub[key] = s
	return nil
}

func GetStore(key HubKey) *GenericStore {
	if s, ok := storeHub[key]; ok {
		return s
	}
	panic(fmt.Sprintf("no store with key: %s", key))
}

func RangeStore(f func(key HubKey, store *GenericStore) bool) {
	for k, s := range storeHub {
		if k != "" && s != nil {
			if !f(k, s) {
				break
			}
		}
	}
}

func InitStores() error {
	err := InitStore(HubKeyConsumer, GenericStoreOption{
		BasePath: conf.ETCDConfig.Prefix + "/consumers",
		ObjType:  reflect.TypeOf(entity.Consumer{}),
		KeyFunc: func(obj interface{}) string {
			r := obj.(*entity.Consumer)
			return r.Username
		},
	})
	if err != nil {
		return err
	}

	err = InitStore(HubKeyRoute, GenericStoreOption{
		BasePath: conf.ETCDConfig.Prefix + "/routes",
		ObjType:  reflect.TypeOf(entity.Route{}),
		KeyFunc: func(obj interface{}) string {
			r := obj.(*entity.Route)
			return utils.InterfaceToString(r.ID)
		},
	})
	if err != nil {
		return err
	}

	err = InitStore(HubKeyService, GenericStoreOption{
		BasePath: conf.ETCDConfig.Prefix + "/services",
		ObjType:  reflect.TypeOf(entity.Service{}),
		KeyFunc: func(obj interface{}) string {
			r := obj.(*entity.Service)
			return utils.InterfaceToString(r.ID)
		},
	})
	if err != nil {
		return err
	}

	err = InitStore(HubKeySsl, GenericStoreOption{
		BasePath: conf.ETCDConfig.Prefix + "/ssl",
		ObjType:  reflect.TypeOf(entity.SSL{}),
		KeyFunc: func(obj interface{}) string {
			r := obj.(*entity.SSL)
			return utils.InterfaceToString(r.ID)
		},
	})
	if err != nil {
		return err
	}

	err = InitStore(HubKeyUpstream, GenericStoreOption{
		BasePath: conf.ETCDConfig.Prefix + "/upstreams",
		ObjType:  reflect.TypeOf(entity.Upstream{}),
		KeyFunc: func(obj interface{}) string {
			r := obj.(*entity.Upstream)
			return utils.InterfaceToString(r.ID)
		},
	})
	if err != nil {
		return err
	}

	err = InitStore(HubKeyScript, GenericStoreOption{
		BasePath: conf.ETCDConfig.Prefix + "/scripts",
		ObjType:  reflect.TypeOf(entity.Script{}),
		KeyFunc: func(obj interface{}) string {
			r := obj.(*entity.Script)
			return r.ID
		},
	})
	if err != nil {
		return err
	}

	err = InitStore(HubKeyGlobalRule, GenericStoreOption{
		BasePath: conf.ETCDConfig.Prefix + "/global_rules",
		ObjType:  reflect.TypeOf(entity.GlobalPlugins{}),
		KeyFunc: func(obj interface{}) string {
			r := obj.(*entity.GlobalPlugins)
			return utils.InterfaceToString(r.ID)
		},
	})
	if err != nil {
		return err
	}

	err = InitStore(HubKeyServerInfo, GenericStoreOption{
		BasePath: conf.ETCDConfig.Prefix + "/data_plane/server_info",
		ObjType:  reflect.TypeOf(entity.ServerInfo{}),
		KeyFunc: func(obj interface{}) string {
			r := obj.(*entity.ServerInfo)
			return utils.InterfaceToString(r.ID)
		},
	})
	if err != nil {
		return err
	}

	err = InitStore(HubKeyPluginConfig, GenericStoreOption{
		BasePath: "/apisix/plugin_configs",
		ObjType:  reflect.TypeOf(entity.PluginConfig{}),
		KeyFunc: func(obj interface{}) string {
			r := obj.(*entity.PluginConfig)
			return utils.InterfaceToString(r.ID)
		},
	})
	if err != nil {
		return err
	}

	err = InitStore(HubKeyDashboardUser, GenericStoreOption{
		BasePath: "/dashboard/users",
		ObjType:  reflect.TypeOf(entity.DashboardUser{}),
		KeyFunc: func(obj interface{}) string {
			r := obj.(*entity.DashboardUser)
			return utils.InterfaceToString(r.ID)
		},
	})
	if err != nil {
		return err
	}

	return nil
}
