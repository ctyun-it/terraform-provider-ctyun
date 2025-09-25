package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVniaDeleteInstanceDiagnosisApi
/* 删除实例诊断
 */type CtvpcVniaDeleteInstanceDiagnosisApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVniaDeleteInstanceDiagnosisApi(client *core.CtyunClient) *CtvpcVniaDeleteInstanceDiagnosisApi {
	return &CtvpcVniaDeleteInstanceDiagnosisApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vnia/delete-instance-diagnosis",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVniaDeleteInstanceDiagnosisApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVniaDeleteInstanceDiagnosisRequest) (*CtvpcVniaDeleteInstanceDiagnosisResponse, error) {
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
	var resp CtvpcVniaDeleteInstanceDiagnosisResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVniaDeleteInstanceDiagnosisRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池 ID  */
	ResourceID   string `json:"resourceID,omitempty"`   /*  资源 ID  */
	ResourceType string `json:"resourceType,omitempty"` /*  资源类型, 仅支持 eip / natgw / elb  */
}

type CtvpcVniaDeleteInstanceDiagnosisResponse struct {
	StatusCode  int32                                              `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                            `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                            `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                            `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcVniaDeleteInstanceDiagnosisReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVniaDeleteInstanceDiagnosisReturnObjResponse struct {
	ResourceID   *string `json:"resourceID,omitempty"`   /*  资源 ID  */
	ResourceType *string `json:"resourceType,omitempty"` /*  资源类型, 仅支持 eip / natgw / elb  */
}
