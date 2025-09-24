package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcGwlbShowTargetGroupApi
/* 查看target_group详情
 */type CtvpcGwlbShowTargetGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGwlbShowTargetGroupApi(client *core.CtyunClient) *CtvpcGwlbShowTargetGroupApi {
	return &CtvpcGwlbShowTargetGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/gwlb/show-target-group",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGwlbShowTargetGroupApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGwlbShowTargetGroupRequest) (*CtvpcGwlbShowTargetGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("targetGroupID", req.TargetGroupID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcGwlbShowTargetGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGwlbShowTargetGroupRequest struct {
	RegionID      string /*  资源池 ID  */
	TargetGroupID string /*  后端服务组 ID  */
}

type CtvpcGwlbShowTargetGroupResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGwlbShowTargetGroupReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcGwlbShowTargetGroupReturnObjResponse struct {
	TargetGroupID     *string `json:"targetGroupID,omitempty"` /*  后端服务组ID  */
	Name              *string `json:"name,omitempty"`          /*  名称  */
	Description       *string `json:"description,omitempty"`   /*  描述  */
	VpcID             *string `json:"vpcID,omitempty"`         /*  vpc id  */
	HealthCheckID     *string `json:"healthCheckID,omitempty"` /*  健康检查 ID  */
	SessionStickyMode int32   `json:"sessionStickyMode"`       /*  0:五元组, 4:二元组, 5:三元组, -1 未开启  */
	FailoverType      int32   `json:"failoverType"`            /*  故障转移类型 1 表示关闭，2 表示再平衡  */
	BypassType        int32   `json:"bypassType"`              /*  旁路类型 1 表示关闭，2 表示自动  */
}
