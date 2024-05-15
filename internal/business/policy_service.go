package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctiam"
	"terraform-provider-ctyun/internal/common"
)

type PolicyService struct {
	meta *common.CtyunMetadata
}

func NewPolicyService(meta *common.CtyunMetadata) *PolicyService {
	return &PolicyService{meta: meta}
}

func (v PolicyService) MustExist(ctx context.Context, policyId string) error {
	resp, err := v.meta.Apis.CtIamApis.PolicyGetApi.Do(ctx, v.meta.Credential, &ctiam.PolicyGetRequest{
		PolicyId: policyId,
	})
	if err != nil {
		return err
	}
	if resp.Id == "" {
		return fmt.Errorf("策略 %s 不存在", policyId)
	}
	return nil
}
