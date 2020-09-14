package codeforces

import (
	"codeforces_cli/pkg/client"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	megaByte int = 1000000
)

func parseProblem(contestID, index string, submit bool) (Problem, error) {
	var problem Problem
	problem.contestID = contestID
	problem.index = index
	problem.url = url.URL{
		Scheme: "https",
		Host:   "codeforces.com",
		Path:   "problemset/problem/" + contestID + "/" + index,
	}

	req, _ := http.NewRequest("GET", problem.url.String(), nil)
	root, err := client.GetRoot(req)
	if err != nil {
		return problem, err
	}

	problemStatement := root.GetNodeByTagAttr("div", "class", "problem-statement")

	// Parse time limit
	timeLimit := strings.Split(problemStatement.GetNodeByTagAttr("div", "class", "time-limit").GetText(), " ")
	if len(timeLimit) != 2 {
		return problem, errors.New("loadProblem: error parsing time limit")
	}
	timeLimitSeconds, err := strconv.ParseFloat(timeLimit[0], 64)
	if err != nil {
		return problem, errors.New("loadProblem: error converting time limit")
	}
	switch timeLimit[1] {
	case "second", "seconds":
		problem.timeLimit = time.Duration(timeLimitSeconds * float64(time.Second))
	default:
		return problem, errors.New("loadProblem: unknown time unit: " + timeLimit[1])
	}

	// Parse memory limit
	memoryLimit := strings.Split(problemStatement.GetNodeByTagAttr("div", "class", "memory-limit").GetText(), " ")
	if problem.memoryLimit, err = strconv.Atoi(memoryLimit[0]); err != nil {
		return problem, errors.New("loadProblem: Error converting memory limit")
	}
	switch memoryLimit[1] {
	case "megabytes":
		problem.memoryLimit *= megaByte
	default:
		return problem, errors.New("loadProblem: unknown memory unit: " + memoryLimit[1])
	}

	problem.inputFile = problemStatement.GetNodeByTagAttr("div", "class", "input-file").GetText()
	if problem.inputFile != "standard input" {
		return problem, errors.New("loadProblem: unknown input file: " + problem.inputFile)
	}

	problem.outputFile = problemStatement.GetNodeByTagAttr("div", "class", "output-file").GetText()
	if problem.outputFile != "standard output" {
		return problem, errors.New("loadProblem: unknown output file: " + problem.outputFile)
	}

	problem.inputs = problemStatement.GetNodesByTagAttr("div", "class", "input").GetNodesByTag("pre").GetText()
	problem.outputs = problemStatement.GetNodesByTagAttr("div", "class", "output").GetNodesByTag("pre").GetText()

	if len(problem.inputs) == 0 || len(problem.inputs) != len(problem.outputs) {
		return problem, fmt.Errorf("loadProblem: unexpected number of inputs and outputs: %d inputs and %d outputs", len(problem.inputs), len(problem.outputs))
	}

	problem.n = len(problem.inputs)

	problem.csrfToken = root.GetNodeByTagAttr("span", "class", "csrf-token").GetAttr("data-csrf")

	if problem.csrfToken == "" {
		return problem, errors.New("loadProblem: error parsing CSRF token")
	}

	problem.submit = submit
	// Load languages
	if problem.submit {
		languages := root.GetNodeByTagAttr("form", "class", "submitForm").GetNodeByTagAttr("select", "name", "programTypeId").GetNodesByTag("option")
		// fmt.Println(root.GetNodeByTagAttr("form", "class", "submitForm")) //.GetNodeByTagAttr("select", "name", "programTypeId").GetNodesByTag("option")

		for _, language := range languages {
			name := language.GetText()
			value := language.GetAttr("value")
			if name == "" || value == "" {
				return problem, errors.New("loadProblem: error parsing language")
			}
			problem.languages = append(problem.languages, Language{
				name: name,
				id:   value,
			})
		}

		if len(problem.languages) == 0 {
			return problem, errors.New("loadProblem: error parsing languages")
		}
	}

	return problem, nil
}
