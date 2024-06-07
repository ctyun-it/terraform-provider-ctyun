package ctyunsdk

const (
	OccurOccasionBeforeRequest  = 0 // 请求前
	OccurOccasionBeforeResponse = 1 // 响应前
	OccurOccasionAfterResponse  = 2 // 响应后
)

type OccurOccasion int

type CtyunRequestError interface {
	error                          // 错误信息
	OccurOccasion() OccurOccasion  // 发生场合
	ErrorCode() string             // 返回码
	CtyunResponse() *CtyunResponse // 标准返回
}

type ctyunRequestError struct {
	errorMessage  string
	occurOccasion OccurOccasion
	ctyunResponse *CtyunResponse
	errorCode     string
}

func (c ctyunRequestError) Error() string {
	return c.errorMessage
}

func (c ctyunRequestError) OccurOccasion() OccurOccasion {
	return c.occurOccasion
}

func (c ctyunRequestError) ErrorCode() string {
	return c.errorCode
}

func (c ctyunRequestError) CtyunResponse() *CtyunResponse {
	return c.ctyunResponse
}

// WrapError 包裹异常
func WrapError(err error, resp *CtyunResponse) CtyunRequestError {
	occurOccasion := 0
	if resp.Request != nil {
		occurOccasion++
	}
	if resp.Response != nil {
		occurOccasion++
	}
	return ctyunRequestError{
		errorMessage:  err.Error(),
		occurOccasion: OccurOccasion(occurOccasion),
		ctyunResponse: resp,
	}
}

// ErrorBeforeRequest 在请求前发生的异常
func ErrorBeforeRequest(err error) CtyunRequestError {
	return ctyunRequestError{
		errorMessage:  err.Error(),
		occurOccasion: OccurOccasionBeforeRequest,
	}
}

// ErrorBeforeResponse 在响应前发生的异常
func ErrorBeforeResponse(err error) CtyunRequestError {
	return ctyunRequestError{
		errorMessage:  err.Error(),
		occurOccasion: OccurOccasionBeforeResponse,
	}
}

// ErrorAfterResponse 在响应后发生的异常
func ErrorAfterResponse(err error, resp *CtyunResponse) CtyunRequestError {
	return WrapError(err, resp)
}

// WrapWithErrorCode 包裹异常
func WrapWithErrorCode(err error, errorCode string, resp *CtyunResponse) CtyunRequestError {
	wrapError := WrapError(err, resp)
	requestError := interface{}(wrapError).(ctyunRequestError)
	requestError.errorCode = errorCode
	return requestError
}
