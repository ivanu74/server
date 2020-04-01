package infrastructure

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"b.yadro.com/sys/ch-server/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewTusd(t *testing.T) {
	composer := NewStoreComposer()
	logger := log.New(os.Stdout, "[test] ", log.LstdFlags)
	hooks := new(interfaces.HooksHandlerMock)
	t.Run("valid New", func(t *testing.T) {
		tusd, err := NewHooksTusdHandler(
			composer,
			hooks,
			logger,
		)
		assert.Nil(t, err)
		assert.NotNil(t, tusd)
	})
	t.Run("invalid argument composer", func(t *testing.T) {
		tusd, err := NewHooksTusdHandler(
			nil,
			hooks,
			logger,
		)
		assert.NotNil(t, err)
		assert.Nil(t, tusd)
	})
	t.Run("invalid argument hooks", func(t *testing.T) {
		tusd, err := NewHooksTusdHandler(
			composer,
			nil,
			logger,
		)
		assert.NotNil(t, err)
		assert.Nil(t, tusd)
	})
	t.Run("invalid argument logger", func(t *testing.T) {
		tusd, err := NewHooksTusdHandler(
			composer,
			hooks,
			nil,
		)
		assert.NotNil(t, err)
		assert.Nil(t, tusd)
	})
}

func TestRestTusd(t *testing.T) {
	composer := NewStoreComposer()
	filepath := "/tmp/test/"
	urlpath := "/test/"
	handler, err := TusdConfig(
		composer,
		filepath,
		urlpath)
	logger := log.New(os.Stdout, "[test] ", log.LstdFlags)
	hooks := new(interfaces.HooksHandlerMock)
	tusd, err := NewHooksTusdHandler(
		composer,
		hooks,
		logger,
	)
	assert.Nil(t, err)
	assert.NotNil(t, tusd)
	if err := os.MkdirAll(filepath, 0777); err != nil {
		assert.FailNow(t, "unable to make dir: %v", err)
	}
	defer os.RemoveAll(filepath)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	hooks.On("Validate", mock.Anything, mock.Anything).Return(nil)
	hooks.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	hooks.On("Complete", mock.Anything).Return(nil)
	hooks.On("Progress", mock.Anything).Return(nil)
	hooks.On("GetChanTerm").Return(make(chan string))
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Start hooks from tusd to usecases
		err := tusd.RunHooks(ctx, handler)
		assert.Nil(t, err, "unable to run hooks.")
	}()
	res := (&httpTest{
		Method: "POST",
		ReqHeader: map[string]string{
			"Tus-Resumable":   "1.0.0",
			"Upload-Offset":   "0",
			"Upload-Length":   "12",
			"Content-Type":    "application/offset+octet-stream",
			"Upload-Metadata": "data aGVsbG8=, filename d29ybGQ=",
		},
		ReqBody: strings.NewReader("hello world!"),
		Code:    http.StatusCreated,
		ResHeader: map[string]string{
			"Upload-Offset": "12",
		},
	}).Run(handler, t)
	time.Sleep(10 * time.Millisecond)
	cancel()
	wg.Wait()
	s := strings.TrimPrefix(res.Header().Get("Location"), "http://tus.io/test/")
	hooks.AssertCalled(t, "Complete", s)
	hooks.AssertCalled(t, "Progress", s)
	// Create(id, metadata.data, metadata.filename)
	hooks.AssertCalled(t, "Create", s, "hello", "world")
	// Validate(id, metadata.data)
	hooks.AssertCalled(t, "Validate", "", "hello")
	// Create a new handler to invoke tusd functions
	invoke, err := NewTusdInvoke(composer)
	assert.Nil(t, err)
	assert.NotNil(t, invoke)
	if err := invoke.Remove(s); err != nil {
		assert.FailNow(t, "unable to delete: %s, error: %v", s, err)
	}
}
