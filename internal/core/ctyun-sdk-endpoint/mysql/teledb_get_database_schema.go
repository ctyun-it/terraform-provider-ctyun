package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetDatabaseSchemaApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetDatabaseSchemaApi(client *ctyunsdk.CtyunClient) *TeledbGetDatabaseSchemaApi {
	return &TeledbGetDatabaseSchemaApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v1/open-api/database/schemas",
		},
	}
}

func (this *TeledbGetDatabaseSchemaApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetDatabaseSchemaRequest, header *TeledbGetDatabaseSchemaRequestHeader) (GetDatabaseSchemaResp *TeledbGetDatabaseSchemaResponse, err error) {
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

	if req.OuterProdInstId == "" {
		err = errors.New("instId 为空")
		return
	}
	builder.AddHeader("inst-id", header.InstID)
	builder.AddParam("outerProdInstId", req.OuterProdInstId)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	GetDatabaseSchemaResp = &TeledbGetDatabaseSchemaResponse{}
	err = resp.Parse(GetDatabaseSchemaResp)
	if err != nil {
		return
	}
	return GetDatabaseSchemaResp, nil
}

type TeledbGetDatabaseSchemaRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` // 外部实例ID，必填
}

type TeledbGetDatabaseSchemaRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}

type TeledbGetDatabaseSchemaResponse struct {
	StatusCode int32                                      `json:"statusCode"`      // 接口状态码
	Error      *string                                    `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                     `json:"message"`         // 描述信息
	ReturnObj  []TeledbGetDatabaseSchemaResponseReturnObj `json:"returnObj"`
}
type TeledbGetDatabaseSchemaResponseReturnObjUserVo struct {
	AccountName string `json:"accountName"`
	ReadOnly    bool   `json:"readOnly"`
	SelectPriv  string `json:"selectPriv"`
	InsertPriv  string `json:"insertPriv"`
}

type TeledbGetDatabaseSchemaResponseReturnObj struct {
	UserVOList  []TeledbGetDatabaseSchemaResponseReturnObjUserVo `json:"userVOList"`
	GrantSchema string                                           `json:"grantSchema"`
}
