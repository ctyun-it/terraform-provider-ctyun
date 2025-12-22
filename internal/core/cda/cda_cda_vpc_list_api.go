package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaVpcListApi
/* 查询用户专线网关下添加的VPC信息 */
type CdaCdaVpcListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaVpcListApi(client *core.CtyunClient) *CdaCdaVpcListApi {
	return &CdaCdaVpcListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/vpc/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaVpcListApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaVpcListRequest) (*CdaCdaVpcListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaVpcListRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CdaCdaVpcListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaVpcListRequest struct {
	GatewayName string `json:"gatewayName"` /*  专线网关名字  */
	Account     string `json:"account"`     /*  天翼云客户邮箱  */
}

type CdaCdaVpcListResponse struct {
	StatusCode            *int32                            `json:"statusCode"`            /*  返回状态码(800为成功，900为失败)  */
	Message               *string                           `json:"message"`               /*  失败时的错误描述，一般为英文描述  */
	Description           *string                           `json:"description"`           /*  失败时的错误描述，一般为中文描述  */
	ErrorCode             *string                           `json:"errorCode"`             /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail           *CdaCdaVpcListErrorDetailResponse `json:"errorDetail"`           /*  错误明细  */
	CdaId                 *string                           `json:"cdaId"`                 /*  专线ID  */
	Account               *string                           `json:"account"`               /*  天翼云账号  */
	VpcId                 *string                           `json:"vpcId"`                 /*  VPC  ID  */
	VrfName               *string                           `json:"vrfName"`               /*  专线网关名字  */
	DcType                *string                           `json:"dcType"`                /*  本参数表示资源池类型。<br>取值范围：<br>MAZ<br>CNP  */
	AzNameIpv6            *string                           `json:"azNameIpv6"`            /*  IPv6 VPC的可用区名字  */
	CdaIdV6               *string                           `json:"cdaIdV6"`               /*  IPv6 VPC的专线ID  */
	VpcSubnet             *string                           `json:"vpcSubnet"`             /*  VPC子网  */
	ResourcePool          *string                           `json:"resourcePool"`          /*  资源池ID  */
	EmazId                *string                           `json:"emazId"`                /*  多可用区ID  */
	IpVersion             *string                           `json:"ipVersion"`             /*  本参数表示包周期类型。<br>取值范围：<br>IPV4<br>IPV6<br>DUALSTACK  */
	VirtualBandwidth      *int32                            `json:"virtualBandwidth"`      /*  虚拟带宽  */
	VpcName               *string                           `json:"vpcName"`               /*  VPC名字  */
	EmazIdV6              *string                           `json:"emazIdV6"`              /*  多可用区ID(IPv6)  */
	GuestGatewayList      []*string                         `json:"guestGatewayList"`      /*  客户侧网关列表  */
	AzName                *string                           `json:"azName"`                /*  可用区名字  */
	VpcSubnetIpv6         *string                           `json:"vpcSubnetIpv6"`         /*  VPC IPv6子网  */
	VpcNetworkSegment     *string                           `json:"vpcNetworkSegment"`     /*  VPC网段  */
	CtUserId              *string                           `json:"ctUserId"`              /*  天翼云用户ID  */
	VpcNetworkSegmentIpv6 *string                           `json:"vpcNetworkSegmentIpv6"` /*  VPC IPv6网段  */
	GuestGatewayListIpv6  []*string                         `json:"guestGatewayListIpv6"`  /*  客户侧IPv6网关列表  */
	Error                 *string                           `json:"error"`                 /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaVpcListErrorDetailResponse struct{}
