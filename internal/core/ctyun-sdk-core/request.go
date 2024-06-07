package ctyunsdk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type CtyunRequest struct {
	method     string      // 请求方法
	urlPath    string      // url路径
	credential Credential  // 密钥信息
	headers    http.Header // 请求头
	params     url.Values  // 请求param参数
	body       []byte      // 请求body
}

type CtyunRequestBuilder struct {
	Method  string // 请求方法
	UrlPath string // url路径
}

// WithCredential 增加请求credential
func (c CtyunRequestBuilder) WithCredential(credential *Credential) *CtyunRequest {
	result := &CtyunRequest{
		method:  c.Method,
		urlPath: c.UrlPath,
		headers: make(http.Header),
		params:  make(url.Values),
	}
	if credential != nil {
		result.credential = *credential
	}
	return result
}

// AddHeader 增加请求头
func (c *CtyunRequest) AddHeader(key, value string) *CtyunRequest {
	c.headers.Add(key, value)
	return c
}

// AddParam 增加参数
func (c *CtyunRequest) AddParam(key, value string) *CtyunRequest {
	c.params.Add(key, value)
	return c
}

// WriteXWwwFormUrlEncoded 以x-www-form-urlencoded方式写入
func (c *CtyunRequest) WriteXWwwFormUrlEncoded(data url.Values) *CtyunRequest {
	encode := data.Encode()
	c.body = []byte(encode)
	c.AddHeader("Content-Type", "application/x-www-form-urlencoded")
	return c
}

// WriteJson 以application/json方式写入
func (c *CtyunRequest) WriteJson(data interface{}) (*CtyunRequest, CtyunRequestError) {
	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, ErrorBeforeRequest(err)
	}
	c.body = marshal
	c.AddHeader("Content-Type", "application/json")
	return c, nil
}

// buildRequest 构造请求
func (c CtyunRequest) buildRequest(endPoint string) (*http.Request, CtyunRequestError) {
	// 构造url
	u := url.URL{}
	u.Scheme = "https"
	u.Host = endPoint
	u.Path = c.urlPath
	u.RawQuery = c.params.Encode()

	// 构造请求头
	tim := time.Now()
	eopDate := tim.Format("20060102T150405Z")
	id := uuid.NewString()
	sign := GetSign(u.RawQuery, c.body, tim, id, c.credential)
	headers := c.headers.Clone()
	headers.Add("ctyun-eop-request-id", id)
	headers.Add("Eop-Authorization", sign)
	headers.Add("Eop-date", eopDate)
	if c.body != nil {
		headers.Add("Content-Length", strconv.Itoa(len(c.body)))
	}

	// 构造实际请求
	req, err := http.NewRequest(c.method, u.String(), bytes.NewReader(c.body))
	if err != nil {
		return nil, ErrorBeforeRequest(err)
	}
	req.Header = headers
	return req, nil
}

// GetSign 加签
func GetSign(query string, body []byte, tim time.Time, uuid string, credential Credential) string {
	hash := sha256.New()
	hash.Write(body)
	sum := hash.Sum(nil)
	calculateContentHash := hex.EncodeToString(sum)
	date := tim.Format("20060102T150405Z")
	sigture := fmt.Sprintf("ctyun-eop-request-id:%s\neop-date:%s\n\n%s\n%s", uuid, date, query, calculateContentHash)
	singerDd := tim.Format("20060102")
	kAk := hmacSHA256(credential.ak, string(hmacSHA256(date, credential.sk)))
	kdate := hmacSHA256(singerDd, string(kAk))
	signaSha256 := hmacSHA256(sigture, string(kdate))
	Signature := base64.StdEncoding.EncodeToString(signaSha256)
	signHeader := credential.ak + " Headers=ctyun-eop-request-id;eop-date Signature=" + Signature
	return signHeader
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
