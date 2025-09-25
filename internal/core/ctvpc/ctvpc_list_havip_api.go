package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcListHavipApi
/* 查询高可用虚IP列表
 */type CtvpcListHavipApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListHavipApi(client *core.CtyunClient) *CtvpcListHavipApi {
	return &CtvpcListHavipApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/havip/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListHavipApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListHavipRequest) (*CtvpcListHavipResponse, error) {
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
	var resp CtvpcListHavipResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListHavipRequest struct {
	ClientToken string  `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  `json:"regionID,omitempty"`    /*  资源池ID  */
	ProjectID   *string `json:"projectID,omitempty"`   /*  企业项目ID，默认为'0'  */
}

type CtvpcListHavipResponse struct {
	StatusCode  int32                              `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                            `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                            `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                            `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcListHavipReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                            `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListHavipReturnObjResponse struct {
	Id           *string                                        `json:"id,omitempty"`       /*  高可用虚IP的ID  */
	Ipv4         *string                                        `json:"ipv4,omitempty"`     /*  IPv4地址  */
	VpcID        *string                                        `json:"vpcID,omitempty"`    /*  虚拟私有云的的id  */
	SubnetID     *string                                        `json:"subnetID,omitempty"` /*  子网id  */
	InstanceInfo []*CtvpcListHavipReturnObjInstanceInfoResponse `json:"instanceInfo"`       /*  绑定实例相关信息  */
	NetworkInfo  []*CtvpcListHavipReturnObjNetworkInfoResponse  `json:"networkInfo"`        /*  绑定弹性 IP 相关信息  */
}

type CtvpcListHavipReturnObjInstanceInfoResponse struct {
	InstanceName *string `json:"instanceName,omitempty"` /*  实例名  */
	Id           *string `json:"id,omitempty"`           /*  实例 ID  */
	PrivateIp    *string `json:"privateIp,omitempty"`    /*  实例私有 IP  */
	PrivateIpv6  *string `json:"privateIpv6,omitempty"`  /*  实例的 IPv6 地址, 可以为空字符串  */
	PublicIp     *string `json:"publicIp,omitempty"`     /*  实例公网 IP  */
}

type CtvpcListHavipReturnObjNetworkInfoResponse struct {
	EipID *string `json:"eipID,omitempty"` /*  弹性 IP ID  */
}
