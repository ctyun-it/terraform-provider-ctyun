package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcAddVpcAccessLogApi
/* 创建VPC访问日志
 */type CtvpcAddVpcAccessLogApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcAddVpcAccessLogApi(client *core.CtyunClient) *CtvpcAddVpcAccessLogApi {
	return &CtvpcAddVpcAccessLogApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/log/add-vpc-accesslog",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcAddVpcAccessLogApi) Do(ctx context.Context, credential core.Credential, req *CtvpcAddVpcAccessLogRequest) (*CtvpcAddVpcAccessLogResponse, error) {
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
	var resp CtvpcAddVpcAccessLogResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcAddVpcAccessLogRequest struct {
	ClientToken    *string `json:"clientToken,omitempty"`    /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID       string  `json:"regionID,omitempty"`       /*  区域ID  */
	Name           string  `json:"name,omitempty"`           /*  名称  */
	ResourceID     string  `json:"resourceID,omitempty"`     /*  资源ID  */
	ResourceType   string  `json:"resourceType,omitempty"`   /*  资源类型 1:vpc 2:subnet 3:port 区分大小写  */
	FlowType       int32   `json:"flowType"`                 /*  流量类型: 0:全部流量 1:被访问控制允许的流量 2:被访问控制拒绝的流量, 默认为0  */
	LogProjectCode string  `json:"logProjectCode,omitempty"` /*  日志项目ID  */
	LogUnitCode    string  `json:"logUnitCode,omitempty"`    /*  日志单元ID  */
	LogCollection  int32   `json:"logCollection"`            /*  0:关闭日志收集 1:启用日志收集,默认0  */
	SampleInterval int32   `json:"sampleInterval"`           /*  采样间隔,单位s 60 或 300  */
	ProjectID      *string `json:"projectID,omitempty"`      /*  企业项目ID,默认0  */
	Description    *string `json:"description,omitempty"`    /*  描述  */
}

type CtvpcAddVpcAccessLogResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcAddVpcAccessLogReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtvpcAddVpcAccessLogReturnObjResponse struct {
	Id *string `json:"id,omitempty"` /*  访问日志ID  */
}
