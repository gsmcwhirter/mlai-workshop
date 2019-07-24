// +build mage

package main

import "github.com/magefile/mage/mg"

type Test mg.Namespace

func (t Test) All() {
	mg.SerialDeps(t.Format, t.Lint, t.Run)
}

func (Test) Run() error {
	return runNoEcho("go", "test", "-cover", "-coverprofile=./cover.profile", "./...")
}

func (Test) Format() error {
	if err := run("bash", "-c", "for d in $$(go list -f {{.Dir}} ./...); do gofmt -s -w $$d/*.go; done"); err != nil {
		return err
	}

	if err := run("bash", "-c", "for d in $$(go list -f {{.Dir}} ./...); do goimports -w -local "+Project+" $$d/*.go; done"); err != nil {
		return err
	}

	return nil
}

func (Test) Lint() error {
	if err := run("golangci-lint", "run", "-E", "golint,gosimple,staticcheck", "./..."); err != nil {
		return err
	}

	if err := run("golangci-lint", "run", "-E", "deadcode,depguard,errcheck,gocritic,gofmt,goimports,gosec,govet,ineffassign,nakedret,prealloc,structcheck,typecheck,unconvert,varcheck", "./..."); err != nil {
		return err
	}

	return nil
}

func (Test) Benchmark() error {
	return runNoEcho("go", "test", "-bench=.", "-benchmem", "./...")
}
