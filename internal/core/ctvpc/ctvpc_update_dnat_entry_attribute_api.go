package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUpdateDnatEntryAttributeApi
/* 修改 dnat 规则
 */type CtvpcUpdateDnatEntryAttributeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUpdateDnatEntryAttributeApi(client *core.CtyunClient) *CtvpcUpdateDnatEntryAttributeApi {
	return &CtvpcUpdateDnatEntryAttributeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/modify-dnat-entry",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUpdateDnatEntryAttributeApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUpdateDnatEntryAttributeRequest) (*CtvpcUpdateDnatEntryAttributeResponse, error) {
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
	var resp CtvpcUpdateDnatEntryAttributeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUpdateDnatEntryAttributeRequest struct {
	RegionID           string  `json:"regionID,omitempty"`         /*  区域id  */
	DNatID             string  `json:"dNatID,omitempty"`           /*  DNAT网关ID  */
	ExternalID         *string `json:"externalID,omitempty"`       /*  弹性公网id  */
	ExternalPort       int32   `json:"externalPort"`               /*  弹性IP公网端口, 1 - 1024  */
	VirtualMachineID   *string `json:"virtualMachineID,omitempty"` /*  云主机  */
	VirtualMachineType int32   `json:"virtualMachineType"`         /*  云主机类型1-选择云主机，serverType 字段必传， 2-自定义，internalIp必传  */
	InternalIp         *string `json:"internalIp,omitempty"`       /*  内部 IP  */
	InternalPort       int32   `json:"internalPort"`               /*  主机内网端口，1-65535  */
	Protocol           string  `json:"protocol,omitempty"`         /*  支持协议：tcp/udp  */
	ClientToken        string  `json:"clientToken,omitempty"`      /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	Description        *string `json:"description,omitempty"`      /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
	ServerType         *string `json:"serverType,omitempty"`       /*  当 virtualMachineType 为 1 时，serverType 必传，支持: VM / BM （仅支持大写）  */
}

type CtvpcUpdateDnatEntryAttributeResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
