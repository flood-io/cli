package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/flood-io/cli/config"
	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	input "github.com/tcnksm/go-input"
)

type payload struct {
	Username  string `url:"username"`
	Password  string `url:"password"`
	GrantType string `url:"grant_type"`
}

/*
 * {
 *   "access_token": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
 *   "token_type": "bearer",
 *   "expires_in": 1209600,
 *   "refresh_token": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
 *   "scope": "admin",
 *   "created_at": 1515462739,
 *   "data": {
 *     "id": "1",
 *     "type": "users",
 *     "links": {
 *       "self": "/api/v3/users/1"
 *     },
 *     "attributes": {
 *       "full-name": "Lachie Cox",
 *       "company-name": "Flood IO"
 *     }
 *   }
 * }
 */

func Login() (err error) {
	config := config.DefaultConfig()

	if config.HasAuthData() {
		fmt.Printf("You're already signed in as %s!\n", config.AuthFullName())
		return
	}

	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	username, err := ui.Ask("What's your username:", &input.Options{
		Default:  "",
		Required: true,
		Loop:     true,
	})
	if err != nil {
		return
	}

	password, err := ui.Ask("What's your password (masked):", &input.Options{
		Default:     "",
		Required:    true,
		Loop:        true,
		Mask:        true,
		MaskDefault: true,
	})

	p := payload{username, password, "password"}
	v, _ := query.Values(p)

	body := strings.NewReader(v.Encode())
	resp, err := http.DefaultClient.Post("https://flood.io/oauth/token", "application/x-www-form-urlencoded", body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		var responseBody []byte
		responseBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		return fmt.Errorf("Authentication failed: %d %s", resp.StatusCode, string(responseBody))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = config.SetAuthTokenData(b)
	if err != nil {
		return errors.Wrapf(err, "unable to set config from JSON response: (body=%s)", string(b))
	}

	fmt.Printf("Welcome back %s!\n", config.AuthFullName())
	return
}
