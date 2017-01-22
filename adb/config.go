package adb

type Config struct {
	Host    string
	Port    int
	AdbPath string
}

func NewConfig(host string, port int, adbPath string) Config {
	return Config{
		Host:    host,
		Port:    port,
		AdbPath: adbPath,
	}
}
