package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
)

type BandwidthService struct {
	meta *common.CtyunMetadata
}

func NewBandwidthService(meta *common.CtyunMetadata) *BandwidthService {
	return &BandwidthService{meta: meta}
}

func (u BandwidthService) MustExist(ctx context.Context, id, regionId, projectId string) error {
	_, err := u.meta.Apis.CtVpcApis.BandwidthDescribeApi.Do(ctx, u.meta.Credential, &ctvpc.BandwidthDescribeRequest{
		BandwidthId: id,
		RegionId:    regionId,
		ProjectId:   projectId,
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiSharedbandwidthNotFound {
			return fmt.Errorf("带宽 %s 不存在", id)
		}
		return err
	}
	return nil
}
