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

	currentTestResult,
	lastTestResult *TestResult

	currentKey   string
	currentValue bytes.Buffer

	// true if start of test has already been reported to listener
	testStartReported bool
	// true if the parser is parsing a line beginning with "INSTRUMENTATION_RESULT"
	inInstrumentationResultKey bool
	// true if the completion of the test run has been detected
	testRunFinished bool

	// stores key-value pairs under INSTRUMENTATION_RESULT header, these are printed at the
	// end of a test run, if applicable
	instrumentationResultBundle map[string]string

	// the number of tests currently run
	numTestsRun,
	// the number of tests expected to run
	testsExpected int

	result chan TestResult
}

func NewParser(lineSeparator string) *Parser {
	return &Parser{
		lineSeparator: lineSeparator,

		testStartReported:          false,
		inInstrumentationResultKey: false,
		testRunFinished:            false,

		currentTestResult: nil,
		lastTestResult:    nil,

		numTestsRun:   0,
		testsExpected: 0,

		result: make(chan TestResult),
	}
}

// processes the instrumentation test output from channel
func (p *Parser) Process(output <-chan string) <-chan TestResult {
	go func() {
		for line := range output {
			fmt.Println("Start processing line:")
			fmt.Println(line)
			p.processLine(line)
		}

		close(p.result)
	}()

	return p.result
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
		//	mTestRunFinished = true;
	} else if strings.HasPrefix(line, prefixTimeReport) {
		// TODO: parse time
	} else {
		if p.currentValueIsEmpty() {
			p.currentValue.WriteString(p.lineSeparator)
			p.currentValue.WriteString(line)
		} else {
			fmt.Println(
				fmt.Sprintf(
					"unrecognized line %s",
					line,
				),
			)
		}
	}
}

// stores the currently parsed key-value pair in the appropriate place
func (p *Parser) submitCurrentKeyValue() {
	if !p.currentValueIsEmpty() && !p.currentKeyIsEmpty() {
		key := p.currentKey
		value := p.currentValue.String()

		if p.inInstrumentationResultKey {
			if !isKnownKey(key) {
				p.instrumentationResultBundle[key] = value
			} else {
				p.handleTestRunFailed(value)
			}
		} else {
			currentTestResult := p.getCurrentTestResult()

			if key == keyClass {
				currentTestResult.TestClass = strings.TrimSpace(value)
			} else if key == keyTest {
				currentTestResult.TestName = strings.TrimSpace(value)
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
					currentTestResult.NumTests = numTests
				}
			} else if key == keyError {
				p.handleTestRunFailed(value)
			}
		}
	}
}

// parses out a status code result
func (p *Parser) parseStatusCode(line string) {
	value := line[len(prefixStatusCode):]
	currentTestResult := p.getCurrentTestResult()

	numTests, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println(
			fmt.Sprintf(
				"unexpected integer number of tests, received: %d",
				numTests,
			),
		)
	} else {
		currentTestResult.NumTests = numTests
	}

	p.reportResult(p.currentTestResult)
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

	if p.lastTestResult != nil &&
		p.lastTestResult.isComplete() &&
		p.lastTestResult.Code == statusStart {

		fmt.Println("test failed " + p.lastTestResult.TestName + message)
	}

	if !p.testStartReported {
		fmt.Println("test run started: 0")
	}

	fmt.Println("test run failed: " + message)
	fmt.Println("test run ended")
}

// reports a test result to the test run listener. Must be called when a individual test
// result has been fully parsed.
func (p *Parser) reportResult(result *TestResult) {
	if !result.isComplete() {
		fmt.Println(
			fmt.Sprintf(
				"invalid instrumentation status bundle: %s",
				*result,
			),
		)
	}

	p.reportTestRunStarted(result)

	switch result.Code {
	case statusStart:
		// TODO: handle testStarted
		fmt.Println("test started " + result.TestName)

	case statusFailure:
		// TODO: testFailed testEnded
		fmt.Println("test failed " + result.TestName + result.getTrace())
		p.numTestsRun++

	case statusError:
		// TODO: testFailed testEnded
		fmt.Println("test error " + result.TestName + result.getTrace())
		p.numTestsRun++

	case statusOk:
		fmt.Println("test passed " + result.TestName)
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
func (p *Parser) reportTestRunStarted(result *TestResult) {
	if !p.testStartReported && result.NumTests != -1 {
		fmt.Println("test run started " + result.TestName)

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
	p.lastTestResult = p.currentTestResult
	p.currentTestResult = nil
}

func (p *Parser) currentValueIsEmpty() bool { return len(p.currentValue.String()) > 0 }
func (p *Parser) currentKeyIsEmpty() bool   { return len(p.currentKey) > 0 }

func (p *Parser) getCurrentTestResult() *TestResult {
	if p.currentTestResult == nil {
		p.currentTestResult = newTestResult()
	}

	return p.currentTestResult
}
