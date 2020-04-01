package interfaces

import (
	"github.com/stretchr/testify/mock"
)

type HooksHandlerMock struct {
	mock.Mock
}

func (m *HooksHandlerMock) Validate(id string, data string) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *HooksHandlerMock) Create(id string, data string, name string) error {
	args := m.Called(id, data, name)
	return args.Error(0)
}
func (m *HooksHandlerMock) Progress(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *HooksHandlerMock) Terminate(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *HooksHandlerMock) Complete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *HooksHandlerMock) GetChanTerm() chan string {
	args := m.Called()
	return args.Get(0).(chan string)
}

type clientHooks struct {
	mock.Mock
}

func (m *clientHooks) GetChanTerm() chan string {
	args := m.Called()
	return args.Get(0).(chan string)
}
