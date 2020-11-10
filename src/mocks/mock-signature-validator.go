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
	pgp "chainsource-gateway/pgp"
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSignatureValidator is a mock of SignatureValidator interface
type MockSignatureValidator struct {
	ctrl     *gomock.Controller
	recorder *MockSignatureValidatorMockRecorder
}

// MockSignatureValidatorMockRecorder is the mock recorder for MockSignatureValidator
type MockSignatureValidatorMockRecorder struct {
	mock *MockSignatureValidator
}

// NewMockSignatureValidator creates a new mock instance
func NewMockSignatureValidator(ctrl *gomock.Controller) *MockSignatureValidator {
	mock := &MockSignatureValidator{ctrl: ctrl}
	mock.recorder = &MockSignatureValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSignatureValidator) EXPECT() *MockSignatureValidatorMockRecorder {
	return m.recorder
}

// Validate mocks base method
func (m *MockSignatureValidator) Validate(arg0 context.Context, arg1 pgp.ValidateArgs) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", arg0, arg1)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Validate indicates an expected call of Validate
func (mr *MockSignatureValidatorMockRecorder) Validate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockSignatureValidator)(nil).Validate), arg0, arg1)
}
