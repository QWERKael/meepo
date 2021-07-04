package executor

import "testing"

func TestExec(t *testing.T) {
	Exec("plugin/show.so", "Net", nil)
}