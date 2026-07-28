package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamInvalidGroupApi
/* 注销用户组 */
type CtiamInvalidGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamInvalidGroupApi(client *core.CtyunClient) *CtiamInvalidGroupApi {
	return &CtiamInvalidGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/userGroup/invalidGroup",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamInvalidGroupApi) Do(ctx context.Context, credential core.Credential, req *CtiamInvalidGroupRequest) (*CtiamInvalidGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamInvalidGroupRequest
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
	var resp CtiamInvalidGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamInvalidGroupRequest struct {
	GroupId string `json:"groupId"` /*  用户组id  */
}

type CtiamInvalidGroupResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息。  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
