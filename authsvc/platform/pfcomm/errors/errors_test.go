package errors

import (
	"testing"
	"fmt"
	"errors"
)

func TestNewError( t *testing.T ) {
	es := NewError(1,errors.New("Error 1062: Duplicate entry '1' for key 'PRIMARY'"))

	fmt.Printf("code=%d,desc=%s\n",es.Code(),es.Error())
}
