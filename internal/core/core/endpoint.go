package core

import (
	"errors"
	"sync"
)

type EndpointName string

var mu sync.RWMutex

type Endpoint struct {
	Name EndpointName // 端点的名称
	Url  string       // 对应的地址
}

type EndpointRegistry struct {
	mapping    map[EndpointName]string
	defaultUrl string
}

// DefaultEndpointRegistry 由默认的url创建
func DefaultEndpointRegistry() *EndpointRegistry {
	return &EndpointRegistry{
		mapping: make(map[EndpointName]string),
	}
}

// GetEndpointUrl 获取端点信息
func (e EndpointRegistry) GetEndpointUrl(endpointName EndpointName) string {
	mu.RLock()
	defer mu.RUnlock()
	url, ok := e.mapping[endpointName]
	if ok {
		return url
	}
	return e.defaultUrl
}

// SetDefaultUrl 设置默认的url
func (e *EndpointRegistry) SetDefaultUrl(defaultUrl string) {
	mu.Lock()
	defer mu.Unlock()
	e.defaultUrl = defaultUrl
}

// Register 注册端点
func (e *EndpointRegistry) Register(endpoint Endpoint) error {
	mu.Lock()
	defer mu.Unlock()
	_, ok := e.mapping[endpoint.Name]
	if ok {
		return errors.New("endpoint name: [" + string(endpoint.Name) + "] already registered!")
	}
	e.mapping[endpoint.Name] = endpoint.Url
	return nil
}
