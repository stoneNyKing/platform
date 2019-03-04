package errors

import (
	"strings"
	"platform/common/utils"
)

type Error	struct {
	ErrType 		int
	ErrCode 		int
	Description 	string
}


type Errors interface {
	Error()			string
	Code()			int
	ErrorType() 	int
}


func CommonError(ei int) string{
	if ret,ok := Comm_Error[ei];ok {
		return ret
	}else{
		return "未知错误"
	}
}


func NewError(etype int,err error) *Error {

	e := new(Error)
	e.ErrType = etype
	switch etype {
	case ERR_TYPE_DBERR :
		e.ErrCode,e.Description = parseDbError(err)
	case ERR_TYPE_COMMON :
		e.ErrCode,e.Description = parseCommonError(err)
	}

	return e
}


func (e *Error)Error() string {
	return e.Description
}

func (e *Error)Code() int {
	return e.ErrCode
}

func (e *Error)ErrorType() int {
	return e.ErrType
}

func parseDbError(err error)(code int,desc string){
	if err == nil {
		return 0,""
	}

	es := err.Error()

	ee :=strings.Split(es,":")
	if len(ee)>0 {
		s1 := ee[0]
		s2 := strings.Split(s1," ")

		if len(s2)>1 {
			code = utils.Convert2Int(s2[1])

			if code == 0 {
				code = 9999
			}

			desc = DB_Error[code]
		}
	}

	return
}

func parseCommonError(err error)(code int,desc string){
	if err == nil {
		return 0,""
	}

	es := err.Error()

	ei := utils.Convert2Int(es)

	if des ,ok := Comm_Error[ei];ok {
		desc = des
		code = ei
	}else{
		desc = "未知错误"
		code = 99999
	}

	return
}