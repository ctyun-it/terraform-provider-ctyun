package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbCheckDatabaseNameAvailableApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbCheckDatabaseNameAvailableApi(client *ctyunsdk.CtyunClient) *TeledbCheckDatabaseNameAvailableApi {
	return &TeledbCheckDatabaseNameAvailableApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v1/open-api/account/account-available",
		},
	}
}

func (this *TeledbCheckDatabaseNameAvailableApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbCheckDatabaseNameAvailableRequest, header *TeledbCheckDatabaseNameAvailableRequestHeader) (CheckDatabaseNameAvailableResp *TeledbCheckDatabaseNameAvailableResponse, err error) {
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

	builder.AddParam("outerProdInstId", req.OuterProdInstId)
	builder.AddParam("dbName", req.DBName)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	CheckDatabaseNameAvailableResp = &TeledbCheckDatabaseNameAvailableResponse{}
	err = resp.Parse(CheckDatabaseNameAvailableResp)
	if err != nil {
		return
	}
	return CheckDatabaseNameAvailableResp, nil
}

type TeledbCheckDatabaseNameAvailableRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` // 实例ID，必填
	DBName          string `json:"dbName"`          // 账户名称
}
type TeledbCheckDatabaseNameAvailableRequestHeader struct {
	ProjectID string `json:"projectID"`
	InstID    string `json:"instId"`    // 实例ID，必填
	RegionID  string `json:"region_id"` // 资源池ID，必填
}
type TeledbCheckDatabaseNameAvailableResponseReturnObj struct {
	Available bool `json:"available"`
}

type TeledbCheckDatabaseNameAvailableResponse struct {
	StatusCode int32                                              `json:"statusCode"`      // 接口状态码
	Error      *string                                            `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                             `json:"message"`         // 描述信息
	ReturnObj  *TeledbCheckDatabaseNameAvailableResponseReturnObj `json:"returnObj"`
}
