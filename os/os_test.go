package os

import (
	"os/exec"
	"testing"
)

func TestGetPsAuxCount(t *testing.T) {
	const in, out = "ps", 1
	if c := GetPsAuxCount(in); c < out {
		t.Errorf("GetPsAuxCount(%s) expected %d+, but actual %d", in, out, c)
	}

	const in2, out2 = "NONE_MATCH PROCESS", 0
	if c := GetPsAuxCount(in2); c != out2 {
		t.Errorf("GetPsAuxCount(%s) expected %d+, but actual %d", in2, out2, c)
	}
}

func TestIsFreeMemoryLessThan100MB(t *testing.T) {
	IsFreeMemoryLessThanMB(100)
}

func TestExecuteCmd(t *testing.T) {
	const expected = "hello\n"
	c := exec.Command("echo", "hello")

	if actual := ExecuteCmd(c, ""); actual != expected {
		t.Log(len(actual), len(expected))
		t.Errorf("ExecuteCmd(%s) expected %s, but actual %s", c.Path, expected, actual)
	}
}
