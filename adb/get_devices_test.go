package adb

import (
	"github.com/dimorinny/golang-adb/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	getDevicesTestLineSeparator = "\n"

	devicesOutput = `List of devices attached
emulator-5556	device
emulator-5554	device
	`

	devicesWithOfflineOutput = `List of devices attached
emulator-5556	device
emulator-5554	device
emulator-5557	offline
	`

	devicesOutputWithStartingDaemon = `List of devices attached
* daemon not running. starting it now at tcp:5037 *
* daemon started successfully *
emulator-5556	device
emulator-5554	device
	`

	emptyDevicesOutput = `List of devices attached

	`
)

func TestDevicesLoadedFromCommonAdbOutput(t *testing.T) {
	identifiers := newDevicesIdentifiersFromOutput(devicesOutput, getDevicesTestLineSeparator)

	assert.Equal(
		t,
		[]model.DeviceIdentifier{
			"emulator-5556",
			"emulator-5554",
		},
		identifiers,
	)
}

func TestDevicesLoadedFromAdbOutputWithOfflineDevice(t *testing.T) {
	identifiers := newDevicesIdentifiersFromOutput(devicesWithOfflineOutput, getDevicesTestLineSeparator)

	assert.Equal(
		t,
		[]model.DeviceIdentifier{
			"emulator-5556",
			"emulator-5554",
		},
		identifiers,
	)
}

func TestDevicesLoadedFromCommonAdbOutputWithStartingAdbDaemon(t *testing.T) {
	identifiers := newDevicesIdentifiersFromOutput(devicesOutputWithStartingDaemon, getDevicesTestLineSeparator)

	assert.Equal(
		t,
		[]model.DeviceIdentifier{
			"emulator-5556",
			"emulator-5554",
		},
		identifiers,
	)
}

func TestDevicesLoadingResultIsEmptyWhenDevicesNotExists(t *testing.T) {
	identifiers := newDevicesIdentifiersFromOutput(emptyDevicesOutput, getDevicesTestLineSeparator)

	assert.Empty(
		t,
		identifiers,
	)
}
