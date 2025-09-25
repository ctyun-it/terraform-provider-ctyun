package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SfsSfsListPermissionRuleSfsApi
/* 根据资源池ID及权限组Fuid或权限组规则Fuid，返回权限组规则描述信息
 */type SfsSfsListPermissionRuleSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsListPermissionRuleSfsApi(client *core.CtyunClient) *SfsSfsListPermissionRuleSfsApi {
	return &SfsSfsListPermissionRuleSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/permission-rule/list-permission-rule",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsListPermissionRuleSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsListPermissionRuleSfsRequest) (*SfsSfsListPermissionRuleSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.PermissionGroupFuid != "" {
		ctReq.AddParam("permissionGroupFuid", req.PermissionGroupFuid)
	}
	if req.PermissionRuleFuid != "" {
		ctReq.AddParam("permissionRuleFuid", req.PermissionRuleFuid)
	}
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
	var resp SfsSfsListPermissionRuleSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsListPermissionRuleSfsRequest struct {
	RegionID            string `json:"regionID,omitempty"`            /*  资源池 ID  */
	PermissionGroupFuid string `json:"permissionGroupFuid,omitempty"` /*  权限组Fuid，permissionGroupFuid和permissionRuleFuid至少存在一个  */
	PermissionRuleFuid  string `json:"permissionRuleFuid,omitempty"`  /*  权限组规则Fuid，permissionGroupFuid和permissionRuleFuid至少存在一个  */
	PageSize            int32  `json:"pageSize,omitempty"`            /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
	PageNo              int32  `json:"pageNo,omitempty"`              /*  页码，取值范围：正整数（≥1），注：默认值为1  */
}

type SfsSfsListPermissionRuleSfsResponse struct {
	StatusCode  int32                                         `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                        `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                        `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsListPermissionRuleSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                        `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码，参考[结果码]  */
	Error       string                                        `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsListPermissionRuleSfsReturnObjResponse struct {
	TotalCount   int32                                               `json:"totalCount"`   /*  权限组下的规则总数  */
	CurrentCount int32                                               `json:"currentCount"` /*  当前查询到的规则个数  */
	List         []*SfsSfsListPermissionRuleSfsReturnObjListResponse `json:"list"`         /*  规则列表  */
	PageSize     int32                                               `json:"pageSize"`     /*  每页记录数目  */
	PageNo       int32                                               `json:"pageNo"`       /*  页码  */
}

type SfsSfsListPermissionRuleSfsReturnObjListResponse struct {
	PermissionRuleFuid      string `json:"permissionRuleFuid"`      /*  权限组规则ID  */
	UpdateTime              string `json:"updateTime"`              /*  更新时间。UTC时间  */
	UserID                  string `json:"userID"`                  /*  租户ID  */
	RegionID                string `json:"regionID"`                /*  资源池ID。中间层使用的  */
	PermissionGroupID       string `json:"permissionGroupID"`       /*  权限组底层ID  */
	PermissionGroupFuid     string `json:"permissionGroupFuid"`     /*  权限组Fuid  */
	PermissionRuleID        string `json:"permissionRuleID"`        /*  权限组规则底层ID  */
	AuthAddr                string `json:"authAddr"`                /*  授权地址，可用于区分子网及具体虚机等  */
	RwPermission            string `json:"rwPermission"`            /*  读写权限控制  */
	UserPermission          string `json:"userPermission"`          /*  用户权限  */
	PermissionRulePriority  int32  `json:"permissionRulePriority"`  /*  优先级  */
	PermissionRuleIsDefault *bool  `json:"permissionRuleIsDefault"` /*  是否为默认规则  */
}
