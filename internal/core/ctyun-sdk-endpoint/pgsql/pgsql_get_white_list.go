package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlGetWhiteListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlGetWhiteListApi(client *ctyunsdk.CtyunClient) *PgsqlGetWhiteListApi {
	return &PgsqlGetWhiteListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/PG/v1/white-list/query",
		},
	}
}

func (this *PgsqlGetWhiteListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlGetWhiteListRequest, header *PgsqlGetWhiteListRequestHeader) (GetAccountInfoResp *PgsqlGetWhiteListResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if req.ProdInstId == "" {
		err = errors.New("instId 为空")
		return
	}
	if header.RegionID == "" {
		err = fmt.Errorf("region_id is required")
		return
	}
	builder.AddHeader("regionId", header.RegionID)
	builder.AddParam("prodInstId", req.ProdInstId)
	if req.IP != nil {
		builder.AddParam("ip", *req.IP)
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	GetAccountInfoResp = &PgsqlGetWhiteListResponse{}
	err = resp.Parse(GetAccountInfoResp)
	if err != nil {
		return
	}
	return GetAccountInfoResp, nil
}

type PgsqlGetWhiteListRequest struct {
	ProdInstId string  `json:"prodInstId"` // 实例ID，必填
	IP         *string `json:"ip"`
}
type PgsqlGetWhiteListRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlGetWhiteListResponse struct {
	StatusCode int32    `json:"statusCode"`      // 接口状态码
	Error      *string  `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string   `json:"message"`         // 描述信息
	ReturnObj  []string `json:"returnObj"`
}
