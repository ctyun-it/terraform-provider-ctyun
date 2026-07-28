package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetSdwanAppGroupApi
/* 应用组查询 */
type SdwanGetSdwanAppGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanAppGroupApi(client *core.CtyunClient) *SdwanGetSdwanAppGroupApi {
	return &SdwanGetSdwanAppGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/app-group/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanAppGroupApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanAppGroupRequest) (*SdwanGetSdwanAppGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Search != nil && *req.Search != "" {
		ctReq.AddParam("search", *req.Search)
	}
	if req.GroupName != nil && *req.GroupName != "" {
		ctReq.AddParam("groupName", *req.GroupName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanAppGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanAppGroupRequest struct {
	PageNo    int32   `json:"pageNo"`              /*  页码  */
	PageSize  int32   `json:"pageSize"`            /*  每页记录数目  */
	Search    *string `json:"search,omitempty"`    /*  模糊查询  */
	GroupName *string `json:"groupName,omitempty"` /*  应用名称  */
}

type SdwanGetSdwanAppGroupResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                 `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                 `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                 `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetSdwanAppGroupReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetSdwanAppGroupReturnObjResponse struct {
	Message      *string                                         `json:"message"`      /*  message  */
	TotalCount   int32                                           `json:"totalCount"`   /*  总数  */
	CurrentCount int32                                           `json:"currentCount"` /*  当前页数量  */
	Code         *string                                         `json:"code"`         /*  状态码  */
	Result       []*SdwanGetSdwanAppGroupReturnObjResultResponse `json:"result"`       /*  列表  */
}

type SdwanGetSdwanAppGroupReturnObjResultResponse struct {
	GroupName   *string `json:"groupName"`   /*  应用组名称  */
	GroupID     *string `json:"groupID"`     /*  应用组ID  */
	GroupType   *string `json:"groupType"`   /*  本参数表示应用组类型<br/><br/>取值范围:<br/>custom:自定义<br/>system:系统  */
	CustomerID  *string `json:"customerID"`  /*  用户ID  */
	CreatedTime *string `json:"createdTime"` /*  创建时间  */
}
