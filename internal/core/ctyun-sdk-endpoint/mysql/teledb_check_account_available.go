package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbCheckAccountAvailableApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbCheckAccountAvailableApi(client *ctyunsdk.CtyunClient) *TeledbCheckAccountAvailableApi {
	return &TeledbCheckAccountAvailableApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v1/open-api/account/account-available",
		},
	}
}

func (this *TeledbCheckAccountAvailableApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbCheckAccountAvailableRequest, header *TeledbCheckAccountAvailableRequestHeader) (CheckAccountAvailableResp *TeledbCheckAccountAvailableResponse, err error) {
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
	builder.AddHeader("inst-id", header.InstID)
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	CheckAccountAvailableResp = &TeledbCheckAccountAvailableResponse{}
	err = resp.Parse(CheckAccountAvailableResp)
	if err != nil {
		return
	}
	return CheckAccountAvailableResp, nil
}

type TeledbCheckAccountAvailableRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` // 实例ID，必填
	AccountName     string `json:"accountName"`     // 账户名称
}
type TeledbCheckAccountAvailableRequestHeader struct {
	ProjectID string `json:"projectID"`
	InstID    string `json:"instId"`    // 实例ID，必填
	RegionID  string `json:"region_id"` // 资源池ID，必填
}
type TeledbCheckAccountAvailableResponseReturnObj struct {
	Available bool `json:"available"`
}

type TeledbCheckAccountAvailableResponse struct {
	StatusCode int32                                         `json:"statusCode"`      // 接口状态码
	Error      *string                                       `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                        `json:"message"`         // 描述信息
	ReturnObj  *TeledbCheckAccountAvailableResponseReturnObj `json:"returnObj"`
}
