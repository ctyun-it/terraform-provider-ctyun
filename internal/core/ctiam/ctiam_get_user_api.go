package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamGetUserApi
/* 根据id查询用户详情 */
type CtiamGetUserApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamGetUserApi(client *core.CtyunClient) *CtiamGetUserApi {
	return &CtiamGetUserApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/user/getUser",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamGetUserApi) Do(ctx context.Context, credential core.Credential, req *CtiamGetUserRequest) (*CtiamGetUserResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("userId", req.UserId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamGetUserResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamGetUserRequest struct {
	UserId string `json:"userId"` /*  用户id  */
}

type CtiamGetUserResponse struct {
	StatusCode *string                        `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamGetUserReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                        `json:"message"`    /*  返回信息，请求失败时会回传错误信息。  */
	Error      *string                        `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamGetUserReturnObjResponse struct {
	LoginEmail   *string                                `json:"loginEmail"`   /*  邮箱  */
	AccountId    *string                                `json:"accountId"`    /*  账号id  */
	MobilePhone  *string                                `json:"mobilePhone"`  /*  电话号码  */
	Groups       []*CtiamGetUserReturnObjGroupsResponse `json:"groups"`       /*  用户组列表  */
	Remark       *string                                `json:"remark"`       /*  描述  */
	UserName     *string                                `json:"userName"`     /*  用户名  */
	VirtualEmail *string                                `json:"virtualEmail"` /*  虚拟邮箱标识  */
	LoginName    *string                                `json:"loginName"`    /*  登录名  */
	UserNickName *string                                `json:"userNickName"` /*  用户登录名自定义前缀  */
}

type CtiamGetUserReturnObjGroupsResponse struct {
	Id *string `json:"id"` /*  用户组id  */
}
