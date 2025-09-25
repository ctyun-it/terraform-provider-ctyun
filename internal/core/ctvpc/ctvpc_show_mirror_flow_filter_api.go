package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowMirrorFlowFilterApi
/* 查看过滤条件详情
 */type CtvpcShowMirrorFlowFilterApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowMirrorFlowFilterApi(client *core.CtyunClient) *CtvpcShowMirrorFlowFilterApi {
	return &CtvpcShowMirrorFlowFilterApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/mirrorflow/show-filter",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowMirrorFlowFilterApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowMirrorFlowFilterRequest) (*CtvpcShowMirrorFlowFilterResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("mirrorFilterID", req.MirrorFilterID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowMirrorFlowFilterResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowMirrorFlowFilterRequest struct {
	RegionID       string /*  区域ID  */
	MirrorFilterID string /*  过滤条件 ID  */
}

type CtvpcShowMirrorFlowFilterResponse struct {
	StatusCode  int32                                       `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowMirrorFlowFilterReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                     `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowMirrorFlowFilterReturnObjResponse struct {
	Name            *string `json:"name,omitempty"`           /*  流量镜像名称  */
	Description     *string `json:"description,omitempty"`    /*  流量镜像描述  */
	MirrorFilterID  *string `json:"mirrorFilterID,omitempty"` /*  流量镜像ID  */
	CreatedTime     *string `json:"createdTime,omitempty"`    /*  创建时间  */
	InRuleCount     int32   `json:"inRuleCount"`              /*  出方向规则数  */
	OutRuleCount    int32   `json:"outRuleCount"`             /*  入方向规则数  */
	MirrorFlowCount int32   `json:"mirrorFlowCount"`          /*  关联流量镜像会话数  */
}
