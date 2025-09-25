package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcSubnetEnableIPv6Api
/* 子网开启 IPv6
 */type CtvpcSubnetEnableIPv6Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcSubnetEnableIPv6Api(client *core.CtyunClient) *CtvpcSubnetEnableIPv6Api {
	return &CtvpcSubnetEnableIPv6Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/subnet-enable-ipv6",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcSubnetEnableIPv6Api) Do(ctx context.Context, credential core.Credential, req *CtvpcSubnetEnableIPv6Request) (*CtvpcSubnetEnableIPv6Response, error) {
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
	var resp CtvpcSubnetEnableIPv6Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcSubnetEnableIPv6Request struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	SubnetID string `json:"subnetID,omitempty"` /*  子网 的 ID  */
}

type CtvpcSubnetEnableIPv6Response struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
