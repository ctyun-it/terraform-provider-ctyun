package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDhcpreplacevpcApi
/* vpc替换dhcpoptionset
 */type CtvpcDhcpreplacevpcApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDhcpreplacevpcApi(client *core.CtyunClient) *CtvpcDhcpreplacevpcApi {
	return &CtvpcDhcpreplacevpcApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/dhcpoptionsets/dhcp_replace_vpc",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDhcpreplacevpcApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDhcpreplacevpcRequest) (*CtvpcDhcpreplacevpcResponse, error) {
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
	var resp CtvpcDhcpreplacevpcResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDhcpreplacevpcRequest struct {
	RegionID         string `json:"regionID,omitempty"`         /*  资源池 ID  */
	DhcpOptionSetsID string `json:"dhcpOptionSetsID,omitempty"` /*  集合ID  */
	VpcID            string `json:"vpcID,omitempty"`            /*  vpc id,必须是已绑定dhcp的vpc才能替换  */
}

type CtvpcDhcpreplacevpcResponse struct {
	StatusCode  int32                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcDhcpreplacevpcReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcDhcpreplacevpcReturnObjResponse struct {
	DhcpOptionSetsID *string `json:"dhcpOptionSetsID,omitempty"` /*  ID  */
}
