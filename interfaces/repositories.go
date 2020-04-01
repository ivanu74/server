package interfaces

import (
	"encoding/binary"
	"encoding/json"

	"b.yadro.com/sys/ch-server/domain"
	"github.com/pkg/errors"
)

const lengthOfKey = 32

// DbHandler - implemented in struct BoltHandler from infrastructure/repository
// to manage database
type DbHandler interface {
	Create(bucket []byte, key []byte, value []byte) error
	Get(bucket []byte, key []byte) ([]byte, error)
	Delete(bucket []byte, key []byte) error
	Keys(bucket []byte) ([][]byte, error)
}

// InvokeHandler - implemented in infrastructure/tusdinvoke to delete files from tusd
type InvokeHandler interface {
	Remove(id string) error
}

// DbDataRepo - Implement interface DataRepository from domain data
type DbDataRepo struct {
	dbHandler     DbHandler
	invokeHandler InvokeHandler
	bucket        string
}

// NewDbDataRepo - create instance of DbDataRepo(interface DataRepository)
func NewDbDataRepo(dbHandler DbHandler, invokeHandler InvokeHandler, bucket string) (*DbDataRepo, error) {
	if dbHandler == nil ||
		invokeHandler == nil ||
		bucket == "" {
		return nil, errors.New("[repositories] [new] bad argument")
	}
	return &DbDataRepo{dbHandler, invokeHandler, bucket}, nil
}

// Store - invoke db methods to store data in database
func (repo *DbDataRepo) Store(data domain.Data) error {
	b, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "[repositories] [store]")
	}
	key := data.SessionID
	if key == "" {
		return errors.New("[repositories] [store] bad data")
	}
	if err := repo.dbHandler.Create([]byte(repo.bucket), []byte(key), b); err != nil {
		return errors.Wrap(err, "[repositories] [store]")
	}
	key = data.SerialNumber + data.LogCollectionTimestamp
	if key == "" {
		return errors.New("[repositories] [store] bad data")
	}
	if err := repo.dbHandler.Create([]byte(repo.bucket), []byte(key), b); err != nil {
		return errors.Wrap(err, "[repositories] [store]")
	}
	return nil
}

// FindById - invoke db methods to find data in database and return data
func (repo *DbDataRepo) FindById(id string) (domain.Data, error) {
	var data domain.Data
	js := make([]byte, 0, 512)
	js, err := repo.dbHandler.Get([]byte(repo.bucket), []byte(id))
	if err != nil {
		return data, errors.Wrap(err, "[repositories] [findById]")
	}
	err = json.Unmarshal(js, &data)
	if err != nil {
		return data, errors.Wrap(err, "[repositories] [findById]")
	}
	return data, nil
}

// Remove - invoke db methods to delete data from database
func (repo *DbDataRepo) Remove(id string) error {
	data, err := repo.FindById(id)
	if err != nil {
		return errors.Wrap(err, "[repositories] [remove]")
	}
	keys := []string{
		data.SessionID,
		data.SerialNumber + data.LogCollectionTimestamp,
	}
	for _, key := range keys {
		if err := repo.dbHandler.Delete([]byte(repo.bucket), []byte(key)); err != nil {
			return errors.Wrap(err, "[repositories] [remove]")
		}
	}
	if err := repo.invokeHandler.Remove(data.SessionID); err != nil {
		return errors.Wrap(err, "[repositories] [remove]")
	}
	return nil
}

// ReadAll - invoke db methods to read all id from database and return string of all id
func (repo *DbDataRepo) ReadAll() ([]string, error) {

	keys, err := repo.dbHandler.Keys([]byte(repo.bucket))
	if err != nil {
		return nil, errors.Wrap(err, "[repositories] [readAll]")
	}
	strings := make([]string, 0, len(keys))
	for _, key := range keys {
		if len(key) == lengthOfKey {
			strings = append(strings, string(key))
		}
	}
	return strings, nil
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
