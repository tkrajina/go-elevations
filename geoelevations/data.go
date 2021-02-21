package geoelevations

import (
	"encoding/json"
	"net/http"
)

type SrtmUrl struct {
	// FileName without extension
	Name string `json:"n"`
	Url  string `json:"u"`
}

// Info (to be (se)serialized) about all the SRTM files and their URLs
type SrtmData struct {
	BaseUrl string            `json:"baseUrl"`
	Files   map[string]string `json:"files"`
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
	return self.GetSrtm3Url(fileName)
}

func (self *SrtmData) GetSrtm3Url(fileName string) (string, *SrtmUrl) {
	return "", nil
	/*
		for _, srtmUrl := range self.Srtm3 {
			if strings.HasPrefix(srtmUrl.Name, fileName) {
				return self.SRTMBaseUrl, &srtmUrl
			}
		}
		return "", nil
	*/
}
