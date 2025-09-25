package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbmMetadataListApi
/* 查询物理机的元数据信息
 */type EbmMetadataListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmMetadataListApi(client *core.CtyunClient) *EbmMetadataListApi {
	return &EbmMetadataListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebm/metadata/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmMetadataListApi) Do(ctx context.Context, credential core.Credential, req *EbmMetadataListRequest) (*EbmMetadataListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("azName", req.AzName)
	ctReq.AddParam("instanceUUID", req.InstanceUUID)
	if req.MetadataKey != nil {
		ctReq.AddParam("metadataKey", *req.MetadataKey)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbmMetadataListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmMetadataListRequest struct {
	RegionID string /*  资源池ID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string /*  可用区名称，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">获取可用区信息</a>，查询结果中zoneList内返回存在可用区名称（即多可用区，本字段填写实际可用区名称），若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */InstanceUUID string /*  实例UUID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9715&data=87&isNormal=1">根据订单号查询uuid</a>获取实例UUID
	 */MetadataKey *string /*  元数据键值，如缺省则查询所有元数据信息
	 */
}

type EbmMetadataListResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  错误码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmMetadataListReturnObjResponse `json:"returnObj"` /*  返回参数，参考表returnObj
	 */Error *string `json:"error"` /*  错误码，为product.module.code三段式码，详见错误码说明
	 */
}

type EbmMetadataListReturnObjResponse struct {
	Metadata *EbmMetadataListReturnObjMetadataResponse `json:"metadata"` /*  元数据，未设置情况下元数据则返回{}
	 */
}

type EbmMetadataListReturnObjMetadataResponse struct{}
