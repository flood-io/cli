package verify

import (
	"fmt"

	pb "github.com/flood-io/go-wrenches/floodchrome"
)

type StatefulTest struct {
	ScriptPath string
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

func NewStatefulTest(scriptPath, channel string) *StatefulTest {
	return &StatefulTest{ScriptPath: scriptPath, Channel: channel}
}

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

func (t *StatefulTest) GetScriptPath() string {
	return t.ScriptPath
}

func (t *StatefulTest) GetChannel() string {
	return t.Channel
}

func (t *StatefulTest) SetContainerVersion(v string) {
	t.ContainerVersion = v
}

func (t *StatefulTest) GetContainerVersion() string {
	return t.ContainerVersion
}

func (t *StatefulTest) SetContainerChannel(c string) {
	t.ContainerChannel = c
}

func (t *StatefulTest) GetContainerChannel() string {
	return t.ContainerChannel
}

func (t *StatefulTest) GetSettings() map[string]string {
	return t.Settings
}

func (t *StatefulTest) GetSetting(key string) string {
	return t.Settings[key]
}

func (t *StatefulTest) SetSteps(steps []string) {
	t.StepKeys = steps
	t.Steps = make(map[string]*Step, len(steps))
	for _, step := range steps {
		t.Steps[step] = &Step{Title: step}
	}
}
func (t *StatefulTest) GetSteps() []*Step {
	steps := make([]*Step, len(t.StepKeys))
	for i, key := range t.StepKeys {
		steps[i] = t.Steps[key]
	}

	return steps
}

func (t *StatefulTest) SetSettings(settings map[string]string) {
	t.Settings = settings
}

func (t *StatefulTest) AssertEnvironmentReady() {
}

func (t *StatefulTest) ScriptError(message string, scriptErr *pb.TestResult_Error) {
}

func (t *StatefulTest) InternalScriptError(message string, scriptErr *pb.TestResult_Error) {
}

func (t *StatefulTest) CompilationError(compErr *pb.TestResult_Error) {
}

func (t *StatefulTest) ScriptLog(level, message string) {
}

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
