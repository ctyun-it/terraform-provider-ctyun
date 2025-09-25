package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcListLogProjectApi
/* 查询日志项目
 */type CtvpcListLogProjectApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListLogProjectApi(client *core.CtyunClient) *CtvpcListLogProjectApi {
	return &CtvpcListLogProjectApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/log/query-log-project-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListLogProjectApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListLogProjectRequest) (*CtvpcListLogProjectResponse, error) {
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
	var resp CtvpcListLogProjectResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListLogProjectRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  区域ID  */
}

type CtvpcListLogProjectResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcListLogProjectReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtvpcListLogProjectReturnObjResponse struct {
	ProjectID   *string `json:"projectID,omitempty"`   /*  日志项目ID  */
	ProjectName *string `json:"projectName,omitempty"` /*  日志项目名称  */
	AliasName   *string `json:"aliasName,omitempty"`   /*  日志项目别名  */
	Description *string `json:"description,omitempty"` /*  描述  */
}
