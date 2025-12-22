package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetAccountInfoApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetAccountInfoApi(client *ctyunsdk.CtyunClient) *TeledbGetAccountInfoApi {
	return &TeledbGetAccountInfoApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v1/open-api/account",
		},
	}
}

func (this *TeledbGetAccountInfoApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetAccountInfoRequest, header *TeledbGetAccountInfoRequestHeader) (GetAccountInfoResp *TeledbGetAccountInfoResponse, err error) {
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
	GetAccountInfoResp = &TeledbGetAccountInfoResponse{}
	err = resp.Parse(GetAccountInfoResp)
	if err != nil {
		return
	}
	return GetAccountInfoResp, nil
}

type TeledbGetAccountInfoRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` // 实例ID，必填
}
type TeledbGetAccountInfoRequestHeader struct {
	ProjectID string `json:"projectID"`
	InstID    string `json:"instId"`    // 实例ID，必填
	RegionID  string `json:"region_id"` // 资源池ID，必填
}
type TeledbGetAccountInfoResponse struct {
	StatusCode int32                                   `json:"statusCode"`      // 接口状态码
	Error      *string                                 `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                  `json:"message"`         // 描述信息
	ReturnObj  []TeledbGetAccountInfoResponseReturnObj `json:"returnObj"`
}

type TeledbGetAccountInfoResponseReturnObj struct {
	SchemaPrivilegeVOList []SchemaPrivilegeVO `json:"schemaPrivilegeVOList"`
	AccountName           string              `json:"accountName"`
}

//GrantSchema  string `json:"grantSchema"`
//ReadOnly     bool   `json:"readOnly,omitempty"`
//DMLPrivilege bool   `json:"dmlPrivilege,omitempty"`
//DDLPrivilege bool   `json:"ddlPrivilege,omitempty"`
//ReadAndWrite bool   `json:"readAndWrite,omitempty"`
