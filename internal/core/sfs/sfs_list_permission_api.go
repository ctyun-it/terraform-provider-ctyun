package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsListPermissionApi
/* 查询文件系统权限组与VPC绑定关系
 */type SfsListPermissionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsListPermissionApi(client *core.CtyunClient) *SfsListPermissionApi {
	return &SfsListPermissionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list-permission",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsListPermissionApi) Do(ctx context.Context, credential core.Credential, req *SfsListPermissionRequest) (*SfsListPermissionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("permissionGroupFuid", req.PermissionGroupFuid)
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sfsUID", req.SfsUID)
	ctReq.AddParam("vpcID", req.VpcID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsListPermissionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsListPermissionRequest struct {
	PermissionGroupFuid string `json:"permissionGroupFuid,omitempty"` /*  权限组ID  */
	RegionID            string `json:"regionID,omitempty"`            /*  资源池 ID  */
	SfsUID              string `json:"sfsUID,omitempty"`              /*  弹性文件功能系统唯一 ID  */
	VpcID               string `json:"vpcID,omitempty"`               /*  vpcID  */
}

type SfsListPermissionResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}
