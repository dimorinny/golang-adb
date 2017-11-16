package adbaster

import (
	"github.com/dimorinny/adbaster/adb"
	"log"
	"testing"
)

func TestMain1(t *testing.T) {
	config := adb.NewConfig("adb")
	client := adb.NewClient(config, true)

	identifiers, err := client.Devices()

	if err != nil {
		log.Fatal(err)
	}

	first := identifiers[0]

	err = client.Install(first, "avito-test.apk")
	if err != nil {
		log.Fatal(err)
	}
}
