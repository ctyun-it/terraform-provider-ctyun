package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcListIPv4GwApi
/* 获取IPv4网关列表
 */type CtvpcListIPv4GwApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListIPv4GwApi(client *core.CtyunClient) *CtvpcListIPv4GwApi {
	return &CtvpcListIPv4GwApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/ipv4-gw/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListIPv4GwApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListIPv4GwRequest) (*CtvpcListIPv4GwResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.VpcID != nil {
		ctReq.AddParam("vpcID", *req.VpcID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListIPv4GwResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListIPv4GwRequest struct {
	RegionID string  /*  区域id  */
	VpcID    *string /*  关联的vpcID  */
}

type CtvpcListIPv4GwResponse struct {
	StatusCode  int32                               `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                             `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                             `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                             `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcListIPv4GwReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                             `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListIPv4GwReturnObjResponse struct {
	Name         *string `json:"name,omitempty"`         /*  IPv4网关名称  */
	Description  *string `json:"description,omitempty"`  /*  IPv4网关描述  */
	VpcID        *string `json:"vpcID,omitempty"`        /*  虚拟私有云 id  */
	Id           *string `json:"id,omitempty"`           /*  IPv4网关id  */
	RouteTableID *string `json:"routeTableID,omitempty"` /*  关联的网关路由表ID  */
	CreatedAt    *string `json:"createdAt,omitempty"`    /*  创建时间  */
	UpdatedAt    *string `json:"updatedAt,omitempty"`    /*  更新时间  */
}
