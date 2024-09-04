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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap/buffer"

	"github.com/apisix/manager-api/internal/conf"
	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/log"
)

type Validator interface {
	Validate(obj interface{}) error
}
type JsonSchemaValidator struct {
	schema *gojsonschema.Schema
}

func NewJsonSchemaValidator(jsonPath string) (Validator, error) {
	bs, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("get abs path failed: %s", err)
	}
	s, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(string(bs)))
	if err != nil {
		return nil, fmt.Errorf("new schema failed: %s", err)
	}
	return &JsonSchemaValidator{
		schema: s,
	}, nil
}

func (v *JsonSchemaValidator) Validate(obj interface{}) error {
	ret, err := v.schema.Validate(gojsonschema.NewGoLoader(obj))
	if err != nil {
		return fmt.Errorf("validate failed: %s", err)
	}

	if !ret.Valid() {
		errString := buffer.Buffer{}
		for i, vErr := range ret.Errors() {
			if i != 0 {
				errString.AppendString("\n")
			}
			errString.AppendString(vErr.String())
		}
		return errors.New(errString.String())
	}
	return nil
}

type APISIXJsonSchemaValidator struct {
	schema    *gojsonschema.Schema
	schemaDef string
}

func NewAPISIXJsonSchemaValidator(jsonPath string) (Validator, error) {
	schemaDef := conf.Schema.Get(jsonPath).String()
	if schemaDef == "" {
		log.Errorf("schema validate failed: schema not found, path: %s", jsonPath)
		return nil, fmt.Errorf("schema validate failed: schema not found, path: %s", jsonPath)
	}

	s, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemaDef))
	if err != nil {
		log.Errorf("new schema failed: %s", err)
		return nil, fmt.Errorf("new schema failed: %s", err)
	}
	return &APISIXJsonSchemaValidator{
		schema:    s,
		schemaDef: schemaDef,
	}, nil
}

func getPlugins(reqBody interface{}) (map[string]interface{}, string) {
	switch bodyType := reqBody.(type) {
	case *entity.Route:
		log.Infof("type of reqBody: %#v", bodyType)
		route := reqBody.(*entity.Route)
		return route.Plugins, "schema"
	case *entity.Service:
		log.Infof("type of reqBody: %#v", bodyType)
		service := reqBody.(*entity.Service)
		return service.Plugins, "schema"
	case *entity.Consumer:
		log.Infof("type of reqBody: %#v", bodyType)
		consumer := reqBody.(*entity.Consumer)
		return consumer.Plugins, "consumer_schema"
	}
	return nil, ""
}

func cHashKeySchemaCheck(upstream *entity.UpstreamDef) error {
	if upstream.HashOn == "consumer" {
		return nil
	}
	if upstream.HashOn != "vars" &&
		upstream.HashOn != "header" &&
		upstream.HashOn != "cookie" {
		return fmt.Errorf("invalid hash_on type: %s", upstream.HashOn)
	}

	var schemaDef string
	if upstream.HashOn == "vars" {
		schemaDef = conf.Schema.Get("main.upstream_hash_vars_schema").String()
		if schemaDef == "" {
			return fmt.Errorf("schema validate failed: schema not found, path: main.upstream_hash_vars_schema")
		}
	}

	if upstream.HashOn == "header" || upstream.HashOn == "cookie" {
		schemaDef = conf.Schema.Get("main.upstream_hash_header_schema").String()
		if schemaDef == "" {
			return fmt.Errorf("schema validate failed: schema not found, path: main.upstream_hash_header_schema")
		}
	}

	s, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemaDef))
	if err != nil {
		return fmt.Errorf("schema validate failed: %s", err)
	}

	ret, err := s.Validate(gojsonschema.NewGoLoader(upstream.Key))
	if err != nil {
		return fmt.Errorf("schema validate failed: %s", err)
	}

	if !ret.Valid() {
		errString := buffer.Buffer{}
		for i, vErr := range ret.Errors() {
			if i != 0 {
				errString.AppendString("\n")
			}
			errString.AppendString(vErr.String())
		}
		return fmt.Errorf("schema validate failed: %s", errString.String())
	}

	return nil
}

func checkUpstream(upstream *entity.UpstreamDef) error {
	if upstream == nil {
		return nil
	}

	if upstream.PassHost == "node" && upstream.Nodes != nil {
		nodes, ok := entity.NodesFormat(upstream.Nodes).([]*entity.Node)
		if !ok {
			return fmt.Errorf("upstrams nodes not support value %v when `pass_host` is `node`", nodes)
		} else if len(nodes) != 1 {
			return fmt.Errorf("only support single node for `node` mode currentlywhen `pass_host` is `node`")
		}
	}

        if upstream.ServiceNameType == "rewrite" && upstream.ServiceName == "" {
                return fmt.Errorf("`service_name` can't be empty when `service_name_type` is `rewrite`")
        }

	if upstream.PassHost == "rewrite" && upstream.UpstreamHost == "" {
		return fmt.Errorf("`upstream_host` can't be empty when `pass_host` is `rewrite`")
	}

	if upstream.Type != "chash" {
		return nil
	}

	//to confirm
	if upstream.HashOn == "" {
		upstream.HashOn = "vars"
	}

	if upstream.HashOn != "consumer" && upstream.Key == "" {
		return fmt.Errorf("missing key")
	}

	if err := cHashKeySchemaCheck(upstream); err != nil {
		return err
	}

	return nil
}

func checkRemoteAddr(remoteAddrs []string) error {
	for _, remoteAddr := range remoteAddrs {
		if remoteAddr == "" {
			return fmt.Errorf("schema validate failed: invalid field remote_addrs")
		}
	}
	return nil
}

func checkConf(reqBody interface{}) error {
	switch bodyType := reqBody.(type) {
	case *entity.Route:
		route := reqBody.(*entity.Route)
		log.Infof("type of reqBody: %#v", bodyType)
		if err := checkUpstream(route.Upstream); err != nil {
			return err
		}
		// todo: this is a temporary method, we'll drop it later
		if err := checkRemoteAddr(route.RemoteAddrs); err != nil {
			return err
		}
	case *entity.Service:
		service := reqBody.(*entity.Service)
		if err := checkUpstream(service.Upstream); err != nil {
			return err
		}
	case *entity.Upstream:
		upstream := reqBody.(*entity.Upstream)
		if err := checkUpstream(&upstream.UpstreamDef); err != nil {
			return err
		}
	}
	return nil
}

func (v *APISIXJsonSchemaValidator) Validate(obj interface{}) error {
	ret, err := v.schema.Validate(gojsonschema.NewGoLoader(obj))
	if err != nil {
		log.Errorf("schema validate failed: %s, s: %v, obj: %v", err, v.schema, obj)
		return fmt.Errorf("schema validate failed: %s", err)
	}

	if !ret.Valid() {
		errString := buffer.Buffer{}
		for i, vErr := range ret.Errors() {
			if i != 0 {
				errString.AppendString("\n")
			}
			errString.AppendString(vErr.String())
		}
		log.Errorf("schema validate failed:s: %v, obj: %#v", v.schemaDef, obj)
		return fmt.Errorf("schema validate failed: %s", errString.String())
	}

	//custom check
	if err := checkConf(obj); err != nil {
		return err
	}

	plugins, schemaType := getPlugins(obj)
	for pluginName, pluginConf := range plugins {
		schemaValue := conf.Schema.Get("plugins." + pluginName + "." + schemaType).Value()
		if schemaValue == nil && schemaType == "consumer_schema" {
			schemaValue = conf.Schema.Get("plugins." + pluginName + ".schema").Value()
		}

		if schemaValue == nil {
			log.Errorf("schema validate failed: schema not found,  %s, %s", "plugins."+pluginName, schemaType)
			return fmt.Errorf("schema validate failed: schema not found, path: %s", "plugins."+pluginName)
		}
		schemaMap := schemaValue.(map[string]interface{})
		schemaByte, err := json.Marshal(schemaMap)
		if err != nil {
			log.Warnf("schema validate failed: schema json encode failed, path: %s, %w", "plugins."+pluginName, err)
			return fmt.Errorf("schema validate failed: schema json encode failed, path: %s, %w", "plugins."+pluginName, err)
		}

		s, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(schemaByte))
		if err != nil {
			log.Errorf("init schema validate failed: %s", err)
			return fmt.Errorf("schema validate failed: %s", err)
		}

		// check property disable, if is bool, remove from json schema checking
		conf := pluginConf.(map[string]interface{})
		var exchange bool
		disable, ok := conf["disable"]
		if ok {
			if fmt.Sprintf("%T", disable) == "bool" {
				delete(conf, "disable")
				exchange = true
			}
		}

		// check schema
		ret, err := s.Validate(gojsonschema.NewGoLoader(conf))
		if err != nil {
			log.Errorf("schema validate failed: %s", err)
			return fmt.Errorf("schema validate failed: %s", err)
		}

		// put the value back to the property disable
		if exchange {
			conf["disable"] = disable
		}

		if !ret.Valid() {
			errString := buffer.Buffer{}
			for i, vErr := range ret.Errors() {
				if i != 0 {
					errString.AppendString("\n")
				}
				errString.AppendString(vErr.String())
			}
			return fmt.Errorf("schema validate failed: %s", errString.String())
		}
	}

	return nil
}

type APISIXSchemaValidator struct {
	schema *gojsonschema.Schema
}

func NewAPISIXSchemaValidator(jsonPath string) (Validator, error) {
	schemaDef := conf.Schema.Get(jsonPath).String()
	if schemaDef == "" {
		log.Warnf("schema validate failed: schema not found, path: %s", jsonPath)
		return nil, fmt.Errorf("schema validate failed: schema not found, path: %s", jsonPath)
	}

	s, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemaDef))
	if err != nil {
		log.Warnf("new schema failed: %w", err)
		return nil, fmt.Errorf("new schema failed: %w", err)
	}
	return &APISIXSchemaValidator{
		schema: s,
	}, nil
}

func (v *APISIXSchemaValidator) Validate(obj interface{}) error {
	ret, err := v.schema.Validate(gojsonschema.NewBytesLoader(obj.([]byte)))
	if err != nil {
		log.Warnf("schema validate failed: %w", err)
		return fmt.Errorf("schema validate failed: %w", err)
	}

	if !ret.Valid() {
		errString := buffer.Buffer{}
		for i, vErr := range ret.Errors() {
			if i != 0 {
				errString.AppendString("\n")
			}
			errString.AppendString(vErr.String())
		}
		return fmt.Errorf("schema validate failed: %s", errString.String())
	}

	return nil
}
