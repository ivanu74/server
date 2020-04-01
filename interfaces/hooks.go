package interfaces

import (
	"encoding/json"
	"time"

	"b.yadro.com/sys/ch-server/usecases"
	"github.com/pkg/errors"
)

const maxLength = 64

type clientHooksI interface {
	GetChanTerm() chan string
}

// dataAgent - interface of dataAgent from usecases
type dataAgent interface {
	Create(data usecases.Data) error
	Read(id string) (usecases.Data, error)
	Update(id string, data usecases.Data) error
	Delete(id string) error
	IsUnique(data usecases.Data) error
	ReadAll() ([]string, error)
	Send(id string) error
}

// hooksHandler - implements interface hooksHandler from TusdHandler(infrastructure)
type hooksHandler struct {
	dataAgent   dataAgent
	clientHooks clientHooksI
	stdout      logger
}

// NewHooksHandler - create new hooksHandler instance
func NewHooksHandler(dataAgent dataAgent, clientHooks clientHooksI, stdlog logger) (*hooksHandler, error) {
	if dataAgent == nil || clientHooks == nil || stdlog == nil {
		return nil, errors.New("[hooks] [new] bad argument")
	}
	return &hooksHandler{dataAgent, clientHooks, stdlog}, nil
}

func (hook *hooksHandler) Validate(id string, data string) error {
	if !json.Valid([]byte(data)) {
		return errors.New("[hooks] [validate] invalid json in metadata")
	}
	meta := usecases.Data{}
	if err := json.Unmarshal([]byte(data), &meta); err != nil {
		return errors.Wrap(err, "[hooks] [validate]")
	}
	if err := preValidate(meta); err != nil {
		return errors.Wrap(err, "[hooks] [validate]")
	}
	if err := hook.dataAgent.IsUnique(meta); err != nil {
		return errors.Wrap(err, "[hooks] [validate]")
	}
	return nil
}

func (hook *hooksHandler) Create(id string, data string, name string) error {
	if !json.Valid([]byte(data)) {
		return errors.New("[hooks] [create] invalid json in metadata")
	}
	meta := usecases.Data{}
	if err := json.Unmarshal([]byte(data), &meta); err != nil {
		return errors.Wrap(err, "Hook Create")
	}
	meta.SessionID = id
	meta.FileName = name
	// Old: meta.StartTimestamp = time.Now().Format("2006/01/02 15:04:05.999")
	meta.StartTimestamp = time.Now().Format(time.UnixDate)
	if err := hook.dataAgent.Create(meta); err != nil {
		return errors.Wrap(err, "[hooks] [create]")
	}
	hook.stdout.Printf("[hooks] [create]: id = %s\n", id)
	return nil
}

func (hook *hooksHandler) Progress(id string) error {
	return nil
}

func (hook *hooksHandler) Terminate(id string) error {
	if err := hook.dataAgent.Delete(id); err != nil {
		return errors.Wrap(err, "[hooks] [terminate]")
	}
	hook.stdout.Printf("[hooks] [terminate]: id = %s\n", id)
	return nil
}

func (hook *hooksHandler) Complete(id string) error {
	meta := usecases.Data{}
	// Old: meta.FinishTimestamp = time.Now().Format("2006/01/02 15:04:05.999")
	meta.FinishTimestamp = time.Now().Format(time.UnixDate)
	if err := hook.dataAgent.Update(id, meta); err != nil {
		return errors.Wrap(err, "[hooks] [complete]")
	}
	hook.stdout.Printf("[hooks] [complete]: id = %s\n", id)
	meta, _ = hook.dataAgent.Read(id)
	hook.stdout.Printf("[hooks] [complete]: metadata = %v\n", meta)

	if err := hook.dataAgent.Send(id); err != nil {
		return errors.Wrap(err, "[hooks] [complete]")
	}
	return nil
}

func (hook *hooksHandler) GetChanTerm() chan string {
	return hook.clientHooks.GetChanTerm()
}
func preValidate(data usecases.Data) error {
	var err error
	if data.SerialNumber == "" {
		return errors.New("A SerialNumber may not be empty")
	}
	if _, err = time.Parse(time.UnixDate, data.LogCollectionTimestamp); err != nil {
		return errors.New("A LogCollectionTimestamp is invalid format")
	}
	if _, err = time.Parse(time.UnixDate, data.ClientStartTimestamp); err != nil {
		return errors.New("A ClientStartTimestamp is invalid format")
	}
	err = maxLen(err, maxLength, data.Service)
	err = maxLen(err, maxLength, data.SerialNumber)
	err = maxLen(err, maxLength, data.LogCollectionTimestamp)
	err = maxLen(err, maxLength, data.ClientStartTimestamp)
	err = maxLen(err, maxLength, data.SystemType)
	err = maxLen(err, maxLength, data.LogLevel)
	err = maxLen(err, maxLength, data.Originator)
	err = maxLen(err, maxLength, data.SessionID)
	err = maxLen(err, maxLength, data.Checksum)
	err = maxLen(err, maxLength, data.Hostname)
	err = maxLen(err, maxLength, data.NotificationManager)
	err = maxLen(err, maxLength, data.Cancel)
	err = maxLen(err, maxLength, data.StartTimestamp)
	err = maxLen(err, maxLength, data.FinishTimestamp)

	return err
}

func maxLen(err error, max int, str string) error {
	if err != nil {
		return err
	}
	if len(str) > maxLength {
		return errors.Errorf("Max length of string in metadata is %d", maxLength)
	}
	return nil
}
