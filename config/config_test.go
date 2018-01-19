package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const tokenResponse = `{
	"access_token":"access-token",
	"token_type":"bearer",
	"expires_in":1209600,
	"refresh_token":"refresh-token",
	"scope":"admin",
	"created_at":1515462909,
	"data": {
		"id":"15099",
		"type":"users",
		"links":{"self":"/api/v3/users/15099"},
		"attributes":{"full-name":"Lachie Cox","company-name":"Flood IO"}
	}
}`

func Test_FileConfig(t *testing.T) {
	as := assert.New(t)
	dir, err := ioutil.TempDir("", "flood-cli-test")
	defer os.RemoveAll(dir)
	as.Nil(err)

	configPath := filepath.Join(dir, "config.json")

	c := LoadFileConfig(configPath)

	as.False(c.HasAuthData())
	as.Equal("", c.AuthToken())
	as.Equal("", c.AuthFullName())

	err = c.SetAuthTokenData([]byte(tokenResponse))
	as.Nil(err)

	b, err := ioutil.ReadFile(c.Path)
	as.Nil(err)
	fmt.Printf("string(b) = %+v\n", string(b))

	as.Contains(string(b), "access-token")

	c2 := LoadFileConfig(configPath)
	as.Equal("access-token", c2.AuthToken())
	as.Equal("Lachie Cox", c.AuthFullName())
}
