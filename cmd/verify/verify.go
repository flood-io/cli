package verify

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/flood-io/cli/cmd/verify/ui"
	pb "github.com/flood-io/go-wrenches/floodchrome"
	fcClient "github.com/flood-io/go-wrenches/floodchrome/client"
	"github.com/pkg/errors"
)

type VerifyCmd struct {
	LaunchDevtoolsMode bool
	FloodChromeChannel string

	Host    string
	DevMode string
}

func (b *VerifyCmd) floodchromeClient(token string) (client *fcClient.Client, err error) {
	client = fcClient.New(b.Host, token)
	return
}

type stateFn func(msg *pb.TestResult) (stateFn, error)
type state struct {
	name      string
	state     stateFn
	stepLabel string
	ui        ui.UI
}

func (b *VerifyCmd) Run(authToken string, scriptFile string) (err error) {
	ui := ui.NewSimpleUI()

	ui.SetStatus("Flood Chrome Verify")

	ui.Log("flood chrome channel: ", b.FloodChromeChannel)
	ui.Log("script file:", scriptFile)

	f, err := os.Open(scriptFile)
	if err != nil {
		err = errors.Wrapf(err, "unable to open script at %s", scriptFile)
		return
	}

	scriptBytes, err := ioutil.ReadAll(f)
	if err != nil {
		err = errors.Wrap(err, "unable to read script contents")
		return
	}

	client, err := b.floodchromeClient(authToken)
	if err != nil {
		err = errors.Wrap(err, "unable to init flood-chrome client")
		return
	}

	test := &pb.TestRequest{
		Script:             string(scriptBytes),
		ScriptFilename:     scriptFile,
		FloodChromeVersion: b.FloodChromeChannel,
	}

	state := &state{
		name: "testy",
		ui:   ui,
	}
	state.state = state.awaitTest

	stream, err := client.Run(context.Background(), test)
	if err != nil {
		err = errors.Wrap(err, "unable to call test.Run")
		return
	}
	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "error on test result stream")
		}

		// state.dumpLife(result)

		// fmt.Printf("result = %+v\n", result)
		// fmt.Printf("result = %T\n", result)
		// fmt.Println("result", result.String())

		err = state.next(result)
		if err != nil {
			return errors.Wrap(err, "error handling test result message")
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

func (s *state) setStepStatus(arrow string, msg ...interface{}) {
	s.ui.SetSlot(fmt.Sprintf("-%s-> %s", arrow, s.stepLabel), fmt.Sprint(msg...))
}

func dump(msg *pb.TestResult) {
	fmt.Printf("msg = %+v\n", msg)
	fmt.Printf("msg = %T\n", msg)
	fmt.Println("msg", msg.String())
}

func (s *state) dumpLife(msg *pb.TestResult) {
	if lifecycleM := msg.GetLifecycle(); lifecycleM != nil {
		s.ui.Logf("[ life] [%10s] %s", lifecycleM.Event.String(), msg.Label)
	}
}

func (s *state) dumpLog(msg *pb.TestResult) {
	if logM := msg.GetServerLog(); logM != nil {
		// TODO make a switch --server-logs
		// s.ui.Logf("[server-%5s] %+v\n", logM.Level, msg.Message)
	} else if logM := msg.GetScriptLog(); logM != nil {
		s.ui.Logf("[script-%5s] %+v", logM.Level, msg.Message)
	}
}

func matchLifecycle(msg *pb.TestResult, event pb.TestResult_Lifecycle_Event) bool {
	lifecycleM := msg.GetLifecycle()
	return lifecycleM != nil && lifecycleM.Event == event
}

func (s *state) awaitTest(msg *pb.TestResult) (next stateFn, err error) {
	next = s.awaitTest

	if msg.Label == "proxy" && matchLifecycle(msg, pb.TestResult_Lifecycle_Setup) {
		s.ui.SetStatus("Proxy starting")

	} else if msg.Label == "floodchrome" {
		if matchLifecycle(msg, pb.TestResult_Lifecycle_Setup) {
			s.ui.SetStatus("setting up test")

		} else if errM := msg.GetError(); errM != nil {
			return s.handleCompilationError(msg, errM)

		}
	} else if msg.Label == "test" && matchLifecycle(msg, pb.TestResult_Lifecycle_BeforeTest) {
		s.ui.SetStatus("Test starting")
		next = s.awaitPlan

	} else {
		s.dumpLog(msg)
	}

	return
}

func (s *state) handleCompilationError(msg *pb.TestResult, errM *pb.TestResult_Error) (next stateFn, err error) {
	next = s.exhaustNothing

	s.ui.SetStatus("compilation error")
	s.ui.Log(errM.Message)

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
	} else if matchError(msg) {
		return s.handleTestError(msg)
	} else {
		s.dumpLog(msg)
	}

	return
}

func matchError(msg *pb.TestResult) bool {
	errM := msg.GetError()
	return errM != nil
}

func (s *state) handleTestError(msg *pb.TestResult) (next stateFn, err error) {
	next = s.exhaustNothing
	errM := msg.GetError()

	if errM.Internal {
		s.ui.SetStatus("Error: internal floodchrome server error")
		return
	}

	s.ui.SetStatus("Error running test script")

	s.ui.Log(errM.Message)
	s.ui.Log()
	s.ui.Log(errM.Callsite.Code)
	s.ui.Logf("%s^", strings.Repeat(" ", int(errM.Callsite.Column)))
	for _, line := range errM.Stack {
		s.ui.Log(line)
	}

	return
}

func (s *state) awaitNext(msg *pb.TestResult) (next stateFn, err error) {
	next = s.awaitNext

	if matchLifecycle(msg, pb.TestResult_Lifecycle_BeforeStep) {
		s.ui.SetStatus("Running step: ", msg.Label)
		s.stepLabel = msg.Label
		s.setStepStatus("-", "running")
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
		s.setStepStatus("!", "failed")
		// s.ui.Log("=!=> failed")
	case pb.TestResult_Lifecycle_StepSucceeded:
		s.setStepStatus("+", "succeeded")
		// s.ui.Log("=v=> succeeded")
	case pb.TestResult_Lifecycle_StepSkipped:
		s.setStepStatus(" ", "skipped")
		// s.ui.Log("= => skipped")
	case pb.TestResult_Lifecycle_BeforeStepAction:
		s.setStepStatus("-", "running ", msg.Label, "()")
		// s.ui.Log("---> running", msg.Label)
	case pb.TestResult_Lifecycle_AfterStep:
		s.ui.SetStatus("Finished step ", msg.Label)
		next = s.awaitNext
	}

	return
}

func (s *state) handleStepError(msg *pb.TestResult, testError *pb.TestResult_Error) (next stateFn, err error) {
	next = s.handleStep

	s.ui.Log("-!-> error !")
	s.ui.Logf("message = %+v", msg.Message)
	s.ui.Logf("testError = %+v", testError.Message)

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
