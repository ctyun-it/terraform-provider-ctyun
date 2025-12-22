package pgsql

import (
	"context"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlCreateAccountApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlCreateAccountApi(client *ctyunsdk.CtyunClient) *PgsqlCreateAccountApi {
	return &PgsqlCreateAccountApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/PG/v1/inst-user/create",
		},
	}
}

func (this *PgsqlCreateAccountApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlCreateAccountRequest, header *PgsqlCreateAccountRequestHeader) (CreateAccountResp *PgsqlCreateAccountResponse, err error) {
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
	if req.Username == "" {
		err = fmt.Errorf("Username is required")
		return
	}
	if req.Password == "" {
		err = fmt.Errorf("Password is required")
		return
	}

	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	CreateAccountResp = &PgsqlCreateAccountResponse{}
	err = resp.Parse(CreateAccountResp)
	if err != nil {
		return
	}
	return CreateAccountResp, nil
}

type PgsqlCreateAccountRequest struct {
	ProdInstId string `json:"prodInstId"` // 实例ID，必填
	Username   string `json:"username"`   // 账户名称
	Password   string `json:"password"`   // 账户密码（安全考虑需要用base64加密后传输）
	UserType   string `json:"userType"`
}
type PgsqlCreateAccountRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlCreateAccountResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
