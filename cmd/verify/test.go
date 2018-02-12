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

	SetSettings(map[string]string)
	GetSettings() map[string]string
	GetSetting(key string) string

	SetSteps([]string)
	GetSteps() []*Step

	AssertEnvironmentReady()

	CompilationError(compError *pb.TestResult_Error)
	ScriptError(message string, scriptErr *pb.TestResult_Error)
	InternalScriptError(message string, scriptErr *pb.TestResult_Error)

	ScriptLog(level, message string)

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
