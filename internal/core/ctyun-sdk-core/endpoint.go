package ctyunsdk

import (
	"errors"
)

type EndpointName string

type Endpoint struct {
	EndpointName EndpointName // 端点的名称
	Url          string       // 对应的地址
}

type EndpointRegistry map[EndpointName]Endpoint

// GetEndpoint 获取端点信息
func (e EndpointRegistry) GetEndpoint(endpointName EndpointName) (Endpoint, bool) {
	target, ok := e[endpointName]
	return target, ok
}

// Register 注册端点
func (e *EndpointRegistry) Register(endpoint Endpoint) error {
	_, ok := e.GetEndpoint(endpoint.EndpointName)
	if ok {
		return errors.New("端点名称：" + string(endpoint.EndpointName) + "已被注册")
	}
	(*e)[endpoint.EndpointName] = endpoint
	return nil
}
