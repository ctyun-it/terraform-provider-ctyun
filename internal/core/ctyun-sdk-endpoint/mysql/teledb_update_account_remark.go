package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbUpdateAccountRemarkApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbUpdateAccountRemarkApi(client *ctyunsdk.CtyunClient) *TeledbUpdateAccountRemarkApi {
	return &TeledbUpdateAccountRemarkApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/RDS2/v1/open-api/account/account-remark",
		},
	}
}

func (this *TeledbUpdateAccountRemarkApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbUpdateAccountRemarkRequest, header *TeledbUpdateAccountRemarkRequestHeader) (UpdateAccountRemarkResp *TeledbUpdateAccountRemarkResponse, err error) {
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
	UpdateAccountRemarkResp = &TeledbUpdateAccountRemarkResponse{}
	err = resp.Parse(UpdateAccountRemarkResp)
	if err != nil {
		return
	}
	return UpdateAccountRemarkResp, nil
}

type TeledbUpdateAccountRemarkRequest struct {
	OuterProdInstId string  `json:"outerProdInstId"` // 实例ID，必填
	AccountName     string  `json:"accountName"`     // 账户名称
	Remark          *string `json:"remark"`
}
type TeledbUpdateAccountRemarkRequestHeader struct {
	ProjectID string `json:"projectID"`
	InstID    string `json:"instId"`    // 实例ID，必填
	RegionID  string `json:"region_id"` // 资源池ID，必填
}
type TeledbUpdateAccountRemarkResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
