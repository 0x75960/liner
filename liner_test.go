package liner

import (
	"bytes"
	"testing"
)

func TestLinesIn(t *testing.T) {
	tt := bytes.NewBufferString("a\nb\nc\n")

	res := []string{}

	for a := range LinesIn(tt) {
		res = append(res, a)
	}

	if len(res) != 3 {
		t.Fatalf("expected 3 items. but result has %d \n %s", len(res), res)
	}

	ans := map[string]bool{}

	for _, i := range res {
		ans[i] = true
	}

	if (!ans["a"]) || (!ans["b"]) || (!ans["c"]) {
		t.Fatalf("expected [a, b, c], but got %s", res)
	}

}

func TestNewLineProcessor(t *testing.T) {
	tt := bytes.NewBufferString("a\nb\nc\n")
	res := []string{}

	handler := func(line string) (err error) {
		res = append(res, line)
		return
	}

	p := NewLineProcessor(handler, IgnoreWhenError)

	p(tt)

	if len(res) != 3 {
		t.Fatalf("expected 3 items. but result has %d \n %s", len(res), res)
	}

	ans := map[string]bool{}

	for _, i := range res {
		ans[i] = true
	}

	if (!ans["a"]) || (!ans["b"]) || (!ans["c"]) {
		t.Fatalf("expected [a, b, c], but got %s", res)
	}
}

func TestNewConcurrentLineProcessor(t *testing.T) {
	tt := bytes.NewBufferString("a\nb\nc\n")
	res := []string{}

	handler := func(line string) (err error) {
		res = append(res, line)
		return
	}

	p := NewConcurrentLineProcessor(handler, 4, IgnoreWhenError)

	p(tt)

	if len(res) != 3 {
		t.Fatalf("expected 3 items. but result has %d \n %s", len(res), res)
	}

	ans := map[string]bool{}

	for _, i := range res {
		ans[i] = true
	}

	if (!ans["a"]) || (!ans["b"]) || (!ans["c"]) {
		t.Fatalf("expected [a, b, c], but got %s", res)
	}
}
