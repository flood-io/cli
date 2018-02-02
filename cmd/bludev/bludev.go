package bludev

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/flood-io/cli/cmd/bludev/ui"
	"github.com/flood-io/cli/config"
	pb "github.com/flood-io/go-wrenches/floodchrome"
	fcClient "github.com/flood-io/go-wrenches/floodchrome/client"
)

type BLUDev struct {
	LaunchDevtoolsMode bool
	FloodChromeChannel string
}

func (b *BLUDev) floodchromeClient() (client *fcClient.Client, err error) {
	host := "http://localhost:5000"

	token := config.DefaultConfig().APIToken()
	client = fcClient.New(host, token)
	return
}

type stateFn func(msg *pb.TestResult) (stateFn, error)
type state struct {
	name      string
	state     stateFn
	stepLabel string
	ui        ui.UI
}

func (b *BLUDev) Run(scriptFile string) (err error) {
	ui := ui.NewSimpleUI()

	ui.SetStatus("Flood Chrome Dev Mode")

	ui.Log("Flood Chrome Dev Mode - starting run")
	ui.Log("flood chrome channel: ", b.FloodChromeChannel)
	ui.Log("script file:", scriptFile)

	f, err := os.Open(scriptFile)
	if err != nil {
		return
	}

	scriptBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	client, err := b.floodchromeClient()
	if err != nil {
		return
	}

	test := &pb.TestRequest{
		Script:             string(scriptBytes),
		FloodChromeVersion: b.FloodChromeChannel,
	}

	state := &state{
		name: "testy",
		ui:   ui,
	}
	state.state = state.initialState

	stream, err := client.Run(context.Background(), test)
	if err != nil {
		log.Fatalf("%v.Run(_) = _, %v", client, err)
	}
	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream %v.Run(_) = _, %v", client, err)
		}

		// state.dumpLife(result)

		// fmt.Printf("result = %+v\n", result)
		// fmt.Printf("result = %T\n", result)
		// fmt.Println("result", result.String())

		err = state.next(result)
		if err != nil {
			log.Fatalf("error handling message %v", err)
		}
	}

	return
}

func writeNetworkTrace(t *pb.TestResult_Trace_Network) (err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	har := filepath.Join(cwd, "har.json")
	fmt.Println("[trace] writing network trace", har)

	f, err := os.Create(har)
	if err != nil {
		return
	}

	_, err = f.WriteString(t.Har)
	if err != nil {
		return
	}

	return
}

func (s *state) next(msg *pb.TestResult) (err error) {
	if s.state == nil {
		panic("next state is nil")
	}
	s.state, err = s.state(msg)
	return
}

func (s *state) exhaustNothing(msg *pb.TestResult) (next stateFn, err error) {
	next = s.exhaustNothing
	return
}
func (s *state) nothingDump(msg *pb.TestResult) (next stateFn, err error) {
	next = s.nothingDump
	dump(msg)
	return
}
func (s *state) exhaustLog(msg *pb.TestResult) (next stateFn, err error) {
	next = s.exhaustLog
	s.dumpLog(msg)
	return
}

func (s *state) setStepStatus(msg ...interface{}) {
	s.ui.SetSlot(s.stepLabel, fmt.Sprint(msg...))
}

func dump(msg *pb.TestResult) {
	fmt.Printf("msg = %+v\n", msg)
	fmt.Printf("msg = %T\n", msg)
	fmt.Println("msg", msg.String())
}

func (s *state) dumpLife(msg *pb.TestResult) {
	if lifecycleM := msg.GetLifecycle(); lifecycleM != nil {
		s.ui.Logf("[ life] [%10s] %s\n", lifecycleM.Event.String(), msg.Label)
	}
}

func (s *state) dumpLog(msg *pb.TestResult) {
	if logM := msg.GetServerLog(); logM != nil {
		// s.ui.Logf("[server-%5s] %+v\n", logM.Level, msg.Message)
	} else if logM := msg.GetScriptLog(); logM != nil {
		s.ui.Logf("[script-%5s] %+v\n", logM.Level, msg.Message)
	}
}

func matchLifecycle(msg *pb.TestResult, event pb.TestResult_Lifecycle_Event) bool {
	lifecycleM := msg.GetLifecycle()
	return lifecycleM != nil && lifecycleM.Event == event
}

func (s *state) initialState(msg *pb.TestResult) (next stateFn, err error) {
	s.ui.SetStatus("~~~~ Flood Chrome ~~~~")

	// dump(msg)

	return s.awaitTest(msg)
}

func (s *state) awaitTest(msg *pb.TestResult) (next stateFn, err error) {
	next = s.awaitTest

	if msg.Label == "proxy" && matchLifecycle(msg, pb.TestResult_Lifecycle_Setup) {
		s.ui.SetStatusAndLog("Proxy starting")

	} else if msg.Label == "floodchrome" && matchLifecycle(msg, pb.TestResult_Lifecycle_Setup) {
		s.ui.SetStatusAndLog("floodchrome starting")

	} else if msg.Label == "test" && matchLifecycle(msg, pb.TestResult_Lifecycle_BeforeTest) {
		s.ui.SetStatusAndLog("Test starting")
		next = s.awaitPlan

	} else {
		s.dumpLog(msg)
	}

	return
}

func (s *state) awaitPlan(msg *pb.TestResult) (next stateFn, err error) {
	next = s.awaitPlan

	plan := msg.GetPlan()
	if plan != nil {
		// fmt.Printf("plan = %+v\n", plan)
		for _, step := range plan.Steps {
			s.ui.AddSlot(step)
		}
		next = s.awaitNext
	} else {
		s.dumpLog(msg)
	}

	return
}

func (s *state) awaitNext(msg *pb.TestResult) (next stateFn, err error) {
	next = s.awaitNext

	if matchLifecycle(msg, pb.TestResult_Lifecycle_BeforeStep) {
		s.ui.SetStatus("Running step ", msg.Label)
		s.ui.HRule()
		s.ui.Log("Running step ", msg.Label)

		s.stepLabel = msg.Label
		s.setStepStatus("running")
		next = s.handleStep
	} else if matchLifecycle(msg, pb.TestResult_Lifecycle_TestSucceeded) {
		s.ui.SetStatus("test succeeded")
	} else if matchLifecycle(msg, pb.TestResult_Lifecycle_TestFailed) {
		s.ui.SetStatus("test failed")
	} else if matchLifecycle(msg, pb.TestResult_Lifecycle_AfterTest) {
		next = s.awaitExit
	} else if matchLifecycle(msg, pb.TestResult_Lifecycle_StepSkipped) {
		s.ui.Log("   > skipped: ", msg.Label)
		next = s.awaitNext
	} else {
		s.dumpLog(msg)
	}
	return
}

func (s *state) handleStep(msg *pb.TestResult) (next stateFn, err error) {
	next = s.handleStep

	lifecycle := msg.GetLifecycle()
	testError := msg.GetError()
	if lifecycle != nil {
		next, err = s.handleStepLifecycle(msg, lifecycle)
	} else if testError != nil {
		next, err = s.handleStepError(msg, testError)
	} else {
		s.dumpLog(msg)
	}
	return
}

func (s *state) handleStepLifecycle(msg *pb.TestResult, lifecycle *pb.TestResult_Lifecycle) (next stateFn, err error) {
	next = s.handleStep

	switch lifecycle.Event {
	case pb.TestResult_Lifecycle_StepFailed:
		s.setStepStatus("failed")
		s.ui.Log("=!=> failed")
	case pb.TestResult_Lifecycle_StepSucceeded:
		s.setStepStatus("succeeded")
		s.ui.Log("=v=> succeeded")
	case pb.TestResult_Lifecycle_StepSkipped:
		s.setStepStatus("skipped")
		s.ui.Log("= => skipped")
	case pb.TestResult_Lifecycle_BeforeStepAction:
		s.setStepStatus("running (", msg.Label, ")")
		s.ui.Log("---> action", msg.Label)
	case pb.TestResult_Lifecycle_AfterStep:
		s.ui.SetStatus("Finished step ", msg.Label)
		next = s.awaitNext
	}

	return
}

func (s *state) handleStepError(msg *pb.TestResult, testError *pb.TestResult_Error) (next stateFn, err error) {
	next = s.handleStep

	s.ui.Log("-!-> error !")
	s.ui.Logf("message = %+v\n", msg.Message)
	s.ui.Logf("testError = %+v\n", testError.Message)

	return
}

func (s *state) awaitExit(msg *pb.TestResult) (next stateFn, err error) {
	next = s.awaitExit

	if msg.Label == "floodchrome" && matchLifecycle(msg, pb.TestResult_Lifecycle_Exit) {
		next = s.exhaustNothing
	} else {
		s.dumpLog(msg)
	}

	return
}
