package core

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type CtyunResponse struct {
	Request  *http.Request
	Response *http.Response
}

// Parse 解析为目标对象
func (c CtyunResponse) Parse(obj interface{}) error {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(c.Response.Body)
	respBody, err := ioutil.ReadAll(c.Response.Body)
	if err != nil {
		return err
	}
	c.Response.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))
	err = json.Unmarshal(respBody, obj)
	if err != nil {
		return err
	}
	return nil
}
