package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtimageCopyImageApi
/* 将 1 份云主机私有镜像经过加密/解密后复制到同一资源池下。<br />准备：<br />1. 在调用前需了解如何构造请求，可参见：如何调用 API - 构造请求。<br />2. OpenAPI 请求需进行加密调用，可参见：如何调用 API - 认证鉴权。<br />注意：在调用前，请您认真阅读此文档，包括但不限于参数描述中的“注意”部分。
 */type CtimageCopyImageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtimageCopyImageApi(client *core.CtyunClient) *CtimageCopyImageApi {
	return &CtimageCopyImageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/image/copy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtimageCopyImageApi) Do(ctx context.Context, credential core.Credential, req *CtimageCopyImageRequest) (*CtimageCopyImageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtimageCopyImageRequest
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
	var resp CtimageCopyImageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtimageCopyImageRequest struct {
	ImageID     string                           `json:"imageID,omitempty"`     /*  源镜像 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89&isNormal=1&vid=83" target="_blank">查询可以使用的镜像资源</a>接口来查询您可使用的镜像资源，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4764&data=89&isNormal=1&vid=83" target="_blank">查询镜像详细信息</a>接口来查询 1 份镜像的详细信息。注意：<br />1. 所指定的镜像应是系统盘镜像或数据盘镜像，且镜像状态为 active 的云主机私有镜像。  */
	ImageName   string                           `json:"imageName,omitempty"`   /*  目标镜像名称。注意：<br />1. 长度为 2~32 个字符，只能由数字、字母、- 组成，不能以数字、- 开头，且不能以 - 结尾。<br />2. 不能与已有的私有镜像的名称重复。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89&isNormal=1&vid=83" target="_blank">查询可以使用的镜像资源</a>接口来查询您在资源池已有的私有镜像。  */
	RegionID    string                           `json:"regionID,omitempty"`    /*  源资源池 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&isNormal=1&vid=81" target="_blank">资源池列表查询</a>接口来查询您可见的资源池的列表。注意：此接口仅支持具备私有镜像加密功能的多可用区资源池。  */
	CmkID       string                           `json:"cmkID,omitempty"`       /*  密钥ID。私有镜像进行同域加密复制若未传cmkID，会使用默认密钥加密。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8342&data=87&isNormal=1&vid=81" target="_blank">查询可以使用的密钥资源</a>接口来查询您在资源池已有的密钥。  */
	Description string                           `json:"description,omitempty"` /*  目标镜像描述信息。注意：长度为 1~128 个字符，不能以空格开头或结尾。  */
	Labels      []*CtimageCopyImageLabelsRequest `json:"labels"`                /*  目标镜像标签列表。注意：<br />1. 列表中最多 10 个标签。<br />2. 标签键不可重复。<br />3. 单个标签键或值应满足长度为 1~32 个字符，不能换行，且不能以空格开头或结尾。  */
	ProjectID   string                           `json:"projectID,omitempty"`   /*  目标镜像企业项目 ID。默认 0（即 default 企业项目）。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=7246&data=114&isNormal=1&vid=107" target="_blank">查询企业项目列表</a>接口来查询您可以使用的企业项目 ID。  */
}

type CtimageCopyImageLabelsRequest struct {
	LabelKey   string `json:"labelKey,omitempty"`   /*  标签键。  */
	LabelValue string `json:"labelValue,omitempty"` /*  标签值。  */
}

type CtimageCopyImageResponse struct {
	StatusCode  int32                              `json:"statusCode,omitempty"`  /*  状态码。取值范围（值：描述）：<br />800：成功。<br />900：失败。  */
	Error       string                             `json:"error,omitempty"`       /*  错误码（product.module.code 三段式码）。  */
	ErrorCode   string                             `json:"errorCode,omitempty"`   /*  同 error 参数。  */
	Message     string                             `json:"message,omitempty"`     /*  响应状态描述（一般为英文）。  */
	Description string                             `json:"description,omitempty"` /*  响应状态描述（一般为中文）。  */
	ReturnObj   *CtimageCopyImageReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据。  */
}

type CtimageCopyImageReturnObjResponse struct{}
