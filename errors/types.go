package errors

var (
	SuccessErrNo      int64 = 10_000 // 请求成功
	BadRequestErrNo   int64 = 400    // 参数错误
	UnauthorizedErrNo int64 = 401    //  没有授权
	UnknownErrNo      int64 = 20_000 // 未知错误
)

// BadRequest new BadRequest error that is mapped to a 400 response.
func BadRequest(errMsg string) *Error {
	return New(400, BadRequestErrNo, errMsg, nil)
}

// Unauthorized new Unauthorized error that is mapped to a 401 response.
func Unauthorized(errMsg string) *Error {
	return New(401, UnauthorizedErrNo, errMsg, nil)
}

func Success(data interface{}) *Error {
	return New(200, SuccessErrNo, "success", data)
}
