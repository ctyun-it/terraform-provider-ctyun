package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsDeletePermissionSfsApi
/* 弹性文件删除权限组
 */type SfsSfsDeletePermissionSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsDeletePermissionSfsApi(client *core.CtyunClient) *SfsSfsDeletePermissionSfsApi {
	return &SfsSfsDeletePermissionSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/permission-group/delete-permission-group",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsDeletePermissionSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsDeletePermissionSfsRequest) (*SfsSfsDeletePermissionSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsDeletePermissionSfsRequest
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
	var resp SfsSfsDeletePermissionSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsDeletePermissionSfsRequest struct {
	RegionID            string `json:"regionID,omitempty"`            /*  资源池 ID  */
	PermissionGroupFuid string `json:"permissionGroupFuid,omitempty"` /*  权限组Fuid  */
}

type SfsSfsDeletePermissionSfsResponse struct {
	StatusCode  int32                                       `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                      `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                      `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsDeletePermissionSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                      `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码，参考[结果码]  */
	Error       string                                      `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsDeletePermissionSfsReturnObjResponse struct {
	OperationID string `json:"operationID"` /*  删除权限组的操作号  */
}
