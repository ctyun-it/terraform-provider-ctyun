package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetSdwanAclRuleApi
/* 查询访问控制规则信息 */
type SdwanGetSdwanAclRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanAclRuleApi(client *core.CtyunClient) *SdwanGetSdwanAclRuleApi {
	return &SdwanGetSdwanAclRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/acl-rule/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanAclRuleApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanAclRuleRequest) (*SdwanGetSdwanAclRuleResponse, error) {
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
	if req.AclID != nil && *req.AclID != "" {
		ctReq.AddParam("aclID", *req.AclID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanAclRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanAclRuleRequest struct {
	PageNo   int32   `json:"pageNo"`           /*  页码  */
	PageSize int32   `json:"pageSize"`         /*  每页记录数目  */
	Search   *string `json:"search,omitempty"` /*  模糊查询  */
	AclID    *string `json:"aclID,omitempty"`  /*  ACL ID  */
}

type SdwanGetSdwanAclRuleResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetSdwanAclRuleReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
type SdwanGetSdwanAclRuleReturnObjResponse struct {
	Result       []SdwanGetSdwanAclRuleReturnObjResponseResult `json:"result"`       /*  rule ID  */
	CurrentCount int32                                         `json:"currentCount"` /*  acl id  */
	TotalCount   int32                                         `json:"totalCount"`   /*  本参数表示控制方向<br/><br/>取值范围:<br/>in:入方向<br/>out:出方向  */
}
type SdwanGetSdwanAclRuleReturnObjResponseResult struct {
	AclRuleID        *string `json:"aclRuleID"`        /*  rule ID  */
	AclID            *string `json:"aclID"`            /*  acl id  */
	Direction        *string `json:"direction"`        /*  本参数表示控制方向<br/><br/>取值范围:<br/>in:入方向<br/>out:出方向  */
	Action           *string `json:"action"`           /*  本参数表示策略类型<br/><br/>取值范围:<br/>allow:允许<br/>deny:拒绝  */
	Protocol         *string `json:"protocol"`         /*  本参数表示协议类型<br/><br/>取值范围:<br/>udp:UDP<br/>icmp:ICMP</br>all:ALL</br>tcp:TCP  */
	FuserLastUpdated *string `json:"fuserLastUpdated"` /*  用户最近更新时间  */
	SrcCidr          *string `json:"srcCidr"`          /*  源网段  */
	DstCidr          *string `json:"dstCidr"`          /*  目的网段  */
	SrcPortRange     *string `json:"srcPortRange"`     /*  源端口范围  */
	DstPortRange     *string `json:"dstPortRange"`     /*  目的端口范围  */
	Priority         int32   `json:"priority"`         /*  priority  */
	Status           *string `json:"status"`           /*  本参数表示状态<br/><br/>取值范围:<br/>normal:状态正常<br/>creating:正在创建</br>deleting:正在删除</br>failed:创建失败</br>error:删除失败  */
}
