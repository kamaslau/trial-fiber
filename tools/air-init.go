package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	// Ensure we always reference ./.air.toml relative to project root
	airToml := filepath.Join("..", ".air.toml")
	// If running from project root, fallback to ./.air.toml
	if _, err := os.Stat(airToml); os.IsNotExist(err) {
		airToml = ".air.toml"
	}

	// Check if .air.toml exists
	if _, err := os.Stat(airToml); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, ".air.toml not found")
		os.Exit(1)
	}

	// Read .air.toml
	content, err := os.ReadFile(airToml)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read .air.toml: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(content), "\n")

	// Prepare new bin/cmd values
	var bin, cmd string
	if runtime.GOOS == "windows" {
		bin = `bin = "tmp\\main.exe"`
		cmd = `cmd = "go build -o ./tmp/main.exe ./src"`
	} else {
		bin = `bin = "tmp/main"`
		cmd = `cmd = "go build -o ./tmp/main ./src"`
	}

	fmt.Printf("Current OS is %s, build.bin should be %q, build.cmd should be %q before replacing\n", runtime.GOOS, bin, cmd)

	// Replace build.bin and build.cmd
	inBuild := false
	changed := false
	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "[build]") {
			inBuild = true
			continue
		}
		if inBuild {
			if strings.HasPrefix(trim, "bin =") {
				if trim == bin {
					fmt.Println("build.bin is already correct, skip")
				} else {
					lines[i] = "  " + bin
					changed = true
				}
			}
			if strings.HasPrefix(trim, "cmd =") {
				if trim == cmd {
					fmt.Println("build.cmd is already correct, skip")
				} else {
					lines[i] = "  " + cmd
					changed = true
				}
				inBuild = false // assume bin/cmd are always together
			}
		}
	}

	if !changed {
		fmt.Println("no changes needed")
		return
	}

	// Write back
	err = os.WriteFile(airToml, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write .air.toml: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(".air.toml updated")
}
