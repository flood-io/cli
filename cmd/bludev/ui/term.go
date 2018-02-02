package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/vbauerster/mpb/cwriter"
)

const (
	DefaultFlushInterval = 200 * time.Millisecond
	DefaultWidth         = 80
)

type TermUI struct {
	cw *cwriter.Writer

	width       int
	pushUpLines int

	flushInterval time.Duration

	status, log string
	slots       []*slot
}

var _ UI = (*TermUI)(nil)

type slot struct {
	label string
	buf   string
}

func NewTermUI() *TermUI {
	u := &TermUI{
		cw:            cwriter.New(os.Stdout),
		flushInterval: DefaultFlushInterval,
	}

	u.updateWidth()

	return u
}

func (u *TermUI) updateWidth() {
	u.width, _, _ = cwriter.TermSize()
}

// func (u *TermUI) start() {
// ticker := time.NewTicker(u.flushInterval)

// // TODO winch
// // TODO shutdown

// go func() {
// for {
// select {
// case <-ticker.C:
// u.render()
// err := u.cw.Flush()
// if err != nil {
// fmt.Fprintf(os.Stderr, "problem flushing: %v", err)
// }
// }
// }
// }()
// }

func (u *TermUI) Flush() {
	u.render()
	err := u.cw.Flush()
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem flushing: %v", err)
	}
}

func (u *TermUI) render() {
	for i := 0; i < u.pushUpLines; i++ {
		fmt.Fprintln(os.Stdout, "") // push the top lines into scrollback
	}
	u.pushUpLines = 0

	u.cw.WriteString(u.log)
	u.cw.WriteString("\n")
	u.cw.WriteString(strings.Repeat("=", u.width))
	u.cw.WriteString(u.status)
	u.cw.WriteString("\n")
	u.renderSlots()
	u.cw.WriteString("\n")
}

func (u *TermUI) renderSlots() {
	// widest:=0
	// for _,slot:=range u

	for _, slot := range u.slots {
		u.cw.WriteString(fmt.Sprintf("%s: %s\n", slot.label, slot.buf))
	}
}

func (u *TermUI) pushUp() {
	u.pushUpLines++
}

func (u *TermUI) SetStatusAndLog(status ...interface{}) {
	u.status = fmt.Sprintln("[", fmt.Sprint(status...), "]")
	u.Log(status...)
}

func (u *TermUI) SetStatus(status ...interface{}) {
	u.status = fmt.Sprintln("[", fmt.Sprint(status...), "]")
	u.Flush()
}

func (u *TermUI) setLog(log string) {
	u.pushUp()
	u.log = log
}

func (u *TermUI) Log(msg ...interface{}) {
	u.setLog(fmt.Sprintln(msg...))
	u.Flush()
}

func (u *TermUI) Logf(format string, args ...interface{}) {
	u.setLog(fmt.Sprintf(format, args...))
	u.Flush()
}

func (u *TermUI) HRule() {
	u.setLog(strings.Repeat("-", u.width))
	u.Flush()
}

func (u *TermUI) AddSlot(label string) {
	u.slots = append(u.slots, &slot{label, ""})
	u.Flush()
}

func (u *TermUI) SetSlot(label, buf string) {
	for _, slot := range u.slots {
		if slot.label == label {
			slot.buf = buf
			return
		}
	}
}
