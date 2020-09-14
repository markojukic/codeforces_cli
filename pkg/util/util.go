package util

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type UsageError struct {
	s string
}

func (e UsageError) Error() string {
	return e.s
}

func NewUsageError(s string) UsageError {
	return UsageError{s: s}
}

func DefineFileFlag(f *flag.FlagSet, t *string) {
	if f.Lookup("file") == nil {
		f.StringVar(t, "file", "", "`file` to run")
	}
}

func Init_file(file *string) error {
	var err error
	// Check if file is provided
	if *file == "" {
		return NewUsageError("File not provided")
	}

	// Convert fileFlag path to absolute
	if *file, err = filepath.Abs(*file); err != nil { // this should be done before chdir
		return err
	}

	// Check if file exists and is a file
	if fileInfo, err := os.Stat(*file); err != nil {
		return fmt.Errorf("%s doesn't exist", *file)
	} else if fileInfo.IsDir() {
		return fmt.Errorf("%s is a directory", *file)
	}

	return nil
}

func Init(tmpDir string) error {
	// Create a directory for temporary files
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return fmt.Errorf("Error creating directory %s: %s", tmpDir, err.Error())
	}

	// Change current working directory to executable's directory
	if executablePath, err := os.Executable(); err != nil {
		return err
	} else if err = os.Chdir(filepath.Dir(executablePath)); err != nil {
		return err
	}

	return nil
}
