package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamGetIdentityProviderInfoApi
/* 查询身份供应商和关联用户或关联委托用户的信息 */
type CtiamGetIdentityProviderInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamGetIdentityProviderInfoApi(client *core.CtyunClient) *CtiamGetIdentityProviderInfoApi {
	return &CtiamGetIdentityProviderInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/identityProvider/getIdentityProviderInfo",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamGetIdentityProviderInfoApi) Do(ctx context.Context, credential core.Credential, req *CtiamGetIdentityProviderInfoRequest) (*CtiamGetIdentityProviderInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamGetIdentityProviderInfoRequest
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
	var resp CtiamGetIdentityProviderInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamGetIdentityProviderInfoRequest struct {
	IdPId      string  `json:"idPId"`                /*  身份供应商ID  */
	EntityId   string  `json:"entityId"`             /*  实体ID，对应IdP元数据文件中EntityId  */
	NameId     string  `json:"nameId"`               /*  身份供应商传递的需要代入的用户ID，用于根据转换规则确定具体用户  */
	LoginEmail *string `json:"loginEmail,omitempty"` /*  登录邮箱,nameId和登录邮箱二选一必填  */
}

type CtiamGetIdentityProviderInfoResponse struct {
	StatusCode *string                                        `json:"statusCode"` /*  兼容性返回码，800标识成功，CTIAM_XXX表示失败  */
	ReturnObj  *CtiamGetIdentityProviderInfoReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    *string                                        `json:"message"`    /*  返回信息，请求失败时会回传错误信息。  */
	Error      *string                                        `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamGetIdentityProviderInfoReturnObjResponse struct {
	IdPId          *string                                                  `json:"idPId"`          /*  身份供应商ID  */
	RawType        *string                                                  `json:"type"`           /*  身份供应商类型， 1 IAM用户， 0 虚用户  */
	EntityId       *string                                                  `json:"entityId"`       /*  身份供应商元数据文件中的实体ID  */
	Signature      *string                                                  `json:"signature"`      /*  身份供应商元数据文件中的签名  */
	UserList       []*CtiamGetIdentityProviderInfoReturnObjUserListResponse `json:"userList"`       /*  用户列表  */
	LoginLocation  *string                                                  `json:"loginLocation"`  /*  身份供应商元数据文件中的登录地址  */
	LogoutLocation *string                                                  `json:"logoutLocation"` /*  身份供应商元数据文件中的登出地址  */
}

type CtiamGetIdentityProviderInfoReturnObjUserListResponse struct {
	AccountId  *string `json:"accountId"`  /*  天翼云账户ID  */
	UserId     *string `json:"userId"`     /*  天翼云用户ID  */
	UserName   *string `json:"userName"`   /*  天翼云用户名  */
	RootUserId *string `json:"rootUserId"` /*  根用户ID  */
	Mobile     *string `json:"mobile"`     /*  手机号  */
	Email      *string `json:"email"`      /*  邮箱  */
}
