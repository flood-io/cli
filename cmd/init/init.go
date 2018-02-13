package initcmd

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	initskeleton "github.com/flood-io/cli/static/init-skeleton"
	au "github.com/logrusorgru/aurora"
	input "github.com/tcnksm/go-input"
)

type InitCmd struct {
	Name            string
	URL             string
	WorkingDir      string
	DestinationPath string
	Force           bool
	Validated       bool
}

func (i *InitCmd) Run(name string) (err error) {
	fmt.Println(au.Blue("~~~ Flood Chrome"), au.Green("Init"), au.Blue("~~~"))

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	i.WorkingDir = dir

	// flood init
	if name == "" {
		i.DestinationPath = i.WorkingDir

		err = i.validateDestination()
		if err != nil {
			return err
		}

		err = i.interactiveConfig(name)
		if err != nil {
			return err
		}

		// flood init name
	} else {
		i.Name = name
		i.DestinationPath = filepath.Join(i.WorkingDir, i.Name)

	}

	fmt.Println("Initialising a test script project named", au.Brown(i.Name), "in folder", au.Gray(i.DestinationPath))

	if i.destinationExists() {
		err = i.validateDestination()
		if err != nil {
			return err
		}
	} else {
		err = i.createDestination()
		if err != nil {
			return err
		}
	}

	err = i.populateDestination()
	if err != nil {
		return err
	}

	fmt.Println(au.Green("done"))
	fmt.Println()
	fmt.Println("Next steps:")
	if i.DestinationPath != i.WorkingDir {
		fmt.Println()
		fmt.Println("    cd", i.DestinationPath)
	}
	fmt.Println(`
    # install node packages using yarn
    yarn
    # or npm
    npm i

    # edit test.ts using e.g. vs code:
    code test.ts

    # verify test.ts against Flood Chrome using flood verify:
    flood verify test.ts
`)

	return
}

func (i *InitCmd) interactiveConfig(name string) (err error) {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	if name == "" {
		name = filepath.Base(i.WorkingDir)
	}

	query := "Project name"
	name, err = ui.Ask(query, &input.Options{
		Default:  name,
		Required: true,
		Loop:     true,
	})
	if err != nil {
		return
	}

	i.Name = name

	query = "Test URL"
	url, err := ui.Ask(query, &input.Options{
		Default:  i.URL,
		Required: true,
		Loop:     true,
	})
	if err != nil {
		return
	}

	i.Name = name
	i.URL = url

	return
}

func (i *InitCmd) destinationExists() bool {
	if _, err := os.Stat(i.DestinationPath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func (i *InitCmd) validateDestination() (err error) {
	if i.Validated {
		return
	}

	dir, err := os.Open(i.DestinationPath)
	if err != nil {
		return
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return
	}

	if len(names) > 0 {
		fmt.Println()
		errMsg := fmt.Sprint(au.Red("destination"), au.Gray(i.DestinationPath), au.Red("isn't empty!"))
		if i.Force {
			fmt.Println(errMsg)
			fmt.Println("However you specified", au.Blue("--force"), "so continuing anyway")
			return nil
		} else {
			return fmt.Errorf(errMsg)
		}
	}

	i.Validated = true

	return nil
}

func (i *InitCmd) createDestination() (err error) {
	return os.MkdirAll(i.DestinationPath, 0755)
}

func (i *InitCmd) populateDestination() (err error) {
	names := initskeleton.AssetNames()

	fmt.Println()
	fmt.Println(au.Green("adding files"))
	for _, name := range names {
		fmt.Println(" -", au.Green(name))
		if name == "test.ts" {
			i.reifyTemplate(name)
		} else {
			i.writeAsset(name)
		}
	}
	return nil
}

func (i *InitCmd) writeAsset(name string) (err error) {
	asset, err := initskeleton.Asset(name)
	if err != nil {
		return
	}

	destPath := filepath.Join(i.DestinationPath, name)
	return ioutil.WriteFile(destPath, asset, 0644)
}

func (i *InitCmd) reifyTemplate(name string) (err error) {
	templateSource, err := initskeleton.Asset(name)
	if err != nil {
		return
	}

	destPath := filepath.Join(i.DestinationPath, name)
	out, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer out.Close()

	tpl, err := template.New(name).Parse(string(templateSource))
	if err != nil {
		return
	}
	return tpl.Execute(out, i)
}
