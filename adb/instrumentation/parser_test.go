package instrumentation

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const (
	parserLineSeparator = "\n"

	noTestsAreFoundedInstrumentationOutput = `INSTRUMENTATION_RESULT: stream=

Time: 0

OK (0 tests)


INSTRUMENTATION_CODE: -1`

	firstTestFailedBeforeRun = `INSTRUMENTATION_RESULT: shortMsg=Process crashed.
INSTRUMENTATION_CODE: 0`

	testCrashedBeforeRunAfterFirstSuccess = `INSTRUMENTATION_STATUS: numtests=2
INSTRUMENTATION_STATUS: stream=
INSTRUMENTATION_STATUS: id=AndroidJUnitRunner
INSTRUMENTATION_STATUS: test=test1
INSTRUMENTATION_STATUS: class=com.example.test.TestClass
INSTRUMENTATION_STATUS: current=1
INSTRUMENTATION_STATUS_CODE: 1
INSTRUMENTATION_STATUS: numtests=2
INSTRUMENTATION_STATUS: stream=.
INSTRUMENTATION_STATUS: id=AndroidJUnitRunner
INSTRUMENTATION_STATUS: test=test1
INSTRUMENTATION_STATUS: class=com.example.test.TestClass
INSTRUMENTATION_STATUS: current=1
INSTRUMENTATION_STATUS_CODE: 0
INSTRUMENTATION_RESULT: shortMsg=Process crashed.
INSTRUMENTATION_CODE: 0`

	testCrashedAfterFirstTestRun = `INSTRUMENTATION_STATUS: numtests=1
INSTRUMENTATION_STATUS: stream=
ru.test.testapplication.ExampleInstrumentedTest:
INSTRUMENTATION_STATUS: id=AndroidJUnitRunner
INSTRUMENTATION_STATUS: test=test1
INSTRUMENTATION_STATUS: class=com.example.test.TestClass
INSTRUMENTATION_STATUS: current=1
INSTRUMENTATION_STATUS_CODE: 1
INSTRUMENTATION_RESULT: shortMsg=Process crashed.
INSTRUMENTATION_CODE: 0`

	secondTestIgnoredAfterOnePassedRun = `INSTRUMENTATION_STATUS: numtests=2
INSTRUMENTATION_STATUS: stream=
com.example.test.TestClass:
INSTRUMENTATION_STATUS: id=AndroidJUnitRunner
INSTRUMENTATION_STATUS: test=test1
INSTRUMENTATION_STATUS: class=com.example.test.TestClass
INSTRUMENTATION_STATUS: current=1
INSTRUMENTATION_STATUS_CODE: 1
INSTRUMENTATION_STATUS: numtests=2
INSTRUMENTATION_STATUS: stream=.
INSTRUMENTATION_STATUS: id=AndroidJUnitRunner
INSTRUMENTATION_STATUS: test=test1
INSTRUMENTATION_STATUS: class=com.example.test.TestClass
INSTRUMENTATION_STATUS: current=1
INSTRUMENTATION_STATUS_CODE: 0
INSTRUMENTATION_STATUS: numtests=2
INSTRUMENTATION_STATUS: stream=
INSTRUMENTATION_STATUS: id=AndroidJUnitRunner
INSTRUMENTATION_STATUS: test=test2
INSTRUMENTATION_STATUS: class=com.example.test.TestClass
INSTRUMENTATION_STATUS: current=2
INSTRUMENTATION_STATUS_CODE: 1
INSTRUMENTATION_STATUS: numtests=2
INSTRUMENTATION_STATUS: stream=
INSTRUMENTATION_STATUS: id=AndroidJUnitRunner
INSTRUMENTATION_STATUS: test=test2
INSTRUMENTATION_STATUS: class=com.example.test.TestClass
INSTRUMENTATION_STATUS: current=2
INSTRUMENTATION_STATUS_CODE: -3
INSTRUMENTATION_RESULT: stream=

Time: 10.073

OK (2 tests)


INSTRUMENTATION_CODE: -1`

	twoTestsPassedRun = `INSTRUMENTATION_STATUS: numtests=2
INSTRUMENTATION_STATUS: stream=
com.example.test.TestClass:
INSTRUMENTATION_STATUS: id=AndroidJUnitRunner
INSTRUMENTATION_STATUS: test=test1
INSTRUMENTATION_STATUS: class=com.example.test.TestClass
INSTRUMENTATION_STATUS: current=1
INSTRUMENTATION_STATUS_CODE: 1
INSTRUMENTATION_STATUS: numtests=2
INSTRUMENTATION_STATUS: stream=.
INSTRUMENTATION_STATUS: id=AndroidJUnitRunner
INSTRUMENTATION_STATUS: test=test1
INSTRUMENTATION_STATUS: class=com.example.test.TestClass
INSTRUMENTATION_STATUS: current=1
INSTRUMENTATION_STATUS_CODE: 0
INSTRUMENTATION_STATUS: numtests=2
INSTRUMENTATION_STATUS: stream=
INSTRUMENTATION_STATUS: id=AndroidJUnitRunner
INSTRUMENTATION_STATUS: test=test2
INSTRUMENTATION_STATUS: class=com.example.test.TestClass
INSTRUMENTATION_STATUS: current=2
INSTRUMENTATION_STATUS_CODE: 1
INSTRUMENTATION_STATUS: numtests=2
INSTRUMENTATION_STATUS: stream=
INSTRUMENTATION_STATUS: id=AndroidJUnitRunner
INSTRUMENTATION_STATUS: test=test2
INSTRUMENTATION_STATUS: class=com.example.test.TestClass
INSTRUMENTATION_STATUS: current=2
INSTRUMENTATION_STATUS_CODE: 0
INSTRUMENTATION_RESULT: stream=

Time: 9.073

OK (2 tests)


INSTRUMENTATION_CODE: -1`

	passInvalidPackageRun = `INSTRUMENTATION_STATUS: id=ActivityManagerService
invalid instrumentation status bundle
INSTRUMENTATION_STATUS: Error=Unable to find instrumentation info for: ComponentInfo{/tmp/go-build506757578/command-line-arguments/_obj/exe/main/com.avito.android.runner.AvitoInstrumentTestRunner}
INSTRUMENTATION_STATUS_CODE: -1
android.util.AndroidException: INSTRUMENTATION_FAILED: /tmp/go-build506757578/command-line-arguments/_obj/exe/main/com.avito.android.runner.AvitoInstrumentTestRunner
        at com.android.commands.am.Am.runInstrument(Am.java:890)
        at com.android.commands.am.Am.onRun(Am.java:400)
        at com.android.internal.os.BaseCommand.run(BaseCommand.java:51)
        at com.android.commands.am.Am.main(Am.java:121)
        at com.android.internal.os.RuntimeInit.nativeFinishInit(Native Method)
        at com.android.internal.os.RuntimeInit.main(RuntimeInit.java:262)`
)

func startInstrumentationOutputParsing(parser *Parser, output string) []Event {
	outputStream := make(chan string, 1000)

	result := []Event{}
	eventsStream, _ := parser.Process(outputStream)

	for _, line := range strings.Split(output, parserLineSeparator) {
		outputStream <- line
	}
	close(outputStream)

	for event := range eventsStream {
		result = append(result, event)
	}

	return result
}

func TestNoTestsPublishedWhenTestsNotFound(t *testing.T) {
	events := startInstrumentationOutputParsing(
		NewParser(parserLineSeparator),
		noTestsAreFoundedInstrumentationOutput,
	)

	assert.Equal(
		t,
		[]Event{
			TestsRunStartedEvent{
				NumberOfTests: 0,
			},
			TestsRunFinishedEvent{
				Time: 0,
			},
		},
		events,
	)
}

func TestFirstTestCrashDetectedWhenItCrashedBeforeRun(t *testing.T) {
	events := startInstrumentationOutputParsing(
		NewParser(parserLineSeparator),
		firstTestFailedBeforeRun,
	)

	assert.Equal(
		t,
		[]Event{
			TestsRunStartedEvent{
				NumberOfTests: 0,
			},
			TestsRunFailedEvent{
				Message: "Process crashed.",
			},
		},
		events,
	)
}

func TestSecondTestCrashDetectedAfterFirstPassedTestRun(t *testing.T) {
	events := startInstrumentationOutputParsing(
		NewParser(parserLineSeparator),
		testCrashedBeforeRunAfterFirstSuccess,
	)

	assert.Equal(
		t,
		[]Event{
			TestsRunStartedEvent{
				NumberOfTests: 2,
			},
			TestStartedEvent{
				Run: TestRun{
					Code:           1,
					NumTests:       2,
					TestName:       "test1",
					TestClass:      "com.example.test.TestClass",
					TestStackTrace: "",
				},
			},
			TestPassedEvent{
				Run: TestRun{
					Code:           0,
					NumTests:       2,
					TestName:       "test1",
					TestClass:      "com.example.test.TestClass",
					TestStackTrace: "",
				},
			},
			TestsRunFailedEvent{
				Message: "Process crashed.",
			},
		},
		events,
	)
}

func TestFirstCrashDetectedAfterStartingTest(t *testing.T) {
	events := startInstrumentationOutputParsing(
		NewParser(parserLineSeparator),
		testCrashedAfterFirstTestRun,
	)

	assert.Equal(
		t,
		[]Event{
			TestsRunStartedEvent{
				NumberOfTests: 1,
			},
			TestStartedEvent{
				Run: TestRun{
					Code:           1,
					NumTests:       1,
					TestName:       "test1",
					TestClass:      "com.example.test.TestClass",
					TestStackTrace: "",
				},
			},
			TestFailedEvent{
				Run: TestRun{
					Code:           -2,
					NumTests:       1,
					TestName:       "test1",
					TestClass:      "com.example.test.TestClass",
					TestStackTrace: "Process crashed.",
				},
			},
			TestsRunFailedEvent{
				Message: "Process crashed.",
			},
		},
		events,
	)
}

func TestSecondIgnoredTestFoundAfterPassedTest(t *testing.T) {
	events := startInstrumentationOutputParsing(
		NewParser(parserLineSeparator),
		secondTestIgnoredAfterOnePassedRun,
	)

	assert.Equal(
		t,
		[]Event{
			TestsRunStartedEvent{
				NumberOfTests: 2,
			},
			TestStartedEvent{
				Run: TestRun{
					Code:           1,
					NumTests:       2,
					TestName:       "test1",
					TestClass:      "com.example.test.TestClass",
					TestStackTrace: "",
				},
			},
			TestPassedEvent{
				Run: TestRun{
					Code:           0,
					NumTests:       2,
					TestName:       "test1",
					TestClass:      "com.example.test.TestClass",
					TestStackTrace: "",
				},
			},
			TestStartedEvent{
				Run: TestRun{
					Code:           1,
					NumTests:       2,
					TestName:       "test2",
					TestClass:      "com.example.test.TestClass",
					TestStackTrace: "",
				},
			},
			TestIgnoredEvent{
				Run: TestRun{
					Code:           -3,
					NumTests:       2,
					TestName:       "test2",
					TestClass:      "com.example.test.TestClass",
					TestStackTrace: "",
				},
			},
			TestsRunFinishedEvent{
				Time: 10.073,
			},
		},
		events,
	)
}

func TestTwoTestsPassedEventsDetected(t *testing.T) {
	events := startInstrumentationOutputParsing(
		NewParser(parserLineSeparator),
		twoTestsPassedRun,
	)

	assert.Equal(
		t,
		[]Event{
			TestsRunStartedEvent{
				NumberOfTests: 2,
			},
			TestStartedEvent{
				Run: TestRun{
					Code:           1,
					NumTests:       2,
					TestName:       "test1",
					TestClass:      "com.example.test.TestClass",
					TestStackTrace: "",
				},
			},
			TestPassedEvent{
				Run: TestRun{
					Code:           0,
					NumTests:       2,
					TestName:       "test1",
					TestClass:      "com.example.test.TestClass",
					TestStackTrace: "",
				},
			},
			TestStartedEvent{
				Run: TestRun{
					Code:           1,
					NumTests:       2,
					TestName:       "test2",
					TestClass:      "com.example.test.TestClass",
					TestStackTrace: "",
				},
			},
			TestPassedEvent{
				Run: TestRun{
					Code:           0,
					NumTests:       2,
					TestName:       "test2",
					TestClass:      "com.example.test.TestClass",
					TestStackTrace: "",
				},
			},
			TestsRunFinishedEvent{
				Time: 9.073,
			},
		},
		events,
	)
}

func TestInvalidPackageErrorDetected(t *testing.T) {
	events := startInstrumentationOutputParsing(
		NewParser(parserLineSeparator),
		passInvalidPackageRun,
	)

	assert.Equal(
		t,
		[]Event{
			TestsRunStartedEvent{
				NumberOfTests: 0,
			},
			TestsRunFailedEvent{
				Message: "Unable to find instrumentation info for: ComponentInfo{/tmp/go-build506757578/command-line-arguments/_obj/exe/main/com.avito.android.runner.AvitoInstrumentTestRunner}",
			},
		},
		events,
	)
}
