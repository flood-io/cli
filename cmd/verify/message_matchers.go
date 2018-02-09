package verify

import pb "github.com/flood-io/go-wrenches/floodchrome"

func matchError(msg *pb.TestResult) bool {
	errM := msg.GetError()
	return errM != nil
}

func matchLifecycle(msg *pb.TestResult, event pb.TestResult_Lifecycle_Event) bool {
	lifecycleM := msg.GetLifecycle()
	return lifecycleM != nil && lifecycleM.Event == event
}
