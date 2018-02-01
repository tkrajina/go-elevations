package geoelevations

import (
	"encoding/json"
	"net/http"
	"strings"
)

type SrtmUrl struct {
	// FileName without extension
	Name string `json:"n"`
	Url  string `json:"u"`

	baseUrl string `json:"-"`
}

// Info (to be (se)serialized) about all the SRTM files and their URLs
type SrtmData struct {
	Srtm1BaseUrl string    `json:"srtm1_base_url"`
	Srtm1        []SrtmUrl `json:"srtm1"`
	Srtm3BaseUrl string    `json:"srtm3_base_url"`
	Srtm3        []SrtmUrl `json:"srtm2"`
}

func newSrtmData(client *http.Client, storage SrtmLocalStorage) (*SrtmData, error) {
	fn := "urls.json"

	bytes, err := storage.LoadFile(fn)
	if err != nil {
		if storage.IsNotExists(err) {
			srtmData, err := LoadSrtmData(client)
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
		} else {
			return nil, err
		}
	}

	srtmData := new(SrtmData)
	if err := json.Unmarshal(bytes, srtmData); err != nil {
		return nil, err
	}

	return srtmData, nil
}

func (self *SrtmData) GetBestSrtmUrl(fileName string) (string, *SrtmUrl) {
	baseUrl, srtm3Url := self.GetSrtm3Url(fileName)
	if srtm3Url != nil {
		return baseUrl, srtm3Url
	}

	return self.GetSrtm1Url(fileName)
}

func (self *SrtmData) GetSrtm1Url(fileName string) (string, *SrtmUrl) {
	for _, srtmUrl := range self.Srtm1 {
		if strings.HasPrefix(fileName, srtmUrl.Name) {
			return self.Srtm1BaseUrl, &srtmUrl
		}
	}
	return "", nil
}

func (self *SrtmData) GetSrtm3Url(fileName string) (string, *SrtmUrl) {
	for _, srtmUrl := range self.Srtm3 {
		if strings.HasPrefix(srtmUrl.Name, fileName) {
			return self.Srtm3BaseUrl, &srtmUrl
		}
	}
	return "", nil
}
