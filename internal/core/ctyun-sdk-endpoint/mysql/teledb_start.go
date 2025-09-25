package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbStartApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbStartApi(client *ctyunsdk.CtyunClient) *TeledbStartApi {
	return &TeledbStartApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v1/open-api/instance/start-instance",
		},
	}
}

func (this *TeledbStartApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbStartRequest, header *TeledbStartRequestHeader) (teledbStartResp *TeledbStartResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
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
	teledbStartResp = &TeledbStartResponse{}
	err = resp.Parse(teledbStartResp)
	if err != nil {
		return
	}
	return teledbStartResp, nil
}

type TeledbStartRequest struct {
	OuterProdInstId string `json:"outerProdInstId"`
}

type TeledbStartResponse struct {
	StatusCode int    `json:"statusCode"` // 返回码
	Message    string `json:"message"`    // 结果信息
	Error      string `json:"error"`      // 错误码，接口失败时返回
	//ReturnObj  interface{} `json:"returnObj"`  // 返回对象
}

type TeledbStartRequestHeader struct {
	ProjectID string `json:"projectId"`
	InstID    string `json:"instId"`
	RegionID  string `json:"regionId"`
}
