//go:build !windows && !plan9
// +build !windows,!plan9

package sys

import (
	"io"
	"testing"

	"src.elv.sh/pkg/testutil"
)

func TestWaitForRead(t *testing.T) {
	r0, w0 := testutil.MustPipe()
	r1, w1 := testutil.MustPipe()
	defer closeAll(r0, w0, r1, w1)

	w0.WriteString("x")
	ready, err := WaitForRead(-1, r0, r1)
	if err != nil {
		t.Error("WaitForRead errors:", err)
	}
	if !ready[0] {
		t.Error("Want ready[0]")
	}
	if ready[1] {
		t.Error("Don't want ready[1]")
	}
}

func closeAll(files ...io.Closer) {
	for _, file := range files {
		file.Close()
	}
}
