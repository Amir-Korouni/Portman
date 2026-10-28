package launcher

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var ErrStartupFailed = errors.New("program failed during startup")

const startupWindow = 2 * time.Second

func RunProgram(name string, port int, args ...string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("program %q not found: %w", name, err)
	}

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "PORT="+strconv.Itoa(port))

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start %q: %w", name, err)
	}
	fmt.Println("started, PID:", cmd.Process.Pid)

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err == nil {
			return nil
		}
		return fmt.Errorf("%w: %s", ErrStartupFailed, describeExit(name, err))
	case <- time.After(startupWindow):
		if err := <-done; err != nil {
			return errors.New(describeExit(name, err))
		}
		return nil
	}
}

func describeExit(name string, err error) string {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return fmt.Sprintf("%q exited with code %d", name, exitErr.ExitCode())
	}
	return fmt.Sprintf("wait %q: %v", name, err)
}