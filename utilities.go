package main

import (
	"encoding/csv"
	"io"
	"os"
)

func isError(err error) bool {
	if err != nil {
		return true
	}
	return false
}

func isEOF(err error) bool {
	if err == io.EOF {
		return true
	}
	return false
}

func getCSVReaderWithoutHeader(filename string) (*csv.Reader, error) {
	file, err := os.Open(filename)
	if isError(err) {
		return nil, err
	}

	r := csv.NewReader(file)
	_, err = r.Read()
	if isError(err) {
		return nil, err
	}

	return r, nil
}
