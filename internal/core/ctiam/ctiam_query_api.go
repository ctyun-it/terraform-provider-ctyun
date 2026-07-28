package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryApi
/* 通过账户ID分页查询权限 */
type CtiamQueryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryApi(client *core.CtyunClient) *CtiamQueryApi {
	return &CtiamQueryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/perm/query",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryRequest) (*CtiamQueryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamQueryRequest
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
	var resp CtiamQueryResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryRequest struct {
	PageNum  int32 `json:"pageNum"`  /*  页面  */
	PageSize int32 `json:"pageSize"` /*  每页条数  */
}

type CtiamQueryResponse struct {
	StatusCode *string                      `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamQueryReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                      `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                      `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamQueryReturnObjResponse struct {
	PageNum  *int32                             `json:"pageNum"`  /*  页码  */
	PageSize *int32                             `json:"pageSize"` /*  每页条数  */
	Total    *int32                             `json:"total"`    /*  总数  */
	Pages    *string                            `json:"pages"`    /*  总页数  */
	List     []*CtiamQueryReturnObjListResponse `json:"list"`     /*  权限列表  */
}

type CtiamQueryReturnObjListResponse struct {
	Id                *string `json:"id"`                /*  权限id  */
	RegionId          *string `json:"regionId"`          /*  资源池id  */
	PolicyId          *string `json:"policyId"`          /*  策略id  */
	AccountGroupId    *string `json:"accountGroupId"`    /*  用户组id  */
	AccountId         *string `json:"accountId"`         /*  账号id  */
	RangeType         *string `json:"rangeType"`         /*  授权范围  */
	PolicyName        *string `json:"policyName"`        /*  策略名称  */
	PolicyDescription *string `json:"policyDescription"` /*  策略描述  */
	CreateTime        *string `json:"createTime"`        /*  创建时间  */
	GroupId           *string `json:"groupId"`           /*  用户组id  */
	GroupName         *string `json:"groupName"`         /*  用户组名称  */
	GroupIntro        *string `json:"groupIntro"`        /*  用户组描述  */
	IsRoot            *string `json:"isRoot"`            /*  是否主用户组  */
}
