package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVniaShowInstanceDiagnosisApi
/* 获取实例诊断详情
 */type CtvpcVniaShowInstanceDiagnosisApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVniaShowInstanceDiagnosisApi(client *core.CtyunClient) *CtvpcVniaShowInstanceDiagnosisApi {
	return &CtvpcVniaShowInstanceDiagnosisApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vnia/show-instance-diagnosis",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVniaShowInstanceDiagnosisApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVniaShowInstanceDiagnosisRequest) (*CtvpcVniaShowInstanceDiagnosisResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("resourceID", req.ResourceID)
	ctReq.AddParam("resourceType", req.ResourceType)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcVniaShowInstanceDiagnosisResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVniaShowInstanceDiagnosisRequest struct {
	RegionID     string /*  资源池 ID  */
	ResourceID   string /*  资源 ID  */
	ResourceType string /*  资源类型, 仅支持 eip / natgw / elb  */
}

type CtvpcVniaShowInstanceDiagnosisResponse struct {
	StatusCode  int32                                            `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                          `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                          `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                          `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcVniaShowInstanceDiagnosisReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVniaShowInstanceDiagnosisReturnObjResponse struct {
	InstanceDiagnosisID *string `json:"instanceDiagnosisID,omitempty"` /*  实例诊断 ID  */
	CreatedAt           *string `json:"createdAt,omitempty"`           /*  创建时间  */
	UpdatedAt           *string `json:"updatedAt,omitempty"`           /*  更新时间  */
	ResourceType        *string `json:"resourceType,omitempty"`        /*  资源类型, 仅支持 eip / natgw / elb  */
	DiagnosisStatus     *string `json:"diagnosisStatus,omitempty"`     /*  诊断状态  */
}
