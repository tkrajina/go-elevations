package main

import (
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	SRTM_BASE_URL = "https://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL3.003/2000.02.11/"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	urls, err := fundURLs(http.DefaultClient, SRTM_BASE_URL, map[string]bool{}, 2)
	panicIfErr(err)

	var urlKeys []string
	for k := range urls {
		urlKeys = append(urlKeys, k)
	}
	sort.Strings(urlKeys)

	f, err := os.Create("geoelevations/urls_generated.go")
	panicIfErr(err)
	f.WriteString("package geoelevations\n")
	f.WriteString(`var urls = SrtmData{
`)
	f.WriteString(`	BaseUrl: "` + SRTM_BASE_URL + `",`)
	f.WriteString(`	Files:   map[string]string{
`)
	for _, urlKey := range urlKeys {
		f.WriteString(`"` + urlKey + `": "` + strings.Replace(strings.ReplaceAll(urls[urlKey], "http:", "https:"), SRTM_BASE_URL, "", 1) + `",
`)
	}
	f.WriteString(`},
}`)
	f.Close()
}

func fundURLs(client *http.Client, url string, visited map[string]bool, depth int) (urls map[string]string, err error) {
	if _, found := visited[url]; found {
		//fmt.Println("already visited", url)
		return
	}
	if depth > 3 {
		//fmt.Println("depth", depth)
		return
	}
	//fmt.Println("Parsing", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	visited[url] = true

	urls = map[string]string{}

	d, err := goquery.NewDocumentFromReader(resp.Body)

	var finalErr error
	d.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.HasSuffix(href, ".hgt.zip") {
			//fmt.Println(s.Text(), href)
			// Example: http://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL3.003/2000.02.11/N12W004.SRTMGL3.hgt.zip
			parts := strings.Split(href, "/")
			name := strings.Split(parts[len(parts)-1], ".")[0]
			//fmt.Println(name)
			urls[name] = href
		} else if strings.Contains(strings.ToLower(href), strings.ToLower(url)) {
			var u map[string]string
			u, err = fundURLs(client, href, visited, depth+1)
			if err != nil {
				finalErr = err
			}
			for k, v := range u {
				urls[k] = v
			}
		}
	})

	err = finalErr
	return
}
