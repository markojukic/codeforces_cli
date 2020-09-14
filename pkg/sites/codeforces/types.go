package codeforces

import (
	"net/url"
	"time"
)

type Language struct {
	name string
	id   string
}

type Problem struct {
	contestID   string
	index       string
	timeLimit   time.Duration
	memoryLimit int
	inputFile   string
	outputFile  string
	n           int
	inputs      []string
	outputs     []string
	csrfToken   string
	languages   []Language
	url         url.URL
	submit      bool
}
