package initcmd

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	initskeleton "github.com/flood-io/cli/static/init-skeleton"
	input "github.com/tcnksm/go-input"
)

type InitCmd struct {
	Name            string
	URL             string
	WorkingDir      string
	DestinationPath string
	Force           bool
}

func (i *InitCmd) Run(name string) (err error) {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	i.WorkingDir = dir

	// flood init
	if name == "" {
		err = i.interactiveConfig(name)
		if err != nil {
			return err
		}

		i.DestinationPath = i.WorkingDir

		// flood init name
	} else {
		i.Name = name
		i.DestinationPath = filepath.Join(i.WorkingDir, i.Name)

	}

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
	fmt.Println("validating destination path", i.DestinationPath)

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
		if i.Force {
			fmt.Printf("destination isn't empty! (%s)\nYou specified --force so continuing anyway...\n", i.DestinationPath)
			return nil
		} else {
			return fmt.Errorf("destination isn't empty! (%s)", i.DestinationPath)
		}
	}

	return nil
}

func (i *InitCmd) createDestination() (err error) {
	return os.MkdirAll(i.DestinationPath, 0755)
}

func (i *InitCmd) populateDestination() (err error) {
	names := initskeleton.AssetNames()

	for _, name := range names {
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
