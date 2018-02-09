package verify

import (
	pb "github.com/flood-io/go-wrenches/floodchrome"
)

func (s *state) awaitTest(msg *pb.TestResult) (next stateFn) {
	next = s.awaitTest

	if msg.Label == "proxy" && matchLifecycle(msg, pb.TestResult_Lifecycle_Setup) {
		s.Test.UpdateStatus("Proxy starting")

	} else if msg.Label == "floodchrome" {

		if matchLifecycle(msg, pb.TestResult_Lifecycle_Setup) {
			s.Test.UpdateStatus("setting up test")

		} else if verM := msg.GetFloodChromeVersion(); verM != nil {
			s.Test.SetContainerVersion(verM.Version)
			s.Test.SetContainerChannel(verM.Channel)

		} else if errM := msg.GetError(); errM != nil {
			return s.handleCompilationError(msg, errM)

		}
	} else if msg.Label == "test" && matchLifecycle(msg, pb.TestResult_Lifecycle_BeforeTest) {
		s.Test.AssertEnvironmentReady()
		next = s.awaitPlan

	} else {
		s.dumpLog(msg)
	}

	return
}

func (s *state) handleCompilationError(msg *pb.TestResult, errM *pb.TestResult_Error) (next stateFn) {
	next = s.exhaustNothing

	s.Test.CompilationError(errM)

	return
}

func (s *state) awaitPlan(msg *pb.TestResult) (next stateFn) {
	next = s.awaitPlan

	plan := msg.GetPlan()
	if plan != nil {
		s.Test.SetSettings(plan.Settings)
		s.Test.SetSteps(plan.Steps)
		s.Test.AssertReady()

		next = s.awaitNext
	} else if matchError(msg) {
		return s.handleTestError(msg)
	} else {
		s.dumpLog(msg)
	}

	return
}

func (s *state) handleTestError(msg *pb.TestResult) (next stateFn) {
	next = s.exhaustNothing
	errM := msg.GetError()

	if errM.Internal {
		s.Test.ScriptError("Error: internal floodchrome server error")
		return
	}

	s.Test.ScriptError("Error running test script", errM)

	return
}

func (s *state) awaitNext(msg *pb.TestResult) (next stateFn) {
	next = s.awaitNext
	lifeM := msg.GetLifecycle()

	if lifeM == nil {
		s.dumpLog(msg)
		return
	}

	switch lifeM.Event {
	case pb.TestResult_Lifecycle_BeforeStep:
		s.Test.StepBefore(msg.Label)
		next = s.handleStep

	case pb.TestResult_Lifecycle_AfterStep:
		s.Test.StepAfter(msg.Label)

	case pb.TestResult_Lifecycle_TestSucceeded:
		s.Test.TestSucceeded(msg.Label)

	case pb.TestResult_Lifecycle_TestFailed:
		s.Test.TestFailed(msg.Label)

	case pb.TestResult_Lifecycle_StepSkipped:
		s.Test.StepSkipped(msg.Label)
		next = s.awaitNext

	case pb.TestResult_Lifecycle_AfterTest:
		s.Test.TestAfter(msg.Label)
		next = s.awaitExit

	default:
		s.dumpLog(msg)
	}
	return
}

func (s *state) handleStep(msg *pb.TestResult) (next stateFn) {
	next = s.handleStep

	lifecycle := msg.GetLifecycle()
	testError := msg.GetError()
	if lifecycle != nil {
		next = s.handleStepLifecycle(msg, lifecycle)
	} else if testError != nil {
		next = s.handleStepError(msg, testError)
	} else {
		s.dumpLog(msg)
	}
	return
}

func (s *state) handleStepLifecycle(msg *pb.TestResult, lifecycle *pb.TestResult_Lifecycle) (next stateFn) {
	next = s.handleStep

	switch lifecycle.Event {
	case pb.TestResult_Lifecycle_StepFailed:
		s.Test.StepFailed(msg.Label)
	case pb.TestResult_Lifecycle_StepSucceeded:
		s.Test.StepSucceeded(msg.Label)
	case pb.TestResult_Lifecycle_StepSkipped:
		s.Test.StepSkipped(msg.Label)
	case pb.TestResult_Lifecycle_BeforeStepAction:
		s.Test.ActionBefore(msg.Label)
	case pb.TestResult_Lifecycle_AfterStep:
		s.Test.StepAfter(msg.Label)
		next = s.awaitNext
	}

	return
}

func (s *state) handleStepError(msg *pb.TestResult, testError *pb.TestResult_Error) (next stateFn) {
	next = s.handleStep

	s.Test.ScriptError("error during step", testError)

	return
}

func (s *state) awaitExit(msg *pb.TestResult) (next stateFn) {
	next = s.awaitExit

	if msg.Label == "floodchrome" && matchLifecycle(msg, pb.TestResult_Lifecycle_Exit) {
		next = s.exhaustNothing
	} else {
		s.dumpLog(msg)
	}

	return
}
