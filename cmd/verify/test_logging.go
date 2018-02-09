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

func (t *LoggingTest) UpdateStatus(status string) {
	t.WrappedTest().UpdateStatus(status)

	fmt.Println(status)
}

func (t *LoggingTest) AssertConfigured() {
	t.WrappedTest().AssertConfigured()

	fmt.Println("~~~ Flood Chrome Verify ~~~")
	fmt.Printf("  - test script       : %s\n", t.GetScriptPath())
	fmt.Printf("  - requested channel : %s\n", t.GetChannel())
}

func (t *LoggingTest) AssertEnvironmentReady() {
	t.WrappedTest().AssertEnvironmentReady()

	fmt.Println("---> environment ready")
	fmt.Printf("   - container version : %s\n", t.GetContainerVersion())
	fmt.Printf("   - container channel : %s\n", t.GetContainerVersion())
}

func (t *LoggingTest) AssertReady() {
	t.WrappedTest().AssertReady()

	fmt.Println("... test ready")
}
