package verify

import (
	pb "github.com/flood-io/go-wrenches/floodchrome"
)

type Test interface {
	Error() error

	UpdateStatus(string)

	GetScriptPath() string
	GetChannel() string
	AssertConfigured()

	SetContainerVersion(string)
	GetContainerVersion() string
	SetContainerChannel(string)
	GetContainerChannel() string

	SetSteps([]string)
	SetSettings(map[string]string)
	AssertEnvironmentReady()

	AssertReady()

	ScriptError(message string, maybeErrors ...*pb.TestResult_Error)
	// XXX merge^
	CompilationError(compError *pb.TestResult_Error)

	AssertStep(msg string, step string) bool

	TestBefore(label string)

	StepBefore(label string)

	ActionBefore(label string)
	ActionAfter(label string)

	StepSucceeded(label string)
	StepFailed(label string)
	StepSkipped(label string)
	StepAfter(label string)

	TestSucceeded(label string)
	TestFailed(label string)
	TestAfter(label string)

	Exit(label string)

	WrappedTest() Test
}
