package fivews

import (
	"testing"
)

func TestNew(t *testing.T) {
	t.Error(New("error"))
}

func TestWrap(t *testing.T) {
	t.Error(Wrap("wrap", New("cause")))
}

func TestJoin(t *testing.T) {
	t.Error(Join("wrap", New("cause1"), Wrap("test", New("cause2"))))
}
