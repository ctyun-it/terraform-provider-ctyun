package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbGwlbUpdateTargetGroupApi
/* 更新target_group
 */type CtelbGwlbUpdateTargetGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbGwlbUpdateTargetGroupApi(client *core.CtyunClient) *CtelbGwlbUpdateTargetGroupApi {
	return &CtelbGwlbUpdateTargetGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/gwlb/update-target-group",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbGwlbUpdateTargetGroupApi) Do(ctx context.Context, credential core.Credential, req *CtelbGwlbUpdateTargetGroupRequest) (*CtelbGwlbUpdateTargetGroupResponse, error) {
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
	var resp CtelbGwlbUpdateTargetGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbGwlbUpdateTargetGroupRequest struct {
	RegionID          string `json:"regionID,omitempty"`          /*  资源池 ID  */
	TargetGroupID     string `json:"targetGroupID,omitempty"`     /*  后端服务组 ID  */
	Name              string `json:"name,omitempty"`              /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	HealthCheckID     string `json:"healthCheckID,omitempty"`     /*  健康检查 ID，传空字符串时表示后端服务关闭健康检查  */
	SessionStickyMode int32  `json:"sessionStickyMode,omitempty"` /*  流保持类型,0:五元组, 4:二元组, 5:三元组  */
	FailoverType      int32  `json:"failoverType,omitempty"`      /*  故障转移类型 1 表示关闭，2 表示再平衡  */
	BypassType        int32  `json:"bypassType,omitempty"`        /*  旁路类型 1 表示关闭，2 表示自动  */
}

type CtelbGwlbUpdateTargetGroupResponse struct {
	StatusCode  int32                                        `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                       `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                       `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbGwlbUpdateTargetGroupReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbGwlbUpdateTargetGroupReturnObjResponse struct {
	TargetGroupID string `json:"targetGroupID,omitempty"` /*  后端服务组ID  */
}
