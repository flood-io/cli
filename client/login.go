package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/tcnksm/go-input"
)

type payload struct {
	Username  string `url:"username"`
	Password  string `url:"password"`
	GrantType string `url:"grant_type"`
}

type PasswordTokenResponse struct {
	AccessToken string `json:"access_token"`
	CreatedAt   int16  `json:"created_at"`
	Data        TokenResponseData
}

type TokenResponseData struct {
	Id         string `json:"id"`
	Attributes struct {
		FullName string `json:"full-name"`
	}
}

func Login() error {
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
		return err
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
		fmt.Println(err.Error())
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	responsePayload := &PasswordTokenResponse{}
	json.Unmarshal(b, responsePayload)

	if resp.StatusCode == 200 {
		fmt.Printf("Welcome back %s!", responsePayload.Data.Attributes.FullName)
	} else {
		fmt.Println("Authentication failed")
	}

	err = ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".flood.json"), b, 0600)

	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
