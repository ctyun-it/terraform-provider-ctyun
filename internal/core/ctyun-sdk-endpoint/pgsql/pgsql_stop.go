package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlStopApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlStopApi(client *ctyunsdk.CtyunClient) *PgsqlStopApi {
	return &PgsqlStopApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/PG/v1/product/stop",
		},
	}
}

func (this *PgsqlStopApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlStopRequest, header *PgsqlStopRequestHeader) (PgsqlStopResp *PgsqlStopResponse, err error) {
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
	PgsqlStopResp = &PgsqlStopResponse{}
	err = resp.Parse(PgsqlStopResp)
	if err != nil {
		return
	}
	return PgsqlStopResp, nil
}

type PgsqlStopRequest struct {
	ProdInstId string `json:"prodInstId"`
}
type PgsqlStopRequestHeader struct {
	ProjectID string `json:"projectId"`
	RegionID  string `json:"regionId"`
}
type PgsqlStopResponse struct {
	StatusCode int    `json:"statusCode"` // 返回码
	Message    string `json:"message"`    // 结果信息
	Error      string `json:"error"`      // 错误码，接口失败时返回
	//ReturnObj  interface{} `json:"returnObj"`  // 返回对象
}
