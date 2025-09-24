package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcGwlbDisableIpv6Api
/* gwlb关闭ipv6
 */type CtvpcGwlbDisableIpv6Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGwlbDisableIpv6Api(client *core.CtyunClient) *CtvpcGwlbDisableIpv6Api {
	return &CtvpcGwlbDisableIpv6Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/gwlb/disable-ipv6",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGwlbDisableIpv6Api) Do(ctx context.Context, credential core.Credential, req *CtvpcGwlbDisableIpv6Request) (*CtvpcGwlbDisableIpv6Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcGwlbDisableIpv6Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGwlbDisableIpv6Request struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	GwLbID   string `json:"gwLbID,omitempty"`   /*  网关负载均衡ID  */
}

type CtvpcGwlbDisableIpv6Response struct {
	StatusCode  int32                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGwlbDisableIpv6ReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcGwlbDisableIpv6ReturnObjResponse struct {
	GwLbID *string `json:"gwLbID,omitempty"` /*  网关负载均衡 ID  */
}
