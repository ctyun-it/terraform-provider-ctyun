package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbmInstanceInterfaceSecurityGroupListApi
/* 查询物理机指定网卡的安全组信息
 */type EbmInstanceInterfaceSecurityGroupListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmInstanceInterfaceSecurityGroupListApi(client *core.CtyunClient) *EbmInstanceInterfaceSecurityGroupListApi {
	return &EbmInstanceInterfaceSecurityGroupListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebm/instance_interface_security_group_list",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmInstanceInterfaceSecurityGroupListApi) Do(ctx context.Context, credential core.Credential, req *EbmInstanceInterfaceSecurityGroupListRequest) (*EbmInstanceInterfaceSecurityGroupListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("azName", req.AzName)
	ctReq.AddParam("instanceUUID", req.InstanceUUID)
	ctReq.AddParam("interfaceUUID", req.InterfaceUUID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbmInstanceInterfaceSecurityGroupListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmInstanceInterfaceSecurityGroupListRequest struct {
	RegionID string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */InstanceUUID string /*  实例UUID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9715&data=87&isNormal=1">根据订单号查询uuid</a>获取实例UUID
	 */InterfaceUUID string /*  实例的网卡UUID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=6940&data=97">查询单台物理机</a>和<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=6941&data=97">批量查询物理机</a>获取网卡UUID
	 */
}

type EbmInstanceInterfaceSecurityGroupListResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmInstanceInterfaceSecurityGroupListReturnObjResponse `json:"returnObj"` /*  返回参数，参考表returnObj
	 */
}

type EbmInstanceInterfaceSecurityGroupListReturnObjResponse struct {
	SecurityGroups []*EbmInstanceInterfaceSecurityGroupListReturnObjSecurityGroupsResponse `json:"securityGroups"` /*  安全组信息列表
	 */
}

type EbmInstanceInterfaceSecurityGroupListReturnObjSecurityGroupsResponse struct {
	SecurityGroupID *string `json:"securityGroupID"` /*  安全组UUID
	 */SecurityGroupName *string `json:"securityGroupName"` /*  安全组名称
	 */
}
