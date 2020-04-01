package infrastructure

import (
	"context"

	"github.com/pkg/errors"
	tusd "github.com/tus/tusd/pkg/handler"
)

var ctx = context.Background() // TODO: pass through ctx from main().

// TusdInvoke - wrapper for tusd.StoreComposer
// implement interface InvokeHandler from interfaces/repositories
type TusdInvoke struct {
	composer *tusd.StoreComposer
}

// Remove - implement delete files from tusd
func (invoke *TusdInvoke) Remove(id string) error {
	// Abort the request handling if the required interface is not implemented
	if !invoke.composer.UsesTerminater {
		return errors.New("[tusdinvoke] [remove] not support terminate")
	}
	upload, err := invoke.composer.Core.GetUpload(ctx, id)
	if err != nil {
		return nil
	}
	var lock tusd.Lock
	if invoke.composer.UsesLocker {
		locker := invoke.composer.Locker
		if lock, err = locker.NewLock(id); err != nil {
			return errors.Wrap(err, "[tusdinvoke] [remove]")
		}
		if err = lock.Lock(); err != nil {
			return errors.Wrap(err, "[tusdinvoke] [remove]")
		}
		defer lock.Unlock()
	}

	if err := invoke.composer.Terminater.AsTerminatableUpload(upload).Terminate(ctx); err != nil {
		return errors.Wrap(err, "[tusdinvoke] [remove]")
	}
	return nil
}

// NewTusdInvoke - create instance of InvokeHandler interface from interfaces/repositories
func NewTusdInvoke(composer *tusd.StoreComposer) (*TusdInvoke, error) {
	if composer == nil {
		return nil, errors.New("[tusdinvoke] [new] bad argument")
	}
	return &TusdInvoke{composer}, nil
}
