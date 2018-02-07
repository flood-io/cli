package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func tokenResponse(createdAt time.Time, expiresIn time.Duration) string {
	return fmt.Sprintf(`{
	"access_token":"access-token",
	"token_type":"bearer",
	"expires_in":%d,
	"refresh_token":"refresh-token",
	"scope":"admin",
	"created_at":%d,
	"data": {
		"id":"15099",
		"type":"users",
		"links":{"self":"/api/v3/users/15099"},
		"attributes":{"full-name":"Lachie Cox","company-name":"Flood IO"}
	}
	}`, int64(expiresIn.Seconds()), createdAt.Unix())
}

func writeTokenResponse(t *testing.T, path string, createdAt time.Time, expiresIn time.Duration) {
	// t.Log(tokenResponse(createdAt, expiresIn))
	err := ioutil.WriteFile(path, []byte(tokenResponse(createdAt, expiresIn)), 0600)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_AuthCache_not_logged_in(t *testing.T) {
	as := assert.New(t)
	dir, err := ioutil.TempDir("", "flood-cli-test")
	defer os.RemoveAll(dir)
	as.Nil(err)

	configPath := filepath.Join(dir, "not-even-auth.json")

	c, err := NewFileAuthCache(configPath)
	as.Nil(err)

	state, err := c.StateE()
	as.Equal(NotLoggedIn, state)
	as.Nil(err)
}

func Test_AuthCache_logged_in(t *testing.T) {
	as := assert.New(t)
	dir, err := ioutil.TempDir("", "flood-cli-test")
	defer os.RemoveAll(dir)
	as.Nil(err)

	configPath := filepath.Join(dir, "auth.json")

	now := time.Now()
	writeTokenResponse(t, configPath, now, 10*time.Minute)

	c, err := NewFileAuthCache(configPath)
	as.Nil(err)

	state, err := c.StateE()
	as.Equal(LoggedIn, state)
	as.Nil(err)
}

func Test_AuthCache_expired(t *testing.T) {
	as := assert.New(t)
	dir, err := ioutil.TempDir("", "flood-cli-test")
	defer os.RemoveAll(dir)
	as.Nil(err)

	configPath := filepath.Join(dir, "auth.json")

	now := time.Now().Add(-15 * time.Minute)
	writeTokenResponse(t, configPath, now, 10*time.Minute)

	c, err := NewFileAuthCache(configPath)
	as.Nil(err)

	state, err := c.StateE()
	as.Equal(Expired, state)
	as.Nil(err)
}

func Test_AuthCache_SetAuthData(t *testing.T) {
	as := assert.New(t)
	dir, err := ioutil.TempDir("", "flood-cli-test")
	defer os.RemoveAll(dir)
	as.Nil(err)

	configPath := filepath.Join(dir, "auth.json")
	c, err := NewFileAuthCache(configPath)

	now := time.Now()
	err = c.SetAuthData([]byte(tokenResponse(now, 10*time.Minute)))
	as.Nil(err)

	c2, err := NewFileAuthCache(configPath)

	as.Equal(LoggedIn, c2.State())
	as.Equal("Lachie Cox", c2.FullName())
	as.Equal("access-token", c2.Token())
}

func Test_AuthCache_ReadBeforeLogin(t *testing.T) {
	as := assert.New(t)
	dir, err := ioutil.TempDir("", "flood-cli-test")
	defer os.RemoveAll(dir)
	as.Nil(err)

	configPath := filepath.Join(dir, "auth.json")
	c, err := NewFileAuthCache(configPath)

	as.Nil(err)

	as.Equal("", c.Token())
	as.Equal("", c.FullName())
}
