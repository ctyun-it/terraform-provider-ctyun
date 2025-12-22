package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbStartAuditApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbStartAuditApi(client *ctyunsdk.CtyunClient) *TeledbStartAuditApi {
	return &TeledbStartAuditApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v2/open-api/audit-log/switch",
		},
	}
}

func (this *TeledbStartAuditApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbStartAuditRequest, header *TeledbStartAuditRequestHeader) (TeledbStartAuditResp *TeledbStartAuditResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if header.InstID != "" {
		builder.AddHeader("inst-id", header.InstID)
	}
	if header.RegionID == "" {
		err = errors.New("missing required field: RegionID")
		return
	}
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	TeledbStartAuditResp = &TeledbStartAuditResponse{}
	err = resp.Parse(TeledbStartAuditResp)
	if err != nil {
		return
	}
	return TeledbStartAuditResp, nil
}

type TeledbStartAuditRequest struct {
	OuterProdInstId string `json:"outerProdInstId"`
	AuditSwitch     bool   `json:"auditSwitch"`
}

type TeledbStartAuditResponse struct {
	StatusCode int    `json:"statusCode"` // 返回码
	Message    string `json:"message"`    // 结果信息
	Error      string `json:"error"`      // 错误码，接口失败时返回
	//ReturnObj  interface{} `json:"returnObj"`  // 返回对象
}

type TeledbStartAuditRequestHeader struct {
	ProjectID *string `json:"project-id"`
	InstID    string  `json:"inst-id"`
	RegionID  string  `json:"regionId"`
}
