package verify

import (
	"fmt"

	pb "github.com/flood-io/go-wrenches/floodchrome"
)

type Test struct {
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

type Step struct {
	Title string
}

func (t *Test) Error() error {
	return t.CurrentError
}

func (t *Test) UpdateStatus(s string) {
	t.Status = s
}

// setters
func (t *Test) SetContainerVersion(v string) {
	t.ContainerVersion = v
}

func (t *Test) SetContainerChannel(c string) {
	t.ContainerChannel = c
}

func (t *Test) SetSteps(steps []string) {
	t.StepKeys = steps
	t.Steps = make(map[string]*Step, len(steps))
	for i, step := range steps {
		t.Steps[step] = &Step{Title: step}
	}
}

func (t *Test) SetSettings(settings map[string]string) {
	t.Settings = settings
}

func (t *Test) EnvironmentReady() {
}

func (t *Test) Ready() {
}

func (t *Test) ScriptError(message string, maybeErrors ...*pb.TestResult_Error) {
	if len(maybeErrors) > 0 {
		err = maybeErrors[0]
	}
}

func (t *Test) CompilationError(compErr *pb.TestResult_Error) {
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

func (t *Test) AssertStep(msg string, step string) bool {
	_, ok := t.Steps[step]
	if !ok {
		t.CurrentError = fmt.Errorf("%s: unknown step '%s'", msg, step)
		return false
	}

	return true
}

func (t *Test) TestBefore(label string) {
}

func (t *Test) StepBefore(step string) {
	if !t.AssertStep("StepBefore", step) {
		return
	}
	t.CurrentStep = t.Steps[step]
}

func (t *Test) ActionBefore(label string) {
}

func (t *Test) ActionAfter(label string) {
}

func (t *Test) StepSucceeded(label string) {
}

func (t *Test) StepFailed(label string) {
}

func (t *Test) StepSkipped(label string) {
}

func (t *Test) StepAfter(label string) {
	if !t.AssertStep("StepAfter", step) {
		return
	}
	t.CurrentStep = nil
}

func (t *Test) TestSucceeded(label string) {
}

func (t *Test) TestFailed(label string) {
}

func (t *Test) TestAfter(label string) {
}

func (t *Test) Exit(label string) {
}
