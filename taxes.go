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

	r, err := getCSVReaderWithoutHeader("SGST.csv")
	if isError(err) {
		return err
	}

	for {
		s, err := r.Read()
		if isEOF(err) {
			break
		}
		if isError(err) {
			return err
		}

		sgst, err := strconv.Atoi(s[1])
		if isError(err) {
			return err
		}
		if !isPositiveInt(sgst) {
			return errors.New("Invalid SGST for the state - " + s[0])
		}
		taxes.SGSTList[s[0]] = sgst
	}
	return err
}
