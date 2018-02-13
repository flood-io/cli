package oauthclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/flood-io/cli/config"
	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

type refreshPayload struct {
	RefreshToken string `url:"refresh_token"`
	GrantType    string `url:"grant_type"`
}

func LoginRefresh(force bool, cache config.AuthCache) (err error) {
	switch cache.State() {
	case config.LoggedIn:
		if !force {
			fmt.Printf("You're already signed in as %s\n", cache.FullName())
			return
		}
	case config.Expired:
		fmt.Println("You're logged in, but your auth token has expired; Attempting to refresh")
	default:
		return errors.New("Unable to refresh. Please flood login first")
	}

	err = Refresh(cache)
	if err != nil {
		return
	}

	fmt.Printf("Successfully refreshed auth token for %s!\n", cache.FullName())
	return
}

func Refresh(cache config.AuthCache) (err error) {
	switch cache.State() {
	case config.LoggedIn, config.Expired:
	default:
		return errors.New("assertion failed: cache state != logged in or expired")
	}

	p := refreshPayload{cache.RefreshToken(), "refresh_token"}
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

	err = cache.SetAuthData(b)
	if err != nil {
		return errors.Wrapf(err, "unable to set cache auth data from JSON response: (body=%s)", string(b))
	}

	if cache.State() != config.LoggedIn {
		return errors.New("Assertion failed: cache state != logged in")
	}

	return
}
