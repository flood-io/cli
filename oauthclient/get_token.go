package oauthclient

import (
	"github.com/flood-io/cli/config"
	"github.com/pkg/errors"
)

func GetAuthToken(cache config.AuthCache) (token string, err error) {
	switch cache.State() {
	case config.LoggedIn:
		token = cache.Token()
	case config.Expired:
		err = Refresh(cache)
		if err != nil {
			return
		}
		token = cache.Token()
	default:
		err = errors.New("Not logged in.")
	}

	return
}
