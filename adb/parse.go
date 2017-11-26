package adb

import (
	"errors"
	"fmt"
	"github.com/dimorinny/adbaster/model"
	"github.com/dimorinny/adbaster/util"
	"strings"
)

func newDevicesIdentifiersFromOutput(output, lineSeparator string) []model.DeviceIdentifier {
	identifiers := []model.DeviceIdentifier{}

	for _, item := range strings.Split(output, lineSeparator) {
		trimmedItem := strings.TrimSpace(item)

		if strings.Contains(trimmedItem, "device") || strings.Contains(trimmedItem, "online") {
			device := strings.Split(item, "\t")

			if len(device) == 2 {
				identifiers = append(identifiers, model.DeviceIdentifier(device[0]))
			}
		}
	}

	return identifiers
}

func newDeviceFromOutput(output, lineSeparator string) *model.Device {
	items := map[string]interface{}{}

	for _, item := range strings.Split(output, lineSeparator) {
		splitItem := strings.Split(item, ": ")
		if len(splitItem) == 2 {
			key := strings.Trim(splitItem[0], "][")
			value := strings.Trim(splitItem[1], "][")
			items[key] = value
		}
	}
	return &model.Device{
		Arch:         util.GetStringWithDefault(items, "ro.product.cpu.abi", ""),
		Timezone:     util.GetStringWithDefault(items, "persist.sys.timezone", ""),
		HeapSize:     util.GetStringWithDefault(items, "dalvik.vm.heapsize", ""),
		Sdk:          util.GetIntWithDefault(items, "ro.build.version.sdk", -1),
		BatteryLevel: util.GetIntWithDefault(items, "status.battery.level_raw", -1),
	}
}

func detectErrorInClearApplicationDataOutput(output string) error {
	if strings.TrimSpace(output) != "Success" {
		return errors.New(
			fmt.Sprintf(
				"error during running clear application data with output: %s",
				output,
			),
		)
	}

	return nil
}
