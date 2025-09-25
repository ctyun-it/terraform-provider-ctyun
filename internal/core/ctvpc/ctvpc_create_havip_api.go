package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateHavipApi
/* 创建高可用虚IP
 */type CtvpcCreateHavipApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateHavipApi(client *core.CtyunClient) *CtvpcCreateHavipApi {
	return &CtvpcCreateHavipApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/havip/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateHavipApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateHavipRequest) (*CtvpcCreateHavipResponse, error) {
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
	var resp CtvpcCreateHavipResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateHavipRequest struct {
	ClientToken string  `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  `json:"regionID,omitempty"`    /*  资源池ID  */
	NetworkID   *string `json:"networkID,omitempty"`   /*  VPC的ID  */
	SubnetID    string  `json:"subnetID,omitempty"`    /*  子网ID  */
	IpAddress   *string `json:"ipAddress,omitempty"`   /*  ip地址  */
	VipType     *string `json:"vipType,omitempty"`     /*  虚拟IP的类型，v4-IPv4类型虚IP，v6-IPv6类型虚IP。默认为v4  */
}

type CtvpcCreateHavipResponse struct {
	StatusCode  int32                                `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                              `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                              `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                              `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcCreateHavipReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                              `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCreateHavipReturnObjResponse struct {
	Uuid *string `json:"uuid,omitempty"` /*  高可用虚IP的ID  */
	Ipv4 *string `json:"ipv4,omitempty"` /*  高可用虚IP的地址  */
	Ipv6 *string `json:"ipv6,omitempty"` /*  高可用虚IP的地址  */
}
