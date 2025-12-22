package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtimageListSupportedDestCrrRegionsApi
/* 查询所指定的源资源池的私有镜像支持复制到的目标资源池。<br />准备：<br />1. 在调用前需了解如何构造请求，可参见：如何调用 API - 构造请求。<br />2. OpenAPI 请求需进行加密调用，可参见：如何调用 API - 认证鉴权。<br />注意：在调用前，请您认真阅读此文档，包括但不限于参数描述中的“注意”部分。
 */type CtimageListSupportedDestCrrRegionsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtimageListSupportedDestCrrRegionsApi(client *core.CtyunClient) *CtimageListSupportedDestCrrRegionsApi {
	return &CtimageListSupportedDestCrrRegionsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/image/list-supported-dest-crr-regions",
			ContentType:  "application/json",
		},
	}
}

func (a *CtimageListSupportedDestCrrRegionsApi) Do(ctx context.Context, credential core.Credential, req *CtimageListSupportedDestCrrRegionsRequest) (*CtimageListSupportedDestCrrRegionsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtimageListSupportedDestCrrRegionsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtimageListSupportedDestCrrRegionsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  源资源池 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&isNormal=1&vid=81" target="_blank">资源池列表查询</a>接口来查询您可见的资源池的列表。注意：此接口仅支持具备镜像跨域复制功能的多可用区资源池。  */
}

type CtimageListSupportedDestCrrRegionsResponse struct {
	StatusCode  int32                                                `json:"statusCode,omitempty"`  /*  状态码。取值范围（值：描述）：<br />800：成功。<br />900：失败。  */
	Error       string                                               `json:"error,omitempty"`       /*  错误码（product.module.code 三段式码）。  */
	ErrorCode   string                                               `json:"errorCode,omitempty"`   /*  同 error 参数。  */
	Message     string                                               `json:"message,omitempty"`     /*  响应状态描述（一般为英文）。  */
	Description string                                               `json:"description,omitempty"` /*  响应状态描述（一般为中文）。  */
	ReturnObj   *CtimageListSupportedDestCrrRegionsReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据。  */
}

type CtimageListSupportedDestCrrRegionsReturnObjResponse struct {
	Regions []string `json:"regions"` /*  目标资源池 ID 列表。  */
}
