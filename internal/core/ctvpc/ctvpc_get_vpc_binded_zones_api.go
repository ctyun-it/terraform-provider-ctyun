package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcGetVpcBindedZonesApi
/* 查询 vpc 绑定的 zone
 */type CtvpcGetVpcBindedZonesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGetVpcBindedZonesApi(client *core.CtyunClient) *CtvpcGetVpcBindedZonesApi {
	return &CtvpcGetVpcBindedZonesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/private-zone/get-vpc-binded-zones",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGetVpcBindedZonesApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGetVpcBindedZonesRequest) (*CtvpcGetVpcBindedZonesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("vpcID", req.VpcID)
	ctReq.AddParam("regionID", req.RegionID)
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
	var resp CtvpcGetVpcBindedZonesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGetVpcBindedZonesRequest struct {
	VpcID    string /*  虚拟私有云 ID  */
	RegionID string /*  资源池ID  */
	PageNo   int32  /*  列表的页码，默认值为 1  */
	PageSize int32  /*  分页查询时每页的行数，最大值为50，默认值为10。  */
}

type CtvpcGetVpcBindedZonesResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGetVpcBindedZonesReturnObjResponse `json:"returnObj"`             /*  object  */
	Error       *string                                  `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcGetVpcBindedZonesReturnObjResponse struct {
	Results      []*CtvpcGetVpcBindedZonesReturnObjResultsResponse `json:"results"`      /*  dns 记录  */
	TotalCount   int32                                             `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                             `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                             `json:"totalPage"`    /*  总页数  */
}

type CtvpcGetVpcBindedZonesReturnObjResultsResponse struct {
	ZoneID       *string `json:"zoneID,omitempty"`       /*  名称  */
	Name         *string `json:"name,omitempty"`         /*  zone名称  */
	Description  *string `json:"description,omitempty"`  /*  描述  */
	ProxyPattern *string `json:"proxyPattern,omitempty"` /*  zone, record  */
	TTL          int32   `json:"TTL"`                    /*  zone ttl, default is 300  */
	CreatedAt    *string `json:"createdAt,omitempty"`    /*  创建时间  */
	UpdatedAt    *string `json:"updatedAt,omitempty"`    /*  更新时间  */
}
