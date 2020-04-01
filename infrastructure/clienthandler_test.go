package infrastructure

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewSYRHandler(t *testing.T) {
	cfg, _ := NewSYRConfig("http://", "/tmp/uploads", "attachment", ".bin")
	auth, _ := NewSYRAuth("http://", "message", "yadro", "test", "test")
	logerr := log.New(os.Stdout, "[test] ", log.LstdFlags)

	t.Run("valid New", func(t *testing.T) {
		s, err := NewSYRHandler(
			cfg,
			auth,
			logerr)
		assert.NotNil(t, s)
		assert.Nil(t, err)
	})

	t.Run("invalid Config", func(t *testing.T) {
		s, err := NewSYRHandler(
			nil,
			auth,
			logerr)
		assert.Nil(t, s)
		assert.NotNil(t, err)
	})
	t.Run("invalid Auth", func(t *testing.T) {
		s, err := NewSYRHandler(
			cfg,
			nil,
			logerr)
		assert.Nil(t, s)
		assert.NotNil(t, err)
	})

	t.Run("invalid Errlog", func(t *testing.T) {
		s, err := NewSYRHandler(
			cfg,
			auth,
			nil)
		assert.Nil(t, s)
		assert.NotNil(t, err)
	})
}

func TestSend(t *testing.T) {
	cfg, _ := NewSYRConfig("http://test", "/tmp", "attachment", "")
	auth, _ := NewSYRAuth("http://test", "message", "yadro", "test", "test")
	logerr := log.New(os.Stdout, "[test] ", log.LstdFlags)
	client, _ := NewSYRHandler(cfg, auth, logerr)
	assert.NotNil(t, client)
	id := "0123456789abcdifg"
	name := "log.tar"
	system := "0123456789"
	testfile := filepath.Join(client.cfg.pathfile, id+client.cfg.fileext)
	defer func() {
		if _, err := os.Stat(testfile); err == nil {
			os.Remove(testfile)
		}
	}()
	file, _ := os.Create(testfile)
	defer file.Close()
	t.Run("valid call", func(t *testing.T) {
		err := client.Send(id, name, system)
		assert.Nil(t, err)
		select {
		case d := <-client.chandata:
			assert.Equal(t,
				&data{id, client.auth.token(), name, "http://test/" + system, testfile, client.cfg.fieldform},
				d)
		default:
			assert.Fail(t, "send: empty channel; expect data")
		}
	})
	t.Run("invalid file name", func(t *testing.T) {
		err := client.Send(id, "", system)
		assert.NotNil(t, err)
		select {
		case <-client.chandata:
			assert.Fail(t, "send: channel should be empty")
		default:

		}
	})
	t.Run("invalid file not exist", func(t *testing.T) {
		err := client.Send("test", name, system)
		assert.NotNil(t, err)
		select {
		case <-client.chandata:
			assert.Fail(t, "send: channel should be empty")
		default:

		}
	})
}

func TestRun(t *testing.T) {
	cfg, _ := NewSYRConfig("http://test", "/tmp", "attachment", "")
	auth := new(SYRAuthMock)
	logerr := log.New(os.Stdout, "[test] ", log.LstdFlags) // TODO: need mock for logerr
	client, _ := NewSYRHandler(cfg, auth, logerr)
	assert.NotNil(t, client)
	id := "0123456789abcdifg"
	name := "log.tar"
	//system := "0123456789"
	//testfile := filepath.Join(client.cfg.pathfile, id+client.cfg.fileext)
	t.Run("invalid authorization fail", func(t *testing.T) {
		e := errors.New("fail")
		ctx, cancel := context.WithCancel(context.Background())
		auth.On("authorize", ctx).Return(e)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := client.Run(ctx, nil)
			assert.Nil(t, err)
		}()
		client.chandata <- &data{id, "", name, "", "", ""}
		time.Sleep(10 * time.Millisecond)
		cancel()
		wg.Wait()
		auth.AssertCalled(t, "authorize", ctx)
	})
}

func TestUpload(t *testing.T) {
	cfg, _ := NewSYRConfig("http://test", "/tmp", "attachment", "")
	//auth := new(SYRAuthMock)
	logerr := log.New(os.Stdout, "[test] ", log.LstdFlags) // TODO: need mock for logerr
	//client, _ := NewSYRHandler(cfg, auth, logerr)
	//assert.NotNil(t, client)
	//hooks := new(HooksClientHandlerMock)
	//_ = client.SetHooksHandler(hooks)
	id := "0123456789abcdifg"
	name := "log.tar"
	ctx := context.Background()
	testfile := filepath.Join(cfg.pathfile, id+cfg.fileext)
	token := "0123456789"
	httpclient := http.DefaultClient
	//multi := new(MultipartfileMock)
	response := httptest.NewRecorder().Result()
	t.Run("valid upload", func(t *testing.T) {
		auth := new(SYRAuthMock)
		client, _ := NewSYRHandler(cfg, auth, logerr)
		assert.NotNil(t, client)
		multi := new(MultipartfileMock)
		response.StatusCode = 201
		response.Status = "201 Created"
		auth.On("token").Return(token)
		auth.On("client").Return(httpclient)
		multi.On("uploadMultipartFile",
			ctx,
			httpclient,
			"http://test/",
			token,
			client.cfg.fieldform,
			testfile,
			name,
		).Return(response, nil)
		d := data{id, token, name, "http://test/", testfile, client.cfg.fieldform}
		err := client.upload(ctx, &d, multi)
		auth.AssertCalled(t, "client")
		auth.AssertNotCalled(t, "token")
		multi.AssertCalled(t, "uploadMultipartFile",
			ctx,
			httpclient,
			"http://test/",
			token,
			client.cfg.fieldform,
			testfile,
			name)
		assert.Nil(t, err)
		select {
		case d := <-client.GetChanTerm():
			assert.Equal(t, id, d)
		default:
			assert.Fail(t, "send: empty channel, expect id")
		}
	})
	t.Run("invalid authorize", func(t *testing.T) {
		auth := new(SYRAuthMock)
		client, _ := NewSYRHandler(cfg, auth, logerr)
		assert.NotNil(t, client)
		multi := new(MultipartfileMock)
		response.StatusCode = 401
		response.Status = "401 Unauthorized"
		auth.On("token").Return(token)
		auth.On("client").Return(httpclient)
		auth.On("authorize", ctx).Return(nil)
		multi.On("uploadMultipartFile",
			ctx,
			httpclient,
			"http://test/",
			token,
			client.cfg.fieldform,
			testfile,
			name,
		).Return(response, nil)
		d := data{id, token, name, "http://test/", testfile, client.cfg.fieldform}
		err := client.upload(ctx, &d, multi)
		auth.AssertCalled(t, "client")
		auth.AssertCalled(t, "token")
		auth.AssertCalled(t, "authorize", ctx)
		multi.AssertCalled(t, "uploadMultipartFile",
			ctx,
			httpclient,
			"http://test/",
			token,
			client.cfg.fieldform,
			testfile,
			name)
		assert.Nil(t, err)
		select {
		case <-client.chanterm:
			assert.Fail(t, "send: channel should be empty")
		default:

		}
	})
}

func TestMultipartFile(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: time.Second * 5}
	token := "yadro0123456789"
	key := "attachment"
	path := "/tmp/testupload.bin"
	defer func() {
		if _, err := os.Stat(path); err == nil {
			os.Remove(path)
		}
	}()
	file, _ := os.Create(path)
	defer file.Close()
	name := "test.bin"
	t.Run("valid multiupload", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(201)
		}))
		defer ts.Close()
		url := ts.URL
		res, err := (&multipartfile{}).uploadMultipartFile(
			ctx,
			client,
			url,
			token,
			key,
			path,
			name,
		)
		assert.Nil(t, err)
		assert.Equal(t, 201, res.StatusCode)
	})
	t.Run("invalid authorize fail", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(401)
		}))
		defer ts.Close()
		url := ts.URL
		res, err := (&multipartfile{}).uploadMultipartFile(
			ctx,
			client,
			url,
			token,
			key,
			path,
			name,
		)
		assert.Nil(t, err)
		assert.Equal(t, 401, res.StatusCode)
	})
}
