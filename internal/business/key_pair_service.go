package business

import (
	"context"
	"errors"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
)

type KeyPairService struct {
	meta *common.CtyunMetadata
}

func NewKeyPairService(meta *common.CtyunMetadata) *KeyPairService {
	return &KeyPairService{meta: meta}
}

func (u KeyPairService) GetKeyPairIDByName(ctx context.Context, keyPairName, regionId, projectId string) (string, error) {
	resp, err := u.meta.Apis.CtEcsApis.KeypairDetailApi.Do(ctx, u.meta.Credential, &ctecs.KeypairDetailRequest{
		RegionId:    regionId,
		KeyPairName: keyPairName,
		ProjectId:   projectId,
		PageNo:      1,
		PageSize:    1,
	})
	if err != nil {
		return "", err
	}
	if len(resp.Results) == 0 {
		return "", errors.New("密钥对不存在：" + keyPairName)
	}
	return resp.Results[0].KeyPairId, nil
}

func (u KeyPairService) GetKeyPairID(ctx context.Context, keyPairName, regionId, projectId string) (string, error) {
	resp, err := u.meta.Apis.SdkCtEcsApis.CtecsDetailsKeypairV41Api.Do(ctx, u.meta.SdkCredential, &ctecs2.CtecsDetailsKeypairV41Request{
		RegionID:    regionId,
		KeyPairName: keyPairName,
		ProjectID:   projectId,
		PageNo:      1,
		PageSize:    1,
	})
	if err != nil {
		return "", err
	}
	if len(resp.ReturnObj.Results) == 0 {
		return "", errors.New("密钥对不存在：" + keyPairName)
	}
	return resp.ReturnObj.Results[0].KeyPairID, nil
}
