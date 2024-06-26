package business

import (
	"context"
	"fmt"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctimage"
)

type ImageService struct {
	meta *common.CtyunMetadata
}

func NewImageService(meta *common.CtyunMetadata) *ImageService {
	return &ImageService{meta: meta}
}

func (u ImageService) MustExist(ctx context.Context, imageId, regionId string) error {
	resp, err := u.meta.Apis.CtImageApis.ImageDetailApi.Do(ctx, u.meta.Credential, &ctimage.ImageDetailRequest{
		ImageId:  imageId,
		RegionId: regionId,
	})
	if err != nil {
		return err
	}
	if len(resp.Images) != 1 {
		return fmt.Errorf("镜像 %s 不存在", imageId)
	}
	return nil
}
