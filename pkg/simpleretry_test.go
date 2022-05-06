package simple

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/wait"
)

var testBackoff = wait.Backoff{
	Steps: 2,
}

func genericTest(t *testing.T, retryable bool, expectedErr error, expectedRuns int) {
	execCount := 0
	err := RetryOnError(testBackoff, func() (bool, error) {
		execCount++
		return retryable, expectedErr
	})

	assert.Equal(t, expectedErr, err)
	if expectedErr != nil {
		assert.Equal(t, expectedErr.Error(), err.Error())
	}
	assert.Equal(t, expectedRuns, execCount)
}

func genericAll(t *testing.T, retryable bool) {
	expectedRetries := 1
	if retryable {
		expectedRetries = testBackoff.Steps
	}

	genericTest(t, retryable, fmt.Errorf("xyz"), expectedRetries)
	genericTest(t, retryable, nil, 1)
	genericTest(t, retryable,
		wrappedError{retryable: true, wrapped: fmt.Errorf("xyz")},
		expectedRetries)
	genericTest(t, retryable,
		wrappedError{retryable: false, wrapped: fmt.Errorf("xyz")},
		expectedRetries)
}

func TestRetryable(t *testing.T) {
	genericAll(t, true)
}

func TestFatal(t *testing.T) {
	genericAll(t, false)
}
