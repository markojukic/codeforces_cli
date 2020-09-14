package main

import (
	"codeforces_cli/pkg/sites/codeforces"
	"codeforces_cli/pkg/sites/udebug"
	"codeforces_cli/pkg/util"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

const (
	tmpDir = "/tmp/codeforces_cli"
)

// https://golang.org/pkg/net/http for more error examples
var (
	// ErrSomethingHappened is returned by ...
	ErrSomethingHappened = errors.New("functionName: Something Happened")
)

var config struct {
	Codeforces codeforces.Config
	Udebug     udebug.Config
}

// type Submission struct {
// 	id      int64
// 	time    time.Time
// 	verdict string
// }

func defineSiteFlag(f *flag.FlagSet, t *string) {
	f.StringVar(t, "site", "", "`site`: codeforces or udebug")
}

/*
func submissions(contest, index string) {
	req, _ := http.NewRequest("GET", "https://codeforces.com/problemset/problem/"+contest+"/"+index, nil)
	doc, err := getRoot(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(doc.Find("#sidebar>div>table.rtable.smaller>tbody>tr").Text())
}
*/

// config flags

//var languageFlag = flag.String("language", "", "language, overwrites value provided in config, for valid languages see <url>")

func loadConfig() error {
	configFile, err := os.Open("config.json")
	if err != nil {
		return fmt.Errorf("Error opening config.json: %s", err)
	}
	defer configFile.Close()
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		return fmt.Errorf("Error parsing config.json: %s", err)
	}
	return nil
}

func main() {
	var err error
	var f = flag.NewFlagSet(os.Args[0], flag.ContinueOnError) // flag.ContinueOnError

	// Print err on return
	defer func() {
		if err != nil {
			// fmt.Println("ERROR:")
			fmt.Println(err)
			// Print usage for UsageErrors
			if _, ok := err.(util.UsageError); ok {
				f.Usage()
			}
		}
	}()

	// Define site flag
	var site string
	defineSiteFlag(f, &site)
	f.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		var f flag.FlagSet
		var t string
		defineSiteFlag(&f, &t)
		f.PrintDefaults()
	}

	// Define other flags
	var codeforcesFlags codeforces.Flags
	var udebugFlags udebug.Flags
	codeforces.DefineFlags(f, &codeforcesFlags)
	udebug.DefineFlags(f, &udebugFlags)

	// Parse flags
	if err = f.Parse(os.Args[1:]); err != nil {
		return
	}

	// Load config
	if err = loadConfig(); err != nil {
		return
	}

	switch site {
	case "codeforces": // Codeforces specific flags
		f.Usage = codeforces.Usage
		if err = codeforces.Run(tmpDir, config.Codeforces, codeforcesFlags); err != nil {
			return
		}
	case "udebug": // Udebug specific flags
		f.Usage = udebug.Usage
		if err = udebug.Run(tmpDir, config.Udebug, udebugFlags); err != nil {
			return
		}
	default: // Invalid site
		err = util.NewUsageError("Invalid site: " + site)
		return
	}
}
