package main

import "strconv"

type Taxes struct {
	SGSTList map[string]int
}

func NewTaxes() *Taxes {
	return &Taxes{
		SGSTList: make(map[string]int),
	}
}

func (taxes *Taxes) getSGSTList() error {

	r, err := getCSVReaderWithoutHeader("sgst.csv")
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

		taxes.SGSTList[s[0]], err = strconv.Atoi(s[1])
		if isError(err) {
			return err
		}
	}
	return err
}
