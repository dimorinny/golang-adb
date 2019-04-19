// this file wrote according with ddmlib class for parsing instrumentation output:
// https://github.com/miracle2k/android-platform_sdk/blob/master/ddms/libs/ddmlib/src/com/android/ddmlib/testrunner/InstrumentationResultParser.java
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
	statusIgnored = -3
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
