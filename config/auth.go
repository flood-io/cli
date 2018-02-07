package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

type AuthState int

const (
	Undefined AuthState = iota
	NotLoggedIn
	Expired
	LoggedIn
)

type AuthCache interface {
	State() AuthState
	StateE() (AuthState, error)
	SetAuthData([]byte) error
	FullName() string
	Token() string
}

type FileAuthCache struct {
	path  string
	state AuthState
	data  *AuthTokenData
}

var _ AuthCache = (*FileAuthCache)(nil)

type AuthTokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int64  `json:"created_at"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Data         struct {
		Id         string `json:"id"`
		Attributes struct {
			FullName string `json:"full-name"`
		}
	}
}

var DefaultAuthCacheFilePath string

func init() {
	DefaultAuthCacheFilePath = filepath.Join(os.Getenv("HOME"), ".cache/flood-io/auth.json")
}

func NewFileAuthCache(path string) (c *FileAuthCache, err error) {
	c = &FileAuthCache{
		path: path,
		data: nil,
	}

	return
}

func (c *FileAuthCache) State() AuthState {
	state, err := c.StateE()
	if err != nil {
		return Undefined
	}

	return state
}

func (c *FileAuthCache) StateE() (AuthState, error) {
	if c.data == nil {
		err := c.readData()
		if os.IsNotExist(err) {
			return NotLoggedIn, nil
		} else if err != nil {
			return Undefined, err
		}
	}

	if c.data.AccessToken == "" {
		return Undefined, errors.New("no access token")
	}

	createdAt := time.Unix(c.data.CreatedAt, 0)
	var expiresIn time.Duration = time.Duration(c.data.ExpiresIn) * time.Second

	if time.Since(createdAt) > expiresIn {
		return Expired, nil
	}

	return LoggedIn, nil
}

func (c *FileAuthCache) readData() (err error) {
	dataBytes, err := ioutil.ReadFile(c.path)
	if err != nil {
		return
	}

	c.data = &AuthTokenData{}
	err = json.Unmarshal(dataBytes, c.data)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal auth cache data")
	}

	return
}

func (c *FileAuthCache) SetAuthData(dataBytes []byte) (err error) {
	var data AuthTokenData
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal auth cache data during set")
	}

	err = ioutil.WriteFile(c.path, dataBytes, 0600)
	if err != nil {
		return errors.Wrapf(err, "unable to write auth cache to '%s'", c.path)
	}

	return
}

func (c *FileAuthCache) FullName() string {
	if c.State() != LoggedIn {
		return ""
	} else {
		return c.data.Data.Attributes.FullName
	}
}

func (c *FileAuthCache) Token() string {
	if c.State() != LoggedIn {
		return ""
	} else {
		return c.data.AccessToken
	}
}
