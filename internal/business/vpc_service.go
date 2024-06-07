package business

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
)

type VpcService struct {
	meta *common.CtyunMetadata
}

func NewVpcService(meta *common.CtyunMetadata) *VpcService {
	return &VpcService{meta: meta}
}

func (v VpcService) MustExist(ctx context.Context, vpcId, regionId, projectId string) error {
	_, err := v.meta.Apis.CtVpcApis.VpcQueryApi.Do(ctx, v.meta.Credential, &ctvpc.VpcQueryRequest{
		RegionId:    regionId,
		ProjectId:   projectId,
		ClientToken: uuid.NewString(),
		VpcId:       vpcId,
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiVpcNotFound {
			return fmt.Errorf("vpc %s 不存在", vpcId)
		}
		return err
	}
	return nil
}
