package launcher

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func RunProgram(name string, args ...string) error {

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start %q: %w", name, err)
	}
	fmt.Println("started, PID:", cmd.Process.Pid)

	if err := cmd.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("%q exited with code %d", name, exitErr.ExitCode())
		}
		return fmt.Errorf("wait %q: %w", name, err)
	}
	return nil
}