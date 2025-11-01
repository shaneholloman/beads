package main

import (
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"rsc.io/script"
	"rsc.io/script/scripttest"
)

func TestScripts(t *testing.T) {
	// Build the beads binary
	exeName := "beads"
	if runtime.GOOS == "windows" {
		exeName += ".exe"
	}
	exe := filepath.Join(t.TempDir(), exeName)
	if err := exec.Command("go", "build", "-o", exe, ".").Run(); err != nil {
		t.Fatal(err)
	}

	// Create minimal engine with default commands plus beads
	// Use longer timeout on Windows for slower process startup and I/O
	timeout := 2 * time.Second
	if runtime.GOOS == "windows" {
		timeout = 5 * time.Second
	}
	engine := script.NewEngine()
	engine.Cmds["beads"] = script.Program(exe, nil, timeout)

	// Run all tests
	scripttest.Test(t, context.Background(), engine, nil, "testdata/*.txt")
}
