package main

import (
	"errors"
	"strconv"
)

type Taxes struct {
	SGSTList map[string]int
}

func NewTaxes() *Taxes {
	return &Taxes{
		SGSTList: make(map[string]int),
	}
}

func (taxes *Taxes) getSGSTList() error {

	reader, err := getCSVReaderWithoutHeader("input/SGST.csv")
	if isError(err) {
		return err
	}

	for {
		sgst, err := reader.Read()
		if isEOF(err) {
			break
		}
		if isError(err) {
			return err
		}

		sgstState := sgst[0]
		if !matchRegex(sgstState, "^[A-Z]*$") {
			return errors.New("Invalid state code in database -" + sgstState)
		}

		sgstValue, err := strconv.Atoi(sgst[1])
		if isError(err) || !isPositiveInt(sgstValue) {
			return errors.New("Invalid SGST value for the state in database - " + sgst[0])
		}
		taxes.SGSTList[sgstState] = sgstValue
	}
	return err
}
