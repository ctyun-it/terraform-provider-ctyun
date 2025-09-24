package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctimage"
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

func (u ImageService) GetImageInfo(ctx context.Context, imageId, regionId string) (image ctimage.ImageDetailImagesResponse, err error) {
	resp, err := u.meta.Apis.CtImageApis.ImageDetailApi.Do(ctx, u.meta.Credential, &ctimage.ImageDetailRequest{
		ImageId:  imageId,
		RegionId: regionId,
	})
	if err != nil {
		return
	}
	if len(resp.Images) != 1 {
		err = fmt.Errorf("镜像 %s 不存在", imageId)
		return
	}
	image = resp.Images[0]
	return
}
