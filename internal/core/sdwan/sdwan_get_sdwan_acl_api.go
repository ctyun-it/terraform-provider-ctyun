package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetSdwanAclApi
/* 查询指定访问控制实例信息 */
type SdwanGetSdwanAclApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanAclApi(client *core.CtyunClient) *SdwanGetSdwanAclApi {
	return &SdwanGetSdwanAclApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/acl/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanAclApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanAclRequest) (*SdwanGetSdwanAclResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.ProjectID != nil && *req.ProjectID != "" {
		ctReq.AddParam("projectID", *req.ProjectID)
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Search != nil && *req.Search != "" {
		ctReq.AddParam("search", *req.Search)
	}
	if req.AclID != nil && *req.AclID != "" {
		ctReq.AddParam("aclID", *req.AclID)
	}
	if req.AclName != nil && *req.AclName != "" {
		ctReq.AddParam("aclName", *req.AclName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanAclResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanAclRequest struct {
	ProjectID *string `json:"projectID,omitempty"` /*  企业项目  */
	PageNo    int32   `json:"pageNo"`              /*  页码  */
	PageSize  int32   `json:"pageSize"`            /*  每页记录数目  */
	Search    *string `json:"search,omitempty"`    /*  模糊查询  */
	AclID     *string `json:"aclID,omitempty"`     /*  ACL ID  */
	AclName   *string `json:"aclName,omitempty"`   /*  名称  */
}

type SdwanGetSdwanAclResponse struct {
	StatusCode  int32                              `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                            `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                            `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                            `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetSdwanAclReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                            `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetSdwanAclReturnObjResponse struct {
	Message      *string                                    `json:"message"`      /*  message  */
	TotalCount   int32                                      `json:"totalCount"`   /*  总数  */
	CurrentCount int32                                      `json:"currentCount"` /*  当前页数量  */
	Code         *string                                    `json:"code"`         /*  code  */
	Result       []*SdwanGetSdwanAclReturnObjResultResponse `json:"result"`       /*  列表  */
}

type SdwanGetSdwanAclReturnObjResultResponse struct {
	AclID          *string `json:"aclID"`          /*  ACL ID  */
	Name           *string `json:"aclName"`        /*  名称  */
	ProjectID      *string `json:"projectID"`      /*  企业项目  */
	RuleCount      int32   `json:"ruleCount"`      /*  规则数量  */
	RuleErrorCount int32   `json:"ruleErrorCount"` /*  错误的规则数量  */
	EdgeCount      int32   `json:"edgeCount"`      /*  绑定的盒子数量  */
	EdgeErrorCount int32   `json:"edgeErrorCount"` /*  绑定失败的盒子数量  */
}
