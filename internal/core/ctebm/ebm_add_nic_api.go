package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbmAddNicApi
/* 添加物理机网卡
 */type EbmAddNicApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmAddNicApi(client *core.CtyunClient) *EbmAddNicApi {
	return &EbmAddNicApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebm/add-nic",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmAddNicApi) Do(ctx context.Context, credential core.Credential, req *EbmAddNicRequest) (*EbmAddNicResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbmAddNicResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmAddNicRequest struct {
	RegionID string `json:"regionID"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string `json:"azName"` /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */InstanceUUID string `json:"instanceUUID"` /*  实例UUID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9715&data=87&isNormal=1">根据订单号查询uuid</a>获取实例UUID
	 */SubnetUUID string `json:"subnetUUID"` /*  子网UUID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10098380">基本概念</a>来查找子网的相关定义<br /> <a href="https://www.ctyun.cn/document/10026755/10040797">查询子网列表</a><br /><a href="https://www.ctyun.cn/document/10026755/10040804">创建子网</a><br/>注：在多可用区类型资源池下，subnetID通常以“subnet-”开头；非多可用区类型资源池subnetID为uuid格式
	 */SecurityGroups string `json:"securityGroups"` /*  安全组id，以逗号,分隔，单次最多允许传入10个安全组id；<br />获取：<br /><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=6940&data=97&isNormal=1&vid=91">查询单台物理机</a><br /><br />获取：<br /><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=6941&data=97&isNormal=1&vid=91">批量查询物理机</a><br />
	 */Ipv4 *string `json:"ipv4"` /*  IPV4地址
	 */
}

type EbmAddNicResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmAddNicReturnObjResponse `json:"returnObj"` /*  返回参数，值为空
	 */
}

type EbmAddNicReturnObjResponse struct{}
