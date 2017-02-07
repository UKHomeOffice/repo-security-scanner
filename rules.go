package main

import (
	"fmt"
	"net/http"

	df "github.com/techjacker/diffence"
)

// Stats analyses a series of diff rule results
type Stats struct {
	Pass bool
	Info string
}

// DiffValidator validates diffs
type DiffValidator interface {
	Check(string) (df.Results, error)
	Stat(df.Results, string) Stats
}

type diffValidator struct {
	rules *[]df.Rule
}

func (d diffValidator) Check(url string) (df.Results, error) {
	resp, err := getHTTP(url)
	defer resp.Body.Close()
	if err != nil {
		return df.Results{}, err
	}
	return df.CheckDiffs(resp.Body, d.rules)
}

func (d diffValidator) Stat(ruleResults df.Results, commitID string) Stats {

	dirty := false
	info := fmt.Sprintf("Commit ID: %s\n", commitID)

	for k, v := range ruleResults {
		if len(v) > 0 {
			dirty = true
			info += fmt.Sprintf("File %s violates %d rules:\n", k, len(v))
			for _, r := range v {
				info += fmt.Sprintf("\n%s\n", r.String())
			}
		}
	}

	return Stats{
		Pass: !dirty,
		Info: info,
	}
}

func getHTTP(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, fmt.Errorf("Could not get diff: %s\n", err)
	}
	req.Header.Set("Accept", "application/vnd.github.VERSION.diff")
	return http.DefaultClient.Do(req)
}
