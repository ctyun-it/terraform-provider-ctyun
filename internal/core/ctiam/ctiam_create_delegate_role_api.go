package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamCreateDelegateRoleApi
/* 创建委托 */
type CtiamCreateDelegateRoleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamCreateDelegateRoleApi(client *core.CtyunClient) *CtiamCreateDelegateRoleApi {
	return &CtiamCreateDelegateRoleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/delegate/createDelegateRole",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamCreateDelegateRoleApi) Do(ctx context.Context, credential core.Credential, req *CtiamCreateDelegateRoleRequest) (*CtiamCreateDelegateRoleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamCreateDelegateRoleRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamCreateDelegateRoleResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamCreateDelegateRoleRequest struct {
	Name               string  `json:"name"`                         /*  委托名称  */
	RawType            int32   `json:"type"`                         /*  委托类型（0账号级委托、1 云服务委托、2身份供应商）  */
	AssumeAccountId    *string `json:"assumeAccountId,omitempty"`    /*  委托账号（委托类型为0账号级委托、1 云服务委托时，必填）  */
	Remark             *string `json:"remark,omitempty"`             /*  描述  */
	IdentityProviderId *string `json:"identityProviderId,omitempty"` /*  身份提供商id，委托类型为2身份供应商时，必填  */
}

type CtiamCreateDelegateRoleResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	ReturnObj  *string `json:"returnObj"`  /*  返回参数  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
