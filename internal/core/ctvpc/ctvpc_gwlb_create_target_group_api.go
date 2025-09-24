package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcGwlbCreateTargetGroupApi
/* 创建target_group
 */type CtvpcGwlbCreateTargetGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGwlbCreateTargetGroupApi(client *core.CtyunClient) *CtvpcGwlbCreateTargetGroupApi {
	return &CtvpcGwlbCreateTargetGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/gwlb/create-target-group",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGwlbCreateTargetGroupApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGwlbCreateTargetGroupRequest) (*CtvpcGwlbCreateTargetGroupResponse, error) {
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
	var resp CtvpcGwlbCreateTargetGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGwlbCreateTargetGroupRequest struct {
	RegionID          string  `json:"regionID,omitempty"`      /*  资源池 ID  */
	Name              string  `json:"name,omitempty"`          /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	VpcID             string  `json:"vpcID,omitempty"`         /*  虚拟私有云 ID  */
	HealthCheckID     *string `json:"healthCheckID,omitempty"` /*  健康检查 ID  */
	SessionStickyMode int32   `json:"sessionStickyMode"`       /*  流保持类型，0:五元组, 4:二元组, 5:三元组  */
	FailoverType      int32   `json:"failoverType"`            /*  故障转移类型 1 表示关闭，2 表示再平衡  */
}

type CtvpcGwlbCreateTargetGroupResponse struct {
	StatusCode  int32                                        `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGwlbCreateTargetGroupReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcGwlbCreateTargetGroupReturnObjResponse struct {
	TargetGroupID *string `json:"targetGroupID,omitempty"` /*  后端服务组ID  */
}
