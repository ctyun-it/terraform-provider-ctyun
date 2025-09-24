package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsModifyPermissionSfsApi
/* 弹性文件支持修改权限组的名称或描述
 */type SfsSfsModifyPermissionSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsModifyPermissionSfsApi(client *core.CtyunClient) *SfsSfsModifyPermissionSfsApi {
	return &SfsSfsModifyPermissionSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/permission-group/modify-permission-group",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsModifyPermissionSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsModifyPermissionSfsRequest) (*SfsSfsModifyPermissionSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsModifyPermissionSfsRequest
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
	var resp SfsSfsModifyPermissionSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsModifyPermissionSfsRequest struct {
	PermissionGroupFuid        string `json:"permissionGroupFuid,omitempty"`        /*  权限组Fuid  */
	RegionID                   string `json:"regionID,omitempty"`                   /*  资源池ID  */
	PermissionGroupName        string `json:"permissionGroupName,omitempty"`        /*  预修改的权限组名称。permissionGroupName和permissionGroupDescription至少输入一个。长度为2-63字符，只能由数字、字母、"-"组成，不能以数字和"-"开头、且不能以"-"结尾  */
	PermissionGroupDescription string `json:"permissionGroupDescription,omitempty"` /*  预修改的权限组描述信息。permissionGroupName和permissionGroupDescription至少输入一个。长度为0-128字符  */
}

type SfsSfsModifyPermissionSfsResponse struct {
	StatusCode  int32                                       `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                      `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                      `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsModifyPermissionSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                      `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码，参考[结果码]  */
	Error       string                                      `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsModifyPermissionSfsReturnObjResponse struct {
	OperationID string `json:"operationID"` /*  修改权限组的操作号  */
}
