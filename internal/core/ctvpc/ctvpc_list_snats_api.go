package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListSnatsApi
/* 获取SNAT列表
 */type CtvpcListSnatsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListSnatsApi(client *core.CtyunClient) *CtvpcListSnatsApi {
	return &CtvpcListSnatsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/list-snats",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListSnatsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListSnatsRequest) (*CtvpcListSnatsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.NatGatewayID != nil {
		ctReq.AddParam("natGatewayID", *req.NatGatewayID)
	}
	if req.SNatID != nil {
		ctReq.AddParam("sNatID", *req.SNatID)
	}
	if req.SubnetID != nil {
		ctReq.AddParam("subnetID", *req.SubnetID)
	}
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListSnatsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListSnatsRequest struct {
	RegionID     string  /*  区域id  */
	NatGatewayID *string /*  要查询的NAT网关的ID。  */
	SNatID       *string /*  snat id  */
	SubnetID     *string /*  子网 id  */
	PageNumber   int32   /*  列表的页码，默认值为1。  */
	PageSize     int32   /*  分页查询时每页的行数，最大值为50，默认值为10。  */
}

type CtvpcListSnatsResponse struct {
	StatusCode  int32                            `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                          `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                          `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                          `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcListSnatsReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtvpcListSnatsReturnObjResponse struct {
	Results      []*CtvpcListSnatsReturnObjResultsResponse `json:"results"`      /*  snat 列表  */
	TotalCount   int32                                     `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                     `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                     `json:"totalPage"`    /*  总页数  */
}

type CtvpcListSnatsReturnObjResultsResponse struct {
	SNatID       *string                                       `json:"sNatID,omitempty"`       /*  snat id  */
	Description  *string                                       `json:"description,omitempty"`  /*  描述信息  */
	SubnetCidr   *string                                       `json:"subnetCidr,omitempty"`   /*  要查询的NAT网关所属VPC子网的cidr  */
	SubnetType   int32                                         `json:"subnetType"`             /*  子网类型：1-有vpcID的子网，0-自定义  */
	CreationTime *string                                       `json:"creationTime,omitempty"` /*  创建时间  */
	Eips         []*CtvpcListSnatsReturnObjResultsEipsResponse `json:"eips"`                   /*  绑定的 eip 信息  */
	SubnetID     *string                                       `json:"subnetID,omitempty"`     /*  子网 ID  */
	NatGatewayID *string                                       `json:"natGatewayID,omitempty"` /*  nat 网关 ID  */
}

type CtvpcListSnatsReturnObjResultsEipsResponse struct {
	EipID     *string `json:"eipID,omitempty"`     /*  弹性 IP id  */
	IpAddress *string `json:"ipAddress,omitempty"` /*  弹性 IP 地址  */
}
