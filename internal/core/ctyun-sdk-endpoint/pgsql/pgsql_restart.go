package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlRestartApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlRestartApi(client *ctyunsdk.CtyunClient) *PgsqlRestartApi {
	return &PgsqlRestartApi{client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/PG/v1/product/restart",
		}}
}

func (this *PgsqlRestartApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlRestartRequest, header *PgsqlRestartRequestHeader) (PgsqlRestartResp *PgsqlRestartResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}

	if header.RegionID == "" {
		err = errors.New("missing required field: RegionID")
		return
	}
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	PgsqlRestartResp = &PgsqlRestartResponse{}
	err = resp.Parse(PgsqlRestartResp)
	if err != nil {
		return
	}
	return PgsqlRestartResp, nil
}

type PgsqlRestartRequest struct {
	ProdInstId string `json:"prodInstId"`
}
type PgsqlRestartRequestHeader struct {
	ProjectID string `json:"projectId"`
	RegionID  string `json:"regionId"`
}
type PgsqlRestartResponse struct {
	StatusCode int    `json:"statusCode"` // 返回码
	Message    string `json:"message"`    // 结果信息
	Error      string `json:"error"`      // 错误码，接口失败时返回
	//ReturnObj  interface{} `json:"returnObj"`  // 返回对象
}
