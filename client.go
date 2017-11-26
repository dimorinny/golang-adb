package adbaster

import (
	"github.com/dimorinny/adbaster/adb/instrumentation"
	"github.com/dimorinny/adbaster/model"
)

type Client interface {
	Devices() ([]model.DeviceIdentifier, error)
	DeviceInfo(device model.DeviceIdentifier) (*model.Device, error)
	Push(device model.DeviceIdentifier, from, to string) error
	Pull(device model.DeviceIdentifier, from, to string) error
	Install(device model.DeviceIdentifier, from string) error
	ClearApplicationData(device model.DeviceIdentifier, applicationPackage string) error
	Logcat(device model.DeviceIdentifier) (
		<-chan string,
		error,
	)
	RunInstrumentationTests(device model.DeviceIdentifier, params model.InstrumentationParams) (
		<-chan instrumentation.Event,
		<-chan string,
		error,
	)
}
