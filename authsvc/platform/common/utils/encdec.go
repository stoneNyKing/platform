package utils

import (
	"errors"
	"github.com/libra9z/mahonia"
)


// All in one ConvertString method, rather than requiring the construction of an iconv.Converter
func ConvertString(input string, fromEncoding string, toEncoding string) (output string, err error) {
	// create a temporary converter

	srcCoder := mahonia.NewDecoder(fromEncoding)

	var sresult string
	var ok bool

	if sresult, ok = srcCoder.ConvertStringOK(input); !ok {
		return "", errors.New("cannot convert string from source Encoding.")
	}

	tagCoder := mahonia.NewEncoder(toEncoding)

	if output, ok = tagCoder.ConvertStringOK(sresult); !ok {
		return "", errors.New("cannot convert string to target Encoding.")
	}

	return output, nil
}
