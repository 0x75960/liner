package liner

import (
	"bufio"
	"io"
	"log"
	"sync"
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
func NewConcurrentLineProcessor(handler func(string) error, concurrency int, strategy ErrStrategy) func(io.Reader) {

	return func(in io.Reader) {

		var wg sync.WaitGroup
		semaphore := make(chan bool, concurrency)

		for l := range LinesIn(in) {

			wg.Add(1)
			semaphore <- true

			go func(line string) {

				defer func() {
					wg.Done()
					<-semaphore
				}()

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
		wg.Wait()
	}
}
