package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SfsSfsListPermissionSfsApi
/* 根据资源池ID和权限组的Fuid返回权限组描述信息
 */type SfsSfsListPermissionSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsListPermissionSfsApi(client *core.CtyunClient) *SfsSfsListPermissionSfsApi {
	return &SfsSfsListPermissionSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/permission-group/list-permission-group",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsListPermissionSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsListPermissionSfsRequest) (*SfsSfsListPermissionSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.PermissionGroupFuid != "" {
		ctReq.AddParam("permissionGroupFuid", req.PermissionGroupFuid)
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
	var resp SfsSfsListPermissionSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsListPermissionSfsRequest struct {
	RegionID            string `json:"regionID,omitempty"`            /*  资源池 ID  */
	PermissionGroupFuid string `json:"permissionGroupFuid,omitempty"` /*  权限组Fuid  */
	PageSize            int32  `json:"pageSize,omitempty"`            /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
	PageNo              int32  `json:"pageNo,omitempty"`              /*  页码，取值范围：正整数（≥1），注：默认值为1  */
}

type SfsSfsListPermissionSfsResponse struct {
	StatusCode  int32                                     `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                    `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                    `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsListPermissionSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                    `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码，参考[结果码]  */
	Error       string                                    `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsListPermissionSfsReturnObjResponse struct {
	TotalCount   int32                                           `json:"totalCount"`   /*  权限组的总个数  */
	CurrentCount int32                                           `json:"currentCount"` /*  当前页记录数目  */
	List         []*SfsSfsListPermissionSfsReturnObjListResponse `json:"list"`         /*  权限组信息列表  */
	PageSize     int32                                           `json:"pageSize"`     /*  每页记录数目  */
	PageNo       int32                                           `json:"pageNo"`       /*  页码  */
}

type SfsSfsListPermissionSfsReturnObjListResponse struct {
	PermissionGroupFuid        string `json:"permissionGroupFuid"`        /*  	权限组Fuid  */
	CreateTime                 string `json:"createTime"`                 /*  创建时间，UTC 时间  */
	UpdateTime                 string `json:"updateTime"`                 /*  更新时间，UTC 时间  */
	UserID                     string `json:"userID"`                     /*  租户 ID  */
	RegionID                   string `json:"regionID"`                   /*  所属资源池 ID  */
	SfsCount                   int32  `json:"sfsCount"`                   /*  绑定的文件系统个数  */
	PermissionRuleCount        int32  `json:"permissionRuleCount"`        /*  权限组规则个数  */
	PermissionGroupIsDefault   *bool  `json:"permissionGroupIsDefault"`   /*  是否为默认权限组  */
	PermissionGroupID          string `json:"permissionGroupID"`          /*  权限组底层 ID  */
	PermissionGroupName        string `json:"permissionGroupName"`        /*  权限组名称  */
	PermissionGroupDescription string `json:"permissionGroupDescription"` /*  权限组描述  */
}
