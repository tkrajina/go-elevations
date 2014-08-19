package geoelevations

import (
	"fmt"
	"strings"
)

type SrtmUrl struct {
	File string
	Url  string
}

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
		if strings.HasPrefix(srtmUrl.File, fileName) {
			return &srtmUrl
		}
	}
	return nil
}

func (self *SrtmData) GetSrtm3Url(fileName string) *SrtmUrl {
	for _, srtmUrl := range self.Srtm3 {
		if strings.HasPrefix(srtmUrl.File, fileName) {
			return &srtmUrl
		}
	}
	return nil
}
