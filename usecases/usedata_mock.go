package usecases

import (
	"b.yadro.com/sys/ch-server/domain"
	"github.com/stretchr/testify/mock"
)

type HttpClientMock struct {
	mock.Mock
}

func (m *HttpClientMock) Send(id string, name string, system string) error {
	args := m.Called(id, name, system)
	return args.Error(0)
}

type DataRepositoryMock struct {
	mock.Mock
}

func (m *DataRepositoryMock) Store(data domain.Data) error {
	args := m.Called(data)
	return args.Error(0)
}
func (m *DataRepositoryMock) FindById(id string) (domain.Data, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Data), args.Error(1)
}
func (m *DataRepositoryMock) Remove(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *DataRepositoryMock) ReadAll() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}
