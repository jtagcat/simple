package simple

import (
	"sync"

	"golang.org/x/sync/errgroup"
)

// Parallel wraps errgroup, and collects results to an UNORDERED slice.
// jc: I don't understand ctx enough to return the first error, dismiss results, and kill any running (sub)goroutines.
//
// fn usually consists of a for loop of some kind (ex: bufio.Scanner).
// Inside the for loop, it shall call g.Go(<something>).
//
// Don't forget to do i, xyz := xyz, hello before g.Go(func() error { return fooBar(i, xyz, returnc) })
// (shadowing) https://golang.org/doc/faq#closures_and_goroutines
func Parallel[T any](fn func(g *errgroup.Group, returnc chan T) error) (output []T, _ error) {
	g, returnc := new(errgroup.Group), make(chan T)

	fnErr := fn(g, returnc)

	var readOK sync.Mutex
	readOK.Lock()
	go func() {
		defer readOK.Unlock()
		for r := range returnc {
			output = append(output, r)
		}
	}()

	err := g.Wait()
	close(returnc)
	readOK.Lock() // wait for returnc to be fully flushed to output

	if err != nil {
		return output, err
	}

	return output, fnErr
}
