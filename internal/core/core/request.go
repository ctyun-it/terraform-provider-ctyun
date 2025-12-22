package core

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
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
	/*	if contentType == "application/x-www-form-urlencoded" {
		encode := data.Encode()
		c.body = []byte(encode)
		c.AddHeader("Content-Type", contentType)
	}*/
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
	tim := time.Now().UTC()
	// 将时间转换为东八区时间（北京时间）
	location, _ := time.LoadLocation("Asia/Shanghai")
	localTime := tim.In(location)
	eopDate := localTime.Format("20060102T150405Z")
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

func readByte(fileName string) []byte {
	result := ""
	// 打开文件
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("文件打开失败 = ", err)
	}
	// 当函数退出时及时关闭file
	// defer在函数结束之前会调用defer后的代码
	// 及时关闭file句柄，否则会有内存泄漏
	defer file.Close()
	// 创建一个 *Reader ， 是带缓冲的
	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString('\n') // 读到一个换行就结束
		result = result + str
		if err == io.EOF { // io.EOF表示文件的末尾
			return []byte(result)
		}
	}
	return nil
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

func PostHttpForFormData(contentType, url, ak, sk string, headerMap map[string]string, fileMap, dataMap map[string]string) string {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	// 请求类型，例：multipart/form-data; boundary=7d8da3f8862ee22367dc1ea5c015d74dedbb2d834be8f1ef0545dc7f3534
	contentType = bodyWriter.FormDataContentType()
	fmt.Println(contentType)
	for temp := range fileMap {
		// key值
		fileNameKey := temp
		// value值
		fileNameVal := fileMap[temp]
		f, err := os.Open(fileNameVal)
		if err != nil {
			return err.Error()
		}
		fileWriter1, err := bodyWriter.CreateFormFile(fileNameKey, filepath.Base(fileNameVal))
		if err != nil {
			return err.Error()
		}
		_, err1 := io.Copy(fileWriter1, f)
		if err1 != nil {
			return err.Error()
		}
		err = bodyWriter.Close()
		if err != nil {
			return err.Error()
		}
	}
	// 拼接加密body
	var byteList [][]byte

	fmt.Println("--------------------------------------------------")
	for key, val := range dataMap {
		fileWriter2, err := bodyWriter.CreateFormField(key)
		if err != nil {
			return err.Error()
		}
		_, errs2 := fileWriter2.Write([]byte(val))
		if errs2 != nil {
			return err.Error()
		}
		// 一定要记着关闭
		err = bodyWriter.Close()
	}
	fmt.Println("--------------------------------------------------")
	newQuery := ""
	/*	var queryStr = strings.Split(query, "&")
		sort.Slice(queryStr, func(i, j int) bool {
			return queryStr[i] < queryStr[j] // 正序
		})
		for _, value := range queryStr {
			newQuery = newQuery + "&" + value
		}*/
	afterQuery := EncodeQueryStr(newQuery)
	if afterQuery != "" {
		url = url + "?" + afterQuery
	}
	// 获取请求头contentType中的boundary
	boundary := substring(contentType, strings.Index(contentType, "=")+1, len(contentType))
	// 获取boundary进行加密
	strings.Split(boundary, ";")
	boundary = "--" + boundary
	fmt.Println(boundary)

	for temp := range fileMap {
		// key值
		fileNameKey := temp
		// value值
		fileNameVal := fileMap[temp]
		// 打印
		fmt.Println(fileNameKey, fileNameVal)
		// 读取文件计算大小
		fi, _ := os.Stat(fileNameVal)
		size := fi.Size()
		fmt.Println(size)
		// 文件大小限制 1048576bytes = 1024kb = 1M
		if size >= 1048576 {
			fmt.Println("----------------文件过大----------------")
		}
		body1 := boundary + "\r\n" +
			"Content-Disposition: form-data; name=\"" + fileNameKey + "\"; filename=\"" + fi.Name() + "\"\r\n" +
			"Content-Type: application/octet-stream" + "\r\n" + "\r\n"
		fmt.Println(body1)
		body3 := "\r\n"
		order1 := []byte(body1)
		fmt.Println("body1:", len(string(order1)))
		order2 := readByte(fileNameVal)
		fmt.Println("body2:", len(string(order2)))
		order3 := []byte(body3)
		fmt.Println("body3:", len(string(order3)))

		// 计算数组长度
		arrLen := len(order1) + len(order2) + len(order3)
		// 声明
		bodyArr := make([]byte, arrLen)

		for i := 0; i < len(order1); i++ {
			bodyArr[i] = order1[i]
		}
		for i := 0; i < len(order2); i++ {
			bodyArr[len(order1)+i] = order2[i]
		}
		for i := 0; i < len(order3); i++ {
			bodyArr[len(order1)+len(order2)+i] = order3[i]
		}
		byteList = append(byteList, bodyArr)
	}

	lengths := 0
	for x := range byteList {
		lengths = lengths + len(byteList[x])
	}
	fmt.Println("byteList 总长度:", lengths)

	bodyArrs := make([]byte, lengths)
	num := 0
	for x := range byteList {
		for i := 0; i < len(byteList[x]); i++ {
			bodyArrs[num] = byteList[x][i]
			num = num + 1
		}
	}
	fmt.Println("bodyArrs 总长度:", len(bodyArrs))

	dataStr := ""
	if dataMap != nil && len(dataMap) > 0 {
		for temp := range dataMap {
			// key值
			dataNameKey := temp
			// value值
			dataNameVal := dataMap[temp]
			body := "Content-Disposition: form-data; name=\"" + dataNameKey + "\"\r\n" + "\r\n"
			nameVal := dataNameVal + "\r\n"
			dataStr = dataStr + boundary + "\r\n" + body + nameVal
		}
	}
	dataStrByte := []byte(dataStr)
	fmt.Println("dataStr 总长度:", len(dataStr))

	body4 := boundary + "--\r\n"
	order4 := []byte(body4)
	fmt.Println("order4 长度:", len(order4))

	lastBodyArr := make([]byte, len(bodyArrs)+len(dataStrByte)+len(order4))

	for i := 0; i < len(bodyArrs); i++ {
		lastBodyArr[i] = bodyArrs[i]
	}
	for i := 0; i < len(dataStrByte); i++ {
		lastBodyArr[len(bodyArrs)+i] = dataStrByte[i]
	}
	for i := 0; i < len(order4); i++ {
		lastBodyArr[len(bodyArrs)+len(dataStrByte)+i] = order4[i]
	}

	fmt.Println("lastBodyArr 长度:", len(lastBodyArr))

	uuId := uuid.New().String()
	fmt.Println("----------------------")
	fmt.Println(string(lastBodyArr))
	fmt.Println("----------------------")
	// 准备: HTTP请求
	reqBody := strings.NewReader(string(lastBodyArr))
	httpReq, err := http.NewRequest(http.MethodPost, url, reqBody)
	fmt.Println(url)
	if err != nil {
		fmt.Printf("NewRequest fail, url: %s, reqBody: %s, err: %v", "", bodyBuf, err)
		return err.Error()
	}
	// tim := time.Now()20221108T53Z
	tim := time.Now()
	eopDate := tim.Format("20060102T150405Z")
	httpReq.Header.Add("ctyun-eop-request-id", uuId)
	httpReq.Header.Add("Eop-Authorization", getSignByByMultipartFormDataBoundary(ak, sk, afterQuery, uuId, tim, lastBodyArr))
	httpReq.Header.Add("Eop-date", eopDate)
	httpReq.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	for temp := range headerMap {
		httpReq.Header.Add(temp, headerMap[temp])
		fmt.Println(temp + ":" + headerMap[temp])
	}
	// DO: HTTP请求
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("do http fail, url: %s, reqBody: %s, err:%v", "", bodyBuf, err)
		return err.Error()
	}
	headers := httpReq.Header
	for k, v := range headers {
		fmt.Println(k, v)
	}

	if httpRsp.StatusCode != http.StatusOK {
		fmt.Println()
		fmt.Println("应答头部 -----: ", httpRsp.Header)
		fmt.Println("请求StatusCode: ", httpRsp.StatusCode)
		fmt.Println()
	}
	defer httpRsp.Body.Close()

	// Read: HTTP结果
	rspBody, err := ioutil.ReadAll(httpRsp.Body)
	fmt.Println("rspBody: " + string(rspBody))
	if err != nil {
		fmt.Printf("ReadAll failed, url: %s, reqBody: %s, err: %v", "", bodyBuf, err)
		return err.Error()
	}
	return string(rspBody)
}

func getSignByByMultipartFormDataBoundary(ak, sk, afterQuery, uuId string, tim time.Time, bodyArr []byte) string {
	fmt.Println("---------------------------开始---------------------------")
	fmt.Println(string(bodyArr))
	fmt.Println("---------------------------结束---------------------------")
	// hash
	hash := sha256.New()
	hash.Write(bodyArr)
	sum1 := hash.Sum(nil)
	calculateContentHash := hex.EncodeToString(sum1)

	var sigture string

	singerDate := tim.Format("20060102T150405Z")
	singerDd := tim.Format("20060102")
	CampmocalHeader := "ctyun-eop-request-id:" + uuId + "\neop-date:" + singerDate + "\n"
	sigture = CampmocalHeader + "\n" + afterQuery + "\n" + calculateContentHash
	fmt.Println("sigture: " + sigture)
	kSecret := sk
	ktime := HmacSHA256(singerDate, kSecret)
	fmt.Println("ktime: " + hex.EncodeToString(ktime))
	kAk := HmacSHA256(ak, string(ktime))
	fmt.Println("kAk: " + hex.EncodeToString(kAk))
	kdate := HmacSHA256(singerDd, string(kAk))
	fmt.Println("kdate: " + hex.EncodeToString(kdate))
	signaSha256 := HmacSHA256(sigture, string(kdate))
	Signature := base64.StdEncoding.EncodeToString(signaSha256)
	fmt.Println("Signature: " + Signature)
	signHeader := ak + " Headers=ctyun-eop-request-id;eop-date Signature=" + Signature
	fmt.Println("signHeader: " + signHeader)
	fmt.Println()
	fmt.Println()
	return signHeader
}

func HmacSHA256(signature, key string) []byte {
	s := []byte(signature)
	k := []byte(key)
	m := hmac.New(sha256.New, k)
	m.Write(s)
	sum1 := m.Sum(nil)
	return sum1
}

func substring(source string, start int, end int) string {
	r := []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

func EncodeQueryStr(query string) string {
	afterQuery := ""
	if len(query) != 0 {
		n := strings.Split(query, "&")
		for _, v := range n {
			if len(afterQuery) < 1 {
				a := strings.Split(v, "=")
				if len(a) >= 2 {
					encodeStr := url.QueryEscape(a[1])
					v = a[0] + "=" + encodeStr
					afterQuery = afterQuery + v
				} else {
					encodeStr := ""
					v = a[0] + "=" + encodeStr
					afterQuery = afterQuery + v
				}
			} else {
				a := strings.Split(v, "=")
				if len(a) >= 2 {
					encodeStr := url.QueryEscape(a[1])
					v = a[0] + "=" + encodeStr
					afterQuery = afterQuery + "&" + v
				} else {
					encodeStr := ""
					v = a[0] + "=" + encodeStr
					afterQuery = afterQuery + "&" + v
				}
			}
		}
	}

	return afterQuery
}

func String2Map(header string) map[string]string {
	data := []byte(header)
	map2 := make(map[string]interface{})
	err := json.Unmarshal(data, &map2)
	if err != nil {
		fmt.Println(err)
	}
	headerMap := MapInterface2String(map2)
	return headerMap
}

func MapInterface2String(inputData map[string]interface{}) map[string]string {
	outputData := map[string]string{}
	for key, value := range inputData {
		switch value.(type) {
		case string:
			outputData[key] = value.(string)
		}
	}
	return outputData
}

func StructToHeader(request interface{}) string {
	var sb strings.Builder
	v := reflect.ValueOf(request)

	// 确保我们是传入的指针
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 遍历结构体字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldValue := v.Field(i)

		// 如果是字符串类型，我们生成类似 JSON 的格式
		if fieldValue.Kind() == reflect.String {
			// 这里动态拼接 JSON 格式的字符串
			if sb.Len() > 0 {
				sb.WriteString(",") // 在多个字段之间添加逗号
			}
			// 拼接字段名和字段值
			sb.WriteString(fmt.Sprintf("\"%s\":\"%s\"", field.Name, fieldValue.String()))
		}
	}

	// 最终生成 header 格式字符串
	return "{" + sb.String() + "}"
}

func StructToFileMap(request interface{}) map[string]string {
	fileMap := make(map[string]string)
	v := reflect.ValueOf(request)

	// 确保我们是传入的指针
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 遍历结构体字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldValue := v.Field(i)

		// 如果字段是字符串类型并且非空，加入到 fileMap
		if fieldValue.Kind() == reflect.String && fieldValue.String() != "" {
			fileMap[strings.ToLower(field.Name)] = fieldValue.String()
		}
	}

	// 返回生成的 fileMap
	return fileMap
}
