package adb

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const (
	successInstrumentation = `
WARNING: linker: libdvm.so has text relocations. This is wasting memory and is a security risk. Please fix.
INSTRUMENTATION_STATUS: numtests=3
INSTRUMENTATION_STATUS: stream=
ru.avito.services.test.auth_agreements.AuthAgreementsTest:
INSTRUMENTATION_STATUS: id=AndroidJUnitRunner
INSTRUMENTATION_STATUS: test=agreementScreen_finishWithResultOK_whenClickOnProceed
INSTRUMENTATION_STATUS: class=ru.avito.services.test.auth_agreements.AuthAgreementsTest
INSTRUMENTATION_STATUS: current=1
INSTRUMENTATION_STATUS_CODE: 1
INSTRUMENTATION_STATUS: numtests=3
INSTRUMENTATION_STATUS: stream=.
INSTRUMENTATION_STATUS: id=AndroidJUnitRunner
INSTRUMENTATION_STATUS: test=agreementScreen_finishWithResultOK_whenClickOnProceed
INSTRUMENTATION_STATUS: class=ru.avito.services.test.auth_agreements.AuthAgreementsTest
INSTRUMENTATION_STATUS: current=1
INSTRUMENTATION_STATUS_CODE: 0

Time: 10.843

OK (3 tests)
	`
)

func TestSuccessInstrumentationResultParsed(t *testing.T) {
	result := newInstrumentationResultFromOutput(successInstrumentation)

	assert.Equal(
		t,
		result,
		&InstrumentationResult{
			Status: "OK",
			Running: 3,
			Passed: 3,
			Failure: 0,
			Output: successInstrumentation,
		},
	)
}
