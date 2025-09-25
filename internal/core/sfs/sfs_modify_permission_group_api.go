package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsModifyPermissionGroupApi
/* 修改权限组
 */type SfsModifyPermissionGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsModifyPermissionGroupApi(client *core.CtyunClient) *SfsModifyPermissionGroupApi {
	return &SfsModifyPermissionGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/permission-group/modify",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsModifyPermissionGroupApi) Do(ctx context.Context, credential core.Credential, req *SfsModifyPermissionGroupRequest) (*SfsModifyPermissionGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsModifyPermissionGroupRequest
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
	var resp SfsModifyPermissionGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsModifyPermissionGroupRequest struct {
	PermissionGroupID string `json:"permissionGroupID,omitempty"` /*  权限组fuid  */
	RegionID          string `json:"regionID,omitempty"`          /*  资源池 ID  */
	Name              string `json:"name,omitempty"`              /*  权限组名字，不可重复；长度为2-63字符，只能由数字、字母、-组成，不能以数字和-开头、且不能以-结尾。  */
	Description       string `json:"description,omitempty"`       /*  权限组描述信息，长度为0-128字符。  */
}

type SfsModifyPermissionGroupResponse struct {
	StatusCode  int32                                      `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                     `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                     `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsModifyPermissionGroupReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                     `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                                     `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsModifyPermissionGroupReturnObjResponse struct {
	OperationID string `json:"operationID"` /*  修改权限组的操作号  */
}
