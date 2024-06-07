package business

import (
	"context"
	"fmt"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctiam"
)

type EnterpriseProjectService struct {
	meta *common.CtyunMetadata
}

func NewEnterpriseProjectService(meta *common.CtyunMetadata) *EnterpriseProjectService {
	return &EnterpriseProjectService{meta: meta}
}

func (u EnterpriseProjectService) MustExist(ctx context.Context, enterpriseProjectId string) error {
	resp, err := u.meta.Apis.CtIamApis.EnterpriseProjectGetApi.Do(ctx, u.meta.Credential, &ctiam.EnterpriseProjectGetRequest{
		Id: enterpriseProjectId,
	})
	if err != nil {
		return err
	}
	if resp.Id == "" {
		return fmt.Errorf("企业项目 %s 不存在", enterpriseProjectId)
	}
	return nil
}
