package main

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	SRTM_BASE_URL = "https://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL3.003/"
)

func main() {
	fundURLs(SRTM_BASE_URL, "")
}

func fundURLs(baseUrl, path string) error {
	url := strings.Trim(baseUrl, "/") + "/" + strings.Trim(path, "/")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d := goquery.NewDocumentFromReader(resp.Body)
	_ = d

	return nil
}
