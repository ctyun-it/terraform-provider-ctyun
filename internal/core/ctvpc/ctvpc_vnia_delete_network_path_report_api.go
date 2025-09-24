package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVniaDeleteNetworkPathReportApi
/* 删除网络路径分析报告
 */type CtvpcVniaDeleteNetworkPathReportApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVniaDeleteNetworkPathReportApi(client *core.CtyunClient) *CtvpcVniaDeleteNetworkPathReportApi {
	return &CtvpcVniaDeleteNetworkPathReportApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vnia/delete-network-path-report",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVniaDeleteNetworkPathReportApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVniaDeleteNetworkPathReportRequest) (*CtvpcVniaDeleteNetworkPathReportResponse, error) {
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
	var resp CtvpcVniaDeleteNetworkPathReportResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVniaDeleteNetworkPathReportRequest struct {
	RegionID            string `json:"regionID,omitempty"`            /*  资源池 ID  */
	NetworkPathReportID string `json:"networkPathReportID,omitempty"` /*  路径分析报告 ID  */
}

type CtvpcVniaDeleteNetworkPathReportResponse struct {
	StatusCode  int32                                              `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                            `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                            `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                            `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcVniaDeleteNetworkPathReportReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVniaDeleteNetworkPathReportReturnObjResponse struct {
	NetworkPathReportID *string `json:"networkPathReportID,omitempty"` /*  路径分析报告 ID  */
}
