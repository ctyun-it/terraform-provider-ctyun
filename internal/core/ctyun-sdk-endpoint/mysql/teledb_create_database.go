package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbCreateDatabaseApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbCreateDatabaseApi(client *ctyunsdk.CtyunClient) *TeledbCreateDatabaseApi {
	return &TeledbCreateDatabaseApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v1/open-api/database",
		},
	}
}

func (this *TeledbCreateDatabaseApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbCreateDatabaseRequest, header *TeledbCreateDatabaseRequestHeader) (CreateDatabaseResp *TeledbCreateDatabaseResponse, err error) {
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
	if req.DBName == "" {
		err = fmt.Errorf("database_name is required")
		return
	}

	builder.AddHeader("inst-id", header.InstID)
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	CreateDatabaseResp = &TeledbCreateDatabaseResponse{}
	err = resp.Parse(CreateDatabaseResp)
	if err != nil {
		return
	}
	return CreateDatabaseResp, nil
}

type TeledbCreateDatabaseRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` // 实例ID，必填
	DBName          string `json:"dbName"`          // 数据库名称
	CharSetName     string `json:"charSetName"`     // 字符集
}
type TeledbCreateDatabaseRequestHeader struct {
	ProjectID string `json:"projectID"`
	InstID    string `json:"instId"`    // 实例ID，必填
	RegionID  string `json:"region_id"` // 资源池ID，必填
}

type TeledbCreateDatabaseResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
	//ReturnObj  *TeledbCreateDatabaseResponseReturnObj `json:"returnObj"`
}
