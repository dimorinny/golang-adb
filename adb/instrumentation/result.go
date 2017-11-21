package instrumentation

type TestResult struct {
	Code     int
	NumTests int
	TestName,
	TestClass,
	TestStackTrace string
}

func newTestResult() *TestResult {
	return &TestResult{
		Code:           -1,
		NumTests:       0,
		TestName:       "",
		TestClass:      "",
		TestStackTrace: "",
	}
}

func (r *TestResult) isComplete() bool {
	return r.Code != -1 && len(r.TestName) > 0 && len(r.TestClass) > 0
}

// returns the stack trace of the current failed test, from the provided testInfo
func (r *TestResult) getTrace() string {
	if len(r.TestStackTrace) > 0 {
		return r.TestStackTrace
	} else {
		return "Unknown failure"
	}
}
