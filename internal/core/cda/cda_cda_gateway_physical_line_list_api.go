package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaGatewayPhysicalLineListApi
/* 查询专线网关已绑定的物理专线 */
type CdaCdaGatewayPhysicalLineListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaGatewayPhysicalLineListApi(client *core.CtyunClient) *CdaCdaGatewayPhysicalLineListApi {
	return &CdaCdaGatewayPhysicalLineListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/gateway/physical-line-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaGatewayPhysicalLineListApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaGatewayPhysicalLineListRequest) (*CdaCdaGatewayPhysicalLineListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaGatewayPhysicalLineListRequest
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
	var resp CdaCdaGatewayPhysicalLineListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaGatewayPhysicalLineListRequest struct {
	GatewayName string `json:"gatewayName"` /*  专线网关名字  */
	Account     string `json:"account"`     /*  天翼云客户邮箱  */
}

type CdaCdaGatewayPhysicalLineListResponse struct {
	StatusCode        *int32  `json:"statusCode"`        /*  返回状态码(800为成功，900为失败)  */
	Message           *string `json:"message"`           /*  失败时的错误描述，一般为英文描述  */
	Description       *string `json:"description"`       /*  失败时的错误描述，一般为中文描述  */
	ErrorCode         *string `json:"errorCode"`         /*  业务细分码，为product.module.code三段式码  */
	LineId            *string `json:"lineId"`            /*  物理专线ID  */
	LineName          *string `json:"lineName"`          /*  物理专线名字  */
	Account           *string `json:"account"`           /*  天翼云客户邮箱  */
	VrfName           *string `json:"vrfName"`           /*  专线网关名字  */
	ResourcePool      *string `json:"resourcePool"`      /*  资源池ID  */
	ResourcePoolName  *string `json:"resourcePoolName"`  /*  资源池名字  */
	IpVersion         *string `json:"ipVersion"`         /*  本参数表示包周期类型。<br>取值范围：<br>IPV4<br>IPV6<br>DUALSTACK  */
	LocalConnectIP    *string `json:"localConnectIP"`    /*  本端互联IP(IPV4或DUALSTACK必填)  */
	RemoteConnectIP   *string `json:"remoteConnectIP"`   /*  远端互联IP(IPV4或DUALSTACK必填)  */
	LocalConnectIPv6  *string `json:"localConnectIPv6"`  /*  本端互联IPv6(IPV6或DUALSTACK必填)  */
	RemoteConnectIPv6 *string `json:"remoteConnectIPv6"` /*  本端互联IPv6(IPV6或DUALSTACK必填)  */
	Hostname          *string `json:"hostname"`          /*  交换机hostname  */
	PortType          *string `json:"portType"`          /*  端口类型(10G、1G)  */
	PortName          *string `json:"portName"`          /*  端口名字  */
	PortNameBus       *string `json:"portNameBus"`       /*  业务口  */
	Bandwidth         *int32  `json:"bandwidth"`         /*  带宽(M)  */
	DeviceIp          *string `json:"deviceIp"`          /*  设备IP  */
	LineType          *string `json:"lineType"`          /*  物理专线类型(PON/IPRAN等)  */
	Tag               *int32  `json:"tag"`               /*  是否带vlan tag, 1 为带， 0为不带  */
	Vlan              *int32  `json:"vlan"`              /*  接入端口放行vlan, tag为1 ，此项必填  */
	IsShared          *int32  `json:"isShared"`          /*  端口类别: 独享(0), 共享(1)默认  */
	AccessPoint       *string `json:"accessPoint"`       /*  接入点，多AZ必填: AP1, AP2  */
	Location          *string `json:"location"`          /*  接入位置  */
	Linecode          *string `json:"linecode"`          /*  电路代号  */
	RegionID          *string `json:"regionID"`          /*  资源池ID  */
	Fuid              *string `json:"fuid"`              /*  物理专线ID  */
	AuthorizedAccount *string `json:"authorizedAccount"` /*  授权的账号  */
	ProjectIDEcs      *string `json:"projectIDEcs"`      /*  Project ID  */
	Layer             *string `json:"layer"`             /*  网络协议栈层级  */
	LineCreateTime    *string `json:"lineCreateTime"`    /*  创建时间  */
	FuserLastUpdated  *string `json:"fuserLastUpdated"`  /*  最近更新时间  */
	DeleteTime        *string `json:"deleteTime"`        /*  删除时间  */
	CustomerId        *string `json:"customerId"`        /*  客户ID  */
	Error             *string `json:"error"`             /*  业务细分码，为product.module.code三段式码  */
}
