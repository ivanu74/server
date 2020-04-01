package domain

import "github.com/pkg/errors"

// DataRepository - interface for save Data, implemented in interfaces/repositories
type DataRepository interface {
	Store(data Data) error
	FindById(id string) (Data, error)
	Remove(id string) error
	ReadAll() ([]string, error)
}

// Data - basic data to identify uploaded log file
type Data struct {
	ID                     int
	Service                string
	SerialNumber           string
	LogCollectionTimestamp string
	ClientStartTimestamp   string
	SystemType             string
	LogLevel               string
	Originator             string
	SessionID              string
	Checksum               string
	Hostname               string
	NotificationManager    string
	Cancel                 string
	StartTimestamp         string
	FinishTimestamp        string
	FileName               string
}

// Validate - test Data on correctness
func (data Data) Validate() error {
	if data.SessionID == "" {
		return errors.New("[data] [validate] a SessionID may not be empty")
	}
	if data.SerialNumber == "" {
		return errors.New("[data] [validate] a SerialNumber may not be empty")
	}
	if data.FileName == "" {
		return errors.New("[data] [validate] a FileName may not be empty")
	}
	return nil
}
