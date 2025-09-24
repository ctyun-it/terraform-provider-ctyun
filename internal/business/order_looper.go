package business

import (
	"context"
	"errors"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
	"strconv"
	"time"
)

type OrderLooper struct {
	api *ctecs.EcsOrderQueryUuidApi
}

func NewOrderLooper(api *ctecs.EcsOrderQueryUuidApi) *OrderLooper {
	return &OrderLooper{
		api: api,
	}
}

// OrderLoop 轮询操作
func (o *OrderLooper) OrderLoop(ctx context.Context, credential ctyunsdk.Credential, masterOrderId string, loopCount ...int) (*LoopOrderResponse, error) {
	var resp *LoopOrderResponse
	var respError error
	c := 60
	if len(loopCount) > 0 {
		c = loopCount[0]
	}
	var cnt int
	retryer, _ := NewRetryer(time.Second*5, c)
	result := retryer.Start(
		func(currentTime int) bool {
			detail, err := o.api.Do(ctx, credential, &ctecs.EcsOrderQueryUuidRequest{
				MasterOrderId: masterOrderId,
			})
			if err != nil {
				respError = err
				return false
			}
			status, err2 := strconv.Atoi(detail.OrderStatus)
			if err2 != nil {
				respError = err2
				return false
			}

			switch status {
			case OrderStatusOpening:
				// 开通中状态
				return true
			case OrderStatusFinish:
				if len(detail.InstanceIDList) == 0 && cnt < 3 {
					cnt++
					return true
				}
				// 开通完成状态
				resp = &LoopOrderResponse{
					Uuid:          detail.InstanceIDList,
					masterOrderId: masterOrderId,
				}
				return false
			default:
				// 其他状态
				sta := OrderStatusName[status]
				respError = errors.New("轮询订购订单状态失败，轮询到的订单状态为：" + sta)
				return false
			}
		},
	)
	if result.ReturnReason == ReachMaxLoopTime {
		// 这里出来的全都是异常的
		return nil, errors.New("轮询订购订单状态失败，已超过最大轮询次数，订单号：" + masterOrderId)
	}
	return resp, respError
}

// RefundLoop 轮询操作
func (o *OrderLooper) RefundLoop(ctx context.Context, credential ctyunsdk.Credential, masterOrderId string) error {
	var respError error
	retryer, _ := NewRetryer(time.Second*5, 60)
	result := retryer.Start(
		func(currentTime int) bool {
			detail, err := o.api.Do(ctx, credential, &ctecs.EcsOrderQueryUuidRequest{
				MasterOrderId: masterOrderId,
			})
			if err != nil {
				respError = err
				return false
			}
			status, err2 := strconv.Atoi(detail.OrderStatus)
			if err2 != nil {
				respError = err2
				return false
			}

			switch status {
			case OrderStatusOpening:
				// 开通中状态
				// 继续执行循环里面的内容
				return true
			case OrderStatusFinish:
				// 开通完成状态
				return false
			default:
				// 其他状态
				sta := OrderStatusName[status]
				respError = errors.New("轮询订购订单状态失败，轮询到的订单状态为：" + sta)
				return false
			}
		},
	)
	if result.ReturnReason == ReachMaxLoopTime {
		// 这里出来的全都是异常的
		return errors.New("轮询退订单状态失败，已超过最大轮询次数，订单号：" + masterOrderId)
	}
	return respError
}

type LoopOrderResponse struct {
	Uuid          []string
	masterOrderId string // 主订单id
}
