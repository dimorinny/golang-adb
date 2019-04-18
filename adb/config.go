package adb

type Config struct {
	AdbPath string
}

func NewConfig(adbPath string) Config {
	return Config{
		AdbPath: adbPath,
	}
}
