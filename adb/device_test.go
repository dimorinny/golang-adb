package adb

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/dimorinny/adbaster/model"
)

const (
	testLineSeparator = "\n"

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
	device := newDeviceFromOutput(successDevice, testLineSeparator)

	assert.Equal(
		t,
		device,
		&model.Device{
			Arch:         "x86",
			Timezone:     "America/New_York",
			Sdk:          19,
			HeapSize:     "48m",
			BatteryLevel: 50,
		},
	)
}

func TestDeviceFieldEmptyIfItNotExists(t *testing.T) {
	deviceWithoutSdk := newDeviceFromOutput(successDeviceWithoutSdk, testLineSeparator)

	assert.Equal(
		t,
		deviceWithoutSdk,
		&model.Device{
			Arch:         "x86",
			Timezone:     "America/New_York",
			Sdk:          -1,
			HeapSize:     "48m",
			BatteryLevel: 50,
		},
	)
}

func TestDeviceEmptyIfResponseIsEmpty(t *testing.T) {
	deviceWithoutSdk := newDeviceFromOutput("", testLineSeparator)

	assert.Equal(
		t,
		deviceWithoutSdk,
		&model.Device{
			Arch:         "",
			Timezone:     "",
			HeapSize:     "",
			Sdk:          -1,
			BatteryLevel: -1,
		},
	)
}
