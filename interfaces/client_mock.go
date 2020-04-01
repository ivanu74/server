package interfaces

import (
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) Send(id string, name string, system string) error {
	args := m.Called(id, name, system)
	return args.Error(0)
}
