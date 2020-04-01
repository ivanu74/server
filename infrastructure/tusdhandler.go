package infrastructure

import (
	"context"

	"github.com/pkg/errors"
	tusd "github.com/tus/tusd/pkg/handler"
)

//type logger interface {
//	Printf(format string, v ...interface{})
//}
type hooksHandler interface {
	Validate(id string, data string) error
	Create(id string, data string, name string) error
	Progress(id string) error
	Terminate(id string) error
	Complete(id string) error
	GetChanTerm() chan string
}

type hookType string

const (
	hookPostFinish    hookType = "post-finish"
	hookPostTerminate hookType = "post-terminate"
	hookPostReceive   hookType = "post-receive"
	hookPostCreate    hookType = "post-create"
	hookPreCreate     hookType = "pre-create"
	hookTerminate     hookType = "pre-terminate"
)

type HooksTusdHandler struct {
	hooks  hooksHandler
	stderr logger
}

type hookDataStore struct {
	tusd.DataStore
	tusdhandler *HooksTusdHandler
}

func (store hookDataStore) NewUpload(ctx context.Context, info tusd.FileInfo) (upload tusd.Upload, err error) {
	if err := store.tusdhandler.invokeHook(hookPreCreate, info); err != nil {
		return nil, errors.Wrapf(err, "hook %s", hookPreCreate)
	}
	return store.DataStore.NewUpload(ctx, info)
}

func (handler *HooksTusdHandler) RunHooks(ctx context.Context, notify *tusd.Handler) error {
	for {
		select {
		case info := <-notify.CompleteUploads:
			if err := handler.invokeHook(hookPostFinish, info.Upload); err != nil {
				handler.stderr.Printf("notify %s: %s", hookPostFinish, err)
			}
		case info := <-notify.TerminatedUploads:
			if err := handler.invokeHook(hookPostTerminate, info.Upload); err != nil {
				handler.stderr.Printf("notify %s: %s", hookPostTerminate, err)
			}
		case info := <-notify.UploadProgress:
			if err := handler.invokeHook(hookPostReceive, info.Upload); err != nil {
				handler.stderr.Printf("notify %s: %s", hookPostReceive, err)
			}
		case info := <-notify.CreatedUploads:
			if err := handler.invokeHook(hookPostCreate, info.Upload); err != nil {
				handler.stderr.Printf("notify %s: %s", hookPostCreate, err)
			}
		case id := <-handler.hooks.GetChanTerm():
			fileinfo := tusd.FileInfo{ID: id}
			if err := handler.invokeHook(hookTerminate, fileinfo); err != nil {
				handler.stderr.Printf("notify %s: %s", hookTerminate, err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func getMetaData(info tusd.FileInfo, key string) string {
	return info.MetaData[key]
}

func NewHooksTusdHandler(
	composer *tusd.StoreComposer,
	handler hooksHandler,
	errlog logger) (*HooksTusdHandler, error) {
	if composer == nil || handler == nil || errlog == nil {
		return nil, errors.New("[tusdhandler] [new] bad argument")
	}
	tusdhandler := &HooksTusdHandler{handler, errlog}
	composer.UseCore(hookDataStore{
		composer.Core,
		tusdhandler,
	})
	return tusdhandler, nil
}

func (handler *HooksTusdHandler) invokeHook(typ hookType, info tusd.FileInfo) error {
	if handler.hooks == nil {
		return errors.New("[tusd] [hooks] Hooks handler is nil")
	}
	switch typ {
	case hookPreCreate:
		return handler.hooks.Validate(info.ID, getMetaData(info, "data"))
	case hookPostCreate:
		return handler.hooks.Create(info.ID, getMetaData(info, "data"), getMetaData(info, "filename"))
	case hookPostFinish:
		return handler.hooks.Complete(info.ID)
	case hookPostTerminate:
		return handler.hooks.Terminate(info.ID)
	case hookPostReceive:
		return handler.hooks.Progress(info.ID)
	case hookTerminate:
		return handler.hooks.Terminate(info.ID)
	}
	return nil
}
