package geoelevations

import (
	//"fmt"
	"strings"
)

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
