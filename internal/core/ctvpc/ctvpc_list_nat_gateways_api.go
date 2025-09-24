package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListNatGatewaysApi
/* 查询 NAT 网关列表
 */type CtvpcListNatGatewaysApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListNatGatewaysApi(client *core.CtyunClient) *CtvpcListNatGatewaysApi {
	return &CtvpcListNatGatewaysApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/describe-nat-gateways",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListNatGatewaysApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListNatGatewaysRequest) (*CtvpcListNatGatewaysResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.NatGatewayID != nil {
		ctReq.AddParam("natGatewayID", *req.NatGatewayID)
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
	var resp CtvpcListNatGatewaysResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListNatGatewaysRequest struct {
	RegionID     string  /*  区域id  */
	NatGatewayID *string /*  要查询的NAT网关的ID。  */
	PageNumber   int32   /*  列表的页码，默认值为1。  */
	PageNo       int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32   /*  分页查询时每页的行数，最大值为50，默认值为10。  */
}

type CtvpcListNatGatewaysResponse struct {
	StatusCode   int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtvpcListNatGatewaysReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	TotalCount   int32                                    `json:"totalCount"`            /*  列表条目数  */
	CurrentCount int32                                    `json:"currentCount"`          /*  分页查询时每页的行数。  */
	TotalPage    int32                                    `json:"totalPage"`             /*  总页数  */
	Error        *string                                  `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListNatGatewaysReturnObjResponse struct {
	ID           *string `json:"ID,omitempty"`           /*  nat 网关 id  */
	Name         *string `json:"name,omitempty"`         /*  nat 网关名字  */
	Description  *string `json:"description,omitempty"`  /*  nat 网关描述  */
	Status       int32   `json:"status"`                 /*  nat 网关状态: 0 表示创建中，2 表示运行中，3 表示冻结  */
	NatGatewayID *string `json:"natGatewayID,omitempty"` /*  nat 网关 id  */
	ZoneID       *string `json:"zoneID,omitempty"`       /*  可用区 ID  */
	State        *string `json:"state,omitempty"`        /*  NAT网关运行状态: running 表示运行中, creating 表示创建中, expired 表示已过期, freeze 表示已冻结  */
	VpcID        *string `json:"vpcID,omitempty"`        /*  虚拟私有云 id  */
	VpcName      *string `json:"vpcName,omitempty"`      /*  虚拟私有云名字  */
	ExpireTime   *string `json:"expireTime,omitempty"`   /*  过期时间  */
	CreationTime *string `json:"creationTime,omitempty"` /*  创建时间  */
	ProjectID    *string `json:"projectID,omitempty"`    /*  项目 ID  */
}
