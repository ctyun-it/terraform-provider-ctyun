package oceanfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// OceanfsVpcChangePermissionApi
/* 根据文件系统ID、VPC及权限组ID ，更换权限组
 */type OceanfsVpcChangePermissionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOceanfsVpcChangePermissionApi(client *core.CtyunClient) *OceanfsVpcChangePermissionApi {
	return &OceanfsVpcChangePermissionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oceanfs/vpc-change-permission",
			ContentType:  "application/json",
		},
	}
}

func (a *OceanfsVpcChangePermissionApi) Do(ctx context.Context, credential core.Credential, req *OceanfsVpcChangePermissionRequest) (*OceanfsVpcChangePermissionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*OceanfsVpcChangePermissionRequest
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
	var resp OceanfsVpcChangePermissionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type OceanfsVpcChangePermissionRequest struct {
	PermissionGroupFuid string `json:"permissionGroupFuid,omitempty"` /*  权限组ID  */
	RegionID            string `json:"regionID,omitempty"`            /*  资源池 ID  */
	SfsUID              string `json:"sfsUID,omitempty"`              /*  海量文件功能系统唯一 ID  */
	VpcID               string `json:"vpcID,omitempty"`               /*  vpc ID  */
}

type OceanfsVpcChangePermissionResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}
