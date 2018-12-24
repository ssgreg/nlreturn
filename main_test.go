package main

import (
	"strings"
	"testing"

	"github.com/ssgreg/logf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var golden = []string{
	"test/main.go:14:3: return with no blank line before\n",
	"test/main.go:23:3: fallthrough with no blank line before\n",
	"test/main.go:27:3: break with no blank line before\n",
	"test/main.go:43:3: return with no blank line before\n",
	"test/main.go:59:3: return with no blank line before\n",
	"test/main.go:77:3: return with no blank line before\n",
	"test/main.go:98:3: return with no blank line before\n",
	"test/main.go:115:3: return with no blank line before\n",
	"test/main.go:132:3: return with no blank line before\n",
}

type TestWriter struct {
	t     *testing.T
	index int
}

func (w *TestWriter) Write(bytes []byte) (int, error) {
	str := string(bytes)
	str = str[strings.Index(str, "test/main.go"):]

	require.True(w.t, len(golden) > w.index, str)
	assert.Equal(w.t, golden[w.index], str)
	w.index++

	return 0, nil
}

func TestCheck(t *testing.T) {
	w := TestWriter{t, 0}
	handleImportPaths(logf.NewDisabledLogger(), &w, []string{"./test"})

	assert.True(t, len(golden) == w.index)
}
