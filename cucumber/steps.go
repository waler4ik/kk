package cucumber

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func InstallKK(ctx context.Context) error {
	return runCmdWithArgs("", "go", "install", ".")
}

func CleanFolder(ctx context.Context, folderName string) error {
	return os.RemoveAll(folderName)
}

func CreateProject(ctx context.Context, router, uri string) error {
	return runCmdWithArgs("", "kk", "init", "-r", router, "rest", uri)
}

func CreateWebSocket(ctx context.Context, path, rootPath string) error {
	return runCmdWithArgs(rootPath, "kk", "add", "ws", path)
}

func CreateRESTResource(ctx context.Context, path, rootPath string) error {
	return runCmdWithArgs(rootPath, "kk", "add", "resource", path)
}

func CompareFolder(ctx context.Context, source, target string) error {
	return diff(source, target)
}

func runCmdWithArgs(rootPath string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if rootPath != "" {
		cmd.Dir = rootPath
	}
	return cmd.Run()
}

func diff(source, target string) error {
	cmd := exec.Command("diff", "-r", source, target)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdoutpipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	slurp, err := io.ReadAll(stdout)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("wait: %w, slurp: %s", err, slurp)
	}
	return nil
}
