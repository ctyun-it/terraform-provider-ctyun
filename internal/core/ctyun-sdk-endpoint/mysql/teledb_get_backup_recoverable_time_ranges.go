package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetRecoverableTimeRangesApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetRecoverableTimeRangesApi(client *ctyunsdk.CtyunClient) *TeledbGetRecoverableTimeRangesApi {
	return &TeledbGetRecoverableTimeRangesApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v2/open-api/recovery/recoveryable-time-range",
		},
	}
}

func (this *TeledbGetRecoverableTimeRangesApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetRecoverableTimeRangesRequest, header *TeledbGetRecoverableTimeRangesRequestHeader) (GetRecoverableTimeRangesResp *TeledbGetRecoverableTimeRangesResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}
	if req.OuterProdInstId == "" || header.InstID == "" {
		err = errors.New("instId 为空")
		return
	}
	if header.RegionID == "" {
		err = fmt.Errorf("region_id is required")
		return
	}
	if header.InstID == "" {
		err = fmt.Errorf("inst_id is required")
		return
	}
	builder.AddHeader("inst-id", header.InstID)
	builder.AddHeader("regionId", header.RegionID)

	builder.AddParam("outerProdInstId", req.OuterProdInstId)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	GetRecoverableTimeRangesResp = &TeledbGetRecoverableTimeRangesResponse{}
	err = resp.Parse(GetRecoverableTimeRangesResp)
	if err != nil {
		return
	}
	return GetRecoverableTimeRangesResp, nil
}

type TeledbGetRecoverableTimeRangesRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` // 实例ID，必填
}
type TeledbGetRecoverableTimeRangesRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}
type TeledbGetRecoverableTimeRangesResponse struct {
	StatusCode int32                                            `json:"statusCode"`      // 接口状态码
	Error      *string                                          `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                           `json:"message"`         // 描述信息
	ReturnObj  *TeledbGetRecoverableTimeRangesResponseReturnObj `json:"returnObj"`
}

type TeledbGetRecoverableTimeRangesResponseReturnObjTimeRange struct {
	StartTimestamp string `json:"startTimestamp"`
	EndTimestamp   string `json:"endTimestamp"`
}

type TeledbGetRecoverableTimeRangesResponseReturnObj struct {
	Data []TeledbGetRecoverableTimeRangesResponseReturnObjTimeRange `json:"data"`
}
