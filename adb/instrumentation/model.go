package instrumentation

// output keys
const (
	keyTest     = "test"
	keyClass    = "class"
	keyStack    = "stack"
	keyNumTests = "numtests"
	keyError    = "Error"
	keyShortMsg = "shortMsg"
)

// output statuses
const (
	statusFailure = -2
	statusStart   = 1
	statusError   = -1
	statusOk      = 0
)

// output prefixes
const (
	prefixStatus       = "INSTRUMENTATION_STATUS: "
	prefixStatusCode   = "INSTRUMENTATION_STATUS_CODE: "
	prefixStatusFailed = "INSTRUMENTATION_FAILED: "
	prefixCode         = "INSTRUMENTATION_CODE: "
	prefixResult       = "INSTRUMENTATION_RESULT: "
	prefixTimeReport   = "Time: "
)

type TestResult struct {
	Code     int
	NumTests int
	TestName,
	TestClass,
	TestStackTrace string
}
