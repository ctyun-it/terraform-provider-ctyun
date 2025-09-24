package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcQueryCtyunServiceApi
/* 查看天翼云终端节点服务
 */type CtvpcQueryCtyunServiceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcQueryCtyunServiceApi(client *core.CtyunClient) *CtvpcQueryCtyunServiceApi {
	return &CtvpcQueryCtyunServiceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpce/query-ctyun-service",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcQueryCtyunServiceApi) Do(ctx context.Context, credential core.Credential, req *CtvpcQueryCtyunServiceRequest) (*CtvpcQueryCtyunServiceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.VpcID != nil {
		ctReq.AddParam("vpcID", *req.VpcID)
	}
	if req.EndpointServiceID != nil {
		ctReq.AddParam("endpointServiceID", *req.EndpointServiceID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcQueryCtyunServiceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcQueryCtyunServiceRequest struct {
	RegionID          string  /*  资源池ID  */
	PageNumber        int32   /*  列表的页码，默认值为 1, 推荐使用该字段, page 后续会废弃  */
	PageSize          int32   /*  每页数据量大小  */
	VpcID             *string /*  虚拟私有云 ID  */
	EndpointServiceID *string /*  终端节点服务 ID  */
}

type CtvpcQueryCtyunServiceResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcQueryCtyunServiceReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcQueryCtyunServiceReturnObjResponse struct {
	TotalCount   int32                                             `json:"totalCount"`   /*  总条数  */
	TotalPage    int32                                             `json:"totalPage"`    /*  总页数  */
	CurrentCount int32                                             `json:"currentCount"` /*  总页数  */
	Results      []*CtvpcQueryCtyunServiceReturnObjResultsResponse `json:"results"`      /*  终端节点服务列表  */
}

type CtvpcQueryCtyunServiceReturnObjResultsResponse struct {
	EndpointServiceID *string                                                `json:"endpointServiceID,omitempty"` /*  终端服务 ID  */
	Name              *string                                                `json:"name,omitempty"`              /*  终端服务名  */
	Description       *string                                                `json:"description,omitempty"`       /*  终端服务描述  */
	RawType           *string                                                `json:"type,omitempty"`              /*  终端服务连接类型：interface / reverse  */
	AutoConnection    *bool                                                  `json:"autoConnection"`              /*  是否自动连接  */
	Rules             []*CtvpcQueryCtyunServiceReturnObjResultsRulesResponse `json:"rules"`                       /*  规则  */
	CreatedAt         *string                                                `json:"createdAt,omitempty"`         /*  创建时间  */
	UpdatedAt         *string                                                `json:"updatedAt,omitempty"`         /*  更新时间  */
	DnsName           *string                                                `json:"dnsName,omitempty"`           /*  域名  */
}

type CtvpcQueryCtyunServiceReturnObjResultsRulesResponse struct {
	Protocol     *string `json:"protocol,omitempty"` /*  协议：TCP / UDP  */
	ServerPort   int32   `json:"serverPort"`         /*  服务端口  */
	EndpointPort int32   `json:"endpointPort"`       /*  终端端口  */
}
