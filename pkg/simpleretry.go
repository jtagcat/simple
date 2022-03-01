package simpleretry

import (
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
)

// based on https://gist.github.com/mikroskeem/cabd826c88717c4ed743a67dea5712ae
type wrappedError struct {
	retryable bool
	wrapped   error
}

func (w wrappedError) Error() string {
	return w.wrapped.Error()
}
func (w wrappedError) Unwrap() error {
	return w.wrapped
}

func OnErrorSimple(backoff wait.Backoff, fn func() (bool, error)) error {
	err := retry.OnError(backoff, func(err error) bool {
		return err.(*wrappedError).retryable
	}, func() error {
		retryable, err := fn()
		if err == nil {
			return nil
		}
		return &wrappedError{retryable: retryable, wrapped: err}
	})

	if err == nil {
		return nil
	}
	return err.(*wrappedError).Unwrap()
}
