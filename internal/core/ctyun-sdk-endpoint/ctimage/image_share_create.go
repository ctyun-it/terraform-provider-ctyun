package ctimage

import (
	"context"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"

	"net/http"
)

// ImageShareCreateApi 共享私有镜像
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=5114&data=89
type ImageShareCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewImageShareCreateApi(client *ctyunsdk.CtyunClient) *ImageShareCreateApi {
	return &ImageShareCreateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/image/shared-image/create",
		},
	}
}

func (this *ImageShareCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *ImageShareCreateRequest) (*ImageShareCreateResponse, ctyunsdk.CtyunRequestError) {
	request := this.WithCredential(&credential)
	_, err := request.WriteJson(&imageShareCreateRealRequest{
		DestinationUser: req.DestinationUser,
		ImageID:         req.ImageId,
		RegionID:        req.RegionId,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtimage, request)
	if err != nil {
		return nil, err
	}

	err = response.ParseByStandardModelWithCheck(nil)
	if err != nil {
		return nil, err
	}
	return &ImageShareCreateResponse{}, nil
}

type imageShareCreateRealRequest struct {
	DestinationUser string `json:"destinationUser"`
	ImageID         string `json:"imageID"`
	RegionID        string `json:"regionID"`
}

type ImageShareCreateRequest struct {
	DestinationUser string
	ImageId         string
	RegionId        string
}

type ImageShareCreateResponse struct {
}
