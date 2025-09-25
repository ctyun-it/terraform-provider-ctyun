package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcGwlbShowTargetApi
/* 查看target详情
 */type CtvpcGwlbShowTargetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGwlbShowTargetApi(client *core.CtyunClient) *CtvpcGwlbShowTargetApi {
	return &CtvpcGwlbShowTargetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/gwlb/show-target",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGwlbShowTargetApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGwlbShowTargetRequest) (*CtvpcGwlbShowTargetResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("targetID", req.TargetID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcGwlbShowTargetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGwlbShowTargetRequest struct {
	RegionID string /*  资源池 ID  */
	TargetID string /*  后端服务 ID  */
}

type CtvpcGwlbShowTargetResponse struct {
	StatusCode  int32                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGwlbShowTargetReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcGwlbShowTargetReturnObjResponse struct {
	TargetID              *string `json:"targetID,omitempty"`              /*  后端服务ID  */
	TargetGroupID         *string `json:"targetGroupID,omitempty"`         /*  后端服务组ID  */
	InstanceType          *string `json:"instanceType,omitempty"`          /*  实例类型，取值有: VM / BMS/ CBM  */
	InstanceID            *string `json:"instanceID,omitempty"`            /*  实例 ID  */
	InstanceVpc           *string `json:"instanceVpc,omitempty"`           /*  实例所在的 vpc  */
	Weight                int32   `json:"weight"`                          /*  权重  */
	HealthCheckStatus     *string `json:"healthCheckStatus,omitempty"`     /*  ipv4 健康检查状态，取值: unknown / online / offline  */
	HealthCheckStatusIpv6 *string `json:"healthCheckStatusIpv6,omitempty"` /*  ipv6 健康检查状态，取值: unknown / online / offline  */
	CreatedAt             *string `json:"createdAt,omitempty"`             /*  创建时间  */
	UpdatedAt             *string `json:"updatedAt,omitempty"`             /*  更新时间  */
}
