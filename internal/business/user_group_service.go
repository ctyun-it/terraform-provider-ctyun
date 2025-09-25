package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctiam"
)

type UserGroupService struct {
	meta *common.CtyunMetadata
}

func NewUserGroupService(meta *common.CtyunMetadata) *UserGroupService {
	return &UserGroupService{meta: meta}
}

func (u UserGroupService) MustExist(ctx context.Context, userGroupId string) error {
	_, err := u.meta.Apis.CtIamApis.UserGroupGetApi.Do(ctx, u.meta.Credential, &ctiam.UserGroupGetRequest{
		GroupId: userGroupId,
	})
	if err != nil {
		if err.ErrorCode() == common.CtiamNoPermission {
			return fmt.Errorf("用户组 %s 不存在", userGroupId)
		}
		return err
	}
	return nil
}
