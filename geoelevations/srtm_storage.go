package geoelevations

import (
	"log"
	"os"
	"path"
)

type SrtmLocalStorage interface {
	// LoadFile loads a file, if not available, then err!=nil and IsNotExists(err) must be true
	LoadFile(fn string) ([]byte, error)
	IsNotExists(err error) bool
	SaveFile(fn string, bytes []byte) error
}

type LocalFileSrtmStorage struct {
	cacheDirectory string
}

func NewLocalFileSrtmStorage(cacheDirectory string) (*LocalFileSrtmStorage, error) {
	if len(cacheDirectory) == 0 {
		cacheDirectory = path.Join(os.Getenv("HOME"), ".geoelevations")
	}
	log.Printf("Using %s to cache SRTM files", cacheDirectory)

	if _, err := os.Stat(cacheDirectory); os.IsNotExist(err) {
		log.Print("Creating", cacheDirectory)

		if err := os.Mkdir(cacheDirectory, os.ModeDir|0700); err != nil {
			return nil, err
		}
	}

	return &LocalFileSrtmStorage{cacheDirectory: cacheDirectory}, nil
}
func (ds LocalFileSrtmStorage) LoadFile(fn string) ([]byte, error) {
	return nil, nil
}
func (ds LocalFileSrtmStorage) IsNotExists(err error) bool {
	return false
}
func (ds LocalFileSrtmStorage) SaveFile(fn string, bytes []byte) error {
	return nil
}

var _ SrtmLocalStorage = new(LocalFileSrtmStorage)
