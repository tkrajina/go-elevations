package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var data = []struct {
	name, description, baseUrl string
}{
	{"SRTMGL3", "The default 3-arc-second data for the world obtained by averaging the 1-arc-second raw data.", "https://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL3.003/2000.02.11/"},
	{"SRTMGL3S", "The sampled 3-arc-second data for the whole world obtained by getting the middle 1-arc-second raw data sample out of a 3×3 matrix.", "https://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL3S.003/2000.02.11/"},
	//{"SRTMGL3N", "The meta-data for the previous 2 data sets explaining the source of each data point.", "https://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL3N.003/2000.02.11/"},
	//{"SRTMGL30", "The 30-arc-second data for the whole world.", "https://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL30.002/2000.02.11/"},
	//{"SRTMSWBD", "The tiled-vector data of the world's coastlines.", "https://e4ftl01.cr.usgs.gov/MEASURES/SRTMSWBD.003/2000.02.11/"},
	//{"SRTMUS1", "The 1-arc-second data for the United States.", "https://e4ftl01.cr.usgs.gov/SRTM/SRTMUS1.003/2000.02.11/"},
	//{"SRTMUS1N", "The meta-data for the United States data set.", "https://e4ftl01.cr.usgs.gov/SRTM/SRTMUS1N.003/2000.02.11/"},
	{"SRTMGL1", "The 1-arc-second data for whole world (NEW!).", "https://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL1.003/2000.02.11/"},
	//{"SRTMGL1N", "The meta-data for this data set.", "https://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL1N.003/2000.02.11/"},
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	f, err := os.Create("geoelevations/urls_generated.go")
	panicIfErr(err)
	f.WriteString("package geoelevations\n")

	for _, d := range data {
		fmt.Println("Generating urls from", d.baseUrl, d.description)
		urls, err := fundURLs(http.DefaultClient, d.baseUrl, map[string]bool{}, 2)
		panicIfErr(err)

		if len(urls) == 0 {
			panic("nor urls for " + d.baseUrl)
		}

		var urlKeys []string
		for k := range urls {
			urlKeys = append(urlKeys, k)
		}
		sort.Strings(urlKeys)

		f.WriteString("\n")
		f.WriteString(`var ` + d.name + ` = SrtmData{
`)
		f.WriteString(`	BaseURL: "` + d.baseUrl + `",
`)
		f.WriteString(`	Files:   map[string]string{
`)
		for _, urlKey := range urlKeys {
			f.WriteString(`"` + urlKey + `": "` + strings.Replace(strings.ReplaceAll(urls[urlKey], "http:", "https:"), d.baseUrl, "", 1) + `",
`)
		}
		f.WriteString(`},
}`)
	}

	panicIfErr(f.Close())
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
	fmt.Println("- parsing", url)
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
