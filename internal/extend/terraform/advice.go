package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"time"
)

type ResourceDecoratorChain[T any, R any] struct {
	Chains []Advice[T, R]
}

// Next 驱动处理链执行下个动作
func (this ResourceDecoratorChain[T, R]) Next(ctx context.Context, request T, response R) error {
	chain := ResourceDecoratorChain[T, R]{
		Chains: this.Chains[1:],
	}
	return this.Chains[0].Around(ctx, request, response, chain)
}

type Advice[T any, R any] interface {
	Around(context.Context, T, R, ResourceDecoratorChain[T, R]) error
}

type adviceWrapper[T any, R any] struct {
	Target func(context.Context, T, R)
}

func (this adviceWrapper[T, R]) Around(ctx context.Context, t T, r R, _ ResourceDecoratorChain[T, R]) error {
	this.Target(ctx, t, r)
	return nil
}

type AopAdvices struct {
	DataSourceReadAopApi []Advice[datasource.ReadRequest, *datasource.ReadResponse]
	ResourceReadAopApi   []Advice[resource.ReadRequest, *resource.ReadResponse]
	ResourceCreateAopApi []Advice[resource.CreateRequest, *resource.CreateResponse]
	ResourceDeleteAopApi []Advice[resource.DeleteRequest, *resource.DeleteResponse]
	ResourceUpdateAopApi []Advice[resource.UpdateRequest, *resource.UpdateResponse]
}

func NewAopAdvices() *AopAdvices {
	return &AopAdvices{
		DataSourceReadAopApi: []Advice[datasource.ReadRequest, *datasource.ReadResponse]{LogAdvice[datasource.ReadRequest, *datasource.ReadResponse]{}},
		ResourceReadAopApi:   []Advice[resource.ReadRequest, *resource.ReadResponse]{LogAdvice[resource.ReadRequest, *resource.ReadResponse]{}},
		ResourceCreateAopApi: []Advice[resource.CreateRequest, *resource.CreateResponse]{LogAdvice[resource.CreateRequest, *resource.CreateResponse]{}},
		ResourceDeleteAopApi: []Advice[resource.DeleteRequest, *resource.DeleteResponse]{LogAdvice[resource.DeleteRequest, *resource.DeleteResponse]{}},
		ResourceUpdateAopApi: []Advice[resource.UpdateRequest, *resource.UpdateResponse]{LogAdvice[resource.UpdateRequest, *resource.UpdateResponse]{}},
	}
}

type LogAdvice[T any, R any] struct {
}

func (LogAdvice[T, R]) Around(ctx context.Context, req T, resp R, chain ResourceDecoratorChain[T, R]) error {
	startTime := time.Now().UnixMilli()
	reqStr, err := json.Marshal(req)
	id := uuid.NewString()
	apiName := fmt.Sprintf("%T", req)
	if err == nil {
		tflog.Info(ctx, "调用插件方法：", map[string]interface{}{
			"id":      id,
			"apiName": apiName,
			"request": string(reqStr),
		})
	}
	err = chain.Next(ctx, req, resp)
	endTime := time.Now().UnixMilli()
	useTime := endTime - startTime
	if err == nil {
		respStr, err := json.Marshal(resp)
		if err == nil {
			tflog.Info(ctx, "调用插件方法成功：", map[string]interface{}{
				"id":       id,
				"apiName":  apiName,
				"costTime": useTime,
				"response": string(respStr),
			})
		}
	} else {
		tflog.Info(ctx, "调用插件方法失败：", map[string]interface{}{
			"id":       id,
			"apiName":  apiName,
			"costTime": useTime,
			"error":    err.Error(),
		})
	}
	return err
}
