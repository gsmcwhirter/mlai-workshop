package magehelpers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Options struct {
	Env    map[string]string
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

func expand(env map[string]string) func(string) string {
	return func(s string) string {
		if strings.HasPrefix(s, "$") {
			return s
		}

		s2, ok := env[s]
		if ok {
			return s2
		}
		return os.Getenv(s)
	}
}

func Exec(opts Options, cmd string, args ...string) (ran bool, err error) {
	exp := expand(opts.Env)
	cmd = os.Expand(cmd, exp)
	for i := range args {
		args[i] = os.Expand(args[i], exp)
	}
	ran, code, err := run(opts, cmd, args...)

	if err == nil {
		return true, nil
	}

	if ran {
		return ran, mg.Fatalf(code, `running "%s %s" failed with exit code %d`, cmd, strings.Join(args, " "), code)
	}

	return ran, fmt.Errorf(`failed to run "%s %s: %v"`, cmd, strings.Join(args, " "), err)
}

func ExecOutput(opts Options, cmd string, args ...string) (string, error) {
	b := &bytes.Buffer{}
	opts.Stdout = b
	_, err := Exec(opts, cmd, args...)
	return strings.TrimSuffix(b.String(), "\n"), err
}

func Output(cmd string, args ...string) (string, error) {
	return ExecOutput(Options{
		Stderr: os.Stderr,
	}, cmd, args...)
}

func OutputWithEnv(env map[string]string, cmd string, args ...string) (string, error) {
	return ExecOutput(Options{
		Env:    env,
		Stderr: os.Stderr,
	}, cmd, args...)
}

func run(opts Options, cmd string, args ...string) (ran bool, code int, err error) {
	c := exec.Command(cmd, args...)

	c.Env = os.Environ()
	for k, v := range opts.Env {
		c.Env = append(c.Env, k+"="+v)
	}

	c.Stderr = opts.Stderr
	c.Stdout = opts.Stdout
	c.Stdin = opts.Stdin

	if opts.Stdin == nil {
		c.Stdin = os.Stdin
	}

	err = c.Run()
	return cmdRan(err), sh.ExitStatus(err), err
}

func cmdRan(err error) bool {
	if err == nil {
		return true
	}
	ee, ok := err.(*exec.ExitError)
	if ok {
		return ee.Exited()
	}
	return false
}
