package main

import (
	"encoding/csv"
	"io"
	"os"
)

func isError(err error) bool {
	return err != nil
}

func isEOF(err error) bool {
	return err == io.EOF
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
