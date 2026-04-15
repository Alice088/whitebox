//go:build windows

package sys

import (
	"os/exec"
)

func SetSysProcAttr(cmd *exec.Cmd) {
}
