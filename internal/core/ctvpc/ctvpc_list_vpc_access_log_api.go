package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcListVpcAccessLogApi
/* 查询VPC访问日志
 */type CtvpcListVpcAccessLogApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListVpcAccessLogApi(client *core.CtyunClient) *CtvpcListVpcAccessLogApi {
	return &CtvpcListVpcAccessLogApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/log/list-vpc-accesslog",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListVpcAccessLogApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListVpcAccessLogRequest) (*CtvpcListVpcAccessLogResponse, error) {
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
	var resp CtvpcListVpcAccessLogResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListVpcAccessLogRequest struct {
	ClientToken *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  `json:"regionID,omitempty"`    /*  区域ID  */
	QueryText   *string `json:"queryText,omitempty"`   /*  查询文本,名称或流日志ID  */
	Page        int32   `json:"page"`                  /*  页号,page>=1  */
	PageSize    int32   `json:"pageSize"`              /*  每页数量,0<=pageSize<=50  */
}

type CtvpcListVpcAccessLogResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcListVpcAccessLogReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtvpcListVpcAccessLogReturnObjResponse struct {
	TotalCount int32 `json:"totalCount"` /*  总共数量  */
	TotalPage  int32 `json:"totalPage"`  /*  总页数  */
}
