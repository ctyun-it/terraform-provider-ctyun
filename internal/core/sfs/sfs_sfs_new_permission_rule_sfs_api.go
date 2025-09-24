package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsNewPermissionRuleSfsApi
/* 根据资源池ID和权限组fuid创建权限组规则
 */type SfsSfsNewPermissionRuleSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsNewPermissionRuleSfsApi(client *core.CtyunClient) *SfsSfsNewPermissionRuleSfsApi {
	return &SfsSfsNewPermissionRuleSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/permission-rule/new-permission-rule",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsNewPermissionRuleSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsNewPermissionRuleSfsRequest) (*SfsSfsNewPermissionRuleSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsNewPermissionRuleSfsRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsNewPermissionRuleSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsNewPermissionRuleSfsRequest struct {
	PermissionGroupFuid    string `json:"permissionGroupFuid,omitempty"`    /*  权限组Fuid  */
	RegionID               string `json:"regionID,omitempty"`               /*  资源池ID  */
	AuthAddr               string `json:"authAddr,omitempty"`               /*  授权地址。支持IPv4和IPv6两种网络类型，可填写单个IP 或者单个网段。同一权限组内，授权地址不能重复  */
	RwPermission           string `json:"rwPermission,omitempty"`           /*  读写权限控制：ro（只读）、rw（读写）。注：当客户端从读写权限改为只读权限再改回读写权限时，需要重新挂载客户端  */
	UserPermission         string `json:"userPermission,omitempty"`         /*  用户权限：no_root_squash（不匿名用户）  */
	PermissionRulePriority int32  `json:"permissionRulePriority,omitempty"` /*  优先级，有效范围为1-400。当同一个权限组内单个 IP 与网段中包含的 IP 的权限有冲突时，会生效优先级高的规则。注：优先级不可重复  */
}

type SfsSfsNewPermissionRuleSfsResponse struct {
	StatusCode  int32                                        `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                       `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                       `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsNewPermissionRuleSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                       `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码，参考[结果码]  */
	Error       string                                       `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsNewPermissionRuleSfsReturnObjResponse struct {
	OperationID string `json:"operationID"` /*  创建权限组规则的操作号  */
}
