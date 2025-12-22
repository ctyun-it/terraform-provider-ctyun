package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryUsersApi
type CtiamQueryUsersApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryUsersApi(client *core.CtyunClient) *CtiamQueryUsersApi {
	return &CtiamQueryUsersApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/openapi/user/getUsers",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryUsersApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryUsersRequest) (*CtiamQueryUsersResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamQueryUsersRequest
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
	var resp CtiamQueryUsersResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryUsersRequest struct {
	PageNum  int32 `json:"pageNum"`  /*  页数  */
	PageSize int32 `json:"pageSize"` /*  每页条数  */
}

type CtiamQueryUsersResponse struct {
	StatusCode *string                           `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamQueryUsersReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                           `json:"message"`    /*  返回信息，请求失败时会回传错误信息。  */
	Error      *string                           `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamQueryUsersReturnObjResponse struct {
	PageNum  *int32                                    `json:"pageNum"`  /*  当前页数  */
	PageSize *int32                                    `json:"pageSize"` /*  每页条数  */
	Pages    *int32                                    `json:"pages"`    /*  总页数  */
	Total    *int64                                    `json:"total"`    /*  总条数  */
	Result   []*CtiamQueryUsersReturnObjResultResponse `json:"result"`   /*  用户信息  */
}

type CtiamQueryUsersReturnObjResultResponse struct {
	UserId       *string                                         `json:"userId"`       /*  用户id  */
	LoginEmail   *string                                         `json:"loginEmail"`   /*  邮箱  */
	AccountId    *string                                         `json:"accountId"`    /*  账号id  */
	MobilePhone  *string                                         `json:"mobilePhone"`  /*  电话号码  */
	Groups       []*CtiamQueryUsersReturnObjResultGroupsResponse `json:"groups"`       /*  用户组  */
	Remark       *string                                         `json:"remark"`       /*  描述信息  */
	UserName     *string                                         `json:"userName"`     /*  用户名  */
	IsRoot       int32                                           `json:"isRoot"`       /*  是否是主用户（1：主用户，0：子用户）  */
	CreateDate   int64                                           `json:"createDate"`   /*  创建时间  */
	Prohibit     int32                                           `json:"prohibit"`     /*  禁用账号，只针对子账号才能是禁用的状态 是否启用( 0启用 , 1 禁用)  */
	LoginName    *string                                         `json:"loginName"`    /*  登录名  */
	UserNickName *string                                         `json:"userNickName"` /*  用户自定义登录名前缀  */
	VirtualEmail *string                                         `json:"virtualEmail"` /*  是否虚拟邮箱，true为虚拟邮箱  */
}

type CtiamQueryUsersReturnObjResultGroupsResponse struct {
	Id         *string `json:"id"`         /*  用户组id  */
	GroupName  *string `json:"groupName"`  /*  用户组名称  */
	GroupIntro *string `json:"groupIntro"` /*  用户组信息  */
}
