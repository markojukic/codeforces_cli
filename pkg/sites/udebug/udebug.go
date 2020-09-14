package udebug

import (
	"codeforces_cli/pkg/util"
	"flag"
	"fmt"
	"os"
)

// Udebug-specific flags
type Config struct {
	Username string
	Password string
}

// Udebug-specific flags
type Flags struct {
	file    string
	judge   string
	problem string
}

// Define Udebug-specific flags
func DefineFlags(f *flag.FlagSet, t *Flags) {
	util.DefineFileFlag(f, &t.file)
	f.StringVar(&t.file, "judge", "", "name of judge")
	if f.Lookup("problem") == nil {
		f.StringVar(&t.problem, "problem", "", "problem")
	}
}

// Print usage of Udebug-specific flags
func Usage() {
	fmt.Printf("Usage of %s --site udebug:\n", os.Args[0])
	var f flag.FlagSet
	var t Flags
	DefineFlags(&f, &t)
	f.PrintDefaults()
}

func Run(tmpDir string, config Config, f Flags) error {
	return nil
}
