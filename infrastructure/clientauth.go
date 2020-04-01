package infrastructure

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// User - struct of user from SYR for authorization
type User struct {
	Login        string   `json:"login" bson:"login"`
	Password     string   `json:"password" bson:"password"`
	FirstName    string   `json:"first_name" bson:"first_name"`
	LastName     string   `json:"last_name" bson:"last_name"`
	Email        string   `json:"email" bson:"email"`
	AccessList   []string `json:"access_list" bson:"access_list"`
	IsPrivileged bool     `json:"is_privileged" bson:"is_privileged"`
	LdapGroups   []string `json:"ldap_groups" bson:"ldap_groups"`
}

type SYRAuth struct {
	urlauth     string
	tokenfield  string
	tokenheader string
	login       string
	password    string
	_token      string
	_client     *http.Client
}

// NewSYRAuth - create new instance of SYR configuration
func NewSYRAuth(
	urlauth string,
	tokenfield string,
	tokenheader string,
	login string,
	password string) (*SYRAuth, error) {

	if urlauth == "" ||
		tokenfield == "" ||
		tokenheader == "" ||
		login == "" ||
		password == "" {
		return nil, errors.New("[clientauth] [new] bad argument")
	}
	return &SYRAuth{
		urlauth,
		tokenfield,
		tokenheader,
		login,
		password,
		"",
		&http.Client{Timeout: time.Second * 30},
	}, nil
}

func (client *SYRAuth) authorize(ctx context.Context) error {

	user := User{
		Login:    client.login,
		Password: client.password,
	}
	body, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "[client] [authorize] user bad json format")
	}
	url := client.urlauth
	payload := strings.NewReader(string(body))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	res, err := client._client.Do(req.WithContext(ctx))
	if err != nil {
		return errors.Wrap(err, "[client] [authorize] request")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "[client] [authorize] [syr] status code = %s ", res.Status)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "[client] [authorize] read res.Body")
	}
	var tokenmap map[string]string
	if err := json.Unmarshal(body, &tokenmap); err != nil {
		return errors.Wrapf(err, "[client] [authorize] unmarshal res.Body: %s", string(body))
	}
	token, ok := tokenmap[client.tokenfield]
	if !ok {
		return errors.New("[client] [authorize] bad token")
	}
	client._token = client.tokenheader + token
	return nil
}

func (client *SYRAuth) token() string {
	return client._token
}
func (client *SYRAuth) client() *http.Client {
	return client._client
}
