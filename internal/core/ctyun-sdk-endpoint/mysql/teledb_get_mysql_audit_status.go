package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetAuditStatusApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetAuditStatusApi(client *ctyunsdk.CtyunClient) *TeledbGetAuditStatusApi {
	return &TeledbGetAuditStatusApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v2/open-api/audit-log/switch",
		},
	}
}

func (this *TeledbGetAuditStatusApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetAuditStatusRequest, header *TeledbGetAuditStatusRequestHeader) (GetAuditStatusResp *TeledbGetAuditStatusResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
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
	GetAuditStatusResp = &TeledbGetAuditStatusResponse{}
	err = resp.Parse(GetAuditStatusResp)
	if err != nil {
		return
	}
	return GetAuditStatusResp, nil
}

type TeledbGetAuditStatusRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` // 外部实例ID，必填
}

type TeledbGetAuditStatusRequestHeader struct {
	ProjectID *string `json:"project-id"`
	InstID    string  `json:"instId"`   // 实例ID，必填
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}

type TeledbGetAuditStatusResponse struct {
	StatusCode int32                                  `json:"statusCode"`      // 接口状态码
	Error      *string                                `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                 `json:"message"`         // 描述信息
	ReturnObj  *TeledbGetAuditStatusResponseReturnObj `json:"returnObj"`
}
type TeledbGetAuditStatusResponseReturnObj struct {
	AuditLogSwitch bool `json:"auditLogSwitch"`
}
