package sanguessuga

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// Report is a value object that holds information about the downloadable links of microdata
type Report struct {
	Href     string
	FilePath string
	Size     int64
}

// ScrapeReports returns a list of downloadable reports from a given page
func ScrapeReports(link string, selector string) ([]Report, error) {
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
	links := allZipFiles.FindAllStringSubmatch(document, -1)
	if links == nil {
		return []Report{}, fmt.Errorf("Unable to match any links")
	}

	reports := []Report{}
	for _, link := range links {
		reports = append(reports, Report{Href: link[1]})
	}

	return reports, nil
}

// DownloadReport downloads the report into a tmp folder in the current directory
func DownloadReport(report Report, dir string) (Report, error) {
	path := report.Href
	segments := strings.Split(path, "/")
	fileName := fmt.Sprintf("%s/%s", dir, segments[len(segments)-1])
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	resp, fetchErr := http.Get(report.Href)
	if fetchErr != nil {
		return Report{}, fetchErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Report{}, fmt.Errorf("Provided link returned unexpected status: %q", resp.StatusCode)
	}

	size, copyErr := io.Copy(file, resp.Body)
	if copyErr != nil {
		return Report{}, fmt.Errorf("Failed to save file into the disk: %q", copyErr)
	}

	report.Size = size
	report.FilePath = fileName

	return report, nil
}
