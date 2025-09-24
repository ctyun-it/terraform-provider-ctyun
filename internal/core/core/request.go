package core

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CtyunRequest struct {
	endpointName EndpointName // 终端名称
	credential   Credential   // 密钥信息
	method       string       // 请求方法
	urlPath      string       // url路径
	headers      http.Header  // 请求头
	params       url.Values   // 请求param参数
	body         []byte       // 请求body
}

type CtyunRequestTemplate struct {
	EndpointName EndpointName // 终端名称
	Method       string       // 请求方法
	UrlPath      string       // url路径
	ContentType  string
}

type CtyunRequestBuilder struct {
	EndpointName EndpointName // 终端名称
	Method       string       // 请求方法
	UrlPath      string       // url路径
	Credential   Credential   // 用户信息
	ContentType  string       // 请求类型
}

func NewCtyunRequestBuilder(template CtyunRequestTemplate) *CtyunRequestBuilder {
	return &CtyunRequestBuilder{
		EndpointName: template.EndpointName,
		Method:       template.Method,
		UrlPath:      template.UrlPath,
		ContentType:  template.ContentType,
	}
}

// ReplaceUrl 替换路径中的目标值，例如把/orders/{masterOrderId}替换为/orders/1
func (c *CtyunRequestBuilder) ReplaceUrl(src string, target interface{}) *CtyunRequestBuilder {
	str := fmt.Sprintf("%v", target)
	str = url.PathEscape(str)
	c.UrlPath = strings.Replace(c.UrlPath, "{"+src+"}", str, -1)
	return c
}

// WithCredential 增加请求credential
func (c *CtyunRequestBuilder) WithCredential(credential Credential) *CtyunRequestBuilder {
	c.Credential = credential
	return c
}

// WithEndpointName 增加请求终端名称
func (c *CtyunRequestBuilder) WithEndpointName(endpointName EndpointName) *CtyunRequestBuilder {
	c.EndpointName = endpointName
	return c
}

// Build 构造
func (c CtyunRequestBuilder) Build() *CtyunRequest {
	return &CtyunRequest{
		endpointName: c.EndpointName,
		method:       c.Method,
		urlPath:      c.UrlPath,
		credential:   c.Credential,
		headers:      make(http.Header),
		params:       make(url.Values),
	}
}

// AddHeader 增加请求头
func (c *CtyunRequest) AddHeader(key, value string) *CtyunRequest {
	c.headers[key] = append(c.headers[key], value)
	return c
}

func (c *CtyunRequest) AddHeaders(key string, value []string) *CtyunRequest {
	c.headers[key] = append(c.headers[key], value...)
	return c
}

// AddParam 增加参数
func (c *CtyunRequest) AddParam(key, value string) *CtyunRequest {
	c.params.Add(key, value)
	return c
}

// AddParams 增加参数
func (c *CtyunRequest) AddParams(key string, value []string) *CtyunRequest {
	for _, v := range value {
		c.AddParam(key, v)
	}
	return c
}

// WriteXWwwFormUrlEncoded 以x-www-form-urlencoded方式写入
// func (c *CtyunRequest) WriteXWwwFormUrlEncoded(data url.Values) *CtyunRequest {
// 	encode := data.Encode()
// 	c.body = []byte(encode)
// 	c.AddHeader("Content-Type", "application/x-www-form-urlencoded")
// 	return c
// }

// WriteJson 以application/json方式写入
func (c *CtyunRequest) WriteJson(data interface{}, contentType string) (*CtyunRequest, error) {
	if contentType == "application/json" {
		marshal, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		c.body = marshal
		c.AddHeader("Content-Type", contentType)
	}
	return c, nil
}

func (c *CtyunRequest) WriteString(data, contentType string) (*CtyunRequest, error) {
	if contentType == "text/plain" {
		c.body = []byte(data)
		c.AddHeader("Content-Type", contentType)
	}
	return c, nil
}

// buildRequest 构造请求
func (c CtyunRequest) buildRequest(endPoint string) (*http.Request, error) {
	// 构造url
	u := endPoint + c.urlPath
	query := c.params.Encode()
	if query != "" {
		u = u + "?" + query
	}

	// 构造请求头
	tim := time.Now()
	eopDate := tim.Format("20060102T150405Z")
	id := uuid.NewString()
	sign := GetSign(query, c.body, eopDate, id, c.credential)
	headers := c.headers.Clone()
	headers.Add("ctyun-eop-request-id", id)
	headers.Add("Eop-Authorization", sign)
	headers.Add("Eop-date", eopDate)
	headers.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	if c.body != nil {
		headers.Add("Content-Length", strconv.Itoa(len(c.body)))
	}

	// 构造实际请求
	req, err := http.NewRequest(c.method, u, bytes.NewReader(c.body))
	if err != nil {
		return nil, err
	}
	req.Header = headers
	return req, nil
}

// GetSign 加签
func GetSign(query string, body []byte, date string, uuid string, credential Credential) string {
	hash := sha256.New()
	hash.Write(body)
	sum := hash.Sum(nil)

	calculateContentHash := hex.EncodeToString(sum)
	signature := fmt.Sprintf("ctyun-eop-request-id:%s\neop-date:%s\n\n%s\n%s", uuid, date, query, calculateContentHash)
	singerDd := date[0:8]
	s := hmacSHA256(date, credential.sk)
	kAk := hmacSHA256(credential.ak, string(s))
	kDate := hmacSHA256(singerDd, string(kAk))
	signatureSha256 := hmacSHA256(signature, string(kDate))
	signatureBase64 := base64.StdEncoding.EncodeToString(signatureSha256)
	return credential.ak + " Headers=ctyun-eop-request-id;eop-date Signature=" + signatureBase64
}

// hmacSHA256 HmacSHA256加密
func hmacSHA256(signature string, key string) []byte {
	s := []byte(signature)
	k := []byte(key)
	m := hmac.New(sha256.New, k)
	m.Write(s)
	sum := m.Sum(nil)
	return sum
}
