package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// ImageShareRejectApi 拒绝共享镜像
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=5227&data=89
type ImageShareRejectApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewimageShareRejectApi(client *ctyunsdk.CtyunClient) *ImageShareRejectApi {
	return &ImageShareRejectApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/image/shared-image/reject",
		},
	}
}

func (this *ImageShareRejectApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *ImageShareRejectRequest) (*ImageShareRejectResponse, ctyunsdk.CtyunRequestError) {
	request := this.WithCredential(&credential)
	_, err := request.WriteJson(&imageShareRejectRealRequest{
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
	return &ImageShareRejectResponse{}, nil
}

type imageShareRejectRealRequest struct {
	ImageID  string `json:"imageID"`
	RegionID string `json:"regionID"`
}

type ImageShareRejectRequest struct {
	ImageId  string
	RegionId string
}

type ImageShareRejectResponse struct {
}
