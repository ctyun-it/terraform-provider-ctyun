package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanQueryAlarmItemApi
/* 查询监控项 */
type SdwanSdwanQueryAlarmItemApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanQueryAlarmItemApi(client *core.CtyunClient) *SdwanSdwanQueryAlarmItemApi {
	return &SdwanSdwanQueryAlarmItemApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/alarm-items/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanQueryAlarmItemApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanQueryAlarmItemRequest) (*SdwanSdwanQueryAlarmItemResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("serialNumbers", req.SerialNumbers)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanSdwanQueryAlarmItemResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanQueryAlarmItemRequest struct {
	SerialNumbers string `json:"serialNumbers"` /*  sn列表，多个盒子用逗号分割  */
}

type SdwanSdwanQueryAlarmItemResponse struct {
	StatusCode  int32                                      `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                    `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                    `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                    `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanQueryAlarmItemReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                    `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanQueryAlarmItemReturnObjResponse struct {
	Result       []*SdwanSdwanQueryAlarmItemReturnObjResultResponse `json:"result"`       /*  获取盒子的监控项  */
	TotalCount   int32                                              `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                              `json:"currentCount"` /*  当前页数量  */
}

type SdwanSdwanQueryAlarmItemReturnObjResultResponse struct {
	AlarmItemName        *string `json:"alarmItemName"`        /*  监控项名称  */
	AlarmItemUnit        *string `json:"alarmItemUnit"`        /*  监控项单位  */
	AlarmItemDescription *string `json:"alarmItemDescription"` /*  监控项描述  */
}
