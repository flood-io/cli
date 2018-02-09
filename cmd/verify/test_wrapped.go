package verify

import (
	pb "github.com/flood-io/go-wrenches/floodchrome"
)

var _ Test = (*WrapperTest)(nil)

type WrapperTest struct {
	wrapped Test
}

func (t *WrapperTest) WrappedTest() Test { return t.wrapped }

func (t *WrapperTest) Error() error {
	return t.WrappedTest().Error()
}

func (t *WrapperTest) UpdateStatus(s string) {
	t.WrappedTest().UpdateStatus(s)
	return
}

// configuration
func (t *WrapperTest) AssertConfigured() {
	t.WrappedTest().AssertConfigured()
	return
}

func (t *WrapperTest) SetContainerVersion(v string) {
	t.WrappedTest().SetContainerVersion(v)
	return
}

func (t *WrapperTest) SetContainerChannel(c string) {
	t.WrappedTest().SetContainerVersion(c)
	return
}

func (t *WrapperTest) SetSteps(steps []string) {
	t.WrappedTest().SetSteps(steps)
	return
}

func (t *WrapperTest) SetSettings(settings map[string]string) {
	t.WrappedTest().SetSettings(settings)
	return
}

func (t *WrapperTest) AssertEnvironmentReady() {
	t.WrappedTest().AssertEnvironmentReady()
	return
}

func (t *WrapperTest) AssertReady() {
	t.WrappedTest().AssertReady()
	return
}

func (t *WrapperTest) ScriptError(message string, maybeErrors ...*pb.TestResult_Error) {
	t.WrappedTest().ScriptError(message, maybeErrors...)
	return
}

func (t *WrapperTest) CompilationError(compErr *pb.TestResult_Error) {
	t.WrappedTest().CompilationError(compErr)
	return
}

func (t *WrapperTest) AssertStep(msg string, step string) bool {
	return t.WrappedTest().AssertStep(msg, step)
}

func (t *WrapperTest) TestBefore(label string) {
	t.WrappedTest().TestBefore(label)
	return
}

func (t *WrapperTest) StepBefore(step string) {
	t.WrappedTest().StepBefore(step)
	return
}

func (t *WrapperTest) ActionBefore(label string) {
	t.WrappedTest().ActionBefore(label)
	return
}

func (t *WrapperTest) ActionAfter(label string) {
	t.WrappedTest().ActionAfter(label)
	return
}

func (t *WrapperTest) StepSucceeded(label string) {
	t.WrappedTest().StepSucceeded(label)
	return
}

func (t *WrapperTest) StepFailed(label string) {
	t.WrappedTest().StepFailed(label)
	return
}

func (t *WrapperTest) StepSkipped(label string) {
	t.WrappedTest().StepSkipped(label)
	return
}

func (t *WrapperTest) StepAfter(step string) {
	t.WrappedTest().StepAfter(step)
	return
}

func (t *WrapperTest) TestSucceeded(label string) {
	t.WrappedTest().TestSucceeded(label)
	return
}

func (t *WrapperTest) TestFailed(label string) {
	t.WrappedTest().TestFailed(label)
	return
}

func (t *WrapperTest) TestAfter(label string) {
	t.WrappedTest().TestAfter(label)
	return
}

func (t *WrapperTest) Exit(label string) {
	t.WrappedTest().Exit(label)
	return
}
