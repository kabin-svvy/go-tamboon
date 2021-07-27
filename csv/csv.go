package csv

import (
	"bytes"
	stdcsv "encoding/csv"
	"errors"
	"go-tamboon/cipher"
	"io/ioutil"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Transaction struct {
	Name           string
	AmountSubunits string
	CCNumber       string
	CVV            string
	ExpMonth       string
	ExpYear        string
	ChargeToken    string
	Ccy            string
	ChargeStatus   string
	ChargeID       string
}

func Read(filename string) ([]Transaction, error) {
	bufferFromFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	reader, err := cipher.NewRot128Reader(bytes.NewBuffer(bufferFromFile))
	if err != nil {
		return nil, err
	}

	size := len(bufferFromFile)

	bufferFromDecrypt := make([]byte, size, size)

	n, err := reader.Read(bufferFromDecrypt)
	if err != nil {
		return nil, err
	}

	if n != size {
		return nil, errors.New("plaintext and ciphertext differ in sizes")
	}

	raw, err := stdcsv.NewReader(bytes.NewBuffer(bufferFromDecrypt)).ReadAll()
	if err != nil {
		return nil, err
	}

	result := []Transaction{}

	for k, v := range raw {
		if k == 0 {
			continue
		}
		t := transform(v)
		result = append(result, *t)
	}
	return result, nil
}

func transform(v []string) *Transaction {
	if len(v) < 6 {
		logrus.Error(v, " string array length less than 6")
		return nil
	}

	name := v[0]
	amout := v[1]
	ccNumber := v[2]
	cvv := v[3]
	expMonth := v[4]
	expYear := v[5]

	if name == "" {
		logrus.Error(v, " Name should not be empty")
		return nil
	}

	amountSubUnits, err := strconv.ParseFloat(amout, 64)
	if err != nil {
		logrus.Error(v, err)
		return nil
	}

	if amountSubUnits <= 0 {
		logrus.Error(v, " AmountSubUnit less than 0")
		return nil
	}

	if ccNumber == "" {
		logrus.Error(v, " CCNumber should not be empty")
		return nil
	}

	if cvv == "" {
		logrus.Error(v, " CVV should not be empty")
		return nil
	}

	if expMonth == "" {
		logrus.Error(v, " ExpMonth should not be empty")
		return nil
	}

	if expYear == "" {
		logrus.Error(v, " ExpYear should not be empty")
		return nil
	}

	return &Transaction{
		Name:           name,
		AmountSubunits: amout,
		CCNumber:       ccNumber,
		CVV:            cvv,
		ExpMonth:       expMonth,
		ExpYear:        expYear,
		Ccy:            "THB",
	}
}
