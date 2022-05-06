package simple_test

import (
	"bufio"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	simple "github.com/jtagcat/simple/v2/pkg"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

var errBoo = errors.New("I don't like boos")

func exampleNestedFn(line string, linenum int, returnc chan string) error {
	if line == "boo" {
		return fmt.Errorf("%d: %w", linenum, errBoo)
	}
	if line == "four" {
		time.Sleep(2 * time.Second) // this is not aborted on error, todo
	}

	returnc <- "hello " + line
	return nil
}

func TestParallel(t *testing.T) {
	input := strings.NewReader("boo\none\ntwo\nthree\nfour")

	output, err := simple.Parallel(func(g *errgroup.Group, returnc chan string) error {
		scanner := bufio.NewScanner(input)
		for i := 1; scanner.Scan(); i++ {
			line, linenum := scanner.Text(), i
			g.Go(func() error {
				return exampleNestedFn(line, linenum, returnc)
			})
		}
		return scanner.Err()
	})
	sort.Strings(output)
	assert.Equal(t, []string{"hello four", "hello one", "hello three", "hello two"}, output)
	assert.ErrorContains(t, err, errBoo.Error())
}
