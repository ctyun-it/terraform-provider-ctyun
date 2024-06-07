package ctimage

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// ImageShareDeleteApi 取消共享私有镜像
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=5229&data=89
type ImageShareDeleteApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewImageShareDeleteApi(client *ctyunsdk.CtyunClient) *ImageShareDeleteApi {
	return &ImageShareDeleteApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/image/shared-image/delete",
		},
	}
}

func (this *ImageShareDeleteApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *ImageShareDeleteRequest) (*ImageShareDeleteResponse, ctyunsdk.CtyunRequestError) {
	request := this.WithCredential(&credential)
	_, err := request.WriteJson(&imageShareDeleteRealRequest{
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
	return &ImageShareDeleteResponse{}, nil
}

type imageShareDeleteRealRequest struct {
	DestinationUser string `json:"destinationUser"`
	ImageID         string `json:"imageID"`
	RegionID        string `json:"regionID"`
}

type ImageShareDeleteRequest struct {
	DestinationUser string
	ImageId         string
	RegionId        string
}

type ImageShareDeleteResponse struct {
}
