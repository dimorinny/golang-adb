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

var (
	knownKeys = map[string]struct{}{
		keyTest:     {},
		keyClass:    {},
		keyStack:    {},
		keyNumTests: {},
		keyError:    {},
		keyShortMsg: {},

		// unused
		"stream":  {},
		"id":      {},
		"current": {},
	}
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

func isKnownKey(key string) bool {
	_, ok := knownKeys[key]
	return ok
}
