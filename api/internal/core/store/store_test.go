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
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shiningrush/droplet/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/apache/apisix-dashboard/api/internal/core/entity"
	"github.com/apache/apisix-dashboard/api/internal/core/storage"
	"github.com/apache/apisix-dashboard/api/internal/utils"
	storage_api "github.com/apache/apisix-dashboard/api/pkg/storage"
)

func TestNewGenericStore(t *testing.T) {
	dfFunc := func(obj interface{}) string { return "" }
	tests := []struct {
		giveOpt   GenericStoreOption
		giveCache map[string]interface{}
		wantStore *GenericStore
		wantErr   error
	}{
		{
			giveOpt: GenericStoreOption{
				BasePath: "test",
				ObjType:  reflect.TypeOf(GenericStoreOption{}),
				KeyFunc:  dfFunc,
			},
			wantStore: &GenericStore{
				Stg: storage.GenEtcdStorage(),
				opt: GenericStoreOption{
					BasePath: "test",
					ObjType:  reflect.TypeOf(GenericStoreOption{}),
					KeyFunc:  dfFunc,
				},
			},
		},
		{
			giveOpt: GenericStoreOption{
				BasePath: "",
				ObjType:  reflect.TypeOf(GenericStoreOption{}),
				KeyFunc:  dfFunc,
			},
			wantErr: fmt.Errorf("base path can not be empty"),
		},
		{
			giveOpt: GenericStoreOption{
				BasePath: "test",
				ObjType:  reflect.TypeOf(""),
				KeyFunc:  dfFunc,
			},
			wantErr: fmt.Errorf("obj type is invalid"),
		},
		{
			giveOpt: GenericStoreOption{
				BasePath: "test",
				ObjType:  nil,
				KeyFunc:  dfFunc,
			},
			wantErr: fmt.Errorf("object type can not be nil"),
		},
		{
			giveOpt: GenericStoreOption{
				BasePath: "test",
				ObjType:  reflect.TypeOf(GenericStoreOption{}),
				KeyFunc:  nil,
			},
			wantErr: fmt.Errorf("key func can not be nil"),
		},
	}
	for _, tc := range tests {
		s, err := NewGenericStore(tc.giveOpt)
		assert.Equal(t, tc.wantErr, err)
		if err != nil {
			continue
		}
		assert.Equal(t, tc.wantStore.Stg, s.Stg)
		assert.Equal(t, tc.wantStore.cache, s.cache)
		assert.Equal(t, tc.wantStore.opt.BasePath, s.opt.BasePath)
		assert.Equal(t, tc.wantStore.opt.ObjType, s.opt.ObjType)
		assert.Equal(t, reflect.TypeOf(tc.wantStore.opt.KeyFunc), reflect.TypeOf(s.opt.KeyFunc))
	}
}

type TestStruct struct {
	entity.BaseInfo
	Field1 string
	Field2 string
}

func TestGenericStore_Init(t *testing.T) {
	tests := []struct {
		caseDesc        string
		giveStore       *GenericStore
		giveListErr     error
		giveListRet     []storage_api.Keypair
		giveWatchCh     chan storage_api.WatchResponse
		giveResp        storage_api.WatchResponse
		wantErr         error
		wantCache       map[string]interface{}
		wantListCalled  bool
		wantWatchCalled bool
	}{
		{
			caseDesc: "sanity",
			giveStore: &GenericStore{
				opt: GenericStoreOption{
					BasePath: "test",
					ObjType:  reflect.TypeOf(TestStruct{}),
					KeyFunc: func(obj interface{}) string {
						return obj.(*TestStruct).Field1
					},
				},
			},
			giveListRet: []storage_api.Keypair{
				{
					Key:   "test/demo1-f1",
					Value: `{"Field1":"demo1-f1", "Field2":"demo1-f2"}`,
				},
				{
					Key:   "test/demo2-f1",
					Value: `{"Field1":"demo2-f1", "Field2":"demo2-f2"}`,
				},
			},
			giveWatchCh: make(chan storage_api.WatchResponse),
			giveResp: storage_api.WatchResponse{
				Events: []storage_api.Event{
					{
						Keypair: storage_api.Keypair{
							Key:   "test/demo3-f1",
							Value: `{"Field1":"demo3-f1", "Field2":"demo3-f2"}`,
						},
						Type: storage_api.EventTypePut,
					},
					{
						Type: storage_api.EventTypeDelete,
						Keypair: storage_api.Keypair{
							Key: "test/demo1-f1",
						},
					},
				},
			},
			wantCache: map[string]interface{}{
				"demo2-f1": &TestStruct{
					BaseInfo: entity.BaseInfo{ID: "demo2-f1"},
					Field1:   "demo2-f1",
					Field2:   "demo2-f2",
				},
				"demo3-f1": &TestStruct{
					BaseInfo: entity.BaseInfo{ID: "demo3-f1"},
					Field1:   "demo3-f1",
					Field2:   "demo3-f2",
				},
			},
			wantListCalled:  true,
			wantWatchCalled: true,
		},
		{
			caseDesc:       "list error",
			giveStore:      &GenericStore{},
			giveListErr:    fmt.Errorf("list error"),
			wantErr:        fmt.Errorf("list error"),
			wantListCalled: true,
		},
		{
			caseDesc: "json error",
			giveStore: &GenericStore{
				opt: GenericStoreOption{
					BasePath: "test",
					ObjType:  reflect.TypeOf(TestStruct{}),
					KeyFunc: func(obj interface{}) string {
						return obj.(*TestStruct).Field1
					},
				},
			},
			giveListRet: []storage_api.Keypair{
				{
					Key:   "test/demo1-f1",
					Value: `{"Field1","demo1-f1", "Field2":"demo1-f2"}`,
				},
				{
					Key:   "test/demo2-f1",
					Value: `{"Field1":"demo2-f1", "Field2":"demo2-f2"}`,
				},
			},
			wantErr: fmt.Errorf("json unmarshal failed\n\tRelated Key:\t\tdemo1-f1\n\tError Description:\t" +
				"invalid character ',' after object key"),
			wantListCalled: true,
		},
	}

	for _, tc := range tests {
		listCalled, watchCalled := false, false
		mStorage := &storage.MockInterface{}
		mStorage.On("List", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			listCalled = true
			assert.Equal(t, tc.giveStore.opt.BasePath, args[1], tc.caseDesc)
		}).Return(tc.giveListRet, tc.giveListErr)
		mStorage.On("Watch", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			watchCalled = true
			assert.Equal(t, tc.giveStore.opt.BasePath, args[1], tc.caseDesc)
		}).Return((<-chan storage_api.WatchResponse)(tc.giveWatchCh))

		tc.giveStore.Stg = mStorage
		err := tc.giveStore.Init()
		assert.Equal(t, tc.wantListCalled, listCalled, tc.caseDesc)
		assert.Equal(t, tc.wantWatchCalled, watchCalled, tc.caseDesc)
		if err != nil {
			assert.Equal(t, tc.wantErr.Error(), err.Error(), tc.caseDesc)
			continue
		}
		tc.giveWatchCh <- tc.giveResp
		time.Sleep(1 * time.Second)
		close(tc.giveWatchCh)
		tc.giveStore.cache.Range(func(key, value interface{}) bool {
			assert.Equal(t, tc.wantCache[key.(string)], value)
			return true
		})
	}
}

func TestGenericStore_Get(t *testing.T) {
	tests := []struct {
		caseDesc  string
		giveId    string
		giveStore *GenericStore
		giveCache map[string]interface{}
		wantRet   interface{}
		wantErr   error
	}{
		{
			caseDesc: "sanity",
			giveId:   "test1",
			giveCache: map[string]interface{}{
				"test2": TestStruct{
					Field1: "test2-f1",
					Field2: "test2-f2",
				},
				"test1": TestStruct{
					Field1: "test1-f1",
					Field2: "test1-f2",
				},
			},
			giveStore: &GenericStore{},
			wantRet: TestStruct{
				Field1: "test1-f1",
				Field2: "test1-f2",
			},
		},
		{
			caseDesc: "not found",
			giveId:   "not",
			giveCache: map[string]interface{}{
				"test2": TestStruct{
					Field1: "test2-f1",
					Field2: "test2-f2",
				},
				"test1": TestStruct{
					Field1: "test1-f1",
					Field2: "test1-f2",
				},
			},
			giveStore: &GenericStore{},
			wantErr:   data.ErrNotFound,
		},
	}

	for _, tc := range tests {
		for k, v := range tc.giveCache {
			tc.giveStore.cache.Store(k, v)
		}
		ret, err := tc.giveStore.Get(context.Background(), tc.giveId)
		assert.Equal(t, tc.wantRet, ret, tc.caseDesc)
		assert.Equal(t, tc.wantErr, err, tc.caseDesc)
	}
}

func TestGenericStore_List(t *testing.T) {
	tests := []struct {
		caseDesc  string
		giveInput ListInput
		giveStore *GenericStore
		giveCache map[string]interface{}
		wantRet   *ListOutput
		wantErr   error
	}{
		{
			caseDesc: "sanity",
			giveInput: ListInput{
				Predicate: func(obj interface{}) bool {
					for _, v := range strings.Split("test1-f2,test3-f2", ",") {
						if v == obj.(*TestStruct).Field2 {
							return true
						}
					}
					return false
				},
			},
			giveCache: map[string]interface{}{
				"test1": &TestStruct{
					Field1: "test1-f1",
					Field2: "test1-f2",
				},
				"test2": &TestStruct{
					Field1: "test2-f1",
					Field2: "test2-f2",
				},
				"test3": &TestStruct{
					Field1: "test3-f1",
					Field2: "test3-f2",
				},
				"test4": &TestStruct{
					Field1: "test4-f1",
					Field2: "test4-f2",
				},
			},
			giveStore: &GenericStore{},
			wantRet: &ListOutput{
				Rows: []interface{}{
					&TestStruct{
						Field1: "test1-f1",
						Field2: "test1-f2",
					},
					&TestStruct{
						Field1: "test3-f1",
						Field2: "test3-f2",
					},
				},
				TotalSize: 2,
			},
		},
		{
			caseDesc: "sanity-page",
			giveInput: ListInput{
				Predicate: func(obj interface{}) bool {
					for _, v := range strings.Split("test1-f2,test3-f2", ",") {
						if v == obj.(*TestStruct).Field2 {
							return true
						}
					}
					return false
				},
				PageSize:   1,
				PageNumber: 2,
			},
			giveCache: map[string]interface{}{
				"test1": &TestStruct{
					Field1: "test1-f1",
					Field2: "test1-f2",
				},
				"test2": &TestStruct{
					Field1: "test2-f1",
					Field2: "test2-f2",
				},
				"test3": &TestStruct{
					BaseInfo: entity.BaseInfo{
						CreateTime: 100,
					},
					Field1: "test3-f1",
					Field2: "test3-f2",
				},
				"test4": &TestStruct{
					Field1: "test4-f1",
					Field2: "test4-f2",
				},
			},
			giveStore: &GenericStore{},
			wantRet: &ListOutput{
				Rows: []interface{}{
					&TestStruct{
						BaseInfo: entity.BaseInfo{
							CreateTime: 100,
						},
						Field1: "test3-f1",
						Field2: "test3-f2",
					},
				},
				TotalSize: 2,
			},
		},
		{
			caseDesc: "page overflow",
			giveInput: ListInput{
				Predicate: func(obj interface{}) bool {
					for _, v := range strings.Split("test1-f2,test3-f2", ",") {
						if v == obj.(*TestStruct).Field2 {
							return true
						}
					}
					return false
				},
				PageSize:   1,
				PageNumber: 33,
			},
			giveCache: map[string]interface{}{
				"test1": &TestStruct{
					Field1: "test1-f1",
					Field2: "test1-f2",
				},
				"test2": &TestStruct{
					Field1: "test2-f1",
					Field2: "test2-f2",
				},
				"test3": &TestStruct{
					Field1: "test3-f1",
					Field2: "test3-f2",
				},
				"test4": &TestStruct{
					Field1: "test4-f1",
					Field2: "test4-f2",
				},
			},
			giveStore: &GenericStore{},
			wantRet: &ListOutput{
				Rows:      []interface{}{},
				TotalSize: 2,
			},
		},
	}

	for _, tc := range tests {
		for k, v := range tc.giveCache {
			tc.giveStore.cache.Store(k, v)
		}

		ret, err := tc.giveStore.List(context.Background(), tc.giveInput)
		assert.Equal(t, tc.wantRet.TotalSize, ret.TotalSize, tc.caseDesc)
		assert.ElementsMatch(t, tc.wantRet.Rows, ret.Rows, tc.caseDesc)
		assert.Equal(t, tc.wantErr, err, tc.caseDesc)
	}
}

func TestGenericStore_ingestValidate(t *testing.T) {
	tests := []struct {
		giveStore       *GenericStore
		giveCache       map[string]interface{}
		giveObj         interface{}
		giveStockCheck  func(obj interface{}, stockObj interface{}) error
		giveValidateErr error
		wantErr         error
	}{
		{
			giveStore: &GenericStore{},
			giveCache: map[string]interface{}{
				"test1-f1": &TestStruct{Field1: "test1-f1", Field2: "test1-f2"},
				"test2-f1": &TestStruct{Field1: "test2-f1", Field2: "test2-f2"},
			},
			giveObj: &TestStruct{
				Field1: "test3-f1",
				Field2: "test2-f2",
			},
			giveStockCheck: func(obj interface{}, stockObj interface{}) error {
				if obj.(*TestStruct).Field2 == stockObj.(*TestStruct).Field2 {
					return fmt.Errorf("field2: %s is conflicted", obj.(*TestStruct).Field2)
				}
				return nil
			},
			wantErr: fmt.Errorf("field2: test2-f2 is conflicted"),
		},
		{
			giveStore:       &GenericStore{},
			giveObj:         &TestStruct{},
			giveValidateErr: fmt.Errorf("validate failed"),
			wantErr:         fmt.Errorf("validate failed"),
		},
	}

	for _, tc := range tests {
		for k, v := range tc.giveCache {
			tc.giveStore.cache.Store(k, v)
		}

		validateCalled := false
		mValidator := &MockValidator{}
		mValidator.On("Validate", mock.Anything).Run(func(args mock.Arguments) {
			validateCalled = true
			assert.Equal(t, tc.giveObj, args.Get(0))
		}).Return(tc.giveValidateErr)

		tc.giveStore.opt.Validator = mValidator
		tc.giveStore.opt.StockCheck = tc.giveStockCheck
		err := tc.giveStore.ingestValidate(tc.giveObj)
		assert.True(t, validateCalled)
		assert.Equal(t, tc.wantErr, err)
	}
}

func TestGenericStore_Create(t *testing.T) {
	tests := []struct {
		caseDesc        string
		giveStore       *GenericStore
		giveCache       map[string]interface{}
		giveObj         *TestStruct
		giveErr         error
		giveValidateErr error
		wantKey         string
		wantErr         error
	}{
		{
			caseDesc: "sanity",
			giveObj: &TestStruct{
				Field1: "test1",
				Field2: "test2",
			},
			giveStore: &GenericStore{
				opt: GenericStoreOption{
					BasePath: "test/path",
					KeyFunc: func(obj interface{}) string {
						return obj.(*TestStruct).Field1
					},
				},
			},
			wantKey: "test/path/test1",
		},
		{
			caseDesc: "create failed",
			giveObj: &TestStruct{
				Field1: "test1",
				Field2: "test2",
			},
			giveStore: &GenericStore{
				opt: GenericStoreOption{
					BasePath: "test/path",
					KeyFunc: func(obj interface{}) string {
						return obj.(*TestStruct).Field1
					},
				},
			},
			giveErr: fmt.Errorf("create failed"),
			wantKey: "test/path/test1",
			wantErr: fmt.Errorf("create failed"),
		},
		{
			caseDesc: "conflicted",
			giveObj: &TestStruct{
				Field1: "test1",
				Field2: "test2",
			},
			giveCache: map[string]interface{}{
				"test1": struct{}{},
			},
			giveStore: &GenericStore{
				opt: GenericStoreOption{
					BasePath: "test/path",
					KeyFunc: func(obj interface{}) string {
						return obj.(*TestStruct).Field1
					},
				},
			},
			wantErr: fmt.Errorf("key: test1 is conflicted"),
		},
		{
			caseDesc: "validate failed",
			giveObj: &TestStruct{
				Field1: "test1",
				Field2: "test2",
			},
			giveCache: map[string]interface{}{
				"test1": struct{}{},
			},
			giveStore: &GenericStore{
				opt: GenericStoreOption{},
			},
			giveValidateErr: fmt.Errorf("validate failed"),
			wantErr:         fmt.Errorf("validate failed"),
		},
	}

	for _, tc := range tests {
		for k, v := range tc.giveCache {
			tc.giveStore.cache.Store(k, v)
		}

		createCalled, validateCalled := false, false
		mStorage := &storage.MockInterface{}
		mStorage.On("Create", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			createCalled = true
			assert.Equal(t, tc.wantKey, args[1], tc.caseDesc)
			input := TestStruct{}
			err := json.Unmarshal([]byte(args[2].(string)), &input)
			id := utils.InterfaceToString(input.ID)
			assert.Nil(t, err)
			assert.Equal(t, tc.giveObj.Field1, input.Field1, tc.caseDesc)
			assert.Equal(t, tc.giveObj.Field2, input.Field2, tc.caseDesc)
			assert.NotEqual(t, 0, len(id), tc.caseDesc)
			assert.NotEqual(t, 0, input.CreateTime, tc.caseDesc)
			assert.NotEqual(t, 0, input.UpdateTime, tc.caseDesc)
		}).Return(tc.giveErr)

		mValidator := &MockValidator{}
		mValidator.On("Validate", mock.Anything).Run(func(args mock.Arguments) {
			validateCalled = true
			assert.Equal(t, tc.giveObj, args.Get(0), tc.caseDesc)
		}).Return(tc.giveValidateErr)

		tc.giveStore.Stg = mStorage
		tc.giveStore.opt.Validator = mValidator
		ret, err := tc.giveStore.Create(context.TODO(), tc.giveObj)
		assert.True(t, validateCalled, tc.caseDesc)
		if err != nil {
			assert.Equal(t, tc.wantErr, err, tc.caseDesc)
			continue
		}
		retTs, ok := ret.(*TestStruct)
		assert.True(t, ok)
		// The returned value (retTs) should be the same as the input (tc.giveObj)
		assert.Equal(t, tc.giveObj.Field1, retTs.Field1, tc.caseDesc)
		assert.Equal(t, tc.giveObj.Field2, retTs.Field2, tc.caseDesc)
		assert.True(t, createCalled, tc.caseDesc)
	}
}

func TestGenericStore_Update(t *testing.T) {
	tests := []struct {
		caseDesc        string
		giveStore       *GenericStore
		giveCache       map[string]interface{}
		giveObj         *TestStruct
		giveErr         error
		giveValidateErr error
		wantKey         string
		wantErr         error
	}{
		{
			caseDesc: "sanity",
			giveObj: &TestStruct{
				Field1: "test1",
				Field2: "test2",
			},
			giveCache: map[string]interface{}{
				"test1": &TestStruct{},
			},
			giveStore: &GenericStore{
				opt: GenericStoreOption{
					BasePath: "test/path",
					KeyFunc: func(obj interface{}) string {
						return obj.(*TestStruct).Field1
					},
				},
			},
			wantKey: "test/path/test1",
		},
		{
			caseDesc: "create failed",
			giveObj: &TestStruct{
				Field1: "test1",
				Field2: "test2",
			},
			giveCache: map[string]interface{}{
				"test1": &TestStruct{},
			},
			giveStore: &GenericStore{
				opt: GenericStoreOption{
					BasePath: "test/path",
					KeyFunc: func(obj interface{}) string {
						return obj.(*TestStruct).Field1
					},
				},
			},
			giveErr: fmt.Errorf("create failed"),
			wantKey: "test/path/test1",
			wantErr: fmt.Errorf("create failed"),
		},
		{
			caseDesc: "not found",
			giveObj: &TestStruct{
				Field1: "test1",
				Field2: "test2",
			},
			giveCache: map[string]interface{}{
				"test2": &TestStruct{},
			},
			giveStore: &GenericStore{
				opt: GenericStoreOption{
					BasePath: "test/path",
					KeyFunc: func(obj interface{}) string {
						return obj.(*TestStruct).Field1
					},
				},
			},
			wantErr: fmt.Errorf("key: test1 is not found"),
		},
	}

	for _, tc := range tests {
		for k, v := range tc.giveCache {
			tc.giveStore.cache.Store(k, v)
		}

		createCalled, validateCalled := false, false
		mStorage := &storage.MockInterface{}
		mStorage.On("Update", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			createCalled = true
			assert.Equal(t, tc.wantKey, args[1], tc.caseDesc)
			input := TestStruct{}
			err := json.Unmarshal([]byte(args[2].(string)), &input)
			assert.Nil(t, err)
			assert.Equal(t, tc.giveObj.Field1, input.Field1, tc.caseDesc)
			assert.Equal(t, tc.giveObj.Field2, input.Field2, tc.caseDesc)
			assert.NotEqual(t, 0, input.UpdateTime, tc.caseDesc)
		}).Return(tc.giveErr)

		mValidator := &MockValidator{}
		mValidator.On("Validate", mock.Anything).Run(func(args mock.Arguments) {
			validateCalled = true
			assert.Equal(t, tc.giveObj, args.Get(0), tc.caseDesc)
		}).Return(tc.giveValidateErr)

		tc.giveStore.Stg = mStorage
		tc.giveStore.opt.Validator = mValidator

		ret, err := tc.giveStore.Update(context.TODO(), tc.giveObj, false)
		assert.True(t, validateCalled, tc.caseDesc)
		if err != nil {
			assert.Equal(t, tc.wantErr, err, tc.caseDesc)
			continue
		}
		retTs, ok := ret.(*TestStruct)
		assert.True(t, ok)
		// The returned value (retTs) should be the same as the input (tc.giveObj)
		assert.Equal(t, tc.giveObj.Field1, retTs.Field1, tc.caseDesc)
		assert.Equal(t, tc.giveObj.Field2, retTs.Field2, tc.caseDesc)
		assert.True(t, createCalled, tc.caseDesc)
	}
}

func TestGenericStore_Delete(t *testing.T) {
	tests := []struct {
		caseDesc  string
		giveStore *GenericStore
		giveKeys  []string
		giveErr   error
		wantKey   []string
		wantErr   error
	}{
		{
			caseDesc: "sanity",
			giveKeys: []string{"test1", "test2"},
			giveStore: &GenericStore{
				opt: GenericStoreOption{
					BasePath: "test/path",
				},
			},
			wantKey: []string{"test/path/test1", "test/path/test2"},
		},
		{
			caseDesc: "delete failed",
			giveKeys: []string{"test1", "test2"},
			giveStore: &GenericStore{
				opt: GenericStoreOption{
					BasePath: "test/path",
				},
			},
			wantKey: []string{"test/path/test1", "test/path/test2"},
			giveErr: fmt.Errorf("delete failed"),
			wantErr: fmt.Errorf("delete failed"),
		},
	}

	for _, tc := range tests {
		createCalled := false
		mStorage := &storage.MockInterface{}
		mStorage.On("BatchDelete", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			createCalled = true
			assert.Equal(t, tc.wantKey, args[1], tc.caseDesc)
		}).Return(tc.giveErr)

		tc.giveStore.Stg = mStorage
		err := tc.giveStore.BatchDelete(context.TODO(), tc.giveKeys)
		assert.True(t, createCalled, tc.caseDesc)
		assert.Equal(t, tc.wantErr, err, tc.caseDesc)
	}
}

func TestGenericStore_StringToObjPtr(t *testing.T) {
	s, err := NewGenericStore(GenericStoreOption{
		BasePath: "test",
		ObjType:  reflect.TypeOf(entity.SSL{}),
		KeyFunc: func(obj interface{}) string {
			r := obj.(*entity.Route)
			return utils.InterfaceToString(r.ID)
		},
	})
	assert.Nil(t, err)
	id := "1"
	sslStr := `{"key":"test_key", "cert":"test_cert"}`
	sslInterface, err := s.StringToObjPtr(sslStr, id)
	ssl := sslInterface.(*entity.SSL)
	assert.Equal(t, id, ssl.ID)
}
