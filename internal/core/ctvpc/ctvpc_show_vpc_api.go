package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowVpcApi
/* 查询用户专有网络
 */type CtvpcShowVpcApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowVpcApi(client *core.CtyunClient) *CtvpcShowVpcApi {
	return &CtvpcShowVpcApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/query",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowVpcApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowVpcRequest) (*CtvpcShowVpcResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ProjectID != nil {
		ctReq.AddParam("projectID", *req.ProjectID)
	}
	ctReq.AddParam("vpcID", req.VpcID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowVpcResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowVpcRequest struct {
	RegionID  string  /*  资源池 ID  */
	ProjectID *string /*  企业项目 ID，默认为0  */
	VpcID     string  /*  VPC 的 ID  */
}

type CtvpcShowVpcResponse struct {
	StatusCode  int32                          `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                        `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                        `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                        `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowVpcReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcShowVpcReturnObjResponse struct {
	VpcID            *string   `json:"vpcID,omitempty"`       /*  vpc 示例 ID  */
	Name             *string   `json:"name,omitempty"`        /*  名称  */
	Description      *string   `json:"description,omitempty"` /*  描述  */
	CIDR             *string   `json:"CIDR,omitempty"`        /*  子网  */
	Ipv6Enabled      *bool     `json:"ipv6Enabled"`           /*  是否开启 ipv6  */
	EnableIpv6       *bool     `json:"enableIpv6"`            /*  是否开启 ipv6  */
	Ipv6CIDRS        []*string `json:"ipv6CIDRS"`             /*  ipv6 子网列表  */
	SubnetIDs        []*string `json:"subnetIDs"`             /*  子网 id 列表  */
	NatGatewayIDs    []*string `json:"natGatewayIDs"`         /*  网关 id 列表  */
	SecondaryCIDRS   []*string `json:"secondaryCIDRS"`        /*  附加网段  */
	ProjectID        *string   `json:"projectID,omitempty"`   /*  企业项目 ID，默认为0  */
	DhcpOptionsSetID []*string `json:"dhcpOptionsSetID"`      /*  VPC关联的dhcp选项集  */
}
