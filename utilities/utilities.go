package utilities

import (
	"encoding/csv"
	"io"
	"os"
	"regexp"
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
		return nil, err
	}

	r := csv.NewReader(file)
	_, err = r.Read()
	if IsError(err) {
		return nil, err
	}

	return r, nil
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

func IsAlphaNumeric(s string) bool {
	return MatchRegex(s, "^[a-zA-Z0-9 ]*$")
}
