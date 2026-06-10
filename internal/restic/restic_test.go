package restic

import (
	"errors"
	stdexec "os/exec"
	"strconv"
	"testing"
)

func exitErr(t *testing.T, code int) error {
	t.Helper()
	err := stdexec.Command("sh", "-c", "exit $1", "sh", strconv.Itoa(code)).Run()
	if err == nil {
		t.Fatalf("expected exit error")
	}
	return err
}

func TestNewExecErrorParsesExitError(t *testing.T) {
	err := newExecError(exitErr(t, 12), []byte(`{"message_type":"exit_error","code":12,"message":"wrong password or no key found"}
`))

	if err.ExitCode() != ExitCodeBadPassword {
		t.Fatalf("exit code = %d", err.ExitCode())
	}
	if !errors.Is(err, ErrBadPassword) {
		t.Fatalf("expected errors.Is(err, ErrBadPassword)")
	}
	if got := err.Error(); got != "restic exited with code 12: wrong password or no key found" {
		t.Fatalf("error = %q", got)
	}
}

func TestNewExecErrorParsesBackupErrorMessage(t *testing.T) {
	err := newExecError(exitErr(t, 1), []byte(`{"message_type":"error","error":{"message":"permission denied"},"during":"scan","item":"/secret"}
`))

	if got := err.Error(); got != "restic exited with code 1: /secret during scan: permission denied" {
		t.Fatalf("error = %q", got)
	}
}
