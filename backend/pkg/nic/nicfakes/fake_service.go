// Code generated by counterfeiter. DO NOT EDIT.
package nicfakes

import (
	"sync"

	"github.com/datianshi/pxeboot/pkg/model"
	"github.com/datianshi/pxeboot/pkg/nic"
)

type FakeService struct {
	CreateServerStub        func(model.ServerConfig) (model.ServerConfig, error)
	createServerMutex       sync.RWMutex
	createServerArgsForCall []struct {
		arg1 model.ServerConfig
	}
	createServerReturns struct {
		result1 model.ServerConfig
		result2 error
	}
	createServerReturnsOnCall map[int]struct {
		result1 model.ServerConfig
		result2 error
	}
	DeleteAllStub        func() error
	deleteAllMutex       sync.RWMutex
	deleteAllArgsForCall []struct {
	}
	deleteAllReturns struct {
		result1 error
	}
	deleteAllReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteServerStub        func(string) error
	deleteServerMutex       sync.RWMutex
	deleteServerArgsForCall []struct {
		arg1 string
	}
	deleteServerReturns struct {
		result1 error
	}
	deleteServerReturnsOnCall map[int]struct {
		result1 error
	}
	FindServerStub        func(string) (model.ServerConfig, error)
	findServerMutex       sync.RWMutex
	findServerArgsForCall []struct {
		arg1 string
	}
	findServerReturns struct {
		result1 model.ServerConfig
		result2 error
	}
	findServerReturnsOnCall map[int]struct {
		result1 model.ServerConfig
		result2 error
	}
	GetServersStub        func() ([]model.ServerConfig, error)
	getServersMutex       sync.RWMutex
	getServersArgsForCall []struct {
	}
	getServersReturns struct {
		result1 []model.ServerConfig
		result2 error
	}
	getServersReturnsOnCall map[int]struct {
		result1 []model.ServerConfig
		result2 error
	}
	UpdateServerStub        func(model.ServerConfig) error
	updateServerMutex       sync.RWMutex
	updateServerArgsForCall []struct {
		arg1 model.ServerConfig
	}
	updateServerReturns struct {
		result1 error
	}
	updateServerReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeService) CreateServer(arg1 model.ServerConfig) (model.ServerConfig, error) {
	fake.createServerMutex.Lock()
	ret, specificReturn := fake.createServerReturnsOnCall[len(fake.createServerArgsForCall)]
	fake.createServerArgsForCall = append(fake.createServerArgsForCall, struct {
		arg1 model.ServerConfig
	}{arg1})
	stub := fake.CreateServerStub
	fakeReturns := fake.createServerReturns
	fake.recordInvocation("CreateServer", []interface{}{arg1})
	fake.createServerMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeService) CreateServerCallCount() int {
	fake.createServerMutex.RLock()
	defer fake.createServerMutex.RUnlock()
	return len(fake.createServerArgsForCall)
}

func (fake *FakeService) CreateServerCalls(stub func(model.ServerConfig) (model.ServerConfig, error)) {
	fake.createServerMutex.Lock()
	defer fake.createServerMutex.Unlock()
	fake.CreateServerStub = stub
}

func (fake *FakeService) CreateServerArgsForCall(i int) model.ServerConfig {
	fake.createServerMutex.RLock()
	defer fake.createServerMutex.RUnlock()
	argsForCall := fake.createServerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeService) CreateServerReturns(result1 model.ServerConfig, result2 error) {
	fake.createServerMutex.Lock()
	defer fake.createServerMutex.Unlock()
	fake.CreateServerStub = nil
	fake.createServerReturns = struct {
		result1 model.ServerConfig
		result2 error
	}{result1, result2}
}

func (fake *FakeService) CreateServerReturnsOnCall(i int, result1 model.ServerConfig, result2 error) {
	fake.createServerMutex.Lock()
	defer fake.createServerMutex.Unlock()
	fake.CreateServerStub = nil
	if fake.createServerReturnsOnCall == nil {
		fake.createServerReturnsOnCall = make(map[int]struct {
			result1 model.ServerConfig
			result2 error
		})
	}
	fake.createServerReturnsOnCall[i] = struct {
		result1 model.ServerConfig
		result2 error
	}{result1, result2}
}

func (fake *FakeService) DeleteAll() error {
	fake.deleteAllMutex.Lock()
	ret, specificReturn := fake.deleteAllReturnsOnCall[len(fake.deleteAllArgsForCall)]
	fake.deleteAllArgsForCall = append(fake.deleteAllArgsForCall, struct {
	}{})
	stub := fake.DeleteAllStub
	fakeReturns := fake.deleteAllReturns
	fake.recordInvocation("DeleteAll", []interface{}{})
	fake.deleteAllMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeService) DeleteAllCallCount() int {
	fake.deleteAllMutex.RLock()
	defer fake.deleteAllMutex.RUnlock()
	return len(fake.deleteAllArgsForCall)
}

func (fake *FakeService) DeleteAllCalls(stub func() error) {
	fake.deleteAllMutex.Lock()
	defer fake.deleteAllMutex.Unlock()
	fake.DeleteAllStub = stub
}

func (fake *FakeService) DeleteAllReturns(result1 error) {
	fake.deleteAllMutex.Lock()
	defer fake.deleteAllMutex.Unlock()
	fake.DeleteAllStub = nil
	fake.deleteAllReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeService) DeleteAllReturnsOnCall(i int, result1 error) {
	fake.deleteAllMutex.Lock()
	defer fake.deleteAllMutex.Unlock()
	fake.DeleteAllStub = nil
	if fake.deleteAllReturnsOnCall == nil {
		fake.deleteAllReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteAllReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeService) DeleteServer(arg1 string) error {
	fake.deleteServerMutex.Lock()
	ret, specificReturn := fake.deleteServerReturnsOnCall[len(fake.deleteServerArgsForCall)]
	fake.deleteServerArgsForCall = append(fake.deleteServerArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.DeleteServerStub
	fakeReturns := fake.deleteServerReturns
	fake.recordInvocation("DeleteServer", []interface{}{arg1})
	fake.deleteServerMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeService) DeleteServerCallCount() int {
	fake.deleteServerMutex.RLock()
	defer fake.deleteServerMutex.RUnlock()
	return len(fake.deleteServerArgsForCall)
}

func (fake *FakeService) DeleteServerCalls(stub func(string) error) {
	fake.deleteServerMutex.Lock()
	defer fake.deleteServerMutex.Unlock()
	fake.DeleteServerStub = stub
}

func (fake *FakeService) DeleteServerArgsForCall(i int) string {
	fake.deleteServerMutex.RLock()
	defer fake.deleteServerMutex.RUnlock()
	argsForCall := fake.deleteServerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeService) DeleteServerReturns(result1 error) {
	fake.deleteServerMutex.Lock()
	defer fake.deleteServerMutex.Unlock()
	fake.DeleteServerStub = nil
	fake.deleteServerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeService) DeleteServerReturnsOnCall(i int, result1 error) {
	fake.deleteServerMutex.Lock()
	defer fake.deleteServerMutex.Unlock()
	fake.DeleteServerStub = nil
	if fake.deleteServerReturnsOnCall == nil {
		fake.deleteServerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteServerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeService) FindServer(arg1 string) (model.ServerConfig, error) {
	fake.findServerMutex.Lock()
	ret, specificReturn := fake.findServerReturnsOnCall[len(fake.findServerArgsForCall)]
	fake.findServerArgsForCall = append(fake.findServerArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.FindServerStub
	fakeReturns := fake.findServerReturns
	fake.recordInvocation("FindServer", []interface{}{arg1})
	fake.findServerMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeService) FindServerCallCount() int {
	fake.findServerMutex.RLock()
	defer fake.findServerMutex.RUnlock()
	return len(fake.findServerArgsForCall)
}

func (fake *FakeService) FindServerCalls(stub func(string) (model.ServerConfig, error)) {
	fake.findServerMutex.Lock()
	defer fake.findServerMutex.Unlock()
	fake.FindServerStub = stub
}

func (fake *FakeService) FindServerArgsForCall(i int) string {
	fake.findServerMutex.RLock()
	defer fake.findServerMutex.RUnlock()
	argsForCall := fake.findServerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeService) FindServerReturns(result1 model.ServerConfig, result2 error) {
	fake.findServerMutex.Lock()
	defer fake.findServerMutex.Unlock()
	fake.FindServerStub = nil
	fake.findServerReturns = struct {
		result1 model.ServerConfig
		result2 error
	}{result1, result2}
}

func (fake *FakeService) FindServerReturnsOnCall(i int, result1 model.ServerConfig, result2 error) {
	fake.findServerMutex.Lock()
	defer fake.findServerMutex.Unlock()
	fake.FindServerStub = nil
	if fake.findServerReturnsOnCall == nil {
		fake.findServerReturnsOnCall = make(map[int]struct {
			result1 model.ServerConfig
			result2 error
		})
	}
	fake.findServerReturnsOnCall[i] = struct {
		result1 model.ServerConfig
		result2 error
	}{result1, result2}
}

func (fake *FakeService) GetServers() ([]model.ServerConfig, error) {
	fake.getServersMutex.Lock()
	ret, specificReturn := fake.getServersReturnsOnCall[len(fake.getServersArgsForCall)]
	fake.getServersArgsForCall = append(fake.getServersArgsForCall, struct {
	}{})
	stub := fake.GetServersStub
	fakeReturns := fake.getServersReturns
	fake.recordInvocation("GetServers", []interface{}{})
	fake.getServersMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeService) GetServersCallCount() int {
	fake.getServersMutex.RLock()
	defer fake.getServersMutex.RUnlock()
	return len(fake.getServersArgsForCall)
}

func (fake *FakeService) GetServersCalls(stub func() ([]model.ServerConfig, error)) {
	fake.getServersMutex.Lock()
	defer fake.getServersMutex.Unlock()
	fake.GetServersStub = stub
}

func (fake *FakeService) GetServersReturns(result1 []model.ServerConfig, result2 error) {
	fake.getServersMutex.Lock()
	defer fake.getServersMutex.Unlock()
	fake.GetServersStub = nil
	fake.getServersReturns = struct {
		result1 []model.ServerConfig
		result2 error
	}{result1, result2}
}

func (fake *FakeService) GetServersReturnsOnCall(i int, result1 []model.ServerConfig, result2 error) {
	fake.getServersMutex.Lock()
	defer fake.getServersMutex.Unlock()
	fake.GetServersStub = nil
	if fake.getServersReturnsOnCall == nil {
		fake.getServersReturnsOnCall = make(map[int]struct {
			result1 []model.ServerConfig
			result2 error
		})
	}
	fake.getServersReturnsOnCall[i] = struct {
		result1 []model.ServerConfig
		result2 error
	}{result1, result2}
}

func (fake *FakeService) UpdateServer(arg1 model.ServerConfig) error {
	fake.updateServerMutex.Lock()
	ret, specificReturn := fake.updateServerReturnsOnCall[len(fake.updateServerArgsForCall)]
	fake.updateServerArgsForCall = append(fake.updateServerArgsForCall, struct {
		arg1 model.ServerConfig
	}{arg1})
	stub := fake.UpdateServerStub
	fakeReturns := fake.updateServerReturns
	fake.recordInvocation("UpdateServer", []interface{}{arg1})
	fake.updateServerMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeService) UpdateServerCallCount() int {
	fake.updateServerMutex.RLock()
	defer fake.updateServerMutex.RUnlock()
	return len(fake.updateServerArgsForCall)
}

func (fake *FakeService) UpdateServerCalls(stub func(model.ServerConfig) error) {
	fake.updateServerMutex.Lock()
	defer fake.updateServerMutex.Unlock()
	fake.UpdateServerStub = stub
}

func (fake *FakeService) UpdateServerArgsForCall(i int) model.ServerConfig {
	fake.updateServerMutex.RLock()
	defer fake.updateServerMutex.RUnlock()
	argsForCall := fake.updateServerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeService) UpdateServerReturns(result1 error) {
	fake.updateServerMutex.Lock()
	defer fake.updateServerMutex.Unlock()
	fake.UpdateServerStub = nil
	fake.updateServerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeService) UpdateServerReturnsOnCall(i int, result1 error) {
	fake.updateServerMutex.Lock()
	defer fake.updateServerMutex.Unlock()
	fake.UpdateServerStub = nil
	if fake.updateServerReturnsOnCall == nil {
		fake.updateServerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateServerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createServerMutex.RLock()
	defer fake.createServerMutex.RUnlock()
	fake.deleteAllMutex.RLock()
	defer fake.deleteAllMutex.RUnlock()
	fake.deleteServerMutex.RLock()
	defer fake.deleteServerMutex.RUnlock()
	fake.findServerMutex.RLock()
	defer fake.findServerMutex.RUnlock()
	fake.getServersMutex.RLock()
	defer fake.getServersMutex.RUnlock()
	fake.updateServerMutex.RLock()
	defer fake.updateServerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeService) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ nic.Service = new(FakeService)
