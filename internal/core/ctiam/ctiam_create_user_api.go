package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamCreateUserApi
/* 创建用户 */
type CtiamCreateUserApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamCreateUserApi(client *core.CtyunClient) *CtiamCreateUserApi {
	return &CtiamCreateUserApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/createUser",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamCreateUserApi) Do(ctx context.Context, credential core.Credential, req *CtiamCreateUserRequest) (*CtiamCreateUserResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamCreateUserRequest
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
	var resp CtiamCreateUserResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamCreateUserRequest struct {
	LoginEmail         *string                              `json:"loginEmail,omitempty"`         /*  邮箱和子用户的登录名自定义前缀，二者必填其一  */
	Password           *string                              `json:"password,omitempty"`           /*  generatePassword字段为false时，password和sourcePassword，必填其中之一，password为加密后的密码，sourcePassword为源密码（无需加密）  */
	MobilePhone        string                               `json:"mobilePhone"`                  /*  电话号码  */
	UserName           string                               `json:"userName"`                     /*  用户名  */
	Groups             []*CtiamCreateUserGroupsRequest      `json:"groups,omitempty"`             /*  用户组信息  */
	Remark             *string                              `json:"remark,omitempty"`             /*  描述  */
	SourcePassword     *string                              `json:"sourcePassword,omitempty"`     /*  generatePassword字段为false时，password和sourcePassword，必填其中之一，password为加密后的密码，sourcePassword为源密码（无需加密）  */
	GeneratePassword   *bool                                `json:"generatePassword,omitempty"`   /*  创建用户是否自动生成密码,true：是, false：否（默认false，自定义密码）  */
	LoginResetPassword *bool                                `json:"loginResetPassword,omitempty"` /*  用户在下次登录时是否重置密码（true:用户在下次登录时必须重置密码, false:无需重置（默认false，无需重置密码））  */
	AccessControl      *CtiamCreateUserAccessControlRequest `json:"accessControl,omitempty"`      /*  控制台访问和编程式访问设置（默认为空，开启）  */
	UserNickName       *string                              `json:"userNickName,omitempty"`       /*  创建子用户的登录名自定义前缀，长度为1-7位，首字母必须是英文字母或汉字，只支持中文、英文大小写字母、数字和下划线，不能以下划线结尾。和邮箱二者必填其一  */
}

type CtiamCreateUserGroupsRequest struct {
	Id *string `json:"id,omitempty"` /*  用户组id  */
}

type CtiamCreateUserAccessControlRequest struct {
	ConsoleAccess *string `json:"consoleAccess,omitempty"` /*  控制台访问(PUB_100_02_0001: 禁止;PUB_100_02_0002: 允许)  */
}

type CtiamCreateUserResponse struct {
	StatusCode *string                           `json:"statusCode"` /*  状态码  */
	ReturnObj  *CtiamCreateUserReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                           `json:"message"`    /*  对于调用失败的补充信息说明  */
	Error      *string                           `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamCreateUserReturnObjResponse struct {
	LoginEmail   *string                                   `json:"loginEmail"`   /*  邮箱  */
	AccountId    *string                                   `json:"accountId"`    /*  账号id  */
	MobilePhone  *string                                   `json:"mobilePhone"`  /*  电话号码  */
	Groups       []*CtiamCreateUserReturnObjGroupsResponse `json:"groups"`       /*  用户组信息  */
	Remark       *string                                   `json:"remark"`       /*  描述  */
	UserName     *string                                   `json:"userName"`     /*  用户名  */
	UserId       *string                                   `json:"userId"`       /*  用户id  */
	Password     *string                                   `json:"password"`     /*  generatePassword为true是返回  */
	LoginName    *string                                   `json:"loginName"`    /*  登录名  */
	VirtualEmail *bool                                     `json:"virtualEmail"` /*  虚拟邮箱  */
}

type CtiamCreateUserReturnObjGroupsResponse struct {
	Id *string `json:"id"` /*  用户组id  */
}
