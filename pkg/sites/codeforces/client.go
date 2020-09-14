package codeforces

import (
	"bytes"
	"codeforces_cli/pkg/client"
	"codeforces_cli/pkg/runner"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func login(username string, password string) error { // make sure loadConfig is called before this
	req, _ := http.NewRequest("GET", "https://codeforces.com/enter", nil)
	root, err := client.GetRoot(req)
	if err != nil {
		return errors.New("login: " + err.Error())
	}
	csrfToken := root.GetNodeByTagAttr("input", "name", "csrf_token").GetAttr("value")
	if csrfToken == "" {
		return errors.New("login: error parsing CSRF token")
	}
	res, err := client.Client.PostForm("https://codeforces.com/enter", url.Values{
		"csrf_token":    {csrfToken},
		"action":        {"enter"},
		"ftaa":          {""},
		"bfaa":          {""},
		"handleOrEmail": {username},
		"password":      {password},
	})
	if err != nil {
		return errors.New("login: " + err.Error())
	}
	if res.StatusCode != http.StatusFound {
		return fmt.Errorf("login: unexpected response status %s, please check login credentials.", res.Status)
	}
	return nil
}

func submit(problem Problem, code *runner.Code) error {
	b := &bytes.Buffer{}
	mpw := multipart.NewWriter(b)
	mpw.WriteField("csrf_token", problem.csrfToken)
	mpw.WriteField("ftaa", "")
	mpw.WriteField("bfaa", "")
	mpw.WriteField("action", "submitSolutionFormSubmitted")
	mpw.WriteField("submittedProblemIndex", problem.index)
	mpw.WriteField("source", "")
	found := false
	for _, language := range problem.languages {
		// Try to find language with prefix code.language
		if strings.HasPrefix(language.name, code.Language) {
			found = true
			mpw.WriteField("programTypeId", language.id)
			break
		}
	}
	if !found {
		return fmt.Errorf("Language %s is not avaiable for this problem", code.Language)
	}
	f, err := os.Open(code.FilePath)
	if err != nil {
		panic(err)
	}
	w, err := mpw.CreateFormFile("sourceFile", filepath.Base(code.FilePath))
	if err != nil {
		panic(err)
	}
	io.Copy(w, f)
	if err := f.Close(); err != nil {
		panic(err)
	}
	mpw.Close()

	formURL := problem.url
	formURL.RawQuery = "csrf_token=" + problem.csrfToken
	req, _ := http.NewRequest("POST", formURL.String(), b)
	req.Header.Set("Content-Type", mpw.FormDataContentType())
	_, err = client.Client.Do(req)
	return err
}
