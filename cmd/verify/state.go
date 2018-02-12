package verify

import (
	"fmt"

	pb "github.com/flood-io/go-wrenches/floodchrome"
)

type stateFn func(msg *pb.TestResult) stateFn
type state struct {
	Test         Test
	state        stateFn
	currentError error
}

func (s *state) next(msg *pb.TestResult) (err error) {
	if s.state == nil {
		panic("next state is nil")
	}
	s.state = s.state(msg)
	return s.err()
}

func (s *state) err() (err error) {
	if s.currentError != nil {
		return s.currentError
	} else if err := s.Test.Error(); err != nil {
		return err
	}

	return nil
}

func (s *state) exhaustNothing(msg *pb.TestResult) (next stateFn) {
	next = s.exhaustNothing
	return
}
func (s *state) nothingDump(msg *pb.TestResult) (next stateFn) {
	next = s.nothingDump
	dump(msg)
	return
}
func (s *state) exhaustLog(msg *pb.TestResult) (next stateFn) {
	next = s.exhaustLog
	s.dumpLog(msg)
	return
}

func dump(msg *pb.TestResult) {
	fmt.Printf("msg = %+v\n", msg)
	fmt.Printf("msg = %T\n", msg)
	fmt.Println("msg", msg.String())
}

func (s *state) dumpLife(msg *pb.TestResult) {
	if lifecycleM := msg.GetLifecycle(); lifecycleM != nil {
		fmt.Printf("[ life] [%10s] %s\n", lifecycleM.Event.String(), msg.Label)
	}
}

func (s *state) dumpLog(msg *pb.TestResult) {
	// fmt.Printf("msg = %+v\n", msg)
	if logM := msg.GetServerLog(); logM != nil {
		// TODO make a switch --server-logs
		// s.ui.Logf("[server-%5s] %+v\n", logM.Level, msg.Message)
	} else if logM := msg.GetScriptLog(); logM != nil {
		s.Test.ScriptLog(logM.Level, msg.Message)
	}
}
