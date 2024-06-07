package business

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctebs"
)

type EbsService struct {
	meta *common.CtyunMetadata
}

func NewEbsService(meta *common.CtyunMetadata) *EbsService {
	return &EbsService{meta: meta}
}

func (u EbsService) MustExist(ctx context.Context, id, regionId string) error {
	_, err := u.meta.Apis.CtEbsApis.EbsShowApi.Do(ctx, u.meta.Credential, &ctebs.EbsShowRequest{
		RegionId: regionId,
		DiskId:   id,
	})
	if err != nil {
		if err.ErrorCode() == common.EbsEbsInfoDataDamaged {
			return fmt.Errorf("云硬盘 %s 不存在", id)
		}
		return err
	}
	return nil
}

// UpdateSize 更新云硬盘大小
func (c *EbsService) UpdateSize(ctx context.Context, diskId string, regionId string, oldSize, newSize int) error {
	// 云硬盘升配 云硬盘不支持降配
	if newSize == oldSize {
		return nil
	} else if newSize < oldSize {
		return errors.New("云硬盘不支持缩容")
	}

	resp, err := c.meta.Apis.CtEbsApis.EbsChangeSizeApi.Do(ctx, c.meta.Credential, &ctebs.EbsChangeSizeRequest{
		RegionId:    regionId,
		DiskId:      diskId,
		DiskSize:    newSize,
		ClientToken: uuid.NewString(),
	})

	var masterOrderId string
	if err == nil {
		masterOrderId = resp.MasterOrderId
	} else {
		if err.ErrorCode() != common.EbsOrderInProgress {
			return err
		}
		id, err := c.getMasterOrderIdIfOrderInProgress(err)
		if err != nil {
			return err
		}
		masterOrderId = id
	}

	helper := NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	_, err2 := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId)
	return err2
}

// getMasterOrderIdIfOrderInProgress 获取masterOrderId
func (c *EbsService) getMasterOrderIdIfOrderInProgress(err ctyunsdk.CtyunRequestError) (string, error) {
	resp := struct {
		MasterOrderId string `json:"masterOrderID"`
		MasterOrderNo string `json:"masterOrderNO"`
	}{}
	if err.CtyunResponse() == nil {
		return "", err
	}
	_, err = err.CtyunResponse().ParseByStandardModel(&resp)
	if err != nil {
		return "", err
	}
	return resp.MasterOrderId, err
}
