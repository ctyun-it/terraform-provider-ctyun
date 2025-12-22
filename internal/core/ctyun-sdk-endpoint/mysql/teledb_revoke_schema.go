package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	http "net/http"
)

type TeledbRevokeSchemaApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbRevokeSchemaApi(client *ctyunsdk.CtyunClient) *TeledbRevokeSchemaApi {
	return &TeledbRevokeSchemaApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v1/open-api/account/privilege-revoke",
		},
	}
}

func (this *TeledbRevokeSchemaApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbRevokeSchemaRequest, header *TeledbRevokeSchemaRequestHeader) (RevokeSchemaResp *TeledbRevokeSchemaResponse, err error) {
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
	RevokeSchemaResp = &TeledbRevokeSchemaResponse{}
	err = resp.Parse(RevokeSchemaResp)
	if err != nil {
		return
	}
	return RevokeSchemaResp, nil
}

type DatabaseVO struct {
	RevokeSchema string `json:"revokeSchema"`
}

type TeledbRevokeSchemaRequest struct {
	OuterProdInstId string       `json:"outerProdInstId"` // 实例ID，必填
	AccountName     string       `json:"accountName"`     // 账户名称
	DatabaseVOList  []DatabaseVO `json:"databaseVOList"`
}

type TeledbRevokeSchemaRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}

type TeledbRevokeSchemaResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
