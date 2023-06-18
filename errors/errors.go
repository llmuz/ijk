package errors

import (
	"fmt"
)

type Error struct {
	Status
	cause error
	Data  interface{} `json:"data"`
}

func (c *Error) GetData() interface{} {
	return c.Data
}

func (c *Error) Error() string {
	return fmt.Sprintf("error: status_code = %d err_no = %d err_msg = %s cause = %v",
		c.Status.HttpCode,
		c.Status.ErrNo,
		c.Status.ErrMsg,
		c.cause,
	)
}

func New(httpCode int32, errNo int64, errMsg string, data interface{}) *Error {
	return &Error{
		Status: Status{
			HttpCode: httpCode,
			ErrNo:    errNo,
			ErrMsg:   errMsg,
		},
		cause: nil,
		Data:  data,
	}
}

func Errorf(httpCode int32, errNo int64, errMsg string, err error, data interface{}) error {
	e := New(httpCode, errNo, errMsg, data)
	e.cause = err
	return e
}
