package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetSdwanAppTrafficApi
/* 查询top5应用的流量 */
type SdwanGetSdwanAppTrafficApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanAppTrafficApi(client *core.CtyunClient) *SdwanGetSdwanAppTrafficApi {
	return &SdwanGetSdwanAppTrafficApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/app/list-traffic",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanAppTrafficApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanAppTrafficRequest) (*SdwanGetSdwanAppTrafficResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.Sn != nil && *req.Sn != "" {
		ctReq.AddParam("sn", *req.Sn)
	}
	if req.StartTime != nil && *req.StartTime != "" {
		ctReq.AddParam("startTime", *req.StartTime)
	}
	if req.EndTime != nil && *req.EndTime != "" {
		ctReq.AddParam("endTime", *req.EndTime)
	}
	if req.Direction != nil && *req.Direction != "" {
		ctReq.AddParam("direction", *req.Direction)
	}
	if req.Limit != nil && *req.Limit != "" {
		ctReq.AddParam("limit", *req.Limit)
	}
	if req.TimeRange != nil && *req.TimeRange != "" {
		ctReq.AddParam("timeRange", *req.TimeRange)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanAppTrafficResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanAppTrafficRequest struct {
	Sn        *string `json:"sn,omitempty"`        /*  设备序列号  */
	StartTime *string `json:"startTime,omitempty"` /*  开始时间 时间戳 秒  */
	EndTime   *string `json:"endTime,omitempty"`   /*  结束时间 时间戳 秒  */
	Direction *string `json:"direction,omitempty"` /*  本参数表示控制方向<br/><br/>取值范围:<br/>in:入方向<br/>out:出方向  */
	Limit     *string `json:"limit,omitempty"`     /*  限制个数 默认值: 5  ，默认值:5  */
	TimeRange *string `json:"timeRange,omitempty"` /*  时间范围 当不传time_from和time_till时必传，认为当前时间为结束时间 单位是分钟  */
}

type SdwanGetSdwanAppTrafficResponse struct {
	StatusCode  int32                                     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetSdwanAppTrafficReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetSdwanAppTrafficReturnObjResponse struct {
	TotalCount   int32                                             `json:"totalCount"`   /*  总数  */
	CurrentCount int32                                             `json:"currentCount"` /*  当前页数量  */
	Code         *string                                           `json:"code"`         /*  状态码  */
	Result       []*SdwanGetSdwanAppTrafficReturnObjResultResponse `json:"result"`       /*  列表  */
}

type SdwanGetSdwanAppTrafficReturnObjResultResponse struct {
	AppName *string `json:"appName"` /*  应用名称  */
	Traffic int32   `json:"traffic"` /*  流量大小单位是M  */
}
