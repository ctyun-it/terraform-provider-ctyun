package ctyunsdk

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
)

type Environment string

const (
	EnvironmentDev  = "dev"
	EnvironmentTest = "test"
	EnvironmentProd = "prod"
)

var Environments = []Environment{
	EnvironmentDev,
	EnvironmentTest,
	EnvironmentProd,
}

type CtyunClient struct {
	Env      Environment
	Config   *CtyunClientConfig
	registry map[Environment]EndpointRegistry
}

// ClientConfigTest 构建测试环境默认的客户端
func ClientConfigTest() *CtyunClientConfig {
	return &CtyunClientConfig{
		ApiHooks: []ApiHook{
			ConsoleLogApiHook,
		},
		HttpHooks: []HttpHook{
			PrintLogHttpHook{},
			AddUserAgentHttpHook{},
		},
		Client: ClientTest(),
	}
}

// ClientTest 测试环境客户端
func ClientTest() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 0,
	}
}

// ClientConfigProd 构建生产环境默认的客户端
func ClientConfigProd() *CtyunClientConfig {
	return &CtyunClientConfig{
		ApiHooks: []ApiHook{
			ConsoleLogApiHook,
		},
		HttpHooks: []HttpHook{
			PrintLogHttpHook{},
			AddUserAgentHttpHook{},
		},
		Client: ClientProd(),
	}
}

// ClientProd 生产环境客户端
func ClientProd() *http.Client {
	return &http.Client{}
}

// NewCtyunClient 新建环境
func NewCtyunClient(env Environment, cfg *CtyunClientConfig) (*CtyunClient, error) {
	client := cfg.Client
	if cfg.Client == nil {
		client = &http.Client{}
	}
	c := &CtyunClient{
		Config: &CtyunClientConfig{
			Client:    client,
			ApiHooks:  cfg.ApiHooks,
			HttpHooks: cfg.HttpHooks,
		},
	}
	c.Config = cfg
	c.registry = make(map[Environment]EndpointRegistry)
	c.Env = env
	return c, nil
}

type CtyunClientConfig struct {
	Client    *http.Client
	ApiHooks  []ApiHook
	HttpHooks []HttpHook
}

// EnvOf 通过指定环境构建
func EnvOf(env Environment) *CtyunClient {
	switch env {
	case EnvironmentDev:
		client, _ := NewCtyunClient(
			EnvironmentDev,
			ClientConfigTest(),
		)
		return client
	case EnvironmentTest:
		client, _ := NewCtyunClient(
			EnvironmentTest,
			ClientConfigTest(),
		)
		return client
	case EnvironmentProd:
		fallthrough
	default:
		client, _ := NewCtyunClient(
			EnvironmentProd,
			ClientConfigProd(),
		)
		return client
	}
}

// RegisterEndpoint 注册端点
func (c *CtyunClient) RegisterEndpoint(env Environment, endpoint Endpoint) {
	endpointRegisty, ok := c.registry[env]
	if !ok {
		endpointRegisty = make(EndpointRegistry)
		c.registry[env] = endpointRegisty
	}
	endpointRegisty[endpoint.EndpointName] = endpoint
}

// RequestToEndpoint 向端点发送请求
func (c CtyunClient) RequestToEndpoint(ctx context.Context, endpointName EndpointName, request *CtyunRequest) (*CtyunResponse, CtyunRequestError) {
	url, ok := c.acquireEndpointUrl(c.Env, endpointName)
	if !ok {
		return nil, ErrorBeforeRequest(errors.New("无法找到端点" + string(endpointName) + "对应的请求URL"))
	}
	req, err := request.buildRequest(url)
	if err != nil {
		return nil, ErrorBeforeResponse(err)
	}
	return c.send(ctx, req)
}

// send 发送请求
func (c CtyunClient) send(ctx context.Context, req *http.Request) (*CtyunResponse, CtyunRequestError) {
	for _, hook := range c.Config.HttpHooks {
		hook.BeforeRequest(ctx, req)
	}
	resp, err := c.Config.Client.Do(req)
	for _, hook := range c.Config.HttpHooks {
		hook.AfterResponse(ctx, resp)
	}
	ctyunResp := &CtyunResponse{Request: req, Response: resp}
	if err != nil {
		return nil, ErrorAfterResponse(err, ctyunResp)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, ErrorAfterResponse(errors.New("响应失败，状态码："+resp.Status), ctyunResp)
	}
	return ctyunResp, nil
}

// acquireEndpointUrl 获取端点url
func (c CtyunClient) acquireEndpointUrl(env Environment, name EndpointName) (string, bool) {
	endpointRegistry, ok := c.registry[env]
	if !ok {
		return "", false
	}
	endpoint, ok := endpointRegistry.GetEndpoint(name)
	if !ok {
		return "", false
	}
	return endpoint.Url, true
}

type HttpHook interface {
	BeforeRequest(context.Context, *http.Request)
	AfterResponse(context.Context, *http.Response)
}

type PrintLogHttpHook struct {
}

func (d PrintLogHttpHook) BeforeRequest(_ context.Context, request *http.Request) {
	dumpRequest, err := httputil.DumpRequest(request, true)
	if err != nil {
		return
	}
	requestContent := string(dumpRequest)
	fmt.Printf("实际请求内容：\n%s\n", requestContent)
}

func (d PrintLogHttpHook) AfterResponse(_ context.Context, response *http.Response) {
	dumpResponse, err := httputil.DumpResponse(response, true)
	if err != nil {
		return
	}
	responseContent := string(dumpResponse)
	fmt.Printf("实际请求返回：\n%s\n", responseContent)
}

type AddUserAgentHttpHook struct {
}

func (h AddUserAgentHttpHook) BeforeRequest(_ context.Context, request *http.Request) {
	// 不添加请求头会出现被风控的现象
	if request.Header.Get("User-Agent") == "" {
		request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	}
}

func (h AddUserAgentHttpHook) AfterResponse(_ context.Context, _ *http.Response) {
}
