package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// ImageUpdateApi 修改私有镜像属性
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=5085&data=89
type ImageUpdateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewImageUpdateApi(client *ctyunsdk.CtyunClient) *ImageUpdateApi {
	return &ImageUpdateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/image/update",
		},
	}
}

func (this *ImageUpdateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *ImageUpdateRequest) (*ImageUpdateResponse, ctyunsdk.CtyunRequestError) {
	request := this.WithCredential(&credential)
	_, err := request.WriteJson(&imageUpdateRealRequest{
		RegionID:    req.RegionId,
		ImageID:     req.ImageId,
		BootMode:    req.BootMode,
		Description: req.Description,
		ImageName:   req.ImageName,
		MaximumRAM:  req.MaximumRam,
		MinimumRAM:  req.MinimumRam,
	})

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtimage, request)
	if err != nil {
		return nil, err
	}

	err = response.ParseByStandardModelWithCheck(nil)
	if err != nil {
		return nil, err
	}
	return &ImageUpdateResponse{}, nil
}

type imageUpdateRealRequest struct {
	ImageID     string `json:"imageID"`
	RegionID    string `json:"regionID"`
	BootMode    string `json:"bootMode,omitempty"`
	Description string `json:"description,omitempty"`
	ImageName   string `json:"imageName,omitempty"`
	MaximumRAM  int    `json:"maximumRAM,omitempty"`
	MinimumRAM  int    `json:"minimumRAM,omitempty"`
}

type ImageUpdateRequest struct {
	ImageId     string
	RegionId    string
	BootMode    string
	Description string
	ImageName   string
	MaximumRam  int
	MinimumRam  int
}

type ImageUpdateResponse struct {
}
