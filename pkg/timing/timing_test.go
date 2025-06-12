package timing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimingElapsed(t *testing.T) {
	t1 := Timestamp()
	t2 := t1 + calc(3, 27, 42)
	assert.Equal(t, Elapsed(t1, t2), "03:27:42")

	t2 = t1 + calc(0, 1, 1)
	assert.Equal(t, Elapsed(t1, t2), "00:01:01")
}

func calc(h, m, s int64) int64 {
	return (h)*3600 + (m)*60 + (s)
}
