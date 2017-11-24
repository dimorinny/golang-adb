package instrumentation

type TestRun struct {
	Code     int
	NumTests int
	TestName,
	TestClass,
	TestStackTrace string
}

func newTestRun() *TestRun {
	return &TestRun{
		Code:           -1,
		NumTests:       0,
		TestName:       "",
		TestClass:      "",
		TestStackTrace: "",
	}
}

func (r *TestRun) isComplete() bool {
	return r.Code != -1 && len(r.TestName) > 0 && len(r.TestClass) > 0
}

// returns the stack trace of the current failed test, from the provided testInfo
func (r *TestRun) getTrace() string {
	if len(r.TestStackTrace) > 0 {
		return r.TestStackTrace
	} else {
		return "Unknown failure"
	}
}
