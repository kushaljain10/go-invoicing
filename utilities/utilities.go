package utilities

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"regexp"
)

var (
	ALPHANUMERIC_REGEX = "^[a-zA-Z0-9 ]*$"
)

func IsError(err error) bool {
	return err != nil
}

func IsEOF(err error) bool {
	return err == io.EOF
}

func GetCSVReaderWithoutHeader(filename string) (*csv.Reader, error) {
	file, err := os.Open(filename)
	if IsError(err) {
		return nil, errors.New("error in opening file - " + filename)
	}

	r := csv.NewReader(file)
	_, err = r.Read()
	if IsError(err) {
		return nil, errors.New("error in reading header from CSV")
	}

	return r, nil
}

func GetCSVReader(filename string) (*csv.Reader, error) {
	file, err := os.Open(filename)
	if IsError(err) {
		return nil, errors.New("error in opening file - " + filename)
	}

	r := csv.NewReader(file)

	return r, nil
}

func GetCSVHeaderIndices(reader *csv.Reader) (map[string]int, error) {

	indices := make(map[string]int)

	header, err := reader.Read()
	if IsError(err) {
		return nil, errors.New("error in reading header from CSV")
	}

	for i, h := range header {
		indices[h] = i
	}

	return indices, nil
}

func CheckErrors(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}

func IsPositiveInt(num int) bool {
	return num > 0
}

func IsPositiveFloat(num float64) bool {
	return num > 0
}

func IsEmptyString(s string) bool {
	return len(s) == 0
}

func MatchRegex(s string, expression string) bool {
	re := regexp.MustCompile(expression)
	return re.MatchString(s)
}

func ValidateString(s string, regex string) (string, error) {
	if IsEmptyString(s) {
		return s, errors.New("no value found")
	}
	if !MatchRegex(s, regex) {
		return s, errors.New("invalid value")
	}
	return s, nil
}
