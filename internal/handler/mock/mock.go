// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go

// Package mock_handler is a generated GoMock package.
package mock_handler

import (
	context "context"
	model "goapi/internal/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, email, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthServiceMockRecorder) Login(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthService)(nil).Login), ctx, email, password)
}

// Register mocks base method.
func (m *MockAuthService) Register(ctx context.Context, email, password string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, email, password)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockAuthServiceMockRecorder) Register(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuthService)(nil).Register), ctx, email, password)
}

// MockProductService is a mock of ProductService interface.
type MockProductService struct {
	ctrl     *gomock.Controller
	recorder *MockProductServiceMockRecorder
}

// MockProductServiceMockRecorder is the mock recorder for MockProductService.
type MockProductServiceMockRecorder struct {
	mock *MockProductService
}

// NewMockProductService creates a new mock instance.
func NewMockProductService(ctrl *gomock.Controller) *MockProductService {
	mock := &MockProductService{ctrl: ctrl}
	mock.recorder = &MockProductServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductService) EXPECT() *MockProductServiceMockRecorder {
	return m.recorder
}

// AddProduct mocks base method.
func (m *MockProductService) AddProduct(ctx context.Context, name string, categoryies []string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProduct", ctx, name, categoryies)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProduct indicates an expected call of AddProduct.
func (mr *MockProductServiceMockRecorder) AddProduct(ctx, name, categoryies interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProduct", reflect.TypeOf((*MockProductService)(nil).AddProduct), ctx, name, categoryies)
}

// DeleteProduct mocks base method.
func (m *MockProductService) DeleteProduct(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProduct", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProduct indicates an expected call of DeleteProduct.
func (mr *MockProductServiceMockRecorder) DeleteProduct(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProduct", reflect.TypeOf((*MockProductService)(nil).DeleteProduct), ctx, id)
}

// EditProductCategory mocks base method.
func (m *MockProductService) EditProductCategory(ctx context.Context, id int64, categoryies []model.Category) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditProductCategory", ctx, id, categoryies)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditProductCategory indicates an expected call of EditProductCategory.
func (mr *MockProductServiceMockRecorder) EditProductCategory(ctx, id, categoryies interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditProductCategory", reflect.TypeOf((*MockProductService)(nil).EditProductCategory), ctx, id, categoryies)
}

// EditProductName mocks base method.
func (m *MockProductService) EditProductName(ctx context.Context, id int64, name string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditProductName", ctx, id, name)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditProductName indicates an expected call of EditProductName.
func (mr *MockProductServiceMockRecorder) EditProductName(ctx, id, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditProductName", reflect.TypeOf((*MockProductService)(nil).EditProductName), ctx, id, name)
}

// GetAllProducts mocks base method.
func (m *MockProductService) GetAllProducts(ctx context.Context, tag string) ([]model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllProducts", ctx, tag)
	ret0, _ := ret[0].([]model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllProducts indicates an expected call of GetAllProducts.
func (mr *MockProductServiceMockRecorder) GetAllProducts(ctx, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllProducts", reflect.TypeOf((*MockProductService)(nil).GetAllProducts), ctx, tag)
}

// GetCategoryProducts mocks base method.
func (m *MockProductService) GetCategoryProducts(ctx context.Context, category string) ([]model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategoryProducts", ctx, category)
	ret0, _ := ret[0].([]model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategoryProducts indicates an expected call of GetCategoryProducts.
func (mr *MockProductServiceMockRecorder) GetCategoryProducts(ctx, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategoryProducts", reflect.TypeOf((*MockProductService)(nil).GetCategoryProducts), ctx, category)
}

// MockCategoryService is a mock of CategoryService interface.
type MockCategoryService struct {
	ctrl     *gomock.Controller
	recorder *MockCategoryServiceMockRecorder
}

// MockCategoryServiceMockRecorder is the mock recorder for MockCategoryService.
type MockCategoryServiceMockRecorder struct {
	mock *MockCategoryService
}

// NewMockCategoryService creates a new mock instance.
func NewMockCategoryService(ctrl *gomock.Controller) *MockCategoryService {
	mock := &MockCategoryService{ctrl: ctrl}
	mock.recorder = &MockCategoryServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCategoryService) EXPECT() *MockCategoryServiceMockRecorder {
	return m.recorder
}

// AddCategory mocks base method.
func (m *MockCategoryService) AddCategory(ctx context.Context, name string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCategory", ctx, name)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCategory indicates an expected call of AddCategory.
func (mr *MockCategoryServiceMockRecorder) AddCategory(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCategory", reflect.TypeOf((*MockCategoryService)(nil).AddCategory), ctx, name)
}

// DeleteCategory mocks base method.
func (m *MockCategoryService) DeleteCategory(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCategory", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCategory indicates an expected call of DeleteCategory.
func (mr *MockCategoryServiceMockRecorder) DeleteCategory(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCategory", reflect.TypeOf((*MockCategoryService)(nil).DeleteCategory), ctx, id)
}

// EditCategory mocks base method.
func (m *MockCategoryService) EditCategory(ctx context.Context, id int64, name string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditCategory", ctx, id, name)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditCategory indicates an expected call of EditCategory.
func (mr *MockCategoryServiceMockRecorder) EditCategory(ctx, id, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditCategory", reflect.TypeOf((*MockCategoryService)(nil).EditCategory), ctx, id, name)
}

// GetAllCategoryies mocks base method.
func (m *MockCategoryService) GetAllCategoryies(ctx context.Context, tag string) ([]model.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCategoryies", ctx, tag)
	ret0, _ := ret[0].([]model.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCategoryies indicates an expected call of GetAllCategoryies.
func (mr *MockCategoryServiceMockRecorder) GetAllCategoryies(ctx, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCategoryies", reflect.TypeOf((*MockCategoryService)(nil).GetAllCategoryies), ctx, tag)
}