package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsListVpcSfsApi
/* 根据文件系统ID查询文件系统下的vpc列表、vpc绑定的权限组
 */type SfsSfsListVpcSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsListVpcSfsApi(client *core.CtyunClient) *SfsSfsListVpcSfsApi {
	return &SfsSfsListVpcSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list-vpc-permission",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsListVpcSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsListVpcSfsRequest) (*SfsSfsListVpcSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sfsUID", req.SfsUID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsListVpcSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsListVpcSfsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  弹性文件功能系统唯一 ID  */
}

type SfsSfsListVpcSfsResponse struct {
	StatusCode  int32                              `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                             `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                             `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsListVpcSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                             `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                             `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsListVpcSfsReturnObjResponse struct {
	List       []*SfsSfsListVpcSfsReturnObjListResponse `json:"list"`       /*  返回的文件列表  */
	TotalCount int32                                    `json:"totalCount"` /*  用户弹性文件下VPC总数  */
}

type SfsSfsListVpcSfsReturnObjListResponse struct {
	UserID                     string `json:"UserID"`                     /*  租户ID  */
	SfsUID                     string `json:"sfsUID"`                     /*  弹性文件功能系统唯一 ID  */
	VpcID                      string `json:"vpcID"`                      /*  vpc底层id  */
	VpcFuid                    string `json:"vpcFuid"`                    /*  vpc fuid  */
	VpcName                    string `json:"vpcName"`                    /*  vpc名称  */
	VpcCidr                    string `json:"vpcCidr"`                    /*  vpc cidr  */
	PermissionGroupID          string `json:"permissionGroupID"`          /*  权限组底层id  */
	PermissionGroupFuid        string `json:"permissionGroupFuid"`        /*  权限组id  */
	PermissionGroupName        string `json:"permissionGroupName"`        /*  权限组名称  */
	PermissionGroupDescription string `json:"permissionGroupDescription"` /*  权限组描述  */
	PermissionGroupIsDefault   *bool  `json:"permissionGroupIsDefault"`   /*  是否为默认权限组  */
}
