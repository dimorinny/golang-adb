package adb

import (
	"github.com/dimorinny/adbaster/util"
	"strings"
)

type (
	Device struct {
		Arch         string
		Timezone     string
		HeapSize     string
		Sdk          int
		BatteryLevel int
	}
)

func newDeviceFromOutput(output, lineSeparator string) *Device {
	items := map[string]interface{}{}

	for _, item := range strings.Split(output, lineSeparator) {
		splitItem := strings.Split(item, ": ")
		if len(splitItem) == 2 {
			key := strings.Trim(splitItem[0], "][")
			value := strings.Trim(splitItem[1], "][")
			items[key] = value
		}
	}
	return &Device{
		Arch:         util.GetStringWithDefault(items, "ro.product.cpu.abi", ""),
		Timezone:     util.GetStringWithDefault(items, "persist.sys.timezone", ""),
		HeapSize:     util.GetStringWithDefault(items, "dalvik.vm.heapsize", ""),
		Sdk:          util.GetIntWithDefault(items, "ro.build.version.sdk", -1),
		BatteryLevel: util.GetIntWithDefault(items, "status.battery.level_raw", -1),
	}
}
