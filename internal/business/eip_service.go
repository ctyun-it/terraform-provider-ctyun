package business

import (
	"context"
	"fmt"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
)

type EipService struct {
	meta *common.CtyunMetadata
}

func NewEipService(meta *common.CtyunMetadata) *EipService {
	return &EipService{meta: meta}
}

func (u EipService) MustExist(ctx context.Context, id, regionId string) error {
	_, err := u.meta.Apis.CtVpcApis.EipShowApi.Do(ctx, u.meta.Credential, &ctvpc.EipShowRequest{
		RegionId: regionId,
		EipId:    id,
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiEipNotFound {
			return fmt.Errorf("eip %s 不存在", id)
		}
		return err
	}
	return nil
}
