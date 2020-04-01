package infrastructure

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type SYRAuthMock struct {
	mock.Mock
}

func (m *SYRAuthMock) authorize(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *SYRAuthMock) token() string {
	args := m.Called()
	return args.String(0)
}

func (m *SYRAuthMock) client() *http.Client {
	args := m.Called()
	return args.Get(0).(*http.Client)
}

/*
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
*/
