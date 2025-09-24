package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDhcpdisassociatevpcApi
/* dhcpoptionsets取消关联vpc
 */type CtvpcDhcpdisassociatevpcApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDhcpdisassociatevpcApi(client *core.CtyunClient) *CtvpcDhcpdisassociatevpcApi {
	return &CtvpcDhcpdisassociatevpcApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/dhcpoptionsets/dhcp_disassociate_vpc",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDhcpdisassociatevpcApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDhcpdisassociatevpcRequest) (*CtvpcDhcpdisassociatevpcResponse, error) {
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
	var resp CtvpcDhcpdisassociatevpcResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDhcpdisassociatevpcRequest struct {
	RegionID         string `json:"regionID,omitempty"`         /*  资源池 ID  */
	DhcpOptionSetsID string `json:"dhcpOptionSetsID,omitempty"` /*  集合ID  */
}

type CtvpcDhcpdisassociatevpcResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcDhcpdisassociatevpcReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcDhcpdisassociatevpcReturnObjResponse struct {
	DhcpOptionSetsID *string `json:"dhcpOptionSetsID,omitempty"` /*  ID  */
}
