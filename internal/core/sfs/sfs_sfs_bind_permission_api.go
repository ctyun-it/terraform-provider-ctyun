package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsBindPermissionApi
/* 根据文件系统ID、VPC及权限组ID ，绑定权限组
 */type SfsSfsBindPermissionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsBindPermissionApi(client *core.CtyunClient) *SfsSfsBindPermissionApi {
	return &SfsSfsBindPermissionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/bind-permission",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsBindPermissionApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsBindPermissionRequest) (*SfsSfsBindPermissionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsBindPermissionRequest
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
	var resp SfsSfsBindPermissionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsBindPermissionRequest struct {
	PermissionGroupID string `json:"permissionGroupID,omitempty"` /*  权限组的fuid  */
	RegionID          string `json:"regionID,omitempty"`          /*  资源池 ID  */
	UID               string `json:"UID,omitempty"`               /*  文件系统ID  */
	VpcID             string `json:"vpcID,omitempty"`             /*  vpcID  */
}

type SfsSfsBindPermissionResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                 `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                 `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsBindPermissionReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                 `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsBindPermissionReturnObjResponse struct {
	OperationID string `json:"operationID"` /*  文件系统VPC绑定权限组的操作号  */
}
