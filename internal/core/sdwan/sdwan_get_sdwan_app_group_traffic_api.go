package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetSdwanAppGroupTrafficApi
/* 查询所有应用组的总流量 */
type SdwanGetSdwanAppGroupTrafficApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanAppGroupTrafficApi(client *core.CtyunClient) *SdwanGetSdwanAppGroupTrafficApi {
	return &SdwanGetSdwanAppGroupTrafficApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/app-group/list-traffic",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanAppGroupTrafficApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanAppGroupTrafficRequest) (*SdwanGetSdwanAppGroupTrafficResponse, error) {
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
	if req.TimeRange != nil && *req.TimeRange != "" {
		ctReq.AddParam("timeRange", *req.TimeRange)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanAppGroupTrafficResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanAppGroupTrafficRequest struct {
	Sn        *string `json:"sn,omitempty"`        /*  设备序列号  */
	StartTime *string `json:"startTime,omitempty"` /*  开始时间 时间戳 秒  */
	EndTime   *string `json:"endTime,omitempty"`   /*  结束时间 时间戳 秒  */
	TimeRange *string `json:"timeRange,omitempty"` /*  时间范围 当不传time_from和time_till时必传，认为当前时间为结束时间 单位是分钟  */
}

type SdwanGetSdwanAppGroupTrafficResponse struct {
	StatusCode  int32                                          `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                        `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                        `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                        `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetSdwanAppGroupTrafficReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                        `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetSdwanAppGroupTrafficReturnObjResponse struct {
	TotalCount   int32                                                  `json:"totalCount"`   /*  总数  */
	CurrentCount int32                                                  `json:"currentCount"` /*  当前页数量  */
	Code         *string                                                `json:"code"`         /*  code  */
	Result       []*SdwanGetSdwanAppGroupTrafficReturnObjResultResponse `json:"result"`       /*  流量列表  */
}

type SdwanGetSdwanAppGroupTrafficReturnObjResultResponse struct {
	GroupName *string `json:"groupName"` /*  应用组名称  */
	Traffic   int32   `json:"traffic"`   /*  流量大小 单位是bytes  */
}
