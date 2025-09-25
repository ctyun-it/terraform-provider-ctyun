package core

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	CtyunEndpointName = "ctyunEndpointName"
)

type CtyunClient struct {
	Config   *CtyunClientConfig
	registry *EndpointRegistry
}

// ClientConfigForTest 构建测试环境默认的客户端
func ClientConfigForTest() *CtyunClientConfig {
	return &CtyunClientConfig{
		HttpHooks: []HttpHook{
			PrintLogHttpHook{},
		},
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: 0,
		},
	}
}

// ClientConfigForProd 构建生产环境默认的客户端
func ClientConfigForProd() *CtyunClientConfig {
	return &CtyunClientConfig{
		HttpHooks: []HttpHook{},
		Client:    &http.Client{},
	}
}

// NewCtyunClient 新建客户端
func NewCtyunClient(cfg *CtyunClientConfig) *CtyunClient {
	client := cfg.Client
	if cfg.Client == nil {
		client = &http.Client{}
	}
	var hooks []HttpHook
	for _, h := range cfg.HttpHooks {
		hooks = append(hooks, h)
	}
	return &CtyunClient{
		Config: &CtyunClientConfig{
			Client:    client,
			HttpHooks: hooks,
		},
		registry: DefaultEndpointRegistry(),
	}
}

// ClientForProd 生产环境客户端
func ClientForProd() *CtyunClient {
	return NewCtyunClient(ClientConfigForProd())
}

// ClientForTest 测试环境客户端
func ClientForTest() *CtyunClient {
	return NewCtyunClient(ClientConfigForTest())
}

// DefaultClient 默认客户端
func DefaultClient() *CtyunClient {
	return ClientForTest()
}

// CtyunClientConfig 自定义配置
type CtyunClientConfig struct {
	Client    *http.Client
	HttpHooks []HttpHook
}

// RegisterEndpoint 注册端点
func (c *CtyunClient) RegisterEndpoint(endpoint Endpoint) {
	_ = c.registry.Register(endpoint)
}

// RequestToEndpoint 向端点发送请求
func (c CtyunClient) RequestToEndpoint(ctx context.Context, request *CtyunRequest) (*CtyunResponse, error) {
	defaultUrl := c.registry.GetEndpointUrl(request.endpointName)
	req, err := request.buildRequest(defaultUrl)
	if err != nil {
		return nil, err
	}
	return c.send(ctx, req)
}

// send 发送请求
func (c CtyunClient) send(ctx context.Context, req *http.Request) (*CtyunResponse, error) {
	for _, hook := range c.Config.HttpHooks {
		hook.BeforeRequest(ctx, req)
	}
	resp, err := c.Config.Client.Do(req)
	for _, hook := range c.Config.HttpHooks {
		hook.AfterResponse(ctx, resp)
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusBadRequest && strings.Contains(req.URL.Path, "/v2/cce") {
		resp.StatusCode = http.StatusOK
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.New("response error, status code: " + strconv.Itoa(resp.StatusCode))
	}
	return &CtyunResponse{Request: req, Response: resp}, nil
}
