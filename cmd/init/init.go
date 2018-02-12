package initcmd

import "fmt"

type InitCmd struct {
}

func (i *InitCmd) Run() (err error) {
	fmt.Println("hello init")
	return nil
}
