package adb

import (
	"github.com/dimorinny/golang-adb/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	deviceInfoTestLineSeparator = "\n"
	identifier                  = "identifier"

	successDevice = `[ARGH]: [ARGH]
[dalvik.vm.heapsize]: [48m]
[dalvik.vm.stack-trace-file]: [/data/anr/traces.txt]
[init.svc.adbd]: [running]
[init.svc.bootanim]: [running]
[init.svc.console]: [running]
[init.svc.debuggerd]: [running]
[init.svc.drm]: [running]
[ro.product.cpu.abi]: [x86]
[init.svc.goldfish-logcat]: [stopped]
[init.svc.goldfish-setup]: [stopped]
[persist.sys.timezone]: [America/New_York]
[ro.build.version.sdk]: [19]
[status.battery.level_raw]: [50]
[dalvik.vm.heapsize]: [48m]`

	successDeviceWithoutSdk = `[ARGH]: [ARGH]
[dalvik.vm.heapsize]: [48m]
[dalvik.vm.stack-trace-file]: [/data/anr/traces.txt]
[init.svc.adbd]: [running]
[init.svc.bootanim]: [running]
[init.svc.console]: [running]
[init.svc.debuggerd]: [running]
[init.svc.drm]: [running]
[ro.product.cpu.abi]: [x86]
[init.svc.goldfish-logcat]: [stopped]
[init.svc.goldfish-setup]: [stopped]
[persist.sys.timezone]: [America/New_York]
[status.battery.level_raw]: [50]
[dalvik.vm.heapsize]: [48m]`
)

func TestDeviceInfoParsed(t *testing.T) {
	device := newDeviceFromOutput(
		identifier,
		successDevice,
		deviceInfoTestLineSeparator,
	)

	assert.Equal(
		t,
		device,
		&model.Device{
			Identifier:   identifier,
			Arch:         "x86",
			Timezone:     "America/New_York",
			Sdk:          19,
			HeapSize:     "48m",
			BatteryLevel: 50,
		},
	)
}

func TestDeviceFieldEmptyIfItNotExists(t *testing.T) {
	deviceWithoutSdk := newDeviceFromOutput(
		identifier,
		successDeviceWithoutSdk,
		deviceInfoTestLineSeparator,
	)

	assert.Equal(
		t,
		deviceWithoutSdk,
		&model.Device{
			Identifier:   identifier,
			Arch:         "x86",
			Timezone:     "America/New_York",
			Sdk:          -1,
			HeapSize:     "48m",
			BatteryLevel: 50,
		},
	)
}

func TestDeviceEmptyIfResponseIsEmpty(t *testing.T) {
	deviceWithoutSdk := newDeviceFromOutput(
		identifier,
		"",
		deviceInfoTestLineSeparator,
	)

	assert.Equal(
		t,
		deviceWithoutSdk,
		&model.Device{
			Identifier:   identifier,
			Arch:         "",
			Timezone:     "",
			HeapSize:     "",
			Sdk:          -1,
			BatteryLevel: -1,
		},
	)
}
