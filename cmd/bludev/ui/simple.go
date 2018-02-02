package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type SimpleUI struct {
	out io.Writer
}

var _ UI = (*SimpleUI)(nil)

func NewSimpleUI() *SimpleUI {
	return &SimpleUI{
		out: os.Stdout,
	}
}

func (s *SimpleUI) SetStatus(status ...interface{}) {
	fmt.Fprintln(s.out, "===> ", fmt.Sprint(status...))
}

func (s *SimpleUI) Log(msg ...interface{}) {
	fmt.Fprintln(s.out, msg...)
}

func (s *SimpleUI) Logf(format string, args ...interface{}) {
	format = format + "\n"
	fmt.Fprintf(s.out, format, args...)
}

func (s *SimpleUI) HRule() {
	fmt.Fprintln(s.out, strings.Repeat("-", 20))
}

func (s *SimpleUI) AddSlot(label string) {}

func (s *SimpleUI) SetSlot(label, buf string) {
	fmt.Fprintf(s.out, "%s: %s\n", label, buf)
}
