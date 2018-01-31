package geoelevations

import (
	"encoding/json"
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

func newSrtmData(storage SrtmLocalStorage) (*SrtmData, error) {
	fn := "urls.json"

	bytes, err := storage.LoadFile(fn)
	if storage.IsNotExists(err) {
		srtmData, err := LoadSrtmData()
		if err != nil {
			return nil, err
		}
		b, err := json.Marshal(srtmData)
		if err != nil {
			return nil, err
		}
		bytes = b

		if err := storage.SaveFile(fn, bytes); err != nil {
			return nil, err
		}
	}

	srtmData := new(SrtmData)
	if err := json.Unmarshal(bytes, srtmData); err != nil {
		return nil, err
	}

	return srtmData, nil
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
