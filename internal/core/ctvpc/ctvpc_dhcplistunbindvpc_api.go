package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDhcplistunbindvpcApi
/* 获取未绑定dhcp的 vpc 列表
 */type CtvpcDhcplistunbindvpcApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDhcplistunbindvpcApi(client *core.CtyunClient) *CtvpcDhcplistunbindvpcApi {
	return &CtvpcDhcplistunbindvpcApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/dhcpoptionsets/dhcp_list_unbind_vpc",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDhcplistunbindvpcApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDhcplistunbindvpcRequest) (*CtvpcDhcplistunbindvpcResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcDhcplistunbindvpcResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDhcplistunbindvpcRequest struct {
	RegionID string /*  资源池 ID  */
}

type CtvpcDhcplistunbindvpcResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcDhcplistunbindvpcReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcDhcplistunbindvpcReturnObjResponse struct {
	Id   *string `json:"id,omitempty"`   /*  vpc id  */
	Name *string `json:"name,omitempty"` /*  vpc 名字  */
}
