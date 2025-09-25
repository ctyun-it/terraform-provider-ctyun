package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbmBatchDeleteInstancesApi
/* 批量删除物理机
 */type EbmBatchDeleteInstancesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmBatchDeleteInstancesApi(client *core.CtyunClient) *EbmBatchDeleteInstancesApi {
	return &EbmBatchDeleteInstancesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebm/batch-delete",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmBatchDeleteInstancesApi) Do(ctx context.Context, credential core.Credential, req *EbmBatchDeleteInstancesRequest) (*EbmBatchDeleteInstancesResponse, error) {
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
	var resp EbmBatchDeleteInstancesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmBatchDeleteInstancesRequest struct {
	RegionID string `json:"regionID"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string `json:"azName"` /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */InstanceUUIDList string `json:"instanceUUIDList"` /*  实例UUID 用,分隔，单次请求最多指定10个实例；您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9715&data=87&isNormal=1">根据订单号查询uuid</a>获取实例UUID
	 */ClientToken string `json:"clientToken"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一，使用同一个ClientToken值，其他请求参数相同时，则代表为同一个请求。保留时间为24小时
	 */
}

type EbmBatchDeleteInstancesResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmBatchDeleteInstancesReturnObjResponse `json:"returnObj"` /*  返回参数，参考表returnObj
	 */
}

type EbmBatchDeleteInstancesReturnObjResponse struct {
	RegionID *string `json:"regionID"` /*  资源池ID
	 */MasterOrderID *string `json:"masterOrderID"` /*  订单ID
	 */MasterOrderNO *string `json:"masterOrderNO"` /*  订单号
	 */
}
