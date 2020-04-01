package usecases

import (
	"testing"

	"b.yadro.com/sys/ch-server/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/pkg/errors"

)
func TestNotUnique(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	meta := Data{}
	repo.On("FindById", meta.SerialNumber).Return(domain.Data{}, nil)
	err := dataAgent.IsUnique(meta)
	assert.NotNil(t, err, "Need to be not unique. FindById is success.")
}
func TestUnique(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	meta := Data{}
	repo.On("FindById", meta.SerialNumber).Return(domain.Data{}, errors.New("not exists"))
	err := dataAgent.IsUnique(meta)
	assert.Nil(t, err, "Need to be unique. FindById is not success.")
}

func TestCreateValidate(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	
	meta := Data{}
	meta.SessionID = "0123456789"
	meta.SerialNumber = "0123456789"
	meta.FileName = ""
	repo.On("Store", mock.Anything).Return(nil)
	err := dataAgent.Create(meta)
	assert.Error(t, err)
	repo.AssertNotCalled(t, "Store", mock.Anything)
}
func TestCreateStoreInvalid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	
	meta := Data{}
	meta.SessionID = "0123456789"
	meta.SerialNumber = "0123456789"
	meta.FileName = "log.tar"
	repo.On("Store", mock.Anything).Return(errors.New("failed to store data"))
	err := dataAgent.Create(meta)
	repo.AssertCalled(t, "Store", mock.Anything)
	assert.Error(t, err)
}
func TestCreateStoreValid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	
	meta := Data{}
	meta.SessionID = "0123456789"
	meta.SerialNumber = "0123456789"
	meta.FileName = "log.tar"
	repo.On("Store", mock.Anything).Return(nil)
	err := dataAgent.Create(meta)
	repo.AssertCalled(t, "Store", mock.Anything)
	assert.Nil(t, err)
}

func TestReadInvalid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	id := "0123456789"
	repo.On("FindById", id).Return(domain.Data{}, errors.New("failed"))
	_, err := dataAgent.Read(id)
	assert.NotNil(t, err, "FindById is not success.")
}

func TestReadValid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	id := "0123456789"
	data := domain.Data{}
	data.SessionID = id
	repo.On("FindById", id).Return(data, nil)
	meta, err := dataAgent.Read(id)
	assert.Nil(t, err, "FindById is success.")
	assert.Equal(t, meta.SessionID, id)
}

func TestUpdateFindInvalid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	id := "0123456789"
	e := errors.New("failed")
	repo.On("FindById", id).Return(domain.Data{}, e)
	err := dataAgent.Update(id, Data{})
	repo.AssertCalled(t, "FindById", id)
	assert.Equal(t, errors.Cause(err), e)
}
func TestUpdateValidateInvalid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	id := "0123456789"
	repo.On("FindById", id).Return(domain.Data{}, nil)
	err := dataAgent.Update(id, Data{})
	repo.AssertCalled(t, "FindById", id)
	assert.NotNil(t, err, "Validate is not success.")
	repo.AssertNotCalled(t, "Store", mock.Anything)
}
func TestUpdateStoreInvalid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	id := "0123456789"
	meta := Data{}
	meta.SessionID = id
	meta.SerialNumber = id
	meta.FileName = "log.tar"
	repo.On("FindById", id).Return(domain.Data{}, nil)
	e := errors.New("failed")
	repo.On("Store", mock.Anything).Return(e)
	err := dataAgent.Update(id, meta)
	repo.AssertCalled(t, "FindById", id)
	repo.AssertCalled(t, "Store", mock.Anything)
	assert.Equal(t, errors.Cause(err), e)
}
func TestUpdateValid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	id := "0123456789"
	meta := Data{}
	meta.SessionID = id
	meta.SerialNumber = id
	meta.FileName = "log.tar"
	repo.On("FindById", id).Return(domain.Data{}, nil)
	repo.On("Store", mock.Anything).Return(nil)
	err := dataAgent.Update(id, meta)
	repo.AssertCalled(t, "FindById", id)
	repo.AssertCalled(t, "Store", mock.Anything)
	assert.Nil(t, err, "Update need to be valid")
}
func TestDeleteRemoveInvalid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	id := "0123456789"
	e := errors.New("failed")
	repo.On("Remove", id).Return(e)
	err := dataAgent.Delete(id)
	repo.AssertCalled(t, "Remove", id)
	assert.Equal(t, errors.Cause(err), e)
}
func TestDeleteValid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	id := "0123456789"
	repo.On("Remove", id).Return(nil)
	err := dataAgent.Delete(id)
	repo.AssertCalled(t, "Remove", id)
	assert.Nil(t, err, "Delete need to be valid")
}
func TestReadAllInvalid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	e := errors.New("failed")
	repo.On("ReadAll").Return([]string{}, e)
	_, err := dataAgent.ReadAll()
	repo.AssertCalled(t, "ReadAll")
	assert.Equal(t, errors.Cause(err), e)
}
func TestReadAllValid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	repo.On("ReadAll").Return([]string{"0123456789",}, nil)
	str, err := dataAgent.ReadAll()
	repo.AssertCalled(t, "ReadAll")
	assert.Equal(t, []string{"0123456789",}, str)
	assert.Nil(t, err, "ReadAll need to be valid")
}
func TestSendReadInvalid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	id := "0123456789"
	e := errors.New("failed")
	repo.On("FindById", id).Return(domain.Data{}, e)
	err := dataAgent.Send(id)
	repo.AssertCalled(t, "FindById", id)
	assert.Equal(t, errors.Cause(err), e)
}
func TestSendInvalid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	id := "0123456789"
	e := errors.New("failed")
	repo.On("FindById", id).Return(domain.Data{}, nil)
	client.On("Send", id, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(e)
	err := dataAgent.Send(id)
	repo.AssertCalled(t, "FindById", id)
	client.AssertCalled(t, "Send", id, mock.AnythingOfType("string"), mock.AnythingOfType("string"))
	assert.Equal(t, errors.Cause(err), e)
}
func TestSendValid(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	dataAgent, _ := NewDataAgent(
		repo,
		client)
	assert.NotNil(t, dataAgent)
	id := "0123456789"
	repo.On("FindById", id).Return(domain.Data{}, nil)
	client.On("Send", 
				id, 
				mock.AnythingOfType("string"), 
				mock.AnythingOfType("string")).Return(nil)
	err := dataAgent.Send(id)
	repo.AssertCalled(t, "FindById", id)
	client.AssertCalled(t, 
						"Send", 
						id, 
						mock.AnythingOfType("string"), 
						mock.AnythingOfType("string"))
	assert.Nil(t, err)
}

func TestNewDataAgentClientNil(t *testing.T) {
	repo := new(DataRepositoryMock)
	d, err := NewDataAgent(
		repo,
		nil)
	assert.Nil(t, d)
	assert.NotNil(t, err)
}
func TestNewDataAgentRepoNil(t *testing.T) {
	client := new(HttpClientMock)
	d, err := NewDataAgent(
		nil,
		client)
	assert.Nil(t, d)
	assert.NotNil(t, err)
}

func TestNewDataAgent(t *testing.T) {
	repo := new(DataRepositoryMock)
	client := new(HttpClientMock)
	d, err := NewDataAgent(
		repo,
		client)
	assert.Equal(t, d, &dataAgent{repo, client})
	assert.Nil(t, err)
}
func TestAssignInt(t *testing.T) {
	test := 55;
	assignInt(0, &test)
	assert.Equal(t, test, 55)
	assignInt(77, &test)
	assert.Equal(t, test, 77)
}
func TestAssignString(t *testing.T) {
	test := "55";
	assignString("", &test)
	assert.Equal(t, "55", test)
	assignString("77", &test)
	assert.Equal(t, "77", test)
}