// this file wrote according with ddmlib class for parsing instrumentation output:
// https://github.com/miracle2k/android-platform_sdk/blob/master/ddms/libs/ddmlib/src/com/android/ddmlib/testrunner/InstrumentationResultParser.java
package instrumentation

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	lineSeparator string

	currentTestRun,
	latestTestRun *TestRun

	currentKey   string
	currentValue bytes.Buffer

	// true if start of test has already been reported to listener
	testStartReported bool
	// true if fail of test run has already been reported to listener
	testFailedReported bool
	// true if the parser is parsing a line beginning with "INSTRUMENTATION_RESULT"
	inInstrumentationResultKey bool
	// true if the completion of the test run has been detected
	testRunFinished bool

	// the number of tests currently run
	numTestsRun,
	// the number of tests expected to run
	testsExpected int

	// stream for handling instrumentation results (like starting, passing, failing)
	resultStream chan Event
	// stream for handling original instrumentation output (for debugging)
	instrumentationOutputStream chan string
}

func NewParser(lineSeparator string) *Parser {
	return &Parser{
		lineSeparator: lineSeparator,
	}
}

// processes the instrumentation test output from channel
func (p *Parser) Process(output <-chan string) (<-chan Event, <-chan string) {
	p.prepareParserStateForRun()

	go func() {
		for line := range output {
			p.instrumentationOutputStream <- line

			p.processLine(line)
		}

		if !p.testFailedReported {
			if !p.testStartReported {
				p.resultStream <- TestsRunStartedEvent{
					NumberOfTests: 0,
				}
			}

			p.resultStream <- TestsRunFinishedEvent{}
		}

		close(p.instrumentationOutputStream)
		close(p.resultStream)
	}()

	return p.resultStream, p.instrumentationOutputStream
}

// parse an individual instrumentation output line
//
// the start of a new status line (starts with Prefixes.STATUS or Prefixes.STATUS_CODE),
// and thus there is a new key=value pair to parse, and the previous key-value pair is
// finished.
//
// a continuation of the previous status (the "value" portion of the key has wrapped
// to the next line).
func (p *Parser) processLine(line string) {
	if strings.HasPrefix(line, prefixStatusCode) {
		p.submitCurrentKeyValue()
		p.inInstrumentationResultKey = false
		p.parseStatusCode(line)
	} else if strings.HasPrefix(line, prefixStatus) {
		p.submitCurrentKeyValue()
		p.inInstrumentationResultKey = false
		p.parseKey(line, len(prefixStatus))
	} else if strings.HasPrefix(line, prefixResult) {
		p.submitCurrentKeyValue()
		p.inInstrumentationResultKey = true
		p.parseKey(line, len(prefixResult))
	} else if strings.HasPrefix(line, prefixStatusFailed) || strings.HasPrefix(line, prefixCode) {
		p.submitCurrentKeyValue()
		p.inInstrumentationResultKey = false
		p.testRunFinished = true
	} else if !p.currentValueIsEmpty() {
		p.currentValue.WriteString(p.lineSeparator)
		p.currentValue.WriteString(line)
	}
}

// stores the currently parsed key-value pair in the appropriate place
func (p *Parser) submitCurrentKeyValue() {
	if !p.currentValueIsEmpty() && !p.currentKeyIsEmpty() {
		key := p.currentKey
		value := p.currentValue.String()

		if p.inInstrumentationResultKey {
			if key == keyShortMsg {
				p.handleTestRunFailed(value)
			}
		} else {
			currentTestRun := p.getCurrentTestRun()

			if key == keyClass {
				currentTestRun.TestClass = strings.TrimSpace(value)
			} else if key == keyTest {
				currentTestRun.TestName = strings.TrimSpace(value)
			} else if key == keyNumTests {
				numTests, err := strconv.Atoi(value)
				if err != nil {
					fmt.Println(
						fmt.Sprintf(
							"unexpected integer number of tests, received: %d",
							numTests,
						),
					)
				} else {
					currentTestRun.NumTests = numTests
				}
			} else if key == keyError {
				p.handleTestRunFailed(value)
			} else if key == keyStack {
				currentTestRun.TestStackTrace = value
			}
		}

		p.currentKey = ""
		p.currentValue = *bytes.NewBufferString("")
	}
}

// parses out a status code result
func (p *Parser) parseStatusCode(line string) {
	value := line[len(prefixStatusCode):]
	currentTestRun := p.getCurrentTestRun()

	code, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println(
			fmt.Sprintf(
				"unexpected integer number of tests, received: %d",
				code,
			),
		)
	} else {
		currentTestRun.Code = code
	}

	p.reportResult(p.currentTestRun)
	p.clearCurrentTestInfo()
}

// process a instrumentation run failure
func (p *Parser) handleTestRunFailed(errorMessage string) {
	var message string

	if len(errorMessage) > 0 {
		message = errorMessage
	} else {
		message = "unknown error"
	}

	if p.latestTestRun != nil && p.latestTestRun.isComplete() && p.latestTestRun.Code == statusStart {

		p.resultStream <- TestFailedEvent{
			Run: TestRun{
				Code:           statusFailure,
				NumTests:       p.latestTestRun.NumTests,
				TestName:       p.latestTestRun.TestName,
				TestClass:      p.latestTestRun.TestClass,
				TestStackTrace: message,
			},
		}
	}

	if !p.testStartReported {
		p.resultStream <- TestsRunStartedEvent{
			NumberOfTests: 0,
		}
	}

	p.resultStream <- TestsRunFailedEvent{
		Message: message,
	}

	p.testFailedReported = true
}

// reports a test result to the test run listener. Must be called when a individual test
// result has been fully parsed.
func (p *Parser) reportResult(result *TestRun) {
	if !result.isComplete() {
		fmt.Println("invalid instrumentation status bundle")
	}

	p.reportTestRunStarted(result)

	switch result.Code {
	case statusStart:
		p.resultStream <- TestStartedEvent{
			Run: *result,
		}

	case statusFailure:
		p.resultStream <- TestFailedEvent{
			Run: *result,
		}
		p.numTestsRun++

	case statusError:
		p.resultStream <- TestFailedEvent{
			Run: *result,
		}
		p.numTestsRun++

	case statusOk:
		p.resultStream <- TestPassedEvent{
			Run: *result,
		}
		p.numTestsRun++

	default:
		fmt.Println(
			fmt.Sprintf(
				"unknown status code received: %d",
				result.Code,
			),
		)
		p.numTestsRun++
	}
}

// reports the start of a test run, and the total test count, if it has not been previously
// reported
func (p *Parser) reportTestRunStarted(result *TestRun) {
	if !p.testStartReported && result.NumTests != -1 {

		p.resultStream <- TestsRunStartedEvent{
			NumberOfTests: result.NumTests,
		}

		p.testsExpected = result.NumTests
		p.testStartReported = true
	}
}

// parses the key from the current line
// expects format of "key=value"
func (p *Parser) parseKey(line string, keyStartPosition int) {
	keyEndPosition := strings.Index(line, "=")

	if keyEndPosition > -1 {
		p.currentKey = strings.TrimSpace(line[keyStartPosition:keyEndPosition])

		p.currentValue = *bytes.NewBufferString("")
		p.currentValue.WriteString(line[keyEndPosition+1:])
	}
}

// clear current test and save it to last test
func (p *Parser) clearCurrentTestInfo() {
	p.latestTestRun = p.currentTestRun
	p.currentTestRun = nil
}

func (p *Parser) currentValueIsEmpty() bool { return len(p.currentValue.String()) == 0 }
func (p *Parser) currentKeyIsEmpty() bool   { return len(p.currentKey) == 0 }

func (p *Parser) getCurrentTestRun() *TestRun {
	if p.currentTestRun == nil {
		p.currentTestRun = newTestRun()
	}

	return p.currentTestRun
}

func (p *Parser) prepareParserStateForRun() {
	p.currentTestRun = nil
	p.latestTestRun = nil

	p.currentKey = ""
	p.currentValue = *bytes.NewBufferString("")

	p.testStartReported = false
	p.testFailedReported = false
	p.inInstrumentationResultKey = false
	p.testRunFinished = false

	p.numTestsRun = 0
	p.testsExpected = 0

	p.resultStream = make(chan Event, 1000)
	p.instrumentationOutputStream = make(chan string, 1000)
}
