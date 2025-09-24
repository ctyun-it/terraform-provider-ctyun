package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"

	"net/http"
)

// ImageShareAcceptApi 接受共享镜像
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=5225&data=89
type ImageShareAcceptApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewimageShareAcceptApi(client *ctyunsdk.CtyunClient) *ImageShareAcceptApi {
	return &ImageShareAcceptApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/image/shared-image/accept",
		},
	}
}

func (this *ImageShareAcceptApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *ImageShareAcceptRequest) (*ImageShareAcceptResponse, ctyunsdk.CtyunRequestError) {
	request := this.WithCredential(&credential)
	_, err := request.WriteJson(&imageShareAcceptRealRequest{
		ImageID:  req.ImageId,
		RegionID: req.RegionId,
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
	return &ImageShareAcceptResponse{}, nil
}

type imageShareAcceptRealRequest struct {
	ImageID  string `json:"imageID"`
	RegionID string `json:"regionID"`
}

type ImageShareAcceptRequest struct {
	ImageId  string
	RegionId string
}

type ImageShareAcceptResponse struct {
}
