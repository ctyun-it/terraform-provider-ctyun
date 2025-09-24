package sdk

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"net/http/httputil"
	"strings"
)

const (
	HeaderConsoleUrl = "consoleUrl"
)

type LogHttpHook struct {
}

func (d LogHttpHook) BeforeRequest(ctx context.Context, request *http.Request) {
	dumpRequest, err := httputil.DumpRequest(request, true)
	if err != nil {
		return
	}
	requestContent := string(dumpRequest)
	tflog.Info(ctx, "实际请求内容：", map[string]interface{}{"request": requestContent})
}

func (d LogHttpHook) AfterResponse(ctx context.Context, response *http.Response) {
	if response != nil {
		dumpResponse, err := httputil.DumpResponse(response, true)
		if err != nil {
			return
		}
		responseContent := string(dumpResponse)
		tflog.Info(ctx, "实际请求返回：\n", map[string]interface{}{"response": responseContent})
		return
	}
	tflog.Info(ctx, "实际请求返回空：\n", map[string]interface{}{"response": response})
}

// MetricHttpHook 使用metric发送埋点日志
type MetricHttpHook struct {
}

func (m MetricHttpHook) BeforeRequest(_ context.Context, request *http.Request) {
	request.Header.Set("From-Terraform-Provider", "true")
}

func (m MetricHttpHook) AfterResponse(_ context.Context, _ *http.Response) {

}

// AddConsoleUrlHook 确定发送console的定位，产线研一那边测试环境通过consoleUrl改造定位
type AddConsoleUrlHook struct {
	consoleUrl       string
	addHeaderHandler addHeaderHandler
}

type addHeaderHandler interface {
	AddHeader(request *http.Request, url string)
}

func NewAddConsoleUrlHook(consoleUrl string, endpoints ...string) *AddConsoleUrlHook {
	if consoleUrl == "" || len(endpoints) == 0 {
		return &AddConsoleUrlHook{
			consoleUrl:       "",
			addHeaderHandler: noOperationAddAddHeaderHandler{},
		}
	}

	var addHeaderHandler addHeaderHandler
	// 拦截所有请求
	for _, endpoint := range endpoints {
		if endpoint == "*" {
			addHeaderHandler = newSimpleAddAddHeaderHandler()
			break
		}
	}
	// 兜底情况
	if addHeaderHandler == nil {
		addHeaderHandler = newKeywordEndpointAddHeaderHandler(endpoints)
	}
	return &AddConsoleUrlHook{
		consoleUrl:       consoleUrl,
		addHeaderHandler: addHeaderHandler,
	}
}

func (m AddConsoleUrlHook) BeforeRequest(_ context.Context, request *http.Request) {
	m.addHeaderHandler.AddHeader(request, m.consoleUrl)
}

func appendHeaderIfNotExist(request *http.Request, url string) {
	_, ok := request.Header[HeaderConsoleUrl]
	if !ok {
		request.Header[HeaderConsoleUrl] = []string{url}
	}
}

func (m AddConsoleUrlHook) AfterResponse(_ context.Context, _ *http.Response) {

}

type noOperationAddAddHeaderHandler struct{}

func (n noOperationAddAddHeaderHandler) AddHeader(_ *http.Request, _ string) {
}

type simpleAddAddHeaderHandler struct{}

func newSimpleAddAddHeaderHandler() *simpleAddAddHeaderHandler {
	return &simpleAddAddHeaderHandler{}
}

func (a simpleAddAddHeaderHandler) AddHeader(request *http.Request, url string) {
	appendHeaderIfNotExist(request, url)
}

type keywordEndpointAddHeaderHandler struct {
	keywords []string
}

func newKeywordEndpointAddHeaderHandler(keywords []string) *keywordEndpointAddHeaderHandler {
	return &keywordEndpointAddHeaderHandler{
		keywords: keywords,
	}
}

func (k keywordEndpointAddHeaderHandler) AddHeader(request *http.Request, url string) {
	for _, keyword := range k.keywords {
		if strings.Contains(request.Host, keyword) {
			appendHeaderIfNotExist(request, url)
			return
		}
	}
}
