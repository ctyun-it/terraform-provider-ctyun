package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowIPv6GatewayApi
/* 调用此接口可查询 IPv6 网关详情。
 */type CtvpcShowIPv6GatewayApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowIPv6GatewayApi(client *core.CtyunClient) *CtvpcShowIPv6GatewayApi {
	return &CtvpcShowIPv6GatewayApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/get-ipv6-gateway-attribute",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowIPv6GatewayApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowIPv6GatewayRequest) (*CtvpcShowIPv6GatewayResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ProjectID != nil {
		ctReq.AddParam("projectID", *req.ProjectID)
	}
	ctReq.AddParam("ipv6GatewayID", req.Ipv6GatewayID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowIPv6GatewayResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowIPv6GatewayRequest struct {
	RegionID      string  /*  资源池ID  */
	ProjectID     *string /*  企业项目 ID，默认为0  */
	Ipv6GatewayID string  /*  IPv6 网关 id  */
}

type CtvpcShowIPv6GatewayResponse struct {
	StatusCode  int32                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowIPv6GatewayReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       *string                                `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowIPv6GatewayReturnObjResponse struct {
	Name          *string `json:"name,omitempty"`          /*  ipv6网关名称  */
	Status        *string `json:"status,omitempty"`        /*  仅有一个状态：ACTIVE（正常）  */
	ProjectIdEcs  *string `json:"projectIdEcs,omitempty"`  /*  企业项目  */
	VpcID         *string `json:"vpcID,omitempty"`         /*  vpcID值  */
	VpcCidr       *string `json:"vpcCidr,omitempty"`       /*  无类别域间路由  */
	VpcName       *string `json:"vpcName,omitempty"`       /*  vpc名称  */
	Ipv6GatewayID *string `json:"ipv6GatewayID,omitempty"` /*  ipv6网关惟一值  */
	CreationTime  *string `json:"creationTime,omitempty"`  /*  创建时间  */
}
