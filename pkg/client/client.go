package client

import (
	"codeforces_cli/pkg/parser"
	"errors"
	"net/http"
	"net/http/cookiejar"
)

var jar, _ = cookiejar.New(nil)

var Client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error { // disable redirects
		return http.ErrUseLastResponse
	},
	Jar: jar,
}

func GetRoot(req *http.Request) (*parser.Node, error) {
	res, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("Response status: " + res.Status)
	}
	return parser.Parse(res.Body)
}
