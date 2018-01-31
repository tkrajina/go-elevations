package geoelevations

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	SRTM_BASE_URL = "http://dds.cr.usgs.gov/srtm"
	SRTM1_URL     = "/version2_1/SRTM1/"
	SRTM3_URL     = "/version2_1/SRTM3/"
)

type Srtm struct {
	cacheDirectory string
	cache          map[string]*SrtmFile
	srtmData       SrtmData
}

func NewSrtm() (*Srtm, error) {
	return NewSrtmWithCustomCacheDir("")
}

func NewSrtmWithCustomCacheDir(cacheDirectory string) (*Srtm, error) {
	if len(cacheDirectory) == 0 {
		// TODO: Windows
		cacheDirectory = path.Join(os.Getenv("HOME"), ".geoelevations")
	}
	log.Printf("Using %s to cache SRTM files", cacheDirectory)

	if _, err := os.Stat(cacheDirectory); os.IsNotExist(err) {
		log.Print("Creating", cacheDirectory)

		if err := os.Mkdir(cacheDirectory, os.ModeDir|0700); err != nil {
			return nil, err
		}
	}

	result := new(Srtm)
	result.cache = make(map[string]*SrtmFile)
	result.cacheDirectory = cacheDirectory
	result.srtmData = *newSrtmData(cacheDirectory)
	return result, nil
}

func (self *Srtm) GetElevation(client *http.Client, latitude, longitude float64) (float64, error) {
	srtmFileName, srtmLatitude, srtmLongitude := self.getSrtmFileNameAndCoordinates(latitude, longitude)
	//log.Printf("srtmFileName for %v,%v: %s", latitude, longitude, srtmFileName)

	srtmFile, ok := self.cache[srtmFileName]
	if !ok {
		srtmFile = newSrtmFile(srtmFileName, "", self.cacheDirectory, srtmLatitude, srtmLongitude)
		srtmFileUrl := self.srtmData.GetBestSrtmUrl(srtmFileName)
		if srtmFileUrl != nil {
			srtmFile = newSrtmFile(srtmFileName, srtmFileUrl.Url, self.cacheDirectory, srtmLatitude, srtmLongitude)
		}
		self.cache[srtmFileName] = srtmFile
	}

	return srtmFile.getElevation(client, latitude, longitude)
}

func (self *Srtm) getSrtmFileNameAndCoordinates(latitude, longitude float64) (string, float64, float64) {
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

	srtmFileName := fmt.Sprintf("%s%02d%s%03d", string(northSouth), latPart, string(eastWest), lonPart)

	return srtmFileName, math.Floor(latitude), math.Floor(longitude)
}

// Struct with contents and some utility methods of a single SRTM file
type SrtmFile struct {
	latitude, longitude float64
	contents            []byte
	name                string
	fileUrl             string
	isValidSrtmFile     bool
	fileRetrieved       bool
	squareSize          int
	cacheDirectory      string
}

func newSrtmFile(name, fileUrl, cacheDirectory string, latitude, longitude float64) *SrtmFile {
	result := SrtmFile{}
	result.name = name
	result.isValidSrtmFile = len(fileUrl) > 0
	result.latitude = latitude
	result.longitude = longitude
	result.cacheDirectory = cacheDirectory

	result.fileUrl = fileUrl
	if !strings.HasSuffix(result.fileUrl, ".zip") {
		result.fileUrl += ".zip"
	}

	return &result
}

func (self *SrtmFile) loadContents(client *http.Client) error {
	if !self.isValidSrtmFile || len(self.fileUrl) == 0 {
		return nil
	}

	fileName := path.Join(self.cacheDirectory, fmt.Sprintf("%s.hgt.zip", self.name))

	// Retrieve if needed:
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Printf("File %s not retrieved => retrieving: %s", fileName, self.fileUrl)
		req, err := http.NewRequest(http.MethodGet, self.fileUrl, nil)
		if err != nil {
			return err
		}
		response, err := client.Do(req)
		if err != nil {
			log.Printf("Error retrieving file: %s", err.Error())
			return err
		}

		responseBytes, _ := ioutil.ReadAll(response.Body)

		f, err := os.Create(fileName)
		if err != nil {
			log.Printf("Error writing file %s: %s", fileName, err.Error())
			return err
		}
		defer f.Close()

		f.Write(responseBytes)
		log.Printf("Written %d bytes to %s", len(responseBytes), fileName)
	}

	contents, err := unzipFile(fileName)
	if err != nil {
		log.Printf("Error loading file %s: %s", fileName, err.Error())
	}
	self.contents = contents

	log.Printf("Loaded %dbytes from %s, squareSize=%d", len(self.contents), fileName, self.squareSize)

	return nil
}

func (self *SrtmFile) getElevation(client *http.Client, latitude, longitude float64) (float64, error) {
	if !self.isValidSrtmFile || len(self.fileUrl) == 0 {
		log.Printf("Invalid file %s", self.name)
		return math.NaN(), nil
	}

	if len(self.contents) == 0 {
		log.Println("load contents")
		err := self.loadContents(client)
		if err != nil {
			return math.NaN(), err
		}
	}

	if self.squareSize <= 0 {
		squareSizeFloat := math.Sqrt(float64(len(self.contents)) / 2.0)
		self.squareSize = int(squareSizeFloat)

		if squareSizeFloat != float64(self.squareSize) || self.squareSize <= 0 {
			return math.NaN(), errors.New(fmt.Sprintf("Invalid size for file %s: %d", self.name, len(self.contents)))
		}
	}

	row, column := self.getRowAndColumn(latitude, longitude)
	//log.Printf("(%f, %f) => row, column = %d, %d", latitude, longitude, row, column)
	elevation := self.getElevationFromRowAndColumn(row, column)

	return elevation, nil
}

func (self SrtmFile) getElevationFromRowAndColumn(row, column int) float64 {
	i := row*self.squareSize + column
	byte1 := self.contents[i*2]
	byte2 := self.contents[i*2+1]
	result := int(byte1)*256 + int(byte2)

	if result > 9000 {
		return math.NaN()
	}

	return float64(result)
	/*
	   i = row * (@square_side) + column

	   i < @square_side ** 2 or raise "Invalid i=#{i}"

	   @file.seek(i * 2)
	   bytes = @file.read(2)
	   byte_1 = bytes[0].ord
	   byte_2 = bytes[1].ord

	   result = byte_1 * 256 + byte_2

	   if result > 9000
	       # TODO(TK) try to detect the elevation from neighbour point:
	       return nil
	   end

	   result
	*/
}

func (self SrtmFile) getRowAndColumn(latitude, longitude float64) (int, int) {
	row := int((self.latitude + 1.0 - latitude) * (float64(self.squareSize - 1.0)))
	column := int((longitude - self.longitude) * (float64(self.squareSize - 1.0)))
	//log.Printf("squareSize=%v", self.squareSize)
	//log.Printf("row, column = %v, %v", row, column)
	return row, column
}

// ----------------------------------------------------------------------------------------------------
// Misc util functions
// ----------------------------------------------------------------------------------------------------

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
