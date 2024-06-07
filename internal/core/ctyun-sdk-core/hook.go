package ctyunsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"time"
)

// ApiHandler api处理器
type ApiHandler interface {
	// Do 执行api动作
	Do(context.Context, *Credential, interface{}) (interface{}, CtyunRequestError)
}

// ApiDecoratorChain 处理链
type ApiDecoratorChain struct {
	Hooks  []ApiHook
	Target interface{}
}

// Next 驱动处理链执行下个动作
func (this ApiDecoratorChain) Next(ctx context.Context, credential *Credential, param interface{}) (interface{}, CtyunRequestError) {
	chain := ApiDecoratorChain{
		Hooks:  this.Hooks[1:],
		Target: this.Target,
	}
	return this.Hooks[0](ctx, credential, param, chain)
}

// ApiHook api钩子
type ApiHook func(context.Context, *Credential, interface{}, ApiDecoratorChain) (interface{}, CtyunRequestError)

// wrapperDecorator 代理类，装饰者
type wrapperDecorator struct {
	target ApiHandler
	hooks  []ApiHook
}

// Do 覆盖执行Do函数，驱动装饰者的hooks执行
func (this wrapperDecorator) Do(ctx context.Context, credential *Credential, t interface{}) (interface{}, CtyunRequestError) {
	var hs []ApiHook
	for _, hook := range this.hooks {
		hs = append(hs, hook)
	}
	// 把需要执行的动作target放在chain的末尾，成为一个hook的装饰着
	hs = append(
		hs,
		func(ct context.Context, c *Credential, t interface{}, _ ApiDecoratorChain) (interface{}, CtyunRequestError) {
			return this.target.Do(ct, c, t)
		},
	)
	chain := ApiDecoratorChain{
		Hooks:  hs,
		Target: interface{}(this.target),
	}
	next, err := chain.Next(ctx, credential, t)
	return next, err
}

type ApiHookBuilder struct {
	hooks []ApiHook
}

// AddHooks 添加钩子
func (this *ApiHookBuilder) AddHooks(hook ApiHook) *ApiHookBuilder {
	this.hooks = append(this.hooks, hook)
	return this
}

func NewApiHookBuilder() *ApiHookBuilder {
	return &ApiHookBuilder{}
}

// Wrap 构造
func Wrap(target ApiHandler, builder ApiHookBuilder) ApiHandler {
	return wrapperDecorator{
		target: target,
		hooks:  builder.hooks,
	}
}

// ConsoleLogApiHook 打印控制台日志的钩子函数
func ConsoleLogApiHook(ctx context.Context, credential *Credential, param interface{}, chain ApiDecoratorChain) (interface{}, CtyunRequestError) {
	startTime := time.Now().UnixMilli()
	req, err := json.Marshal(param)
	id := uuid.NewString()
	apiName := fmt.Sprintf("%T", param)
	if err == nil {
		fmt.Printf("开始执行请求动作，id：%s，请求：%s，内容：%s\n", id, apiName, string(req))
	}
	result, ctyunRequestError := chain.Next(ctx, credential, param)
	endTime := time.Now().UnixMilli()
	useTime := endTime - startTime
	if ctyunRequestError == nil {
		resp, err := json.Marshal(result)
		if err == nil {
			fmt.Printf("执行请求动作成功，id：%s，请求：%s，花费时间：%d毫秒，返回信息：%s\n", id, apiName, useTime, string(resp))
		}
	} else {
		fmt.Printf("执行请求动作失败，id：%s，请求：%s，花费时间：%d毫秒，返回信息：%s\n", id, apiName, useTime, ctyunRequestError.Error())
	}
	return result, ctyunRequestError
}

type StubFunc func(ctx context.Context, credential *Credential, param interface{}, target interface{}) (interface{}, CtyunRequestError)

type StubApiHook struct {
	stubs map[string]StubFunc
}

func NewStubApiHook() *StubApiHook {
	return &StubApiHook{stubs: make(map[string]StubFunc)}
}

// RegisterStubWithApiHandler 注册桩
func (this *StubApiHook) RegisterStubWithApiHandler(handler interface{}, f StubFunc) {
	if handler == nil || f == nil {
		return
	}
	this.stubs[this.getApiHandlerTargetName(handler)] = f
}

// RegisterStubWithName 注册桩
func (this *StubApiHook) RegisterStubWithName(stubName string, f StubFunc) {
	if stubName == "" || f == nil {
		return
	}
	this.stubs[stubName] = f
}

// getApiHandlerTargetName 获取名称
func (this StubApiHook) getApiHandlerTargetName(handler interface{}) string {
	to := reflect.TypeOf(handler)
	if to.Kind() == reflect.Ptr {
		return to.Elem().Name()
	}
	return to.Name()
}

// Stub 打桩
func (this StubApiHook) Stub(ctx context.Context, credential *Credential, param interface{}, chain ApiDecoratorChain) (interface{}, CtyunRequestError) {
	name := this.getApiHandlerTargetName(chain.Target)
	stubFunc := this.stubs[name]
	if stubFunc == nil {
		return chain.Next(ctx, credential, param)
	}
	return stubFunc(ctx, credential, param, chain.Target)
}
