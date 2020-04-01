package main_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"b.yadro.com/sys/ch-server/infrastructure"
	"b.yadro.com/sys/ch-server/infrastructure/repository"
	"b.yadro.com/sys/ch-server/interfaces"
	"b.yadro.com/sys/ch-server/usecases"
	"github.com/stretchr/testify/assert"
)

const data = `
ewogICAgIlNlcnZpY2UiICAgICAgICAgICAgICAgOiAidHJ1ZSIsCiAgICAi
U2VyaWFsTnVtYmVyIiAgICAgICAgICA6ICIwMTIzNDU2Nzg5IiwKICAgICJM
b2dDb2xsZWN0aW9uVGltZXN0YW1wIjogIlRodSBEZWMgMjMgMTQ6MDA6MTIg
TVNLIDIwMTkiLAogICAgIkNsaWVudFN0YXJ0VGltZXN0YW1wIiAgOiAiVGh1
IE5vdiAwNiAxNDowMDowNiBNU0sgMjAxOSIsCiAgICAiU3lzdGVtVHlwZSIg
ICAgICAgICAgICA6ICJ0YXRsaW4iLAogICAgIkxvZ0xldmVsIiAgICAgICAg
ICAgICAgOiAiU3lzdGVtIiwKICAgICJPcmlnaW5hdG9yIiAgICAgICAgICAg
IDogInN5c3RlbSIsCiAgICAiU2Vzc2lvbklEIiAgICAgICAgICAgICA6ICIi
LAogICAgIkNoZWNrc3VtIiAgICAgICAgICAgICAgOiAiZGM5OGZjNjc1Y2Zl
N2FiZmFkMmUwZDA2YjU2YWNlMzAiLAogICAgIkhvc3RuYW1lIiAgICAgICAg
ICAgICAgOiAiZXhhbXBsZS5jb20iLAogICAgIk5vdGlmaWNhdGlvbk1hbmFn
ZXIiICAgOiIiLAogICAgIkNhbmNlbCIgICAgICAgICAgICAgICAgOiAiTm8i
LAogICAgIlN0YXJ0VGltZXN0YW1wIiAgICAgICAgOiAiIiwKICAgICJGaW5p
c2hUaW1lc3RhbXAiICAgICAgIDogIiIKICAgIH0=`

func TestMain(t *testing.T) {
	// New composer of tus
	composer := infrastructure.NewStoreComposer()
	assert.NotNil(t, composer)
	// New handler of tus
	filepath := "/tmp/"
	urlpath := "/test/"
	tusdHandler, err := infrastructure.TusdConfig(
		composer,
		filepath,
		urlpath)
	assert.Nil(t, err)
	assert.NotNil(t, tusdHandler)
	// New boltHandler for store metadata
	var testbuffer bytes.Buffer
	logger := log.New(&testbuffer, "[test] ", log.LstdFlags)
	dbfile := "/tmp/test.db"
	dbHandler, err := repository.NewBoltHandler(
		dbfile,
		logger)
	assert.Nil(t, err)
	assert.NotNil(t, dbHandler)
	defer os.Remove(dbfile)
	// Emulate response of authorization from SYR server
	tokenfield := "message"
	token := "0123456789"
	tsAuth := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"%s\": \"%s\"}\n", tokenfield, token)
	}))
	defer tsAuth.Close()
	// New SYR authorize
	urlauth := tsAuth.URL
	tokenheader := "yadro"
	login := "test"
	password := "test"
	syrAuth, _ := infrastructure.NewSYRAuth(
		urlauth,
		tokenfield,
		tokenheader,
		login,
		password)
	assert.NotNil(t, syrAuth)
	// Emulate response of upload file from SYR server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()
	// New SYR config
	url := ts.URL
	syrConfig, _ := infrastructure.NewSYRConfig(url, "/tmp", "attachment", "")
	assert.NotNil(t, syrConfig)
	// New logerr for catch errors from services is running in goroutines
	var errorbuffer bytes.Buffer
	logerr := log.New(&errorbuffer, "[error] ", log.LstdFlags)
	// New SYR handler for handle upload files to SYR server
	syrHandler, _ := infrastructure.NewSYRHandler(syrConfig, syrAuth, logerr)
	assert.NotNil(t, syrHandler)
	// Create a new handler to invoke tusd functions
	invokeHandler, err := infrastructure.NewTusdInvoke(composer)
	assert.Nil(t, err)
	assert.NotNil(t, invokeHandler)
	// New repository handler to save metadata
	repositoryHandler, err := interfaces.NewDbDataRepo(dbHandler, invokeHandler, "root")
	assert.Nil(t, err)
	assert.NotNil(t, repositoryHandler)
	// New http client to send file to SYR
	httpClientHandler, err := interfaces.NewHTTPClient(syrHandler, logger)
	assert.Nil(t, err)
	assert.NotNil(t, httpClientHandler)
	// New data agent to handle metadata
	dataAgent, err := usecases.NewDataAgent(repositoryHandler, httpClientHandler)
	assert.Nil(t, err)
	assert.NotNil(t, dataAgent)
	// New hooks handler to invoke functions from SYR handler
	hooksHandler, err := interfaces.NewHooksHandler(dataAgent, syrHandler, logger)
	assert.Nil(t, err)
	assert.NotNil(t, hooksHandler)
	// New hooks handler to manage notice from tusd
	hooksTusdHandler, err := infrastructure.NewHooksTusdHandler(composer, hooksHandler, logerr)
	assert.Nil(t, err)
	assert.NotNil(t, hooksTusdHandler)
	// Test with success pipeline
	t.Run("valid test", func(t *testing.T) {
		t.Log("Integration success test")
		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		// Run hooks from tusd
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Start hooks from tusd to usecases
			err := hooksTusdHandler.RunHooks(ctx, tusdHandler)
			assert.Nil(t, err, "unable to run tusd hooks.")
			t.Log("exit hooksTusdHandler.RunHooks()")
			assert.Emptyf(t, errorbuffer.String(), "tusd error: %s", errorbuffer.String())
		}()
		// Run goroutine to handle syr client
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := syrHandler.Run(ctx, nil)
			assert.Nil(t, err, "unable to run syr handler.")
			t.Log("exit syrHandler.Run()")
			assert.Emptyf(t, errorbuffer.String(), "syr error: %s", errorbuffer.String())
		}()
		// Emulate send file to tusd server
		const body = "hello world!"
		res := (&httpTest{
			Method: http.MethodPost,
			ReqHeader: map[string]string{
				"Tus-Resumable":   "1.0.0",
				"Upload-Offset":   "0",
				"Upload-Length":   strconv.Itoa(len(body)),
				"Content-Type":    "application/offset+octet-stream",
				"Upload-Metadata": "data " + data + ", filename d29ybGQ=",
			},
			// check response
			ReqBody: strings.NewReader(body),
			Code:    http.StatusCreated,
			ResHeader: map[string]string{
				"Upload-Offset": strconv.Itoa(len(body)),
			},
		}).Run(tusdHandler, t)
		t.Log("exit request Tusd server")
		time.Sleep(100 * time.Millisecond)
		cancel()
		wg.Wait()
		assert.NotNil(t, res)
		// check log output
		assert.Containsf(t, testbuffer.String(), string("[hooks] [terminate]"), "log: %s", testbuffer.String())
	})
}

type httpTest struct {
	Name string

	Method string
	URL    string

	ReqBody   io.Reader
	ReqHeader map[string]string

	Code      int
	ResBody   string
	ResHeader map[string]string
}

func (test *httpTest) Run(handler http.Handler, t *testing.T) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(test.Method, test.URL, test.ReqBody)
	req.RequestURI = test.URL

	// Add headers
	for key, value := range test.ReqHeader {
		req.Header.Set(key, value)
	}

	req.Host = "tus.io"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != test.Code {
		t.Errorf("Expected %v %s as status code (got %v %s), body: %s", test.Code, http.StatusText(test.Code), w.Code, http.StatusText(w.Code), w.Body.String())
	}

	for key, value := range test.ResHeader {
		header := w.Header().Get(key)

		if value != header {
			t.Errorf("Expected '%s' as '%s' (got '%s')", value, key, header)
		}
	}

	if test.ResBody != "" && w.Body.String() != test.ResBody {
		t.Errorf("Expected '%s' as body (got '%s'", test.ResBody, w.Body.String())
	}

	return w
}
