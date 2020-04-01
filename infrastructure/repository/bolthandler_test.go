package repository

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoltHandler(t *testing.T) {
	dbfile := "test.db"
	logerr := log.New(os.Stdout, "[test] ", log.LstdFlags)

	t.Run("valid New", func(t *testing.T) {
		d, err := NewBoltHandler(
			dbfile,
			logerr)
		assert.Equal(t, d, &BoltHandler{dbfile, logerr})
		assert.Nil(t, err)
	})

	t.Run("invalid DbFile", func(t *testing.T) {
		d, err := NewBoltHandler(
			"",
			logerr)
		assert.Nil(t, d)
		assert.NotNil(t, err)
	})

	t.Run("invalid Errlog", func(t *testing.T) {
		d, err := NewBoltHandler(
			dbfile,
			nil)
		assert.Nil(t, d)
		assert.NotNil(t, err)
	})
}

func TestBoltHandler(t *testing.T) {
	dbfile := "/tmp/test.db"
	logerr := log.New(os.Stdout, "[test] ", log.LstdFlags)
	d, err := NewBoltHandler(
		dbfile,
		logerr)
	assert.Nil(t, err)
	defer func() {
		if _, err := os.Stat(dbfile); err == nil {
			os.Remove(dbfile)
		}
	}()

	t.Run("valid create get keys delete", func(t *testing.T) {
		err := d.Create([]byte("test"), []byte("testkey"), []byte("testvalue"))
		assert.Nil(t, err)
		value, err := d.Get([]byte("test"), []byte("testkey"))
		assert.Nil(t, err)
		assert.Equal(t, []byte("testvalue"), value)
		keys, err := d.Keys([]byte("test"))
		assert.Nil(t, err)
		assert.Equal(t, [][]byte{[]byte("testkey")}, keys)
		err = d.Delete([]byte("test"), []byte("testkey"))
		assert.Nil(t, err)
		_, err = d.Get([]byte("test"), []byte("testkey"))
		assert.NotNil(t, err)
	})
}
