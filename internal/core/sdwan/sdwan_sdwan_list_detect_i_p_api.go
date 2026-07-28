package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanListDetectIPApi
/* 查询链路检测目的IP列表 */
type SdwanSdwanListDetectIPApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanListDetectIPApi(client *core.CtyunClient) *SdwanSdwanListDetectIPApi {
	return &SdwanSdwanListDetectIPApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/detect-ip/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanListDetectIPApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanListDetectIPRequest) (*SdwanSdwanListDetectIPResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	ctReq.AddParam("edgePort", req.EdgePort)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanSdwanListDetectIPResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanListDetectIPRequest struct {
	EdgeID   string `json:"edgeID"`   /*  智能网关ID  */
	EdgePort string `json:"edgePort"` /*  智能网关端口  */
}

type SdwanSdwanListDetectIPResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanListDetectIPReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanListDetectIPReturnObjResponse struct {
	Result       []*SdwanSdwanListDetectIPReturnObjResultResponse `json:"result"`       /*  返回结果  */
	TotalCount   int32                                            `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                            `json:"currentCount"` /*  当前页数量  */
}

type SdwanSdwanListDetectIPReturnObjResultResponse struct {
	LinkDetectionList []*string `json:"linkDetectionList"` /*  链路检测目的IP  */
}
