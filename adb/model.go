package adb

import (
	"github.com/dimorinny/adbaster/model"
	"github.com/dimorinny/adbaster/util"
	"regexp"
	"strconv"
	"strings"
)

const (
	instrumentalOk   = "OK"
	instrumentalFail = "FAIL"
)

var (
	okMatcher      = regexp.MustCompile(`OK \((\d+) tests\)`)
	failureMatcher = regexp.MustCompile(`Tests run: (\d+),  Failures: (\d+)`)
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

func newInstrumentationResultFromOutput(output string) *model.InstrumentationResult {
	result := model.InstrumentationResult{}

	if okMatcher.MatchString(output) {

		testsPassed, _ := strconv.Atoi(okMatcher.FindStringSubmatch(output)[1])
		result.Status = instrumentalOk
		result.Failure = 0
		result.Passed = testsPassed
		result.Running = testsPassed
	} else if failureMatcher.MatchString(output) {

		matcherValue := failureMatcher.FindStringSubmatch(output)[1:]
		testsRunning, _ := strconv.Atoi(matcherValue[0])
		testsFailed, _ := strconv.Atoi(matcherValue[1])
		result.Status = instrumentalFail
		result.Failure = testsFailed
		result.Passed = testsRunning - testsFailed
		result.Running = testsRunning
	}

	result.Output = output
	return &result
}
