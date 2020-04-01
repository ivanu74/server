package interfaces

import (
	"log"
	"os"
	"testing"
	"time"

	"b.yadro.com/sys/ch-server/domain"
	"b.yadro.com/sys/ch-server/usecases"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func newDataAgent() dataAgent {
	repo := new(usecases.DataRepositoryMock)
	client := new(usecases.HttpClientMock)
	dataAgent, err := usecases.NewDataAgent(
		repo,
		client)
	if err != nil {
		log.Fatal("Fail to create dataAgent")
	}
	return dataAgent
}
func TestNewHooksHandlerDataAgentNil(t *testing.T) {
	h, err := NewHooksHandler(
		nil,
		new(clientHooks),
		log.New(os.Stdout, "[test] ", log.LstdFlags))
	assert.Nil(t, h)
	assert.NotNil(t, err)
}
func TestNewHooksHandlerClientNil(t *testing.T) {
	h, err := NewHooksHandler(
		newDataAgent(),
		nil,
		log.New(os.Stdout, "[test] ", log.LstdFlags))
	assert.Nil(t, h)
	assert.NotNil(t, err)
}
func TestNewHooksHandlerLoggerNil(t *testing.T) {
	h, err := NewHooksHandler(
		newDataAgent(),
		new(clientHooks),
		nil)
	assert.Nil(t, h)
	assert.NotNil(t, err)
}
func TestNewHooksHandlerValid(t *testing.T) {
	h, err := NewHooksHandler(
		newDataAgent(),
		new(clientHooks),
		log.New(os.Stdout, "[test] ", log.LstdFlags))
	assert.NotNil(t, h)
	assert.Nil(t, err)
}
func TestValidate(t *testing.T) {
	repo := new(usecases.DataRepositoryMock)
	client := new(usecases.HttpClientMock)
	dataAgent, err := usecases.NewDataAgent(
		repo,
		client)
	assert.Nil(t, err)
	hooksHandler, _ := NewHooksHandler(
		dataAgent,
		new(clientHooks),
		log.New(os.Stdout, "[test] ", log.LstdFlags))
	id := "0123456789"
	data := `{
		"SerialNumber"          : "0123456789",
		"LogCollectionTimestamp": "Thu Aug 17 14:00:06 MSK 2019",
		"ClientStartTimestamp"  : "Thu Oct 17 14:00:06 MSK 2019"
		}`
	unique := "0123456789" + "Thu Aug 17 14:00:06 MSK 2019"
	repo.On("FindById", unique).Return(domain.Data{}, errors.New("fail"))
	err = hooksHandler.Validate(id, data)
	repo.AssertCalled(t, "FindById", unique)
	assert.Nil(t, err)
}

func TestCreate(t *testing.T) {
	repo := new(usecases.DataRepositoryMock)
	client := new(usecases.HttpClientMock)
	dataAgent, err := usecases.NewDataAgent(
		repo,
		client)
	assert.Nil(t, err)
	hooksHandler, _ := NewHooksHandler(
		dataAgent,
		new(clientHooks),
		log.New(os.Stdout, "[test] ", log.LstdFlags))
	id := "0123456789"
	data := `{
		"SerialNumber"          : "0123456789",
		"LogCollectionTimestamp": "Thu Aug 17 14:00:06 MSK 2019",
		"ClientStartTimestamp"  : "Thu Oct 17 14:00:06 MSK 2019"
		}`
	name := "logs.tar"
	meta := domain.Data{}
	meta.SerialNumber = id
	meta.SessionID = id
	meta.FileName = name
	meta.LogCollectionTimestamp = "Thu Aug 17 14:00:06 MSK 2019"
	meta.ClientStartTimestamp = "Thu Oct 17 14:00:06 MSK 2019"
	meta.StartTimestamp = time.Now().Format(time.UnixDate)
	repo.On("Store", meta).Return(nil)
	err = hooksHandler.Create(id, data, name)
	repo.AssertCalled(t, "Store", meta)
	assert.Nil(t, err)
}
func TestTerminate(t *testing.T) {
	repo := new(usecases.DataRepositoryMock)
	client := new(usecases.HttpClientMock)
	dataAgent, err := usecases.NewDataAgent(
		repo,
		client)
	assert.Nil(t, err)
	hooksHandler, _ := NewHooksHandler(
		dataAgent,
		new(clientHooks),
		log.New(os.Stdout, "[test] ", log.LstdFlags))
	id := "0123456789"
	repo.On("Remove", id).Return(nil)
	err = hooksHandler.Terminate(id)
	repo.AssertCalled(t, "Remove", id)
	assert.Nil(t, err)
}
func TestComplete(t *testing.T) {
	repo := new(usecases.DataRepositoryMock)
	client := new(usecases.HttpClientMock)
	dataAgent, err := usecases.NewDataAgent(
		repo,
		client)
	assert.Nil(t, err)
	hooksHandler, _ := NewHooksHandler(
		dataAgent,
		new(clientHooks),
		log.New(os.Stdout, "[test] ", log.LstdFlags))
	id := "0123456789"
	name := "logs.tar"
	meta := domain.Data{}
	meta.SerialNumber = id
	meta.SessionID = id
	meta.FileName = name
	meta.LogCollectionTimestamp = "Thu Aug 17 14:00:06 MSK 2019"
	meta.ClientStartTimestamp = "Thu Oct 17 14:00:06 MSK 2019"
	meta.StartTimestamp = time.Now().Format(time.UnixDate)
	meta.FinishTimestamp = time.Now().Format(time.UnixDate)
	repo.On("FindById", id).Return(meta, nil)
	repo.On("Store", meta).Return(nil)
	client.On("Send", id, name, id).Return(nil)
	err = hooksHandler.Complete(id)
	repo.AssertCalled(t, "FindById", id)
	repo.AssertCalled(t, "Store", meta)
	client.AssertCalled(t, "Send", id, name, id)
	assert.Nil(t, err)
}
