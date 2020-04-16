package geoelevations

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"log"

	"io/ioutil"
)

func gzipBytes(b *[]byte) (*[]byte, error) {
	buf := new(bytes.Buffer)
	gz := gzip.NewWriter(buf)

	_, err := gz.Write(*b)
	if err != nil {
		return nil, err
	}

	err = gz.Close()
	if err != nil {
		return nil, err
	}

	out := buf.Bytes()

	return &out, nil
}

func ungzipBytes(b *[]byte) (*[]byte, error) {
	r, err := gzip.NewReader(ioutil.NopCloser(bytes.NewBuffer(*b)))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	bb, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &bb, nil
}

func unzipBytes(byts []byte) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(byts), int64(len(byts)))
	if err != nil {
		return nil, err
	}

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			log.Printf("Error reading %s: %s", f.Name, err.Error())
			return nil, err
		}
		defer rc.Close()

		bytes, err := ioutil.ReadAll(rc)
		if err != nil {
			log.Printf("Error reading %s: %s", f.Name, err.Error())
			return nil, err
		}

		return bytes, nil
	}

	return nil, errors.New(fmt.Sprintf("No file in .zip"))
}
