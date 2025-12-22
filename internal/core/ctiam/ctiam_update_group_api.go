package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamUpdateGroupApi
/* 修改用户组 */
type CtiamUpdateGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamUpdateGroupApi(client *core.CtyunClient) *CtiamUpdateGroupApi {
	return &CtiamUpdateGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/userGroup/updateGroup",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamUpdateGroupApi) Do(ctx context.Context, credential core.Credential, req *CtiamUpdateGroupRequest) (*CtiamUpdateGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamUpdateGroupRequest
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
	var resp CtiamUpdateGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamUpdateGroupRequest struct {
	Id         string `json:"id"`         /*  用户组id  */
	GroupName  string `json:"groupName"`  /*  用户组名称  */
	GroupIntro string `json:"groupIntro"` /*  用户组描述  */
}

type CtiamUpdateGroupResponse struct {
	StatusCode *string                            `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamUpdateGroupReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                            `json:"message"`    /*  返回信息，请求失败时会回传错误信息。  */
	Error      *string                            `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamUpdateGroupReturnObjResponse struct {
	Id         *string `json:"id"`         /*  用户组id  */
	GroupName  *string `json:"groupName"`  /*  用户组名称  */
	AccountId  *string `json:"accountId"`  /*  账号id  */
	GroupIntro *string `json:"groupIntro"` /*  用户组描述  */
}
