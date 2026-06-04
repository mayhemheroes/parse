package fuzz

import (
	"testing"

	"github.com/tdewolff/parse/v2"
)

// FuzzNumber is a native Go fuzz target for the parse.Number function.
func FuzzNumber(f *testing.F) {
	f.Add([]byte("1234"))
	f.Add([]byte("3.14"))
	f.Add([]byte("-1e10"))
	f.Fuzz(func(t *testing.T, data []byte) {
		_ = parse.Number(data)
	})
}
