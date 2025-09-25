package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcShowEndpointServiceConnectionsApi
/* 终端节点服务连接查询
 */type CtvpcShowEndpointServiceConnectionsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowEndpointServiceConnectionsApi(client *core.CtyunClient) *CtvpcShowEndpointServiceConnectionsApi {
	return &CtvpcShowEndpointServiceConnectionsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpce/show-endpoint-service-connections",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowEndpointServiceConnectionsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowEndpointServiceConnectionsRequest) (*CtvpcShowEndpointServiceConnectionsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("endpointServiceID", req.EndpointServiceID)
	if req.Page != 0 {
		ctReq.AddParam("page", strconv.FormatInt(int64(req.Page), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowEndpointServiceConnectionsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowEndpointServiceConnectionsRequest struct {
	RegionID          string /*  资源池ID  */
	EndpointServiceID string /*  终端节点服务id  */
	Page              int32  /*  分页参数  */
	PageNo            int32  /*  列表的页码，默认值为 1, 推荐使用该字段, page 后续会废弃  */
	PageSize          int32  /*  每页数据量大小  */
}

type CtvpcShowEndpointServiceConnectionsResponse struct {
	StatusCode  int32                                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcShowEndpointServiceConnectionsReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                                 `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowEndpointServiceConnectionsReturnObjResponse struct {
	EndpointID        *string `json:"endpointID,omitempty"`        /*  终端节点ID  */
	EndpointName      *string `json:"endpointName,omitempty"`      /*  终端节点名称  */
	EndpointServiceID *string `json:"endpointServiceID,omitempty"` /*  终端节点服务ID  */
	EndpointType      *string `json:"endpointType,omitempty"`      /*  接口还反向  */
	ConnectionStatus  *string `json:"connectionStatus,omitempty"`  /*  连接状态  */
	Status            *string `json:"status,omitempty"`            /*  状态：UP/DOWN  */
	Account           *string `json:"account,omitempty"`           /*  账户信息  */
}
