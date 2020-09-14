package codeforces

import (
	"codeforces_cli/pkg/runner"
	"codeforces_cli/pkg/util"
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

// Codeforces config
type Config struct {
	Username string
	Password string
}

// Codeforces-specific flags
type Flags struct {
	file     string
	contest  string
	problem  string
	username string
	password string
	submit   bool
}

// Define Codeforces-specific flags
func DefineFlags(f *flag.FlagSet, t *Flags) {
	util.DefineFileFlag(f, &t.file)
	f.StringVar(&t.contest, "contest", "", "`contest` (from problem's url)")
	if f.Lookup("problem") == nil {
		f.StringVar(&t.problem, "problem", "", "`problem` index (from problem's url)")
	}
	f.StringVar(&t.username, "username", "", "username, overwrites value provided in config")
	f.StringVar(&t.password, "password", "", "password, overwrites value provided in config")
	f.BoolVar(&t.submit, "submit", false, "`submit` the file if it passes test cases")
}

// Print usage of Codeforces-specific flags
func Usage() {
	fmt.Printf("Usage of %s --site codeforces:\n", os.Args[0])
	var f flag.FlagSet
	var t Flags
	DefineFlags(&f, &t)
	f.PrintDefaults()
}

func Run(tmpDir string, config Config, f Flags) error {
	if f.username != "" {
		config.Username = f.username
	}
	if f.password != "" {
		config.Password = f.password
	}
	if err := util.Init_file(&f.file); err != nil {
		return err
	}
	if config.Username == "" {
		return errors.New("Username not provided")
	}
	if config.Password == "" {
		return errors.New("Password not provided")
	}
	if f.contest == "" {
		return errors.New("Contest not provided")
	}
	if f.problem == "" {
		return errors.New("Problem not provided")
	}
	if err := util.Init(tmpDir); err != nil {
		return err
	}

	// Login, required for loading Problem.languages
	if f.submit {
		if err := login(config.Username, config.Password); err != nil {
			return err
		}
	}

	// Parse problem
	problem, err := parseProblem(f.contest, f.problem, f.submit)
	if err != nil {
		return err
	}

	// fmt.Println(problem.languages)

	fmt.Printf("Problem %s %s\n", problem.contestID, problem.index)
	fmt.Printf("Time limit: 	%v\n", problem.timeLimit)
	fmt.Printf("Memory limit: 	%dMB\n", problem.memoryLimit/megaByte)

	verdicts := make([]string, problem.n)
	times := make([]time.Duration, problem.n)

	// WaitGroup for parallel testing
	wg := &sync.WaitGroup{}

	code, err := NewCode(f.file)
	if err != nil {
		fmt.Println(err)
	}

	// Prepare for running
	r, ok := Runners[code.Language]
	if !ok {
		return errors.New("Invalid language.")
	}
	r.Init(code.FilePath, tmpDir)

	for i := 0; i < problem.n; i++ {
		wg.Add(1)
		go r.RunTestCase(problem.inputs, problem.outputs, problem.memoryLimit, problem.timeLimit, verdicts, times, i, wg)
	}
	wg.Wait()

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(tw, "Test case\tVerdict\tTime\t")
	for i := 0; i < problem.n; i++ {
		fmt.Fprintf(tw, "%d\t%s\t%v\t\n", i+1, verdicts[i], times[i])
	}
	tw.Flush()

	if f.submit {
		accepted := true
		for _, verdict := range verdicts {
			if verdict != runner.VerdictOK {
				accepted = false
			}
		}
		if accepted {
			if err := submit(problem, code); err != nil {
				return err
			}
		} else {
			fmt.Println("NOT ACCEPTED")
		}
	}
	return nil
}
