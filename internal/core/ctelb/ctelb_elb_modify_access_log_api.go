package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbElbModifyAccessLogApi
/* 修改elb流日志
 */type CtelbElbModifyAccessLogApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbElbModifyAccessLogApi(client *core.CtyunClient) *CtelbElbModifyAccessLogApi {
	return &CtelbElbModifyAccessLogApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/modify-access-log",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbElbModifyAccessLogApi) Do(ctx context.Context, credential core.Credential, req *CtelbElbModifyAccessLogRequest) (*CtelbElbModifyAccessLogResponse, error) {
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
	var resp CtelbElbModifyAccessLogResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbElbModifyAccessLogRequest struct {
	ClientToken   string `json:"clientToken,omitempty"`   /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID      string `json:"regionID,omitempty"`      /*  区域ID  */
	AccesslogID   string `json:"accesslogID,omitempty"`   /*  访问日志ID  */
	LogCollection int32  `json:"logCollection,omitempty"` /*  0:关闭日志收集 1:启用日志收集  */
	ProjectCode   string `json:"projectCode,omitempty"`   /*  日志项目code  */
	UnitCode      string `json:"unitCode,omitempty"`      /*  日志单元code  */
}

type CtelbElbModifyAccessLogResponse struct {
	StatusCode  int32                                       `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbElbModifyAccessLogReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtelbElbModifyAccessLogReturnObjResponse struct {
	Code string `json:"Code,omitempty"` /*  返回码  */
}
