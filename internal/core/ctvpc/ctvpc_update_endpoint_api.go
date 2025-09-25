package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUpdateEndpointApi
/* 更新终端节点
 */type CtvpcUpdateEndpointApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUpdateEndpointApi(client *core.CtyunClient) *CtvpcUpdateEndpointApi {
	return &CtvpcUpdateEndpointApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/update-endpoint",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUpdateEndpointApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUpdateEndpointRequest) (*CtvpcUpdateEndpointResponse, error) {
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
	var resp CtvpcUpdateEndpointResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUpdateEndpointRequest struct {
	ClientToken       string    `json:"clientToken,omitempty"`       /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID          string    `json:"regionID,omitempty"`          /*  资源池ID  */
	EndpointID        string    `json:"endpointID,omitempty"`        /*  终端节点id  */
	EndpointName      *string   `json:"endpointName,omitempty"`      /*  终端节点名称  */
	EnableWhitelist   *bool     `json:"enableWhitelist"`             /*  白名单开关 True.开启 False.关闭  */
	EnableDns         *bool     `json:"enableDns"`                   /*  是否开启dns, true:开启,false:关闭  */
	Whitelist         []*string `json:"whitelist,omitempty"`         /*  白名单列表，最多支持同时添加 20 个 ip  */
	DeleteProtection  *bool     `json:"deleteProtection"`            /*  是否开启删除保护, true:开启,false:关闭，不传默认关闭  */
	ProtectionService *string   `json:"protectionService,omitempty"` /*  删除保护使能服务  */
}

type CtvpcUpdateEndpointResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
