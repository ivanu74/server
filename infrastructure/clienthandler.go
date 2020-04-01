package infrastructure

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

type authHandler interface {
	authorize(ctx context.Context) error
	token() string
	client() *http.Client
}

type data struct {
	id        string
	token     string
	filename  string
	urlpath   string
	filepath  string
	fieldform string
}

type SYRConfig struct {
	url       string
	pathfile  string
	fieldform string
	fileext   string
}

// NewSYRConfig - create new instance of SYR configuration
func NewSYRConfig(
	url string,
	path string,
	fieldform string,
	fileext string) (*SYRConfig, error) {

	if url == "" ||
		path == "" ||
		fieldform == "" {
		return nil, errors.New("[clienthandler] [new config] bad argument")
	}
	return &SYRConfig{
		url,
		path,
		fieldform,
		fileext}, nil
}

// SYRHandler - implements interface clientHandler from client(interfaces)
type SYRHandler struct {
	cfg      *SYRConfig
	chandata chan *data
	chanterm chan string
	stderr   logger
	auth     authHandler
}

// NewSYRHandler - create new instance of SYRHandler for HTTPclient
func NewSYRHandler(cfg *SYRConfig, auth authHandler, errlog logger) (*SYRHandler, error) {
	if cfg == nil || errlog == nil || auth == nil {
		return nil, errors.New("[clienthandler] [new handler] bad argument")
	}
	return &SYRHandler{
		cfg,
		make(chan *data, 10),
		make(chan string, 1),
		errlog,
		auth,
	}, nil
}

// Send - implement func to send file to SYR server
func (client *SYRHandler) Send(id string, name string, system string) error {

	var fieldform = client.cfg.fieldform
	u, err := url.Parse(client.cfg.url)
	if err != nil {
		return errors.Wrap(err, "[client] [send]")
	}
	u.Path = path.Join(u.Path, system)
	pathfile := filepath.Join(client.cfg.pathfile, id+client.cfg.fileext)
	if _, err := os.Stat(pathfile); os.IsNotExist(err) {
		return errors.Wrap(err, "[client] [send]")
	}
	var namefile = name
	if namefile == "" {
		return errors.New("[client] [send] filename is empty")
	}
	client.chandata <- &data{id, client.auth.token(), namefile, u.String(), pathfile, fieldform}
	// client.errlog.Printf("[client] [send]: id = %s; url = %s\n", id, u.String())
	return nil
}

// Run - goroutine to handle connect to SYR
func (client *SYRHandler) Run(ctx context.Context, multi imultipartfile) error {
	if multi == nil {
		multi = &multipartfile{}
	}
	for {
		select {
		case data := <-client.chandata:
			err := client.upload(ctx, data, multi)
			if err != nil {
				client.stderr.Printf("[client] [run]: %s\n", err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// GetChanTerm - return channel of id(string) to delete file with that id
func (client *SYRHandler) GetChanTerm() chan string {
	return client.chanterm
}

type imultipartfile interface {
	uploadMultipartFile(ctx context.Context, client *http.Client, url, token, key, path string, name string) (*http.Response, error)
}

func (client *SYRHandler) upload(ctx context.Context, data *data, multi imultipartfile) error {
	if data.token == "" {
		if err := client.auth.authorize(ctx); err != nil {
			return errors.Wrap(err, "[client] [upload]: ")
		}
		data.token = client.auth.token()
	}
	res, err := multi.uploadMultipartFile(
		ctx,
		client.auth.client(),
		data.urlpath,
		data.token,
		data.fieldform,
		data.filepath,
		data.filename,
	)
	if err != nil {
		return errors.Wrap(err, "[client] [upload]: ")
	}
	if res.StatusCode == 401 {
		if err := client.auth.authorize(ctx); err != nil {
			return errors.Wrap(err, "[client] [upload]: ")
		}
		data.token = client.auth.token()
		client.chandata <- data
		return nil
	}
	if res.StatusCode != 201 {
		return errors.Errorf("[client] [upload]: [syr] status code = %s; request %v ", res.Status, res.Request)
	}
	if res.StatusCode == 201 {
		select {
		case client.chanterm <- data.id:
		default:
			return errors.Errorf("[client] [upload]: file isn't deleted: %s ", data.id)
		}
	}
	return nil
}

type multipartfile struct{}

// TODO: implement test with http-server, test when authorize is fail
func (m *multipartfile) uploadMultipartFile(ctx context.Context, client *http.Client, url, token, key, path string, name string) (*http.Response, error) {
	body, writer := io.Pipe()
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	multiwriter := multipart.NewWriter(writer)
	req.Header.Add("Content-Type", multiwriter.FormDataContentType())
	req.Header.Add("Authorization", token)
	errchan := make(chan error)
	go func() {
		defer close(errchan)
		defer writer.Close()
		defer multiwriter.Close()
		w, err := multiwriter.CreateFormFile(key, name)
		if err != nil {
			errchan <- err
			return
		}
		in, err := os.Open(path)
		if err != nil {
			errchan <- err
			return
		}
		defer in.Close()
		if written, err := io.Copy(w, in); err != nil {
			errchan <- errors.Wrapf(err, "(%d bytes written)", written)
			return
		}
		if err := multiwriter.Close(); err != nil {
			errchan <- err
			return
		}
	}()
	resp, err := client.Do(req.WithContext(ctx))
	multierr := <-errchan
	if err != nil || multierr != nil {
		return resp, errors.Errorf("http error: %v, multipart error: %v", err, multierr)
	}
	return resp, nil
}
