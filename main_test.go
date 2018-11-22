package main

import (
	"strings"
	"testing"

	"github.com/ssgreg/logf"
	"github.com/stretchr/testify/assert"
)

var golden = []string{
	"test/main.go:14:3: return with no blank line before\n",
	"test/main.go:23:3: fallthrough with no blank line before\n",
	"test/main.go:27:3: break with no blank line before\n",
	"test/main.go:43:3: return with no blank line before\n",
}

type TestWriter struct {
	t     *testing.T
	index int
}

func (w *TestWriter) Write(bytes []byte) (int, error) {
	str := string(bytes)
	str = str[strings.Index(str, "test/main.go"):]

	assert.Equal(w.t, golden[w.index], str)
	w.index++

	return 0, nil
}

func TestCheck(t *testing.T) {
	handleImportPaths(logf.NewDisabledLogger(), &TestWriter{t, 0}, []string{"./test"})
}
