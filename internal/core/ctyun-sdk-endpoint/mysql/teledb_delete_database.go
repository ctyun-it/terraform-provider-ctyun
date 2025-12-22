package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbDeleteDatabaseApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbDeleteDatabaseApi(client *ctyunsdk.CtyunClient) *TeledbDeleteDatabaseApi {
	return &TeledbDeleteDatabaseApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodDelete,
			UrlPath: "/RDS2/v1/open-api/database",
		},
	}
}

func (this *TeledbDeleteDatabaseApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbDeleteDatabaseRequest, header *TeledbDeleteDatabaseRequestHeader) (DeleteDatabaseResp *TeledbDeleteDatabaseResponse, err error) {
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
	builder.AddParam("dbName", fmt.Sprintf("%s", req.DBName))

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	DeleteDatabaseResp = &TeledbDeleteDatabaseResponse{}
	err = resp.Parse(DeleteDatabaseResp)
	if err != nil {
		return
	}
	return DeleteDatabaseResp, nil
}

type TeledbDeleteDatabaseRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` // 实例ID，必填
	DBName          string `json:"dbName"`
}
type TeledbDeleteDatabaseRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}
type TeledbDeleteDatabaseResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
