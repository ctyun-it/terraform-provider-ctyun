package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowEndpointApi
/* 查看终端节点详情
 */type CtvpcShowEndpointApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowEndpointApi(client *core.CtyunClient) *CtvpcShowEndpointApi {
	return &CtvpcShowEndpointApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpce/show-endpoint",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowEndpointApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowEndpointRequest) (*CtvpcShowEndpointResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("endpointID", req.EndpointID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowEndpointResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowEndpointRequest struct {
	RegionID   string /*  资源池ID  */
	EndpointID string /*  终端节点 ID  */
}

type CtvpcShowEndpointResponse struct {
	StatusCode  int32                               `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                             `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                             `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                             `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowEndpointReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                             `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowEndpointReturnObjResponse struct {
	Endpoint *CtvpcShowEndpointReturnObjEndpointResponse `json:"endpoint"` /*  终端节点列表数据  */
}

type CtvpcShowEndpointReturnObjEndpointResponse struct {
	ID                *string                                                `json:"ID,omitempty"`                /*  终端节点ID  */
	EndpointServiceID *string                                                `json:"endpointServiceID,omitempty"` /*  终端节点服务ID  */
	RawType           *string                                                `json:"type,omitempty"`              /*  接口还是反向，interface:接口，reverse:反向  */
	Name              *string                                                `json:"name,omitempty"`              /*  终端节点名称  */
	VpcID             *string                                                `json:"vpcID,omitempty"`             /*  所属的专有网络id  */
	VpcAddress        *string                                                `json:"vpcAddress,omitempty"`        /*  私网地址  */
	Whitelist         *string                                                `json:"whitelist,omitempty"`         /*  白名单  */
	Status            int32                                                  `json:"status"`                      /*  endpoint状态, 1 表示已链接，2 表示未链接  */
	Description       *string                                                `json:"description,omitempty"`       /*  描述  */
	EndpointObj       *CtvpcShowEndpointReturnObjEndpointEndpointObjResponse `json:"endpointObj"`                 /*  后端节点信息，可能为 null  */
	CreatedTime       *string                                                `json:"createdTime,omitempty"`       /*  创建时间  */
	UpdatedTime       *string                                                `json:"updatedTime,omitempty"`       /*  更新时间  */
}

type CtvpcShowEndpointReturnObjEndpointEndpointObjResponse struct {
	SubnetID  *string   `json:"subnetID,omitempty"` /*  子网id  */
	PortID    *string   `json:"portID,omitempty"`   /*  端口id  */
	Ip        *string   `json:"ip,omitempty"`       /*  私网地址  */
	EnableDns int32     `json:"enableDns"`          /*  是否开启 dns  */
	DnsNames  []*string `json:"dnsNames"`           /*  dns名称列表  */
}
