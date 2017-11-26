package adb

type Config struct {
	AdbPath       string
	LineSeparator string
}

func NewConfig(adbPath, lineSeparator string) Config {
	return Config{
		AdbPath:       adbPath,
		LineSeparator: lineSeparator,
	}
}
