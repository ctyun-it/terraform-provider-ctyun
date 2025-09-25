package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowEndpointServiceApi
/* 查看终端节点服务详情
 */type CtvpcShowEndpointServiceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowEndpointServiceApi(client *core.CtyunClient) *CtvpcShowEndpointServiceApi {
	return &CtvpcShowEndpointServiceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpce/show-endpoint-service",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowEndpointServiceApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowEndpointServiceRequest) (*CtvpcShowEndpointServiceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("endpointServiceID", req.EndpointServiceID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowEndpointServiceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowEndpointServiceRequest struct {
	RegionID          string /*  资源池 ID  */
	EndpointServiceID string /*  终端节点服 ID  */
}

type CtvpcShowEndpointServiceResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800 为成功，900 为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为 success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode 为 900 时为业务细分错误码，三段式：product.module.code; statusCode 为 800 时为 SUCCESS  */
	ReturnObj   *CtvpcShowEndpointServiceReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcShowEndpointServiceReturnObjResponse struct {
	ID             *string                                              `json:"ID,omitempty"`          /*  终端节点服务 ID  */
	Name           *string                                              `json:"name,omitempty"`        /*  终端节点服务名称  */
	VpcID          *string                                              `json:"vpcID,omitempty"`       /*  所属的专有网络 id  */
	Description    *string                                              `json:"description,omitempty"` /*  描述  */
	RawType        *string                                              `json:"type,omitempty"`        /*  接口还是反向，interface:接口，reverse:反向  */
	AutoConnection *bool                                                `json:"autoConnection"`        /*  是否自动连接  */
	Rules          []*CtvpcShowEndpointServiceReturnObjRulesResponse    `json:"rules"`                 /*  接口规则数据  */
	Backends       []*CtvpcShowEndpointServiceReturnObjBackendsResponse `json:"backends"`              /*  后端数据  */
	CreatedAt      *string                                              `json:"createdAt,omitempty"`   /*  创建时间  */
	UpdatedAt      *string                                              `json:"updatedAt,omitempty"`   /*  更新时间  */
	DnsName        *string                                              `json:"dnsName,omitempty"`     /*  域名  */
}

type CtvpcShowEndpointServiceReturnObjRulesResponse struct {
	Protocol     *string `json:"protocol,omitempty"` /*  协议，TCP:TCP 协议,UDP:UDP 协议  */
	ServerPort   int32   `json:"serverPort"`         /*  服务端口(用于创建 backend 传入)  */
	EndpointPort int32   `json:"endpointPort"`       /*  节点端口(用于创建 rule 传入)  */
}

type CtvpcShowEndpointServiceReturnObjBackendsResponse struct {
	InstanceType *string `json:"instanceType,omitempty"` /*  vm:虚机类型,bm:智能网卡类型,vip:vip 类型,lb:负载均衡类型  */
	InstanceID   *string `json:"instanceID,omitempty"`   /*  实例 id  */
	ProtocolPort int32   `json:"protocolPort"`           /*  端口  */
	Weight       int32   `json:"weight"`                 /*  权重  */
}
