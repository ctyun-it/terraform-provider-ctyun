package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCheckLogAuthApi
/* 检查日志权限
 */type CtvpcCheckLogAuthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCheckLogAuthApi(client *core.CtyunClient) *CtvpcCheckLogAuthApi {
	return &CtvpcCheckLogAuthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/log/check-log-auth",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCheckLogAuthApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCheckLogAuthRequest) (*CtvpcCheckLogAuthResponse, error) {
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
	var resp CtvpcCheckLogAuthResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCheckLogAuthRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  区域ID  */
}

type CtvpcCheckLogAuthResponse struct {
	StatusCode  int32                               `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                             `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                             `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                             `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCheckLogAuthReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtvpcCheckLogAuthReturnObjResponse struct {
	Check_result *bool `json:"check_result"` /*  true:已开通,false:未开通  */
}
