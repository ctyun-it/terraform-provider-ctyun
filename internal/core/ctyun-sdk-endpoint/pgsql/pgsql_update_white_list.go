package pgsql

import (
	"context"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlUpdateWhiteListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlUpdateWhiteListApi(client *ctyunsdk.CtyunClient) *PgsqlUpdateWhiteListApi {
	return &PgsqlUpdateWhiteListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/PG/v1/white-list/modify",
		},
	}
}

func (this *PgsqlUpdateWhiteListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlUpdateWhiteListRequest, header *PgsqlUpdateWhiteListRequestHeader) (CreateWhiteListResp *PgsqlUpdateWhiteListResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}

	if header.RegionID == "" {
		err = fmt.Errorf("region_id is required")
		return
	}
	if req.ProdInstId == "" {
		err = fmt.Errorf("ProdInstId is required")
		return
	}
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	CreateWhiteListResp = &PgsqlUpdateWhiteListResponse{}
	err = resp.Parse(CreateWhiteListResp)
	if err != nil {
		return
	}
	return CreateWhiteListResp, nil
}

type PgsqlUpdateWhiteListRequest struct {
	ProdInstId string   `json:"prodInstId"` // 实例ID，必填
	Mode       string   `json:"mode"`       // 账户名称
	IpList     []string `json:"ipList"`     // ip白名单列表
}
type PgsqlUpdateWhiteListRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlUpdateWhiteListResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
