package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctiam"
	"terraform-provider-ctyun/internal/common"
)

type UserService struct {
	meta *common.CtyunMetadata
}

func NewUserService(meta *common.CtyunMetadata) *UserService {
	return &UserService{meta: meta}
}

func (u UserService) MustExist(ctx context.Context, userId string) error {
	_, err := u.meta.Apis.CtIamApis.UserGetApi.Do(ctx, u.meta.Credential, &ctiam.UserGetRequest{
		UserId: userId,
	})
	if err != nil {
		if err.ErrorCode() == common.CtiamNoPermission {
			return fmt.Errorf("用户 %s 不存在", userId)
		}
		return err
	}
	return nil
}
