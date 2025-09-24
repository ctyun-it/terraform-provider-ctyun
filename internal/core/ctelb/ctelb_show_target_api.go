package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbShowTargetApi
/* 查看后端服务详情
 */type CtelbShowTargetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbShowTargetApi(client *core.CtyunClient) *CtelbShowTargetApi {
	return &CtelbShowTargetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/show-target",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbShowTargetApi) Do(ctx context.Context, credential core.Credential, req *CtelbShowTargetRequest) (*CtelbShowTargetResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ID != "" {
		ctReq.AddParam("ID", req.ID)
	}
	ctReq.AddParam("targetID", req.TargetID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbShowTargetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbShowTargetRequest struct {
	RegionID string /*  区域ID  */
	ID       string /*  后端服务ID, 该字段后续废弃  */
	TargetID string /*  后端服务ID, 推荐使用该字段, 当同时使用 ID 和 targetID 时，优先使用 targetID  */
}

type CtelbShowTargetResponse struct {
	StatusCode  int32                             `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                            `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                            `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                            `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbShowTargetReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbShowTargetReturnObjResponse struct {
	RegionID              string `json:"regionID,omitempty"`              /*  区域ID  */
	AzName                string `json:"azName,omitempty"`                /*  可用区名称  */
	ProjectID             string `json:"projectID,omitempty"`             /*  项目ID  */
	ID                    string `json:"ID,omitempty"`                    /*  后端服务ID  */
	TargetGroupID         string `json:"targetGroupID,omitempty"`         /*  后端服务组ID  */
	Description           string `json:"description,omitempty"`           /*  描述  */
	InstanceType          string `json:"instanceType,omitempty"`          /*  实例类型: VM / BM  */
	InstanceID            string `json:"instanceID,omitempty"`            /*  实例ID  */
	ProtocolPort          int32  `json:"protocolPort,omitempty"`          /*  协议端口  */
	Weight                int32  `json:"weight,omitempty"`                /*  权重  */
	HealthCheckStatus     string `json:"healthCheckStatus,omitempty"`     /*  IPv4的健康检查状态: offline / online / unknown  */
	HealthCheckStatusIpv6 string `json:"healthCheckStatusIpv6,omitempty"` /*  IPv6的健康检查状态: offline / online / unknown  */
	Status                string `json:"status,omitempty"`                /*  状态: DOWN / ACTIVE  */
	CreatedTime           string `json:"createdTime,omitempty"`           /*  创建时间，为UTC格式  */
	UpdatedTime           string `json:"updatedTime,omitempty"`           /*  更新时间，为UTC格式  */
}
