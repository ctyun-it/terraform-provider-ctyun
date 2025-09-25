package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsUnbindVpcSfsApi
/* 根据文件系统ID、VPC及权限组ID ，解绑权限组
 */type SfsSfsUnbindVpcSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsUnbindVpcSfsApi(client *core.CtyunClient) *SfsSfsUnbindVpcSfsApi {
	return &SfsSfsUnbindVpcSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/vpc-unbind-permission",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsUnbindVpcSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsUnbindVpcSfsRequest) (*SfsSfsUnbindVpcSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsUnbindVpcSfsRequest
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
	var resp SfsSfsUnbindVpcSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsUnbindVpcSfsRequest struct {
	PermissionGroupFuid string `json:"permissionGroupFuid,omitempty"` /*  权限组ID  */
	RegionID            string `json:"regionID,omitempty"`            /*  资源池 ID  */
	SfsUID              string `json:"sfsUID,omitempty"`              /*  弹性文件功能系统唯一 ID  */
	VpcID               string `json:"vpcID,omitempty"`               /*  Vpc ID  */
}

type SfsSfsUnbindVpcSfsResponse struct {
	StatusCode  int32                                `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                               `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                               `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsUnbindVpcSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                               `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                               `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsUnbindVpcSfsReturnObjResponse struct {
	OperationID string `json:"operationID"` /*  文件系统VPC解绑权限组的操作号  */
}
