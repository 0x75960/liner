package liner

import (
	"bufio"
	"io"
	"log"
)

// ErrStrategy when encounting error
type ErrStrategy int

const (
	// IgnoreWhenError (no logging and continue)
	IgnoreWhenError ErrStrategy = iota

	// LogWhenError (continue)
	LogWhenError

	// StopWhenError (not continue)
	StopWhenError
)

// NewLineEnumerater from io.Reader
func NewLineEnumerater(r io.Reader) (o chan string) {

	o = make(chan string)

	scanner := bufio.NewScanner(r)

	go func(s *bufio.Scanner) {

		for s.Scan() {
			o <- s.Text()
		}

		close(o)

	}(scanner)

	return
}

// NewLineProcessor from handler and strategy
func NewLineProcessor(handler func(string) error, strategy ErrStrategy) func(io.Reader) {

	return func(in io.Reader) {

		for line := range NewLineEnumerater(in) {

			err := handler(line)

			if err == nil {
				// when no error
				continue
			}

			switch strategy {
			case IgnoreWhenError:
				continue

			case LogWhenError:
				log.Println(err)

			case StopWhenError:
				log.Fatalln(err)

			default:
				log.Println(err)
			}

		}

	}

}
