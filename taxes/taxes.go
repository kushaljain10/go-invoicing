package taxes

import (
	"errors"
	"strconv"

	"github.com/kushaljain/go-invoicing/utilities"
)

type Taxes struct {
	SGSTList map[string]int
}

func NewTaxes() *Taxes {
	return &Taxes{
		SGSTList: make(map[string]int),
	}
}

func (taxes *Taxes) GetSGSTList() error {

	reader, err := utilities.GetCSVReaderWithoutHeader("input/SGST.csv")
	if utilities.IsError(err) {
		return errors.New("error in reading from SGST File")
	}

	for {
		sgst, err := reader.Read()
		if utilities.IsEOF(err) {
			break
		}
		if utilities.IsError(err) {
			return errors.New("error while reading from SGST file")
		}

		sgstState := sgst[0]
		if !utilities.MatchRegex(sgstState, "^[A-Z]*$") {
			return errors.New("Invalid state code in database -" + sgstState)
		}

		sgstValue, err := strconv.Atoi(sgst[1])
		if utilities.IsError(err) || !utilities.IsPositiveInt(sgstValue) {
			return errors.New("Invalid SGST value for the state in database - " + sgst[0])
		}
		taxes.SGSTList[sgstState] = sgstValue
	}
	return nil
}
