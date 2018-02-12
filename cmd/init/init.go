package initcmd

import (
	"fmt"
	"html/template"
	"os"

	initskeleton "github.com/flood-io/cli/static/init-skeleton"
)

type InitCmd struct {
	Title string
	URL   string
}

func (i *InitCmd) Run() (err error) {
	i.URL = "https://challenge.flood.io"
	i.Title = "titley"

	fmt.Println("hello init")
	testScript, err := initskeleton.Asset("test.ts")
	if err != nil {
		return
	}

	tpl, err := template.New("test.ts").Parse(string(testScript))
	if err != nil {
		return
	}
	err = tpl.Execute(os.Stdout, i)
	return
}
