package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsModifyPermissionRuleSfsApi
/* 根据资源池ID和权限组规则fuid修改权限组规则
 */type SfsSfsModifyPermissionRuleSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsModifyPermissionRuleSfsApi(client *core.CtyunClient) *SfsSfsModifyPermissionRuleSfsApi {
	return &SfsSfsModifyPermissionRuleSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/permission-rule/modify-permission-rule",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsModifyPermissionRuleSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsModifyPermissionRuleSfsRequest) (*SfsSfsModifyPermissionRuleSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsModifyPermissionRuleSfsRequest
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
	var resp SfsSfsModifyPermissionRuleSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsModifyPermissionRuleSfsRequest struct {
	PermissionRuleFuid     string `json:"permissionRuleFuid,omitempty"`     /*  权限组规则的fuid  */
	RegionID               string `json:"regionID,omitempty"`               /*  资源池 ID  */
	AuthAddr               string `json:"authAddr,omitempty"`               /*  授权地址，可用于区分子网及具体虚机等。有效值范围：ipv4、ipv6  */
	RwPermission           string `json:"rwPermission,omitempty"`           /*  读写权限控制。有效值范围：rw、ro  */
	UserPermission         string `json:"userPermission,omitempty"`         /*  nfs 访问用户映射。有效值范围：no_root_squash  */
	PermissionRulePriority int32  `json:"permissionRulePriority,omitempty"` /*  优先级 1 为最高 1-400。不同规则优先级不能重复  */
}

type SfsSfsModifyPermissionRuleSfsResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码，参考[结果码]  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}
