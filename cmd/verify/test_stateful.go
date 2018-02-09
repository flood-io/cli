package verify

import (
	"fmt"

	pb "github.com/flood-io/go-wrenches/floodchrome"
)

type StatefulTest struct {
	ScriptFile string
	Channel    string

	ContainerVersion string
	ContainerChannel string

	Title       string
	Description string

	StepKeys []string
	Steps    map[string]*Step

	CurrentStep *Step

	Settings map[string]string

	Status string

	CurrentError error
}

var _ Test = (*StatefulTest)(nil)

type Step struct {
	Title string
}

func (t *StatefulTest) WrappedTest() Test { return nil }

func (t *StatefulTest) Error() error {
	return t.CurrentError
}

func (t *StatefulTest) UpdateStatus(s string) {
	t.Status = s
}

// configuration
func (t *StatefulTest) AssertConfigured() {
	fmt.Println("StatefulTest: AssertConfigured")
}

func (t *StatefulTest) SetContainerVersion(v string) {
	t.ContainerVersion = v
}

func (t *StatefulTest) SetContainerChannel(c string) {
	t.ContainerChannel = c
}

func (t *StatefulTest) SetSteps(steps []string) {
	t.StepKeys = steps
	t.Steps = make(map[string]*Step, len(steps))
	for _, step := range steps {
		t.Steps[step] = &Step{Title: step}
	}
}

func (t *StatefulTest) SetSettings(settings map[string]string) {
	t.Settings = settings
}

func (t *StatefulTest) AssertEnvironmentReady() {
}

func (t *StatefulTest) AssertReady() {
}

func (t *StatefulTest) ScriptError(message string, maybeErrors ...*pb.TestResult_Error) {
	var scriptError *pb.TestResult_Error
	if len(maybeErrors) > 0 {
		scriptError = maybeErrors[0]
	}

	if scriptError == nil {
		// synthesise one...
	}

	fmt.Printf("scriptError = %+v\n", scriptError)
}

func (t *StatefulTest) CompilationError(compErr *pb.TestResult_Error) {
}

// s.ui.Log(errM.Message)
// s.ui.Log()
// s.ui.Log(errM.Callsite.Code)
// s.ui.Logf("%s^", strings.Repeat(" ", int(errM.Callsite.Column)))
// for _, line := range errM.Stack {
// s.ui.Log(line)
// }

// s.ui.Log("-!-> error !")
// s.ui.Logf("message = %+v", msg.Message)

// s.ui.Logf("%s", testError.Message)

// c := testError.Callsite
// s.ui.Logf("%s:%d\n%s\n%s^", c.File, c.Line, strings.Replace(c.Code, "\t", "    ", -1), strings.Repeat(" ", int(c.Column)))
// for _, line := range testError.Stack {
// s.ui.Logf("  %s", line)
// }

func (t *StatefulTest) AssertStep(msg string, step string) bool {
	_, ok := t.Steps[step]
	if !ok {
		t.CurrentError = fmt.Errorf("%s: unknown step '%s'", msg, step)
		return false
	}

	return true
}

func (t *StatefulTest) TestBefore(label string) {
}

func (t *StatefulTest) StepBefore(step string) {
	if !t.AssertStep("StepBefore", step) {
		return
	}
	t.CurrentStep = t.Steps[step]
}

func (t *StatefulTest) ActionBefore(label string) {
}

func (t *StatefulTest) ActionAfter(label string) {
}

func (t *StatefulTest) StepSucceeded(label string) {
}

func (t *StatefulTest) StepFailed(label string) {
}

func (t *StatefulTest) StepSkipped(label string) {
}

func (t *StatefulTest) StepAfter(step string) {
	if !t.AssertStep("StepAfter", step) {
		return
	}
	t.CurrentStep = nil
}

func (t *StatefulTest) TestSucceeded(label string) {
}

func (t *StatefulTest) TestFailed(label string) {
}

func (t *StatefulTest) TestAfter(label string) {
}

func (t *StatefulTest) Exit(label string) {
}
