package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetEndpointApi
/* 查询可访问资源池的endpoint。
 */type ZosGetEndpointApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetEndpointApi(client *core.CtyunClient) *ZosGetEndpointApi {
	return &ZosGetEndpointApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-endpoint",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetEndpointApi) Do(ctx context.Context, credential core.Credential, req *ZosGetEndpointRequest) (*ZosGetEndpointResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetEndpointResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetEndpointRequest struct {
	RegionID string /*  区域ID  */
}

type ZosGetEndpointResponse struct {
	StatusCode  int64                            `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                           `json:"message,omitempty"`     /*  状态描述  */
	Description string                           `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetEndpointReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                           `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                           `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetEndpointReturnObjResponse struct {
	Ipv6Endpoint     []string `json:"ipv6Endpoint"`     /*  内网 ipv6 列表，仅当没有内网域名时返回此参数，且无可用 ipv6 地址时为空数组   */
	IntranetEndpoint []string `json:"intranetEndpoint"` /*  内网 endpoint 列表，无可用地址时为空。有内网域名时使用域名，否则为 IPv4 地址  */
	InternetEndpoint []string `json:"internetEndpoint"` /*  外网 endpoint 列表，无可用地址时为空  */
}
