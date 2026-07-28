package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanSdwanListAclEdgeApi
/* 查询绑定好acl的盒子信息 */
type SdwanSdwanListAclEdgeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanListAclEdgeApi(client *core.CtyunClient) *SdwanSdwanListAclEdgeApi {
	return &SdwanSdwanListAclEdgeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/acl-edge/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanListAclEdgeApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanListAclEdgeRequest) (*SdwanSdwanListAclEdgeResponse, error) {
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
	ctReq.AddParam("aclID", req.AclID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanSdwanListAclEdgeResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanListAclEdgeRequest struct {
	PageNo   int32   `json:"pageNo"`           /*  页码  */
	PageSize int32   `json:"pageSize"`         /*  每页计算数目  */
	Search   *string `json:"search,omitempty"` /*  模糊查询  */
	AclID    string  `json:"aclID"`            /*  acl ID  */
}

type SdwanSdwanListAclEdgeResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                 `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                 `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                 `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanListAclEdgeReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanListAclEdgeReturnObjResponse struct {
	Message      *string                                         `json:"message"`      /*  message  */
	TotalCount   int32                                           `json:"totalCount"`   /*  总数  */
	CurrentCount int32                                           `json:"currentCount"` /*  当前页数量  */
	Code         *string                                         `json:"code"`         /*  code  */
	Result       []*SdwanSdwanListAclEdgeReturnObjResultResponse `json:"result"`       /*  结果列表  */
}

type SdwanSdwanListAclEdgeReturnObjResultResponse struct {
	EdgeName  *string `json:"edgeName"`  /*  盒子名称  */
	EdgeID    *string `json:"edgeID"`    /*  盒子ID  */
	SdwanName *string `json:"sdwanName"` /*  sdwan名称  */
	Status    *string `json:"status"`    /*  状态  */
	BindTime  *string `json:"bindTime"`  /*  盒子绑定acl时间  */
}
