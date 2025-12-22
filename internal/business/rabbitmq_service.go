package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/amqp"
)

type RabbitmqService struct {
	meta *common.CtyunMetadata
}

func NewRabbitmqService(meta *common.CtyunMetadata) *RabbitmqService {
	return &RabbitmqService{meta: meta}
}

func (c RabbitmqService) CheckVhostExist(ctx context.Context, vhost, instanceID, regionID string) (exist bool, err error) {
	params := &amqp.AmqpVhostQueryV3Request{
		RegionId:   regionID,
		ProdInstId: instanceID,
	}

	resp, err := c.meta.Apis.SdkAmqpApis.AmqpVhostQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.ReturnObj.Data == nil {
		err = common.InvalidReturnObjResultsError
		return
	}
	for _, v := range resp.ReturnObj.Data.Vhosts {
		if v == vhost {
			return true, err
		}
	}
	return
}

func (c RabbitmqService) CheckExchangeExist(ctx context.Context, name, vhost, instanceID, regionID string) (exist bool, err error) {
	params := &amqp.AmqpExchangeQueryV3Request{
		RegionId:   regionID,
		ProdInstId: instanceID,
		Vhost:      vhost,
		Name:       name,
	}

	resp, err := c.meta.Apis.SdkAmqpApis.AmqpExchangeQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		if resp.Message == "交换器不存在" {
			exist = false
			err = nil
		} else {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
		}
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.ReturnObj.Data == nil {
		err = common.InvalidReturnObjResultsError
		return
	} else if len(resp.ReturnObj.Data.Items) == 0 {
		exist = false
		return
	}
	for _, item := range resp.ReturnObj.Data.Items {
		if item.Name == name {
			exist = true
			return
		}
	}
	return
}

func (c RabbitmqService) CheckQueueExist(ctx context.Context, name, vhost, instanceID, regionID string) (exist bool, err error) {
	params := &amqp.AmqpQueueQueryV3Request{
		RegionId:   regionID,
		ProdInstId: instanceID,
		Vhost:      vhost,
		Name:       name,
	}

	resp, err := c.meta.Apis.SdkAmqpApis.AmqpQueueQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.ReturnObj.Data == nil {
		err = common.InvalidReturnObjResultsError
		return
	} else if len(resp.ReturnObj.Data.Items) == 0 {
		exist = false
		return
	}
	for _, item := range resp.ReturnObj.Data.Items {
		if item.Name == name {
			exist = true
			return
		}
	}
	return
}
