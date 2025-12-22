package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbUpdateDatabaseRemarkApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbUpdateDatabaseRemarkApi(client *ctyunsdk.CtyunClient) *TeledbUpdateDatabaseRemarkApi {
	return &TeledbUpdateDatabaseRemarkApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v1/open-api/database/edit-remark",
		},
	}
}

func (this *TeledbUpdateDatabaseRemarkApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbUpdateDatabaseRemarkRequest, header *TeledbUpdateDatabaseRemarkRequestHeader) (UpdateDatabaseRemarkResp *TeledbUpdateDatabaseRemarkResponse, err error) {
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
	builder.AddHeader("inst-id", header.InstID)
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	UpdateDatabaseRemarkResp = &TeledbUpdateDatabaseRemarkResponse{}
	err = resp.Parse(UpdateDatabaseRemarkResp)
	if err != nil {
		return
	}
	return UpdateDatabaseRemarkResp, nil
}

type TeledbUpdateDatabaseRemarkRequest struct {
	OuterProdInstId string  `json:"outerProdInstId"` // 实例ID，必填
	DatabaseName    string  `json:"databaseName"`    // 账户名称
	Remark          *string `json:"remark"`
}
type TeledbUpdateDatabaseRemarkRequestHeader struct {
	ProjectID string `json:"Project-Id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}
type TeledbUpdateDatabaseRemarkResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
