package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
)

type SecurityGroupService struct {
	meta *common.CtyunMetadata
}

func NewSecurityGroupService(meta *common.CtyunMetadata) *SecurityGroupService {
	return &SecurityGroupService{meta: meta}
}

func (s SecurityGroupService) MustExist(ctx context.Context, SecurityGroupId, regionId string) error {
	_, err := s.meta.Apis.CtVpcApis.SecurityGroupDescribeAttributeApi.Do(ctx, s.meta.Credential, &ctvpc.SecurityGroupDescribeAttributeRequest{
		RegionId:        regionId,
		SecurityGroupId: SecurityGroupId,
		Direction:       "all",
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiSecurityGroupNotFound {
			return fmt.Errorf("安全组 %s 不存在", SecurityGroupId)
		}
		return err
	}
	return nil
}

func (s SecurityGroupService) MustExistInVpc(ctx context.Context, vpcId, securityGroupId, regionId string) error {
	resp, err := s.meta.Apis.CtVpcApis.SecurityGroupDescribeAttributeApi.Do(ctx, s.meta.Credential, &ctvpc.SecurityGroupDescribeAttributeRequest{
		RegionId:        regionId,
		SecurityGroupId: securityGroupId,
		Direction:       "all",
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiSecurityGroupNotFound {
			return fmt.Errorf("安全组 %s 不存在", securityGroupId)
		}
		return err
	}
	if resp.VpcId != vpcId {
		return fmt.Errorf("安全组 %s 不属于 %s", vpcId)
	}
	return nil
}
