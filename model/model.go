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
)

type (
	InstrumentationArguments map[string]string

	InstrumentationParams struct {
		TestPackage, Runner string
		Arguments           InstrumentationArguments
	}
)
