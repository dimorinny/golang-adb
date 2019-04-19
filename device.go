package golangadb

import (
	"github.com/dimorinny/golang-adb/adb/instrumentation"
	"github.com/dimorinny/golang-adb/model"
)

type Device interface {
	Info() (*model.Device, error)
	Push(from, to string) error
	Pull(from, to string) error
	Install(from string) error
	ClearApplicationData(applicationPackage string) error
	Logcat() (
		<-chan string,
		error,
	)
	RunInstrumentationTests(params model.InstrumentationParams) (
		<-chan instrumentation.Event,
		<-chan string,
		error,
	)
}
