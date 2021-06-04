package main

import (
	"encoding/csv"
	"io"
	"os"
	"regexp"
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

func isPositiveInt(num int) bool {
	return num > 0
}

func isPositiveFloat(num float64) bool {
	return num > 0
}

func isEmptyString(s string) bool {
	return len(s) == 0
}

func matchRegex(s string, expression string) bool {
	re := regexp.MustCompile(expression)
	return re.MatchString(s)
}

func isAlphaNumeric(s string) bool {
	return matchRegex(s, "^[a-zA-Z0-9 ]*$")
}
