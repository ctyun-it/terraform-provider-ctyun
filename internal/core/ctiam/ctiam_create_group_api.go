package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamCreateGroupApi
/* 创建用户组 */
type CtiamCreateGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamCreateGroupApi(client *core.CtyunClient) *CtiamCreateGroupApi {
	return &CtiamCreateGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/userGroup/createGroup",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamCreateGroupApi) Do(ctx context.Context, credential core.Credential, req *CtiamCreateGroupRequest) (*CtiamCreateGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamCreateGroupRequest
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
	var resp CtiamCreateGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamCreateGroupRequest struct {
	GroupName  string `json:"groupName"`  /*  用户组名称  */
	GroupIntro string `json:"groupIntro"` /*  用户组描述  */
}

type CtiamCreateGroupResponse struct {
	StatusCode *string                            `json:"statusCode"` /*  状态码  */
	ReturnObj  *CtiamCreateGroupReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                            `json:"message"`    /*  对于调用失败的补充信息说明  */
	Error      *string                            `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamCreateGroupReturnObjResponse struct {
	Id         *string `json:"id"`         /*  用户组id  */
	GroupName  *string `json:"groupName"`  /*  用户组名称  */
	AccountId  *string `json:"accountId"`  /*  账号id  */
	GroupIntro *string `json:"groupIntro"` /*  用户组描述  */
}
