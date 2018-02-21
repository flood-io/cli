package verify

import (
	"bytes"
	"fmt"
	"strings"

	pb "github.com/flood-io/go-wrenches/floodchrome"
	au "github.com/logrusorgru/aurora"
)

type indentLevel int

const (
	outside indentLevel = iota
	insideTest
	insideStep
)

type LoggingTest struct {
	*WrapperTest
	currentIndent indentLevel
}

var _ Test = (*LoggingTest)(nil)

func NewLoggingTest(innerTest Test) *LoggingTest {
	return &LoggingTest{
		WrapperTest:   &WrapperTest{innerTest},
		currentIndent: outside,
	}
}

func (t *LoggingTest) UpdateStatus(status string) {
	t.WrappedTest().UpdateStatus(status)

	fmt.Println(status)
}

func (t *LoggingTest) CompilationError(compErr *pb.TestResult_Error) {
	t.WrappedTest().CompilationError(compErr)
	t.Println(au.Red("x Unable to compile test script "), t.GetScriptPath())
	t.Println()

	t.Println(compErr.Message)
}

func (t *LoggingTest) ScriptError(message string, scriptErr *pb.TestResult_Error) {
	t.WrappedTest().ScriptError(message, scriptErr)
	t.Println(au.Red("Script Error"))

	t.Println()

	// print unindented
	fmt.Println(scriptErr.Message)

	// NOTE this is assuming 4-space-tabs
	// Column is in spaces even when the code is in tabs
	if c := scriptErr.Callsite; c != nil {
		fmt.Println(strings.Replace(c.Code, "\t", "    ", -1))
		fmt.Printf("%s^\n", strings.Repeat(" ", int(c.Column)))

		for _, line := range scriptErr.Stack {
			fmt.Println(line)
		}
	}

	t.Println()
}

func (t *LoggingTest) InternalScriptError(message string, scriptErr *pb.TestResult_Error) {
	t.WrappedTest().InternalScriptError(message, scriptErr)
	t.Println(au.Red("Internal Flood Chrome Error"))
	fmt.Println("An internal Flood Chrome error has occurred.")
	fmt.Println("This could be a bug in Flood Chrome or a temporary infrastructural issue.")
	fmt.Println("Please retry your script and if the error persists, please report to flood via ...")
}

func scriptLogPrefix(level string) string {
	switch level {
	case "log":
		return au.Blue("console.log").String()
	case "info":
		return au.Blue("console.info").String()
	case "dir":
		return au.Blue("console.dir").String()
	case "trace":
		return au.Blue("console.trace").String()
	case "error":
		return au.Red("console.error").String()
	default:
		return fmt.Sprintf("console.%s", level)
	}

}

func (t *LoggingTest) ScriptLog(level, message string) {
	t.WrappedTest().ScriptLog(level, message)
	t.Printf("%s: %s\n", scriptLogPrefix(level), message)
}

func (t *LoggingTest) AssertConfigured() {
	t.WrappedTest().AssertConfigured()

	t.Printfln("%-17s : %s", au.Gray("test script"), t.GetScriptPath())
	t.Printfln("%-17s : %s", au.Gray("requested channel"), t.GetChannel())
	t.Println()
}

func (t *LoggingTest) AssertEnvironmentReady() {
	t.WrappedTest().AssertEnvironmentReady()

	t.Println(au.Blue("--->"), " Verification environment ready")
	t.Printfln("%-20s : %s", au.Gray("container version"), t.GetContainerVersion())
	t.Printfln("%-20s : %s", au.Gray("container channel"), t.GetContainerVersion())
	t.Println()
}

func (t *LoggingTest) TestBefore(label string) {
	t.WrappedTest().TestBefore(label)

	name := t.WrappedTest().GetSetting("name")
	if name == "" {
		name = "untitled"
	}
	t.Println(au.Blue("--->"), " Starting test ", au.Brown(name))
	t.Println()
	t.Println("Test Plan:")
	for i, step := range t.GetSteps() {
		t.Printfln("%2d. %s", au.Gray(i+1), step.Title)
	}
	t.Println()

	t.Println("Test Settings:")
	for k, v := range t.GetSettings() {
		t.Printfln("  %-20s : %s", au.Gray(k), v)
	}
	t.Println()

	t.currentIndent = insideTest
}
func (t *LoggingTest) TestSucceeded(label string) {
	t.WrappedTest().TestSucceeded(label)
	t.Println(au.Green("+ test succeeded"))
}
func (t *LoggingTest) TestFailed(label string) {
	t.WrappedTest().TestFailed(label)
	t.Println(au.Red("x test failed"))
}

func (t *LoggingTest) TestAfter(label string) {
	t.WrappedTest().TestAfter(label)
	t.currentIndent = outside
}

func (t *LoggingTest) StepBefore(label string) {
	t.WrappedTest().StepBefore(label)
	t.Println(au.Gray(label).Bold())
	t.currentIndent = insideStep
}
func (t *LoggingTest) StepSucceeded(label string) {
	t.WrappedTest().StepSucceeded(label)
	t.Println(au.Green("+ succeeded"))
}
func (t *LoggingTest) StepFailed(label string) {
	t.WrappedTest().StepFailed(label)
	t.Println(au.Red("x step failed"))
}
func (t *LoggingTest) StepSkipped(label string) {
	t.WrappedTest().StepSkipped(label)
	t.Println(au.Brown("- skipped"))
}
func (t *LoggingTest) StepAfter(label string) {
	t.WrappedTest().StepAfter(label)
	t.Println()
	t.currentIndent = insideTest
}

func (t *LoggingTest) ActionBefore(label string) {
	t.WrappedTest().ActionBefore(label)
	t.Printfln("%s()", au.Gray(label).Bold())
}
func (t *LoggingTest) ActionAfter(label string) {
	t.WrappedTest().ActionAfter(label)
}

func (t *LoggingTest) indent() string {
	switch t.currentIndent {
	case outside:
		return ""
	case insideTest:
		return "    "
	case insideStep:
		return "        "
	default:
		return ""
	}
}

func (t *LoggingTest) Printf(format string, args ...interface{}) {
	var buf bytes.Buffer
	buf.WriteString(t.indent())
	buf.WriteString(format)
	fmt.Printf(buf.String(), args...)
}

func (t *LoggingTest) Printfln(format string, args ...interface{}) {
	var buf bytes.Buffer
	buf.WriteString(t.indent())
	buf.WriteString(format)
	buf.WriteString("\n")
	fmt.Printf(buf.String(), args...)
}

func (t *LoggingTest) Println(args ...interface{}) {
	var buf bytes.Buffer
	buf.WriteString(t.indent())
	buf.WriteString(fmt.Sprint(args...))
	buf.WriteString("\n")
	fmt.Print(buf.String())
}
