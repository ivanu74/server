package usecases

import (
	"b.yadro.com/sys/ch-server/domain"
	"github.com/pkg/errors"
)

// Data - struct for use in usecases, protect Data from domain package
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

type httpClient interface {
	Send(id string, name string, system string) error
	// TODO: - define Download or Send func
}

// dataAgent - Implement dataAgent interface from interfaces hooks.
type dataAgent struct {
	DataRepository domain.DataRepository
	DataClient     httpClient
}

func (agent *dataAgent) IsUnique(data Data) error {
	unique := data.SerialNumber + data.LogCollectionTimestamp
	if _, err := agent.DataRepository.FindById(unique); err == nil {
		return errors.New("MetaData is not unique") // TODO: convert error to [usedata] [unique] for all file
	}
	return nil
}
func (agent *dataAgent) Create(data Data) error {
	d := domain.Data{
		ID:                     data.ID,
		Service:                data.Service,
		SerialNumber:           data.SerialNumber,
		LogCollectionTimestamp: data.LogCollectionTimestamp,
		ClientStartTimestamp:   data.ClientStartTimestamp,
		SystemType:             data.SystemType,
		LogLevel:               data.LogLevel,
		Originator:             data.Originator,
		SessionID:              data.SessionID,
		Checksum:               data.Checksum,
		Hostname:               data.Hostname,
		NotificationManager:    data.NotificationManager,
		Cancel:                 data.Cancel,
		StartTimestamp:         data.StartTimestamp,
		FinishTimestamp:        data.FinishTimestamp,
		FileName:               data.FileName,
	}
	if err := d.Validate(); err != nil {
		return errors.Wrap(err, "Create data")
	}
	if err := agent.DataRepository.Store(d); err != nil {
		return errors.Wrap(err, "Create data")
	}
	return nil
}

func (agent *dataAgent) Read(id string) (Data, error) {
	d, err := agent.DataRepository.FindById(id)
	if err != nil {
		return Data{}, errors.Wrap(err, "Read data")
	}
	data := Data{
		ID:                     d.ID,
		Service:                d.Service,
		SerialNumber:           d.SerialNumber,
		LogCollectionTimestamp: d.LogCollectionTimestamp,
		ClientStartTimestamp:   d.ClientStartTimestamp,
		SystemType:             d.SystemType,
		LogLevel:               d.LogLevel,
		Originator:             d.Originator,
		SessionID:              d.SessionID,
		Checksum:               d.Checksum,
		Hostname:               d.Hostname,
		NotificationManager:    d.NotificationManager,
		Cancel:                 d.Cancel,
		StartTimestamp:         d.StartTimestamp,
		FinishTimestamp:        d.FinishTimestamp,
		FileName:               d.FileName,
	}
	return data, nil
}

func (agent *dataAgent) Update(id string, data Data) error {
	d, err := agent.DataRepository.FindById(id)
	if err != nil {
		return errors.Wrap(err, "Update data")
	}
	assignInt(data.ID, &d.ID)
	assignString(data.Service, &d.Service)
	assignString(data.SerialNumber, &d.SerialNumber)
	assignString(data.LogCollectionTimestamp, &d.LogCollectionTimestamp)
	assignString(data.ClientStartTimestamp, &d.ClientStartTimestamp)
	assignString(data.SystemType, &d.SystemType)
	assignString(data.LogLevel, &d.LogLevel)
	assignString(data.Originator, &d.Originator)
	assignString(data.SessionID, &d.SessionID)
	assignString(data.Checksum, &d.Checksum)
	assignString(data.Hostname, &d.Hostname)
	assignString(data.NotificationManager, &d.NotificationManager)
	assignString(data.Cancel, &d.Cancel)
	assignString(data.StartTimestamp, &d.StartTimestamp)
	assignString(data.FinishTimestamp, &d.FinishTimestamp)
	assignString(data.FileName, &d.FileName)

	if err := d.Validate(); err != nil {
		return errors.Wrap(err, "Update data")
	}
	if err := agent.DataRepository.Store(d); err != nil {
		return errors.Wrap(err, "Update data")
	}
	return nil
}

func (agent *dataAgent) Delete(id string) error {
	err := agent.DataRepository.Remove(id)
	return errors.Wrap(err, "Delete data")
}

func (agent *dataAgent) ReadAll() ([]string, error) {
	strings, err := agent.DataRepository.ReadAll()
	return strings, errors.Wrap(err, "Delete data")
}

func (agent *dataAgent) Send(id string) error {
	meta, err := agent.Read(id)
	if err == nil {
		err = agent.DataClient.Send(id, meta.FileName, meta.SerialNumber)
	}
	return errors.Wrap(err, "[usedata] [send]")
}

// NewDataAgent - create dataAgent for invoke from hooksHandler interfaces
func NewDataAgent(repo domain.DataRepository, client httpClient) (*dataAgent, error) {
	if repo == nil || client == nil {
		return nil, errors.New("[usedata] [new] bad argument")
	}
	return &dataAgent{repo, client}, nil
}

func assignInt(src int, dst *int) {
	if src != 0 {
		*dst = src
	}
}
func assignString(src string, dst *string) {
	if src != "" {
		*dst = src
	}
}
