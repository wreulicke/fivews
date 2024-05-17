package fivews

import (
	"strings"
	"testing"
)

func firstLine(err error) string {
	m := err.Error()
	return m[:strings.Index(m, "\n")]
}

func TestNew(t *testing.T) {
	if firstLine(New("error")) != "error" {
		t.Error("New error message is not correct")
	}
}

func TestWrap(t *testing.T) {
	if firstLine(Wrap("wrap", New("cause"))) != "wrap: cause" {
		t.Error("Wrap error message is not correct")
	}
}

func TestJoin(t *testing.T) {
	if firstLine(Join("wrap", New("cause1"), Wrap("test", New("cause2")))) != "wrap: cause1: test: cause2" {
		t.Error("Join error message is not correct")
	}
}

func TestLastMessage(t *testing.T) {
	if _, m := LastMessage(Wrap("wrap", New("error"))); m != "wrap" {
		t.Error("LastMessage error message is not correct")
	}
}
