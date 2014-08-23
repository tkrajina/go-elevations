package geoelevations

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
)

const (
	SRTM_BASE_URL = "http://dds.cr.usgs.gov/srtm"
	SRTM1_URL     = "/version2_1/SRTM1/"
	SRTM3_URL     = "/version2_1/SRTM3/"
)

type Srtm struct {
}

func NewSrtm() *Srtm {
	return new(Srtm)
}

func (self *Srtm) GetElevation(latitude, longitude float64) (float64, error) {
	srtmFileName := self.getSrtmFileName(latitude, longitude)
	log.Printf("srtmFileName for %v,%v: %s", latitude, longitude, srtmFileName)

	srtmData := GetSrtmData()

	// TODO Cache files...
	srtmFile := newSrtmFile(srtmFileName, "")
	srtmFileUrl := srtmData.GetBestSrtmUrl(srtmFileName)
	if srtmFileUrl != nil {
		srtmFile = newSrtmFile(srtmFileName, srtmFileUrl.Url)
	}

	return srtmFile.getElevation(latitude, longitude)
}

func (self *Srtm) getSrtmFileName(latitude, longitude float64) string {
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

// Struct with contents and some utility methods of a single SRTM file
type SrtmFile struct {
	contents        []byte
	fileName        string
	fileUrl         string
	isValidSrtmFile bool
	fileRetrieved   bool
}

func newSrtmFile(fileName, fileUrl string) *SrtmFile {
	result := SrtmFile{}
	result.fileName = fileName
	result.isValidSrtmFile = len(fileUrl) > 0

	result.fileUrl = fileUrl
	if !strings.HasSuffix(result.fileUrl, ".zip") {
		result.fileUrl += ".zip"
	}

	return &result
}

func (self SrtmFile) loadContents() error {
	if !self.isValidSrtmFile || len(self.fileUrl) == 0 {
		return nil
	}

	// Retrieve if needed:
	if _, err := os.Stat(self.fileName); os.IsNotExist(err) {
		log.Printf("Retrieving: %s", self.fileUrl)
		response, err := http.Get(self.fileUrl)
		if err != nil {
			log.Printf("Error retrieving file: %s", err.Error())
			return err
		}

		responseBytes, _ := ioutil.ReadAll(response.Body)

		f, err := os.Create(self.fileName)
		if err != nil {
			log.Printf("Error writing file %s: %s", self.fileName, err.Error())
			return err
		}
		defer f.Close()

		f.Write(responseBytes)
		log.Printf("Written %d bytes to %s", len(responseBytes), self.fileName)
	}

	f, err := os.Open(self.fileName)
	if err != nil {
		log.Printf("Error loading file %s: %s", self.fileName, err.Error())
	}
	defer f.Close()

	self.contents, err = ioutil.ReadAll(f)
	if err != nil {
		log.Printf("Error loading file %s: %s", self.fileName, err.Error())
	}

	return nil
}

func (self SrtmFile) getElevation(latitude, longitude float64) (float64, error) {
	if !self.isValidSrtmFile || len(self.fileUrl) == 0 {
		return math.NaN(), nil
	}

	if len(self.contents) == 0 {
		err := self.loadContents()
		if err != nil {
			return math.NaN(), err
		}
	}

	return 0.0, nil
}

// ----------------------------------------------------------------------------------------------------
// Misc util functions
// ----------------------------------------------------------------------------------------------------

func GetSrtmData() *SrtmData {
	f, err := os.Open("urls.json")
	if err != nil {
		panic("Can't find srtm urls")
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic("Can't find srtm urls")
	}

	srtmData := new(SrtmData)
	json.Unmarshal(bytes, srtmData)

	return srtmData
}

func LoadSrtmData() (*SrtmData, error) {
	result := new(SrtmData)

	var err error
	result.Srtm1, err = getLinksFromUrl(SRTM_BASE_URL+SRTM1_URL, 0)
	if err != nil {
		return nil, err
	}

	result.Srtm3, err = getLinksFromUrl(SRTM_BASE_URL+SRTM3_URL, 0)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getLinksFromUrl(url string, depth int) ([]SrtmUrl, error) {
	if depth >= 2 {
		return []SrtmUrl{}, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	result := make([]SrtmUrl, 0)

	urls := getLinksFromHtmlDocument(resp.Body)
	for _, tmpUrl := range urls {
		urlLowercase := strings.ToLower(tmpUrl)
		if strings.HasSuffix(urlLowercase, ".hgt.zip") {
			parts := strings.Split(tmpUrl, "/")
			name := parts[len(parts)-1]
			name = strings.Replace(name, ".hgt.zip", "", -1)
			srtmUrl := SrtmUrl{Name: name, Url: fmt.Sprintf("%s/%s", url, tmpUrl)}
			result = append(result, srtmUrl)
			log.Printf("> %s/%s -> %s\n", url, tmpUrl, tmpUrl)
		} else if len(urlLowercase) > 0 && urlLowercase[0] != '/' && !strings.HasPrefix(urlLowercase, "http") && !strings.HasSuffix(urlLowercase, ".jpg") {
			newLinks, err := getLinksFromUrl(fmt.Sprintf("%s/%s", url, tmpUrl), depth+1)
			if err != nil {
				return nil, err
			}
			result = append(result, newLinks...)
			log.Printf("> %s\n", tmpUrl)
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
