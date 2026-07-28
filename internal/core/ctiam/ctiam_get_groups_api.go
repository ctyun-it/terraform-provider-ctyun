package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamGetGroupsApi
/* 分页查询用户组 */
type CtiamGetGroupsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamGetGroupsApi(client *core.CtyunClient) *CtiamGetGroupsApi {
	return &CtiamGetGroupsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/userGroup/getGroups",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamGetGroupsApi) Do(ctx context.Context, credential core.Credential, req *CtiamGetGroupsRequest) (*CtiamGetGroupsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamGetGroupsRequest
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
	var resp CtiamGetGroupsResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamGetGroupsRequest struct {
	PageNum   int32   `json:"pageNum"`             /*  页码  */
	PageSize  int32   `json:"pageSize"`            /*  每页条数  */
	GroupName *string `json:"groupName,omitempty"` /*  用户组名称  */
}

type CtiamGetGroupsResponse struct {
	StatusCode *string                          `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamGetGroupsReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                          `json:"message"`    /*  返回信息，请求失败时会回传错误信息。  */
	Error      *string                          `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamGetGroupsReturnObjResponse struct {
	PageNum  *int32                                   `json:"pageNum"`  /*  页码  */
	PageSize *int32                                   `json:"pageSize"` /*  每页条数  */
	Total    *int32                                   `json:"total"`    /*  总条数  */
	Pages    *int32                                   `json:"pages"`    /*  总页数  */
	Result   []*CtiamGetGroupsReturnObjResultResponse `json:"result"`   /*  用户组列表  */
}

type CtiamGetGroupsReturnObjResultResponse struct {
	Id         *string `json:"id"`         /*  用户组id  */
	GroupName  *string `json:"groupName"`  /*  用户组名称  */
	AccountId  *string `json:"accountId"`  /*  账号id  */
	GroupIntro *string `json:"groupIntro"` /*  用户组信息  */
	IsRoot     *string `json:"isRoot"`     /*  是否ROOT组  */
	IsValid    *string `json:"isValid"`    /*  是否有效  */
	HwGroupId  *string `json:"hwGroupId"`  /*  华为用户组id  */
	UserCount  *string `json:"userCount"`  /*  用户数量  */
	CreateTime *string `json:"createTime"` /*  创建时间  */
	UpdateTime *string `json:"updateTime"` /*  更新时间  */
}
