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
	agent "chainsource-gateway/agent"
	"context"
	"errors"
	gomock "github.com/golang/mock/gomock"
	"net/http"
	reflect "reflect"
)

// MockProvider is a mock of Provider interface
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// GetAgentConfigForRepo mocks base method
func (m *MockProvider) GetAgentConfigForRepo(arg0 string) (agent.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAgentConfigForRepo", arg0)
	ret0, _ := ret[0].(agent.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgentConfigForRepo indicates an expected call of GetAgentConfigForRepo
func (mr *MockProviderMockRecorder) GetAgentConfigForRepo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgentConfigForRepo", reflect.TypeOf((*MockProvider)(nil).GetAgentConfigForRepo), arg0)
}

// NewAgent mocks base method
func (m *MockProvider) NewAgent(arg0 *agent.Config) agent.Agent {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewAgent", arg0)
	ret0, _ := ret[0].(agent.Agent)
	return ret0
}

// NewAgent indicates an expected call of NewAgent
func (mr *MockProviderMockRecorder) NewAgent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewAgent", reflect.TypeOf((*MockProvider)(nil).NewAgent), arg0)
}

// InjectAgentProviderIntoRequest injects agentProvider into a request
func InjectAgentProviderIntoRequest(r *http.Request, ctrl *gomock.Controller, repoAgentMap map[string]agent.Agent) (reqWithContext *http.Request) {
	ctx := InjectAgentProviderIntoContext(r.Context(), ctrl, repoAgentMap)
	reqWithContext = r.WithContext(ctx)
	return
}

// injectAgentProviderIntoContext injects agentProvider into a context
func InjectAgentProviderIntoContext(ctx context.Context, ctrl *gomock.Controller, repoAgentMap map[string]agent.Agent) (context.Context) {
	mockProvider := NewMockProvider(ctrl)
	if repoAgentMap != nil {
		for repoID, mappedAgent := range repoAgentMap {
			mockProvider.EXPECT().GetAgentConfigForRepo(gomock.Eq(repoID)).
				Return(agent.Config{RepoID: repoID}, nil).AnyTimes()
			mockProvider.EXPECT().NewAgent(AgentRequestFor(repoID)).
				Return(mappedAgent).AnyTimes()
		}
	}
	mockProvider.EXPECT().GetAgentConfigForRepo(gomock.Eq("NIL_REPO")).
		Return(agent.Config{}, errors.New("")).AnyTimes()
	ctx = context.WithValue(ctx, "agentProvider", mockProvider)
	return ctx
}

