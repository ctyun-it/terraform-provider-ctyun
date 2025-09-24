package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowHavipApi
/* 查看高可用虚 IP 详情
 */type CtvpcShowHavipApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowHavipApi(client *core.CtyunClient) *CtvpcShowHavipApi {
	return &CtvpcShowHavipApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/havip/show",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowHavipApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowHavipRequest) (*CtvpcShowHavipResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("haVipID", req.HaVipID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowHavipResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowHavipRequest struct {
	RegionID string /*  资源池 ID  */
	HaVipID  string /*  高可用虚 IP 的 ID  */
}

type CtvpcShowHavipResponse struct {
	StatusCode  int32                              `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                            `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                            `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                            `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcShowHavipReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                            `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowHavipReturnObjResponse struct {
	Id           *string                                        `json:"id,omitempty"`       /*  高可用虚 IP 的 ID  */
	Ipv4         *string                                        `json:"ipv4,omitempty"`     /*  IPv4 地址  */
	VpcID        *string                                        `json:"vpcID,omitempty"`    /*  虚拟私有云的的 id  */
	SubnetID     *string                                        `json:"subnetID,omitempty"` /*  子网 id  */
	InstanceInfo []*CtvpcShowHavipReturnObjInstanceInfoResponse `json:"instanceInfo"`       /*  绑定实例相关信息  */
	NetworkInfo  []*CtvpcShowHavipReturnObjNetworkInfoResponse  `json:"networkInfo"`        /*  绑定弹性 IP 相关信息  */
	BindPorts    []*CtvpcShowHavipReturnObjBindPortsResponse    `json:"bindPorts"`          /*  绑定网卡信息  */
}

type CtvpcShowHavipReturnObjInstanceInfoResponse struct {
	InstanceName *string `json:"instanceName,omitempty"` /*  实例名  */
	Id           *string `json:"id,omitempty"`           /*  实例 ID  */
	PrivateIp    *string `json:"privateIp,omitempty"`    /*  实例私有 IP  */
	PrivateIpv6  *string `json:"privateIpv6,omitempty"`  /*  实例的 IPv6 地址, 可以为空字符串  */
	PublicIp     *string `json:"publicIp,omitempty"`     /*  实例公网 IP  */
}

type CtvpcShowHavipReturnObjNetworkInfoResponse struct {
	EipID *string `json:"eipID,omitempty"` /*  弹性 IP ID  */
}

type CtvpcShowHavipReturnObjBindPortsResponse struct {
	PortID    *string `json:"portID,omitempty"`    /*  网卡 ID  */
	Role      *string `json:"role,omitempty"`      /*  keepalive 角色: master / slave  */
	CreatedAt *string `json:"createdAt,omitempty"` /*  创建时间  */
}
