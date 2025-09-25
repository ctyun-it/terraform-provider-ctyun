package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SfsPermissionRuleListApi
/* 返回权限组规则描述信息
 */type SfsPermissionRuleListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsPermissionRuleListApi(client *core.CtyunClient) *SfsPermissionRuleListApi {
	return &SfsPermissionRuleListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/permission-rule/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsPermissionRuleListApi) Do(ctx context.Context, credential core.Credential, req *SfsPermissionRuleListRequest) (*SfsPermissionRuleListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("permissionGroupID", req.PermissionGroupID)
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsPermissionRuleListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsPermissionRuleListRequest struct {
	RegionID          string `json:"regionID,omitempty"`          /*  资源池 ID  */
	PermissionGroupID string `json:"permissionGroupID,omitempty"` /*  权限组的fuid  */
	PageSize          int32  `json:"pageSize,omitempty"`          /*  每页个数。默认为10  */
	PageNo            int32  `json:"pageNo,omitempty"`            /*  页数。默认为1  */
}

type SfsPermissionRuleListResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                  `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                  `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsPermissionRuleListReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                  `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsPermissionRuleListReturnObjResponse struct {
	TotalCount   int32                                         `json:"totalCount"`   /*  权限组下的规则总数  */
	CurrentCount int32                                         `json:"currentCount"` /*  当前查询到的规则个数  */
	List         []*SfsPermissionRuleListReturnObjListResponse `json:"list"`         /*  查询到的权限组规则列表。参考list  */
	PageSize     int32                                         `json:"pageSize"`     /*  每页个数  */
	PageNo       int32                                         `json:"pageNo"`       /*  页数  */
}

type SfsPermissionRuleListReturnObjListResponse struct {
	Fuid                  string `json:"fuid"`                  /*  规则的 ID  */
	Fuser_last_updated    string `json:"fuser_last_updated"`    /*  更新时间。utc 时间  */
	User_id               string `json:"user_id"`               /*  租户 ID  */
	Region_id             string `json:"region_id"`             /*  资源池 ID。中间层用的  */
	Permission_group_id   string `json:"permission_group_id"`   /*  权限组底层 id。yacos用的id  */
	Permission_group_fuid string `json:"permission_group_fuid"` /*  权限组 fuid。权限组查询接口中的fuid  */
	Rule_id               string `json:"rule_id"`               /*  权限组规则底层 id  */
	Rule_auth_address     string `json:"rule_auth_address"`     /*  授权地址，可用于区分子网及具体虚机等  */
	Rule_rw_permission    string `json:"rule_rw_permission"`    /*  读写权限控制  */
	Rule_user_permission  string `json:"rule_user_permission"`  /*  nfs 访问用户映射  */
	Rule_priority_id      int32  `json:"rule_priority_id"`      /*  优先级  */
	Rule_is_default       *bool  `json:"rule_is_default"`       /*  优先级  */
}
