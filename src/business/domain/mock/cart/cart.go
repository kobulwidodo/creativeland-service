// Code generated by MockGen. DO NOT EDIT.
// Source: src/business/domain/cart/cart.go

// Package mock_cart is a generated GoMock package.
package mock_cart

import (
	entity "go-clean/src/business/entity"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockInterface) Create(cart entity.Cart) (entity.Cart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", cart)
	ret0, _ := ret[0].(entity.Cart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockInterfaceMockRecorder) Create(cart interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockInterface)(nil).Create), cart)
}

// Delete mocks base method.
func (m *MockInterface) Delete(param entity.CartParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", param)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockInterfaceMockRecorder) Delete(param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockInterface)(nil).Delete), param)
}

// Get mocks base method.
func (m *MockInterface) Get(param entity.CartParam) (entity.Cart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", param)
	ret0, _ := ret[0].(entity.Cart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockInterfaceMockRecorder) Get(param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockInterface)(nil).Get), param)
}

// GetList mocks base method.
func (m *MockInterface) GetList(param entity.CartParam) ([]entity.Cart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", param)
	ret0, _ := ret[0].([]entity.Cart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetList indicates an expected call of GetList.
func (mr *MockInterfaceMockRecorder) GetList(param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockInterface)(nil).GetList), param)
}

// GetListInByID mocks base method.
func (m *MockInterface) GetListInByID(ids []int64) ([]entity.Cart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListInByID", ids)
	ret0, _ := ret[0].([]entity.Cart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListInByID indicates an expected call of GetListInByID.
func (mr *MockInterfaceMockRecorder) GetListInByID(ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListInByID", reflect.TypeOf((*MockInterface)(nil).GetListInByID), ids)
}

// Update mocks base method.
func (m *MockInterface) Update(selectParam entity.CartParam, updateParam entity.UpdateCartParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", selectParam, updateParam)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockInterfaceMockRecorder) Update(selectParam, updateParam interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockInterface)(nil).Update), selectParam, updateParam)
}

// UpdatesByIDs mocks base method.
func (m *MockInterface) UpdatesByIDs(ids []uint, updateParam entity.UpdateCartParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatesByIDs", ids, updateParam)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatesByIDs indicates an expected call of UpdatesByIDs.
func (mr *MockInterfaceMockRecorder) UpdatesByIDs(ids, updateParam interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatesByIDs", reflect.TypeOf((*MockInterface)(nil).UpdatesByIDs), ids, updateParam)
}
