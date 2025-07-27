package utils

import (
    "fmt"
    "os"
    "os/exec"
)

// RunCommand runs an external command in the given directory and streams its
// output to the parent process' stdio. The working directory can be left empty
// to inherit the current directory.
func RunCommand(dir string, name string, args ...string) error {
    Log("▶️  Running: %s %v", name, args)
    cmd := exec.Command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    if dir != "" {
        cmd.Dir = dir
    }
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("command %s failed: %w", name, err)
    }
    return nil
}
