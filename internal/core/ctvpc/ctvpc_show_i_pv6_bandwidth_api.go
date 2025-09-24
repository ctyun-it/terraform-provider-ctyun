package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowIPv6BandwidthApi
/* 查看 IPv6 带宽详情。
 */type CtvpcShowIPv6BandwidthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowIPv6BandwidthApi(client *core.CtyunClient) *CtvpcShowIPv6BandwidthApi {
	return &CtvpcShowIPv6BandwidthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ipv6_bandwidth/show",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowIPv6BandwidthApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowIPv6BandwidthRequest) (*CtvpcShowIPv6BandwidthResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("bandwidthID", req.BandwidthID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowIPv6BandwidthResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowIPv6BandwidthRequest struct {
	RegionID    string /*  资源池 ID  */
	BandwidthID string /*  IPv6 带宽 ID  */
}

type CtvpcShowIPv6BandwidthResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcShowIPv6BandwidthReturnObjResponse `json:"returnObj"`             /*  返回查询的共享带宽实例的详细信息。  */
	Error       *string                                    `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowIPv6BandwidthReturnObjResponse struct {
	Id            *string `json:"id,omitempty"`            /*  IPv6 带宽 ID  */
	Status        *string `json:"status,omitempty"`        /*  IPv6 带宽状态: ACTIVE（正常） / EXPIRED（过期） / FREEZING（冻结） /CREATEING（创建中）  */
	Name          *string `json:"name,omitempty"`          /*  IPv6 带宽名字  */
	Bandwidth     int32   `json:"bandwidth"`               /*  IPv6 带宽峰值 mbps  */
	ResourceSpec  *string `json:"resourceSpec,omitempty"`  /*  独享 / 共享  */
	PaymentType   *string `json:"paymentType,omitempty"`   /*  计费类型  */
	CreatedTime   *string `json:"createdTime,omitempty"`   /*  IPv6 带宽创建时间  */
	ExpiredTime   *string `json:"expiredTime,omitempty"`   /*  IPv6 带宽过期时间  */
	IpAddress     *string `json:"ipAddress,omitempty"`     /*  IP 地址  */
	Ipv6GatewayID *string `json:"ipv6GatewayID,omitempty"` /*  IPv6 网关 ID  */
}
