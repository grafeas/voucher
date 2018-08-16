package voucher

import (
	"bytes"
	"os/exec"
)

// RunShellCommand runs shell commands
func RunShellCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return stdout.String(), err
}
