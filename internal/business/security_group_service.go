package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctvpc"
	"terraform-provider-ctyun/internal/common"
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
