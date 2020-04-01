package interfaces

import "github.com/pkg/errors"

// clientHandler - interface of SYRHandler from infrastructure/clienthandler
type clientHandler interface {
	Send(id string, name string, system string) error
}

// HTTPClient - implement interface httpClient from usecases
type HTTPClient struct {
	clientHandler clientHandler
	stdout        logger
}

// NewHTTPClient - create new instance of HTTPClient for NewDataAgent
func NewHTTPClient(client clientHandler, stdlog logger) (*HTTPClient, error) {
	if client == nil || stdlog == nil {
		return nil, errors.New("[client] [new] bad argument")
	}
	return &HTTPClient{client, stdlog}, nil
}

// Send - implement func to send file to extern server
func (client *HTTPClient) Send(id string, name string, system string) error {
	client.stdout.Printf("[client] [send]: id = %s; filename = %s\n", id, name)
	if err := client.clientHandler.Send(id, name, system); err != nil {
		return errors.Wrap(err, "[client] [send]")
	}
	return nil
}
