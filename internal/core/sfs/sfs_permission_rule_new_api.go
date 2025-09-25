package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsPermissionRuleNewApi
/* 创建权限组规则
 */type SfsPermissionRuleNewApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsPermissionRuleNewApi(client *core.CtyunClient) *SfsPermissionRuleNewApi {
	return &SfsPermissionRuleNewApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/permission-rule/new",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsPermissionRuleNewApi) Do(ctx context.Context, credential core.Credential, req *SfsPermissionRuleNewRequest) (*SfsPermissionRuleNewResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsPermissionRuleNewRequest
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
	var resp SfsPermissionRuleNewResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsPermissionRuleNewRequest struct {
	PermissionGroupID string `json:"permissionGroupID,omitempty"` /*  权限组的fuid  */
	RegionID          string `json:"regionID,omitempty"`          /*  资源池 ID  */
	AuthAddr          string `json:"authAddr,omitempty"`          /*  授权地址。有效值范围：ipv4、ipv6  */
	RwPermission      string `json:"rwPermission,omitempty"`      /*  读写权限控制。有效值范围：rw、ro  */
	UserPermission    string `json:"userPermission,omitempty"`    /*  nfs 访问用户映射。有效值范围：no_root_squash  */
	Priority          int32  `json:"priority,omitempty"`          /*  优先级。有效值范围：1-400  */
}

type SfsPermissionRuleNewResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                 `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                 `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsPermissionRuleNewReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                 `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsPermissionRuleNewReturnObjResponse struct {
	OperationID string `json:"operationID"` /*  创建权限组规则的操作号  */
}
