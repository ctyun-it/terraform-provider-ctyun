package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcListLogUnitApi
/* 查询日志单元
 */type CtvpcListLogUnitApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListLogUnitApi(client *core.CtyunClient) *CtvpcListLogUnitApi {
	return &CtvpcListLogUnitApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/log/query-log-unit-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListLogUnitApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListLogUnitRequest) (*CtvpcListLogUnitResponse, error) {
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
	var resp CtvpcListLogUnitResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListLogUnitRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  区域ID  */
	ProjectCode string `json:"projectCode,omitempty"` /*  日志项目Code  */
	Page        int32  `json:"page"`                  /*  页数page=1开始 默认page=1  */
	PageSize    string `json:"pageSize,omitempty"`    /*  每页数量 pageSize<=300 默认pageSize=300  */
}

type CtvpcListLogUnitResponse struct {
	StatusCode  int32                                `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                              `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                              `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                              `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcListLogUnitReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	TotalCount  int32                                `json:"totalCount"`            /*  总数  */
	TotalPage   int32                                `json:"totalPage"`             /*  总页数  */
}

type CtvpcListLogUnitReturnObjResponse struct {
	UnitCode    *string `json:"unitCode,omitempty"`    /*  日志单元code  */
	UnitName    *string `json:"unitName,omitempty"`    /*  日志单元名称  */
	AliasName   *string `json:"aliasName,omitempty"`   /*  日志单元别名  */
	Description *string `json:"description,omitempty"` /*  描述  */
	ProjectCode *string `json:"projectCode,omitempty"` /*  日志项目code  */
	ProjectName *string `json:"projectName,omitempty"` /*  日志项目名称  */
}
