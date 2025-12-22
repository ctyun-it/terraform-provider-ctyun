package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanCreateLinkDetectApi
/* 开启链路监控 */
type SdwanSdwanCreateLinkDetectApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanCreateLinkDetectApi(client *core.CtyunClient) *SdwanSdwanCreateLinkDetectApi {
	return &SdwanSdwanCreateLinkDetectApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/link-detect/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanCreateLinkDetectApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanCreateLinkDetectRequest) (*SdwanSdwanCreateLinkDetectResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanCreateLinkDetectRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanSdwanCreateLinkDetectResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanCreateLinkDetectRequest struct {
	LinkID           string                                               `json:"linkID"`                     /*  链路Id  */
	DetectPolicyList []*SdwanSdwanCreateLinkDetectDetectPolicyListRequest `json:"detectPolicyList,omitempty"` /*  专线链路探测对象  */
}

type SdwanSdwanCreateLinkDetectDetectPolicyListRequest struct {
	Direction string `json:"direction"` /*  本参数表示探测方向<br/><br/>取值范围:<br/>to:下到源盒子<br/>from:下到目的盒子  */
	DetectIP  string `json:"detectIP"`  /*  探测IP  */
	NextHopIP string `json:"nextHopIP"` /*  下一跳地址  */
}

type SdwanSdwanCreateLinkDetectResponse struct {
	StatusCode  int32                                        `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                      `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                      `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                      `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanCreateLinkDetectReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                      `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanCreateLinkDetectReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作日志Id  */
}
