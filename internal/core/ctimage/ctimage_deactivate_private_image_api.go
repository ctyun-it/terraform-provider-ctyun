package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtimageDeactivatePrivateImageApi
/* 弃用 1 份云主机私有镜像<br />准备：<br />1. 在调用前需了解如何构造请求，可参见：如何调用 API - 构造请求<br />2. OpenAPI 请求需进行加密调用，可参见：如何调用 API - 认证鉴权<br />注意：在调用前，请您认真阅读此文档，包括但不限于参数描述中的“注意”部分
 */type CtimageDeactivatePrivateImageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtimageDeactivatePrivateImageApi(client *core.CtyunClient) *CtimageDeactivatePrivateImageApi {
	return &CtimageDeactivatePrivateImageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/image/deactivate",
			ContentType:  "application/json",
		},
	}
}

func (a *CtimageDeactivatePrivateImageApi) Do(ctx context.Context, credential core.Credential, req *CtimageDeactivatePrivateImageRequest) (*CtimageDeactivatePrivateImageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("imageID", req.ImageID)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtimageDeactivatePrivateImageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtimageDeactivatePrivateImageRequest struct {
	ImageID  string `json:"imageID,omitempty"`  /*  镜像 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89&isNormal=1&vid=83" target="_blank">查询可以使用的镜像资源</a>接口来查询您可使用的镜像资源，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4764&data=89&isNormal=1&vid=83" target="_blank">查询镜像详细信息</a>接口来查询 1 份镜像的详细信息，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=6764&data=89&isNormal=1&vid=83" target="_blank">查询私有镜像的共享列表</a>接口来查询 1 份私有镜像的共享列表。注意：所指定的镜像应是私有镜像的共享数量为 0、镜像状态为 active、镜像类型不为 iso_image 的私有镜像。此镜像在非多可用区资源池中还应是镜像类型不为 full_ecs_image 的镜像  */
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&isNormal=1&vid=81" target="_blank">资源池列表查询</a>接口来查询您可见的资源池的列表。注意：此接口仅支持具备弃用/取消弃用私有镜像的功能的资源池  */
}

type CtimageDeactivatePrivateImageResponse struct {
	StatusCode  int32                                           `json:"statusCode,omitempty"`  /*  状态码。取值范围（值：描述）：<br />800：成功，<br />900：失败  */
	Error       string                                          `json:"error,omitempty"`       /*  错误码（product.module.code 三段式码）  */
	ErrorCode   string                                          `json:"errorCode,omitempty"`   /*  同 error 参数  */
	Message     string                                          `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                          `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtimageDeactivatePrivateImageReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtimageDeactivatePrivateImageReturnObjResponse struct{}
