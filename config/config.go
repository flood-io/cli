package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config interface {
	HasData() bool

	HasAuthData() bool
	AuthToken() string
	AuthFullName() string

	SetAuthTokenData([]byte) error
}

var _ Config = (*FileConfig)(nil)

var DefaultConfigRef *FileConfig
var DefaultConfigFilePath string

func init() {
	DefaultConfigFilePath = filepath.Join(os.Getenv("HOME"), ".config/flood-io/config.json")
}

func DefaultConfig() *FileConfig {
	if DefaultConfigRef == nil {
		DefaultConfigRef = LoadFileConfig(DefaultConfigFilePath)
	}

	return DefaultConfigRef
}

type AuthTokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int    `json:"created_at"`
	Data         struct {
		Id         string `json:"id"`
		Attributes struct {
			FullName string `json:"full-name"`
		}
	}
}

type ConfigData struct {
	Null          bool `json:"-"`
	AuthTokenData AuthTokenData
	APIToken      string
}

var NullConfigData = &ConfigData{
	Null: true,
}

func NewConfigData() *ConfigData {
	return &ConfigData{}
}

type FileConfig struct {
	Path string
	Data *ConfigData
}

func LoadFileConfig(path string) (f *FileConfig) {
	f = &FileConfig{Path: path, Data: NullConfigData}

	f.Read()

	return
}

func (f *FileConfig) Read() error {
	data, err := ioutil.ReadFile(f.Path)
	if err != nil {
		f.Data = NullConfigData
		return err
	}

	f.Data = NewConfigData()
	err = json.Unmarshal(data, f.Data)
	if err != nil {
		f.Data = NullConfigData
		return err
	}

	return nil
}

func (f *FileConfig) Write() (err error) {
	if !f.HasData() {
		return
	}

	configHome := filepath.Dir(f.Path)
	err = os.MkdirAll(configHome, 0700)
	if err != nil {
		return
	}

	b, err := json.Marshal(f.Data)
	if err != nil {
		return
	}

	return ioutil.WriteFile(f.Path, b, 0600)
}

func (f *FileConfig) HasData() bool {
	return !f.Data.Null
}

func (f *FileConfig) HasAuthData() bool {
	return f.Data.AuthTokenData.AccessToken != ""
}

func (f *FileConfig) APIToken() string {
	return f.Data.APIToken
}

func (f *FileConfig) AuthToken() string {
	return f.Data.AuthTokenData.AccessToken
}

func (f *FileConfig) AuthFullName() string {
	return f.Data.AuthTokenData.Data.Attributes.FullName
}

func (f *FileConfig) SetAuthTokenData(b []byte) (err error) {
	data := f.ReplaceData()
	err = json.Unmarshal(b, &data.AuthTokenData)
	if err != nil {
		return
	}
	return f.Write()
}

func (f *FileConfig) ReplaceData() (data *ConfigData) {
	if f.Data.Null {
		f.Data = NewConfigData()
	}
	return f.Data
}
