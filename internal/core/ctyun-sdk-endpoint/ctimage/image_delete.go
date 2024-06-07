package ctimage

import (
	"context"

	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// ImageDeleteApi 删除私有镜像
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4766&data=89
type ImageDeleteApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewImageDeleteApi(client *ctyunsdk.CtyunClient) *ImageDeleteApi {
	return &ImageDeleteApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/image/delete",
		},
	}
}

func (this *ImageDeleteApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *ImageDeleteRequest) (*ImageDeleteResponse, ctyunsdk.CtyunRequestError) {
	request := this.WithCredential(&credential)
	_, err := request.WriteJson(&imageDeleteRealRequest{
		RegionID: req.RegionId,
		ImageID:  req.ImageId,
	})

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtimage, request)
	if err != nil {
		return nil, err
	}

	err = response.ParseByStandardModelWithCheck(nil)
	if err != nil {
		return nil, err
	}
	return &ImageDeleteResponse{}, nil
}

type imageDeleteRealRequest struct {
	RegionID string `json:"regionID"`
	ImageID  string `json:"imageID"`
}

type ImageDeleteRequest struct {
	RegionId string // 资源池id
	ImageId  string // 镜像id
}

type ImageDeleteResponse struct {
}
