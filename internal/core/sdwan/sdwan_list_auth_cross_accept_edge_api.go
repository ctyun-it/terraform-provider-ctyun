package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanListAuthCrossAcceptEdgeApi
/* 查找跨账号被授权站点信息 */
type SdwanListAuthCrossAcceptEdgeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanListAuthCrossAcceptEdgeApi(client *core.CtyunClient) *SdwanListAuthCrossAcceptEdgeApi {
	return &SdwanListAuthCrossAcceptEdgeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/auth-cross-accept-edge/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanListAuthCrossAcceptEdgeApi) Do(ctx context.Context, credential core.Credential, req *SdwanListAuthCrossAcceptEdgeRequest) (*SdwanListAuthCrossAcceptEdgeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.ProjectID != nil && *req.ProjectID != "" {
		ctReq.AddParam("projectID", *req.ProjectID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanListAuthCrossAcceptEdgeResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanListAuthCrossAcceptEdgeRequest struct {
	PageNo    int32   `json:"pageNo"`              /*  页码  */
	PageSize  int32   `json:"pageSize"`            /*  每页记录数目  */
	ProjectID *string `json:"projectID,omitempty"` /*  企业项目  */
}

type SdwanListAuthCrossAcceptEdgeResponse struct {
	StatusCode         int32   `json:"statusCode"`         /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode          *string `json:"errorCode"`          /*  业务细分码，为product.module.code三段式码  */
	Message            *string `json:"message"`            /*  失败时的错误描述，一般为英文描述  */
	Description        *string `json:"description"`        /*  失败时的错误描述，一般为中文描述  */
	TotalCount         int32   `json:"totalCount"`         /*  总数  */
	CurrentCount       int32   `json:"currentCount"`       /*  当前页数量  */
	EdgeID             *string `json:"edgeID"`             /*  站点ID  */
	EdgeName           *string `json:"edgeName"`           /*  站点名称  */
	SdwanID            *string `json:"sdwanID"`            /*  sdwan id  */
	SdwanName          *string `json:"sdwanName"`          /*  sdwan名称  */
	GrantCustomerEmail *string `json:"grantCustomerEmail"` /*  授权方账号邮箱  */
	GrantCustomerID    *string `json:"grantCustomerID"`    /*  授权方账号ID  */
	GrantTime          *string `json:"grantTime"`          /*  授权时间  */
	Error              *string `json:"error"`              /*  业务细分码，为product.module.code三段式码  */
}
