//go:build windows

package main

import "os/exec"

func setProcessGroup(cmd *exec.Cmd) {}
