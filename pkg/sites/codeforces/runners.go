package codeforces

import (
	"codeforces_cli/pkg/runner"
)

// Flags used on codeforces: http://codeforces.com/blog/entry/79
var Runners = map[string]runner.Runner{
	GNUGPP11: runner.CompilerRunner("g++", "-std=gnu++11", "-o"),
	GNUGPP14: runner.CompilerRunner("g++", "-std=gnu++14", "-o"),
	GNUGPP17: runner.CompilerRunner("g++", "-std=gnu++17", "-o"),
	GO:       runner.CompilerRunner("go", "build", "-o"),
	PYTHON2:  runner.InterpreterRunner("python2"),
	PYTHON3:  runner.InterpreterRunner("python3"),
}
