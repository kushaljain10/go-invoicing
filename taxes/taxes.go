package taxes

import (
	"errors"
	"strconv"

	"github.com/kushaljain/go-invoicing/cache"
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

func GetTaxes(cache *cache.RedisCache) (*Taxes, error) {
	tax := getTax(cache, "SGST")
	if tax != nil {
		return tax, nil
	}

	tax = NewTaxes()
	sgst, err := GetSGSTList(tax)
	if err != nil {
		return nil, err
	}
	tax.SetSGSTList(sgst)

	setTax(cache, "SGST", tax)

	return tax, nil
}

func (taxes *Taxes) SetSGSTList(sgst map[string]int) {
	taxes.SGSTList = sgst
}

func GetSGSTList(taxes *Taxes) (map[string]int, error) {

	reader, err := utilities.GetCSVReader("input/SGST.csv")
	if utilities.IsError(err) {
		return nil, errors.New("error in reading from SGST File")
	}
	headerIndices, err := utilities.GetCSVHeaderIndices(reader)
	if utilities.IsError(err) {
		return nil, errors.New("error while fetching SGST details - " + err.Error())
	}

	for {
		sgst, err := reader.Read()
		if utilities.IsEOF(err) {
			break
		}
		if utilities.IsError(err) {
			return nil, errors.New("error while reading from SGST file")
		}

		sgstState := sgst[headerIndices["StateCode"]]
		if !utilities.MatchRegex(sgstState, "^[A-Z]*$") {
			return nil, errors.New("Invalid state code in database -" + sgstState)
		}

		sgstValue, err := strconv.Atoi(sgst[headerIndices["SGST"]])
		if utilities.IsError(err) || !utilities.IsPositiveInt(sgstValue) {
			return nil, errors.New("Invalid SGST value for the state in database - " + sgst[0])
		}
		taxes.SGSTList[sgstState] = sgstValue
	}
	return taxes.SGSTList, nil
}
