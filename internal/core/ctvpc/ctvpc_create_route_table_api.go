package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateRouteTableApi
/* 创建路由表
 */type CtvpcCreateRouteTableApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateRouteTableApi(client *core.CtyunClient) *CtvpcCreateRouteTableApi {
	return &CtvpcCreateRouteTableApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/route-table/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateRouteTableApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateRouteTableRequest) (*CtvpcCreateRouteTableResponse, error) {
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
	var resp CtvpcCreateRouteTableResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateRouteTableRequest struct {
	ClientToken string  `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  `json:"regionID,omitempty"`    /*  区域id  */
	VpcID       string  `json:"vpcID,omitempty"`       /*  关联的vpcID  */
	Name        string  `json:"name,omitempty"`        /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	ProjectID   *string `json:"projectID,omitempty"`   /*  企业项目 ID，默认为0  */
}

type CtvpcCreateRouteTableResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCreateRouteTableReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       *string                                 `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCreateRouteTableReturnObjResponse struct {
	Name            *string `json:"name,omitempty"`        /*  路由表名字  */
	Description     *string `json:"description,omitempty"` /*  路由表描述  */
	VpcID           *string `json:"vpcID,omitempty"`       /*  虚拟私有云 id  */
	Id              *string `json:"id,omitempty"`          /*  路由 id  */
	RouteRulesCount int32   `json:"routeRulesCount"`       /*  路由表中的路由数  */
	CreatedAt       *string `json:"createdAt,omitempty"`   /*  创建时间  */
	UpdatedAt       *string `json:"updatedAt,omitempty"`   /*  更新时间  */
}
