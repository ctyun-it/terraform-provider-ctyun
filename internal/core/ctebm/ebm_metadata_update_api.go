package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbmMetadataUpdateApi
/* 物理机更新元数据信息
 */type EbmMetadataUpdateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmMetadataUpdateApi(client *core.CtyunClient) *EbmMetadataUpdateApi {
	return &EbmMetadataUpdateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebm/metadata/update",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmMetadataUpdateApi) Do(ctx context.Context, credential core.Credential, req *EbmMetadataUpdateRequest) (*EbmMetadataUpdateResponse, error) {
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
	var resp EbmMetadataUpdateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmMetadataUpdateRequest struct {
	RegionID string `json:"regionID"` /*  资源池ID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string `json:"azName"` /*  可用区名称，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">获取可用区信息</a>，查询结果中zoneList内返回存在可用区名称（即多可用区，本字段填写实际可用区名称），若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */InstanceUUID string `json:"instanceUUID"` /*  实例UUID
	 */MetadataKey string `json:"metadataKey"` /*  元数据键值
	 */MetadataValue string `json:"metadataValue"` /*  元数据值
	 */
}

type EbmMetadataUpdateResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  错误码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmMetadataUpdateReturnObjResponse `json:"returnObj"` /*  返回参数，值为空
	 */Error *string `json:"error"` /*  错误码，为product.module.code三段式码，详见错误码说明
	 */
}

type EbmMetadataUpdateReturnObjResponse struct{}
