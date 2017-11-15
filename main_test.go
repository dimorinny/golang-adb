package adbaster

import (
	"github.com/dimorinny/adbaster/adb"
	"log"
	"testing"
)

func TestMain1(t *testing.T) {
	config := adb.NewConfig("adb")
	client := adb.NewClient(config)

	identifiers, err := client.Devices()

	if err != nil {
		log.Fatal(err)
	}

	t.Logf("\n\n %s \n\n", identifiers)
}
