package geoelevations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func reloadJsonUrls(destinationFilename string) error {
	srtmData, err := LoadSrtmData()
	if err != nil {
		return err
	}

	srtmDataJson, err := json.MarshalIndent(srtmData, "", "\t")
	if err != nil {
		return err
	}

	f, err := os.Create(destinationFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(srtmDataJson)

	log.Print("Written ", len(srtmDataJson), " bytes to ", destinationFilename)

	return nil
}

type SrtmUrl struct {
	// FileName without extension
	Name string
	Url  string
}

// Info (to be (se)serialized) about all the SRTM files and their URLs
type SrtmData struct {
	Srtm1 []SrtmUrl
	Srtm3 []SrtmUrl
}

func newSrtmData(cacheDirectory string) *SrtmData {
	urlsFilename := path.Join(cacheDirectory, "urls.json")
	f, err := os.Open(urlsFilename)
	if err != nil {
		if os.IsNotExist(err) {
			reloadJsonUrls(urlsFilename)
		} else {
			panic(fmt.Sprintf("Can't find srtm urls in \"%s\"", cacheDirectory))
		}
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("Can't load srtm urls in \"%s\"", cacheDirectory))
	}

	srtmData := new(SrtmData)
	json.Unmarshal(bytes, srtmData)

	return srtmData
}

func (self *SrtmData) GetBestSrtmUrl(fileName string) *SrtmUrl {
	srtm3Url := self.GetSrtm3Url(fileName)
	if srtm3Url != nil {
		return srtm3Url
	}

	return self.GetSrtm1Url(fileName)
}

func (self *SrtmData) GetSrtm1Url(fileName string) *SrtmUrl {
	for _, srtmUrl := range self.Srtm1 {
		if strings.HasPrefix(fileName, srtmUrl.Name) {
			return &srtmUrl
		}
	}
	return nil
}

func (self *SrtmData) GetSrtm3Url(fileName string) *SrtmUrl {
	for _, srtmUrl := range self.Srtm3 {
		if strings.HasPrefix(srtmUrl.Name, fileName) {
			return &srtmUrl
		}
	}
	return nil
}
