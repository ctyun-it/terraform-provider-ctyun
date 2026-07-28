package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryDelegateRolesApi
/* 查询委托角色分页信息 */
type CtiamQueryDelegateRolesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryDelegateRolesApi(client *core.CtyunClient) *CtiamQueryDelegateRolesApi {
	return &CtiamQueryDelegateRolesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/delegate/queryDelegateRoles",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryDelegateRolesApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryDelegateRolesRequest) (*CtiamQueryDelegateRolesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamQueryDelegateRolesRequest
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
	var resp CtiamQueryDelegateRolesResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryDelegateRolesRequest struct {
	PageNum  int32 `json:"pageNum"`  /*  页码  */
	PageSize int32 `json:"pageSize"` /*  每页条数  */
}

type CtiamQueryDelegateRolesResponse struct {
	StatusCode *string                                   `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                   `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	ReturnObj  *CtiamQueryDelegateRolesReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Error      *string                                   `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamQueryDelegateRolesReturnObjResponse struct {
	List     []*CtiamQueryDelegateRolesReturnObjListResponse `json:"list"`     /*  委托列表  */
	PageNum  *int32                                          `json:"pageNum"`  /*  页码  */
	PageSize *int32                                          `json:"pageSize"` /*  每页条数  */
	Total    *int32                                          `json:"total"`    /*  总条数  */
	Pages    *int32                                          `json:"pages"`    /*  总页数  */
}

type CtiamQueryDelegateRolesReturnObjListResponse struct {
	Id              *int32  `json:"id"`              /*  id  */
	Name            *string `json:"name"`            /*  名称  */
	AccountId       *string `json:"accountId"`       /*  账号id  */
	AssumeAccountId *string `json:"assumeAccountId"` /*  委托账号id  */
	AssumeUserId    *string `json:"assumeUserId"`    /*  委托用户id  */
	RawType         *int32  `json:"type"`            /*  类型  */
	Remark          *string `json:"remark"`          /*  描述  */
	Status          *int32  `json:"status"`          /*  状态  */
	CreateTime      *string `json:"createTime"`      /*  创建时间  */
	UpdateTime      *string `json:"updateTime"`      /*  更新时间  */
}
