package adbaster

import (
	"github.com/dimorinny/adbaster/adb"
	"github.com/dimorinny/adbaster/model"
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

	client.RunInstrumentationTests(
		first,
		model.InstrumentationParams{
			Package: "com.lol",
			Runner:  "test",
			Arguments: model.InstrumentationArguments{
				"test":  "test2",
				"test2": "test3",
			},
		},
	)
}
