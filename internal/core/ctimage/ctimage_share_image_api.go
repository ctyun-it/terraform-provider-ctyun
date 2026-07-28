package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtimageShareImageApi
/* 与指定用户共享 1 份私有镜像<br />准备：<br />1. 在调用前需了解如何构造请求，可参见：如何调用 API - 构造请求<br />2. OpenAPI 请求需进行加密调用，可参见：如何调用 API - 认证鉴权<br />注意：在调用前，请您认真阅读此文档，包括但不限于参数描述中的“注意”部分
 */type CtimageShareImageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtimageShareImageApi(client *core.CtyunClient) *CtimageShareImageApi {
	return &CtimageShareImageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/image/shared-image/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtimageShareImageApi) Do(ctx context.Context, credential core.Credential, req *CtimageShareImageRequest) (*CtimageShareImageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtimageShareImageRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtimageShareImageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtimageShareImageRequest struct {
	DestinationAccountID string `json:"destinationAccountID,omitempty"` /*  共享镜像接受者的账号 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=13017&data=114&isNormal=1&vid=107" target="_blank">分页查询用户</a>接口来查询用户信息。注意：所指定的共享镜像接受者不能是传入的 imageID 参数所指定的镜像的拥有者，也不能在此镜像的共享列表中。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=6764&data=89&isNormal=1&vid=83" target="_blank">查询私有镜像的共享列表</a>接口来查询 1 份私有镜像的共享列表  */
	ImageID              string `json:"imageID,omitempty"`              /*  镜像 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89&isNormal=1&vid=83" target="_blank">查询可以使用的镜像资源</a>接口来查询您可使用的镜像资源，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4764&data=89&isNormal=1&vid=83" target="_blank">查询镜像详细信息</a>接口来查询 1 份镜像的详细信息。注意：<br />1. 所指定的镜像应是镜像状态为 active、镜像类型不为 iso_image 的私有镜像。此镜像在非多可用区资源池中还应是镜像类型不为 full_ecs_image 的镜像<br />2. 所指定的镜像的共享镜像接受者配额余量足够，即可添加的共享镜像接受者数量未达上线  */
	RegionID             string `json:"regionID,omitempty"`             /*  资源池 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&isNormal=1&vid=81" target="_blank">资源池列表查询</a>接口来查询您可见的资源池的列表。注意：此接口仅支持具备共享/取消共享私有镜像的功能的资源池  */
}

type CtimageShareImageResponse struct {
	StatusCode  int32                               `json:"statusCode,omitempty"`  /*  状态码。取值范围（值：描述）：<br />800：成功，<br />900：失败  */
	Error       string                              `json:"error,omitempty"`       /*  错误码（product.module.code 三段式码）  */
	ErrorCode   string                              `json:"errorCode,omitempty"`   /*  同 error 参数  */
	Message     string                              `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                              `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtimageShareImageReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtimageShareImageReturnObjResponse struct{}
