package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVniaCreateNetworkPathAnalysisApi
/* 创建网络路径分析
 */type CtvpcVniaCreateNetworkPathAnalysisApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVniaCreateNetworkPathAnalysisApi(client *core.CtyunClient) *CtvpcVniaCreateNetworkPathAnalysisApi {
	return &CtvpcVniaCreateNetworkPathAnalysisApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vnia/create-network-path-analysis",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVniaCreateNetworkPathAnalysisApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVniaCreateNetworkPathAnalysisRequest) (*CtvpcVniaCreateNetworkPathAnalysisResponse, error) {
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
	var resp CtvpcVniaCreateNetworkPathAnalysisResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVniaCreateNetworkPathAnalysisRequest struct {
	RegionID      string `json:"regionID,omitempty"`      /*  资源池 ID  */
	NetworkPathID string `json:"networkPathID,omitempty"` /*  网络路径 ID  */
}

type CtvpcVniaCreateNetworkPathAnalysisResponse struct {
	StatusCode  int32                                                `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                              `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                              `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                              `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcVniaCreateNetworkPathAnalysisReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVniaCreateNetworkPathAnalysisReturnObjResponse struct {
	AnalysisID *string `json:"analysisID,omitempty"` /*  路径分析 ID  */
}
