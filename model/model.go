package model

type (
	DeviceIdentifier string

	Device struct {
		Identifier   DeviceIdentifier
		Arch         string
		Timezone     string
		HeapSize     string
		Sdk          int
		BatteryLevel int
	}

	InstrumentationParams struct {
		From, Runner, TestClass string
	}

	InstrumentationResult struct {
		Status  string
		Running int
		Passed  int
		Failure int
		Output  string
	}
)
