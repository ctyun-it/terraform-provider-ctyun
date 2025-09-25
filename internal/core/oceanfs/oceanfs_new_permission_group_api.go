package oceanfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// OceanfsNewPermissionGroupApi
/* 创建权限组
 */type OceanfsNewPermissionGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOceanfsNewPermissionGroupApi(client *core.CtyunClient) *OceanfsNewPermissionGroupApi {
	return &OceanfsNewPermissionGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oceanfs/permission-group/new-permission-group",
			ContentType:  "application/json",
		},
	}
}

func (a *OceanfsNewPermissionGroupApi) Do(ctx context.Context, credential core.Credential, req *OceanfsNewPermissionGroupRequest) (*OceanfsNewPermissionGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*OceanfsNewPermissionGroupRequest
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
	var resp OceanfsNewPermissionGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type OceanfsNewPermissionGroupRequest struct {
	RegionID                   string `json:"regionID,omitempty"`                   /*  资源池 ID  */
	PermissionGroupName        string `json:"permissionGroupName,omitempty"`        /*  权限组名字。名字不能重复  */
	NetworkType                string `json:"networkType,omitempty"`                /*  权限组网络类型。有效值范围：private_network  */
	PermissionGroupDescription string `json:"permissionGroupDescription,omitempty"` /*  描述信息  */
}

type OceanfsNewPermissionGroupResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}
