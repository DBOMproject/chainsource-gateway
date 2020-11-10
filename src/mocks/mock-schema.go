/*
 * Copyright 2020 Unisys Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package mocks

import (
	context "context"
	"errors"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAssetSchema is a mock of AssetSchema interface
type MockAssetSchema struct {
	ctrl     *gomock.Controller
	recorder *MockAssetSchemaMockRecorder
}

// MockAssetSchemaMockRecorder is the mock recorder for MockAssetSchema
type MockAssetSchemaMockRecorder struct {
	mock *MockAssetSchema
}

// NewMockAssetSchemaAlwaysValid creates a new mock instance that always returns valid
func NewMockAssetSchemaAlwaysValid(ctrl *gomock.Controller) *MockAssetSchema {
	mock := &MockAssetSchema{ctrl: ctrl}
	mock.recorder = &MockAssetSchemaMockRecorder{mock}
	mock.EXPECT().ValidateAsset(gomock.Any(), gomock.Any()).Return("", true, nil).AnyTimes()
	mock.EXPECT().ValidateDetachSubasset(gomock.Any(), gomock.Any()).Return("", true, nil).AnyTimes()
	mock.EXPECT().ValidateAttachSubasset(gomock.Any(), gomock.Any()).Return("", true, nil).AnyTimes()
	mock.EXPECT().ValidateTransferAsset(gomock.Any(), gomock.Any()).Return("", true, nil).AnyTimes()
	return mock
}

// NewMockAssetSchemaAlwaysInvalid creates a new mock instance that always returns invalid
func NewMockAssetSchemaAlwaysInvalid(ctrl *gomock.Controller) *MockAssetSchema {
	mock := &MockAssetSchema{ctrl: ctrl}
	mock.recorder = &MockAssetSchemaMockRecorder{mock}
	mock.EXPECT().ValidateAsset(gomock.Any(), gomock.Any()).Return("errors", false, nil).AnyTimes()
	mock.EXPECT().ValidateDetachSubasset(gomock.Any(), gomock.Any()).Return("errors", false, nil).AnyTimes()
	mock.EXPECT().ValidateAttachSubasset(gomock.Any(), gomock.Any()).Return("errors", false, nil).AnyTimes()
	mock.EXPECT().ValidateTransferAsset(gomock.Any(), gomock.Any()).Return("errors", false, nil).AnyTimes()
	return mock
}

// NewMockAssetSchemaAlwaysFailure creates a new mock instance that always returns an error
func NewMockAssetSchemaAlwaysFailure(ctrl *gomock.Controller) *MockAssetSchema {
	mock := &MockAssetSchema{ctrl: ctrl}
	mock.recorder = &MockAssetSchemaMockRecorder{mock}
	mock.EXPECT().ValidateAsset(gomock.Any(), gomock.Any()).Return("", false, errors.New("")).AnyTimes()
	mock.EXPECT().ValidateDetachSubasset(gomock.Any(), gomock.Any()).Return("", false,  errors.New("")).AnyTimes()
	mock.EXPECT().ValidateAttachSubasset(gomock.Any(), gomock.Any()).Return("", false,  errors.New("")).AnyTimes()
	mock.EXPECT().ValidateTransferAsset(gomock.Any(), gomock.Any()).Return("", false,  errors.New("")).AnyTimes()
	return mock
}



// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAssetSchema) EXPECT() *MockAssetSchemaMockRecorder {
	return m.recorder
}

// ValidateAsset mocks base method
func (m *MockAssetSchema) ValidateAsset(arg0 context.Context, arg1 map[string]interface{}) (string, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAsset", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ValidateAsset indicates an expected call of ValidateAsset
func (mr *MockAssetSchemaMockRecorder) ValidateAsset(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAsset", reflect.TypeOf((*MockAssetSchema)(nil).ValidateAsset), arg0, arg1)
}

// ValidateAttachSubasset mocks base method
func (m *MockAssetSchema) ValidateAttachSubasset(arg0 context.Context, arg1 map[string]interface{}) (string, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAttachSubasset", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ValidateAttachSubasset indicates an expected call of ValidateAttachSubasset
func (mr *MockAssetSchemaMockRecorder) ValidateAttachSubasset(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAttachSubasset", reflect.TypeOf((*MockAssetSchema)(nil).ValidateAttachSubasset), arg0, arg1)
}

// ValidateDetachSubasset mocks base method
func (m *MockAssetSchema) ValidateDetachSubasset(arg0 context.Context, arg1 map[string]interface{}) (string, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateDetachSubasset", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ValidateDetachSubasset indicates an expected call of ValidateDetachSubasset
func (mr *MockAssetSchemaMockRecorder) ValidateDetachSubasset(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateDetachSubasset", reflect.TypeOf((*MockAssetSchema)(nil).ValidateDetachSubasset), arg0, arg1)
}

// ValidateTransferAsset mocks base method
func (m *MockAssetSchema) ValidateTransferAsset(arg0 context.Context, arg1 map[string]interface{}) (string, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateTransferAsset", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ValidateTransferAsset indicates an expected call of ValidateTransferAsset
func (mr *MockAssetSchemaMockRecorder) ValidateTransferAsset(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateTransferAsset", reflect.TypeOf((*MockAssetSchema)(nil).ValidateTransferAsset), arg0, arg1)
}
