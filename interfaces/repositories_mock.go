package interfaces

import (
	"github.com/stretchr/testify/mock"
)

type InvokeHandlerMock struct {
	mock.Mock
}

func (m *InvokeHandlerMock) Remove(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type DbHandlerMock struct {
	mock.Mock
}

func (m *DbHandlerMock) Create(bucket []byte, key []byte, value []byte) error {
	args := m.Called(bucket, key, value)
	return args.Error(0)
}

func (m *DbHandlerMock) Get(bucket []byte, key []byte) ([]byte, error) {
	args := m.Called(bucket, key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *DbHandlerMock) Delete(bucket []byte, key []byte) error {
	args := m.Called(bucket, key)
	return args.Error(0)
}

func (m *DbHandlerMock) Keys(bucket []byte) ([][]byte, error) {
	args := m.Called(bucket)
	return args.Get(0).([][]byte), args.Error(1)
}
