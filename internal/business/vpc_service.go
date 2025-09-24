package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
	"github.com/google/uuid"
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

func (v VpcService) GetVpcSubnet(ctx context.Context, vpcId, regionId, projectId string) (map[string]ctvpc.SubnetListSubnetsResponse, error) {
	params := &ctvpc.VpcQueryRequest{
		RegionId:    regionId,
		ProjectId:   projectId,
		ClientToken: uuid.NewString(),
		VpcId:       vpcId,
	}
	if projectId != "" {
		params.ProjectId = projectId
	}
	resp, err := v.meta.Apis.CtVpcApis.VpcQueryApi.Do(ctx, v.meta.Credential, params)
	if err != nil {
		if err.ErrorCode() == common.OpenapiVpcNotFound {
			return nil, fmt.Errorf("vpc %s 不存在", vpcId)
		}
		return nil, err
	}
	r, err := v.meta.Apis.CtVpcApis.SubnetListApi.Do(ctx, v.meta.Credential, &ctvpc.SubnetListRequest{
		RegionId:   regionId,
		VpcId:      vpcId,
		SubnetIds:  resp.SubnetIds,
		PageNumber: 1,
		PageSize:   100,
	})
	if err != nil {
		return nil, err
	}
	result := make(map[string]ctvpc.SubnetListSubnetsResponse)
	for _, subnet := range r.Subnets {
		result[subnet.SubnetId] = subnet
	}
	return result, nil
}
