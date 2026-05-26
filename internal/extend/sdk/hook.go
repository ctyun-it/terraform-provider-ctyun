package sdk

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"net/http/httputil"
	"path/filepath"
	"runtime"
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
	TerraformVersion string
	ProviderVersion  string
}

func (m MetricHttpHook) BeforeRequest(ctx context.Context, request *http.Request) {
	typ, name, action := GetResourceInfoFromStack()
	request.Header.Set("X-Terraform-Source-Type", typ)
	request.Header.Set("X-Terraform-Source-Name", name)
	request.Header.Set("X-Terraform-Source-Action", action)

	request.Header.Set("User-Agent", fmt.Sprintf("Terraform/%s terraform-provider-ctyun/%s", m.TerraformVersion, m.ProviderVersion))
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

// GetResourceInfoFromStack 从调用栈中获取 Terraform 资源类型、资源名称和操作类型
// 返回值: typ (resource/datasource/unknown), name (资源名), action (Create/Update/Delete/Read/ImportState/unknown)
func GetResourceInfoFromStack() (typ string, name string, action string) {
	typ = "unknown"
	action = "unknown"
	// 从 i=3 开始跳过自身，性能更好
	for i := 3; i < 30; i++ {
		pc, file, _, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// 尝试从函数名中提取操作类型
		if a := getActionFromPC(pc); a != "unknown" {
			action = a
		}

		filename := filepath.Base(file)

		// 匹配 resource_*.go，记录但不立即返回，继续向外遍历寻找 action
		if typ == "unknown" && strings.HasPrefix(filename, "resource_") && strings.HasSuffix(filename, ".go") {
			resName := strings.TrimSuffix(strings.TrimPrefix(filename, "resource_"), ".go")
			typ = "resource"
			name = resName
		}

		// 匹配 datasource_*.go
		if typ == "unknown" && strings.HasPrefix(filename, "datasource_") && strings.HasSuffix(filename, ".go") {
			dataName := strings.TrimSuffix(strings.TrimPrefix(filename, "datasource_"), ".go")
			typ = "datasource"
			name = dataName
		}

		// typ 和 action 都已找到，提前退出
		if typ != "unknown" && action != "unknown" {
			break
		}
	}

	// datasource 只有读操作
	if typ == "datasource" && action == "unknown" {
		action = "Read"
	}

	return typ, name, action
}

// getActionFromPC 从 program counter 中提取 Terraform 操作类型
func getActionFromPC(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	funcName := fn.Name()
	// 提取函数名（去掉包路径）
	parts := strings.Split(funcName, ".")
	if len(parts) == 0 {
		return "unknown"
	}
	methodName := parts[len(parts)-1]
	switch methodName {
	case "Create", "Update", "Delete", "Read", "ImportState":
		return methodName
	}
	return "unknown"
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
