package bludev

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
	name  string
	state stateFn
}

func (b *BLUDev) Run(scriptFile string) (err error) {
	fmt.Println("running dev-blu")
	fmt.Println("flood chrome channel: ", b.FloodChromeChannel)
	fmt.Println("script file:", scriptFile)

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

	state := &state{name: "testy"}
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

		// fmt.Printf("result = %+v\n", result)
		// fmt.Printf("result = %T\n", result)
		// fmt.Println("result", result.String())

		err = state.next(result)
		if err != nil {
			log.Fatalf("error handling message %v", err)
		}

		// dispatch
		// if lifecycleM := result.GetLifecycle(); lifecycleM != nil {
		// fmt.Printf("[ life] [%10s] %s\n", lifecycleM.Event.String(), lifecycleM.Label)

		// } else if logM := result.GetLog(); logM != nil {
		// fmt.Printf("[%5s] %+v\n", logM.Level, result.Message)

		// } else if measurementM := result.GetMeasurement(); measurementM != nil {
		// if currentStep != measurementM.Label {
		// currentStep = measurementM.Label
		// fmt.Println("")
		// fmt.Println("==================")
		// fmt.Println("[ step]", currentStep)
		// }
		// fmt.Printf("[ meas] %s - %s - %v\n", measurementM.Label, measurementM.Measurement, measurementM.Value)

		// } else if traceM := result.GetTrace(); traceM != nil {
		// fmt.Printf("[trace] %s - response code %s\n", result.Message, traceM.ResponseCode)
		// if networkT := traceM.GetNetwork(); networkT != nil {
		// err = writeNetworkTrace(networkT)
		// if err != nil {
		// return err
		// }
		// }

		// } else if errorM := result.GetError(); errorM != nil {
		// fmt.Printf("[error] %s %s\n", errorM.Label, result.Message)
		// fmt.Println(errorM.Detail)

		// } else if completeM := result.GetComplete(); completeM != nil {
		// break

		// } else {
		// fmt.Println("--- unhandled type ---")
		// fmt.Printf("result = %+T\n", result)
		// fmt.Println(result.Message)
		// fmt.Println(result.String())
		// }
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
	s.state, err = s.state(msg)
	return
}

func (s *state) nothing(msg *pb.TestResult) (next stateFn, err error) {
	next = s.nothing
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

func dump(msg *pb.TestResult) {
	fmt.Printf("msg = %+v\n", msg)
	fmt.Printf("msg = %T\n", msg)
	fmt.Println("msg", msg.String())
}

func (s *state) dumpLog(msg *pb.TestResult) {
	if logM := msg.GetLog(); logM != nil {
		fmt.Printf("[%5s] %+v\n", logM.Level, msg.Message)
	} else if lifecycleM := msg.GetLifecycle(); lifecycleM != nil {
		fmt.Printf("[ life] [%10s] %s\n", lifecycleM.Event.String(), lifecycleM.Label)
	}
}

func matchLifecycle(msg *pb.TestResult, label string, event pb.TestResult_Lifecycle_Event) bool {
	lifecycleM := msg.GetLifecycle()
	return lifecycleM != nil && lifecycleM.Label == label && lifecycleM.Event == event
}

func matchLifecycleEvent(msg *pb.TestResult, event pb.TestResult_Lifecycle_Event) bool {
	lifecycleM := msg.GetLifecycle()
	return lifecycleM != nil && lifecycleM.Event == event
}

func lifecycleEventLabel(msg *pb.TestResult, event pb.TestResult_Lifecycle_Event) string {
	lifecycleM := msg.GetLifecycle()
	if lifecycleM != nil && lifecycleM.Event == event {
		return lifecycleM.Label
	} else {
		return ""
	}
}

func (s *state) initialState(msg *pb.TestResult) (next stateFn, err error) {
	fmt.Println("~~~~ Flood Chrome ~~~~")

	// dump(msg)

	return s.awaitProxyStart(msg)
}

func (s *state) awaitProxyStart(msg *pb.TestResult) (next stateFn, err error) {
	next = s.awaitProxyStart

	if matchLifecycle(msg, "proxy", pb.TestResult_Lifecycle_Setup) {
		fmt.Println("==== Proxy starting ===")
		next = s.awaitFCStart
	} else {
		s.dumpLog(msg)
	}

	return
}

func (s *state) awaitFCStart(msg *pb.TestResult) (next stateFn, err error) {
	next = s.awaitFCStart

	if matchLifecycle(msg, "floodchrome", pb.TestResult_Lifecycle_Setup) {
		fmt.Println("==== FC starting ===")
		next = s.awaitTest
	} else {
		s.dumpLog(msg)
	}
	return
}

func (s *state) awaitTest(msg *pb.TestResult) (next stateFn, err error) {
	next = s.awaitTest

	if matchLifecycle(msg, "test", pb.TestResult_Lifecycle_BeforeTest) {
		fmt.Println("==== Test starting ===")
		next = s.awaitStep
	} else {
		s.dumpLog(msg)
	}
	return
}

func (s *state) awaitStep(msg *pb.TestResult) (next stateFn, err error) {
	next = s.awaitStep

	if label := lifecycleEventLabel(msg, pb.TestResult_Lifecycle_BeforeStep); label != "" {
		fmt.Println("---> Step: ", label)
		next = s.handleStep
	} else if label := lifecycleEventLabel(msg, pb.TestResult_Lifecycle_StepSkipped); label != "" {
		fmt.Println("   > skipped step: ", label)
		next = s.awaitStep
	} else {
		s.dumpLog(msg)
	}
	return
}

func (s *state) handleStep(msg *pb.TestResult) (next stateFn, err error) {
	next = s.handleStep

	lifecycle := msg.GetLifecycle()
	if lifecycle != nil {
		next, err = s.handleStepLifecycle(lifecycle)
	} else {
		s.dumpLog(msg)
	}
	return
}

func (s *state) handleStepLifecycle(lifecycle *pb.TestResult_Lifecycle) (next stateFn, err error) {
	next = s.handleStep

	switch lifecycle.Event {
	case pb.TestResult_Lifecycle_StepFailed:
		fmt.Println("=!=> failed")
		// next = s.failedStepAwaitTestFinished
	case pb.TestResult_Lifecycle_StepSucceeded:
		fmt.Println("=v=> succeeded")
	case pb.TestResult_Lifecycle_StepSkipped:
		fmt.Println("= => skipped")
	case pb.TestResult_Lifecycle_AfterStep:
		next = s.awaitStep
	}

	return
}

func (s *state) failedStepAwaitTestFinished(msg *pb.TestResult) (next stateFn, err error) {
	next = s.failedStepAwaitTestFinished

	if matchLifecycle(msg, "test", pb.TestResult_Lifecycle_AfterTest) {
		fmt.Println("<=== test finished ===>")
	} else {
		// s.dumpLog(msg)
	}

	return
}
