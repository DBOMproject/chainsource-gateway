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

// Package mocks contains the mock implementations generated by mockgen for gomock to facilitate unit testing
package mocks

import (
	agent "chainsource-gateway/agent"
	context "context"
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAgent is a mocks of Agent interface
type MockAgent struct {
	ctrl     *gomock.Controller
	recorder *MockAgentMockRecorder
}

// MockAgentMockRecorder is the mocks recorder for MockAgent
type MockAgentMockRecorder struct {
	mock *MockAgent
}

// NewMockAgent creates a new mocks instance
func NewMockAgent(ctrl *gomock.Controller) *MockAgent {
	mock := &MockAgent{ctrl: ctrl}
	mock.recorder = &MockAgentMockRecorder{mock}
	mock.EXPECT().GetPort().Return(3000).AnyTimes()
	mock.EXPECT().GetHost().Return("TestHost").AnyTimes()
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAgent) EXPECT() *MockAgentMockRecorder {
	return m.recorder
}

// Commit mocks base method
func (m *MockAgent) Commit(arg0 context.Context, arg1 agent.CommitArgs) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", arg0, arg1)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Commit indicates an expected call of Commit
func (mr *MockAgentMockRecorder) Commit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockAgent)(nil).Commit), arg0, arg1)
}

// GetHost mocks base method
func (m *MockAgent) GetHost() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHost")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetHost indicates an expected call of GetHost
func (mr *MockAgentMockRecorder) GetHost() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHost", reflect.TypeOf((*MockAgent)(nil).GetHost))
}

// GetPort mocks base method
func (m *MockAgent) GetPort() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPort")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetPort indicates an expected call of GetPort
func (mr *MockAgentMockRecorder) GetPort() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPort", reflect.TypeOf((*MockAgent)(nil).GetPort))
}

// QueryAuditTrail mocks base method
func (m *MockAgent) QueryAuditTrail(arg0 context.Context, arg1 agent.QueryArgs) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAuditTrail", arg0, arg1)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAuditTrail indicates an expected call of QueryAuditTrail
func (mr *MockAgentMockRecorder) QueryAuditTrail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAuditTrail", reflect.TypeOf((*MockAgent)(nil).QueryAuditTrail), arg0, arg1)
}

// QueryStream mocks base method
func (m *MockAgent) QueryStream(arg0 context.Context, arg1 agent.QueryArgs) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryStream", arg0, arg1)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryStream indicates an expected call of QueryStream
func (mr *MockAgentMockRecorder) QueryStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryStream", reflect.TypeOf((*MockAgent)(nil).QueryStream), arg0, arg1)
}

// QueryAssets mocks base method
func (m *MockAgent) QueryAssets(arg0 context.Context, arg1 agent.QueryArgs, arg2 agent.RichQueryArgs) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAssets", arg0, arg1, arg2)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAssets indicates an expected call of QueryAssets
func (mr *MockAgentMockRecorder) QueryAssets(arg0, arg1 interface{}, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAssets", reflect.TypeOf((*MockAgent)(nil).QueryAssets), arg0, arg1, arg2)
}

// ListChannels mocks base method
func (m *MockAgent) ListChannels(arg0 context.Context, arg1 agent.QueryArgs) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChannels", arg0, arg1)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListChannels indicates an expected call of ListChannels
func (mr *MockAgentMockRecorder) ListChannels(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChannels", reflect.TypeOf((*MockAgent)(nil).QueryStream), arg0, arg1)
}

// ListAssets mocks base method
func (m *MockAgent) ListAssets(arg0 context.Context, arg1 agent.QueryArgs) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAssets", arg0, arg1)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAssets indicates an expected call of ListAssets
func (mr *MockAgentMockRecorder) ListAssets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAssets", reflect.TypeOf((*MockAgent)(nil).ListAssets), arg0, arg1)
}
