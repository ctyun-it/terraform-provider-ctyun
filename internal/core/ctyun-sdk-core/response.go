package ctyunsdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	StatusCodeSuccessString = "800"
)

type CtyunResponse struct {
	Request  *http.Request
	Response *http.Response
}

type CtyunResponseModel struct {
	StatusCode  interface{} `json:"statusCode"`
	ErrorCode   interface{} `json:"errorCode"`
	Message     string      `json:"message"`
	Description string      `json:"description"`
	ErrorMsg    string      `json:"error_msg"`
	ReturnObj   interface{} `json:"returnObj,omitempty"`
}

// IsSuccess 判断响应是否成功
func (c CtyunResponseModel) IsSuccess() bool {
	return c.ParseStatusCode() == StatusCodeSuccessString
}

// ParseErrorCode 解析ErrorCode到string
func (c CtyunResponseModel) ParseErrorCode() string {
	return c.parse(c.ErrorCode)
}

// ParseStatusCode 解析StatusCode到string
func (c CtyunResponseModel) ParseStatusCode() string {
	return c.parse(c.StatusCode)
}

// parse 解析code
func (c CtyunResponseModel) parse(code interface{}) string {
	codeString, okString := code.(string)
	if okString {
		return codeString
	}
	codeInt, okInt := code.(int)
	if okInt {
		return strconv.Itoa(codeInt)
	}
	codeFloat64, okFloat64 := code.(float64)
	if okFloat64 {
		return strconv.FormatFloat(codeFloat64, 'f', -1, 64)
	}
	return ""
	// return fmt.Sprintf("%v", code)
}

// ParseByStandardModelWithCheck 解析并且判断
func (c CtyunResponse) ParseByStandardModelWithCheck(obj interface{}) CtyunRequestError {
	model, err := c.ParseByStandardModel(obj)
	if err != nil {
		return err
	}

	if !model.IsSuccess() {
		var errInfos []string
		if model.Description != "" {
			errInfos = append(errInfos, model.Description)
		}
		if model.Message != "" {
			errInfos = append(errInfos, model.Message)
		}
		if model.ErrorMsg != "" {
			errInfos = append(errInfos, model.ErrorMsg)
		}
		code := model.ParseErrorCode()
		if code == "" {
			code = model.ParseStatusCode()
		}
		err := fmt.Errorf("API return error: " + strings.Join(errInfos, " "))
		wrapError := WrapWithErrorCode(err, code, &c)
		return wrapError
	}
	return nil
}

// Parse 解析为目标对象
func (c CtyunResponse) Parse(obj interface{}) CtyunRequestError {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(c.Response.Body)
	respBody, err := io.ReadAll(c.Response.Body)
	if err != nil {
		return WrapError(err, &c)
	}
	c.Response.Body = io.NopCloser(bytes.NewBuffer(respBody))
	err = json.Unmarshal(respBody, obj)
	if err != nil {
		return WrapError(err, &c)
	}
	return nil
}

// ParseByStandardModel 解析为标准模型对象
func (c CtyunResponse) ParseByStandardModel(obj interface{}) (*CtyunResponseModel, CtyunRequestError) {
	var model CtyunResponseModel
	model.ReturnObj = obj
	err := c.Parse(&model)
	if err != nil {
		return &CtyunResponseModel{}, err
	}
	return &model, nil
}
