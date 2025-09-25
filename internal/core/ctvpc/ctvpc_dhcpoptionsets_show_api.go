package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDhcpoptionsetsShowApi
/* 查询dhcpoptionsets详情
 */type CtvpcDhcpoptionsetsShowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDhcpoptionsetsShowApi(client *core.CtyunClient) *CtvpcDhcpoptionsetsShowApi {
	return &CtvpcDhcpoptionsetsShowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/dhcpoptionsets/show",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDhcpoptionsetsShowApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDhcpoptionsetsShowRequest) (*CtvpcDhcpoptionsetsShowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("dhcpOptionSetsID", req.DhcpOptionSetsID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcDhcpoptionsetsShowResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDhcpoptionsetsShowRequest struct {
	RegionID         string /*  资源池 ID  */
	DhcpOptionSetsID string /*  集合ID  */
}

type CtvpcDhcpoptionsetsShowResponse struct {
	StatusCode  int32                                     `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcDhcpoptionsetsShowReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcDhcpoptionsetsShowReturnObjResponse struct {
	DhcpOptionSetsID *string   `json:"dhcpOptionSetsID,omitempty"` /*  dhcpoptionsets  ID  */
	Name             *string   `json:"name,omitempty"`             /*  名字  */
	Description      *string   `json:"description,omitempty"`      /*  描述  */
	DomainName       []*string `json:"domainName"`                 /*  域名  */
	DnsList          []*string `json:"dnsList"`                    /*  ip 列表  */
	VpcList          []*string `json:"vpcList"`                    /*  vpc 列表  */
	CreatedAt        *string   `json:"createdAt,omitempty"`        /*  创建时间  */
	UpdatedAt        *string   `json:"updatedAt,omitempty"`        /*  更新时间  */
}
