package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetSdwanOverloadApi
/* 过载保护列表查询 */
type SdwanGetSdwanOverloadApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanOverloadApi(client *core.CtyunClient) *SdwanGetSdwanOverloadApi {
	return &SdwanGetSdwanOverloadApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/overload/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanOverloadApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanOverloadRequest) (*SdwanGetSdwanOverloadResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.EdgeID != nil && *req.EdgeID != "" {
		ctReq.AddParam("edgeID", *req.EdgeID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanOverloadResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanOverloadRequest struct {
	EdgeID *string `json:"edgeID,omitempty"` /*  edge id  */
}

type SdwanGetSdwanOverloadResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                 `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                 `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                 `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetSdwanOverloadReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetSdwanOverloadReturnObjResponse struct {
	TotalCount   int32                                           `json:"totalCount"`   /*  总数  */
	CurrentCount int32                                           `json:"currentCount"` /*  当前页数量  */
	Code         *string                                         `json:"code"`         /*  code  */
	Result       []*SdwanGetSdwanOverloadReturnObjResultResponse `json:"result"`       /*  列表  */
}

type SdwanGetSdwanOverloadReturnObjResultResponse struct {
	Port          *string `json:"port"`          /*  本参数表示端口名称<br/>取值范围:<br/>WAN1:WAN1<br/>WAN2:WAN2<br/>LAN:LAN<br/>LTE:LTE  */
	TransportNet  *string `json:"transportNet"`  /*  本参数表示传输网络类型<br/><br/>取值范围:<br/>internet:internet<br/>mpls:mpls<br/>lte:lte  */
	LinkDetection *string `json:"linkDetection"` /*  本参数表示链路探测目标IP，多个用逗号分割  */
	InRate        int32   `json:"inRate"`        /*  输入速率，单位mbps  */
	OutRate       int32   `json:"outRate"`       /*  输出速率，单位mbps  */
	InProtect     int32   `json:"inProtect"`     /*  输入保护速率，单位mbps  */
	OutProtect    int32   `json:"outProtect"`    /*  输出保护速率，单位mbps  */
}
