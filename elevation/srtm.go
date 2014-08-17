package geoelevations

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
)

const (
	SRTM_BASE_URL = "http://dds.cr.usgs.gov/srtm"
	SRTM1_URL     = "/version2_1/SRTM1/"
	SRTM3_URL     = "/version2_1/SRTM3/"
)

type Srtm struct {
}

func (self *Srtm) GetElevation(latitude, longitude float64) float64 {
	return 0
}

func getSrtmFileName(latitude, longitude float64) string {
	northSouth := 'S'
	if latitude >= 0 {
		northSouth = 'N'
	}

	eastWest := 'W'
	if longitude >= 0 {
		eastWest = 'E'
	}

	latPart := int(math.Abs(math.Floor(latitude)))
	lonPart := int(math.Abs(math.Floor(longitude)))

	return fmt.Sprintf("%s%02d%s%03d.hgt", string(northSouth), latPart, string(eastWest), lonPart)
}

func getLinksFromUrl(url string, depth int) ([]string, error) {

	if depth >= 2 {
		return make([]string, 0), nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	result := make([]string, 10)

	urls := getLinksFromHtmlDocument(resp.Body)
	for _, tmpUrl := range urls {
		urlLowercase := strings.ToLower(tmpUrl)
		if strings.HasSuffix(urlLowercase, ".hgt.zip") {
			fmt.Printf("hgt>%s/%s -> %s...\n", url, tmpUrl, tmpUrl)
		} else if len(urlLowercase) > 0 && urlLowercase[0] != '/' && !strings.HasPrefix(urlLowercase, "http") && !strings.HasSuffix(urlLowercase, ".jpg") {
			getLinksFromUrl(fmt.Sprintf("%s/%s", url, tmpUrl), depth+1)
			fmt.Printf(">%s...\n", tmpUrl)
		}
	}

	return result, nil
}

func getLinksFromHtmlDocument(html io.ReadCloser) []string {
	result := make([]string, 10)

	decoder := xml.NewDecoder(html)
	for token, _ := decoder.Token(); token != nil; token, _ = decoder.Token() {
		switch typedToken := token.(type) {
		case xml.StartElement:
			for _, attr := range typedToken.Attr {
				if strings.ToLower(attr.Name.Local) == "href" {
					result = append(result, strings.Trim(attr.Value, " \r\t\n"))
				}
			}
		}
	}

	return result
}
