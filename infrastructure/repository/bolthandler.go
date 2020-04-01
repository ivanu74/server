package repository

import (
	"bytes"
	"os"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type logger interface {
	Fatalf(string, ...interface{})
}

// BoltHandler - implements interface DBHandler from repositories(interfaces)
type BoltHandler struct {
	dbname string
	errlog logger
}

func (handler *BoltHandler) Create(bucket []byte, key []byte, value []byte) error {
	conn := handler.open()
	defer handler.close(conn)
	err := conn.Update(
		func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				return errors.Wrap(err, "Create bolthandler")
			}
			if err := b.Put(key, value); err != nil {
				return errors.Wrap(err, "Create bolthandler")
			}
			return nil
		})
	return err
}

func (handler *BoltHandler) Get(bucket []byte, key []byte) ([]byte, error) {
	conn := handler.open()
	defer handler.close(conn)
	var value bytes.Buffer
	err := conn.View(
		func(tx *bolt.Tx) error {
			b := tx.Bucket(bucket)
			if b == nil {
				return errors.New("Name of bucket is wrong")
			}
			v := b.Get(key)
			if v == nil {
				return errors.New("Name of key is wrong")
			}
			value.Write(v)
			return nil
		})
	if err != nil {
		return nil, errors.Wrap(err, "Get bolthandler")
	}
	return value.Bytes(), nil
}

func (handler *BoltHandler) Delete(bucket []byte, key []byte) error {
	conn := handler.open()
	defer handler.close(conn)
	err := conn.Update(
		func(tx *bolt.Tx) error {
			b := tx.Bucket(bucket)
			if b == nil {
				return errors.New("Name of bucket is wrong")
			}
			if err := b.Delete(key); err != nil {
				return errors.Wrap(err, "Delete bolthandler")
			}
			return nil
		})
	return err
}

func (handler *BoltHandler) Keys(bucket []byte) ([][]byte, error) {
	conn := handler.open()
	defer handler.close(conn)
	values := make([][]byte, 0, 100)
	err := conn.View(
		func(tx *bolt.Tx) error {
			b := tx.Bucket(bucket)
			if b == nil {
				return errors.New("Name of bucket is wrong")
			}
			b.ForEach(func(k, _ []byte) error {
				key := make([]byte, len(k))
				copy(key, k)
				values = append(values, key)
				return nil
			})
			return nil
		})
	if err != nil {
		return nil, errors.Wrap(err, "Keys bolthandler")
	}
	return values, nil
}

func NewBoltHandler(dbfilename string, errlog logger) (*BoltHandler, error) {
	if dbfilename == "" || errlog == nil {
		return nil, errors.New("[bolthandler] [new handler] bad argument")
	}
	return &BoltHandler{dbfilename, errlog}, nil
}

func (handler *BoltHandler) open() *bolt.DB {
	conn, err := bolt.Open(handler.dbname, os.FileMode(0664), nil)
	if err != nil {
		handler.errlog.Fatalf("from open db: %v\n", err)
	}
	return conn
}

func (handler *BoltHandler) close(conn *bolt.DB) {
	conn.Close()
}
