package verify

import "fmt"

type LoggingTest struct {
	*WrapperTest
}

var _ Test = (*LoggingTest)(nil)

func NewLoggingTest(innerTest Test) *LoggingTest {
	return &LoggingTest{
		WrapperTest: &WrapperTest{innerTest},
	}
}

// func (t *LoggingTest) WrappedTest() Test {
// return t.wrappedTest
// }

func (t *LoggingTest) AssertConfigured() {
	fmt.Println("AssertConfigured")
	t.WrappedTest().AssertConfigured()
}
