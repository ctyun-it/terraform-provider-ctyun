package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbmRaidTypeListApi
/* 查询物理机本地系统盘或数据盘可选择的raid类型
 */type EbmRaidTypeListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmRaidTypeListApi(client *core.CtyunClient) *EbmRaidTypeListApi {
	return &EbmRaidTypeListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebm/raid-type-list",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmRaidTypeListApi) Do(ctx context.Context, credential core.Credential, req *EbmRaidTypeListRequest) (*EbmRaidTypeListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("azName", req.AzName)
	ctReq.AddParam("deviceType", req.DeviceType)
	ctReq.AddParam("volumeType", req.VolumeType)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbmRaidTypeListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmRaidTypeListRequest struct {
	RegionID string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */DeviceType string /*  物理机套餐类型<br/><a href="https://www.ctyun.cn/document/10027724/10040107">查询资源池内物理机套餐</a><br /><a href="https://www.ctyun.cn/document/10027724/10040124">查询指定物理机的套餐信息</a>
	 */VolumeType string /*  本地磁盘类型，取值范围：system（系统盘），data（数据盘）
	 */
}

type EbmRaidTypeListResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmRaidTypeListReturnObjResponse `json:"returnObj"` /*  返回参数，参考表returnObj
	 */
}

type EbmRaidTypeListReturnObjResponse struct {
	TotalCount int32 `json:"totalCount"` /*  总记录数
	 */Results []*EbmRaidTypeListReturnObjResultsResponse `json:"results"` /*  分页明细，元素类型是results，定义请参考表results
	 */
}

type EbmRaidTypeListReturnObjResultsResponse struct {
	DeviceType *string `json:"deviceType"` /*  套餐类型
	 */VolumeType *string `json:"volumeType"` /*  磁盘类型
	 */Uuid *string `json:"uuid"` /*  raid uuid
	 */NameEn *string `json:"nameEn"` /*  raid英文名称
	 */NameZh *string `json:"nameZh"` /*  raid中文名称
	 */VolumeDetail *string `json:"volumeDetail"` /*  对应套餐磁盘描述
	 */DescriptionEn *string `json:"descriptionEn"` /*  raid英文介绍
	 */DescriptionZh *string `json:"descriptionZh"` /*  raid中文介绍
	 */
}
