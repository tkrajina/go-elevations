package geoelevations

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"strings"
)

const (
	SRTM_BASE_URL = "https://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL3.003/2000.02.11/"
)

type Srtm struct {
	cache map[string]*SrtmFile

	client             *httpClient
	username, password string

	srtmData SrtmData
	storage  SrtmLocalStorage
}

func NewSrtm() (*Srtm, error) {
	return NewSrtmWithCustomCacheDir("")
}

func NewSrtmWithCustomStorage(storage SrtmLocalStorage) (*Srtm, error) {
	return &Srtm{
		cache:    make(map[string]*SrtmFile),
		storage:  storage,
		srtmData: *&SRTMGL3S,
	}, nil
}

func NewSrtmWithCustomCacheDir(cacheDirectory string) (*Srtm, error) {
	storage, err := NewLocalFileSrtmStorage(cacheDirectory)
	if err != nil {
		return nil, err
	}
	return NewSrtmWithCustomStorage(storage)
}

func (s *Srtm) SetAuth(username, password string) {
	s.username = username
	s.password = password
}

func (self *Srtm) GetElevation(latitude, longitude float64) (float64, error) {
	srtmFileName, srtmLatitude, srtmLongitude := self.getSrtmFileNameAndCoordinates(latitude, longitude)
	//log.Printf("srtmFileName for %v,%v: %s", latitude, longitude, srtmFileName)

	if self.client == nil {
		c, err := newHTTPClient(self.username, self.password)
		if err != nil {
			return 0, err
		}
		self.client = c
	}

	_, found := self.cache[srtmFileName]
	if !found {
		srtmURL, found := SRTMGL3S.getFileURL(srtmFileName)
		if !found {
			return 0, fmt.Errorf("no SRTM url for (%f,%f) (%s)", latitude, longitude, srtmFileName)
		}
		self.cache[srtmFileName] = newSrtmFile(srtmFileName, srtmURL, srtmLatitude, srtmLongitude)
	}

	return self.cache[srtmFileName].getElevation(self.client, self.storage, latitude, longitude)
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
}

func newSrtmFile(name, fileUrl string, latitude, longitude float64) *SrtmFile {
	result := SrtmFile{}
	result.name = name
	result.isValidSrtmFile = len(fileUrl) > 0
	result.latitude = latitude
	result.longitude = longitude

	result.fileUrl = fileUrl
	if !strings.HasSuffix(result.fileUrl, ".zip") {
		result.fileUrl += ".zip"
	}

	return &result
}

func (self *SrtmFile) loadContents(client *httpClient, storage SrtmLocalStorage) error {
	if !self.isValidSrtmFile || len(self.fileUrl) == 0 {
		return nil
	}

	fileName := fmt.Sprintf("%s.hgt.zip", self.name)

	bytes, err := storage.LoadFile(fileName)
	if err != nil {
		if storage.IsNotExists(err) {
			log.Printf("File %s not retrieved => retrieving: %s", fileName, self.fileUrl)
			responseBytes, err := client.downloadSrtmURL(self.fileUrl)
			if err != nil {
				return err
			}

			if err := storage.SaveFile(fileName, responseBytes); err != nil {
				return err
			}
			log.Printf("Written %d bytes to %s", len(responseBytes), fileName)

			bytes = responseBytes
		} else {
			return err
		}
	}

	contents, err := unzipBytes(bytes)
	if err != nil {
		log.Printf("Error loading file %s: %s", fileName, err.Error())
	}
	self.contents = contents

	log.Printf("Loaded %dbytes from %s, squareSize=%d", len(self.contents), fileName, self.squareSize)

	return nil
}

func (self *SrtmFile) getElevation(client *httpClient, storage SrtmLocalStorage, latitude, longitude float64) (float64, error) {
	if !self.isValidSrtmFile || len(self.fileUrl) == 0 {
		log.Printf("Invalid file %s", self.name)
		return math.NaN(), nil
	}

	if len(self.contents) == 0 {
		log.Println("load contents")
		err := self.loadContents(client, storage)
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
