package runner

import (
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	VerdictOK  = "OK"
	VerdictWA  = "Wrong Answer"
	VerdictTLE = "Time Limit Exceeded"
)

type Code struct {
	FilePath string
	Language string
}

// Runner:
//	- Init: initializes the program for running, e.g. compiling. Returns output and errors
//	- Command: returns the command that will be used to run the program
// 	- Cleanup: cleans up, e.g. deletes the compiled binary
type Runner struct {
	Init    func(codePath string, directory string) ([]byte, error)
	Command func() []string
	Cleanup func() error
}

// Runner for compiled languages:
//	- Init: compiles by appending executablePath and codePath to the given command
//	- Command: runs compiled binary
//	- Cleanup: deletes compiled binary
func CompilerRunner(name string, arg ...string) Runner {
	var executablePath string
	return Runner{
		Init: func(codePath string, directory string) ([]byte, error) {
			executablePath = filepath.Join(directory, "run")
			arg = append(arg, executablePath, codePath)
			return exec.Command(name, arg...).Output()
		},
		Command: func() []string {
			return []string{executablePath}
		},
		Cleanup: func() error {
			return os.Remove(executablePath)
		},
	}
}

// Runner for interpreted languages:
//	- Init: does nothing
//	- Command: runs by appending codePath to the given command
//	- Cleanup: does nothing
func InterpreterRunner(cmd ...string) Runner {
	return Runner{
		Init: func(codePath string, directory string) ([]byte, error) {
			cmd = append(cmd, codePath)
			return nil, nil
		},
		Command: func() []string {
			interpreter, err := exec.LookPath(cmd[0])
			if err != nil {
				panic(err)
			}
			return append([]string{interpreter}, cmd[1:]...)
		},
		Cleanup: func() error {
			return nil
		},
	}
}

// Compares non-whitespace parts of two strings
func outputEqual(output1, output2 string) bool {
	o1 := strings.Fields(output1)
	o2 := strings.Fields(output2)
	n := len(o1)
	if len(o2) != n {
		return false
	}
	for i := 0; i < n; i++ {
		if o1[i] != o2[i] {
			return false
		}
	}
	return true
}

func (runner *Runner) RunTestCase(inputs []string, outputs []string, memoryLimit int, timeLimit time.Duration, verdicts []string, times []time.Duration, i int, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := exec.Command(
		"./executor",
		append(
			[]string{
				strconv.Itoa(memoryLimit),
				strconv.FormatFloat(math.Ceil(timeLimit.Seconds()), 'f', -1, 64),
			},
			runner.Command()...,
		)...,
	)

	if stdin, err := cmd.StdinPipe(); err != nil {
		log.Fatal(err)
	} else {
		stdin.Write([]byte(inputs[i]))
		stdin.Close()
	}

	output, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	verdict := ""
	totalTime := cmd.ProcessState.UserTime() + cmd.ProcessState.SystemTime()
	if err != nil {
		signal := cmd.ProcessState.Sys().(syscall.WaitStatus).Signal() // not portable
		if signal == syscall.SIGKILL {
			verdict = VerdictTLE
		} else {
			verdict = signal.String()
		}
	} else if totalTime > timeLimit {
		verdict = VerdictTLE
	} else {
		if outputEqual(string(output), outputs[i]) {
			verdict = VerdictOK
		} else {
			verdict = VerdictWA
		}
	}

	verdicts[i] = verdict
	times[i] = totalTime
}
