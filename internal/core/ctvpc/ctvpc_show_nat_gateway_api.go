package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowNatGatewayApi
/* 查询 NAT 网关详情
 */type CtvpcShowNatGatewayApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowNatGatewayApi(client *core.CtyunClient) *CtvpcShowNatGatewayApi {
	return &CtvpcShowNatGatewayApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/get-nat-gateway-attribute",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowNatGatewayApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowNatGatewayRequest) (*CtvpcShowNatGatewayResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("natGatewayID", req.NatGatewayID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowNatGatewayResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowNatGatewayRequest struct {
	RegionID     string /*  区域id  */
	NatGatewayID string /*  要查询的NAT网关的ID。  */
}

type CtvpcShowNatGatewayResponse struct {
	StatusCode  int32                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowNatGatewayReturnObjResponse `json:"returnObj"`             /*  object  */
	Error       *string                               `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowNatGatewayReturnObjResponse struct {
	Name         *string     `json:"name,omitempty"`         /*  NAT网关实例名称  */
	Status       int32       `json:"status"`                 /*  NAT 网关状态: 0 表示创建中，2 表示运行中，3 表示冻结  */
	State        *string     `json:"state,omitempty"`        /*  NAT网关运行状态: running 表示运行中, creating 表示创建中, expired 表示已过期, freeze 表示已冻结  */
	Specs        *string     `json:"specs,omitempty"`        /*  规格  */
	ZoneID       *string     `json:"zoneID,omitempty"`       /*  NAT网关所在的可用区ID。  */
	VpcID        *string     `json:"vpcID,omitempty"`        /*  要查询的NAT网关所属VPC的ID。  */
	ProjectID    *string     `json:"projectID,omitempty"`    /*  项目类型：default-企业项目  */
	VpcName      *string     `json:"vpcName,omitempty"`      /*  NAT所属的专有网络名字  */
	VpcCidr      *string     `json:"vpcCidr,omitempty"`      /*  当前网关所属的vpc cidr  */
	CreationTime *string     `json:"creationTime,omitempty"` /*  NAT网关的创建时间  */
	ExpiredTime  *string     `json:"expiredTime,omitempty"`  /*  NAT网关实例的过期时间  */
	NatGatewayID *string     `json:"natGatewayID,omitempty"` /*  NAT网关的ID  */
	Description  *string     `json:"description,omitempty"`  /*  NAT网关实例的描述  */
	SnatTable    interface{} `json:"snatTable"`              /*  SNAT列表信息  */
	DnatTable    interface{} `json:"dnatTable"`              /*  DNAT列表的信息  */
}

type CtvpcShowNatGatewayReturnObjSnatTableResponse struct {
	Id           *string                                                   `json:"id,omitempty"`           /*  snat id  */
	SNatID       *string                                                   `json:"sNatID,omitempty"`       /*  snat id  */
	Description  *string                                                   `json:"description,omitempty"`  /*  描述信息  */
	Status       *string                                                   `json:"status,omitempty"`       /*  状态: ACTIVE 表示运行中，Creating 表示创建中，Freezing 表示冻结  */
	VpcID        *string                                                   `json:"vpcID,omitempty"`        /*  要查询的NAT网关所属VPC的ID  */
	VpcName      *string                                                   `json:"vpcName,omitempty"`      /*  要查询的NAT网关所属VPC的名称  */
	VpcCidr      *string                                                   `json:"vpcCidr,omitempty"`      /*  要查询的NAT网关所属VPC的cidr  */
	SubnetCidr   *string                                                   `json:"subnetCidr,omitempty"`   /*  要查询的NAT网关所属VPC子网的cidr  */
	SubnetType   int32                                                     `json:"subnetType"`             /*  子网类型：1-有vpcID的子网，0-自定义  */
	CreationTime *string                                                   `json:"creationTime,omitempty"` /*  创建时间  */
	IpAddress    []*CtvpcShowNatGatewayReturnObjSnatTableIpAddressResponse `json:"ipAddress"`              /*  eip地址信息  */
	SubnetName   *string                                                   `json:"subnetName,omitempty"`   /*  子网名字  */
}

type CtvpcShowNatGatewayReturnObjDnatTableResponse struct {
	CreationTime       *string `json:"creationTime,omitempty"`       /*  创建时间  */
	Description        *string `json:"description,omitempty"`        /*  描述信息  */
	Id                 *string `json:"id,omitempty"`                 /*  dnat id  */
	DNatID             *string `json:"dNatID,omitempty"`             /*  dnat id  */
	IpExpireTime       *string `json:"ipExpireTime,omitempty"`       /*  ip到期时间  */
	ExternalID         *string `json:"externalID,omitempty"`         /*  弹性公网id  */
	ExternalPort       *string `json:"externalPort,omitempty"`       /*  公网端口  */
	ExternalIp         *string `json:"externalIp,omitempty"`         /*  弹性公网ip  */
	InternalPort       *string `json:"internalPort,omitempty"`       /*  私网端口  */
	InternalIp         *string `json:"internalIp,omitempty"`         /*  内网 IP 地址  */
	Protocol           *string `json:"protocol,omitempty"`           /*  TCP:转发TCP协议的报文 UDP：转发UDP协议的报文。  */
	State              *string `json:"state,omitempty"`              /*  运行状态: ACTIVE / FREEZING / CREATING  */
	VirtualMachineName *string `json:"virtualMachineName,omitempty"` /*  虚拟机名称  */
	VirtualMachineID   *string `json:"virtualMachineID,omitempty"`   /*  虚拟机id  */
}

type CtvpcShowNatGatewayReturnObjSnatTableIpAddressResponse struct {
	ID        *string `json:"ID,omitempty"`        /*  弹性 IP id  */
	IpAddress *string `json:"ipAddress,omitempty"` /*  eip所属的ip地址。  */
}
