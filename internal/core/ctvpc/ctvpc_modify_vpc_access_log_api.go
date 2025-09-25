package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcModifyVpcAccessLogApi
/* 修改VPC访问日志
 */type CtvpcModifyVpcAccessLogApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcModifyVpcAccessLogApi(client *core.CtyunClient) *CtvpcModifyVpcAccessLogApi {
	return &CtvpcModifyVpcAccessLogApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/log/modify-vpc-accesslog",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcModifyVpcAccessLogApi) Do(ctx context.Context, credential core.Credential, req *CtvpcModifyVpcAccessLogRequest) (*CtvpcModifyVpcAccessLogResponse, error) {
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
	var resp CtvpcModifyVpcAccessLogResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcModifyVpcAccessLogRequest struct {
	ClientToken    *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID       string  `json:"regionID,omitempty"`    /*  区域ID  */
	FlowID         string  `json:"flowID,omitempty"`      /*  访问日志ID  */
	LogCollection  int32   `json:"logCollection"`         /*  0:关闭日志收集 1:启用日志收集  */
	SampleInterval int32   `json:"sampleInterval"`        /*  采样间隔 60或300  */
	Name           *string `json:"name,omitempty"`        /*  名称  */
}

type CtvpcModifyVpcAccessLogResponse struct {
	StatusCode  int32                                       `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcModifyVpcAccessLogReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtvpcModifyVpcAccessLogReturnObjResponse struct {
	ID *string `json:"ID,omitempty"` /*  流日志ID  */
}
