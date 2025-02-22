//go:build !windows && !plan9 && !js
// +build !windows,!plan9,!js

package daemon

import "golang.org/x/sys/unix"

// Make sure that files created by the daemon is not accessible to other users.
func setUmaskForDaemon() { unix.Umask(0077) }
