package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamUpdateUserApi
/* 修改用户 */
type CtiamUpdateUserApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamUpdateUserApi(client *core.CtyunClient) *CtiamUpdateUserApi {
	return &CtiamUpdateUserApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/updateUser",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamUpdateUserApi) Do(ctx context.Context, credential core.Credential, req *CtiamUpdateUserRequest) (*CtiamUpdateUserResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamUpdateUserRequest
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
	var resp CtiamUpdateUserResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamUpdateUserRequest struct {
	UserId       string  `json:"userId"`                 /*  用户id  */
	Remark       *string `json:"remark,omitempty"`       /*  描述信息  */
	LoginEmail   *string `json:"loginEmail,omitempty"`   /*  邮箱，和用户登录名自定义前缀必填其一  */
	MobilePhone  string  `json:"mobilePhone"`            /*  电话号码  */
	UserName     string  `json:"userName"`               /*  用户名  */
	UserNickName *string `json:"userNickName,omitempty"` /*  用户登录名自定义前缀，邮箱和用户登录名自定义前缀必填其一  */
	Prohibit     int32   `json:"prohibit"`               /*  是否禁用，0启用，1禁用  */
}

type CtiamUpdateUserResponse struct {
	StatusCode *string                           `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                           `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	ReturnObj  *CtiamUpdateUserReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Error      *string                           `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamUpdateUserReturnObjResponse struct {
	LoginEmail   *string `json:"loginEmail"`   /*  邮箱  */
	AccountId    *string `json:"accountId"`    /*  账号id  */
	MobilePhone  *string `json:"mobilePhone"`  /*  电话号码  */
	Remark       *string `json:"remark"`       /*  描述  */
	UserName     *string `json:"userName"`     /*  用户名  */
	LoginName    *string `json:"loginName"`    /*  登录名  */
	VirtualEmail *bool   `json:"virtualEmail"` /*  虚拟邮箱，true为是  */
}
