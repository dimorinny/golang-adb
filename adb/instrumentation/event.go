package instrumentation

type (
	// base type for instrumentation run events
	Event interface{}

	// tests running events
	TestsRunStartedEvent  struct{ NumberOfTests int }
	TestsRunFailedEvent   struct{ Message string }
	TestsRunFinishedEvent struct{}

	// test running events
	TestStartedEvent struct{ Run TestRun }
	TestPassedEvent  struct{ Run TestRun }
	TestFailedEvent  struct{ Run TestRun }
	TestIgnoredEvent struct{ Run TestRun }
)
