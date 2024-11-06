package util

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func RunCmd(timeout int, command string, args ...string) error {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return fmt.Errorf("cmd.CombinedOutput error: %w\n%s", err, out)
	}

	return nil
}
