package business

import (
	"context"
	"errors"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
)

type KeyPairService struct {
	meta *common.CtyunMetadata
}

func NewKeyPairService(meta *common.CtyunMetadata) *KeyPairService {
	return &KeyPairService{meta: meta}
}

func (u KeyPairService) MustExist(ctx context.Context, keyPairName, regionId, projectId string) error {
	resp, err := u.meta.Apis.CtEcsApis.KeypairDetailApi.Do(ctx, u.meta.Credential, &ctecs.KeypairDetailRequest{
		RegionId:    regionId,
		KeyPairName: keyPairName,
		ProjectId:   projectId,
		PageNo:      1,
		PageSize:    1,
	})
	if err != nil {
		return err
	}
	if len(resp.Results) == 0 {
		return errors.New("密钥对不存在：" + keyPairName)
	}
	return nil
}
