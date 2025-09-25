package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVniaShowInstanceDiagnosisRecordApi
/* 获取实例诊断记录详情
 */type CtvpcVniaShowInstanceDiagnosisRecordApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVniaShowInstanceDiagnosisRecordApi(client *core.CtyunClient) *CtvpcVniaShowInstanceDiagnosisRecordApi {
	return &CtvpcVniaShowInstanceDiagnosisRecordApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vnia/show-instance-diagnosis-record",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVniaShowInstanceDiagnosisRecordApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVniaShowInstanceDiagnosisRecordRequest) (*CtvpcVniaShowInstanceDiagnosisRecordResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("diagnosisRecordID", req.DiagnosisRecordID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcVniaShowInstanceDiagnosisRecordResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVniaShowInstanceDiagnosisRecordRequest struct {
	RegionID          string /*  资源池 ID  */
	DiagnosisRecordID string /*  实例诊断记录 ID  */
}

type CtvpcVniaShowInstanceDiagnosisRecordResponse struct {
	StatusCode  int32                                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcVniaShowInstanceDiagnosisRecordReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVniaShowInstanceDiagnosisRecordReturnObjResponse struct {
	CreatedAt         *string `json:"createdAt,omitempty"`         /*  创建时间  */
	UpdatedAt         *string `json:"updatedAt,omitempty"`         /*  更新时间  */
	DiagnosisRecordID *string `json:"diagnosisRecordID,omitempty"` /*  实例诊断记录 ID  */
	ResourceID        *string `json:"resourceID,omitempty"`        /*  资源 ID  */
	ResourceType      *string `json:"resourceType,omitempty"`      /*  资源类型, 仅支持 eip / natgw / elb  */
	ErrMsg            *string `json:"errMsg,omitempty"`            /*  错误信息  */
	DiagnosisStatus   *string `json:"diagnosisStatus,omitempty"`   /*  诊断状态  */
	HealthStatus      int32   `json:"healthStatus"`                /*  健康检查状态  */
}
