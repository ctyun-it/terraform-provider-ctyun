package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDeleteRouteTableApi
/* 删除路由表，其中自定义路由表可以删除，默认路由表随 VPC 删除时一起删除。
 */type CtvpcDeleteRouteTableApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDeleteRouteTableApi(client *core.CtyunClient) *CtvpcDeleteRouteTableApi {
	return &CtvpcDeleteRouteTableApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/route-table/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDeleteRouteTableApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDeleteRouteTableRequest) (*CtvpcDeleteRouteTableResponse, error) {
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
	var resp CtvpcDeleteRouteTableResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDeleteRouteTableRequest struct {
	ClientToken  string `json:"clientToken,omitempty"`  /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID     string `json:"regionID,omitempty"`     /*  区域id  */
	RouteTableID string `json:"routeTableID,omitempty"` /*  路由表 id  */
}

type CtvpcDeleteRouteTableResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
