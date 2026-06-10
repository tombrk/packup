package restic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	stdexec "os/exec"
	"path/filepath"
	"strings"

	pkgexec "github.com/sh0rez/packup/internal/exec"
)

type Repository struct {
	Addr     string
	Password string
}

func Open(addr, pass string) (*Repository, error) {
	_, err := os.Stat(addr)
	if err != nil && !os.IsNotExist(err) {
		addr, err = filepath.Abs(addr)
		if err != nil {
			return nil, err
		}
	}

	r := &Repository{
		Addr:     addr,
		Password: pass,
	}
	return r, nil
}

const Bin = "restic"

const (
	EnvPrefix = "RESTIC_"

	EnvPass = EnvPrefix + "PASSWORD"
	EnvRepo = EnvPrefix + "REPOSITORY"
)

const (
	ExitCodeSourceNotFound     = 3
	ExitCodeRepositoryNotFound = 10
	ExitCodeRepositoryLocked   = 11
	ExitCodeBadPassword        = 12
	ExitCodeInterrupted        = 130
)

func (r *Repository) cmd(action string, argv []string) pkgexec.Cmd {
	cmd := pkgexec.Command(Bin, append([]string{action}, argv...)...)
	cmd.Env[EnvRepo] = r.Addr
	cmd.Env[EnvPass] = r.Password

	return cmd
}

type ExecError struct {
	err     error
	stderr  string
	code    int
	message string
}

func (e ExecError) Error() string {
	msg := e.message
	if msg == "" {
		msg = strings.TrimSpace(e.stderr)
	}
	if msg == "" && e.err != nil {
		msg = e.err.Error()
	}

	if e.code != 0 {
		return strings.TrimSpace(fmt.Sprintf("restic exited with code %d: %s", e.code, msg))
	}
	if e.err == nil {
		return msg
	}
	return strings.TrimSpace(fmt.Sprintf("%s.\n\n%s", e.err.Error(), e.stderr))
}

func (e ExecError) Unwrap() error  { return e.err }
func (e ExecError) Stderr() string { return e.stderr }
func (e ExecError) ExitCode() int  { return e.code }

func (e ExecError) Is(target error) bool {
	t, ok := target.(ExecError)
	return ok && t.code != 0 && e.code == t.code
}

var (
	ErrSourceNotFound     = ExecError{code: ExitCodeSourceNotFound}
	ErrRepositoryNotFound = ExecError{code: ExitCodeRepositoryNotFound}
	ErrRepositoryLocked   = ExecError{code: ExitCodeRepositoryLocked}
	ErrBadPassword        = ExecError{code: ExitCodeBadPassword}
	ErrInterrupted        = ExecError{code: ExitCodeInterrupted}
)

type resticExitError struct {
	MessageType string `json:"message_type"`
	Code        int    `json:"code"`
	Message     string `json:"message"`
}

type resticErrorMessage struct {
	MessageType string `json:"message_type"`
	Error       struct {
		Message string `json:"message"`
	} `json:"error"`
	During string `json:"during"`
	Item   string `json:"item"`
}

func newExecError(err error, stderr []byte) ExecError {
	ee := ExecError{
		err:    err,
		stderr: string(stderr),
	}

	var exitErr *stdexec.ExitError
	if errors.As(err, &exitErr) {
		if code := exitErr.ExitCode(); code >= 0 {
			ee.code = code
		}
	}

	for _, line := range bytes.Split(bytes.TrimSpace(stderr), []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] != '{' {
			continue
		}

		var exit resticExitError
		if json.Unmarshal(line, &exit) == nil && exit.MessageType == "exit_error" {
			if exit.Code != 0 {
				ee.code = exit.Code
			}
			ee.message = exit.Message
			continue
		}

		var msg resticErrorMessage
		if json.Unmarshal(line, &msg) == nil && msg.MessageType == "error" && msg.Error.Message != "" {
			if msg.Item != "" && msg.During != "" {
				ee.message = fmt.Sprintf("%s during %s: %s", msg.Item, msg.During, msg.Error.Message)
			} else {
				ee.message = msg.Error.Message
			}
		}
	}

	return ee
}

func (r *Repository) exec(action string, argv []string) ([]byte, error) {
	cmd := r.cmd(action, argv)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, newExecError(err, stderr.Bytes())
	}

	return stdout.Bytes(), nil
}
