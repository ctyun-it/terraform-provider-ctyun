package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowPortApi
/* 查询网卡信息
 */type CtvpcShowPortApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowPortApi(client *core.CtyunClient) *CtvpcShowPortApi {
	return &CtvpcShowPortApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ports/show",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowPortApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowPortRequest) (*CtvpcShowPortResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("networkInterfaceID", req.NetworkInterfaceID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowPortResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowPortRequest struct {
	RegionID           string /*  资源池ID  */
	NetworkInterfaceID string /*  虚拟网卡id  */
}

type CtvpcShowPortResponse struct {
	StatusCode  int32                           `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                         `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                         `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                         `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowPortReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                         `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowPortReturnObjResponse struct {
	NetworkInterfaceName *string                                      `json:"networkInterfaceName,omitempty"` /*  虚拟网名称  */
	NetworkInterfaceID   *string                                      `json:"networkInterfaceID,omitempty"`   /*  虚拟网id  */
	VpcID                *string                                      `json:"vpcID,omitempty"`                /*  所属vpc  */
	SubnetID             *string                                      `json:"subnetID,omitempty"`             /*  所属子网id  */
	Role                 int32                                        `json:"role"`                           /*  网卡类型: 0 主网卡， 1 弹性网卡  */
	MacAddress           *string                                      `json:"macAddress,omitempty"`           /*  mac地址  */
	PrimaryPrivateIp     *string                                      `json:"primaryPrivateIp,omitempty"`     /*  主ip  */
	Ipv6Addresses        []*string                                    `json:"ipv6Addresses"`                  /*  ipv6地址  */
	InstanceID           *string                                      `json:"instanceID,omitempty"`           /*  关联的设备id  */
	InstanceType         *string                                      `json:"instanceType,omitempty"`         /*  设备类型 VM, BM, Other  */
	Description          *string                                      `json:"description,omitempty"`          /*  描述  */
	SecurityGroupIds     []*string                                    `json:"securityGroupIds"`               /*  安全组ID列表  */
	SecondaryPrivateIps  []*string                                    `json:"secondaryPrivateIps"`            /*  辅助私网IP  */
	AdminStatus          *string                                      `json:"adminStatus,omitempty"`          /*  是否启用DOWN, UP  */
	AssociatedEip        *CtvpcShowPortReturnObjAssociatedEipResponse `json:"associatedEip"`                  /*  关联的eip信息  */
}

type CtvpcShowPortReturnObjAssociatedEipResponse struct {
	Id   *string `json:"id,omitempty"`   /*  eip id  */
	Name *string `json:"name,omitempty"` /*  eip名称  */
}
