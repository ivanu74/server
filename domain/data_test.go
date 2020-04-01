package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var data = Data{}

func fillData() {
	data.SessionID = "dc98fc675cfe7abfad2e0d06b56ace30"
	data.SerialNumber = "0123456789"
	data.FileName = "log.tar"
}
func TestValidData(t *testing.T) {
	fillData()
	assert.Nil(t, data.Validate())
}
func TestValidateSessionID(t *testing.T) {
	fillData()
	data.SessionID = ""
	assert.Error(t, data.Validate())
}
func TestValidateSerialNumber(t *testing.T) {
	fillData()
	data.SerialNumber = ""
	assert.Error(t, data.Validate())
}

func TestValidateFileName(t *testing.T) {
	fillData()
	data.FileName = ""
	assert.Error(t, data.Validate())
}
