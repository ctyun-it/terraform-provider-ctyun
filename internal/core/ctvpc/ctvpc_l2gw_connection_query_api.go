package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcL2gwConnectionQueryApi
/* 查询l2gw_connection列表
 */type CtvpcL2gwConnectionQueryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcL2gwConnectionQueryApi(client *core.CtyunClient) *CtvpcL2gwConnectionQueryApi {
	return &CtvpcL2gwConnectionQueryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/l2gw_connection/query",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcL2gwConnectionQueryApi) Do(ctx context.Context, credential core.Credential, req *CtvpcL2gwConnectionQueryRequest) (*CtvpcL2gwConnectionQueryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.L2gwID != nil {
		ctReq.AddParam("l2gwID", *req.L2gwID)
	}
	if req.L2ConnectionID != nil {
		ctReq.AddParam("l2ConnectionID", *req.L2ConnectionID)
	}
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
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
	var resp CtvpcL2gwConnectionQueryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcL2gwConnectionQueryRequest struct {
	RegionID       string  /*  资源池 ID  */
	L2gwID         *string /*  传入的l2gwID,如果传入此参数则表示查询某个l2gw下属的l2gw_connection，3.0资源池此参数必填  */
	L2ConnectionID *string /*  传入的l2gwConnectionID,如果传入此参数则表示查询某个l2gw_connection  */
	PageNumber     int32   /*  列表的页码，默认值为 1。  */
	PageNo         int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize       int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcL2gwConnectionQueryResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcL2gwConnectionQueryReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcL2gwConnectionQueryReturnObjResponse struct {
	Results      []*CtvpcL2gwConnectionQueryReturnObjResultsResponse `json:"results"`      /*  l2gw_connection组  */
	TotalCount   int32                                               `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                               `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                               `json:"totalPage"`    /*  总页数  */
}

type CtvpcL2gwConnectionQueryReturnObjResultsResponse struct {
	ID          *string `json:"ID,omitempty"`          /*  l2gw_connection 示例 ID  */
	Name        *string `json:"name,omitempty"`        /*  名字  */
	Description *string `json:"description,omitempty"` /*  描述  */
	L2conIp     *string `json:"l2conIp,omitempty"`     /*  ip  */
	SubnetID    *string `json:"subnetID,omitempty"`    /*  子网ID  */
	SubnetCidr  *string `json:"subnetCidr,omitempty"`  /*  子网cidr  */
	SubnetName  *string `json:"subnetName,omitempty"`  /*  子网name  */
	TunnelID    *string `json:"tunnelID,omitempty"`    /*  隧道号  */
	TunnelIp    *string `json:"tunnelIp,omitempty"`    /*  隧道ip  */
	CreatedAt   *string `json:"createdAt,omitempty"`   /*  创建时间  */
}
