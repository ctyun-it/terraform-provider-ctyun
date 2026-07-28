package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanListAclEdgeRelationApi
/* 查询acl与盒子之间的关系信息 */
type SdwanListAclEdgeRelationApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanListAclEdgeRelationApi(client *core.CtyunClient) *SdwanListAclEdgeRelationApi {
	return &SdwanListAclEdgeRelationApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/acl-edge-relation/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanListAclEdgeRelationApi) Do(ctx context.Context, credential core.Credential, req *SdwanListAclEdgeRelationRequest) (*SdwanListAclEdgeRelationResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("relation", req.Relation)
	if req.SdwanID != nil && *req.SdwanID != "" {
		ctReq.AddParam("sdwanID", *req.SdwanID)
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
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanListAclEdgeRelationResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanListAclEdgeRelationRequest struct {
	Relation string  `json:"relation"`          /*  本参数表示智能网关和acl之间的关系<br/><br/>取值范围:<br/>bind:绑定<br/>unbind:解除绑定  */
	SdwanID  *string `json:"sdwanID,omitempty"` /*  sdwan的id  */
	PageNo   int32   `json:"pageNo"`            /*  页码  */
	PageSize int32   `json:"pageSize"`          /*  每页数目  */
	Search   *string `json:"search,omitempty"`  /*  模糊查询  */
}

type SdwanListAclEdgeRelationResponse struct {
	StatusCode   int32   `json:"statusCode"`   /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode    *string `json:"errorCode"`    /*  业务细分码，为product.module.code三段式码  */
	Message      *string `json:"message"`      /*  失败时的错误描述，一般为英文描述  */
	Description  *string `json:"description"`  /*  失败时的错误描述，一般为中文描述  */
	TotalCount   int32   `json:"totalCount"`   /*  总数  */
	CurrentCount int32   `json:"currentCount"` /*  当前页数量  */
	EdgeName     *string `json:"edgeName"`     /*  盒子名称  */
	EdgeID       *string `json:"edgeID"`       /*  盒子ID  */
	SdwanName    *string `json:"sdwanName"`    /*  sdwan名称  */
	SdwanID      *string `json:"sdwanID"`      /*  sdwan的ID  */
	AclID        *string `json:"aclID"`        /*  ACL ID  */
	AclName      *string `json:"aclName"`      /*  acl规则名称  */
	Error        *string `json:"error"`        /*  业务细分码，为product.module.code三段式码  */
}
