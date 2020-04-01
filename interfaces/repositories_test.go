package interfaces

import (
	"testing"
	"encoding/json"

	"b.yadro.com/sys/ch-server/domain"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewDbDataRepoDbHandlerNil(t *testing.T) {
	db, err := NewDbDataRepo(
		nil,
		new(InvokeHandlerMock),
		"root")
	assert.Nil(t, db)
	assert.NotNil(t, err)
}

func TestNewDbDataRepoInvokeHandlerNil(t *testing.T) {
	db, err := NewDbDataRepo(
		new(DbHandlerMock),
		nil,
		"root")
	assert.Nil(t, db)
	assert.NotNil(t, err)
}
func TestNewDbDataRepoBucketNil(t *testing.T) {
	db, err := NewDbDataRepo(
		new(DbHandlerMock),
		new(InvokeHandlerMock),
		"")
	assert.Nil(t, db)
	assert.NotNil(t, err)
}

func TestNewDbDataRepoValid(t *testing.T) {
	db, err := NewDbDataRepo(
		new(DbHandlerMock),
		new(InvokeHandlerMock),
		"root")
	assert.NotNil(t, db)
	assert.Nil(t, err)
}

func TestStoreValid(t *testing.T) {
	db := new(DbHandlerMock)
	invoke := new(InvokeHandlerMock)
	repo, _ := NewDbDataRepo(
		db,
		invoke,
		"root")

	meta := domain.Data{}
	meta.SessionID = "0123456789"
	meta.SerialNumber = "0123456789"
	meta.LogCollectionTimestamp = "Thu Aug 17 14:00:06 MSK 2019"
	db.On("Create",
		[]byte(repo.bucket),
		[]byte(meta.SessionID),
		mock.Anything).Return(nil)
	db.On("Create",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp),
		mock.Anything).Return(nil)
	err := repo.Store(meta)
	assert.Nil(t, err)
	db.AssertCalled(t, "Create",
		[]byte(repo.bucket),
		[]byte(meta.SessionID),
		mock.Anything)
	db.AssertCalled(t, "Create",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp),
		mock.Anything)
}

func TestStoreDataInvalid(t *testing.T) {
	db := new(DbHandlerMock)
	invoke := new(InvokeHandlerMock)
	repo, _ := NewDbDataRepo(db, invoke, "root")
	meta := domain.Data{}
	db.On("Create",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp),
		mock.Anything).Return(nil)
	// test empty SessionID
	err := repo.Store(meta)
	assert.NotNil(t, err)
	db.AssertNotCalled(t, "Create",
		[]byte(repo.bucket),
		[]byte(meta.SessionID),
		mock.Anything)
	// test empty meta.SerialNumber+meta.LogCollectionTimestamp
	meta.SessionID = "0123456789"
	db.On("Create",
		[]byte(repo.bucket),
		[]byte(meta.SessionID),
		mock.Anything).Return(nil)
	err = repo.Store(meta)
	assert.NotNil(t, err)
	db.AssertNotCalled(t, "Create",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp),
		mock.Anything)
}

func TestStoreFirstCreateInvalid(t *testing.T) {
	db := new(DbHandlerMock)
	invoke := new(InvokeHandlerMock)
	repo, _ := NewDbDataRepo(
		db,
		invoke,
		"root")
	meta := domain.Data{}
	meta.SessionID = "0123456789"
	meta.SerialNumber = "0123456789"
	meta.LogCollectionTimestamp = "Thu Aug 17 14:00:06 MSK 2019"
	e := errors.New("fail")
	db.On("Create",
		[]byte(repo.bucket),
		[]byte(meta.SessionID),
		mock.Anything).Return(e)
	db.On("Create",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp),
		mock.Anything).Return(nil)
	err := repo.Store(meta)
	assert.Equal(t, errors.Cause(err), e)
	db.AssertCalled(t, "Create",
		[]byte(repo.bucket),
		[]byte(meta.SessionID),
		mock.Anything)
	db.AssertNotCalled(t, "Create",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp),
		mock.Anything)
}

func TestStoreSecondCreateInvalid(t *testing.T) {
	db := new(DbHandlerMock)
	invoke := new(InvokeHandlerMock)
	repo, _ := NewDbDataRepo(
		db,
		invoke,
		"root")
	meta := domain.Data{}
	meta.SessionID = "0123456789"
	meta.SerialNumber = "0123456789"
	meta.LogCollectionTimestamp = "Thu Aug 17 14:00:06 MSK 2019"
	e := errors.New("fail")
	db.On("Create",
		[]byte(repo.bucket),
		[]byte(meta.SessionID),
		mock.Anything).Return(nil)
	db.On("Create",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp),
		mock.Anything).Return(e)
	err := repo.Store(meta)
	assert.Equal(t, errors.Cause(err), e)
	db.AssertCalled(t, "Create",
		[]byte(repo.bucket),
		[]byte(meta.SessionID),
		mock.Anything)
	db.AssertCalled(t, "Create",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp),
		mock.Anything)
}

func TestFindByIdValid(t *testing.T) {
	db := new(DbHandlerMock)
	invoke := new(InvokeHandlerMock)
	repo, _ := NewDbDataRepo(db, invoke, "root")
	meta := domain.Data{}
	meta.SessionID = "0123456789"
	meta.SerialNumber = "0123456789"
	meta.LogCollectionTimestamp = "Thu Aug 17 14:00:06 MSK 2019"
	data, _ := json.Marshal(meta)
	db.On("Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID)).Return([]byte(data), nil)
	respmeta, err := repo.FindById(meta.SessionID)
	assert.Nil(t, err)
	assert.Equal(t, respmeta, meta)
	db.AssertCalled(t, "Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID))
}
func TestFindByIdGetInvalid(t *testing.T) {
	db := new(DbHandlerMock)
	invoke := new(InvokeHandlerMock)
	repo, _ := NewDbDataRepo(db, invoke, "root")
	meta := domain.Data{}
	meta.SessionID = "0123456789"
	meta.SerialNumber = "0123456789"
	meta.LogCollectionTimestamp = "Thu Aug 17 14:00:06 MSK 2019"
	data, _ := json.Marshal(meta)
	e := errors.New("fail")
	db.On("Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID)).Return([]byte(data), e)
	respmeta, err := repo.FindById(meta.SessionID)
	assert.Equal(t, errors.Cause(err), e)
	assert.Equal(t, respmeta, domain.Data{})
	db.AssertCalled(t, "Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID))
}

func TestRemoveValid(t *testing.T) {
	db := new(DbHandlerMock)
	invoke := new(InvokeHandlerMock)
	repo, _ := NewDbDataRepo(db, invoke, "root")
	meta := domain.Data{}
	meta.SessionID = "0123456789"
	meta.SerialNumber = "0123456789"
	meta.LogCollectionTimestamp = "Thu Aug 17 14:00:06 MSK 2019"
	data, _ := json.Marshal(meta)
	db.On("Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID)).Return([]byte(data), nil)
	db.On("Delete",
		[]byte(repo.bucket),
		[]byte(meta.SessionID)).Return(nil)
	db.On("Delete",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp)).Return(nil)
	invoke.On("Remove", meta.SessionID).Return(nil)
	err := repo.Remove(meta.SessionID)
	assert.Nil(t, err)
	db.AssertCalled(t, "Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID))
	db.AssertCalled(t, "Delete",
		[]byte(repo.bucket),
		[]byte(meta.SessionID))
	db.AssertCalled(t, "Delete",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp))
	invoke.AssertCalled(t, "Remove", meta.SessionID)
}
func TestRemoveFindByIdInvalid(t *testing.T) {
	db := new(DbHandlerMock)
	invoke := new(InvokeHandlerMock)
	repo, _ := NewDbDataRepo(db, invoke, "root")
	meta := domain.Data{}
	meta.SessionID = "0123456789"
	meta.SerialNumber = "0123456789"
	meta.LogCollectionTimestamp = "Thu Aug 17 14:00:06 MSK 2019"
	data, _ := json.Marshal(meta)
	e := errors.New("fail")
	db.On("Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID)).Return([]byte(data), e)
	db.On("Delete",
		[]byte(repo.bucket),
		[]byte(meta.SessionID)).Return(nil)
	db.On("Delete",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp)).Return(nil)
	invoke.On("Remove", meta.SessionID).Return(nil)
	err := repo.Remove(meta.SessionID)
	assert.Equal(t, errors.Cause(err), e)
	db.AssertCalled(t, "Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID))
	db.AssertNotCalled(t, "Delete",
		[]byte(repo.bucket),
		[]byte(meta.SessionID))
	db.AssertNotCalled(t, "Delete",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp))
	invoke.AssertNotCalled(t, "Remove", meta.SessionID)
}
func TestRemoveFirstDeleteInvalid(t *testing.T) {
	db := new(DbHandlerMock)
	invoke := new(InvokeHandlerMock)
	repo, _ := NewDbDataRepo(db, invoke, "root")
	meta := domain.Data{}
	meta.SessionID = "0123456789"
	meta.SerialNumber = "0123456789"
	meta.LogCollectionTimestamp = "Thu Aug 17 14:00:06 MSK 2019"
	data, _ := json.Marshal(meta)
	e := errors.New("fail")
	db.On("Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID)).Return([]byte(data), nil)
	db.On("Delete",
		[]byte(repo.bucket),
		[]byte(meta.SessionID)).Return(e)
	db.On("Delete",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp)).Return(nil)
	invoke.On("Remove", meta.SessionID).Return(nil)
	err := repo.Remove(meta.SessionID)
	assert.Equal(t, errors.Cause(err), e)
	db.AssertCalled(t, "Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID))
	db.AssertCalled(t, "Delete",
		[]byte(repo.bucket),
		[]byte(meta.SessionID))
	db.AssertNotCalled(t, "Delete",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp))
	invoke.AssertNotCalled(t, "Remove", meta.SessionID)
}

func TestRemoveSecondDeleteInvalid(t *testing.T) {
	db := new(DbHandlerMock)
	invoke := new(InvokeHandlerMock)
	repo, _ := NewDbDataRepo(db, invoke, "root")
	meta := domain.Data{}
	meta.SessionID = "0123456789"
	meta.SerialNumber = "0123456789"
	meta.LogCollectionTimestamp = "Thu Aug 17 14:00:06 MSK 2019"
	data, _ := json.Marshal(meta)
	e := errors.New("fail")
	db.On("Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID)).Return([]byte(data), nil)
	db.On("Delete",
		[]byte(repo.bucket),
		[]byte(meta.SessionID)).Return(nil)
	db.On("Delete",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp)).Return(e)
	invoke.On("Remove", meta.SessionID).Return(nil)
	err := repo.Remove(meta.SessionID)
	assert.Equal(t, errors.Cause(err), e)
	db.AssertCalled(t, "Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID))
	db.AssertCalled(t, "Delete",
		[]byte(repo.bucket),
		[]byte(meta.SessionID))
	db.AssertCalled(t, "Delete",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp))
	invoke.AssertNotCalled(t, "Remove", meta.SessionID)
}

func TestRemoveRemoveInvalid(t *testing.T) {
	db := new(DbHandlerMock)
	invoke := new(InvokeHandlerMock)
	repo, _ := NewDbDataRepo(db, invoke, "root")
	meta := domain.Data{}
	meta.SessionID = "0123456789"
	meta.SerialNumber = "0123456789"
	meta.LogCollectionTimestamp = "Thu Aug 17 14:00:06 MSK 2019"
	data, _ := json.Marshal(meta)
	e := errors.New("fail")
	db.On("Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID)).Return([]byte(data), nil)
	db.On("Delete",
		[]byte(repo.bucket),
		[]byte(meta.SessionID)).Return(nil)
	db.On("Delete",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp)).Return(nil)
	invoke.On("Remove", meta.SessionID).Return(e)
	err := repo.Remove(meta.SessionID)
	assert.Equal(t, errors.Cause(err), e)
	db.AssertCalled(t, "Get",
		[]byte(repo.bucket),
		[]byte(meta.SessionID))
	db.AssertCalled(t, "Delete",
		[]byte(repo.bucket),
		[]byte(meta.SessionID))
	db.AssertCalled(t, "Delete",
		[]byte(repo.bucket),
		[]byte(meta.SerialNumber+meta.LogCollectionTimestamp))
	invoke.AssertCalled(t, "Remove", meta.SessionID)
}
func TestReadAllValid(t *testing.T) {
	db := new(DbHandlerMock)
	invoke := new(InvokeHandlerMock)
	repo, _ := NewDbDataRepo(db, invoke, "root")
	id := "01234567890123456789012345678901"
	data := [][]byte{[]byte(id), []byte(id)}
	strings := []string{id, id}
	db.On("Keys", []byte(repo.bucket)).Return(data, nil)
	s, err := repo.ReadAll()
	assert.Nil(t, err)
	assert.Equal(t, strings, s)
	db.AssertCalled(t, "Keys",
		[]byte(repo.bucket))
}
func TestReadAllInvalid(t *testing.T) {
	db := new(DbHandlerMock)
	invoke := new(InvokeHandlerMock)
	repo, _ := NewDbDataRepo(db, invoke, "root")
	id := "01234567890123456789012345678901"
	data := [][]byte{[]byte(id), []byte(id)}
	e := errors.New("fail")
	db.On("Keys", []byte(repo.bucket)).Return(data, e)
	_, err := repo.ReadAll()
	assert.Equal(t, errors.Cause(err), e)
	db.AssertCalled(t, "Keys",
		[]byte(repo.bucket))
}