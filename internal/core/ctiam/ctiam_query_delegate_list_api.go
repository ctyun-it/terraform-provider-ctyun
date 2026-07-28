package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtiamQueryDelegateListApi
/* 查询指定账号下的云服务委托或内联委托列表 */
type CtiamQueryDelegateListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryDelegateListApi(client *core.CtyunClient) *CtiamQueryDelegateListApi {
	return &CtiamQueryDelegateListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/delegate/queryDelegateList",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryDelegateListApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryDelegateListRequest) (*CtiamQueryDelegateListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("accountId", req.AccountId)
	if req.ServiceCode != nil && *req.ServiceCode != "" {
		ctReq.AddParam("serviceCode", *req.ServiceCode)
	}
	if *req.RawType != 0 {
		ctReq.AddParam("type", strconv.FormatInt(int64(*req.RawType), 10))
	}
	if req.Name != nil && *req.Name != "" {
		ctReq.AddParam("name", *req.Name)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamQueryDelegateListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryDelegateListRequest struct {
	AccountId   string  `json:"accountId"`             /*  账户id  */
	ServiceCode *string `json:"serviceCode,omitempty"` /*  服务编码，查询特定服务的委托时可以填写该参数，例如查询ctiam的委托时，填写对应的服务编码为ctiam。一般产品加载ctiam或openAPI时会分配产品编码。  */
	RawType     *int32  `json:"type,omitempty"`        /*  类型 1：云服务委托, 3：服务内联委托，不传：默认为云服务委托类型 1。只支持查询云服务委托或内联委托  */
	Name        *string `json:"name,omitempty"`        /*  委托名称  */
}

type CtiamQueryDelegateListResponse struct {
	StatusCode *string                                  `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                  `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                  `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamQueryDelegateListReturnObjResponse `json:"returnObj"`  /*  返回体  */
}

type CtiamQueryDelegateListReturnObjResponse struct {
	Result []*CtiamQueryDelegateListReturnObjResultResponse `json:"result"` /*  返回结果  */
}

type CtiamQueryDelegateListReturnObjResultResponse struct {
	Name         *string `json:"name"`         /*  委托角色名称  */
	AccountId    *string `json:"accountId"`    /*  委托账号  */
	AssumeUserId *string `json:"assumeUserId"` /*  被委托用户id  */
	RawType      *int32  `json:"type"`         /*  类型 0账号级委托、1 云服务委托、2身份供应商、3服务内联委托  */
	CreateTime   *string `json:"createTime"`   /*  创建时间  */
	UpdateTime   *string `json:"updateTime"`   /*  更新时间  */
}
