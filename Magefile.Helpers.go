// +build mage

package main

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gsmcwhirter/mlai-workshop/tools/magehelpers"
)

func envOrDefault(key, backup string) string {
	ret, ok := os.LookupEnv(key)
	if !ok {
		ret = backup
	}

	return ret
}

var announce = color.New(color.FgMagenta, color.Bold)
var cmdColor = color.New(color.FgCyan, color.Bold)

func run(cmd string, args ...string) error {
	return runWithEnv(nil, cmd, args...)
}

func runWithEnv(env map[string]string, cmd string, args ...string) error {
	_, _ = cmdColor.Printf("%s %s\n", cmd, strings.Join(args, " "))
	_, err := magehelpers.Exec(magehelpers.Options{
		Env:    env,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, cmd, args...)
	return err
}

func runNoEcho(cmd string, args ...string) error {
	_, err := magehelpers.Exec(magehelpers.Options{
		Env:    nil,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, cmd, args...)
	return err
}

func ok() {
	color.Green("--- ok\n\n")
}
