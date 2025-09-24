package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowIPv6Api
/* 查询 IPv6 详情
 */type CtvpcShowIPv6Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowIPv6Api(client *core.CtyunClient) *CtvpcShowIPv6Api {
	return &CtvpcShowIPv6Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ipv6/ipv6-show",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowIPv6Api) Do(ctx context.Context, credential core.Credential, req *CtvpcShowIPv6Request) (*CtvpcShowIPv6Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("ipv6ID", req.Ipv6ID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowIPv6Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowIPv6Request struct {
	RegionID string /*  资源池ID  */
	Ipv6ID   string /*  ipv6 id  */
}

type CtvpcShowIPv6Response struct {
	StatusCode  int32                             `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                           `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                           `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                           `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcShowIPv6ReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                           `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowIPv6ReturnObjResponse struct {
	ID              *string `json:"ID,omitempty"`            /*  ipv6 id  */
	VpcID           *string `json:"vpcID,omitempty"`         /*  vpc id  */
	SubnetID        *string `json:"subnetID,omitempty"`      /*  子网id  */
	IpAddress       *string `json:"ipAddress,omitempty"`     /*  ipv6地址  */
	AssociationID   *string `json:"associationID,omitempty"` /*  绑定实例的id  */
	AssociationType int32   `json:"associationType"`         /*  绑定实例类型 0: port, 1:ha_vip, 2:elb  */
	BandwidthID     *string `json:"bandwidthID,omitempty"`   /*  公网带宽id  */
	CreatedAt       *string `json:"createdAt,omitempty"`     /*  创建时间  */
	UpdatedAt       *string `json:"updatedAt,omitempty"`     /*  更新时间  */
}
