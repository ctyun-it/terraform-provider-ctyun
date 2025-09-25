package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EbmDeviceStockListApi
/* 查询物理机库存
 */type EbmDeviceStockListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmDeviceStockListApi(client *core.CtyunClient) *EbmDeviceStockListApi {
	return &EbmDeviceStockListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebm/device-stock-list",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmDeviceStockListApi) Do(ctx context.Context, credential core.Credential, req *EbmDeviceStockListRequest) (*EbmDeviceStockListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("azName", req.AzName)
	if req.DeviceType != nil {
		ctReq.AddParam("deviceType", *req.DeviceType)
	}
	if req.Count != 0 {
		ctReq.AddParam("count", strconv.FormatInt(int64(req.Count), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbmDeviceStockListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmDeviceStockListRequest struct {
	RegionID string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */DeviceType *string /*  物理机套餐类型<br/><a href="https://www.ctyun.cn/document/10027724/10040107">查询资源池内物理机套餐</a><br /><a href="https://www.ctyun.cn/document/10027724/10040124">查询指定物理机的套餐信息</a>
	 */Count int32 /*  所需库存数，必须为大于0的正整数
	 */
}

type EbmDeviceStockListResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmDeviceStockListReturnObjResponse `json:"returnObj"` /*  返回参数，参考表returnObj
	 */
}

type EbmDeviceStockListReturnObjResponse struct {
	TotalCount int32 `json:"totalCount"` /*  总记录数
	 */Results []*EbmDeviceStockListReturnObjResultsResponse `json:"results"` /*  分页明细，元素类型是results，定义请参考表results
	 */
}

type EbmDeviceStockListReturnObjResultsResponse struct {
	Available int32 `json:"available"` /*  可用物理机数量
	 */Success *bool `json:"success"` /*  是否库存足够
	 */Stocks []*EbmDeviceStockListReturnObjResultsStocksResponse `json:"stocks"` /*  套餐库存详情
	 */
}

type EbmDeviceStockListReturnObjResultsStocksResponse struct {
	DeviceType *string `json:"deviceType"` /*  套餐名称
	 */Available int32 `json:"available"` /*  可用物理机数量
	 */Success *bool `json:"success"` /*  是否库存足够
	 */
}
