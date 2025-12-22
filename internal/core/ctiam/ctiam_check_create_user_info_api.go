package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamCheckCreateUserInfoApi
type CtiamCheckCreateUserInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamCheckCreateUserInfoApi(client *core.CtyunClient) *CtiamCheckCreateUserInfoApi {
	return &CtiamCheckCreateUserInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/checkCreateUserInfo",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamCheckCreateUserInfoApi) Do(ctx context.Context, credential core.Credential, req *CtiamCheckCreateUserInfoRequest) (*CtiamCheckCreateUserInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamCheckCreateUserInfoRequest
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
	var resp CtiamCheckCreateUserInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamCheckCreateUserInfoRequest struct {
	UserList []*CtiamCheckCreateUserInfoUserListRequest `json:"userList,omitempty"` /*  用户列表  */
}

type CtiamCheckCreateUserInfoUserListRequest struct {
	Id          int32  `json:"id"`          /*  序列号  */
	LoginEmail  string `json:"loginEmail"`  /*  邮箱  */
	MobilePhone string `json:"mobilePhone"` /*  手机号  */
}

type CtiamCheckCreateUserInfoResponse struct {
	StatusCode *string                                    `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                    `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                    `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamCheckCreateUserInfoReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamCheckCreateUserInfoReturnObjResponse struct {
	CheckResultList []*CtiamCheckCreateUserInfoReturnObjCheckResultListResponse `json:"checkResultList"` /*  返回结果  */
}

type CtiamCheckCreateUserInfoReturnObjCheckResultListResponse struct {
	Id                  *int32 `json:"id"`                  /*  序列号  */
	CheckMobilePhoneRes *bool  `json:"checkMobilePhoneRes"` /*  手机号是否可使用：true：是，false：否  */
	CheckLoginEmailRes  *bool  `json:"checkLoginEmailRes"`  /*  邮箱是否可使用：true：是，false：否  */
}
