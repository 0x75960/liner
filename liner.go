package liner

import (
	"bufio"
	"io"
	"log"

	"github.com/0x75960/lmttr"
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

// LinesIn from io.Reader
func LinesIn(r io.Reader) (o chan string) {

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

		for line := range LinesIn(in) {

			err := handler(line)

			if err == nil {
				// when no error
				continue
			}

			switch strategy {
			case IgnoreWhenError:
				continue

			case LogWhenError:
				log.Printf("got error \"%s\" in processing trailing line.\n%s\ncontinue processing.", err, line)

			case StopWhenError:
				log.Fatalf("got error \"%s\" in processing trailing line.\n%s\nstop processing.", err, line)

			default:
				log.Println(err)
			}

		}

	}

}

// NewConcurrentLineProcessor from handler and strategy
func NewConcurrentLineProcessor(handler func(string) error, concurrency uint, strategy ErrStrategy) func(io.Reader) {

	return func(in io.Reader) {

		lmt, _ := lmttr.NewLimitter(concurrency)

		for l := range LinesIn(in) {

			lmt.Start()

			go func(line string) {

				defer lmt.End()

				err := handler(line)

				if err == nil {
					// when no error
					return
				}

				switch strategy {
				case IgnoreWhenError:
					return

				case LogWhenError:
					log.Printf("got error \"%s\" in processing trailing line.\n%s\ncontinue processing.", err, line)

				case StopWhenError:
					log.Fatalf("got error \"%s\" in processing trailing line.\n%s\nstop processing.", err, line)

				default:
					log.Println(err)
				}

			}(l)
		}

		lmt.Wait()
	}
}
