package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcListDnatEntriesApi
/* 查询 dnat 列表
 */type CtvpcListDnatEntriesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListDnatEntriesApi(client *core.CtyunClient) *CtvpcListDnatEntriesApi {
	return &CtvpcListDnatEntriesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/describe-dnat-entries",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListDnatEntriesApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListDnatEntriesRequest) (*CtvpcListDnatEntriesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("natGatewayID", req.NatGatewayID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListDnatEntriesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListDnatEntriesRequest struct {
	RegionID     string /*  区域id  */
	NatGatewayID string /*  要查询的NAT网关的ID。  */
}

type CtvpcListDnatEntriesResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []CtvpcListDnatEntriesReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       *string                                 `json:"error,omitempty"`       /*  ingress  */
}

type CtvpcListDnatEntriesReturnObjResponse struct {
	CreationTime              *string `json:"creationTime,omitempty"`              /*  创建时间。  */
	Description               *string `json:"description,omitempty"`               /*  描述信息  */
	Id                        *string `json:"id,omitempty"`                        /*  dnatID 值。  */
	DNatID                    *string `json:"dNatID,omitempty"`                    /*  dnatID 值。  */
	IpExpireTime              *string `json:"ipExpireTime,omitempty"`              /*  ip到期时间。  */
	ExternalID                *string `json:"externalID,omitempty"`                /*  弹性 IP  */
	ExternalIp                *string `json:"externalIp,omitempty"`                /*  弹性 IP 地址  */
	ExternalPort              int32   `json:"externalPort"`                        /*  外部访问端口  */
	InternalPort              int32   `json:"internalPort"`                        /*  内部访问端口  */
	InternalIp                *string `json:"internalIp,omitempty"`                /*  内网 IP 地址  */
	Protocol                  *string `json:"protocol,omitempty"`                  /*  TCP:转发TCP协议的报文 UDP：转发UDP协议的报文。  */
	State                     *string `json:"state,omitempty"`                     /*  运行状态: ACTIVE / FREEZING / CREATING  */
	VirtualMachineDisplayName *string `json:"virtualMachineDisplayName,omitempty"` /*  虚拟机展示名称。  */
	VirtualMachineID          *string `json:"virtualMachineID,omitempty"`          /*  虚拟机id。  */
	VirtualMachineName        *string `json:"virtualMachineName,omitempty"`        /*  虚拟机名称。  */
}
