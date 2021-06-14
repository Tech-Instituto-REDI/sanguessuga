package sanguessuga

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

// Report is a value object that holds information about the downloadable links of microdata
type Report struct {
	Href    string
	Content string
}

// FetchReports returns a list of downloadable reports from a given page
func FetchReports(link string, selector string) ([]Report, error) {
	u, invalidURIErr := url.Parse(link)
	if invalidURIErr != nil {
		return []Report{}, invalidURIErr
	}

	resp, fetchErr := http.Get(u.String())
	if fetchErr != nil {
		return []Report{}, fetchErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return []Report{}, fmt.Errorf("Provided link returned unexpected status: %q", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Report{}, fmt.Errorf("Unable to load the page content into memory: %q", err)
	}

	document := string(body)

	allZipFiles := regexp.MustCompile("href=\"(.*\\.zip)\"")
	links := allZipFiles.FindAllString(document, -1)
	if links == nil {
		return []Report{}, fmt.Errorf("Unable to match any links")
	}

	reports := []Report{}
	for _, link := range links {
		reports = append(reports, Report{link, link})
	}

	return reports, nil
}
