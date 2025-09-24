package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcGwQueryApi
/* 查询gateway列表
 */type CtvpcGwQueryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGwQueryApi(client *core.CtyunClient) *CtvpcGwQueryApi {
	return &CtvpcGwQueryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/l2gw/gw_query",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGwQueryApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGwQueryRequest) (*CtvpcGwQueryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("vpcID", req.VpcID)
	ctReq.AddParam("linkGwType", req.LinkGwType)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcGwQueryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGwQueryRequest struct {
	RegionID   string /*  资源池 ID  */
	VpcID      string /*  vpcid  */
	LinkGwType string /*  隧道连接方式 linegw：云专线  vpn：VPN  */
}

type CtvpcGwQueryResponse struct {
	StatusCode  int32                          `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                        `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                        `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                        `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGwQueryReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcGwQueryReturnObjResponse struct {
	Gws []*CtvpcGwQueryReturnObjGwsResponse `json:"gws"` /*  网关列表  */
}

type CtvpcGwQueryReturnObjGwsResponse struct {
	Id   *string `json:"id,omitempty"`   /*  vpn或云专线id  */
	Name *string `json:"name,omitempty"` /*  vpn或云专线名称  */
}
