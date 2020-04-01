package infrastructure

import (
	"github.com/pkg/errors"
	"github.com/tus/tusd/pkg/filelocker"
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
)

// NewStoreComposer - wrapper of customized StoreComposer
func NewStoreComposer() *tusd.StoreComposer {
	return tusd.NewStoreComposer()
}

// TusdConfig - instance of tusd configure
func TusdConfig(composer *tusd.StoreComposer, filepath string, urlpath string) (*tusd.Handler, error) {

	// Create a new FileStore instance.
	store := filestore.New(filepath)
	// Composer use the file store.
	store.UseIn(composer)
	locker := filelocker.New(filepath)
	locker.UseIn(composer)

	// Create a new HTTP handler for the tusd server by providing a configuration.
	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:                urlpath,
		StoreComposer:           composer,
		NotifyCompleteUploads:   true,
		NotifyTerminatedUploads: true,
		NotifyUploadProgress:    true,
		NotifyCreatedUploads:    true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "[tusd] Unable to create handler: ")
	}
	return handler, nil
}
