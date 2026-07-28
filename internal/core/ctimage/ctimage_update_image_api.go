package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtimageUpdateImageApi
/* 修改 1 份私有镜像的名称、描述信息等属性<br />准备：<br />1. 在调用前需了解如何构造请求，可参见：如何调用 API - 构造请求<br />2. OpenAPI 请求需进行加密调用，可参见：如何调用 API - 认证鉴权<br />注意：在调用前，请您认真阅读此文档，包括但不限于参数描述中的“注意”部分
 */type CtimageUpdateImageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtimageUpdateImageApi(client *core.CtyunClient) *CtimageUpdateImageApi {
	return &CtimageUpdateImageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/image/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtimageUpdateImageApi) Do(ctx context.Context, credential core.Credential, req *CtimageUpdateImageRequest) (*CtimageUpdateImageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtimageUpdateImageRequest
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
	var resp CtimageUpdateImageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtimageUpdateImageRequest struct {
	ImageID     string `json:"imageID,omitempty"`     /*  镜像 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89&isNormal=1&vid=83" target="_blank">查询可以使用的镜像资源</a>接口来查询您可使用的镜像资源，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4764&data=89&isNormal=1&vid=83" target="_blank">查询镜像详细信息</a>接口来查询 1 份镜像的详细信息。注意：所指定的镜像应是镜像状态为 active 的私有镜像  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&isNormal=1&vid=81" target="_blank">资源池列表查询</a>接口来查询您可见的资源池的列表  */
	BootMode    string `json:"bootMode,omitempty"`    /*  x86_64 架构非数据盘镜像的启动方式。默认传入的 imageID 参数所指定的镜像的启动方式。取值范围（值：描述）：<br />bios：BIOS 启动方式，<br />uefi：UEFI 启动方式<br />注意：仅在以下条件均满足时生效<br />1. 传入的 imageID 参数所指定的镜像是系统架构为 x86_64、镜像类型为空（系统盘镜像）或 full_ecs_image 的云主机镜像<br />2. 传入的 regionID 参数所指定的资源池具备 x86_64 UEFI 启动方式功能  */
	Description string `json:"description,omitempty"` /*  描述信息。注意：<br />1. 长度为 1~128 个字符，不能以空格开头或结尾<br />2. 若此参数值等效为空，则描述信息将被清空  */
	ImageName   string `json:"imageName,omitempty"`   /*  镜像名称。默认传入的 imageID 参数所指定的镜像的镜像名称。注意：<br />1. 长度为 2~32 个字符，只能由数字、字母、- 组成，不能以数字、- 开头，且不能以 - 结尾<br />2. 不能与您已有的私有镜像（传入的 imageID 参数所指定的镜像除外）的名称重复。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89&isNormal=1&vid=83" target="_blank">查询可以使用的镜像资源</a>接口来查询您已有的私有镜像  */
	MaximumRAM  int32  `json:"maximumRAM,omitempty"`  /*  最大内存。单位为 GiB。默认传入的 imageID 参数所指定的镜像的最大内存。取值范围：0（不限制）/1/2/4/8/16/32/64/128/256/512。注意：<br />1. 若同时限制最大和最小内存，则最大内存不能小于最小内存<br />2. 仅在传入的 imageID 参数所指定的镜像是镜像类型为空（系统盘镜像）的云主机镜像且传入的 regionID 参数所指定的资源池具备设置最大/最小内存的功能时生效  */
	MinimumRAM  int32  `json:"minimumRAM,omitempty"`  /*  最小内存。单位为 GiB。默认传入的 imageID 参数所指定的镜像的最小内存。取值范围：0（不限制）/1/2/4/8/16/32/64/128/256/512。注意：<br />1. 若同时限制最大和最小内存，则最大内存不能小于最小内存<br />2. 仅在传入的 imageID 参数所指定的镜像是镜像类型为空（系统盘镜像）的云主机镜像且传入的 regionID 参数所指定的资源池具备设置最大/最小内存的功能时生效  */
	SupportXSSD *bool  `json:"supportXSSD"`           /*  用于表示是否支持 XSSD 类型盘的标识。默认传入的 imageID 参数所指定的镜像的对应标识值。注意：仅在传入的 imageID 参数所指定的镜像是磁盘格式为 raw、镜像类型为空（系统盘镜像）或 full_ecs_image 的云主机镜像且传入的 regionID 参数所指定的资源池具备设置用于表示是否支持 XSSD 类型盘的标识的功能时生效  */
}

type CtimageUpdateImageResponse struct {
	StatusCode  int32                                `json:"statusCode,omitempty"`  /*  状态码。取值范围（值：描述）：<br />800：成功，<br />900：失败  */
	Error       string                               `json:"error,omitempty"`       /*  错误码（product.module.code 三段式码）  */
	ErrorCode   string                               `json:"errorCode,omitempty"`   /*  同 error 参数  */
	Message     string                               `json:"message,omitempty"`     /*  英文描述信息   */
	Description string                               `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtimageUpdateImageReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtimageUpdateImageReturnObjResponse struct{}
