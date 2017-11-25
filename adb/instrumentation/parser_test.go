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
			TestsRunFinishedEvent{},
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
