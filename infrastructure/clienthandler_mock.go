package infrastructure

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type HooksClientHandlerMock struct {
	mock.Mock
}

func (m *HooksClientHandlerMock) Terminate(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MultipartfileMock struct {
	mock.Mock
}

func (m *MultipartfileMock) uploadMultipartFile(
	ctx context.Context,
	client *http.Client,
	url,
	token,
	key,
	path string,
	name string) (*http.Response, error) {
	args := m.Called(ctx, client, url, token, key, path, name)
	return args.Get(0).(*http.Response), args.Error(1)
}
