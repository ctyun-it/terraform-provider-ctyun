package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowRouteTableApi
/* 查询路由表详情
 */type CtvpcShowRouteTableApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowRouteTableApi(client *core.CtyunClient) *CtvpcShowRouteTableApi {
	return &CtvpcShowRouteTableApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/route-table/show",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowRouteTableApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowRouteTableRequest) (*CtvpcShowRouteTableResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("routeTableID", req.RouteTableID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowRouteTableResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowRouteTableRequest struct {
	RegionID     string /*  区域id  */
	RouteTableID string /*  路由表 id  */
}

type CtvpcShowRouteTableResponse struct {
	StatusCode  int32                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowRouteTableReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       *string                               `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowRouteTableReturnObjResponse struct {
	Name            *string                                             `json:"name,omitempty"`        /*  路由表名字  */
	Description     *string                                             `json:"description,omitempty"` /*  路由表描述  */
	VpcID           *string                                             `json:"vpcID,omitempty"`       /*  虚拟私有云 id  */
	Id              *string                                             `json:"id,omitempty"`          /*  路由 id  */
	Freezing        *bool                                               `json:"freezing"`              /*  是否冻结  */
	RouteRulesCount int32                                               `json:"routeRulesCount"`       /*  路由表中的路由数  */
	CreatedAt       *string                                             `json:"createdAt,omitempty"`   /*  创建时间  */
	UpdatedAt       *string                                             `json:"updatedAt,omitempty"`   /*  更新时间  */
	RouteRules      []*string                                           `json:"routeRules"`            /*  路由规则 id 列表  */
	SubnetDetail    []*CtvpcShowRouteTableReturnObjSubnetDetailResponse `json:"subnetDetail"`          /*  子网配置详情  */
	RawType         int32                                               `json:"type"`                  /*  路由表类型:0-子网路由表，2-网关路由表  */
	Origin          *string                                             `json:"origin,omitempty"`      /*  路由表来源：default-系统默认; user-用户创建  */
}

type CtvpcShowRouteTableReturnObjSubnetDetailResponse struct {
	Id       *string `json:"id,omitempty"`       /*  路由下子网 id  */
	Name     *string `json:"name,omitempty"`     /*  路由下子网名字  */
	Cidr     *string `json:"cidr,omitempty"`     /*  ipv4 无类别域间路由  */
	Ipv6Cidr *string `json:"ipv6Cidr,omitempty"` /*  ipv6 无类别域间路由  */
}
