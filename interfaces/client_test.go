package interfaces

import (
	"log"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPClientClientHandlerNil(t *testing.T) {
	h, err := NewHTTPClient(
		nil,
		log.New(os.Stdout, "[test] ", log.LstdFlags))
	assert.Nil(t, h)
	assert.NotNil(t, err)
}

func TestNewHTTPClientLoggerNil(t *testing.T) {
	h, err := NewHTTPClient(
		new(ClientMock),
		nil)
	assert.Nil(t, h)
	assert.NotNil(t, err)
}
func TestNewHTTPClientValid(t *testing.T) {
	h, err := NewHTTPClient(
		new(ClientMock),
		log.New(os.Stdout, "[test] ", log.LstdFlags))
	assert.NotNil(t, h)
	assert.Nil(t, err)
}

func TestSendInvalid(t *testing.T) {
	client := new(ClientMock)
	httpClient, _ := NewHTTPClient(
		client,
		log.New(os.Stdout, "[test] ", log.LstdFlags))
	assert.NotNil(t, httpClient)
	id := "0123456789"
	name := "logs.tar"
	system := "0123456789"
	e := errors.New("failed")
	client.On("Send", id, name, system).Return(e)
	err := httpClient.Send(id, name, system)
	client.AssertCalled(t, "Send", id, name, system)
	assert.Equal(t, errors.Cause(err), e)
}

func TestSendValid(t *testing.T) {
	client := new(ClientMock)
	httpClient, _ := NewHTTPClient(
		client,
		log.New(os.Stdout, "[test] ", log.LstdFlags))
	assert.NotNil(t, httpClient)
	id := "0123456789"
	name := "logs.tar"
	system := "0123456789"
	client.On("Send", id, name, system).Return(nil)
	err := httpClient.Send(id, name, system)
	client.AssertCalled(t, "Send", id, name, system)
	assert.Nil(t, err)
}
